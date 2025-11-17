package checker

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type HealthStatus struct {
	URL          string
	Status       string
	ResponseTime time.Duration
	LastCheck    time.Time
	Error        string
}

type HealthChecker struct {
	urls          []string
	cache         *Cache
	checkInterval time.Duration
	timeout       time.Duration
	mu            sync.RWMutex
	client        *http.Client
	workers       int
}

func New(checkIntercal, timeout time.Duration) *HealthChecker {
	return &HealthChecker{
		urls:          make([]string, 0),
		cache:         NewCache(),
		checkInterval: checkIntercal,
		timeout:       timeout,
		workers:       5,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (hc *HealthChecker) AddURL(url string) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.urls = append(hc.urls, url)
}

func (hc *HealthChecker) Start(ctx context.Context) {
	urlChan := make(chan string, len(hc.urls))
	resultChan := make(chan HealthStatus, len(hc.urls))

	var wg sync.WaitGroup
	for i := 0; i < hc.workers; i++ {
		wg.Add(1)
		go hc.worker(ctx, &wg, urlChan, resultChan)
	}

	go hc.collectResults(ctx, resultChan)

	go func() {
		ticker := time.NewTicker(hc.checkInterval)
		defer ticker.Stop()

		hc.scheduleChecks(urlChan)
		for {
			select {
			case <-ctx.Done():
				close(urlChan)
				wg.Wait()
				close(resultChan)
				return

			case <-ticker.C:
				hc.scheduleChecks(urlChan)
			}
		}
	}()
}

func (hc *HealthChecker) worker(ctx context.Context, wg *sync.WaitGroup, urlChan <-chan string, resultChan chan<- HealthStatus) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case url, ok := <-urlChan:
			if !ok {
				return
			}
			status := hc.checkURL(url)

			select {
			case resultChan <- status:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (hc *HealthChecker) checkURL(url string) HealthStatus {
	start := time.Now()
	status := HealthStatus{
		URL:       url,
		LastCheck: start,
	}

	resp, err := hc.client.Get(url)
	if err != nil {
		status.Status = "down"
		status.Error = err.Error()
		return status
	}

	defer resp.Body.Close()
	status.ResponseTime = time.Since(start)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		status.Status = "up"
	} else {
		status.Status = "degraded"
		status.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}

	return status
}

func (hc *HealthChecker) scheduleChecks(urlChan chan<- string) {
	hc.mu.RLock()
	urls := make([]string, len(hc.urls))
	copy(urls, hc.urls)
	hc.mu.RUnlock()

	for _, url := range urls {
		urlChan <- url
	}
}

func (hc *HealthChecker) collectResults(ctx context.Context, resultChan <-chan HealthStatus) {
	for {
		select {
		case <-ctx.Done():
			return
		case status, ok := <-resultChan:
			if !ok {
				return
			}
			hc.cache.Set(status.URL, status)
		}
	}
}

func (hc *HealthChecker) GetStatus(url string) (HealthStatus, bool) {
	return hc.cache.Get(url)
}

func (hc *HealthChecker) GetAllStatuses() []HealthStatus {
	hc.mu.RLock()
	urls := make([]string, len(hc.urls))
	copy(urls, hc.urls)
	hc.mu.RUnlock()

	statuses := make([]HealthStatus, 0, len(urls))
	for _, url := range urls {
		if status, ok := hc.cache.Get(url); ok {
			statuses = append(statuses, status)
		}
	}

	return statuses
}
