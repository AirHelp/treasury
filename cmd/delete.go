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

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete KEY",
	Short: "Remove secrets from Treasury",
	Long:  `Remove secret with the given key from Treasury.`,
	RunE:  delete,
}

func init() {
	RootCmd.AddCommand(deleteCmd)
	deleteCmd.SuggestFor = []string{"remove"}
}

func delete(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Missing Key to delete.")
	}
	key := args[0]

	treasury, err := client.New(&client.Options{
		Region:       s3Region,
		S3BucketName: treasuryS3,
	})
	if err != nil {
		return err
	}

	err = treasury.Delete(key)
	if err != nil {
		return err
	}

	fmt.Printf("Key %s has been successfully deleted\n", key)
	return nil
}
