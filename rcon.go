package main

import (
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/logging"
	rcon "github.com/gtaylor/factorio-rcon"
	"github.com/jpillora/backoff"
)

var errPlaceholder = errors.New("placeholder")

// dialAndAuth creates the RCON client and authenticates with the server.
func dialAndAuth(logger *logging.Logger, cfg configRcon) *rcon.RCON {
	// Set up exponential backoff
	b := &backoff.Backoff{
		Min:    cfg.BackoffMin,
		Max:    cfg.BackoffMax,
		Factor: cfg.BackoffFactor,
		Jitter: true,
	}

	var r *rcon.RCON

	// Set placeholder errors before going into loop
	errDial, errAuth := fmt.Errorf("%w", errPlaceholder), fmt.Errorf("%w", errPlaceholder)

	for errDial != nil || errAuth != nil {
		r, errDial = rcon.Dial(rconAddress)
		if errDial != nil {
			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("error dialling address '%s': %v", rconAddress, errDial),
				Severity: logging.Error,
			})
			time.Sleep(b.Duration())

			continue
		}

		errAuth = r.Authenticate(mustGetPassword(logger))
		if errAuth != nil {
			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("error authenticating to address '%s': %v", rconAddress, errAuth),
				Severity: logging.Error,
			})
			r.Close()
			time.Sleep(b.Duration())
		}
	}

	return r
}
