package client

import (
	"github.com/AirHelp/treasury/backend"
	"github.com/aws/aws-sdk-go-v2/aws"
)

const (
	// version is the version of this client
	version = "0.0.1"
)

// Client is the client that performs all operations against a treasury backend
type Client struct {
	// version of the Client
	version string

	// Backend interface
	Backend backend.BackendAPI
}

// Options for client
type Options struct {
	Version string
	// backend region where we keep secrets
	Region       string
	S3BucketName string
	Backend      backend.BackendAPI
	AWSConfig    aws.Config
}

// New initializes a new client for the given AWS account with S3 bucket
func New(options *Options) (*Client, error) {
	if options.Version == "" {
		options.Version = version
	}

	// backend connection
	var err error
	if options.Backend == nil {
		options.Backend, err = backend.New(backend.Options{
			Region:       options.Region,
			S3BucketName: options.S3BucketName,
			AWSConfig:    options.AWSConfig,
		})
		if err != nil {
			return nil, err
		}
	}

	return &Client{
		version: options.Version,
		Backend: options.Backend,
	}, nil
}
