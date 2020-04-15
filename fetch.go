package main

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/andybalholm/cascadia"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

type titlePuller struct {
	Name      string
	Selector  cascadia.Selector
	Extractor func(*html.Node) (string, error)
}

func attrValueFor(name string) func(n *html.Node) (string, error) {
	return func(n *html.Node) (string, error) {
		for _, attr := range n.Attr {
			if attr.Key == name {
				return attr.Val, nil
			}
		}
		return "", fmt.Errorf("Missing value attribute")
	}
}

func textContentFor(n *html.Node) (string, error) {
	if n.Type == html.TextNode {
		return n.Data, nil
	}

	if n.FirstChild == nil {
		return "", fmt.Errorf("No TextNodes found")
	}

	return textContentFor(n.FirstChild)
}

var selectors []titlePuller = []titlePuller{
	{"og:title", cascadia.MustCompile("meta[property=\"og:title\"]"), attrValueFor("content")},
	{"title", cascadia.MustCompile("title"), textContentFor},
	{"h1", cascadia.MustCompile("h1"), textContentFor},
	{"h2", cascadia.MustCompile("h2"), textContentFor},
	{"h3", cascadia.MustCompile("h3"), textContentFor},
	{".title", cascadia.MustCompile(".title"), textContentFor},
}

func getTitleCandidates(res *html.Node) ([]string, error) {
	var titles []string = []string{}

	for _, sel := range selectors {
		titleNode := cascadia.Query(res, sel.Selector)

		if titleNode != nil {
			titleText, err := sel.Extractor(titleNode)
			if err != nil {
				return nil, fmt.Errorf("Error extracting title: %w", err)
			}
			trimmedTitle := strings.TrimSpace(titleText)
			if len(trimmedTitle) > 0 {
				titles = append(titles, trimmedTitle)
				log.Debugf("Found title using %s: '%s'", sel.Name, trimmedTitle)
			}
		}
	}

	return titles, nil
}

type rankedTitle struct {
	Title string
	Score int
}

const (
	exactMatch        = 3
	includedByAnother = 2
	includesAnother   = 1 // this is most likely a site name suffixed title
)

func rankTitleCandidates(titles []string) ([]string, error) {
	var scoredTitles []rankedTitle = []rankedTitle{}

	if len(titles) == 0 {
		return nil, fmt.Errorf("No titles to rank")
	}

	for i, t := range titles {
		rt := rankedTitle{t, 0}
		for j, t2 := range titles {
			if i == j {
				continue
			}
			// Increase a title's rank when:
			// - It exactly matches another title
			if t == t2 {
				rt.Score += exactMatch
			}
			// - It is contained by another title
			if strings.Contains(t2, t) {
				rt.Score += includedByAnother
			}
			// - It contains another title
			if strings.Contains(t, t2) {
				rt.Score += includesAnother
			}
		}
		scoredTitles = append(scoredTitles, rt)
	}

	sort.SliceStable(scoredTitles, func(a, b int) bool {
		// Return a > b for descending rank order
		return scoredTitles[a].Score > scoredTitles[b].Score
	})

	log.Debugf("Titles scored: %+v", scoredTitles)

	// Return titles only
	rankedTitles := []string{}
	for _, t := range scoredTitles {
		rankedTitles = append(rankedTitles, t.Title)
	}

	return rankedTitles, nil
}

func uniqueTitles(allTitles []string) []string {
	if len(allTitles) == 0 {
		return allTitles
	}

	titles := []string{}
	for _, title := range allTitles {
		seen := false
		for _, t2 := range titles {
			if t2 == title {
				seen = true
				break
			}
		}
		if !seen {
			titles = append(titles, title)
		}
	}

	return titles
}

type fetchResult struct {
	URL    string
	Titles []string
}

func fetch(url string) (*fetchResult, error) {
	var err error
	var res *http.Response

	if res, err = http.Get(url); err != nil {
		return nil, err
	} else if res.StatusCode != 200 {
		return nil, fmt.Errorf("Error %v while fetching %s", res.StatusCode, url)
	}

	var body *html.Node
	body, err = html.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	var titleCandidates, rankedCandidates, uniqueCandidates []string

	titleCandidates, err = getTitleCandidates(body)
	rankedCandidates, err = rankTitleCandidates(titleCandidates)
	uniqueCandidates = uniqueTitles(rankedCandidates)

	log.Debugf("Found %d overall and %d ranked title candidates", len(titleCandidates), len(uniqueCandidates))
	for _, title := range uniqueCandidates {
		log.Debugf("- %s", title)
	}

	return &fetchResult{url, uniqueCandidates}, nil
}
