package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"cloud.google.com/go/logging"
)

// mustGetPassword fetches the RCON password.
func mustGetPassword(l *logging.Logger) string {
	bPassword, errRF := ioutil.ReadFile("/opt/factorio/config/rconpw")
	if errRF != nil {
		l.Log(logging.Entry{
			Payload:  fmt.Errorf("error reading password file: %w", errRF),
			Severity: logging.Critical,
		})

		if err := l.Flush(); err != nil {
			fmt.Fprintf(os.Stderr, "could not flush logger: %v", err)
		}

		os.Exit(1)
	}

	bPasswordTrimmed := bytes.TrimSpace(bPassword)
	sPassword := string(bPasswordTrimmed)

	return sPassword
}
