package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/spf13/pflag"
)

type Result struct {
	status int
	err    error
}

func main() {
	var (
		url        string
		totalReqs  int
		concurrent int
	)

	pflag.StringVar(&url, "url", "", "Target URL")
	pflag.IntVar(&totalReqs, "requests", 100, "Total number of requests")
	pflag.IntVar(&concurrent, "concurrency", 10, "Number of concurrent requests")
	pflag.Parse()

	if url == "" {
		fmt.Println("Error: --url is required")
		os.Exit(1)
	}

	start := time.Now()
	results := make(chan Result, totalReqs)
	wg := sync.WaitGroup{}

	sem := make(chan struct{}, concurrent)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for i := 0; i < totalReqs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			resp, err := client.Get(url)
			if err != nil {
				results <- Result{err: err}
			} else {
				results <- Result{status: resp.StatusCode}
				resp.Body.Close()
			}
			<-sem
		}()
	}

	wg.Wait()
	close(results)
	duration := time.Since(start)

	// Report
	count200 := 0
	statusDist := map[int]int{}
	errors := 0

	for res := range results {
		if res.err != nil {
			errors++
			continue
		}
		if res.status == 200 {
			count200++
		} else {
			statusDist[res.status]++
		}
	}

	fmt.Println("📊 Load Test Report")
	fmt.Printf("URL: %s\n", url)
	fmt.Printf("Total requests: %d\n", totalReqs)
	fmt.Printf("Concurrent requests: %d\n", concurrent)
	fmt.Printf("Total time: %s\n", duration)
	fmt.Printf("Successful requests (200): %d\n", count200)
	fmt.Printf("Other status codes:\n")
	for code, count := range statusDist {
		fmt.Printf("  - %d: %d\n", code, count)
	}
	fmt.Printf("Errors: %d\n", errors)
}
