package main

import (
	"fmt"
)

// Take ldflags from GoReleaser
var (
	//nolint
	version, commit, date, builtBy string
)

func versionDetails() string {
	return fmt.Sprintf("gopukku v%s from commit %s, built %s by %s.", version, commit, date, builtBy)
}
