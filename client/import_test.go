package client_test

import (
	"testing"

	"github.com/AirHelp/treasury/aws"
	"github.com/AirHelp/treasury/client"
	"github.com/AirHelp/treasury/test"
)

func TestImport(t *testing.T) {
	prefix := "test/webapp/"
	dummyClientOptions := &client.Options{
		AwsClient: &aws.Client{
			S3Svc: &test.MockS3Client{},
		},
	}
	treasury, err := client.New("fake_s3_bucket", dummyClientOptions)
	if err != nil {
		t.Error(err)
	}
	if err := treasury.Import(prefix, "../test/resources/import.env.test", false); err != nil {
		t.Error("Could not import secrets. Error: ", err.Error())
	}
}
