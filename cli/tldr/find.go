package main

import (
	"fmt"
	"strings"
	"tldr/storage"

	log "github.com/sirupsen/logrus"
)

type findCmd struct{}

type findFilter func(entry *storage.Entry) bool
type findFilters []findFilter

type filterFlags struct {
	Unread bool
	Read   bool
}

func unreadFilter(state bool) findFilter {
	return func(e *storage.Entry) bool {
		return e.Unread == state
	}
}

var (
	showRelated = false
)

var (
	errUnreadReadConflict  = fmt.Errorf("Cannot use both --unread and --read")
	errArgsAfterSearchTerm = fmt.Errorf("Cannot pass flags after the search term")
	errUnknownArg          = fmt.Errorf("Unknown argument")
	errInvalidArgument     = fmt.Errorf("Invalid argument")
	errJunkAfterNeedle     = fmt.Errorf("Found junk after search term")
)

func (c *findCmd) Execute(subcommand string, args ...string) error {
	filters := []findFilter{}
	providedFlags := &filterFlags{}

	needleFound := false
	var needle string = ""

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") || strings.HasPrefix(arg, "--") {
			if needleFound {
				return errArgsAfterSearchTerm
			}

			switch arg {
			case "-u", "--unread":
				if providedFlags.Read == true {
					return errUnreadReadConflict
				}
				filters = append(filters, unreadFilter(true))
				providedFlags.Unread = true

			case "-r", "--read":
				if providedFlags.Unread == true {
					return errUnreadReadConflict
				}
				filters = append(filters, unreadFilter(false))
				providedFlags.Read = true

			case "-rel", "--related":
				log.Debugf("related, Set `showRelated=true`")
				showRelated = true

			default:
				log.Debugf("Unknown argument: %s", arg)
				return fmt.Errorf("%w: %s", errUnknownArg, arg)
			}
		} else if !needleFound {
			needle = strings.ToLower(arg)
			needleFound = true
		} else {
			return fmt.Errorf("%w: %s", errInvalidArgument, arg)
		}
	}

	if !needleFound {
		return fmt.Errorf("You must provide a search term")
	}

	log.Debugf("Find filters: +%v", filters)

	source, err := stor.Load()
	if err != nil {
		log.Debugf("Unhandled error loading storage: %v", err)
		return err
	}

	switch subcommand {
	case "one", "first":
		findFirst(source, needle, &filters)

	case "all":
		fallthrough

	default:
		findAll(source, needle, &filters)
	}

	return nil
}

func (c *findCmd) Help(subcommand string, args ...string) {
	log.Debugf("Help for %s, %v", subcommand, args)
}

type searchResult struct {
	Entry  *storage.Entry
	Record *storage.Record
}

func filtersMatch(entry *storage.Entry, filters *[]findFilter) bool {
	if filters == nil || len(*filters) == 0 {
		return true
	}

	for _, f := range *filters {
		if !f(entry) {
			return false
		}
	}

	return true
}

func locateMatches(stor *[]storage.Record, needle string, stopAfter int, filters *[]findFilter) []searchResult {
	searched := 0
	results := []searchResult{}
	for _, record := range *stor {
		for _, entry := range record.Entries {
			searched++
			if entry.Contains(needle) && filtersMatch(&entry, filters) {
				log.Debugf("Found needle (%s) in %+v added on %v", needle, entry, record.Date)
				e := entry
				r := record
				results = append(results, searchResult{&e, &r})

				if stopAfter > 0 && len(results) >= stopAfter {
					return results
				}
			}
		}
	}

	if len(results) == 0 {
		log.Debugf("Searched %d records, but found no match for needle '%s'", searched, needle)
	} else {
		log.Debugf("Searched %d records, found %d matches", searched, len(results))
	}

	return results
}

func noMatches(needle string) {
	fmt.Printf("No match found for \"%s\"\n", needle)
}

func oneMatch(sr searchResult, needle string) {
	entry := sr.Entry

	log.Debugf("Showing entry: %+v, Related: %v", entry, showRelated)

	if !entry.Unread {
		fmt.Printf("[x] %s\n%s\n", entry.Title, entry.URL)
	} else {
		fmt.Printf("[ ] %s\n%s\n", entry.Title, entry.URL)
	}

	if showRelated && len(entry.RelatedURLs) > 0 {
		fmt.Print("See also:\n")
		for _, rel := range entry.RelatedURLs {
			fmt.Printf("- %s\n", rel)
		}
	}
}

func findAll(source *storage.Source, needle string, filters *[]findFilter) {
	results := locateMatches(source.Records, needle, 0, filters)
	if len(results) == 0 {
		noMatches(needle)
		return
	}

	matches := "match"
	if len(results) > 1 {
		matches = "matches"
	}

	fmt.Printf("Found %d %s for \"%s\"\n", len(results), matches, needle)
	for _, rs := range results {
		oneMatch(rs, needle)
	}
}

func findFirst(source *storage.Source, needle string, filters *[]findFilter) {
	results := locateMatches(source.Records, needle, 1, filters)
	if len(results) == 0 {
		noMatches(needle)
		return
	}

	fmt.Printf("Found match for \"%s\" from %s\n", needle, results[0].Record.Date)
	oneMatch(results[0], needle)
}
