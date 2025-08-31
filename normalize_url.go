package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"golang.org/x/net/html"
)

// TODO: urls should all be lowercase, there can be multiple / at the end of path
func normalizeURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	p := path.Clean(parsedURL.Path)
	if p == "/" {
		p = ""
	}
	lowerCaseHost := strings.ToLower(parsedURL.Host)
	lowerCasePath := strings.ToLower(p)
	normalizedURL := fmt.Sprintf("%s%s", lowerCaseHost, lowerCasePath)
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

func isAsset(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	ext := strings.ToLower(path.Ext(u.Path))
	deny := map[string]struct{}{
		".png": {}, ".jpg": {}, ".jpeg": {}, ".gif": {}, ".svg": {},
		".css": {}, ".js": {}, ".mp4": {}, ".webm": {}, ".ico": {},
		".pdf": {},
	}
	_, blocked := deny[ext]
	return blocked
}

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	n, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return []string{}, err
	}
	rawURLs := make([]string, 0)

	var traverse func(node *html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, atr := range node.Attr {
				if atr.Key == "href" && !strings.HasPrefix(atr.Val, "#") {
					appendURL := appendBaseURL(rawBaseURL, atr.Val)
					if isAsset(appendURL) {
						continue
					}
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
