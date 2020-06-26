package main

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	errInvalidArg        = fmt.Errorf("Invalid argument")
	errExpectedNumberArg = fmt.Errorf("Expected numeric argument")
)

type listCmd struct{}

func (f *listCmd) ParseArgs(subcommand string, args ...string) error {
	return nil
}

func (f *listCmd) Execute(subcommand string, args ...string) error {
	log.Debugf("list:%s, args=%v", subcommand, args)

	// FIXME: default page size should be in config
	settings := struct {
		num    int
		offset int
	}{1, 0}

	argsCopy := args[0:]

	for len(argsCopy) > 0 {
		arg := argsCopy[0]

		switch strings.ToLower(arg) {
		case "-n", "--num":
			num, err := strconv.Atoi(argsCopy[1])
			if err != nil {
				return fmt.Errorf("%w: %s", errExpectedNumberArg, argsCopy[1])
			}
			settings.num = num
			argsCopy = argsCopy[1:]

		case "-o", "--offset", "--skip":
			offset, err := strconv.Atoi(argsCopy[1])
			if err != nil {
				return fmt.Errorf("%w: %s", errExpectedNumberArg, argsCopy[1])
			}
			settings.offset = offset
			argsCopy = argsCopy[1:]

		default:
			return fmt.Errorf("%w: %s", errInvalidArg, arg)
		}

		// Shift args
		argsCopy = argsCopy[1:]
	}

	source, err := stor.Load()
	if err != nil {
		return err
	}

	displayed := 0
	skipped := 0
	for _, d := range *source.Records {
		for i := 0; i < len(d.Entries) && displayed < settings.num; i++ {
			e := d.Entries[i]
			if skipped < settings.offset {
				skipped++
				continue
			}

			fmt.Printf("👉 %v, %+v\n", d.Date, e)

			displayed++
		}

		if displayed >= settings.num {
			break
		}
	}

	return nil
}

func (f *listCmd) Help(subcommand string, args ...string) {
	fmt.Printf("Help for list: subcommand=%s, args=%v", subcommand, args)
}
