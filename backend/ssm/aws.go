package ssm

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// Client with AWS services
type Client struct {
	svc *ssm.Client
}

// New returns clients for AWS services
func New(region string, awsConfig aws.Config) (*Client, error) {
	// Load the default configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS configuration. Error: %s", err)
	}

	// Create a SSM client with additional configuration
	svc := ssm.NewFromConfig(cfg)

	return &Client{
		svc: svc,
	}, nil
}
