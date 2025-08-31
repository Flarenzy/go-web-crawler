package main

import (
	"fmt"
	"net/url"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	rawURL := args[0]
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Printf("couldnt parse url %v\n", rawURL)
		os.Exit(1)
	}
	fmt.Printf("starting crawl of: %s://%s\n", parsedURL.Scheme, parsedURL.Host)
	urlCount := make(map[string]int)
	crawlPage(rawURL, rawURL, urlCount)

	for k, v := range urlCount {
		fmt.Printf("Found URL %s %v many times\n", k, v)
	}
}
