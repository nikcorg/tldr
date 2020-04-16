package extract

import (
	"fmt"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

var (
	errNoTextNodes = fmt.Errorf("No TextNodes found")
)

type titlePuller struct {
	Name      string
	Selector  cascadia.Selector
	Extractor func(*html.Node) (string, error)
}

type rankedTitle struct {
	Title string
	Score int
}
