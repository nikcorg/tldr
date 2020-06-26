package fetch

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"

	"github.com/nikcorg/tldr-cli/extract"

	log "github.com/sirupsen/logrus"
)

func any(url string) (*Details, error) {
	var err error
	var res *http.Response

	if res, err = http.Get(url); err != nil {
		return nil, err
	} else if res.StatusCode != 200 {
		return nil, fmt.Errorf("Error %v while fetching %s", res.StatusCode, url)
	}

	log.Debugf("Fetched URL: %s -> %s\n", url, res.Request.URL)

	var body *html.Node
	body, err = html.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	var titles, _ = extract.Titles(body)

	return &Details{res.Request.URL.String(), titles}, nil
}
