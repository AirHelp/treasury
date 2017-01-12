package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
)

// Client for AWS shared session
type Client struct {
	sess *session.Session
}

// New returns aws session
func New() (*Client, error) {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("Failed to create AWS session,", err)
		return &Client{}, err
	}
	return &Client{sess: sess}, nil
}
