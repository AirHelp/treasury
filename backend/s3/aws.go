package s3

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Client with AWS services
type S3ClientInterface interface {
	PutObject(context.Context, *s3.PutObjectInput, ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObject(context.Context, *s3.GetObjectInput, ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	ListObjects(context.Context, *s3.ListObjectsInput, ...func(*s3.Options)) (*s3.ListObjectsOutput, error)
	DeleteObject(context.Context, *s3.DeleteObjectInput, ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

type Client struct {
    S3Svc  S3ClientInterface
    bucket string
}

func New(region, bucket string) (*Client, error) {
	if bucket == "" {
		return nil, errors.New("S3 bucket name is missing")
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))

	if err != nil {
		return &Client{}, errors.Join(
			fmt.Errorf("unable to load SDK config with region %s", region),
			err,
		)
	}

	return &Client{
		S3Svc: s3.NewFromConfig(cfg),
	}, err
}
