package main

import (
	"fmt"

	"github.com/nikcorg/tldr-cli/config/rotation"
	"github.com/nikcorg/tldr-cli/storage"
	"github.com/nikcorg/tldr-cli/sync"
)

type syncCmd struct{}

func (s *syncCmd) ParseArgs(subcommand string, args ...string) error {
	return nil
}

func (s *syncCmd) Init() {}

func (s *syncCmd) Help(subcommand string, args ...string) {}

func (s *syncCmd) Execute(subcommand string, args ...string) error {
	if runtimeConfig.Sync.Exec == "" && runtimeConfig.Sync.Remote != "" {
		return fmt.Errorf("git sync not yet implemented")
	}

	if runtimeConfig.Rotation == rotation.None {
		return s.simpleSync()
	}

	return s.multiSync()
}

func (s *syncCmd) simpleSync() error {
	syncer := sync.NewSync(runtimeConfig)

	source, err := stor.Load()
	if err != nil {
		return err
	}

	return syncer.WithCommand([]*storage.Source{source})
}

func (s *syncCmd) multiSync() error {
	return nil
}
