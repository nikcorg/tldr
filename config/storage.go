package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/nikcorg/tldr-cli/config/rotation"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	configFileName = "config.yaml"
	tldrConfigDir  = "TLDR_CONFIG_DIR"
)

var configDir string = ""

// error outcomes
var (
	ErrConfigFileNotFound = fmt.Errorf("no config file found")
)

func init() {
	envConfigDir := os.Getenv(tldrConfigDir)
	if envConfigDir != "" {
		configDir = envConfigDir
	} else {
		userConfigDir, _ := os.UserConfigDir()
		configDir = path.Join(userConfigDir, "tldr")
	}

	log.Debugf("Resolved configDir=%s\n", configDir)
}

// NewWithDefaults initalises a new Settings with default values set
func NewWithDefaults() *Settings {
	return withDefaults(&Settings{})
}

// Load reads and unserialises the runtime config from disk
func (s *Settings) Load(configFile string) error {
	fullPath := defaultConfigFilename(configFile, path.Join(configDir, configFileName))

	bytes, err := ioutil.ReadFile(fullPath)
	if os.IsNotExist(err) {
		return ErrConfigFileNotFound
	} else if err != nil {
		return err
	}

	log.Debugf("Loaded %d bytes of config from %s", len(bytes), fullPath)
	if len(bytes) == 0 {
		return nil
	}

	err = yaml.Unmarshal(bytes, &s.Configuration)
	if err != nil {
		return err
	}

	return nil
}

// Save serialises and writes the runtime config to disk
func (s *Settings) Save(configFile string) error {
	fullPath := defaultConfigFilename(configFile, path.Join(configDir, configFileName))

	bytes, err := yaml.Marshal(s.Configuration)
	if err != nil {
		return err
	}

	err = os.MkdirAll(configDir, 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("Error creating data dir: %s %w", configDir, err)
	}

	log.Debugf("writing %d bytes to %s", len(bytes), fullPath)

	err = ioutil.WriteFile(fullPath, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// StorageFilePath returns the fully qualified path to the actual storage file
func (s *Settings) StorageFilePath() string {
	return path.Join(s.Storage.Path, sourceFileName(s))
}

func sourceFileName(cfg *Settings) string {
	rotationSuffix := currentStorageNameForRota(cfg)

	if len(rotationSuffix) == 0 {
		return cfg.Storage.Name
	}

	baseName := strings.TrimSuffix(cfg.Storage.Name, ".yaml")
	return fmt.Sprintf("%s.%s.yaml", baseName, rotationSuffix)
}

func currentStorageNameForRota(cfg *Settings) string {
	switch cfg.Rotation {
	case rotation.None:
		return ""
	case rotation.Monthly:
		return time.Now().Format("2006-01")
	case rotation.Daily:
		return time.Now().Format("2006-01-02")
	case rotation.Weekly:
		year, week := time.Now().ISOWeek()
		return fmt.Sprintf("%d-%d", year, week)
	case rotation.Yearly:
		return time.Now().Format("2006")
	}

	panic(fmt.Sprintf("Invalid rotation: %v", cfg.Rotation))
}
