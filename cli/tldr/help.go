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
		helpFocus, helpFocusSubcommand := splitCommand(args[0])
		runnable := runnableForCommand(helpFocus)
		runnable.Help(helpFocusSubcommand, args[1:]...)
	}

	return nil
}

func (c *helpCmd) Help(subcommand string, args ...string) {}
