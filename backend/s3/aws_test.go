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
		{"", "fakeBuckeName", false},
		{"eu-west-1", "fakeBuckeName", false},
		{"eu-west-1", "", true},
	}

	for _, test := range tests {
		_, err := New(test.region, test.bucket)
		if (err != nil) != test.wantErr {
			t.Errorf("New() error = %v, wantErr %v", err, test.wantErr)
			return
		}
	}
}
