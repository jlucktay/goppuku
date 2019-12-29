package main

import (
	"fmt"
	"time"

	"cloud.google.com/go/logging"
	rcon "github.com/gtaylor/factorio-rcon"
)

func monitor(r *rcon.RCON, l *logging.Logger) {
	// Keep track of how long the server has been empty for
	minutesEmpty := 0

	// Main monitoring loop
	for {
		time.Sleep(time.Minute)

		players, errCP := r.CmdPlayers()
		if errCP != nil {
			l.Log(logging.Entry{
				Payload:  fmt.Sprintf("error fetching player count: %v", errCP),
				Severity: logging.Error,
			})

			continue
		}

		l.Log(logging.Entry{
			Payload:  fmt.Sprintf("%+v", players),
			Severity: logging.Info,
		})

		anyOnline := false

		for _, player := range players {
			if player.Online {
				anyOnline = true
				minutesEmpty = 0

				break
			}
		}

		if !anyOnline {
			minutesEmpty++
			l.Log(logging.Entry{
				Payload:  fmt.Sprintf("Minutes without any online players: %d", minutesEmpty),
				Severity: logging.Info,
			})
		}

		if minutesEmpty >= shutdownMinutes {
			l.Log(logging.Entry{
				Payload:  fmt.Sprintf("Threshold reached; %d minutes elapsed without any online players", shutdownMinutes),
				Severity: logging.Notice,
			})

			break
		}
	}
}
