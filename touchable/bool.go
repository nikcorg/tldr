package touchable

import (
	log "github.com/sirupsen/logrus"
)

// Bool is a boolean container that tracks its untouched state
type Bool struct {
	touched bool
	value   bool
}

// NewBool creates a new Bool and sets the initial value
func NewBool(init bool) Bool {
	return Bool{
		touched: false,
		value:   init,
	}
}

// Set updates the value and sets the touched flag
func (b Bool) Set(v bool) bool {
	b.touched = true
	b.value = v

	log.Debugf("setting value and flagging as touched, %+v", b)

	return b.value
}

// Val returns the set value
func (b Bool) Val() bool {
	return b.value
}

// ValOrDefault returns the set value or the default if untouched
func (b Bool) ValOrDefault(def bool) bool {
	if !b.touched {
		return def
	}

	return b.value
}
