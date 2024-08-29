package s3

import (
	"bytes"
	"context"
	"fmt"

	"github.com/AirHelp/treasury/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
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
		ACL:    s3Types.ObjectCannedACLPrivate,
		Body:   bytes.NewReader([]byte(object.Value)),
		Metadata: map[string]string{
			ApplicatonMetaKey:  object.Application,
			EnvironmentMetaKey: object.Environment,
		},
		ServerSideEncryption: s3Types.ServerSideEncryptionAwsKms,
		SSEKMSKeyId:          aws.String("alias/" + object.Environment),
		Tagging:              aws.String(tags),
	}

	_, err := c.S3Svc.PutObject(context.Background(), params)

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

	resp, err := c.S3Svc.GetObject(context.Background(), params)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return &types.GetObjectOutput{Value: buf.String()}, nil
}

// GetObjects returns key value map for given pattern
func (c *Client) GetObjects(object *types.GetObjectsInput) (*types.GetObjectsOuput, error) {
	params := &s3.ListObjectsInput{
		Bucket: aws.String(c.bucket),
		Prefix: aws.String(object.Prefix),
	}

	resp, err := c.S3Svc.ListObjects(context.Background(), params)
	if err != nil {
		return nil, err
	}

	keyValuePairs := make(map[string]string, len(resp.Contents))
	for _, keyObject := range resp.Contents {
		key := *keyObject.Key
		object, err := c.GetObject(&types.GetObjectInput{Key: key})
		if err != nil {
			return nil, err
		}
		keyValuePairs[key] = object.Value
	}
	return &types.GetObjectsOuput{Secrets: keyValuePairs}, nil
}
func (c *Client) DeleteObject(object *types.DeleteObjectInput) error {
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(object.Key),
	}
	_, err := c.S3Svc.DeleteObject(context.Background(), params)
	return err
}
