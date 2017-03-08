package aws_test

import (
	"testing"

	"github.com/AirHelp/treasury/aws"
)

func TestAws(t *testing.T) {
	var tests = []struct {
		region string
	}{
		{""},
		{"eu-west-1"},
	}

	for _, test := range tests {
		if _, got := aws.New(aws.Options{Region: test.region}); got != nil {
			t.Error(got)
		}
	}
}
