package pokecache

import (
	"time"
	"sync"
	"fmt"
)

type Cache struct {
	mu			sync.Mutex // MAPS NOT THREAD SAFE, LOCK WHEN IN USE
    entries 	map[string]cacheEntry 
	interval	time.Duration
}

func (c *Cache) Write(key string, val cacheEntry) {
    c.mu.Lock()         // Acquire exclusive access
    defer c.mu.Unlock() // Ensure unlock happens at the end
    c.entries[key] = val
}

func (c *Cache) Read(key string) ([]byte, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    val, ok := c.entries[key]
	return val.val, ok
}

type cacheEntry struct {
	createdAt 	time.Time
	val 		[]byte
}

func NewCache(interval time.Duration) (*Cache) {
	cache := Cache{
    	entries:  make(map[string]cacheEntry),
    	interval: interval,
    	mu:      sync.Mutex{},
	}
	go cache.reapLoop()
	return &cache
}

func (c *Cache) Add(key string, val []byte) {
	entry := cacheEntry{
		createdAt: 		time.Now(),
		val:		 	val,
	}
	c.Write(key, entry)
}

func (c *Cache) Get(key string) ([]byte, bool) {
	val, ok := c.Read(key)
	return val, ok
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()   
		for key, val := range c.entries {
			cutoff := time.Now().Add(-c.interval) //The cutoff is "now minus the interval"
			if val.createdAt.Before(cutoff) {
				delete(c.entries, key)
				fmt.Printf("This cache entry missed the cutoff:%v it was deleted!!", key)
			}
		}
		c.mu.Unlock()
	}

}