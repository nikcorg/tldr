package extract

import (
	"fmt"
	"sort"
	"strings"

	"github.com/andybalholm/cascadia"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

var selectors []titlePuller = []titlePuller{
	{"og:title", cascadia.MustCompile("meta[property=\"og:title\"]"), attrValueFor("content")},
	{"twitter:title", cascadia.MustCompile("meta[property=\"twitter:title\"]"), attrValueFor("content")},
	{"title", cascadia.MustCompile("title"), textContentFor},
	{"h1", cascadia.MustCompile("h1"), textContentFor},
	{"h2", cascadia.MustCompile("h2"), textContentFor},
	{"h3", cascadia.MustCompile("h3"), textContentFor},
	{".title", cascadia.MustCompile(".title"), textContentFor},
}

// Titles find title candidates from a root html.Node,
// and returns a ranked, unique set
func Titles(root *html.Node) ([]string, error) {
	var titleCandidates, rankedCandidates, uniqueCandidates []string

	titleCandidates, _ = getTitleCandidates(root)
	rankedCandidates, _ = rankTitleCandidates(titleCandidates)
	uniqueCandidates = uniqueTitles(rankedCandidates)

	log.Debugf("Found %d overall and %d ranked title candidates", len(titleCandidates), len(uniqueCandidates))
	for _, title := range uniqueCandidates {
		log.Debugf("- %s", title)
	}

	return uniqueCandidates, nil
}

func getTitleCandidates(res *html.Node) ([]string, error) {
	var titles []string = []string{}

	for _, sel := range selectors {
		titleNode := cascadia.Query(res, sel.Selector)

		if titleNode != nil {
			titleText, err := sel.Extractor(titleNode)

			if err != nil {
				switch err {
				case errNoTextNodes:
					continue
				default:
					log.Debugf("Error extracting title using %s: %s", sel.Name, err.Error())
					return nil, fmt.Errorf("Error extracting title: %w", err)
				}
			}

			trimmedTitle := strings.TrimSpace(titleText)
			if len(trimmedTitle) > 0 {
				titles = append(titles, trimmedTitle)
				log.Debugf("Found title using %s: '%s'", sel.Name, trimmedTitle)
			}
		}
	}

	log.Debugf("Returning candidates: %v", titles)

	return titles, nil
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
