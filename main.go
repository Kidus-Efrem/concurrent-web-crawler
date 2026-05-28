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

	for len(queue) > 0 {
		// pop first element (BFS)
		url := queue[0]
		queue = queue[1:]

		if visited[url] {
			continue
		}
		visited[url] = true

		fmt.Println("Visiting:", url)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		doc, err := html.Parse(resp.Body)
		resp.Body.Close()

		if err != nil {
			fmt.Println("Parse error:", err)
			continue
		}

		links := extractLinks(doc)

		for _, link := range links {
			queue = append(queue, link)
		}
	}
}
func extractLinks(n *html.Node) []string {
	var links []string

	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				links = append(links, attr.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, extractLinks(c)...)
	}

	return links
}