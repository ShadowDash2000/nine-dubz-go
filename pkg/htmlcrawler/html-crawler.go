package htmlcrawler

import (
	"golang.org/x/net/html"
)

func CrawlByTag(tagName string, node *html.Node) *html.Node {
	if node.Type == html.ElementNode && node.Data == tagName {
		return node
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if node := CrawlByTag(tagName, c); node != nil {
			return node
		}
	}

	return nil
}
