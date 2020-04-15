package main

import "time"

type record struct {
	Date    time.Time
	Entries []entry
}

type entry struct {
	RelatedURLs []string `yaml:"related_urls"`
	SourceURL   string   `yaml:"source_url"`
	Tags        []string
	Title       string
	URL         string
	Unread      bool
}
