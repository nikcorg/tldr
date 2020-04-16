package extract

import (
	"fmt"

	"golang.org/x/net/html"
)

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
	for currentNode := n; currentNode != nil; currentNode = currentNode.FirstChild {
		if currentNode.Type == html.TextNode {
			return currentNode.Data, nil
		}
	}

	return "", errNoTextNodes
}
