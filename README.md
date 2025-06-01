# Cache Library for Go

This library provides efficient caching solutions for various usage patterns, such as write-heavy or read-heavy scenarios, and advanced features like expiration handling, integer-specific operations, and easy-to-use locking.

## Features

- **WriteHeavyCache**: Optimized for frequent write operations.
- **ReadHeavyCache**: Optimized for frequent read operations.
- **Expiration Support**: Built-in expiration for cache entries in `WriteHeavyCacheExpired` and `ReadHeavyCacheExpired`.
- **Integer-Specific Caches**: Specialized caches (`WriteHeavyCacheInteger` and `ReadHeavyCacheInteger`) with increment operations.
- **RollingCache**: A dynamically growing slice-based cache with efficient `Append` and `Rotate` operations, suitable for maintaining ordered sequences of elements.
- **Faster Singleflight**: A custom implementation that is up to **2x faster** than the standard Singleflight. See detailed results in the [benchmark](/benchmark) directory.
- **LockManager**: Simplified locking for managing concurrency in your applications.

## Documentation

For full API documentation, visit [pkg.go.dev](https://pkg.go.dev/github.com/catatsuy/cache).

## Installation

Install the library using `go get`:

```sh
go get github.com/catatsuy/cache
```

## Usage

### WriteHeavyCache

The `WriteHeavyCache` uses a `sync.Mutex` to lock for both read and write operations. This makes it suitable for scenarios with frequent writes.

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
		fmt.Println("Found:", value) // Output: Found: apple
	} else {
		fmt.Println("Not found")
	}

	// GetItems example
	items := c.GetItems()
	fmt.Println("Items:", items) // Output: Items: map[1:apple]

	// SetItems example
	c.SetItems(map[int]string{2: "banana", 3: "cherry"})
	fmt.Println("Items after SetItems:", c.GetItems()) // Output: Items after SetItems: map[2:banana 3:cherry]

	// Size example
	fmt.Println("Size:", c.Size()) // Output: Size: 2
}
```

### ReadHeavyCache

The `ReadHeavyCache` uses a `sync.RWMutex`, allowing multiple readers while locking only for writes. This makes it ideal for read-heavy scenarios.

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
		fmt.Println("Found:", value) // Output: Found: orange
	} else {
		fmt.Println("Not found")
	}

	// GetItems example
	items := c.GetItems()
	fmt.Println("Items:", items) // Output: Items: map[1:orange]

	// SetItems example
	c.SetItems(map[int]string{2: "peach", 3: "plum"})
	fmt.Println("Items after SetItems:", c.GetItems()) // Output: Items after SetItems: map[2:peach 3:plum]

	// Size example
	fmt.Println("Size:", c.Size()) // Output: Size: 2
}
```

### Expiration Support

The `WriteHeavyCacheExpired` and `ReadHeavyCacheExpired` caches provide expiration functionality, allowing you to specify a duration for each cache entry.

#### WriteHeavyCacheExpired Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/catatsuy/cache"
)

func main() {
	c := cache.NewWriteHeavyCacheExpired[int, string]()

	c.Set(1, "apple", 1*time.Second)
	fmt.Println("Before expiration:", c.Get(1)) // Output: Found: apple

	time.Sleep(2 * time.Second)
	_, found := c.Get(1)
	fmt.Println("After expiration:", found) // Output: After expiration: false
}
```

#### ReadHeavyCacheExpired Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/catatsuy/cache"
)

func main() {
	c := cache.NewReadHeavyCacheExpired[int, string]()

	c.Set(1, "orange", 1*time.Second)
	fmt.Println("Before expiration:", c.Get(1)) // Output: Found: orange

	time.Sleep(2 * time.Second)
	_, found := c.Get(1)
	fmt.Println("After expiration:", found) // Output: After expiration: false
}
```

### Integer-Specific Caches

For scenarios where increment operations are common, `WriteHeavyCacheInteger` and `ReadHeavyCacheInteger` are available.

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
	c.Incr(1, 10)

	value, _ := c.Get(1)
	fmt.Println("Incremented Value:", value) // Output: Incremented Value: 110
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
	c.Incr(1, 5)

	value, _ := c.Get(1)
	fmt.Println("Incremented Value:", value) // Output: Incremented Value: 55
}
```

### RollingCache

The `RollingCache` provides a simple, dynamically growing cache for ordered elements. It supports appending new items and rotating the cache, which resets its contents while returning the previous state.

```go
package main

import (
	"fmt"

	"github.com/catatsuy/cache"
)

func main() {
	// Create a RollingCache with an initial capacity of 10
	c := cache.NewRollingCache[int](10)
	c.Append(1)
	c.Append(2)
	c.Append(3)

	// Get the current items
	fmt.Println("Current items:", c.GetItems()) // Output: Current items: [1 2 3]

	// Check the size of the cache
	fmt.Println("Current size:", c.Size()) // Output: Current size: 3

	// Rotate the cache and get the rotated items
	rotated := c.Rotate()
	fmt.Println("Rotated items:", rotated) // Output: Rotated items: [1 2 3]

	// The cache should now be empty
	fmt.Println("Items after rotation:", c.GetItems()) // Output: Items after rotation: []
}
```

### LockManager

The `LockManager` is designed for managing locks associated with unique keys.

```go
package main

import (
	"fmt"
	"time"

	"github.com/catatsuy/cache"
)

func main() {
	lm := cache.NewLockManager[int]()

	lm.Lock(1)
	go func() {
		defer lm.GetAndLock(1).Unlock()
		fmt.Println("Goroutine locked and released")
	}()

	time.Sleep(1 * time.Second)
	lm.Unlock(1)
	fmt.Println("Main goroutine released lock")
}
```

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
	value, err, _ := sf.Do("key", func() (string, error) {
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
	v, err, _ := sf.Do(fmt.Sprintf("cacheGet_%d", key), func() (int, error) {
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
