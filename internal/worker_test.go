package internal

import (
	"testing"

	"github.com/nindalf/crawler/queue"
)

func TestWorker(t *testing.T) {
	startingURL := serveDirectory(t, "./example-websites/example.com")

	workQueue := queue.NewChannelQueue(100)
	resultsQueue := queue.NewChannelQueue(100)
	workQueue.Write(startingURL)
	var activeWorkers int64
	for i := 0; i < 4; i++ {
		worker := NewWorker(&activeWorkers, workQueue, resultsQueue)
		go worker.Start()
	}

	expected := "https://www.iana.org/domains/example"
	result := resultsQueue.BlockingRead()
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}

	if !workQueue.Empty() || !resultsQueue.Empty() || activeWorkers != 0 {
		t.Fatalf("Nonzero activity")
	}
}
