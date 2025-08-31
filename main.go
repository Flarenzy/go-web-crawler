package main

import (
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"sync"
)

type pageStat struct {
	pageURL string
	count   int
}

func printReport(pages map[string]int, baseURL string) {
	fmt.Println("=============================")
	fmt.Printf("REPORT for %s\n", baseURL)
	fmt.Println("=============================")

	stats := make([]pageStat, 0, len(pages))
	for url, c := range pages {
		stats = append(stats, pageStat{pageURL: url, count: c})
	}

	sort.Slice(stats, func(i, j int) bool {
		if stats[i].count != stats[j].count {
			return stats[i].count > stats[j].count // higher first
		}
		return stats[i].pageURL < stats[j].pageURL // tie-breaker: alphabetical
	})

	for _, s := range stats {
		fmt.Printf("Found %d internal links to %s\n", s.count, s.pageURL)
	}
}

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
		maxPages = 10
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
	fmt.Printf("Len of map %v\n", len(cfg.pages))
	printReport(cfg.pages, rawURL)
}
