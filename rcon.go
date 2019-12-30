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
	rconPassword := mustGetPassword(l)

	// Set up exponential backoff
	b := &backoff.Backoff{
		Max:    10 * time.Minute,
		Jitter: true,
	}

	var r *rcon.RCON

	errDial := errors.New("placeholder")
	errAuth := errors.New("placeholder")

	for errDial != nil || errAuth != nil {
		r, errDial = rcon.Dial(rconAddress)
		if errDial != nil {
			l.Log(logging.Entry{
				Payload:  fmt.Sprintf("error dialling: %v", errDial),
				Severity: logging.Error,
			})
			time.Sleep(b.Duration())

			continue
		}

		errAuth = r.Authenticate(rconPassword)
		if errAuth != nil {
			l.Log(logging.Entry{
				Payload:  fmt.Sprintf("error authenticating: %v", errAuth),
				Severity: logging.Error,
			})
			r.Close()
			time.Sleep(b.Duration())
		}
	}

	return r
}