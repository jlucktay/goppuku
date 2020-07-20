package cmd

import "time"

type configRcon struct {
	BackoffMin    time.Duration `env:"GOPPUKU_BACKOFF_MIN_DURATION" env-default:"10s"`
	BackoffMax    time.Duration `env:"GOPPUKU_BACKOFF_MAX_DURATION" env-default:"5m"`
	BackoffFactor float64       `env:"GOPPUKU_BACKOFF_FACTOR" env-default:"1.2"`

	// MonitorInterval is the amount of time between checks on whether the server is empty.
	MonitorInterval time.Duration `env:"GOPPUKU_MONITOR_INTERVAL" env-default:"1m"`

	// ShutdownLimit is the amount of time for the server to be empty before shutting down.
	ShutdownLimit time.Duration `env:"GOPPUKU_SHUTDOWN_LIMIT" env-default:"15m"`
}
