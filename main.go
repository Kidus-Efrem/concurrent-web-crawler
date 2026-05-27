package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	url := "https://example.com"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	fmt.Println(string(body))
}
