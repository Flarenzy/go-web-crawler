package main

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
)

type config struct {
	pages              map[string]int
	maxPages           int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func (cfg *config) isInternal(raw string) bool {
	cu, err := url.Parse(raw)
	if err != nil {
		return false
	}
	return strings.EqualFold(cfg.baseURL.Hostname(), cu.Hostname())
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	if cfg.pages[normalizedURL] > 0 {
		cfg.pages[normalizedURL]++
		return false
	}
	if len(cfg.pages) >= cfg.maxPages {
		return false
	}
	cfg.pages[normalizedURL] = 1
	return true
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.wg.Add(1)
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	if !cfg.isInternal(rawCurrentURL) {
		return
	}

	key, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	if !cfg.addPageVisit(key) {
		return
	}

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		return
	}
	fmt.Printf("Starting crawl of page %v\n", rawCurrentURL)
	links, err := getURLsFromHTML(html, rawCurrentURL)
	if err != nil {
		return
	}
	for _, u := range links {
		go cfg.crawlPage(u)
	}
}
