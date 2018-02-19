package client_test

import (
	"testing"

	"github.com/AirHelp/treasury/client"
)

func TestClient(t *testing.T) {
	tests := []struct {
		options *client.Options
	}{
		{&client.Options{S3BucketName: "fake_s3_bucket_name"}},
	}

	for _, test := range tests {
		if _, got := client.New(test.options); got != nil {
			t.Fatalf("Could not initialize client. Error:%s", got)
		}
	}
}
