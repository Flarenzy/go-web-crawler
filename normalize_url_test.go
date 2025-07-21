package main

import (
	"os"
	"testing"
)

func TestNormalizeUrl(t *testing.T) {
	tests := []struct {
		name        string
		inputURL    string
		expectedURL string
	}{
		{
			name:        "remove https scheme",
			inputURL:    "https://blog.boot.dev/path",
			expectedURL: "blog.boot.dev/path",
		}, {
			name:        "remove http scheme",
			inputURL:    "http://blog.boot.dev/path",
			expectedURL: "blog.boot.dev/path",
		}, {
			name:        "remove slash at end",
			inputURL:    "https://blog.boot.dev/path/",
			expectedURL: "blog.boot.dev/path",
		},
	}

	for _, te := range tests {
		t.Run(te.name, func(t *testing.T) {
			t.Parallel()
			resultURL, err := normalizeURL(te.inputURL)
			if err != nil {
				t.Errorf("Expected nil error but got %v", err)
				return
			}

			if resultURL != te.expectedURL {
				t.Errorf("Expected %s but got %s", te.expectedURL, resultURL)
			}
		})
	}
}

func TestGetUrlFromBody(t *testing.T) {
	rawBaseURL := "https://debian-handbook.info/"
	fp, err := os.ReadFile("debian-handbook.txt")
	if err != nil {
		t.Fatalf("Expected debian-handbook.txt but got err %v", err)
	}
	rawHTML := string(fp)
	rawURLs, err := getURLsFromHTML(rawHTML, rawBaseURL)
	if err != nil {
		t.Fatalf("Expected nil but got %v", err)
	}
	_ = rawURLs
}
