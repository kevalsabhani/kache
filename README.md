# Kache

Kache is a generic, thread-safe key-value cache implementation in Go with expiration support. It provides a simple and efficient way to cache any type of key-value pairs with automatic cleanup of expired entries.

## Features

- 🎯 Generic implementation supporting any comparable key type and any value type
- ⏰ Automatic cleanup of expired items
- 🔒 Thread-safe operations through sync.RWMutex
- 🚀 Simple and intuitive API
- 💡 Zero external dependencies

## Installation

```bash
go get github.com/yourusername/kache
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "time"
    "github.com/yourusername/kache"
)

func main() {
    // Create a new cache instance with string keys and int values
    cache := kache.New[string, int]()

    // Set a value with 1 minute expiration
    cache.Set("counter", 42, 1*time.Minute)

    // Retrieve a value
    if value, exists := cache.Get("counter"); exists {
        fmt.Printf("Value: %d\n", value)
    }

    // Remove and get a value
    if value, exists := cache.Pop("counter"); exists {
        fmt.Printf("Popped value: %d\n", value)
    }
}
```

### Using Different Types

```go
// Cache with integer keys and custom struct values
type User struct {
    Name  string
    Email string
}

cache := kache.New[int, User]()
cache.Set(1, User{Name: "John", Email: "john@example.com"}, 30*time.Minute)
```

## API Reference

### `New[K comparable, V any]() *Kache[K, V]`
Creates a new cache instance with the specified key and value types.

### `Get(key K) (V, bool)`
Retrieves a value from the cache. Returns the value and a boolean indicating if the value was found and not expired.

### `Set(key K, value V, expiry time.Duration)`
Adds a value to the cache with an expiration duration. If no expiry is provided, the value will not expire.

### `Delete(key K)`
Removes a value from the cache.

### `Pop(key K) (V, bool)`
Removes and returns a value from the cache. Returns the value and a boolean indicating if the value was found and not expired.

## Features Details

- **Generic Implementation**: Use any comparable type as keys and any type as values
- **Automatic Cleanup**: Background goroutine removes expired items every 5 seconds
- **Thread-Safety**: All operations are protected by sync.RWMutex
- **Zero Dependencies**: Only uses Go standard library

## License

[Add your chosen license here]

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

