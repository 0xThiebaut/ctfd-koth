package conf

import (
	"errors"
	"time"
)

// Flag contains the needed information to monitor flag files.
type Flag struct {
	Interval time.Duration
	Path     string
	Award    *Award
}

// Check if the flag state seems credible.
func (f *Flag) Check() error {
	// A negative interval will cause a panic
	if f.Interval <= 0 {
		return errors.New("flag interval should be greater than zero (>0)")
	}
	// A non-existent path will never resolve
	if len(f.Path) == 0 {
		return errors.New("missing path for flag")
	}
	// Flags without awards are useless
	if f.Award == nil {
		return errors.New("missing award for flag")
	} else if err := f.Award.Check(); err != nil {
		return err
	}
	return nil
}
