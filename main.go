package main

import (
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

func main() {
	startURL := "https://example.com"

	queue := []string{startURL}
	visited := make(map[string]bool)

	base, err := url.Parse(startURL)
	if err != nil {
		fmt.Println("Base URL parse error:", err)
		return
	}

	for len(queue) > 0 {
		// BFS pop front
		currentURL := queue[0]
		queue = queue[1:]

		// Skip already visited URLs
		if visited[currentURL] {
			continue
		}
		visited[currentURL] = true

		fmt.Println("Visiting:", currentURL)

		resp, err := http.Get(currentURL)
		if err != nil {
			fmt.Println("Request error:", err)
			continue
		}

		doc, err := html.Parse(resp.Body)
		resp.Body.Close()

		if err != nil {
			fmt.Println("HTML parse error:", err)
			continue
		}

		links := extractLinks(doc)

		for _, link := range links {

			// Skip invalid links
			if link == "" {
				continue
			}

			if link[0] == '#' {
				continue
			}

			if len(link) > 10 && link[:10] == "javascript" {
				continue
			}

			if len(link) > 6 && link[:6] == "mailto" {
				continue
			}

			parsedLink, err := url.Parse(link)
			if err != nil {
				continue
			}

			resolved := base.ResolveReference(parsedLink)

			queue = append(queue, resolved.String())
		}
	}
}

func extractLinks(n *html.Node) []string {
	var links []string

	// Check if node is an <a> tag
	if n.Type == html.ElementNode && n.Data == "a" {

		// Look through attributes
		for _, attr := range n.Attr {

			// Find href
			if attr.Key == "href" {
				links = append(links, attr.Val)
			}
		}
	}

	// Traverse child nodes recursively
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		links = append(links, extractLinks(child)...)
	}

	return links
}