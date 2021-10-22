package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nikcorg/tldr-cli/input/entry"
	"github.com/nikcorg/tldr-cli/storage"

	log "github.com/sirupsen/logrus"
)

var (
	errMultipleMatches        = fmt.Errorf("multiple matches found")
	errNoEntryFound           = fmt.Errorf("no matching entry could be found")
	errNothingSelected        = fmt.Errorf("nothing selected")
	errSelectionOutsideBounds = fmt.Errorf("selection outside bounds")
)

type editCmd struct{}

func (e *editCmd) Init() {}

func (e *editCmd) ParseArgs(subcommand string, args ...string) error {
	return nil
}

func (e *editCmd) Execute(subcommand string, args ...string) error {
	var err error
	var matchedEntry *storage.Entry
	source, err := stor.Load()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		matchedEntry = source.FirstRecord().MostRecentEntry()
	} else {
		needle := args[0]
		matchedEntry, err = matchEntry(source, needle)
	}

	if err != nil {
		return err
	} else if matchedEntry == nil {
		return errNoEntryFound
	} else if err = edit(matchedEntry); err != nil {
		return err
	}

	log.Debugf("before save, %+v", source.FirstRecord())

	if err := stor.Save(source); err != nil {
		return err
	}

	return nil
}

func (e *editCmd) Help(subcommand string, args ...string) {
	log.Debugf("Help for %s, %v", subcommand, args)
}

func matchEntry(source *storage.Source, needle string) (*storage.Entry, error) {
	var err error
	var matchedEntry *storage.Entry

	results := locateMatches(source.Records, needle, 0, &[]findFilter{})
	switch len(results) {
	case 0:
		fmt.Printf("no matching entries were found for '%s'\n", needle)
		return nil, errNoEntryFound

	case 1:
		matchedEntry := results[0].Entry
		return matchedEntry, nil

	default:
		fmt.Printf("Multiple matches found for '%s', using interactive mode\n", needle)

		if matchedEntry, err = selectEntry(results); err != nil {
			return nil, err
		}

		return matchedEntry, nil
	}
}

func edit(matchedEntry *storage.Entry) error {
	if matchedEntry == nil {
		return errNoEntryFound
	}

	if err := entry.Edit(matchedEntry, &entry.EditContext{Titles: []string{matchedEntry.Title}}); err != nil {
		return err
	}

	return nil
}

func selectEntry(results []searchResult) (*storage.Entry, error) {
	var err error

	reader := bufio.NewReader(os.Stdin)

	for {
		for n, r := range results {
			e := r.Entry
			fmt.Printf("%d) %s (%v)\n", n, e.Title, r.Record.Date)
		}

		fmt.Printf("Q) Quit\n")
		fmt.Print("Your selection: ")

		var selection string
		selection, err = reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		selection = strings.ToUpper(strings.TrimSpace(selection))

		if len(selection) == 0 {
			fmt.Println("Invalid selection, try again")
			continue
		}

		switch selection {
		case "Q":
			fmt.Printf("Ok, quitting without saving.\n")
			os.Exit(0)
		default:
			if idx, err := strconv.ParseUint(selection, 10, 32); err != nil {
				return nil, err
			} else if int(idx) > len(results) {
				return nil, errSelectionOutsideBounds
			} else {
				return results[idx].Entry, nil
			}
		}
	}
}
