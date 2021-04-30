package crawler

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

func Extract(baseUrl string, body io.Reader) ([]string, error) {
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
					urls = append(urls, formatUrl(baseUrl, a.Val))
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

func formatUrl(baseUrl string, url string) string {
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	if strings.HasPrefix(url, "/") {
		url = strings.TrimSuffix(baseUrl, "/") + url
	}
	return url
}
