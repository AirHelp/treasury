package version_test

import (
	"testing"

	"github.com/AirHelp/treasury/version"
)

func TestGet(t *testing.T) {
	version := version.Get()
	if version.Version == "" {
		t.Errorf("No default version given")
	}
}
