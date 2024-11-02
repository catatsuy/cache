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

// Example for WriteHeavyCacheExpired
func ExampleWriteHeavyCacheExpired() {
	c := cache.NewWriteHeavyCacheExpired[int, string]()

	c.Set(1, "apple", 1*time.Second)

	if value, found := c.Get(1); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}

	time.Sleep(2 * time.Second)
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

	time.Sleep(2 * time.Second)
	if _, found := c.Get(1); !found {
		fmt.Println("Item has expired")
	}
	// Output:
	// Found: orange
	// Item has expired
}

// Example for LockManager with GetAndLock
func ExampleLockManager_GetAndLock() {
	var lm = cache.NewLockManager[int]()

	heavyOperation := func(id int) {
		defer lm.GetAndLock(id).Unlock()
		fmt.Printf("Starting heavy operation on resource %d\n", id)
		time.Sleep(2 * time.Second)
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
	wg.Add(3)

	// Simulate concurrent access to the same key
	key := 1

	// First goroutine locks and performs some work
	go func() {
		defer wg.Done()
		lm.Lock(key)
		fmt.Println("Goroutine 1: Locked")
		// Simulate some work
		fmt.Println("Goroutine 1: Doing work")
		lm.Unlock(key)
		fmt.Println("Goroutine 1: Unlocked")
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond)
		defer lm.GetAndLock(key).Unlock()
		fmt.Println("Goroutine 2: Locked")
		// Simulate some work
		fmt.Println("Goroutine 2: Doing work")
		fmt.Println("Goroutine 2: Unlocked")
	}()

	lm.Lock(key)
	fmt.Println("Locked")
	// Simulate some work
	fmt.Println("Doing work")
	lm.Unlock(key)
	fmt.Println("Unlocked")
	wg.Done()
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
