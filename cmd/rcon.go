package cmd

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

	logger.Log(logging.Entry{Payload: fmt.Sprintf("Dialling '%s' and authing...", rconAddress)})

	// Set placeholder errors before going into loop
	errDial, errAuth := errPlaceholder, errPlaceholder

	for errDial != nil || errAuth != nil {
		r, errDial = rcon.Dial(rconAddress)
		if errDial != nil {
			logger.Log(logging.Entry{
				Payload:  fmt.Errorf("error dialling address '%s': %w", rconAddress, errDial),
				Severity: logging.Error,
			})
			time.Sleep(b.Duration())

			continue
		}

		errAuth = r.Authenticate(cfg.Password)
		if errAuth != nil {
			logger.Log(logging.Entry{
				Payload:  fmt.Errorf("error authenticating to address '%s': %w", rconAddress, errAuth),
				Severity: logging.Error,
			})
			r.Close()
			time.Sleep(b.Duration())
		}
	}

	logger.Log(logging.Entry{Payload: "Online!"})

	return r
}
