package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// Client with AWS services
type Client struct {
	sess  *session.Session
	S3Svc s3iface.S3API
}

// Options for AWS services
type Options struct {
	Region string
}

// New returns clients for AWS services
func New(options Options) (*Client, error) {

	config := aws.Config{}

	if options.Region != "" {
		config.Region = aws.String(options.Region)
	}

	sess, err := session.NewSession(&config)

	if err != nil {
		fmt.Println("Failed to create AWS session,", err)
		return nil, err
	}
	s3Svc := s3.New(sess)
	return &Client{
		sess:  sess,
		S3Svc: s3Svc,
	}, nil
}
