package cmd

import (
	"fmt"
	"runtime"
)

// Take ldflags from GoReleaser.
var version, commit, date, builtBy string //nolint:gochecknoglobals

func versionDetails() string {
	return fmt.Sprintf("goppuku %s built from commit %s with %s on %s by %s.",
		version, commit, runtime.Version(), date, builtBy)
}
