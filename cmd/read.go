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

	"github.com/AirHelp/treasury/client"
	"github.com/spf13/cobra"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read KEY",
	Short: "Read secrets from Treasury",
	Long:  `Reads secret for the given key from Treasury.`,
	RunE:  read,
}

func init() {
	RootCmd.AddCommand(readCmd)
	readCmd.SuggestFor = []string{"get"}
}

func read(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Missing Key to read.")
	}
	key := args[0]

	treasury, err := client.New(treasuryS3, s3Region, &client.Options{})
	if err != nil {
		return err
	}
	secret, err := treasury.Read(key)
	if err != nil {
		return err
	}

	fmt.Println(secret.Value)
	return nil
}
