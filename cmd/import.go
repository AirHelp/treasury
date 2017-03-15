package cmd

import (
	"errors"
	"fmt"

	"github.com/AirHelp/treasury/client"
	"github.com/spf13/cobra"
)

const (
	success = "Import successful"
)

var importCmd = &cobra.Command{
	Use:   "import ENVIRONMENT/APPLICATION/ secrets/file/path",
	Short: "command importing secrets form file",
	Long:  `command importing secrets to properties file`,
	RunE:  importFunc,
}

func init() {
	RootCmd.AddCommand(importCmd)
	importCmd.PersistentFlags().Bool("force", false, "Force overwrite secret value")
}

func importFunc(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("Missing environment/application/ and secrets/file/path")
	}
	keyPrefix := args[0]
	secretsFilePath := args[1]
	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return err
	}

	treasury, err := client.New(treasuryS3, &client.Options{Region: s3Region})
	if err != nil {
		return err
	}
	if err := treasury.Import(keyPrefix, secretsFilePath, force); err != nil {
		return err
	}
	fmt.Println(success)
	return nil
}
