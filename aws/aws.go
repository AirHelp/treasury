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

// New returns clients for AWS services
func New(s3Region string) (*Client, error) {
	var sess *session.Session
	var err error

	if s3Region != "" {
		sess, err = session.NewSession(&aws.Config{Region: aws.String(s3Region)})
	} else {
		sess, err = session.NewSession()
	}

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
