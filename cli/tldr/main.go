package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/nikcorg/tldr-cli/config"
	"github.com/nikcorg/tldr-cli/storage"

	log "github.com/sirupsen/logrus"
)

var configFile string = ""
var sourceDir string = ""
var sourceFile string = ""

var debugLogging bool = false
var verboseLogging bool = false

var runtimeConfig *config.Settings

var stor *storage.Storage

func main() {
	handleFlags()
	setLogLevel()

	runtimeConfig = config.NewWithDefaults()
	stor = storage.New(runtimeConfig)

	if err := mainWithErr(flag.Args()...); err != nil {
		log.Fatalf("Error running cmd: %s", err.Error())
	}
}

var cmdAdd = &addCmd{}
var cmdConfig = &configCmd{}
var cmdEdit = &editCmd{}
var cmdFind = &findCmd{}
var cmdHelp = &helpCmd{}
var cmdList = &listCmd{}

func splitCommand(cmd string) (string, string) {
	if cmd == "" {
		return "", ""
	} else if !strings.Contains(cmd, ":") {
		return cmd, ""
	}
	cmds := strings.SplitN(cmd, ":", 2)

	return cmds[0], cmds[1]
}

func runnableForCommand(firstArg string, args []string) (runnable, string, string, []string) {
	var (
		runnableCommand runnable
		nextArgs        []string = args
	)

	command, subcommand := splitCommand(firstArg)

	switch command {
	case "config":
		runnableCommand = cmdConfig
	case "edit":
		runnableCommand = cmdEdit
	case "find":
		runnableCommand = cmdFind
	case "list", "show":
		runnableCommand = cmdList
	case "help":
		runnableCommand = cmdHelp
	case "add":
		runnableCommand = cmdAdd
	default:
		subcommand = ""
		runnableCommand, nextArgs = defaultRunnable(firstArg)
	}

	return runnableCommand, command, subcommand, nextArgs
}

func defaultRunnable(arg string) (runnable, []string) {
	if strings.HasPrefix(arg, "http") {
		return cmdAdd, []string{arg}
	}

	return cmdHelp, []string{}
}

func mainWithErr(args ...string) error {
	var err error
	if err = runtimeConfig.Load(configFile); err != nil {
		return err
	}

	log.Debugf("Runtime config after Load %+v", runtimeConfig)

	var firstArg string = ""
	var restArgs []string = []string{}

	if len(args) > 0 {
		firstArg = args[0]
	}

	runnableCommand, command, subcommand, restArgs := runnableForCommand(firstArg, args[1:])

	if err = runnableCommand.Execute(subcommand, restArgs...); err != nil {
		if subcommand != "" {
			return fmt.Errorf("Error running %s:%s: %w", command, subcommand, err)

		}
		return fmt.Errorf("Error running %s: %w", command, err)
	}

	return nil
}

func handleFlags() {
	flag.StringVar(&configFile, "c", "", "Override config file")
	flag.StringVar(&sourceDir, "d", "", "Override storage location")
	flag.StringVar(&sourceFile, "f", "tldr.yaml", "Override storage file name (stem)")
	flag.BoolVar(&verboseLogging, "v", false, "Show verbose output")
	flag.BoolVar(&debugLogging, "vv", false, "Show debug output")
	flag.Parse()
}

func setLogLevel() {
	log.SetLevel(log.ErrorLevel)
	if debugLogging {
		log.SetLevel(log.DebugLevel)
	} else if verboseLogging {
		log.SetLevel(log.InfoLevel)
	}
}
