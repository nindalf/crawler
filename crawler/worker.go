package crawler

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type Worker struct {
	client *retryablehttp.Client
	wg     *sync.WaitGroup
}

func NewWorker(wg *sync.WaitGroup) *Worker {
	client := retryablehttp.NewClient()
	client.RetryMax = 3
	client.HTTPClient.Timeout = 10 * time.Second
	tr := http.Transport{
		MaxIdleConnsPerHost: 20,
	}
	client.HTTPClient.Transport = &tr
	return &Worker{client, wg}
}

func (w *Worker) Start(urlsToCrawl <-chan string, results chan<- string) {
	for url := range urlsToCrawl {
		extractedUrls := w.visit(url)
		for _, extractedUrl := range extractedUrls {
			// handle closed channel
			results <- extractedUrl
		}
	}
	w.wg.Done()
}

func (w *Worker) visit(url string) []string {
	fmt.Printf("Visiting url %s\n", url)
	resp, err := w.client.Get(url)
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
