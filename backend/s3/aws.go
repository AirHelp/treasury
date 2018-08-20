package s3

import (
	"errors"
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
	if bucket == "" {
		return nil, errors.New("S3 bucket name is missing")
	}
	sessionOpts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}

	if region != "" {
		sessionOpts.Config = aws.Config{Region: aws.String(region)}
	}

	sess, err := session.NewSessionWithOptions(sessionOpts)
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
