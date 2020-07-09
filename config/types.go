package config

import (
	"github.com/nikcorg/tldr-cli/config/rotation"
	"github.com/nikcorg/tldr-cli/config/sync"
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
	Sync     SyncConfig      `yaml:"sync"`
}

// StorageConfig represents storage settings
type StorageConfig struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

// SyncConfig represents synchronisation settings
type SyncConfig struct {
	Mode   sync.Mode `yaml:"mode"`
	Exec   string    `yaml:"exec,omitempty"`
	Remote string    `yaml:"remote,omitempty"`
}
