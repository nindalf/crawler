package crawler

import (
	"net/url"
	"sync"
	"testing"
)

func TestWorker(t *testing.T) {
	startingUrl := "https://example.com"
	url, _ := url.Parse(startingUrl)

	urlsToCrawl := make(chan string, 100)
	results := make(chan string, 100)
	urlsToCrawl <- startingUrl
	var wg sync.WaitGroup
	for i := 0; i < NUM_WORKERS; i++ {
		wg.Add(1)
		worker := NewWorker(url.Host, &wg)
		go worker.Start(urlsToCrawl, results)
	}

	close(urlsToCrawl)
	wg.Wait()

	if len(results) != 1 {
		t.Fatalf("Expected 1 result but got %d", len(results))
	}

	expected := "https://www.iana.org/domains/example/"
	result := <-results
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}

}
