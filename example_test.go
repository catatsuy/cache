package cache_test

import (
	"fmt"
	"sync"
	"time"

	"github.com/catatsuy/cache"
)

// Example for WriteHeavyCache
func ExampleWriteHeavyCache() {
	c := cache.NewWriteHeavyCache[int, string]()

	c.Set(1, "apple")
	value, found := c.Get(1)

	if found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}
	// Output: Found: apple
}

// Example for WriteHeavyCache GetItems
func ExampleWriteHeavyCache_GetItems() {
	c := cache.NewWriteHeavyCache[int, string]()

	c.Set(1, "apple")
	c.Set(2, "banana")

	items := c.GetItems()
	fmt.Println("Items:", items)
	// Output: Items: map[1:apple 2:banana]
}

// Example for WriteHeavyCache SetItems
func ExampleWriteHeavyCache_SetItems() {
	c := cache.NewWriteHeavyCache[int, string]()

	c.SetItems(map[int]string{
		1: "grape",
		2: "cherry",
	})

	items := c.GetItems()
	fmt.Println("Items after SetItems:", items)
	// Output: Items after SetItems: map[1:grape 2:cherry]
}

// Example for WriteHeavyCache Size
func ExampleWriteHeavyCache_Size() {
	c := cache.NewWriteHeavyCache[int, string]()

	c.Set(1, "apple")
	c.Set(2, "banana")

	fmt.Println("Size:", c.Size())
	// Output: Size: 2
}

// Example for ReadHeavyCache
func ExampleReadHeavyCache() {
	c := cache.NewReadHeavyCache[int, string]()

	c.Set(1, "orange")
	value, found := c.Get(1)

	if found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}
	// Output: Found: orange
}

// Example for ReadHeavyCache GetItems
func ExampleReadHeavyCache_GetItems() {
	c := cache.NewReadHeavyCache[int, string]()

	c.Set(1, "orange")
	c.Set(2, "lemon")

	items := c.GetItems()
	fmt.Println("Items:", items)
	// Output: Items: map[1:orange 2:lemon]
}

// Example for ReadHeavyCache SetItems
func ExampleReadHeavyCache_SetItems() {
	c := cache.NewReadHeavyCache[int, string]()

	c.SetItems(map[int]string{
		1: "peach",
		2: "plum",
	})

	items := c.GetItems()
	fmt.Println("Items after SetItems:", items)
	// Output: Items after SetItems: map[1:peach 2:plum]
}

// Example for ReadHeavyCache Size
func ExampleReadHeavyCache_Size() {
	c := cache.NewReadHeavyCache[int, string]()

	c.Set(1, "orange")
	c.Set(2, "lemon")

	fmt.Println("Size:", c.Size())
	// Output: Size: 2
}

// Example for WriteHeavyCacheExpired
func ExampleWriteHeavyCacheExpired() {
	c := cache.NewWriteHeavyCacheExpired[int, string]()

	c.Set(1, "apple", 1*time.Second)

	if value, found := c.Get(1); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}

	// Expire immediately without waiting
	c.Set(1, "apple", -1*time.Second)
	if _, found := c.Get(1); !found {
		fmt.Println("Item has expired")
	}
	// Output:
	// Found: apple
	// Item has expired
}

// Example for ReadHeavyCacheExpired
func ExampleReadHeavyCacheExpired() {
	c := cache.NewReadHeavyCacheExpired[int, string]()

	c.Set(1, "orange", 1*time.Second)

	if value, found := c.Get(1); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}

	// Expire immediately without waiting
	c.Set(1, "orange", -1*time.Second)
	if _, found := c.Get(1); !found {
		fmt.Println("Item has expired")
	}
	// Output:
	// Found: orange
	// Item has expired
}

// Example for RollingCache Append and GetItems
func ExampleRollingCache() {
	c := cache.NewRollingCache[int](10)

	// Append values
	c.Append(1)
	c.Append(2)
	c.Append(3)

	// Get the items
	items := c.GetItems()
	fmt.Println("Items:", items)
	// Output: Items: [1 2 3]
}

