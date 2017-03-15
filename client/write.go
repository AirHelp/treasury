package client

import (
	"strings"

	"github.com/AirHelp/treasury/aws"
	"github.com/AirHelp/treasury/utils"
)

const (
	noSuchKey = "NoSuchKey"
)

// Write secret to Treasure
func (c *Client) Write(key, secret string, force bool) error {
	environment, application, err := utils.FindEnvironmentApplicationName(key)
	if err != nil {
		return err
	}

	if !force {
		secretObject, err := c.Read(key)
		if err != nil && !strings.Contains(err.Error(), noSuchKey) {
			return err
		} else if secret == secretObject.Value {
			return nil
		}
	}

	body := &aws.PutObjectInput{
		Bucket:      c.bucketName,
		Key:         key,
		Value:       secret,
		Application: application,
		Environment: environment,
	}

	err = c.AwsClient.PutObject(body)
	if err != nil {
		return err
	}
	return nil
}
