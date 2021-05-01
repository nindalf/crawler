package internal

import (
	"net/url"
	"strings"
	"time"

	"github.com/nindalf/crawler/storage"
)

const NUM_WORKERS = 4

type Crawler struct {
	store         storage.Storage
	startingUrl   string
	baseUrl       string
	host          string
	urlsToCrawl   chan string
	results       chan string
	numWorkers    int64
	activeWorkers int64
}

func NewCrawler(store storage.Storage, startingUrl string, numWorkers int64) (Crawler, error) {
	urlParts, err := url.Parse(startingUrl)
	if err != nil {
		return Crawler{}, err
	}
	startingUrl = urlParts.String()
	baseUrl := urlParts.Scheme + "://" + urlParts.Host
	host := urlParts.Host

	urlsToCrawl := make(chan string, 100)
	results := make(chan string, 100)
	var activeWorkers int64

	return Crawler{store, startingUrl, baseUrl, host, urlsToCrawl, results, numWorkers, activeWorkers}, nil
}

func (c Crawler) ListUrls() []string {
	return c.store.List()
}

func (c Crawler) Crawl() {

	c.urlsToCrawl <- c.startingUrl
	c.store.Add(c.startingUrl)

	for i := 0; i < NUM_WORKERS; i++ {
		worker := NewWorker(&c.activeWorkers, c.urlsToCrawl, c.results)
		go worker.Start()
	}
loop:
	for {
		select {
		case <-time.After(time.Second * 1):
			// If all workers are idle, and both channels are empty, break
			if len(c.urlsToCrawl) == 0 && len(c.results) == 0 && c.activeWorkers == 0 {
				close(c.urlsToCrawl)
				close(c.results)
				break loop
			}
			continue
		case result := <-c.results:
			result = c.normalizeUrl(result)
			resultUrl, err := url.Parse(result)
			if err != nil {
				// Invalid URL
				continue
			}

			result = resultUrl.String()
			if c.store.Contains(result) {
				// we've already visited this URL.
				continue
			}

			c.store.Add(result)
			if resultUrl.Host == c.host {
				// it's on the same Domain, visit it
				c.urlsToCrawl <- result
			}
		}
	}
}

func (c Crawler) normalizeUrl(url string) string {
	if strings.HasPrefix(url, "/") {
		url = c.baseUrl + url
	}
	return url
}
