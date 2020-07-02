package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/nikcorg/tldr-cli/fetch"
	"github.com/nikcorg/tldr-cli/input/entry"
	"github.com/nikcorg/tldr-cli/storage"

	log "github.com/sirupsen/logrus"
)

var (
	errURLNotFound error = fmt.Errorf("no URL argument found")
)

type addCmd struct {
	interactive bool
	url         string
	sourceURL   string
	relatedURLs []string
}

func (c *addCmd) Execute(subcommand string, args ...string) error {
	log.Debugf("add:%s, args=%v", subcommand, args)

	source, err := stor.Load()
	if err != nil {
		return err
	}

	switch subcommand {
	case "title":
		err = amendTitle(source, c.url)

	case "source":
		err = amendSource(source, c.url)

	case "related":
		err = amendRelated(source, c.url)

	default:
		err = c.addEntry(source)
	}

	if err != nil {
		return err
	}

	err = stor.Save(source)
	if err != nil {
		return err
	}

	return nil
}

func (c *addCmd) Help(subcommand string, args ...string) {
	log.Debugf("Help for %s, %v", subcommand, args)
}

func (c *addCmd) ParseArgs(subcommand string, args ...string) error {
	urlFound := false
	nextArg := 0
	for nextArg < len(args) {
		arg := args[nextArg]

		if strings.HasPrefix(arg, "-") || strings.HasPrefix(arg, "--") {
			switch arg {
			case "-i", "--interactive":
				c.interactive = true

			case "-r", "-rel", "--related":
				nextArg++
				url := args[nextArg]
				log.Debugf("found -r: %s", url)
				if !strings.HasPrefix(url, "http") {
					log.Debugf("invalid url: %s", url)
					return fmt.Errorf("%w: %s", errInvalidArg, url)
				}
				c.relatedURLs = append(c.relatedURLs, url)

			case "-s", "--source":
				nextArg++
				url := args[nextArg]
				log.Debugf("found -s: %s", url)
				if !strings.HasPrefix(url, "http") {
					return fmt.Errorf("%w: %s", errInvalidArg, url)
				}
				c.sourceURL = url
			}
		} else if strings.HasPrefix(arg, "http") {
			urlFound = true
			c.url = arg
		}

		nextArg++
	}

	if !urlFound && subcommand == "" {
		return errURLNotFound
	}

	return nil
}

///

func (c *addCmd) addEntry(source *storage.Source) error {
	log.Debugf("Fetching %v", c.url)
	var res *fetch.Details
	var err error
	if res, err = fetch.Preview(c.url); err != nil {
		return fmt.Errorf("Error fetching (%s): %w", c.url, err)
	}
	log.Debugf("Fetch result: %+v", res)

	var title string = ""

	if len(res.Titles) > 0 {
		title = res.Titles[0]
	}

	var newEntry = &storage.Entry{
		URL:         strings.ToLower(res.URL),
		Title:       title,
		Unread:      true,
		RelatedURLs: c.relatedURLs,
		SourceURL:   c.sourceURL,
	}

	if c.interactive {
		entry.Create(newEntry, &entry.EditContext{Titles: res.Titles})
	} else {
		fmt.Printf("Title: %s\n", newEntry.Title)
		fmt.Printf("URL: %s\n", newEntry.URL)
		fmt.Printf("Unread: %v\n", newEntry.Unread)
		if newEntry.SourceURL != "" {
			fmt.Printf("Source: %s\n", newEntry.SourceURL)
		}
		if len(newEntry.RelatedURLs) > 0 {
			fmt.Println("Related:")
			for _, u := range newEntry.RelatedURLs {
				fmt.Printf("- %s\n", u)
			}
		}
	}

	addEntryToTLDR(newEntry, source)

	return nil
}

func addEntryToTLDR(newEntry *storage.Entry, source *storage.Source) {
	y1, m1, d1 := time.Now().Date()

	if source.Size() == 0 {
		source.Records = &[]storage.Record{
			{
				Date:    time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC),
				Entries: []storage.Entry{*newEntry},
			},
		}
		return
	}

	y2, m2, d2 := (*source.Records)[0].Date.Date()

	if y1 == y2 && m1 == m2 && d1 == d2 {
		log.Debug("Entry for today already exists, appending")
		(*source.Records)[0].Entries = append((*source.Records)[0].Entries, *newEntry)
	} else {
		log.Debug("Entry for today doesn't exist, creating")
		newRecords := append([]storage.Record{
			{
				Date:    time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC),
				Entries: []storage.Entry{*newEntry},
			},
		}, (*source.Records)...)
		source.Records = &newRecords
	}
}

///

func amendRelated(source *storage.Source, url string) error {
	r := (*source.Records)[0]
	e := &r.Entries[len(r.Entries)-1]

	e.RelatedURLs = append(e.RelatedURLs, url)

	return nil
}

///

func amendSource(source *storage.Source, url string) error {
	r := (*source.Records)[0]
	e := &r.Entries[len(r.Entries)-1]

	e.SourceURL = url

	return nil
}

func amendTitle(source *storage.Source, title string) error {
	r := (*source.Records)[0]
	e := &r.Entries[len(r.Entries)-1]

	e.Title = title

	return nil
}
