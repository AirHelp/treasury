package s3

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
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
	s3iface.S3API
}

func (m *MockS3Client) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	if _, ok := KeyValueMap[*input.Key]; !ok {
		return nil, errors.New(fmt.Sprintf("Missing key:%s in KeyValue map", *input.Key))
	}
	return &s3.PutObjectOutput{}, nil
}

func (m *MockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return &s3.GetObjectOutput{
		Body: ioutil.NopCloser(
			bytes.NewReader(
				[]byte(KeyValueMap[*input.Key]))),
	}, nil
}

func (m *MockS3Client) ListObjects(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	var contents []*s3.Object
	for key := range KeyValueMap {
		if strings.Contains(key, *input.Prefix) {
			keyToAdd := key
			contents = append(contents, &s3.Object{Key: &keyToAdd})
		}
	}
	return &s3.ListObjectsOutput{
		Contents: contents,
	}, nil
}
