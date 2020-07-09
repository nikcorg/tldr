package main

import (
	"fmt"
	"strings"

	"github.com/nikcorg/tldr-cli/config/rotation"
	log "github.com/sirupsen/logrus"
)

var (
	errUnknownSetting = fmt.Errorf("unknown setting")
)

type configCmd struct {
	forced bool
}

func (c *configCmd) Init() {
	c.forced = false
}

func (c *configCmd) ParseArgs(subcommand string, args ...string) error {
	for _, arg := range args {
		switch arg {
		case "-f", "--force":
			if subcommand != "init" {
				return fmt.Errorf("%w: %s can only be used with `init`", errInvalidArg, arg)
			}
			c.forced = true

		default:
			return fmt.Errorf("%w: %s", errUnknownArg, arg)
		}
	}

	return nil
}

func (c *configCmd) Execute(subcommand string, args ...string) error {
	log.Debugf("config:%s, args=%v", subcommand, strings.Join(args, "|"))

	var changed bool
	var err error

	switch subcommand {
	case "set":
		changed, err = c.set(args[0], args[1])

	case "get":
		err = c.get(args[0])

	case "init":
		changed = c.forced || !configWasLoadedFromDisk
	}

	if err == nil && changed {
		return runtimeConfig.Save(configFile)
	}

	return err
}

func (c *configCmd) Help(subcommand string, args ...string) {
	log.Debugf("Help for %s, %v", subcommand, args)
}

func (c *configCmd) set(key, value string) (bool, error) {
	switch strings.ToLower(key) {
	case "rotation":
		if rot := rotation.NewFromString(value); rot != runtimeConfig.Rotation {
			runtimeConfig.Rotation = rot
			return true, nil
		}

	case "storage.path":
		if runtimeConfig.Storage.Path != value {
			runtimeConfig.Storage.Path = value
			return true, nil
		}

	case "storage.name":
		if runtimeConfig.Storage.Name != value {
			runtimeConfig.Storage.Name = value
			return true, nil
		}

	default:
		return false, fmt.Errorf("%w: %s", errUnknownSetting, key)
	}

	return false, nil
}

func (c *configCmd) get(key string) error {
	switch strings.ToLower(key) {
	case "rotation":
		fmt.Printf("rotation=%s\n", runtimeConfig.Rotation.String())

	case "storage.path":
		fmt.Printf("storage.path=%s\n", runtimeConfig.Storage.Path)

	case "storage.name":
		fmt.Printf("storage.name=%s\n", runtimeConfig.Storage.Name)

	default:
		log.Debugf("Unknown setting: %s", key)
		return fmt.Errorf("%w: %s", errUnknownSetting, key)
	}

	return nil
}
