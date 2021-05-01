package internal

import (
	"net/http"
	"testing"
)

func TestWorker(t *testing.T) {
	go http.ListenAndServe(":8001", http.FileServer(http.Dir("./example-websites/example.com")))
	startingUrl := "http://localhost:8001/"

	urlsToCrawl := make(chan string, 100)
	results := make(chan string, 100)
	urlsToCrawl <- startingUrl
	var activeWorkers int64
	for i := 0; i < NUM_WORKERS; i++ {
		worker := NewWorker(&activeWorkers, urlsToCrawl, results)
		go worker.Start()
	}

	expected := "https://www.iana.org/domains/example"
	result := <-results
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}

	if len(urlsToCrawl) != 0 || len(results) != 0 || activeWorkers != 0 {
		t.Fatalf("Nonzero activity")
	}
}
