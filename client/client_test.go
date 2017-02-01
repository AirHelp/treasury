package client_test

import (
	"testing"

	"github.com/AirHelp/treasury/client"
)

const (
	testKey    = "test/webapp/cocpit_api_pass"
	testSecret = "as9@#$%^&*(/2hdiwnf"
)

func TestClient(t *testing.T) {
	_, err := client.New("testBucketName", &client.Options{})
	if err != nil {
		t.Fatalf("Could not initialize client. Error:%s", err)
	}
}
