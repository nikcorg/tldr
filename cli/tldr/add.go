package main

import (
	"fmt"
	"strings"
	"time"

	"tldr/fetch"
	"tldr/input/entry"
	"tldr/storage"

	log "github.com/sirupsen/logrus"
)

type addCmd struct{}

func (c *addCmd) Execute(subcommand string, args ...string) error {
	log.Debugf("add:%s, args=%v", subcommand, args)

	source, err := stor.Load()
	if err != nil {
		return err
	}

	switch subcommand {
	case "title":
		amendTitle(source, args[0])

	case "source":
		amendSource(source, args[0])

	case "related":
		amendRelated(source, args[0])

	default:
		addEntry(args[0], source)
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

///

func addEntry(url string, source *storage.Source) error {
	log.Debugf("Fetching %v", url)
	var res *fetch.Details
	var err error
	if res, err = fetch.Preview(url); err != nil {
		return fmt.Errorf("Error fetching (%s): %w", url, err)
	}
	log.Debugf("Fetch result: %+v", res)

	var title string = ""

	if len(res.Titles) > 0 {
		title = res.Titles[0]
	}

	var newEntry = &storage.Entry{
		URL:    strings.ToLower(res.URL),
		Title:  title,
		Unread: true,
	}

	entry.Create(newEntry, &entry.EditContext{Titles: res.Titles})

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
