package ssm_test

import (
	"testing"

	"github.com/AirHelp/treasury/backend/ssm"
	"github.com/aws/aws-sdk-go/aws"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		region  string
		config  aws.Config
		wantErr bool
	}{
		{
			name:    "empty region in args",
			region:  "",
			config:  aws.Config{},
			wantErr: false,
		},
		{
			name:    "region not defined",
			config:  aws.Config{},
			wantErr: false,
		},
		{
			name:    "valid region",
			region:  "eu-west-1",
			config:  aws.Config{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ssm.New(tt.region, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
