package client

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"

	"github.com/AirHelp/treasury/aws"
)

// Write secret to Treasure
func (c *Client) Write(key, secret, kmsKey string) (string, error) {
	// AWS connection
	awsClient, err := aws.New()
	if err != nil {
		return "", err
	}

	// get user name from AWS IAM credentials
	username, err := awsClient.GetUserName()
	if err != nil {
		return "", err
	}

	// convert plain text secret into encrypted blob
	kmsResponse, err := awsClient.Encrypt(kmsKey, secret)
	if err != nil {
		return "", err
	}

	body := Secret{
		Key:    key,
		Value:  base64.StdEncoding.EncodeToString(kmsResponse.CiphertextBlob),
		KmsARN: *kmsResponse.KeyId,
		Author: username,
	}

	// API rq
	req, err := c.NewRequest("POST", "secret", body)
	if err != nil {
		return "", err
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		// TO DO: return something more
		return "", errors.New(res.Status)
	}

	// unmarshal API response
	msg := WriteMessage{}
	err = json.NewDecoder(res.Body).Decode(&msg)
	if err != nil && err != io.EOF {
		// ignore EOF errors caused by empty response body
		return "", err
	}

	return msg.Message, nil
}
