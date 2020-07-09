package config

import (
	"os"
	"path"

	"github.com/nikcorg/tldr-cli/config/rotation"
)

func withDefaults(config *Settings) *Settings {
	return &Settings{
		Configuration: Configuration{
			Rotation: defaultRotation(config.Configuration.Rotation, rotation.None),
			Storage: StorageConfig{
				Name: defaultStorageName(config.Configuration.Storage.Name),
				Path: defaultStoragePath(config.Configuration.Storage.Path),
			},
		},
	}
}

func defaultStoragePath(candidatePath string) string {
	if candidatePath == "" {
		home, _ := os.UserHomeDir()
		return path.Join(home, "tldr")
	}

	return candidatePath
}

func defaultRotation(candidate rotation.Period, def rotation.Period) rotation.Period {
	if candidate == rotation.Unset {
		return def
	}

	return candidate
}

func defaultConfigFilename(override string, fallback string) string {
	if override == "" {
		return fallback
	}
	return override
}

func defaultStorageName(candidateName string) string {
	if candidateName == "" {
		return "tldr.yaml"
	}

	return candidateName
}
