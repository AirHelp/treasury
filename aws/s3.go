package aws

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// PutObjectInput structure for PutObject
type PutObjectInput struct {
	Bucket      string
	Key         string
	Value       string
	Application string
	Environment string
}

// GetObjectInput structure for GetObject
type GetObjectInput struct {
	Bucket  string
	Key     string
	Version string
}

const (
	// http://docs.aws.amazon.com/AmazonS3/latest/dev/acl-overview.html
	s3ACL = "private"
	// http://docs.aws.amazon.com/AmazonS3/latest/dev/UsingServerSideEncryption.html
	s3ServerSideEncryption = "AES256"
	// ApplicatonMetaKey is used as a Key for s3 object's metadata and tag
	ApplicatonMetaKey = "Application"
	// EnvironmentMetaKey is used as a Key for s3 object's metadata and tag
	EnvironmentMetaKey = "Environment"
)

// PutObject copy secret data on S3 bucket
// https://docs.aws.amazon.com/AmazonS3/latest/API/RESTObjectPUT.html
func (c *Client) PutObject(object *PutObjectInput) error {

	tags := fmt.Sprintf(
		"%s=%s&%s=%s",
		ApplicatonMetaKey,
		object.Application,
		EnvironmentMetaKey,
		object.Environment,
	)

	params := &s3.PutObjectInput{
		Bucket: aws.String(object.Bucket),
		Key:    aws.String(object.Key),
		ACL:    aws.String(s3ACL),
		Body:   bytes.NewReader([]byte(object.Value)),
		Metadata: map[string]*string{
			ApplicatonMetaKey:  aws.String(object.Application),
			EnvironmentMetaKey: aws.String(object.Environment),
		},
		ServerSideEncryption: aws.String(s3ServerSideEncryption),
		Tagging:              aws.String(tags),
	}

	_, err := c.S3Svc.PutObject(params)

	return err
}

// GetObject reads secret from S3 bucket
// http://docs.aws.amazon.com/sdk-for-go/api/service/s3/#S3.GetObject
// https://docs.aws.amazon.com/goto/WebAPI/s3-2006-03-01/GetObject
func (c *Client) GetObject(object *GetObjectInput) (*s3.GetObjectOutput, error) {

	params := &s3.GetObjectInput{
		Bucket: aws.String(object.Bucket),
		Key:    aws.String(object.Key),
	}

	resp, err := c.S3Svc.GetObject(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
