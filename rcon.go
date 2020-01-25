package main

import (
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/logging"
	rcon "github.com/gtaylor/factorio-rcon"
	"github.com/jpillora/backoff"
)

// dialAndAuth creates the RCON client and authenticates with the server.
func dialAndAuth(l *logging.Logger) *rcon.RCON {
	// Set up exponential backoff
	b := &backoff.Backoff{
		Min:    3 * time.Second,
		Max:    5 * time.Minute,
		Factor: 3,
		Jitter: true,
	}

	var r *rcon.RCON

	errDial := errors.New("placeholder")
	errAuth := errors.New("placeholder")

	for errDial != nil || errAuth != nil {
		r, errDial = rcon.Dial(rconAddress)
		if errDial != nil {
			l.Log(logging.Entry{
				Payload:  fmt.Sprintf("error dialling address '%s': %v", rconAddress, errDial),
				Severity: logging.Error,
			})
			time.Sleep(b.Duration())

			continue
		}

		errAuth = r.Authenticate(mustGetPassword(l))
		if errAuth != nil {
			l.Log(logging.Entry{
				Payload:  fmt.Sprintf("error authenticating to address '%s': %v", rconAddress, errAuth),
				Severity: logging.Error,
			})
			r.Close()
			time.Sleep(b.Duration())
		}
	}

	return r
}
