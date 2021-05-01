package crawler

import (
	"net/http"
	"testing"

	"github.com/nindalf/crawler/storage"
)

func TestCrawl(t *testing.T) {
	go http.ListenAndServe(":8000", http.FileServer(http.Dir("./example-websites/blog.nindalf.com")))
	startingUrl := "http://localhost:8000/"
	storage := storage.NewMapStorage()
	crawler, err := NewCrawler(storage, startingUrl)
	if err != nil {
		t.Fatalf("Error creating crawler - %v\n", err)
	}
	crawler.Crawl()
	urls := crawler.ListUrls()
	if len(urls) != 205 {
		t.Fatalf("Expected 205 URLs, found %d - ", len(urls))
	}
}
