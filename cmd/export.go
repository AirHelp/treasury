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

// exportCmd represents the read command
var exportCmd = &cobra.Command{
	Use:   "export ENVIRONMENT/APPLICATION/[KEY]",
	Short: "Return secrets found by given pattern",
	Long:  `Exports pairs KEY=VALUE from Treasury.`,
	RunE:  export,
}

func init() {
	RootCmd.AddCommand(exportCmd)
	readCmd.SuggestFor = []string{"get"}
}

func export(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Missing Key pattern to export.")
	}
	key := args[0]

	treasury, err := client.New(treasuryS3, &client.Options{Region: s3Region})
	if err != nil {
		return err
	}
	exportCommand, err := treasury.Export(key)
	if err != nil {
		return err
	}
	fmt.Println(exportCommand)
	return nil
}
