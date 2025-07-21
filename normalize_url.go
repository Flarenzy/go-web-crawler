package main

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	normalizedURL, _ := strings.CutSuffix(fmt.Sprintf("%s%s", parsedURL.Host, parsedURL.Path), "/")
	return normalizedURL, nil
}

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	return []string{}, fmt.Errorf("not implemented")
}
