package config

import (
	"github.com/nikcorg/tldr-cli/config/rotation"
)

// Settings contains the runtime configuration
type Settings struct {
	ConfigPath string
	Configuration
}

// Configuration represents persistently configurable settings
type Configuration struct {
	Rotation rotation.Period `yaml:"rotation"`
	Storage  StorageConfig   `yaml:"storage"`
}

// StorageConfig represents storage settings
type StorageConfig struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}
