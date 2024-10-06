package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func fetchURL(ctx context.Context, url string, resultChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Set up an HTTP client with a timeout
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		resultChan <- fmt.Sprintf("URL: %s\nError: %v", url, err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		resultChan <- fmt.Sprintf("URL: %s\nError: %v", url, err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resultChan <- fmt.Sprintf("URL: %s\nError: %v", url, err)
		return
	}

	content := string(body)
	if len(content) > 200 {
		content = content[:200] // Limiting output to first 200 characters
	}

	resultChan <- fmt.Sprintf("URL: %s\nContent: %s", url, content)
}

func main() {
	urls := []string{
		"https://example.com",
		"https://golang.org",
		"https://github.com",
		"https://stackoverflow.com",
		"https://wikipedia.org",
		"https://reddit.com",
		"https://news.ycombinator.com",
		"https://bbc.com",
		"https://cnn.com",
		"https://nytimes.com",
	}

	resultChan := make(chan string, len(urls))
	var wg sync.WaitGroup
	maxGoroutines := 5
	sem := make(chan struct{}, maxGoroutines) // Semaphore to limit goroutines

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			sem <- struct{}{}        // Acquire
			defer func() { <-sem }() // Release
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			fetchURL(ctx, url, resultChan, &wg)
		}(url)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		fmt.Println("---------------------------------")
		fmt.Println(result)
	}
}
