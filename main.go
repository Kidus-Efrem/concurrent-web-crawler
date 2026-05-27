package main

import (
	"fmt"
	"net/http"
	"golang.org/x/net/html"
)

func main() {
	url := "https://example.com"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	links := extractLinks(doc)

	for _, link := range links{
		fmt.Println((link))
	}

}

func extractLinks(n *html.Node)[]string{
	var links []string
	if n.Type == html.ElementNode && n.Data == "a"{
		for _, attr := range n.Attr{
			if attr.Key=="href"{
				links = append(links, attr.Val)
			}		}
	}
	for child := n.FirstChild; child !=nil; child = child.NextSibling{
		links = append(links, extractLinks(child)...)
	}
	return links
}