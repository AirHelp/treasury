package client_test

import (
	"testing"

	"github.com/AirHelp/treasury/client"
)

func TestClient(t *testing.T) {
	tests := []struct {
		bucketName string
		options    *client.Options
	}{
		{"testBucketName", &client.Options{}},
	}

	for _, test := range tests {
		if _, got := client.New(test.bucketName, test.options); got != nil {
			t.Fatalf("Could not initialize client. Error:%s", got)
		}
	}
}
