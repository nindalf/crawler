package internal

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type Worker struct {
	client        *retryablehttp.Client
	activeCounter *int64
	urlsToCrawl   <-chan string
	results       chan<- string
}

func NewWorker(activeCounter *int64, urlsToCrawl <-chan string, results chan<- string) *Worker {
	client := retryablehttp.NewClient()
	client.RetryMax = 3
	client.HTTPClient.Timeout = 10 * time.Second
	tr := http.Transport{
		MaxIdleConnsPerHost: 20,
	}
	client.HTTPClient.Transport = &tr
	client.Logger = nil
	return &Worker{client, activeCounter, urlsToCrawl, results}
}

func (w *Worker) Start() {
	for url := range w.urlsToCrawl {
		extractedUrls := w.visit(url)
		for _, extractedUrl := range extractedUrls {
			w.results <- extractedUrl
		}
	}
}

func (w *Worker) visit(url string) []string {
	atomic.AddInt64(w.activeCounter, 1)
	resp, err := w.client.Get(url)
	atomic.AddInt64(w.activeCounter, -1)
	if err != nil {
		log.Fatalf("Non-retriable error while visiting URL - %v\n", err)
	}
	defer resp.Body.Close()
	urls, err := Extract(resp.Body)
	if err != nil {
		log.Printf("Failed to parse HTML on url %s. Continuing...\n", url)
	}
	return urls
}
