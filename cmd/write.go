// Copyright © 2016 AirHelp Inc. devops@airhelp.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/AirHelp/treasury/client"
	"github.com/spf13/cobra"
)

// writeCmd represents the write command
var writeCmd = &cobra.Command{
	Use:   "write ENVIRONMENT/APPLICATION/KEY or write user/USER.NAME/KEY",
	Short: "Write secrets into Treasury",
	Long: `Write sends data into Treasury at the given key (path).

The secret value is never given as a command line argument, so it does not end
up in the shell history nor in the process list. When run in a terminal treasury
asks for the secret and hides what is typed, otherwise the secret is read from
the standard input:

  treasury write development/webapp/cockpit_api_pass
  echo "${SECRET}" | treasury write development/webapp/cockpit_api_pass

The trailing newline added by echo is stripped. Pipe a variable or a command,
never the secret itself - that would leak into the shell history again.

With --file the second argument is a path to a file with the content to store.`,
	RunE: write,
}

func init() {
	RootCmd.AddCommand(writeCmd)
	writeCmd.SuggestFor = []string{"put"}
	writeCmd.PersistentFlags().Bool("force", false, "Force overwrite secret value")
	writeCmd.PersistentFlags().Bool("file", false, "Save file content into Treasury")
}

func write(cmd *cobra.Command, args []string) error {
	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return err
	}

	file, err := cmd.Flags().GetBool("file")
	if err != nil {
		return err
	}

	if file {
		return writeFile(cmd, args, force)
	}

	if len(args) == 0 {
		return errors.New("missing key to write")
	}
	if len(args) > 1 {
		return errors.New("too many arguments, the secret is not passed as an argument anymore - treasury asks for it or reads it from the standard input")
	}
	key := args[0]

	secret, err := readSecret(cmd)
	if err != nil {
		return err
	}

	treasury, err := newClient()
	if err != nil {
		return err
	}

	if err := treasury.Write(key, secret, force); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Success! Data written to: %s (%d characters)\n", key, utf8.RuneCountInString(secret))
	return nil
}

func writeFile(cmd *cobra.Command, args []string, force bool) error {
	if len(args) != 2 {
		return errors.New("missing key and file path to write")
	}
	key, filePath := args[0], args[1]

	treasury, err := newClient()
	if err != nil {
		return err
	}

	if err := treasury.WriteFile(key, filePath, force); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Success! Data written to: %s\n", key)
	return nil
}

func newClient() (*client.Client, error) {
	return client.New(&client.Options{
		Region:       s3Region,
		S3BucketName: treasuryS3,
	})
}
