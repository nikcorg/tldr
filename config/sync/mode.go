package sync

import (
	"fmt"
	"strconv"
)

// Mode defines the means of remote sync
type Mode int

// Known error outcomes
var (
	ErrUnknownMode = fmt.Errorf("unknown mode")
	ErrInvalidMode = fmt.Errorf("invalid mode")
)

// Rotation
const (
	Unset Mode = iota
	Command
	Git
)

const (
	strUnset   = "unset"
	strCommand = "command"
	strGit     = "git"
)

// NewFromString maps from a string to a Mod or panic
func NewFromString(s string) Mode {
	switch s {
	case strUnset, strUnset[0:1]:
		return Unset
	case strCommand, strCommand[0:1]:
		return Command
	case strGit, strGit[0:1]:
		return Git
	}

	panic(fmt.Errorf("%w: %s", ErrUnknownMode, s))
}

func (p Mode) String() string {
	switch p {
	case Unset:
		return strUnset
	case Command:
		return strCommand
	case Git:
		return strGit
	}

	panic(fmt.Errorf("%w: %d", ErrInvalidMode, p))
}

// MarshalYAML implements the YAML Marshaler interface
func (p Mode) MarshalYAML() (interface{}, error) {
	return p.String(), nil
}

// UnmarshalYAML implements the YAML Unmarshaler interface
func (p *Mode) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string

	if err := unmarshal(&s); err != nil {
		return err
	}

	if i, err := strconv.Atoi(s); err == nil {
		*p = Mode(i)
	} else {
		*p = NewFromString(s)
	}

	return nil
}
