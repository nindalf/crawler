package internal

import (
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/nindalf/crawler/queue"
	"github.com/nindalf/crawler/storage"
)

type Crawler struct {
	store         storage.Storage
	startingURL   string
	baseURL       string
	host          string
	workQueue     queue.Queue
	resultsQueue  queue.Queue
	numWorkers    int64
	activeWorkers int64
}

func NewCrawler(store storage.Storage, workQueue queue.Queue, resultsQueue queue.Queue, startingURL string, numWorkers int64) (Crawler, error) {
	urlParts, err := url.Parse(startingURL)
	if err != nil {
		return Crawler{}, err
	}
	startingURL = urlParts.String()
	baseURL := urlParts.Scheme + "://" + urlParts.Host
	host := urlParts.Host

	var activeWorkers int64

	return Crawler{store, startingURL, baseURL, host, workQueue, resultsQueue, numWorkers, activeWorkers}, nil
}

func (c Crawler) Crawl() {
	log.Println("Crawl starting")
	c.workQueue.Write(c.startingURL)
	c.store.Add(c.startingURL)

	for i := 0; i < int(c.numWorkers); i++ {
		worker := NewWorker(&c.activeWorkers, c.workQueue, c.resultsQueue)
		go worker.Start()
	}

	for {
		result, err := c.resultsQueue.ReadWithTimeout(time.Second)
		if err != nil {
			if err != queue.ErrTimeout {
				log.Fatalf("Unexpected error - %v\nTerminating ...", err)
			}
			if c.resultsQueue.Empty() && c.workQueue.Empty() && c.activeWorkers == 0 {
				log.Println("Crawl complete")
				break
			}
			continue
		}
		result = c.normalizeURL(result)
		resultURL, err := url.Parse(result)
		if err != nil {
			// Invalid URL
			continue
		}

		result = resultURL.String()
		if c.store.Contains(result) {
			// we've already visited this URL.
			continue
		}

		c.store.Add(result)
		if resultURL.Host == c.host {
			// it's on the same Domain, visit it
			c.workQueue.Write(result)
		}
	}
}

func (c Crawler) normalizeURL(url string) string {
	if strings.HasPrefix(url, "/") {
		url = c.baseURL + url
	}
	return url
}
