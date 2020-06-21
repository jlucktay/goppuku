package main

import "time"

type configRcon struct {
	BackoffMin    time.Duration `env:"GOPUKKU_BACKOFF_MIN_DURATION" env-default:"10s" env-description:"RCON backoff minimum duration"`
	BackoffMax    time.Duration `env:"GOPUKKU_BACKOFF_MAX_DURATION" env-default:"5m" env-description:"RCON backoff maximum duration"`
	BackoffFactor float64       `env:"GOPUKKU_BACKOFF_FACTOR" env-default:"1.2" env-description:"RCON backoff factor"`
}
