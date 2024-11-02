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

## API

### WriteHeavyCache

- **`Set(key K, value V)`**: Stores the given key-value pair in the cache.
- **`Get(key K) (V, bool)`**: Retrieves the value associated with the key, returning a boolean indicating whether the key exists.
- **`Clear()`**: Removes all key-value pairs from the cache.

### ReadHeavyCache

- **`Set(key K, value V)`**: Stores the given key-value pair in the cache.
- **`Get(key K) (V, bool)`**: Retrieves the value associated with the key, allowing concurrent reads.
- **`Clear()`**: Removes all key-value pairs from the cache.

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

## Acknowledgements

This library makes use of Go's powerful concurrency mechanisms (`sync.Mutex` and `sync.RWMutex`) to achieve thread-safe caching for various scenarios.
