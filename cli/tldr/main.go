package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/nikcorg/tldr-cli/config"
	"github.com/nikcorg/tldr-cli/storage"

	log "github.com/sirupsen/logrus"
)

var (
	configFile string = ""
	sourceDir  string = ""
	sourceFile string = ""

	debugLogging   bool = false
	verboseLogging bool = false

	configWasLoadedFromDisk bool = false
	runtimeConfig           *config.Settings

	stor *storage.Storage

	cmdAdd    = &addCmd{}
	cmdConfig = &configCmd{}
	cmdEdit   = &editCmd{}
	cmdFind   = &findCmd{}
	cmdHelp   = &helpCmd{}
	cmdList   = &listCmd{}
)

func main() {
	handleFlags()
	setLogLevel()

	runtimeConfig = config.NewWithDefaults()
	stor = storage.New(runtimeConfig)

	if err := mainWithErr(flag.Args()...); err != nil {
		log.Fatalf("Error running cmd: %s", err.Error())
	}
}

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
	case "add":
		runnableCommand = cmdAdd
	case "amend":
		runnableCommand = cmdAdd
		subcommand = "amend"
	case "config":
		runnableCommand = cmdConfig
	case "edit":
		runnableCommand = cmdEdit
	case "find":
		runnableCommand = cmdFind
	case "help":
		runnableCommand = cmdHelp
	case "list", "show":
		runnableCommand = cmdList
	default:
		subcommand = ""
		runnableCommand, nextArgs = defaultRunnable(firstArg, args)
	}

	return runnableCommand, command, subcommand, nextArgs
}

func defaultRunnable(firstArg string, args []string) (runnable, []string) {
	if strings.HasPrefix(firstArg, "http") {
		return cmdAdd, append([]string{firstArg}, args...)
	}

	return cmdHelp, append([]string{firstArg}, args...)
}

func mainWithErr(args ...string) error {
	var err error
	if err = runtimeConfig.Load(configFile); err != nil && err != config.ErrConfigFileNotFound {
		return err
	}

	configWasLoadedFromDisk = err != config.ErrConfigFileNotFound

	log.Debugf("Runtime config after Load (from disk? %v) %+v", configWasLoadedFromDisk, runtimeConfig)

	var (
		runnableCommand     runnable
		command, subcommand string
		restArgs            []string
	)

	if len(args) == 0 {
		runnableCommand, command, subcommand, restArgs = runnableForCommand("", []string{})
	} else if len(args) == 1 {
		runnableCommand, command, subcommand, restArgs = runnableForCommand(args[0], []string{})
	} else {
		runnableCommand, command, subcommand, restArgs = runnableForCommand(args[0], args[1:])
	}

	runnableCommand.Init()

	if err = runnableCommand.ParseArgs(subcommand, restArgs...); err != nil {
		return fmt.Errorf("%w: %s", errInvalidArg, err)
	}

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
