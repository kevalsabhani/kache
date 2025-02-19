// Package kache provides a generic, thread-safe key-value cache implementation with expiration support.
//
// The package offers a simple interface (Kacher) and implementation (Kache) for caching arbitrary
// key-value pairs with automatic cleanup of expired entries. The cache is safe for concurrent use
// through sync.RWMutex protection.
//
// Example usage:
//
//	cache := kache.New[string, int]()
//	cache.Set("key", 42)
//	value, exists := cache.Get("key")
//
// Key features:
//   - Generic implementation supporting any comparable key type and any value type
//   - Automatic cleanup of expired items
//   - Thread-safe operations
//   - Basic operations: Get, Set, Delete, Pop

package kache

import (
	"sync"
	"time"
)

// Kacher is a generic key-value cache interface that supports basic operations
// like Get, Set, Delete and Pop with keys of type K and values of type V.
type Kacher[K comparable, V any] interface {
	Get(K) (V, bool)
	Set(K, V)
	Delete(K)
	Pop(K) (V, bool)
}

// item represents a cached value with an expiration time.
// V is the type of the value being stored.
type item[V any] struct {
	value  V
	expiry time.Time
}

// isExpired checks if the item has expired by comparing its expiry time
// with the current time. Returns true if the item has expired.
func (i item[V]) isExpired() bool {
	return i.expiry.After(time.Now())
}

// Kache is a generic key-value cache implementation that supports basic operations
// like Get, Set, Delete and Pop with keys of type K and values of type V.
type Kache[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]item[V]
}

// New creates a new Kache instance.
func New[K comparable, V any]() *Kache[K, V] {
	c := &Kache[K, V]{
		data: make(map[K]item[V]),
	}

	// Start a goroutine to periodically clean up expired items.
	go func() {
		for range time.Tick(5 * time.Second) {
			c.mu.Lock()
			for k, v := range c.data {
				if v.isExpired() {
					delete(c.data, k)
				}
			}
			c.mu.Unlock()
		}
	}()

	return c
}

// Get retrieves a value from the cache.
// Returns the value and a boolean indicating if the value was found.
func (c *Kache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.data[key]

	if !found {
		return item.value, false
	}

	if item.isExpired() {
		delete(c.data, key)
		return item.value, false
	}

	return item.value, found
}

// Set adds a value to the cache with an optional expiry time.
// If no expiry is provided, the value will not expire.
func (c *Kache[K, V]) Set(key K, value V, expiry time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = item[V]{
		value:  value,
		expiry: time.Now().Add(expiry),
	}
}

// Delete removes a value from the cache.
func (c *Kache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// Pop removes a value from the cache and returns it.
func (c *Kache[K, V]) Pop(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, found := c.data[key]

	if !found {
		return item.value, false
	}

	delete(c.data, key)

	if item.isExpired() {
		return item.value, false
	}

	return item.value, found
}
