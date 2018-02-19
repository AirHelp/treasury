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
			name:    "ssm backend with no options",
			args:    backend.Options{},
			wantErr: false,
		},
		{
			name: "bo backend with bucket",
			args: backend.Options{
				S3BucketName: "fake_bucket_name",
			},
			wantErr: false,
		},
		{
			name: "ssm backend with Region",
			args: backend.Options{
				Backend: "ssm",
				Region:  "eu-west-1",
			},
			wantErr: false,
		},
		{
			name: "s3 backend without Bucket",
			args: backend.Options{
				Backend: "s3",
				Region:  "",
			},
			wantErr: true,
		},
		{
			name: "s3 backend with Bucket",
			args: backend.Options{
				Backend:      "s3",
				S3BucketName: "fake_bucket_name",
			},
			wantErr: false,
		},
		{
			name: "s3 backend with Bucket and Region",
			args: backend.Options{
				Backend:      "s3",
				S3BucketName: "fake_bucket_name",
				Region:       "eu-west-1",
			},
			wantErr: false,
		},
		{
			name: "invalid backend",
			args: backend.Options{
				Backend: "consul",
			},
			wantErr: true,
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
