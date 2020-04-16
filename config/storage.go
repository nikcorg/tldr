package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const configFileName = "config.yaml"

var configDir string = ""

func init() {
	userConfigDir, _ := os.UserConfigDir()
	configDir = path.Join(userConfigDir, "tldr")

	log.Debugf("Resolved configDir=%s (%v)\n", configDir, userConfigDir)
}

// NewWithDefaults initalises a new Settings with default values set
func NewWithDefaults() *Settings {
	return withDefaults(&Settings{})
}

// Load reads and unserialises the runtime config from disk
func (s *Settings) Load(configFile string) error {
	fullPath := defaultConfigFilename(configFile, path.Join(configDir, configFileName))

	bytes, err := ioutil.ReadFile(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	log.Debugf("Loaded %d bytes of config from %s", len(bytes), fullPath)
	if len(bytes) == 0 {
		return nil
	}

	err = yaml.Unmarshal(bytes, s)
	if err != nil {
		return err
	}

	return nil
}

// Save serialises and writes the runtime config to disk
func (s *Settings) Save(configFile string) error {
	fullPath := defaultConfigFilename(configFile, path.Join(configDir, configFileName))

	bytes, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fullPath, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// StorageFilePath returns the fully qualified path to the actual storage file
func (s *Settings) StorageFilePath() string {
	return path.Join(s.StoragePath, sourceFileName(s))
}

func sourceFileName(cfg *Settings) string {
	rotationSuffix := currentStorageNameForRota(cfg)

	if len(rotationSuffix) == 0 {
		return cfg.StorageName
	}

	baseName := strings.TrimSuffix(cfg.StorageName, ".yaml")
	return fmt.Sprintf("%s.%s.yaml", baseName, rotationSuffix)
}

func currentStorageNameForRota(cfg *Settings) string {
	switch cfg.Rotation {
	case None:
		return ""
	case Monthly:
		return time.Now().Format("2006-01")
	case Daily:
		return time.Now().Format("2006-01-02")
	case Weekly:
		year, week := time.Now().ISOWeek()
		return fmt.Sprintf("%d-%d", year, week)
	case Yearly:
		return time.Now().Format("2006")
	}

	panic(fmt.Sprintf("Invalid rotation: %v", cfg.Rotation))
}
