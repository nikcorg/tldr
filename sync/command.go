package sync

import (
	"syscall"

	"github.com/nikcorg/tldr-cli/storage"
)

func (s *Sync) WithCommand(sources []*storage.Source) error {
	var args []string = []string{s.config.Storage.Path}

	for _, source := range sources {
		args = append(args, source.SourceFile)
	}

	if err := syscall.Exec(s.config.Sync.Exec, args, []string{}); err != nil {
		return err
	}

	return nil
}
