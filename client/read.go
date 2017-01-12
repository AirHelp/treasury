package client

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/AirHelp/treasury/aws"
)

// Read returns decrypted secret for given key
func (c *Client) Read(key string) (*Secret, error) {
	var data Secret

	// AWS connection
	awsClient, err := aws.New()
	if err != nil {
		return &Secret{}, err
	}

	// API rq
	context := fmt.Sprintf("secret?key=%s", key)
	req, err := c.NewRequest("GET", context, nil)
	if err != nil {
		return &Secret{}, err
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return &Secret{}, err
	}
	if res.StatusCode != 200 {
		return &Secret{}, errors.New(res.Status)
	}

	// unmarshal response
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil && err != io.EOF {
		return &Secret{}, err
	}

	// decrypt secret
	blobSecret, err := base64.StdEncoding.DecodeString(data.Value)
	if err != nil {
		return &Secret{}, err
	}
	plainTextSecret, err := awsClient.Decrypt(blobSecret)
	if err != nil {
		return &Secret{}, err
	}
	data.Value = plainTextSecret

	return &data, err
}
