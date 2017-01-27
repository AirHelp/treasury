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

// GetObjectInput structure for PutObject
type GetObjectInput struct {
	Bucket  string
	Key     string
	Version string
}

// PutObject copy secret data on S3 bucket
// https://docs.aws.amazon.com/AmazonS3/latest/API/RESTObjectPUT.html
func (c *Client) PutObject(obcject *PutObjectInput) error {

	tags := fmt.Sprintf(
		"Application=%s&Environment=%s",
		obcject.Application,
		obcject.Environment,
	)

	params := &s3.PutObjectInput{
		Bucket: aws.String(obcject.Bucket),
		Key:    aws.String(obcject.Key),
		ACL:    aws.String("private"),
		Body:   bytes.NewReader([]byte(obcject.Value)),
		Metadata: map[string]*string{
			"Application": aws.String(obcject.Application),
			"Environment": aws.String(obcject.Environment),
		},
		ServerSideEncryption: aws.String("AES256"),
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
		// IfMatch:           aws.String("IfMatch"),
		// IfModifiedSince:   aws.Time(time.Now()),
		// IfUnmodifiedSince: aws.Time(time.Now()),
		// VersionId: aws.String(object.Version),
	}

	resp, err := c.S3Svc.GetObject(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
