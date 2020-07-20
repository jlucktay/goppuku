package cmd

import "time"

type configRcon struct {
	BackoffMin    time.Duration `env:"GOPPUKU_BACKOFF_MIN_DURATION" env-default:"10s"`
	BackoffMax    time.Duration `env:"GOPPUKU_BACKOFF_MAX_DURATION" env-default:"5m"`
	BackoffFactor float64       `env:"GOPPUKU_BACKOFF_FACTOR" env-default:"1.2"`
}
