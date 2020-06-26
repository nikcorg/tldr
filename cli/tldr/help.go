package main

import (
	"flag"
	"strings"

	log "github.com/sirupsen/logrus"
)

type helpCmd struct{}

func (c *helpCmd) Execute(subcommand string, args ...string) error {
	log.Debugf("help:%s, args=%v", subcommand, strings.Join(args, "|"))
	flag.PrintDefaults()

	if len(args) > 0 && args[0] != "" {
		firstArg := args[0]
		runnableCommand, helpFocus, helpFocusSubcommand, restArgs := runnableForCommand(firstArg, args[1:])

		log.Debugf("focused help: %s:%s", helpFocus, helpFocusSubcommand)

		runnableCommand.Help(helpFocusSubcommand, restArgs...)
	}

	return nil
}

func (c *helpCmd) ParseArgs(subcommand string, args ...string) error {
	return nil
}

func (c *helpCmd) Help(subcommand string, args ...string) {}
