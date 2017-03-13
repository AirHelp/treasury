package aws

import (
	"path/filepath"
	"testing"

	"github.com/AirHelp/treasury/test"
)

func TestPutObject(t *testing.T) {
	dummyClient := &Client{S3Svc: &test.MockS3Client{}}
	object := &PutObjectInput{
		Bucket:      "dummyBucket",
		Key:         test.Key1,
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
		Key:    test.Key1,
	}
	resp, err := dummyClient.GetObject(s3objectInput)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal()
	}
}

func TestGetObjects(t *testing.T) {
	dummyClient := &Client{S3Svc: &test.MockS3Client{}}
	scenarios := []struct {
		input         *GetObjectsInput
		responsePairs int
	}{
		{
			input: &GetObjectsInput{
				Bucket: "dummyBucket",
				Prefix: filepath.Dir(test.Key1) + "/",
			},
			responsePairs: 2,
		},
		{
			input: &GetObjectsInput{
				Bucket: "dummyBucket",
				Prefix: test.Key1,
			},
			responsePairs: 1,
		},
	}
	for _, scenario := range scenarios {
		resp, err := dummyClient.GetObjects(scenario.input)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal()
		}
		if len(resp.Secrets) != scenario.responsePairs {
			t.Errorf("wrong number of returned secrets, expected:%d, got:%d", scenario.responsePairs, len(resp.Secrets))
		}
	}
}
