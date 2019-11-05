package client

import (
	"github.com/AirHelp/treasury/types"
	"github.com/AirHelp/treasury/utils"
	"github.com/aws/aws-sdk-go/aws/awserr"

	"bytes"
	"compress/gzip"
	b64 "encoding/base64"
	"io/ioutil"
)

const (
	noSuchKey       = "NoSuchKey"
	noSuchParameter = "ParameterNotFound"
)

// Write secret to Treasure
func (c *Client) Write(key, secret string, force bool) error {
	environment, application, err := utils.FindEnvironmentApplicationName(key)
	if err != nil {
		return err
	}

	if !force {
		secretObject, err := c.Read(key)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				// in this case 404 is ok for us
				// so we'd proceed if 404 occurs
				if aerr.Code() != noSuchKey && aerr.Code() != noSuchParameter {
					return err
				}
			} else {
				return err
			}
		} else if secret == secretObject.Value {
			return nil
		}
	}

	body := &types.PutObjectInput{
		Key:         key,
		Value:       secret,
		Application: application,
		Environment: environment,
	}

	err = c.Backend.PutObject(body)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) WriteFile(key, file string, force bool) error {

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	var gzipped bytes.Buffer
	gz := gzip.NewWriter(&gzipped)
	if _, err := gz.Write(data); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}

	return c.Write(key, b64.StdEncoding.EncodeToString(gzipped.Bytes()), force)
}
