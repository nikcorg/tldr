package config

import (
	"os"
	"path"
)

func withDefaults(config *Settings) *Settings {
	return &Settings{
		Rotation:    defaultRotation(config.Rotation, None),
		StorageName: defaultStorageName(config.StorageName),
		StoragePath: defaultStoragePath(config.StoragePath),
	}
}

func defaultStoragePath(candidatePath string) string {
	if candidatePath == "" {
		home, _ := os.UserHomeDir()
		return path.Join(home, "tldr")
	}

	return candidatePath
}

func defaultRotation(candidate StorageGranularity, def StorageGranularity) StorageGranularity {
	if candidate == Unset {
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
