package internal

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/nindalf/crawler/queue"
	"github.com/nindalf/crawler/storage"
)

const NumWorkers = 4

func TestCrawl(t *testing.T) {
	storage := storage.NewMapStorage()
	workQueue := queue.NewChannelQueue(100)
	resultsQueue := queue.NewChannelQueue(100)
	startingURL := serveDirectory(t, "./example-websites/blog.nindalf.com")
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

func serveDirectory(t *testing.T, directory string) string {
	// Listen on a randomly assigned port
	conn, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to serve test directory - %v\n", err)
	}
	go func() {
		http.Serve(conn, http.FileServer(http.Dir(directory)))
	}()
	return fmt.Sprintf("http://%s/", conn.Addr().String())
}
