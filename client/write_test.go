package client_test

import (
	"testing"

	"github.com/AirHelp/treasury/client"
	test "github.com/AirHelp/treasury/test/backend"
)

func TestWrite(t *testing.T) {
	dummyClientOptions := &client.Options{
		Backend:      &test.MockBackendClient{},
		S3BucketName: "fake_s3_bucket",
	}
	treasury, err := client.New(dummyClientOptions)
	if err != nil {
		t.Error(err)
	}
	err = treasury.Write(test.Key1, test.KeyValueMap[test.Key1], false)
	if err != nil {
		t.Error(err)
	}
}
