package client_test

import (
	"testing"

	"github.com/AirHelp/treasury/aws"
	"github.com/AirHelp/treasury/client"
	"github.com/AirHelp/treasury/test"
)

func TestRead(t *testing.T) {
	dummyClientOptions := &client.Options{
		AwsClient: &aws.Client{
			S3Svc: &test.MockS3Client{},
		},
	}
	treasury, err := client.New("fake_s3_bucket", "", dummyClientOptions)
	if err != nil {
		t.Error(err)
	}
	secret, err := treasury.Read(test.TestKey)
	if err != nil {
		t.Error(err)
	}
	if secret.Value != test.TestSecret {
		t.Errorf("Reads returns wrong secret")
	}
}
