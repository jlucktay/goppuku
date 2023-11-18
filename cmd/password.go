package cmd

import (
	"bytes"
	"fmt"
	"os"

	"cloud.google.com/go/logging"
)

// mustGetPassword fetches the RCON password.
func mustGetPassword(logger *logging.Logger) string {
	bPassword, errRF := os.ReadFile("/opt/factorio/config/rconpw")
	if errRF != nil {
		logger.Log(logging.Entry{
			Payload:  fmt.Errorf("error reading password file: %w", errRF),
			Severity: logging.Critical,
		})

		if err := logger.Flush(); err != nil {
			fmt.Fprintf(os.Stderr, "could not flush logger: %v", err)
		}

		os.Exit(1)
	}

	bPasswordTrimmed := bytes.TrimSpace(bPassword)
	sPassword := string(bPasswordTrimmed)

	return sPassword
}
