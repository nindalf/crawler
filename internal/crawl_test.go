package internal

import (
	"net/http"
	"testing"

	"github.com/nindalf/crawler/queue"
	"github.com/nindalf/crawler/storage"
)

const NumWorkers = 4

func TestCrawl(t *testing.T) {
	go http.ListenAndServe(":8000", http.FileServer(http.Dir("./example-websites/blog.nindalf.com")))
	storage := storage.NewMapStorage()
	workQueue := queue.NewChannelQueue(100)
	resultsQueue := queue.NewChannelQueue(100)
	startingURL := "http://localhost:8000/"
	crawler, err := NewCrawler(storage, workQueue, resultsQueue, startingURL, NumWorkers)
	if err != nil {
		t.Fatalf("Error creating crawler - %v\n", err)
	}
	crawler.Crawl()
	urls := storage.List()
	if len(urls) != 205 {
		t.Fatalf("Expected 205 URLs, found %d - ", len(urls))
	}
}
