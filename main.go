package main

import (
	"flag"
	"log"

	"github.com/nindalf/crawler/crawler"
	"github.com/nindalf/crawler/storage"
)

var (
	url     = flag.String("u", "https://example.com", "The URL to be crawled")
	workers = flag.Int64("w", 4, "The number of workers to run simultaneously")
)

func main() {
	storage := storage.NewMapStorage()
	crawler, err := crawler.NewCrawler(storage, *url, *workers)
	if err != nil {
		log.Fatalf("Error creating crawler - %v\n", err)
	}
	crawler.Crawl()
	urls := crawler.ListUrls()
	log.Println(urls)
}
