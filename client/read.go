package client

import (
	"bytes"

	"github.com/AirHelp/treasury/aws"
	"github.com/AirHelp/treasury/utils"
)

// Read returns decrypted secret for given key
func (c *Client) Read(key string) (*Secret, error) {
	if err := utils.ValidateInputKey(key); err != nil {
		return nil, err
	}

	s3objectInput := &aws.GetObjectInput{
		Bucket: c.bucketName,
		Key:    key,
	}
	s3object, err := c.AwsClient.GetObject(s3objectInput)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(s3object.Body)

	secret := &Secret{
		Key:   key,
		Value: buf.String(),
	}
	return secret, nil
}
