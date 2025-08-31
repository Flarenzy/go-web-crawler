package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	rawURL := args[0]
	var maxConcurency int
	var maxPages int
	var err error
	if len(args) >= 2 {
		maxConcurency, err = strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Arg: %v is not a number!", args[1])
			os.Exit(1)
		}
	} else {
		maxConcurency = 5
	}
	if len(args) == 3 {
		maxPages, err = strconv.Atoi(args[2])
		if err != nil {
			fmt.Printf("Arg: %v is not a number!", args[1])
			os.Exit(1)
		}
	} else {
		maxPages = 200
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Printf("couldnt parse url %v\n", rawURL)
		os.Exit(1)
	}
	fmt.Printf("starting crawl of: %s://%s\n", parsedURL.Scheme, parsedURL.Host)
	concControlCH := make(chan struct{}, maxConcurency)
	cfg := config{
		pages:              make(map[string]int),
		maxPages:           maxPages,
		baseURL:            parsedURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: concControlCH,
		wg:                 &sync.WaitGroup{},
	}
	cfg.crawlPage(rawURL)
	cfg.wg.Wait()
	fmt.Printf("Crawled %v pages\n", len(cfg.pages))
	for k, v := range cfg.pages {
		fmt.Printf("Found URL %s %v many times\n", k, v)
	}
}
