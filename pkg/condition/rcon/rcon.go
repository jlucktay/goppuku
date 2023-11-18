// Package rcon exports an RCON-based condition trigger.
package rcon

import (
	"go.jlucktay.dev/goppuku/pkg/condition"
)

var _ = []condition.Trigger{}

// Trigger is an idea waiting to happen.
type Trigger struct{}

// New makes one.
func New() (Trigger, error) {
	panic("not implemented")
}

// Dial dials.
func (t *Trigger) Dial() error {
	panic("not implemented")
}

// Check checks.
func (t *Trigger) Check() (bool, error) {
	panic("not implemented")
}

// HangUp hangups.
func (t *Trigger) HangUp() error {
	panic("not implemented")
}
