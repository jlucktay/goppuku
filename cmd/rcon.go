package cmd

import (
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/logging"
	factorioRCON "github.com/gtaylor/factorio-rcon"
	"github.com/jpillora/backoff"
)

var errPlaceholder = errors.New("placeholder")

// dialAndAuth creates the RCON client and authenticates with the server.
func dialAndAuth(logger *logging.Logger, cfg configRcon) *factorioRCON.RCON {
	// Set up exponential backoff
	backoff := &backoff.Backoff{
		Min:    cfg.BackoffMin,
		Max:    cfg.BackoffMax,
		Factor: cfg.BackoffFactor,
		Jitter: true,
	}

	var rcon *factorioRCON.RCON

	logger.Log(logging.Entry{Payload: fmt.Sprintf("Dialling '%s' and authing...", rconAddress)})

	// Set placeholder errors before going into loop
	errDial, errAuth := errPlaceholder, errPlaceholder

	for errDial != nil || errAuth != nil {
		rcon, errDial = factorioRCON.Dial(rconAddress)
		if errDial != nil {
			logger.Log(logging.Entry{
				Payload:  fmt.Errorf("error dialling address '%s': %w", rconAddress, errDial),
				Severity: logging.Error,
			})
			time.Sleep(backoff.Duration())

			continue
		}

		errAuth = rcon.Authenticate(cfg.Password)
		if errAuth != nil {
			logger.Log(logging.Entry{
				Payload:  fmt.Errorf("error authenticating to address '%s': %w", rconAddress, errAuth),
				Severity: logging.Error,
			})

			if errClose := rcon.Close(); errClose != nil {
				logger.Log(logging.Entry{
					Payload:  fmt.Errorf("error closing RCON: %w", errClose),
					Severity: logging.Error,
				})
			}

			time.Sleep(backoff.Duration())
		}
	}

	logger.Log(logging.Entry{Payload: "Online!"})

	return rcon
}
