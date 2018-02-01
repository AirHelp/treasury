package s3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// Client with AWS services
type Client struct {
	sess   *session.Session
	S3Svc  s3iface.S3API
	bucket string
}

// New returns clients for AWS services
func New(region, bucket string) (*Client, error) {
	config := aws.Config{}
	if region != "" {
		config.Region = aws.String(region)
	}

	sess, err := session.NewSession(&config)
	if err != nil {
		return nil, fmt.Errorf("Failed to create AWS session. Error: %s", err)
	}

	s3Svc := s3.New(sess)

	return &Client{
		sess:   sess,
		S3Svc:  s3Svc,
		bucket: bucket,
	}, nil
}
