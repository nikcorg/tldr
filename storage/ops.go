package storage

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nikcorg/tldr-cli/config"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// Storage represents the current file
type Storage struct {
	config *config.Settings
}

// New constructs a Storage
func New(cfg *config.Settings) *Storage {
	return &Storage{cfg}
}

// Load retrieves a storage from disk
func (s *Storage) Load() (*Source, error) {
	storageName := s.sourceFileName()

	source, err := s.readSource(storageName)
	if err != nil {
		return nil, err
	}

	records := new([]*Record)
	if err := yaml.Unmarshal(source, records); err != nil {
		return nil, err
	}

	log.Debugf("Unmarshalled %d bytes into %d entries", len(source), len(*records))

	return &Source{
		SourceFile: storageName,
		Records:    *records,
	}, nil
}

// Save serialises a storage to disk
func (s *Storage) Save(source *Source) error {
	yamlString, err := yaml.Marshal(source.Records)
	if err != nil {
		return fmt.Errorf("Error serialising yaml: %w", err)
	}

	log.Debugf("Marshalled %d entries into %d bytes", source.Size(), len(yamlString))

	return s.writeSource(s.sourceFileName(), yamlString)
}

func (s *Storage) sourceFileName() string {
	return s.config.StorageFilePath()
}

func (s *Storage) readSource(fullSourcePath string) ([]byte, error) {
	var err error
	var source []byte

	_, err = os.Stat(fullSourcePath)
	if err != nil && os.IsNotExist(err) {
		log.Debugf("Source file does not exist: %s", fullSourcePath)
		return []byte{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("Error reading source: %w", err)
	}

	source, err = ioutil.ReadFile(fullSourcePath)
	if err != nil {
		return nil, fmt.Errorf("Error reading %s: %w", fullSourcePath, err)
	}

	log.Debugf("Read %d bytes from %s", len(source), fullSourcePath)

	return source, nil
}

func (s *Storage) writeSource(fullSourcePath string, b []byte) error {
	var err error

	err = os.MkdirAll(s.config.Storage.Path, 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("Error creating data dir: %s %w", s.config.Storage.Path, err)
	}

	err = ioutil.WriteFile(fullSourcePath, b, 0644)
	if err != nil && !os.IsNotExist(err) {
		log.Debugf("Error writing %s: %s", fullSourcePath)
		return fmt.Errorf("Error writing %s: %w", fullSourcePath, err)
	}

	log.Debugf("Wrote %d bytes into %s", len(b), fullSourcePath)

	return nil
}
