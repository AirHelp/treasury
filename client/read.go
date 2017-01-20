package client

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"

	"github.com/AirHelp/treasury/aws"
)

// Read returns decrypted secret for given key
func (c *Client) Read(key string) (*Secret, error) {
	var data Secret

	// AWS connection
	awsClient, err := aws.New()
	if err != nil {
		return nil, err
	}

	// API rq
	req, err := c.NewRequest("GET", "secret", nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Add("key", key)
	req.URL.RawQuery = query.Encode()

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}

	// unmarshal response
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil && err != io.EOF {
		return nil, err
	}

	// decrypt secret
	blobSecret, err := base64.StdEncoding.DecodeString(data.Value)
	if err != nil {
		return nil, err
	}
	plainTextSecret, err := awsClient.Decrypt(blobSecret)
	if err != nil {
		return nil, err
	}
	data.Value = plainTextSecret

	return &data, err
}
