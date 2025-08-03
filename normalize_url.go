package main

import (
	"fmt"
	"io"
	"net/http"
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
	if hrefURL == "" {
		return baseURL
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}
	ref, err := url.Parse(hrefURL)
	if err != nil {
		return ""
	}
	return base.ResolveReference(ref).String()
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

func getHTML(rawURL string) (string, error) {
	res, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}
	if res.StatusCode >= 400 {
		return "", fmt.Errorf("error getting link, http code %v", res.Status)
	}
	contentType := res.Header.Get("content-type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("expected html content type but got %v", contentType)
	}

	buff, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	return string(buff), nil
}
