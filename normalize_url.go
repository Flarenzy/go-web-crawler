package main

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func normalizeURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	normalizedURL, _ := strings.CutSuffix(fmt.Sprintf("%s%s", parsedURL.Host, parsedURL.Path), "/")
	return normalizedURL, nil
}

func appendBaseURL(baseURL, hrefURL string) string {
	if strings.HasPrefix(hrefURL, "/") {
		return fmt.Sprintf("%s%s", baseURL, hrefURL)
	}
	return hrefURL
}

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	n, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return []string{}, err
	}
	rawURLs := make([]string, 0)

	var traverse func(node *html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode {
			for _, atr := range node.Attr {
				if atr.Key == "href" {
					appendURL := appendBaseURL(rawBaseURL, atr.Val)
					rawURLs = append(rawURLs, appendURL)
				}
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)

	return rawURLs, nil
}
