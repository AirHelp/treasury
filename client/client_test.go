package client

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestNew(t *testing.T) {
	type args struct {
		options *Options
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		{
			name:    "empty options",
			args:    args{options: &Options{}},
			want:    &Client{},
			wantErr: false,
		},
		{
			name:    "add aws config",
			args:    args{options: &Options{AWSConfig: aws.Config{Region: "eu-west-1"}}},
			want:    &Client{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
