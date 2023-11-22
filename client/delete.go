package client

import (
	"github.com/AirHelp/treasury/types"
	"github.com/AirHelp/treasury/utils"
)

// Delete removeds specified secret for given key
func (c *Client) Delete(key string) error {
	if err := utils.ValidateInputKey(key); err != nil {
		return err
	}

	secretObject := &types.DeleteObjectInput{
		Key: key,
	}

	if err := c.Backend.DeleteObject(secretObject); err != nil {
		return err
	}

	return nil
}
