package rotation

import (
	"fmt"
	"strconv"
)

// Period defines the time between file rotation
type Period int

// Known error outcomes
var (
	ErrUnknownPeriod = fmt.Errorf("unknown period")
	ErrInvalidPeriod = fmt.Errorf("invalid period")
)

// Rotation
const (
	Unset Period = iota
	None
	Daily
	Weekly
	Monthly
	Yearly
)

const (
	strUnset   = "unset"
	strNone    = "none"
	strDaily   = "daily"
	strWeekly  = "weekly"
	strMonthly = "monthly"
	strYearly  = "yearly"
)

// NewFromString maps from a string to a Period or panic
func NewFromString(s string) Period {
	switch s {
	case strUnset, strUnset[0:1]:
		return Unset
	case strNone, strNone[0:1]:
		return None
	case strDaily, strDaily[0:1]:
		return Daily
	case strWeekly, strWeekly[0:1]:
		return Weekly
	case strMonthly, strMonthly[0:1]:
		return Monthly
	case strYearly, strYearly[0:1]:
		return Yearly
	}

	panic(fmt.Errorf("%w: %s", ErrUnknownPeriod, s))
}

func (p Period) String() string {
	switch p {
	case Unset:
		return strUnset
	case None:
		return strNone
	case Daily:
		return strDaily
	case Weekly:
		return strWeekly
	case Monthly:
		return strMonthly
	case Yearly:
		return strYearly
	}

	panic(fmt.Errorf("%w: %d", ErrInvalidPeriod, p))
}

// MarshalYAML implements the YAML Marshaler interface
func (p Period) MarshalYAML() (interface{}, error) {
	return p.String(), nil
}

// UnmarshalYAML implements the YAML Unmarshaler interface
func (p *Period) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string

	if err := unmarshal(&s); err != nil {
		return err
	}

	if i, err := strconv.Atoi(s); err == nil {
		*p = Period(i)
	} else {
		*p = NewFromString(s)
	}

	return nil
}
