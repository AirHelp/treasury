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

const (
	// ExportString format of single export string
	ExportString = "export %s='%s'\n"
)

// exportCmd represents the read command
var exportCmd = &cobra.Command{
	Use:   "export ENVIRONMENT/APPLICATION/[KEY]",
	Short: "Returns command exporting found secretes",
	Long:  `Returns command exporting found secretes to environment variables.`,
	RunE:  export,
}

func init() {
	RootCmd.AddCommand(exportCmd)
}

func export(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Missing Key pattern to export.")
	}
	key := args[0]

	treasury, err := client.New(&client.Options{
		Region:       s3Region,
		S3BucketName: treasuryS3,
	})
	if err != nil {
		return err
	}
	exportCommand, err := treasury.Export(key, ExportString, map[string]string{})
	if err != nil {
		return err
	}
	fmt.Println(exportCommand)
	return nil
}
