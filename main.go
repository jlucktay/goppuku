package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"cloud.google.com/go/logging"
)

// GCP project to send Stackdriver logs to
const projectID = "jlucktay-factorio"

// Sets the name of the log to write to
const logName = "goppuku"

// Address of RCON server
const rconAddress = "127.0.0.1:27015"

// Number of minutes for the server to be empty before shutting down
const shutdownMinutes = 15

func main() {
	// Create a logger client
	ctx := context.Background()

	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	logger := client.Logger(logName)
	logger.Log(logging.Entry{
		Payload:  fmt.Sprintf("%s online", logName),
		Severity: logging.Notice,
	})

	r := dialAndAuth(logger)
	defer r.Close()
	monitor(r, logger)

	// Server seppuku
	cmd := exec.Command("shutdown", "--poweroff", "now")

	logger.Log(logging.Entry{
		Payload:  fmt.Sprintf("Calling shutdown command: '%+v'", cmd),
		Severity: logging.Notice,
	})
	logger.Flush()

	_ = cmd.Start()
}

// TODO: add log prefix with PID, e.g.:
// systemd[7691]:
