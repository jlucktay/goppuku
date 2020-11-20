package cmd

import (
	"fmt"
	"time"

	"cloud.google.com/go/logging"
)

func monitor(logger *logging.Logger, cfg config) {
	// Keep track of how long the server has been empty for
	emptySince := time.Now().UTC()

	// Main monitoring loop
	for {
		time.Sleep(cfg.Monitor.MonitorInterval)

		r := dialAndAuth(logger, cfg.RCON)

		players, errCP := r.CmdPlayers()
		if errCP != nil {
			logger.Log(logging.Entry{
				Payload:  fmt.Errorf("error fetching player count: %w", errCP),
				Severity: logging.Error,
			})

			continue
		}

		r.Close()

		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("%+v", players),
			Severity: logging.Info,
		})

		anyOnline := false

		for _, player := range players {
			if player.Online {
				anyOnline = true
				emptySince = time.Now().UTC()

				break
			}
		}

		if !anyOnline {
			logger.Log(logging.Entry{
				Payload: fmt.Sprintf("Time without any online players: %s",
					time.Now().UTC().Sub(emptySince).Truncate(time.Second)),
				Severity: logging.Info,
			})
		}

		if emptySince.Add(cfg.Monitor.MonitorLimit).Before(time.Now().UTC()) {
			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("Threshold reached; %s elapsed without any online players", cfg.Monitor.MonitorLimit),
				Severity: logging.Notice,
			})

			return
		}
	}
}
