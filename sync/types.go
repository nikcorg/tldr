package sync

import (
	"github.com/nikcorg/tldr-cli/config"
)

type Sync struct {
	config *config.Settings
}

func NewSync(config *config.Settings) Sync {
	return Sync{config}
}
