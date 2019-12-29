package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"

	"cloud.google.com/go/logging"
	rcon "github.com/gtaylor/factorio-rcon"
	"github.com/jpillora/backoff"
)

// Number of minutes for the server to be empty before shutting down
const shutdownMinutes = 15

func main() {
	// GCP project to send Stackdriver logs to
	const projectID = "jlucktay-factorio"
	// Sets the name of the log to write to
	const logName = "goppuku"
	// Address of RCON server
	const rconAddress = "127.0.0.1:27015"

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

	rconPassword := mustGetPassword(logger)

	// Set up exponential backoff
	b := &backoff.Backoff{
		Max:    10 * time.Minute,
		Jitter: true,
	}

	// Creates the RCON client and authenticates with the server
	r, errDial := rcon.Dial(rconAddress)
	for errDial != nil {
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("error dialing: %v", errDial),
			Severity: logging.Error,
		})
		time.Sleep(b.Duration())

		r, errDial = rcon.Dial(rconAddress)
	}
	b.Reset()

	defer r.Close()

	errAuth := r.Authenticate(rconPassword)
	for errAuth != nil && errDial != nil {
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("error authenticating: %v", errAuth),
			Severity: logging.Error,
		})
		time.Sleep(b.Duration())

		r.Close()

		r, errDial = rcon.Dial(rconAddress)
		if errDial != nil {
			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("error redialing: %v", errAuth),
				Severity: logging.Critical,
			})
			time.Sleep(b.Duration())

			continue
		}

		errAuth = r.Authenticate(rconPassword)
	}
	b.Reset()

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
