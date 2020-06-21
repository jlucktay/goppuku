package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/logging"
	"github.com/ilyakaznacheev/cleanenv"
)

// Sets the name of the log to write to.
const logName = "goppuku"

// Address of RCON server.
const rconAddress = "127.0.0.1:27015"

// Number of minutes for the server to be empty before shutting down.
const shutdownMinutes = 15

func main() {
	// GCP project to send Stackdriver logs to.
	projectID, err := metadata.ProjectID()
	if err != nil {
		log.Fatalf("could not fetch project ID from metadata: %v", err)
	}

	// Create a logger client
	ctx := context.Background()

	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	logger := client.Logger(logName)

	// Handle incoming signals
	cleanup := make(chan os.Signal)
	signal.Notify(cleanup, os.Interrupt, syscall.SIGTERM)

	go func(l *logging.Logger) {
		<-cleanup
		l.Flush()
		os.Exit(1)
	}(logger)

	// Finish setting up logger
	notice := logger.StandardLogger(logging.Notice)
	notice.SetPrefix(fmt.Sprintf("%s[%d]: ", logName, os.Getpid()))
	notice.Print(versionDetails())

	notice.Printf("Configuring RCON from environment...")

	var cfg configRcon

	errReadConf := cleanenv.ReadEnv(&cfg)
	if errReadConf != nil {
		envDesc, _ := cleanenv.GetDescription(&cfg, nil)

		logger.Log(logging.Entry{
			Payload:  envDesc,
			Severity: logging.Warning,
		})

		logger.Log(logging.Entry{
			Payload:  fmt.Errorf("could not configure RCON from environment: %w", errReadConf),
			Severity: logging.Error,
		})

		cleanup <- syscall.SIGTERM
	}

	notice.Printf("Dialling '%s' and authing...", rconAddress)

	// Main loop
	r := dialAndAuth(logger, cfg)
	defer r.Close()
	notice.Print("Online!")
	monitor(r, logger)

	// Server seppuku
	cmd := exec.Command("shutdown", "--poweroff", "now")
	notice.Printf("Calling shutdown command: '%+v'", cmd)
	logger.Flush()

	_ = cmd.Start()
}
