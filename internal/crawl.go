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
	startingUrl   string
	baseUrl       string
	host          string
	workQueue     queue.Queue
	resultsQueue  queue.Queue
	numWorkers    int64
	activeWorkers int64
}

func NewCrawler(store storage.Storage, workQueue queue.Queue, resultsQueue queue.Queue, startingUrl string, numWorkers int64) (Crawler, error) {
	urlParts, err := url.Parse(startingUrl)
	if err != nil {
		return Crawler{}, err
	}
	startingUrl = urlParts.String()
	baseUrl := urlParts.Scheme + "://" + urlParts.Host
	host := urlParts.Host

	var activeWorkers int64

	return Crawler{store, startingUrl, baseUrl, host, workQueue, resultsQueue, numWorkers, activeWorkers}, nil
}

func (c Crawler) ListUrls() []string {
	return c.store.List()
}

func (c Crawler) Crawl() {
	log.Println("Crawl starting")
	c.workQueue.Write(c.startingUrl)
	c.store.Add(c.startingUrl)

	for i := 0; i < int(c.numWorkers); i++ {
		worker := NewWorker(&c.activeWorkers, c.workQueue, c.resultsQueue)
		go worker.Start()
	}

	for {
		result, err := c.resultsQueue.ReadWithTimeout(time.Second)
		if err != nil {
			if err != queue.ERROR_TIMEOUT {
				log.Fatalf("Unexpected error - %v\nTerminating ...", err)
			}
			if c.resultsQueue.Empty() && c.workQueue.Empty() && c.activeWorkers == 0 {
				log.Println("Crawl complete")
				break
			}
			continue
		}
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
			c.workQueue.Write(result)
		}
	}
}

func (c Crawler) normalizeUrl(url string) string {
	if strings.HasPrefix(url, "/") {
		url = c.baseUrl + url
	}
	return url
}