// Example for RollingCache with dynamic growth
func ExampleRollingCache_dynamicGrowth() {
	c := cache.NewRollingCache[int](3)

	// Append more values than the initial length
	c.Append(1)
	c.Append(2)
	c.Append(3)
	c.Append(4) // This grows the slice beyond the initial length

	items := c.GetItems()
	fmt.Println("Items after appending more:", items)
	// Output: Items after appending more: [1 2 3 4]
}

// Example for RollingCache Rotate
func ExampleRollingCache_Rotate() {
	c := cache.NewRollingCache[int](10)

	// Append values
	c.Append(1)
	c.Append(2)
	c.Append(3)

	// Rotate the cache
	rotated := c.Rotate()
	fmt.Println("Rotated items:", rotated)

	// Cache should now be empty
	fmt.Println("Items after rotation:", c.GetItems())
	// Output:
	// Rotated items: [1 2 3]
	// Items after rotation: []
}

// Example for RollingCache Size
func ExampleRollingCache_Size() {
	c := cache.NewRollingCache[int](10)

	// Initially empty
	fmt.Println("Size initially:", c.Size())

	// Append values
	c.Append(1)
	c.Append(2)
	fmt.Println("Size after appending:", c.Size())

	// Append another value
	c.Append(3)
	fmt.Println("Size after appending more:", c.Size())

	// Output:
	// Size initially: 0
	// Size after appending: 2
	// Size after appending more: 3
}

// Example for LockManager with GetAndLock
func ExampleLockManager_GetAndLock() {
	var lm = cache.NewLockManager[int]()

	heavyOperation := func(id int) {
		defer lm.GetAndLock(id).Unlock()
		fmt.Printf("Starting heavy operation on resource %d\n", id)
		// simulate heavy work without slow test
		// time.Sleep(2 * time.Second)
		fmt.Printf("Completed heavy operation on resource %d\n", id)
	}

	heavyOperation(1)
	heavyOperation(2)
	// Output:
	// Starting heavy operation on resource 1
	// Completed heavy operation on resource 1
	// Starting heavy operation on resource 2
	// Completed heavy operation on resource 2
}

// Example for LockManager with Lock and Unlock
func ExampleLockManager_withLockAndUnlock() {
	lm := cache.NewLockManager[int]()

	lm.Lock(1)
	fmt.Println("Resource 1 is locked")

	fmt.Println("Resource 1 is being used")

	lm.Unlock(1)
	fmt.Println("Resource 1 is unlocked")
	// Output:
	// Resource 1 is locked
	// Resource 1 is being used
	// Resource 1 is unlocked
}

// ExampleLockManager provides an example usage of LockManager.
func ExampleLockManager() {
	// Create a new LockManager for integer keys
	lm := cache.NewLockManager[int]()
	var wg sync.WaitGroup
	firstDone := make(chan struct{})

	// Simulate concurrent access to the same key
	key := 1

	// Main goroutine does work first deterministically
	lm.Lock(key)
	fmt.Println("Locked")
	// Simulate some work
	fmt.Println("Doing work")
	lm.Unlock(key)
	fmt.Println("Unlocked")

	// First goroutine locks and performs some work
	wg.Go(func() {
		lm.Lock(key)
		fmt.Println("Goroutine 1: Locked")
		// Simulate some work
		fmt.Println("Goroutine 1: Doing work")
		lm.Unlock(key)
		fmt.Println("Goroutine 1: Unlocked")
		close(firstDone)
	})

	wg.Go(func() {
		<-firstDone // ensure goroutine 1 runs first
		defer lm.GetAndLock(key).Unlock()
		fmt.Println("Goroutine 2: Locked")
		// Simulate some work
		fmt.Println("Goroutine 2: Doing work")
		fmt.Println("Goroutine 2: Unlocked")
	})

	wg.Wait()
	// Output:
	// Locked
	// Doing work
	// Unlocked
	// Goroutine 1: Locked
	// Goroutine 1: Doing work
	// Goroutine 1: Unlocked
	// Goroutine 2: Locked
	// Goroutine 2: Doing work
	// Goroutine 2: Unlocked
}

func ExampleSingleflightGroup() {
	sf := cache.NewSingleflightGroup[string]()

	v, err := sf.Do("example_key", func() (string, error) {
		return "result", nil
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Value:", v)
	// Output: Value: result
}
