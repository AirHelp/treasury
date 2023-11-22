package client_test

import (
	"testing"

	"github.com/AirHelp/treasury/client"
	test "github.com/AirHelp/treasury/test/backend"
)

func TestDelete(t *testing.T) {

	dummyClientOptions := &client.Options{
		Backend:      &test.MockBackendClient{},
		S3BucketName: "fake_s3_bucket",
	}

	treasury, err := client.New(dummyClientOptions)
	if err != nil {
		t.Error(err)
	}

	err = treasury.Delete(test.Key9)
	if err != nil {
		t.Error(err)
	}
}
