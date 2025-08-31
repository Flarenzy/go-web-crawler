package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"
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
	concControlCH := make(chan struct{}, 5)
	cfg := config{
		pages:              make(map[string]int),
		baseURL:            parsedURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: concControlCH,
		wg:                 &sync.WaitGroup{},
	}
	cfg.crawlPage(rawURL)
	cfg.wg.Wait()
	for k, v := range cfg.pages {
		fmt.Printf("Found URL %s %v many times\n", k, v)
	}
}
