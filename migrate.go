package main

// import "time"

// import (
// 	"fmt"
// 	"regexp"
// 	"strings"
// 	"time"

// 	log "github.com/sirupsen/logrus"
// 	yaml "gopkg.in/yaml.v2"
// )

// type record struct {
// 	Date    time.Time
// 	Entries []entry
// }

// type entry struct {
// 	RelatedURLs []string `yaml:"related_urls"`
// 	SourceURL   string   `yaml:"source_url"`
// 	Tags        []string
// 	Title       string
// 	URL         string
// 	Unread      bool
// }

// type mode int

// const (
// 	initial mode = iota
// 	foundRecord
// 	foundEntryStart

// 	foundSeeAlso
// )

// var matchRecordDate = regexp.MustCompile(`###.(\d{4})-(\d{2})-(\d{2})`)
// var matchURL = regexp.MustCompile(`https?://.*$`)

// func parseEntry(hunk []string) entry {
// 	var parserMode mode = initial

// 	unread := strings.Contains(hunk[0], "[ ]")
// 	title := strings.TrimSpace(hunk[0][5:])
// 	url := strings.TrimSpace(hunk[1])

// 	ret := entry{
// 		Title:  title,
// 		Unread: unread,
// 		URL:    url,
// 	}

// 	if len(hunk) > 2 {
// 		for ln, line := range hunk[2:] {
// 			trimmed := strings.TrimSpace(line)
// 			if strings.HasPrefix(trimmed, "Source") {
// 				matches := matchURL.FindStringSubmatch(trimmed)
// 				ret.SourceURL = matches[0]
// 			}

// 			if strings.HasPrefix(trimmed, "See also") {
// 				parserMode = foundSeeAlso
// 				matches := matchURL.FindStringSubmatch(trimmed)
// 				if len(matches) > 0 {
// 					ret.RelatedURLs = append(ret.RelatedURLs, strings.TrimSpace(matches[0]))
// 				} else {
// 					for _, line := range hunk[2+ln:] {
// 						matches := matchURL.FindStringSubmatch(line)
// 						if len(matches) > 0 {
// 							ret.RelatedURLs = append(ret.RelatedURLs, strings.TrimSpace(matches[0]))
// 						}
// 					}
// 				}
// 			}

// 			// See also always ends a block
// 			if parserMode == foundSeeAlso {
// 				break
// 			}
// 		}
// 	}

// 	return ret
// }

// func migrate(sourceLines []string) {
// 	var tldr []record = []record{}
// 	var parserStatus mode = initial
// 	var currRecord record = record{}
// 	var recordHunk []string

// 	for _, line := range sourceLines {
// 		if (parserStatus == initial || parserStatus == foundEntryStart) && strings.HasPrefix(line, "###") {
// 			if parserStatus == foundEntryStart {
// 				currRecord.Entries = append(currRecord.Entries, parseEntry(recordHunk))
// 			}
// 			if parserStatus != initial {
// 				tldr = append(tldr, currRecord)
// 			}

// 			parserStatus = foundRecord
// 			currRecord = record{}

// 			if parsed, err := time.Parse("### 2006-01-02", line); err == nil {
// 				currRecord.Date = parsed
// 			} else {
// 				currRecord.Date = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
// 			}

// 			continue
// 		} else if parserStatus == foundRecord && strings.HasPrefix(line, "-") {
// 			parserStatus = foundEntryStart
// 			recordHunk = make([]string, 0)
// 			recordHunk = append(recordHunk, line)
// 			continue
// 		} else if parserStatus == foundEntryStart && strings.HasPrefix(line, "-") {
// 			currRecord.Entries = append(currRecord.Entries, parseEntry(recordHunk))
// 			// entry := parseEntry(recordHunk)
// 			// log.Debugf("Found entry %+v", entry)
// 			recordHunk = make([]string, 0)
// 			recordHunk = append(recordHunk, line)
// 			continue
// 		} else if parserStatus == foundEntryStart {
// 			recordHunk = append(recordHunk, line)
// 		}
// 	}

// 	currRecord.Entries = append(currRecord.Entries, parseEntry(recordHunk))
// 	tldr = append(tldr, currRecord)

// 	yamlText, err := yaml.Marshal(&tldr)
// 	if err != nil {
// 		log.Fatalf("Error marshaling yaml: %v", err.Error())
// 	}

// 	fmt.Print(string(yamlText))
// }
