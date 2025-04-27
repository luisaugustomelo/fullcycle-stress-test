package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoadTester(t *testing.T) {
	// Start a simple HTTP server for testing
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/notfound" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tests := []struct {
		name        string
		url         string
		requests    int
		concurrency int
		expect200   bool
	}{
		{
			name:        "Successful requests",
			url:         server.URL,
			requests:    10,
			concurrency: 2,
			expect200:   true,
		},
		{
			name:        "404 responses",
			url:         server.URL + "/notfound",
			requests:    5,
			concurrency: 1,
			expect200:   false,
		},
		{
			name:        "Invalid host (should error)",
			url:         "http://10.255.255.1", // timeout IP
			requests:    5,
			concurrency: 1,
			expect200:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			success, errors, statusDist := simulateLoad(tt.url, tt.requests, tt.concurrency)

			if tt.expect200 && success == 0 {
				t.Errorf("Expected some successful requests, got 0")
			}
			if !tt.expect200 && success != 0 {
				t.Errorf("Expected no successful requests, got %d", success)
			}
			if len(statusDist) > 0 {
				t.Logf("Status distribution: %+v", statusDist)
			}
			if errors > 0 {
				t.Logf("Encountered %d errors", errors)
			}
		})
	}
}

// Helper to simulate load without exiting main
func simulateLoad(url string, totalReqs, concurrency int) (success int, errors int, statusDist map[int]int) {
	results := make(chan Result, totalReqs)
	sem := make(chan struct{}, concurrency)
	done := make(chan struct{})

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	go func() {
		for i := 0; i < totalReqs; i++ {
			sem <- struct{}{}
			go func() {
				defer func() { <-sem }()
				resp, err := client.Get(url)
				if err != nil {
					results <- Result{err: err}
				} else {
					results <- Result{status: resp.StatusCode}
					resp.Body.Close()
				}
			}()
		}
		// wait all goroutines finished
		for i := 0; i < cap(sem); i++ {
			sem <- struct{}{}
		}
		close(done)
	}()

	statusDist = make(map[int]int)

	go func() {
		<-done
		close(results)
	}()

	for res := range results {
		if res.err != nil {
			errors++
			continue
		}
		if res.status == 200 {
			success++
		} else {
			statusDist[res.status]++
		}
	}
	return
}
