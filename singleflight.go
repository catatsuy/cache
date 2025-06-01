package cache

import (
	"sync"
)

// SingleflightGroup manages single concurrent requests per key, ensuring that
// only one execution of a function occurs for a given key at a time.
//
// This implementation is simplified compared to the official singleflight package
// and lacks advanced error handling and other features, such as:
//   - Panic and runtime.Goexit handling: This implementation does not handle cases
//     where the function panics or terminates abnormally. In the official implementation,
//     errors from panic and Goexit are handled to prevent blocked goroutines.
//   - Immediate synchronous cleanup: In this implementation, the completed result is
//     removed from the map asynchronously. In the official implementation, cleanup
//     is handled synchronously within the doCall function to ensure immediate memory release.
type SingleflightGroup[V any] struct {
	mu sync.Mutex
	m  map[string]*call[V]
}

// call represents a single execution result for a specific key, holding the
// value, any error encountered, and whether the execution is completed.
type call[V any] struct {
	mu    sync.Mutex
	value V
	err   error
	done  bool
}

// NewSingleflightGroup creates a new instance of SingleflightGroup, initialized
// with an empty map to store calls by key.
func NewSingleflightGroup[V any]() *SingleflightGroup[V] {
	return &SingleflightGroup[V]{
		m: make(map[string]*call[V]),
	}
}

// Do ensures that for a given key, only one execution of fn occurs at a time.
// If a call for the key is already in progress, other calls wait for its completion
// and return the same result. Once complete, the result is stored and used for
// subsequent calls until it's removed from the map.
//
// Unlike the official singleflight, this function does not provide:
//   - Panic and Goexit handling
func (sf *SingleflightGroup[V]) Do(key string, fn func() (V, error)) (v V, err error, shared bool) {
	// Lock to check if a call is already in progress for the given key
	sf.mu.Lock()
	c, ok := sf.m[key]
	if !ok {
		// If no call exists for the key, create a new one
		c = &call[V]{}
		sf.m[key] = c
	}
	sf.mu.Unlock()

	// Lock the call to ensure only one execution of fn
	c.mu.Lock()
	if !c.done {
		// If fn has not been executed, run it and store the result
		c.value, c.err = fn()
		c.done = true

		// Schedule the deletion of the completed call asynchronously
		go func() {
			sf.mu.Lock()
			delete(sf.m, key)
			sf.mu.Unlock()
		}()

		c.mu.Unlock()

		return c.value, c.err, false
	}
	c.mu.Unlock()

	return c.value, c.err, true
}
