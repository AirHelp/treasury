package ssm

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// Client with AWS services
type Client struct {
	svc ssmiface.SSMAPI
}

// New returns clients for AWS services
func New(region string) (*Client, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return nil, fmt.Errorf("Failed to create AWS session. Error: %s", err)
	}

	// Create a SSM client with additional configuration
	svc := ssm.New(sess)

	return &Client{
		svc: svc,
	}, nil
}
