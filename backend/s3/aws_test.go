package s3

import (
	"testing"
)

func TestAws(t *testing.T) {
	var tests = []struct {
		region  string
		bucket  string
		wantErr bool
	}{
		{
			region:  "",
			bucket:  "fakeBuckeName",
			wantErr: false,
		},
		{
			region:  "eu-west-1",
			bucket:  "fakeBuckeName",
			wantErr: false,
		},
		{
			region:  "eu-west-1",
			wantErr: true,
		},
		{
			wantErr: true,
		},
	}

	for _, test := range tests {
		_, err := New(test.region, test.bucket)
		if (err != nil) != test.wantErr {
			t.Errorf("New() error = %v, wantErr %v", err, test.wantErr)
			return
		}
	}
}
