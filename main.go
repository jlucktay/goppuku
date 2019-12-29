package main

import (
	"context"
	"fmt"
	"log"
	"os"
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

	notice := logger.StandardLogger(logging.Notice)
	notice.SetPrefix(fmt.Sprintf("%s[%d]", logName, os.Getpid()))

	notice.Printf("Dialling '%s' and authing...", rconAddress)

	r := dialAndAuth(logger)
	defer r.Close()

	notice.Print("Online!")
	monitor(r, logger)

	// Server seppuku
	cmd := exec.Command("shutdown", "--poweroff", "now")

	notice.Printf("Calling shutdown command: '%+v'", cmd)
	logger.Flush()

	_ = cmd.Start()
}
