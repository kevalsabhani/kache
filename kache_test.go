package kache

import (
	"sync"
	"testing"
	"time"
)

func TestKacheBasicOperations(t *testing.T) {
	t.Run("Set and Get operations", func(t *testing.T) {
		cache := New[string, int]()

		// Test Set and Get
		cache.Set("key1", 42, 1*time.Minute)
		value, exists := cache.Get("key1")

		if !exists {
			t.Error("Expected key to exist")
		}
		if value != 42 {
			t.Errorf("Expected value 42, got %d", value)
		}

		// Test non-existent key
		_, exists = cache.Get("nonexistent")
		if exists {
			t.Error("Expected key to not exist")
		}
	})

	t.Run("Delete operation", func(t *testing.T) {
		cache := New[string, int]()

		cache.Set("key1", 42, 1*time.Minute)
		cache.Delete("key1")

		_, exists := cache.Get("key1")
		if exists {
			t.Error("Key should have been deleted")
		}

		// Delete non-existent key should not panic
		cache.Delete("nonexistent")
	})

	t.Run("Pop operation", func(t *testing.T) {
		cache := New[string, int]()

		cache.Set("key1", 42, 1*time.Minute)
		value, exists := cache.Pop("key1")

		if !exists {
			t.Error("Expected key to exist")
		}
		if value != 42 {
			t.Errorf("Expected value 42, got %d", value)
		}

		// Verify key was removed
		_, exists = cache.Get("key1")
		if exists {
			t.Error("Key should have been removed after Pop")
		}

		// Pop non-existent key
		_, exists = cache.Pop("nonexistent")
		if exists {
			t.Error("Expected Pop on non-existent key to return false")
		}
	})
}

func TestExpiration(t *testing.T) {
	t.Run("Item expiration", func(t *testing.T) {
		cache := New[string, int]()

		// Set item with very short expiration
		cache.Set("key1", 42, 50*time.Millisecond)

		// Verify item exists initially
		_, exists := cache.Get("key1")
		if !exists {
			t.Error("Item should exist before expiration")
		}

		// Wait for expiration
		time.Sleep(100 * time.Millisecond)

		// Verify item has expired
		_, exists = cache.Get("key1")
		if exists {
			t.Error("Item should have expired")
		}
	})

	t.Run("Automatic cleanup", func(t *testing.T) {
		cache := New[string, int]()

		// Set multiple items with short expiration
		cache.Set("key1", 1, 50*time.Millisecond)
		cache.Set("key2", 2, 50*time.Millisecond)

		// Wait for cleanup cycle (>5 seconds)
		time.Sleep(6 * time.Second)

		// Verify items were cleaned up
		if len(cache.data) != 0 {
			t.Error("Expected all items to be cleaned up")
		}
	})

	t.Run("Zero expiration", func(t *testing.T) {
		cache := New[string, int]()

		// Set item with zero expiration
		cache.Set("key1", 42, 0)

		// Wait some time
		time.Sleep(100 * time.Millisecond)

		// Verify item still exists
		_, exists := cache.Get("key1")
		if !exists {
			t.Error("Item with zero expiration should not expire")
		}
	})

	t.Run("Pop with expired item", func(t *testing.T) {
		cache := New[string, int]()

		cache.Set("key1", 42, 50*time.Millisecond)
		time.Sleep(100 * time.Millisecond)

		value, exists := cache.Pop("key1")
		if exists {
			t.Error("Expected item to be expired")
		}
		if value != 0 {
			t.Errorf("Expected value 0, got %d", value)
		}
	})
}

func TestConcurrency(t *testing.T) {
	t.Run("Concurrent access", func(t *testing.T) {
		cache := New[int, int]()
		var wg sync.WaitGroup

		// Number of concurrent goroutines
		workers := 100
		operations := 1000

		// Start multiple goroutines performing operations
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				for j := 0; j < operations; j++ {
					key := workerID*operations + j

					// Perform random operations
					switch j % 4 {
					case 0:
						cache.Set(key, key, 1*time.Minute)
					case 1:
						cache.Get(key)
					case 2:
						cache.Delete(key)
					case 3:
						cache.Pop(key)
					}
				}
			}(i)
		}

		wg.Wait()
	})
}

func TestDifferentTypes(t *testing.T) {
	t.Run("String keys and struct values", func(t *testing.T) {
		type User struct {
			Name  string
			Email string
		}

		cache := New[string, User]()
		user := User{Name: "John", Email: "john@example.com"}

		cache.Set("user1", user, 1*time.Minute)

		retrieved, exists := cache.Get("user1")
		if !exists {
			t.Error("Expected user to exist")
		}
		if retrieved.Name != user.Name || retrieved.Email != user.Email {
			t.Error("Retrieved user does not match original")
		}
	})

	t.Run("Integer keys and interface values", func(t *testing.T) {
		cache := New[int, interface{}]()

		cache.Set(1, "string value", 1*time.Minute)
		cache.Set(2, 42, 1*time.Minute)
		cache.Set(3, true, 1*time.Minute)

		value, exists := cache.Get(1)
		if !exists || value.(string) != "string value" {
			t.Error("Wrong value for string")
		}

		value, exists = cache.Get(2)
		if !exists || value.(int) != 42 {
			t.Error("Wrong value for int")
		}

		value, exists = cache.Get(3)
		if !exists || value.(bool) != true {
			t.Error("Wrong value for bool")
		}
	})
}

func TestItemExpiry(t *testing.T) {
	t.Run("isExpired method", func(t *testing.T) {
		// Test expired item
		expiredItem := item[string]{
			value:  "test",
			expiry: time.Now().Add(-1 * time.Minute).UnixNano(),
		}
		if !expiredItem.isExpired() {
			t.Error("Item should be expired")
		}

		// Test non-expired item
		validItem := item[string]{
			value:  "test",
			expiry: time.Now().Add(1 * time.Minute).UnixNano(),
		}
		if validItem.isExpired() {
			t.Error("Item should not be expired")
		}
	})
}
