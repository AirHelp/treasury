package aws_test

import (
	"testing"

	"github.com/AirHelp/treasury/aws"
)

func TestAws(t *testing.T) {
	_, err := aws.New("")
	if err != nil {
		t.Error(err)
	}
}
