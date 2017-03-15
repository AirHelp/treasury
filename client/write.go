package client

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AirHelp/treasury/aws"
	"github.com/AirHelp/treasury/utils"
)

const (
	equalSecrets = "Secrets are equal"
	noSuchKey    = "NoSuchKey"
)

var skipErrors = []string{equalSecrets, noSuchKey}

// Write secret to Treasure
func (c *Client) Write(key, secret string, force bool) error {
	environment, application, err := utils.FindEnvironmentApplicationName(key)
	if err != nil {
		return err
	}

	if err = c.skipValue(key, secret, force); err != nil {
		for _, errorMessage := range skipErrors {
			if strings.Contains(err.Error(), errorMessage) {
				fmt.Printf("Success! Skipped secret write:%s - %s\n", key, err.Error())
				return nil
			}
		}
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
	fmt.Println("Success! Data written to: ", key)
	return nil
}

func (c Client) skipValue(key, value string, force bool) error {
	if force {
		return nil
	}
	secretObject, err := c.Read(key)
	if err != nil {
		return err
	}
	if value == secretObject.Value {
		return errors.New(equalSecrets)
	}
	return nil
}
