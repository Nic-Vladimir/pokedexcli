package internal

import (
	"fmt"
	"sync"
	"time"
)

// --- Data structures ---
type Cache struct {
	entries map[string]cacheEntry
	mutex   sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// --- Functions ---
func NewCache(cleanTime time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]cacheEntry),
	}
	go c.reapLoop(cleanTime)
	return c
}

// Cache.Add() method
func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := c.entries[key]; !ok {
		return []byte{}, false
	}
	return c.entries[key].val, true
}

func (c *Cache) reapLoop(cleanTime time.Duration) {
	ticker := time.NewTicker(cleanTime)
	defer ticker.Stop()
	for range ticker.C {
		c.mutex.Lock()
		for k, entry := range c.entries {
			if time.Since(entry.createdAt) > cleanTime {
				delete(c.entries, k)
				fmt.Println("\nExpired cache: ", k)
			}
		}
		c.mutex.Unlock()
	}
}
