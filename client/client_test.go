package client_test

import (
	"testing"

	"github.com/AirHelp/treasury/client"
)

func TestClient(t *testing.T) {
	tests := []struct {
		bucketName, region string
		options            *client.Options
	}{
		{"testBucketName", "", &client.Options{}},
		{"testBucketName", "eu-west-1", &client.Options{}},
	}

	for _, test := range tests {
		if _, got := client.New(test.bucketName, test.region, test.options); got != nil {
			t.Fatalf("Could not initialize client. Error:%s", got)
		}
	}
}
