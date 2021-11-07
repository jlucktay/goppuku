package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/logging"
	"github.com/davecgh/go-spew/spew"
	"github.com/ilyakaznacheev/cleanenv"
	"go.jlucktay.dev/version"
)

// Sets the name of the log to write to.
const logName = "goppuku"

// Address of RCON server.
const rconAddress = "127.0.0.1:27015"

// Run is the core of goppuku's logic/flow.
func Run(_ []string, stderr io.Writer) error {
	// GCP project to send Stackdriver logs to.
	projectID, err := metadata.ProjectID()
	if err != nil {
		return fmt.Errorf("could not fetch project ID from metadata: %w", err)
	}

	// Create a logger client
	ctx := context.Background()

	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to create logging client: %w", err)
	}
	defer client.Close() //nolint:errcheck // Don't let the door hit you on the way out
	logger := client.Logger(logName)

	// Handle incoming signals
	cleanup := make(chan os.Signal, 1)
	signal.Notify(cleanup, os.Interrupt, syscall.SIGTERM)

	go func(l *logging.Logger) {
		<-cleanup

		if errFlush := l.Flush(); errFlush != nil {
			fmt.Fprintf(stderr, "could not flush logger: %v", errFlush)
		}
	}(logger)

	// Finish setting up logger
	notice := logger.StandardLogger(logging.Notice)
	notice.SetPrefix(fmt.Sprintf("%s[%d]: ", logName, os.Getpid()))

	verDets, err := version.Details()
	if err != nil {
		cleanup <- syscall.SIGTERM

		return fmt.Errorf("could not get version details: %w", err)
	}

	notice.Print(verDets)
	notice.Print("Loading configuration from environment...")

	var cfg config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		envDesc, errGetDesc := cleanenv.GetDescription(&cfg, nil)
		if errGetDesc != nil {
			return fmt.Errorf("could not get description of environment variables: %w", errGetDesc)
		}

		logger.Log(logging.Entry{
			Payload:  envDesc,
			Severity: logging.Warning,
		})

		logger.Log(logging.Entry{
			Payload:  fmt.Errorf("could not load configuration from environment: %w", err),
			Severity: logging.Error,
		})

		cleanup <- syscall.SIGTERM

		return fmt.Errorf("could not read environment variables: %w", err)
	}

	notice.Printf("Configuration loaded:\n%s", spew.Sdump(cfg))

	// Store password in config struct _after_ spewing out all of the other settings
	cfg.RCON.Password = mustGetPassword(logger)

	// Main loop
	monitor(logger, cfg)

	// Server seppuku
	cmd := exec.Command("shutdown", "--poweroff", "now")
	notice.Printf("Calling shutdown command: '%s'", cmd)

	if err := logger.Flush(); err != nil {
		return fmt.Errorf("could not flush logger: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting '%s' command: %w", cmd, err)
	}

	return nil
}
