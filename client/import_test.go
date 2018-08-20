package client_test

import (
	"testing"

	"github.com/AirHelp/treasury/client"
	test "github.com/AirHelp/treasury/test/backend"
)

func TestImportS3(t *testing.T) {
	prefix := "test/webapp/"
	dummyClientOptions := &client.Options{
		Backend:      &test.MockBackendClient{},
		S3BucketName: "fake_s3_bucket",
	}
	treasury, err := client.New(dummyClientOptions)
	if err != nil {
		t.Error(err)
	}
	if err := treasury.Import(prefix, "../test/resources/import.env.test", false); err != nil {
		t.Error("Could not import secrets. Error: ", err.Error())
	}
}

func TestImportSSM(t *testing.T) {
	prefix := "test/webapp/"
	dummyClientOptions := &client.Options{
		Backend:      &test.MockBackendClient{},
	}
	treasury, err := client.New(dummyClientOptions)
	if err != nil {
		t.Error(err)
	}
	if err := treasury.Import(prefix, "../test/resources/import.env.test", false); err != nil {
		t.Error("Could not import secrets. Error: ", err.Error())
	}
}
