package client

import (
	"errors"

	"github.com/AirHelp/treasury/aws"
)

const (
	// version is the version of this client
	version = "0.0.1"
)

// Client is the client that performs all operations against a treasury backend
type Client struct {
	// s3 bucket name
	bucketName string

	// version of the Client
	version string

	// AWS shared client
	AwsClient *aws.Client
}

// Options for client
type Options struct {
	Version   string
	AwsClient *aws.Client
}

// NewClient initializes a new client for the given AWS account with S3 bucket
func NewClient(bucketName string, options *Options) (*Client, error) {
	if bucketName == "" {
		return nil, errors.New("S3 bucket name is missing")
	}

	if options.Version == "" {
		options.Version = version
	}

	// AWS connection
	var err error
	if options.AwsClient == nil {
		options.AwsClient, err = aws.New()
		if err != nil {
			return nil, err
		}
	}

	return &Client{
		bucketName: bucketName,
		version:    options.Version,
		AwsClient:  options.AwsClient,
	}, nil
}
