# Cache Library for Go

This library provides two types of caches optimized for different usage patterns: **WriteHeavyCache** and **ReadHeavyCache**. Additionally, it supports caches specifically designed for integer-like types, offering flexibility for handling various data types.

## Features

- **WriteHeavyCache**: Optimized for scenarios where write operations are frequent.
- **ReadHeavyCache**: Optimized for scenarios where read operations are frequent.
- Supports **generic types** for flexible key-value storage.
- Specialized caches for **integer-like types** to support operations like increments.

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

## API

### WriteHeavyCache

- **`Set(key K, value V)`**: Stores the given key-value pair in the cache.
- **`Get(key K) (V, bool)`**: Retrieves the value associated with the key, returning a boolean indicating whether the key exists.

### ReadHeavyCache

- **`Set(key K, value V)`**: Stores the given key-value pair in the cache.
- **`Get(key K) (V, bool)`**: Retrieves the value associated with the key, allowing concurrent reads.

### WriteHeavyCacheInteger

- **`Set(key K, value V)`**: Stores the given key-value pair in the cache.
- **`Get(key K) (V, bool)`**: Retrieves the value associated with the key.
- **`Incr(key K, value V)`**: Increments the value by the given amount. If the key does not exist, it sets the value.

### ReadHeavyCacheInteger

- **`Set(key K, value V)`**: Stores the given key-value pair in the cache.
- **`Get(key K) (V, bool)`**: Retrieves the value associated with the key.
- **`Incr(key K, value V)`**: Increments the value by the given amount. If the key does not exist, it sets the value.

## Acknowledgements

This library makes use of Go's powerful concurrency mechanisms (`sync.Mutex` and `sync.RWMutex`) to achieve thread-safe caching for various scenarios.
