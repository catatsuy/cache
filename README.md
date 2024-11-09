# Cache Library for Go

This library provides two types of caches optimized for different usage patterns: **WriteHeavyCache** and **ReadHeavyCache**. Additionally, it supports caches specifically designed for integer-like types, offering flexibility for handling various data types.

## Features

- **WriteHeavyCache**: Optimized for scenarios where write operations are frequent.
- **ReadHeavyCache**: Optimized for scenarios where read operations are frequent.
- Supports **generic types** for flexible key-value storage.
- Specialized caches for **integer-like types** to support operations like increments.

## Documentation

For full API documentation, visit [pkg.go.dev](https://pkg.go.dev/github.com/catatsuy/cache).

## Usage

### WriteHeavyCache

The `WriteHeavyCache` uses a `sync.Mutex` to lock for both read and write operations, making it suitable for cases where writing to the cache is frequent.

```go
package main

import (
	"fmt"

	"github.com/catatsuy/cache"
)

func main() {
	c := cache.NewWriteHeavyCache[int, string]()

	c.Set(1, "apple")
	value, found := c.Get(1)

	if found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}
}
```

### ReadHeavyCache

The `ReadHeavyCache` uses a `sync.RWMutex`, allowing multiple concurrent readers while locking only for write operations. This is ideal for read-heavy workloads.

```go
package main

import (
	"fmt"

	"github.com/catatsuy/cache"
)

func main() {
	c := cache.NewReadHeavyCache[int, string]()

	c.Set(1, "orange")
	value, found := c.Get(1)

	if found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}
}
```

### WriteHeavyCacheExpired and ReadHeavyCacheExpired

The `WriteHeavyCacheExpired` and `ReadHeavyCacheExpired` caches provide expiration functionality for stored items. You can specify an expiration duration for each item, after which it will no longer be accessible.

- **WriteHeavyCacheExpired**: Optimized for write-heavy scenarios, using `sync.Mutex` for all operations.
- **ReadHeavyCacheExpired**: Optimized for read-heavy scenarios, using `sync.RWMutex` to allow multiple readers concurrently while still locking for writes.

#### WriteHeavyCacheExpired Example

The `WriteHeavyCacheExpired` cache is designed for situations where write operations are more frequent.

```go
package main

import (
	"fmt"
	"time"

	"github.com/catatsuy/cache"
)

func main() {
	c := cache.NewWriteHeavyCacheExpired[int, string]()

	// Set an item with a 1-second expiration
	c.Set(1, "apple", 1*time.Second)

	// Retrieve the item immediately
	if value, found := c.Get(1); found {
		fmt.Println("Found:", value) // Output: Found: apple
	} else {
		fmt.Println("Not found")
	}

	// Wait for the item to expire
	time.Sleep(2 * time.Second)
	if _, found := c.Get(1); !found {
		fmt.Println("Item has expired") // Output: Item has expired
	}
}
```

#### ReadHeavyCacheExpired Example

The `ReadHeavyCacheExpired` cache is designed for scenarios where read operations are more frequent than writes.

```go
package main

import (
	"fmt"
	"time"

	"github.com/catatsuy/cache"
)

func main() {
	c := cache.NewReadHeavyCacheExpired[int, string]()

	// Set an item with a 1-second expiration
	c.Set(1, "orange", 1*time.Second)

	// Retrieve the item immediately
	if value, found := c.Get(1); found {
		fmt.Println("Found:", value) // Output: Found: orange
	} else {
		fmt.Println("Not found")
	}

	// Wait for the item to expire
	time.Sleep(2 * time.Second)
	if _, found := c.Get(1); !found {
		fmt.Println("Item has expired") // Output: Item has expired
	}
}
```

### Integer-Specific Caches

For scenarios where you need to increment values stored in the cache, the library provides `WriteHeavyCacheInteger` and `ReadHeavyCacheInteger` for integer-like types.

#### WriteHeavyCacheInteger Example

```go
package main

import (
	"fmt"

	"github.com/catatsuy/cache"
)

func main() {
	c := cache.NewWriteHeavyCacheInteger[int, int]()

	c.Set(1, 100)
	c.Incr(1, 10) // Increment the value by 10

	value, found := c.Get(1)

	if found {
		fmt.Println("New Value:", value) // Output: New Value: 110
	} else {
		fmt.Println("Not found")
	}
}
```

#### ReadHeavyCacheInteger Example

```go
package main

import (
	"fmt"

	"github.com/catatsuy/cache"
)

func main() {
	c := cache.NewReadHeavyCacheInteger[int, int]()

	c.Set(1, 50)
	c.Incr(1, 5) // Increment the value by 5

	value, found := c.Get(1)

	if found {
		fmt.Println("New Value:", value) // Output: New Value: 55
	} else {
		fmt.Println("Not found")
	}
}
```

### LockManager

The `LockManager` is useful for managing distributed locks associated with unique keys, especially when multiple goroutines need synchronized access to specific resources. `LockManager` provides a convenient `GetAndLock` function, which retrieves and locks a mutex in one step.

#### LockManager Example with `GetAndLock`

```go
package main

import (
	"fmt"
	"github.com/catatsuy/cache"
	"time"
)

// Create a global LockManager instance
var lm = cache.NewLockManager[int]()

func main() {
	// Simulate calling heavyOperation on a shared resource
	heavyOperation(1)
	heavyOperation(2)
}

func heavyOperation(id int) {
	// Lock the resource with the specified key and defer unlocking
	defer lm.GetAndLock(id).Unlock()

	// Simulate a time-consuming process
	fmt.Printf("Starting heavy operation on resource %d\n", id)
	time.Sleep(2 * time.Second)
	fmt.Printf("Completed heavy operation on resource %d\n", id)
}
```

The `GetAndLock` function can be used with `defer` to ensure the mutex is automatically unlocked at the end of the function.

#### LockManager Example with `Lock` and `Unlock`

You can also explicitly call `Lock` and `Unlock` for more control over the locking process.

```go
package main

import (
	"fmt"
	"github.com/catatsuy/cache"
)

func main() {
	lm := cache.NewLockManager[int]()

	// Lock the resource with a specific key
	lm.Lock(1)
	fmt.Println("Resource 1 is locked")

	// Perform some work with the locked resource
	fmt.Println("Resource 1 is being used")

	// Unlock the resource
	lm.Unlock(1)
	fmt.Println("Resource 1 is unlocked")
}
```

Using `Lock` and `Unlock` directly allows you to control when the lock is held and released, making it suitable for cases where `defer` might not be appropriate.

### SingleflightGroup

`SingleflightGroup` ensures that only one function call happens at a time for a given key. If multiple calls are made with the same key, only the first one runs, while others wait and receive the same result. This is useful when you want to avoid running the same operation multiple times simultaneously.

#### SingleflightGroup Example

```go
package main

import (
	"fmt"

	"github.com/catatsuy/cache"
)

func main() {
	sf := cache.NewSingleflightGroup[string]()

	// Define a function to load data only if it's not already cached
	loadData := func(key string) (string, error) {
		// Simulate data fetching or updating
		return fmt.Sprintf("Data for key %s", key), nil
	}

	// Use SingleflightGroup to ensure only one call for the same key at a time
	value, err := sf.Do("key", func() (string, error) {
		return loadData("key")
	})

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", value) // Output: Result: Data for key 1
	}
}
```

#### SingleflightGroup and Caching Example

This example demonstrates how to use `SingleflightGroup` with `WriteHeavyCache` to prevent duplicate data retrieval requests for the same key. When a key is requested multiple times simultaneously, `SingleflightGroup` ensures that only one retrieval function (`HeavyGet`) executes, while the other requests wait and receive the cached result once it completes.

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/catatsuy/cache"
)

