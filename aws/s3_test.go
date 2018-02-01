package aws

import (
	"path/filepath"
	"testing"

	test "github.com/AirHelp/treasury/test/s3"
	"github.com/AirHelp/treasury/types"
)

func TestPutObject(t *testing.T) {
	dummyClient := &Client{S3Svc: &test.MockS3Client{}, bucket: "dummyBucket"}
	object := &types.PutObjectInput{
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
	dummyClient := &Client{S3Svc: &test.MockS3Client{}, bucket: "dummyBucket"}
	s3objectInput := &types.GetObjectInput{Key: test.Key1}
	resp, err := dummyClient.GetObject(s3objectInput)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal()
	}
}

func TestGetObjects(t *testing.T) {
	dummyClient := &Client{S3Svc: &test.MockS3Client{}, bucket: "dummyBucket"}
	scenarios := []struct {
		input         *types.GetObjectsInput
		responsePairs int
	}{
		{
			input: &types.GetObjectsInput{
				Prefix: filepath.Dir(test.Key1) + "/",
			},
			responsePairs: 2,
		},
		{
			input: &types.GetObjectsInput{
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
