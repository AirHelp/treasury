package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	// "github.com/aws/aws-sdk-go-v2/service/s3/s3iface"
)

const (
	Key1 = "test/webapp/cocpit_api_pass"
	Key2 = "test/webapp/user_api_pass"
	Key3 = "test/cockpit/user_api_pass"
)

var KeyValueMap = map[string]string{
	Key1: "as9@#$%^&*(/2hdiwnf",
	Key2: "as9@#$&*(/2saddsahdiwnf",
	Key3: "#$&*(/2saddsah&as",
}

// MockS3Client fake S3API
type MockS3Client struct {
}

func (m *MockS3Client) PutObject(ctx context.Context, input *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	if _, ok := KeyValueMap[*input.Key]; !ok {
		return nil, errors.New(fmt.Sprintf("Missing key:%s in KeyValue map", *input.Key))
	}
	return &s3.PutObjectOutput{}, nil
}

func (m *MockS3Client) GetObject(ctx context.Context, input *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return &s3.GetObjectOutput{
		Body: ioutil.NopCloser(
			bytes.NewReader(
				[]byte(KeyValueMap[*input.Key]))),
	}, nil
}

func (m *MockS3Client) ListObjects(ctx context.Context, input *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
	var contents []types.Object
	for key := range KeyValueMap {
		if strings.Contains(key, *input.Prefix) {
			keyToAdd := key
			contents = append(contents, types.Object{Key: &keyToAdd})
		}
	}
	return &s3.ListObjectsOutput{
		Contents: contents,
	}, nil
}

func (m *MockS3Client) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	return &s3.DeleteObjectOutput{}, nil
}
