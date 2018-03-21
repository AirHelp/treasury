// Copyright © 2018 AirHelp Inc. devops@airhelp.com
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

// listCmd represents the read command
var listCmd = &cobra.Command{
	Use:   "list PATH",
	Short: "List the secrets set for a path",
	Long:  `List the secrets set for a path`,
	RunE:  list,
}

func init() {
	RootCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Missing key path to list")
	}
	path := args[0]

	treasury, err := client.New(&client.Options{
		Region:       s3Region,
		S3BucketName: treasuryS3,
	})
	if err != nil {
		return err
	}
	secrets, err := treasury.ReadGroup(path)
	if err != nil {
		return err
	}

	for _, secret := range secrets {
		fmt.Println(secret.Key)
	}
	return nil
}
