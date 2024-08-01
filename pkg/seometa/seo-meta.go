package seometa

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"nine-dubz/pkg/htmlcrawler"
)

func Set(node *html.Node, meta map[string]string) {
	for key, val := range meta {
		if val == "" {
			continue
		}

		metaNode := &html.Node{
			Type: html.ElementNode,
			Data: "meta",
			Attr: []html.Attribute{
				{Key: "property", Val: "og:" + key},
				{Key: "content", Val: val},
			},
		}
		node.AppendChild(metaNode)

		if key == "title" {
			titleNode := htmlcrawler.CrawlByTag("title", node)
			if titleNode != nil {
				if titleTextNode := titleNode.FirstChild; titleTextNode != nil {
					titleTextNode.Data = val
				}
			} else {
				titleNode = &html.Node{
					Type:     html.ElementNode,
					Data:     "title",
					DataAtom: atom.Title,
				}
				titleTextNode := &html.Node{
					Type: html.TextNode,
					Data: val,
				}
				titleNode.AppendChild(titleTextNode)
				node.AppendChild(titleNode)
			}
		}
	}
}
