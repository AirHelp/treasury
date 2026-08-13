package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/AirHelp/treasury/client"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
)

const (
	runCommandEnvFileArgument = "env-file"
	runCommandProfileArgument = "profile"
	defaultEnvFile            = ".env.treasury"

	// exit codes follow the shell convention
	commandNotFoundExitCode      = 127
	commandNotExecutableExitCode = 126
	signalExitCodeBase           = 128
)

// signals relayed to the command being run, the ones a process is expected to
// act on. The rest is left to the default behaviour.
var forwardedSignals = []os.Signal{
	syscall.SIGINT,
	syscall.SIGTERM,
	syscall.SIGQUIT,
	syscall.SIGHUP,
	syscall.SIGUSR1,
	syscall.SIGUSR2,
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [flags] COMMAND [ARGS...]",
	Short: "Runs a command with secrets loaded into environment variables",
	Long: `Run loads secrets described in an environment file and executes the given
command with them exported as environment variables. Secrets are kept in memory
only, they are never written to disk.

  treasury run bundle exec rake db:migrate
  treasury run --env-file .env.staging -- rails server

The environment file (.env.treasury by default) accepts:

  {{ export "development/webapp/" }}                 all secrets from the path,
                                                     named after the last path part
  AUTH_API_PASSWORD={{ read "development/auth/PASSWORD" }}  a single secret
  API_TOKEN=test                                     a plain value

The AWS profile is taken from AWS_PROFILE, use --profile to pick another one.`,
	RunE: run,
}

func init() {
	RootCmd.AddCommand(runCmd)
	// everything after the command name belongs to the command being run
	runCmd.Flags().SetInterspersed(false)
	runCmd.Flags().String(runCommandEnvFileArgument, defaultEnvFile, "path to the environment file with secrets")
	runCmd.Flags().String(runCommandProfileArgument, "", "AWS profile to use, defaults to AWS_PROFILE")
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("missing command to run")
	}
	// the arguments are fine, whatever fails next is not a usage problem
	cmd.SilenceUsage = true

	envFile, err := cmd.Flags().GetString(runCommandEnvFileArgument)
	if err != nil {
		return err
	}
	profile, err := cmd.Flags().GetString(runCommandProfileArgument)
	if err != nil {
		return err
	}

	treasury, err := newClientWithProfile(profile)
	if err != nil {
		return err
	}

	environment, err := treasury.EnvFile(envFile)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("environment file %s not found, create it or point to another one with --%s", envFile, runCommandEnvFileArgument)
		}
		return err
	}

	return execute(cmd.Context(), args, environment)
}

// execute runs the command with the given environment and leaves treasury with
// the same exit code the command ended with
func execute(ctx context.Context, args, environment []string) error {
	command := exec.CommandContext(ctx, args[0], args[1:]...) // #nosec G204
	command.Env = append(os.Environ(), environment...)
	command.Stdin, command.Stdout, command.Stderr = os.Stdin, os.Stdout, os.Stderr

	if err := command.Start(); err != nil {
		switch {
		case errors.Is(err, exec.ErrNotFound), errors.Is(err, fs.ErrNotExist):
			fmt.Fprintf(os.Stderr, "treasury: %s: command not found\n", args[0])
			os.Exit(commandNotFoundExitCode)
		case errors.Is(err, fs.ErrPermission):
			fmt.Fprintf(os.Stderr, "treasury: %s: cannot execute\n", args[0])
			os.Exit(commandNotExecutableExitCode)
		}
		return err
	}

	stopForwarding := forwardSignals(command.Process)
	err := command.Wait()
	stopForwarding()

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		os.Exit(exitCode(exitErr.ProcessState))
	}
	return err
}

// forwardSignals relays the signals treasury receives to the command it runs,
// so that a SIGTERM from a supervisor reaches the process that does the work.
// The command shares the process group with treasury, so signals coming from
// the terminal reach it on their own, a Ctrl-C simply arrives twice.
// The returned function stops the forwarding.
func forwardSignals(process *os.Process) func() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, forwardedSignals...)

	done := make(chan struct{})
	go func() {
		for {
			select {
			case receivedSignal := <-signals:
				// the command may be gone already, nothing to do about it
				_ = process.Signal(receivedSignal)
			case <-done:
				return
			}
		}
	}()

	return func() {
		signal.Stop(signals)
		close(done)
	}
}

// exitCode of a finished command, 128 + signal number when it was killed
func exitCode(state *os.ProcessState) int {
	if status, ok := state.Sys().(syscall.WaitStatus); ok && status.Signaled() {
		return signalExitCodeBase + int(status.Signal())
	}
	return state.ExitCode()
}

func newClientWithProfile(profile string) (*client.Client, error) {
	options := &client.Options{
		Region:       s3Region,
		S3BucketName: treasuryS3,
	}
	if profile != "" {
		awsConfig, err := config.LoadDefaultConfig(context.Background(),
			config.WithSharedConfigProfile(profile),
			config.WithRegion(s3Region),
		)
		if err != nil {
			return nil, err
		}
		options.AWSConfig = awsConfig
	}
	return client.New(options)
}
