package s3

import (
	"bytes"
	"fmt"

	"github.com/AirHelp/treasury/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	// http://docs.aws.amazon.com/AmazonS3/latest/dev/acl-overview.html
	s3ACL = "private"
	// http://docs.aws.amazon.com/AmazonS3/latest/dev/UsingServerSideEncryption.html
	s3ServerSideEncryption = "aws:kms"
	// ApplicatonMetaKey is used as a Key for s3 object's metadata and tag
	ApplicatonMetaKey = "Application"
	// EnvironmentMetaKey is used as a Key for s3 object's metadata and tag
	EnvironmentMetaKey = "Environment"
)

// PutObject copy secret data on S3 bucket
// https://docs.aws.amazon.com/AmazonS3/latest/API/RESTObjectPUT.html
func (c *Client) PutObject(object *types.PutObjectInput) error {

	tags := fmt.Sprintf(
		"%s=%s&%s=%s",
		ApplicatonMetaKey,
		object.Application,
		EnvironmentMetaKey,
		object.Environment,
	)

	params := &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(object.Key),
		ACL:    aws.String(s3ACL),
		Body:   bytes.NewReader([]byte(object.Value)),
		Metadata: map[string]*string{
			ApplicatonMetaKey:  aws.String(object.Application),
			EnvironmentMetaKey: aws.String(object.Environment),
		},
		ServerSideEncryption: aws.String(s3ServerSideEncryption),
		SSEKMSKeyId:          aws.String("alias/" + object.Environment),
		Tagging:              aws.String(tags),
	}

	_, err := c.S3Svc.PutObject(params)

	return err
}

// GetObject reads secret from S3 bucket
// http://docs.aws.amazon.com/sdk-for-go/api/service/s3/#S3.GetObject
// https://docs.aws.amazon.com/goto/WebAPI/s3-2006-03-01/GetObject
func (c *Client) GetObject(object *types.GetObjectInput) (*types.GetObjectOutput, error) {

	params := &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(object.Key),
	}

	resp, err := c.S3Svc.GetObject(params)
	if err != nil {
		return nil, err
	}

	return &types.GetObjectOutput{Body: resp.Body}, nil
}

// GetObjects returns key value map for given pattern
func (c *Client) GetObjects(object *types.GetObjectsInput) (*types.GetObjectsOuput, error) {
	params := &s3.ListObjectsInput{
		Bucket: aws.String(c.bucket),
		Prefix: aws.String(object.Prefix),
	}

	resp, err := c.S3Svc.ListObjects(params)
	if err != nil {
		return nil, err
	}

	keyValuePairs := make(map[string]string, len(resp.Contents))
	for _, keyObject := range resp.Contents {
		key := *keyObject.Key
		s3Object, err := c.GetObject(&types.GetObjectInput{Key: key})
		if err != nil {
			return nil, err
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(s3Object.Body)
		keyValuePairs[key] = buf.String()
	}
	return &types.GetObjectsOuput{Secrets: keyValuePairs}, nil
}
