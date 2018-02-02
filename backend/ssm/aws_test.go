package ssm_test

import (
	"testing"

	"github.com/AirHelp/treasury/backend/ssm"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		region  string
		wantErr bool
	}{
		{
			name:    "empty region in args",
			region:  "",
			wantErr: false,
		},
		{
			name:    "valid region",
			region:  "eu-west-1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ssm.New(tt.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
