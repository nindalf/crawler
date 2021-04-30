package crawler

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Worker struct {
	baseUrl string
	client  *http.Client
	wg      *sync.WaitGroup
}

func NewWorker(baseUrl string, wg *sync.WaitGroup) *Worker {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
		Timeout: 10 * time.Second,
	}
	return &Worker{baseUrl: baseUrl, client: client, wg: wg}
}

func (w *Worker) Start(urlsToCrawl <-chan string, results chan<- string) {
	for url := range urlsToCrawl {
		extractedUrls, err := w.visit(url)
		if err != nil {
			log.Fatalln(err)
		}
		for _, extractedUrl := range extractedUrls {
			// handle closed channel
			results <- extractedUrl
		}
	}
	w.wg.Done()
}

func (w *Worker) visit(url string) ([]string, error) {
	fmt.Printf("Visiting url %s\n", url)
	resp, err := w.client.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	return Extract(w.baseUrl, resp.Body)
}
