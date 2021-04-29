package commands

import (
	"runtime"
	"strings"
)

const mod = `module %s

go %s
`

var (
	//Version define software version
	Version string
	// CommitHash represents the Git commit hash at built time
	CommitHash string
	// BuildDate represents the date when this tool was built
	BuildDate string
	// GoVersion represents the version of build go
	GoVersion string
)

func getGoVersion() string {
	v := runtime.Version()
	v = strings.TrimPrefix(v, "go")
	if strings.Count(v, ".") > 1 {
		index := strings.LastIndex(v, ".")
		v = v[:index]
	}
	return v
}
