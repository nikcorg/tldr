package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	errInvalidArg        = fmt.Errorf("Invalid argument")
	errExpectedNumberArg = fmt.Errorf("Expected numeric argument")
)

type listCmd struct {
	num       int
	offset    int
	newerThan *time.Time
}

func (f *listCmd) ParseArgs(subcommand string, args ...string) error {
	// FIXME: default page size should be in config
	f.num = -1
	f.offset = 0

	argsCopy := args[0:]

	for len(argsCopy) > 0 {
		arg := argsCopy[0]

		switch strings.ToLower(arg) {
		case "-t", "--today":
			now := time.Now()
			limit := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
			f.newerThan = &limit
			log.Debugf("newer than %v", f.newerThan)

		case "-n", "--num":
			num, err := strconv.Atoi(argsCopy[1])
			if err != nil {
				return fmt.Errorf("%w: %s", errExpectedNumberArg, argsCopy[1])
			}
			f.num = num

			// Shift args
			argsCopy = argsCopy[1:]

		case "-o", "--offset", "--skip":
			offset, err := strconv.Atoi(argsCopy[1])
			if err != nil {
				return fmt.Errorf("%w: %s", errExpectedNumberArg, argsCopy[1])
			}
			f.offset = offset

			// Shift args
			argsCopy = argsCopy[1:]

		default:
			return fmt.Errorf("%w: %s", errInvalidArg, arg)
		}

		// Shift args
		argsCopy = argsCopy[1:]
	}

	if f.newerThan == nil && f.num < 0 {
		return fmt.Errorf("unlimited listing not yet supported")
	}

	return nil
}

func (f *listCmd) Execute(subcommand string, args ...string) error {
	log.Debugf("list:%s, args=%v", subcommand, args)

	source, err := stor.Load()
	if err != nil {
		return err
	}

	displayed := 0
	skipped := 0
	for _, d := range *source.Records {
		if !d.Date.Equal(*f.newerThan) && !d.Date.After(*f.newerThan) {
			log.Debugf("%v < %v", d.Date, f.newerThan)
			break
		}

		for i := 0; i < len(d.Entries) && (f.num < 0 || displayed < f.num); i++ {
			e := d.Entries[i]
			if skipped < f.offset {
				skipped++
				continue
			}

			fmt.Printf("👉 %v, %+v\n", d.Date, e)

			displayed++
		}

		if f.num > -1 && displayed >= f.num {
			break
		}
	}

	return nil
}

func (f *listCmd) Help(subcommand string, args ...string) {
	fmt.Printf("Help for list: subcommand=%s, args=%v", subcommand, args)
}
