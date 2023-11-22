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

	tests := []struct {
		name    string
		key     string
		want    string
		wantErr bool
	}{
		{
			name:    "test non existing key",
			key:     test.Key9,
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := treasury.ReadValue(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ReadValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.ReadValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
