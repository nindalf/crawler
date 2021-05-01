package internal

import (
	"os"
	"testing"
)

func TestExtract(t *testing.T) {
	testcases := []struct {
		string
		int
	}{
		{"./example-websites/blog.nindalf.com/index.html", 42},
		{"./example-websites/example.com/index.html", 1},
	}
	for _, testcase := range testcases {

		index_html, err := os.Open(testcase.string)
		if err != nil {
			t.Fatalf("Error opening %s - %v\n", testcase.string, err)
		}

		urls, err := Extract(index_html)
		if err != nil {
			t.Fatalf("Error extracting URLs from %s - %v\n", testcase.string, err)
		}
		if len(urls) != testcase.int {
			t.Fatalf("Expected %d urls in blog.nindalf.com, found - %d", testcase.int, len(urls))
		}
	}
}
