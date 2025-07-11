package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mu      sync.Mutex
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	newEntry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.entries[key] = newEntry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	content := make([]byte, 0)
	value, keyExists := c.entries[key]

	if !keyExists {
		return content, keyExists
	}

	return value.val, keyExists
}

func (c *Cache) readLoop(duration time.Duration) {
	ticker := time.NewTicker(duration)

	for range ticker.C {

		func() {
			c.mu.Lock()
			defer c.mu.Unlock()
			now := time.Now()
			for key, value := range c.entries {
				elapsedTime := now.Sub(value.createdAt)
				if elapsedTime > duration {
					delete(c.entries, key)
				}
			}
		}()
	}

}

func NewCache(interval time.Duration) *Cache {

	cacheInstance := &Cache{
		mu:      sync.Mutex{},
		entries: make(map[string]cacheEntry),
	}

	go cacheInstance.readLoop(interval)

	return cacheInstance
}
