package client_test

import (
	"testing"

	"github.com/AirHelp/treasury/client"
)

func TestClient(t *testing.T) {
	_, err := client.New("testBucketName", "", &client.Options{})
	if err != nil {
		t.Fatalf("Could not initialize client. Error:%s", err)
	}
}
