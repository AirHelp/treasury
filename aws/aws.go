package aws

import (
	"fmt"

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
func New(externalSess ...*session.Session) (*Client, error) {
	if len(externalSess) != 0 {
		return &Client{sess: externalSess[0]}, nil
	}
	sess, err := session.NewSession()
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
