package client

import (
	"github.com/AirHelp/treasury/aws"
	"github.com/AirHelp/treasury/utils"
)

// Write secret to Treasure
func (c *Client) Write(key, secret string) error {
	environment, application, err := utils.FindEnvironmentApplicationName(key)
	if err != nil {
		return err
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
