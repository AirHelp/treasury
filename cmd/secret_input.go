package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const secretPrompt = "Please paste your secret: "

var errEmptySecret = errors.New("empty secret, nothing was written")

// readSecret gets the secret value without exposing it in the command line.
// On a terminal it prompts for the secret with echo disabled, otherwise it
// reads the secret from the piped stdin, which keeps scripts and CI working.
func readSecret(cmd *cobra.Command) (string, error) {
	in := cmd.InOrStdin()

	//nolint:gosec // G115: a file descriptor always fits in an int
	if file, ok := in.(*os.File); ok && term.IsTerminal(int(file.Fd())) {
		return readSecretFromTerminal(cmd, file)
	}

	data, err := io.ReadAll(in)
	if err != nil {
		return "", err
	}

	return validateSecret(trimEOL(string(data)))
}

func readSecretFromTerminal(cmd *cobra.Command, file *os.File) (string, error) {
	out := cmd.OutOrStdout()

	_, _ = fmt.Fprint(out, secretPrompt)
	//nolint:gosec // G115: a file descriptor always fits in an int
	secret, err := term.ReadPassword(int(file.Fd()))
	_, _ = fmt.Fprintln(out)
	if err != nil {
		return "", err
	}

	return validateSecret(string(secret))
}

// trimEOL removes a single trailing line ending, the one shells and editors add
// to piped input. Any other whitespace is a part of the secret.
func trimEOL(secret string) string {
	return strings.TrimSuffix(strings.TrimSuffix(secret, "\n"), "\r")
}

func validateSecret(secret string) (string, error) {
	if secret == "" {
		return "", errEmptySecret
	}
	return secret, nil
}
