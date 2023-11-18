// Package condition defines an interface for trigger conditions.
package condition

// Trigger condition upon which to shut the VM down.
type Trigger interface {
	Dial() error
	Check() (bool, error)
	HangUp() error
}
