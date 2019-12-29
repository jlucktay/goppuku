package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"cloud.google.com/go/logging"
)

// mustGetPassword prepares the RCON password.
func mustGetPassword(l *logging.Logger) string {
	pwBytes, errRF := ioutil.ReadFile("/opt/factorio/config/rconpw")
	if errRF != nil {
		l.Log(logging.Entry{
			Payload:  fmt.Sprintf("error reading password file: %v", errRF),
			Severity: logging.Critical,
		})
		l.Flush()
		os.Exit(1)
	}

	return strings.TrimSpace(string(pwBytes))
}
