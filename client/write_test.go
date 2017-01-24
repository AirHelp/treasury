package client_test

import (
	"os"
	"testing"

	"github.com/AirHelp/treasury/client"
)

func TestWrite(t *testing.T) {
	treasuryURL := os.Getenv("TREASURY_URL")
	if treasuryURL == "" {
		t.Fatalf("TREASURY_URL environment variable is missing")
	}

	treasury, err := client.NewClient(treasuryURL, client.Options{})
	if err != nil {
		t.Error(err)
	}
	_, err = treasury.Write(testKey, testSecret)
	if err != nil {
		t.Error(err)
	}
}
