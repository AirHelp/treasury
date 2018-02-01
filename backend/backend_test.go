package backend_test

import (
	"testing"

	"github.com/AirHelp/treasury/backend"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		args    backend.Options
		wantErr bool
	}{
		{
			name: "happy path with all options",
			args: backend.Options{
				Region:       "",
				S3BucketName: "fake_bucket_name",
			},
			wantErr: false,
		},
		{
			name: "no S3BucketName option set",
			args: backend.Options{
				Region: "",
			},
			wantErr: true,
		},
		{
			name: "no region in options",
			args: backend.Options{
				S3BucketName: "fake_bucket_name",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := backend.New(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