// Global cache and singleflight group
var (
	c  = cache.NewWriteHeavyCache[int, int]()
	sf = cache.NewSingleflightGroup[int]()
)

// Get retrieves a value from the cache or loads it with HeavyGet if not cached.
// Singleflight ensures only one HeavyGet call per key at a time.
func Get(key int) int {
	// Attempt to retrieve the item from cache
	if value, found := c.Get(key); found {
		return value
	}

	// Use SingleflightGroup to prevent duplicate HeavyGet calls for the same key
	v, err := sf.Do(fmt.Sprintf("cacheGet_%d", key), func() (int, error) {
		// Load the data and store it in the cache
		value := HeavyGet(key)
		c.Set(key, value)
		return value, nil
	})
	if err != nil {
		panic(err)
	}

	return v
}

// HeavyGet simulates a time-consuming data retrieval operation.
// Here, it sleeps for 1 second and returns twice the input key as the result.
func HeavyGet(key int) int {
	log.Printf("call HeavyGet %d\n", key)
	time.Sleep(time.Second)
	return key * 2
}

func main() {
	// Simulate concurrent access to Get with keys 0 through 9, each key accessed 10 times
	for i := 0; i < 100; i++ {
		go func(i int) {
			Get(i % 10)
		}(i)
	}

	// Wait for all concurrent fetches to complete
	time.Sleep(2 * time.Second)

	// Print cached values for keys 0 through 9
	for i := 0; i < 10; i++ {
		log.Println(Get(i))
	}
}
```

This example shows how multiple simultaneous requests for the same key result in only a single call to `HeavyGet`, while other requests wait for the result. The results are then cached, preventing repeated retrievals for the same data.

## API

### WriteHeavyCache

- **`Set(key K, value V)`**: Stores the given key-value pair in the cache.
- **`Get(key K) (V, bool)`**: Retrieves the value associated with the key, returning a boolean indicating whether the key exists.
- **`Clear()`**: Removes all key-value pairs from the cache.

### ReadHeavyCache

- **`Set(key K, value V)`**: Stores the given key-value pair in the cache.
- **`Get(key K) (V, bool)`**: Retrieves the value associated with the key, allowing concurrent reads.
- **`Clear()`**: Removes all key-value pairs from the cache.

### WriteHeavyCacheExpired

- **`Set(key K, value V, duration time.Duration)`**: Stores the given key-value pair in the cache with an expiration duration.
- **`Get(key K) (V, bool)`**: Retrieves the value associated with the key if it exists and is not expired. Returns a boolean indicating whether the key is still valid.

### ReadHeavyCacheExpired

- **`Set(key K, value V, duration time.Duration)`**: Stores the given key-value pair in the cache with an expiration duration.
- **`Get(key K) (V, bool)`**: Retrieves the value associated with the key if it exists and is not expired. Returns a boolean indicating whether the key is still valid.

### WriteHeavyCacheInteger

- **`Set(key K, value V)`**: Stores the given key-value pair in the cache.
- **`Get(key K) (V, bool)`**: Retrieves the value associated with the key.
- **`Incr(key K, value V)`**: Increments the value by the given amount. If the key does not exist, it sets the value.
- **`Clear()`**: Removes all key-value pairs from the cache.

### ReadHeavyCacheInteger

- **`Set(key K, value V)`**: Stores the given key-value pair in the cache.
- **`Get(key K) (V, bool)`**: Retrieves the value associated with the key.
- **`Incr(key K, value V)`**: Increments the value by the given amount. If the key does not exist, it sets the value.
- **`Clear()`**: Removes all key-value pairs from the cache.

### LockManager

- **`NewLockManager[K comparable]()`**: Creates a new instance of `LockManager`.
- **`Lock(id K)`**: Locks the mutex associated with the given key.
- **`Unlock(id K)`**: Unlocks the mutex associated with the given key.
- **`GetAndLock(id K) *sync.Mutex`**: Retrieves and locks the mutex associated with the given key, returning the locked mutex. Useful with `defer` to automatically release the lock when the function exits.

### SingleflightGroup

- **`NewSingleflightGroup[V any]()`**: Creates a new instance of `SingleflightGroup`.
- **`Do(key string, fn func() (V, error)) (V, error)`**: Executes the provided function `fn` for the given key only if it is not already in progress for that key. If a duplicate request is made with the same key while the function is still running, the duplicate request waits and receives the same result when the function completes.

## Acknowledgements

This library makes use of Go's powerful concurrency mechanisms (`sync.Mutex` and `sync.RWMutex`) to achieve thread-safe caching for various scenarios.
