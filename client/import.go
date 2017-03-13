package client

import (
	"github.com/AirHelp/treasury/utils"
)

// Import imports secrets from file into s3, if value does not change it is not overridden
func (c *Client) Import(prefix, secretsFilePath string) error {
	secrets, err := utils.ReadSecrets(secretsFilePath)
	if err != nil {
		return err
	}
	for key, value := range secrets {
		if err = c.Write(prefix+key, value); err != nil {
			return err
		}
	}
	return nil
}
