package touchable

import (
	log "github.com/sirupsen/logrus"
)

// String is a string container that tracks its untouched state
type String struct {
	touched bool
	value   string
}

// NewString creates a new String and sets the initial value
func NewString(init string) *String {
	return &String{
		touched: false,
		value:   init,
	}
}

// Set updates the value and sets the touched flag
func (s *String) Set(v string) string {
	s.touched = true
	s.value = v

	log.Debugf("value set, %v, %+v", v, s)

	return v
}

// SetUnlessTouched updates the value unless it has been touched
func (s *String) SetUnlessTouched(v string) string {
	if !s.touched {
		return s.Set(v)
	}

	return v
}

// Val returns the set value
func (s *String) Val() string {
	return s.value
}

// ValOrDefault returns the set value or the default if untouched
func (s *String) ValOrDefault(def string) string {
	if !s.touched {
		return def
	}

	return s.value
}
