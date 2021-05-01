package internal

import (
	"io"

	"golang.org/x/net/html"
)

func Extract(body io.Reader) ([]string, error) {
	urls := make([]string, 0, 4)
	doc, err := html.Parse(body)
	if err != nil {
		return urls, err
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					urls = append(urls, a.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return urls, nil
}
