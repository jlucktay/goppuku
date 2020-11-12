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
)

// Sets the name of the log to write to.
const logName = "goppuku"

// Address of RCON server.
const rconAddress = "127.0.0.1:27015"

func Run(_ []string, _ io.Writer) error {
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

	notice.Printf("Loading configuration from environment...")

	var cfg config

	errReadConf := cleanenv.ReadEnv(&cfg)
	if errReadConf != nil {
		envDesc, _ := cleanenv.GetDescription(&cfg, nil)

		logger.Log(logging.Entry{
			Payload:  envDesc,
			Severity: logging.Warning,
		})

		logger.Log(logging.Entry{
			Payload:  fmt.Errorf("could not load configuration from environment: %w", errReadConf),
			Severity: logging.Error,
		})

		cleanup <- syscall.SIGTERM
	}

	notice.Printf("Configuration loaded:\n%s", spew.Sdump(cfg))
	notice.Printf("Dialling '%s' and authing...", rconAddress)

	// Main loop
	r := dialAndAuth(logger, cfg.RCON)
	defer r.Close()
	notice.Print("Online!")
	monitor(r, logger, cfg.Monitor)

	// Server seppuku
	cmd := exec.Command("shutdown", "--poweroff", "now")
	notice.Printf("Calling shutdown command: '%+v'", cmd)
	logger.Flush()

	_ = cmd.Start()

	return nil
}
