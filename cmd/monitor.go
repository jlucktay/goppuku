package cmd

import (
	"fmt"
	"time"

	"cloud.google.com/go/logging"
	rcon "github.com/gtaylor/factorio-rcon"
)

func monitor(r *rcon.RCON, logger *logging.Logger, cfg configMonitor) {
	// Keep track of how long the server has been empty for
	emptySince := time.Now().UTC()

	// Main monitoring loop
	for {
		time.Sleep(cfg.MonitorInterval)

		players, errCP := r.CmdPlayers()
		if errCP != nil {
			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("error fetching player count: %v", errCP),
				Severity: logging.Error,
			})

			continue
		}

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
				Payload:  fmt.Sprintf("Time without any online players: %s", time.Now().UTC().Sub(emptySince)),
				Severity: logging.Info,
			})
		}

		if emptySince.Add(cfg.MonitorLimit).Before(time.Now().UTC()) {
			logger.Log(logging.Entry{
				Payload:  fmt.Sprintf("Threshold reached; %s elapsed without any online players", cfg.MonitorLimit),
				Severity: logging.Notice,
			})

			return
		}
	}
}
