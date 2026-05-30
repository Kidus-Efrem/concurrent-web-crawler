package main

import (
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

type CrawlItem struct {
	URL   string
	Depth int
}

func main() {
	startURL := "https://example.com"
	maxDepth := 2

	startParsed, err := url.Parse(startURL)
	if err != nil {
		fmt.Println("Start URL parse error:", err)
		return
	}

	queue := []CrawlItem{
		{
			URL:   startURL,
			Depth: 0,
		},
	}

	visited := make(map[string]bool)

	for len(queue) > 0 {

		current := queue[0]
		queue = queue[1:]

		currentURL := current.URL
		currentDepth := current.Depth

		if currentDepth > maxDepth {
			continue
		}

		if visited[currentURL] {
			continue
		}

		visited[currentURL] = true

		fmt.Printf("Visiting: %s (depth=%d)\n", currentURL, currentDepth)

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

		base, err := url.Parse(currentURL)
		if err != nil {
			continue
		}

		links := extractLinks(doc)

		for _, link := range links {

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

			// Stay inside original domain
			if resolved.Host != startParsed.Host {
				continue
			}

			queue = append(queue, CrawlItem{
				URL:   resolved.String(),
				Depth: currentDepth + 1,
			})
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

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		links = append(links, extractLinks(child)...)
	}

	return links
}