package aws

import (
	"testing"
)

func TestAws(t *testing.T) {
	var tests = []struct {
		region string
	}{
		{""},
		{"eu-west-1"},
	}

	for _, test := range tests {
		if _, got := New(test.region, "fakeBuckeName"); got != nil {
			t.Error(got)
		}
	}
}
