package main

import (
	"fmt"
	"runtime"
)

// Take ldflags from GoReleaser.
var version, commit, date, builtBy string //nolint:gochecknoglobals

func versionDetails() string {
	return fmt.Sprintf("gopukku %s from commit %s, built %s by %s with %s.",
		version, commit, date, builtBy, runtime.Version())
}
