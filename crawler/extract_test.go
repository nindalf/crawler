package crawler

import (
	"os"
	"testing"
)

func TestExtract(t *testing.T) {
	nindalfIndex, _ := os.Open("./example-websites/blog.nindalf.com/index.html")
	urls, err := Extract("https://blog.nindalf.com/", nindalfIndex)
	if err != nil {
		t.Fatal(err)
	}
	if len(urls) != 42 {
		t.Fatalf("Expected 42 urls in blog.nindalf.com, found - %d", len(urls))
	}
	exampleIndex, _ := os.Open("./example-websites/example.com/index.html")
	urls, err = Extract("example.com", exampleIndex)
	if err != nil {
		t.Fatal(err)
	}
	if len(urls) != 1 {
		t.Fatalf("Expected 1 urls in example.com, found - %d", len(urls))
	}
}
