package test

import (
	"bytes"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// MockS3Client fake S3API
type MockS3Client struct {
	s3iface.S3API
}

func (m *MockS3Client) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return &s3.PutObjectOutput{}, nil
}

func (m *MockS3Client) GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return &s3.GetObjectOutput{
		Body: ioutil.NopCloser(
			bytes.NewReader(
				[]byte("secret"))),
	}, nil
}
