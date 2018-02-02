package s3

import (
	"testing"
)

func TestAws(t *testing.T) {
	var tests = []struct {
		region  string
		wantErr bool
	}{
		{"", false},
		{"eu-west-1", false},
	}

	for _, test := range tests {
		_, err := New(test.region, "fakeBuckeName")
		if (err != nil) != test.wantErr {
			t.Errorf("New() error = %v, wantErr %v", err, test.wantErr)
			return
		}
	}
}
