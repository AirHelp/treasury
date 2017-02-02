package aws

import (
	"testing"

	"github.com/AirHelp/treasury/test"
)

func TestPutObject(t *testing.T) {
	dummyClient := &Client{S3Svc: &test.MockS3Client{}}
	object := &PutObjectInput{
		Bucket:      "dummyBucket",
		Key:         "test/application/key",
		Value:       "secret",
		Application: "application",
		Environment: "test",
	}
	if err := dummyClient.PutObject(object); err != nil {
		t.Fatal(err)
	}
}

func TestGetObject(t *testing.T) {
	dummyClient := &Client{S3Svc: &test.MockS3Client{}}
	s3objectInput := &GetObjectInput{
		Bucket: "dummyBucket",
		Key:    "test/application/key",
	}
	resp, err := dummyClient.GetObject(s3objectInput)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal()
	}
}
