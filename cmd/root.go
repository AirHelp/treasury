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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const treasuryS3EnvVariable = "TREASURY_S3"

var (
	// TO DO: CLI config
	// if we decided to use config file we should use https://github.com/spf13/viper
	treasuryS3 string
	s3Region   string
	addToArray []string
)

var RootCmd = &cobra.Command{
	Use:   "treasury",
	Short: "A tool for managing secrets.",
	Long: `
Treasury is a very simple and easy to use tool for managing secrets.
`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVarP(&s3Region, "region", "r", "", "s3 region")
	RootCmd.PersistentFlags().StringArrayVar(&addToArray, "addto", []string{}, "variable suffix, e.g: --addto \"DATABASE_URL:?pool=10\"")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	treasuryS3 = os.Getenv(treasuryS3EnvVariable)
}
