package config

import (
	"fmt"
	"strings"
)

// StorageGranularity defines the rotation frequency for storage
type StorageGranularity int

// Links storage rotation granularity
const (
	Unset StorageGranularity = iota
	None
	Daily
	Weekly
	Monthly
	Yearly
)

var (
	errUnknownGranularity = fmt.Errorf("Unknown granularity")
)

// Settings contains the runtime configuration
type Settings struct {
	ConfigPath  string
	Rotation    StorageGranularity
	StorageName string
	StoragePath string
}

func (s StorageGranularity) String() string {
	switch s {
	case None:
		return "none"
	case Daily:
		return "daily"
	case Weekly:
		return "weekly"
	case Monthly:
		return "monthly"
	case Yearly:
		return "yearly"
	}

	return ""
}

// RotationFromString maps a string to a StorageGranularity or panics
func RotationFromString(rot string) StorageGranularity {
	switch strings.ToLower(rot) {
	case "d", "daily":
		return Daily
	case "w", "weekly":
		return Weekly
	case "m", "monthly":
		return Monthly
	case "y", "yearly":
		return Yearly
	case "n", "none":
		return None
	}

	panic(fmt.Errorf("%w: %s", errUnknownGranularity, rot))
}
