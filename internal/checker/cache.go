package checker

import (
	"sync"
)

type Cache struct {
	mu   sync.RWMutex
	data map[string]HealthStatus
	once sync.Once
}

func NewCache() *Cache {
	c := &Cache{}

	c.once.Do(func() {
		c.data = make(map[string]HealthStatus)
	})

	return c
}

func (c *Cache) Set(url string, status HealthStatus) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[url] = status
}

func (c *Cache) Get(url string) (HealthStatus, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	status, ok := c.data[url]
	return status, ok
}

func (c *Cache) GetAll() map[string]HealthStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]HealthStatus, len(c.data))
	for k, v := range c.data {
		result[k] = v
	}

	return result
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]HealthStatus)
}
