package main

import (
	"net/url"
	"strings"
)

func isInternal(base, raw string) bool {
	bu, _ := url.Parse(base)
	cu, _ := url.Parse(raw)
	return strings.EqualFold(bu.Hostname(), cu.Hostname())
}

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	if !isInternal(rawBaseURL, rawCurrentURL) {
		return
	}

	key, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	if pages[key] > 0 { // already seen, skip fetch
		pages[key]++
		return
	}
	pages[key] = 1

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		return
	}

	links, err := getURLsFromHTML(html, rawCurrentURL)
	if err != nil {
		return
	}
	for _, u := range links {
		crawlPage(rawBaseURL, u, pages)
	}
}
