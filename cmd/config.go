package cmd

import "time"

type config struct {
	RCON    configRcon
	Monitor configMonitor
}

type configRcon struct {
	Password string `env:"GOPPUKU_RCON_PASSWORD"`

	BackoffMin time.Duration `env:"GOPPUKU_BACKOFF_MIN_DURATION" env-default:"10s"`
	BackoffMax time.Duration `env:"GOPPUKU_BACKOFF_MAX_DURATION" env-default:"5m"`

	BackoffFactor float64 `env:"GOPPUKU_BACKOFF_FACTOR" env-default:"1.2"`
}

type configMonitor struct {
	// MonitorInterval is the amount of time between checks on whether the server is empty.
	MonitorInterval time.Duration `env:"GOPPUKU_MONITOR_INTERVAL" env-default:"1m"`

	// MonitorLimit is the total amount of time for the server to be empty before shutting down.
	MonitorLimit time.Duration `env:"GOPPUKU_MONITOR_LIMIT" env-default:"15m"`
}
