package client

import (
	"fmt"

	"github.com/AirHelp/treasury/types"
	"github.com/AirHelp/treasury/utils"
)

// Read returns decrypted secret for given key
func (c *Client) Read(key string) (*Secret, error) {
	if err := utils.ValidateInputKey(key); err != nil {
		return nil, err
	}

	s3objectInput := &types.GetObjectInput{
		Key: key,
	}
	s3object, err := c.Backend.GetObject(s3objectInput)
	if err != nil {
		return nil, err
	}

	return &Secret{
		Key:   key,
		Value: s3object.Value,
	}, nil
}

// ReadValue returns secret as a string for given key.
func (c *Client) ReadValue(key string) (string, error) {
	secret, err := c.Read(key)
	if err != nil {
		return "", err
	}
	return secret.Value, nil
}

// ReadFromEnv returns value of given key in specified env.
func (c *Client) ReadFromEnv(env, key string) (string, error) {
	return c.ReadValue(fmt.Sprintf("%s/%s", env, key))
}

// ReadKeys returns the values of the given keys only, fetched in as few calls
// as the backend allows. Nothing outside the listed keys is read or decrypted.
func (c *Client) ReadKeys(keys []string) (map[string]string, error) {
	if len(keys) == 0 {
		return nil, nil
	}
	for _, key := range keys {
		if err := utils.ValidateInputKey(key); err != nil {
			return nil, err
		}
	}
	resp, err := c.Backend.GetObjects(&types.GetObjectsInput{Keys: keys})
	if err != nil {
		return nil, err
	}
	return resp.Secrets, nil
}

// ReadGroup returns list of secrets for given key prefix
func (c *Client) ReadGroup(keyPrefix string) ([]*Secret, error) {
	if err := utils.ValidateInputKeyPattern(keyPrefix); err != nil {
		return nil, err
	}
	params := &types.GetObjectsInput{
		Prefix: keyPrefix,
	}
	resp, err := c.Backend.GetObjects(params)
	if err != nil {
		return nil, err
	}

	secrets := make([]*Secret, 0, len(resp.Secrets))

	for key, value := range resp.Secrets {
		secret := &Secret{
			Key:   key,
			Value: value,
		}
		secrets = append(secrets, secret)
	}
	return secrets, nil
}
