package client

import (
	"github.com/AirHelp/treasury/utils"
)

// Import imports secrets from file into treasury store, if value does not change it is not overridden
func (c *Client) Import(prefix, secretsFilePath string, force bool) error {
	secrets, err := utils.ReadSecrets(secretsFilePath)
	if err != nil {
		return err
	}
	for key, value := range secrets {
		if err = c.Write(prefix+key, value, force); err != nil {
			return err
		}
	}
	return nil
}
