package version

import (
	"fmt"
	"runtime"
)

// treasury version should be changed here
const version = "0.4.3"

// This will be filled in by the compiler.
var (
	// GitCommit of builded package
	gitCommit string
	// state of git tree, either "clean" or "dirty"
	gitTreeState string
	// BuildTime has time when binary is builded
	buildDate string = "1970-01-01T00:00:00Z"
)

type Version struct {
	Version, GitCommit, GitTreeState, BuildDate, GoVersion, Compiler, Platform string
}

// Get returns Version in struct
func Get() Version {
	return Version{
		Version:      version,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
