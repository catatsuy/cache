package cache

import (
	"sync"
	"time"
)

// WriteHeavyCache is a cache optimized for write-heavy operations.
// It uses a Mutex to synchronize access to the cache items.
type WriteHeavyCache[K comparable, V any] struct {
	sync.Mutex // WriteHeavyCache uses Mutex for all operations
	items      map[K]V
}

// ReadHeavyCache is a cache optimized for read-heavy operations.
// It uses an RWMutex to allow concurrent reads and synchronized writes.
type ReadHeavyCache[K comparable, V any] struct {
	sync.RWMutex // ReadHeavyCache allows concurrent read access with RWMutex
	items        map[K]V
}

// Set sets a value in WriteHeavyCache, locking for the write operation
func (c *WriteHeavyCache[K, V]) Set(key K, value V) {
	c.Lock()
	c.items[key] = value
	c.Unlock()
}

// Get retrieves a value from WriteHeavyCache, locking for read as well
func (c *WriteHeavyCache[K, V]) Get(key K) (V, bool) {
	c.Lock()
	v, found := c.items[key]
	c.Unlock()
	return v, found
}

// Clear removes all items from WriteHeavyCache
func (c *WriteHeavyCache[K, V]) Clear() {
	c.Lock()
	c.items = make(map[K]V)
	c.Unlock()
}

// Set sets a value in ReadHeavyCache, locking for the write operation
func (c *ReadHeavyCache[K, V]) Set(key K, value V) {
	c.Lock()
	c.items[key] = value
	c.Unlock()
}

// Get retrieves a value from ReadHeavyCache, using a read lock
func (c *ReadHeavyCache[K, V]) Get(key K) (V, bool) {
	c.RLock()
	v, found := c.items[key]
	c.RUnlock()
	return v, found
}

// Clear removes all items from ReadHeavyCache
func (c *ReadHeavyCache[K, V]) Clear() {
	c.Lock()
	c.items = make(map[K]V)
	c.Unlock()
}

// NewWriteHeavyCache creates a new instance of WriteHeavyCache
func NewWriteHeavyCache[K comparable, V any]() *WriteHeavyCache[K, V] {
	return &WriteHeavyCache[K, V]{
		items: make(map[K]V),
	}
}

// NewReadHeavyCache creates a new instance of ReadHeavyCache
func NewReadHeavyCache[K comparable, V any]() *ReadHeavyCache[K, V] {
	return &ReadHeavyCache[K, V]{
		items: make(map[K]V),
	}
}

// expiredValue represents a cached value with an expiration time.
type expiredValue[V any] struct {
	value  V
	expire time.Time
}

// WriteHeavyCacheExpired is a cache optimized for write-heavy operations with expiration support.
// It uses a Mutex to synchronize access and stores values with expiration times.
type WriteHeavyCacheExpired[K comparable, V any] struct {
	sync.Mutex
	items map[K]expiredValue[V]
}

// ReadHeavyCacheExpired is a cache optimized for read-heavy operations with expiration support.
// It uses an RWMutex to allow concurrent reads and synchronized writes, storing values with expiration times.
type ReadHeavyCacheExpired[K comparable, V any] struct {
	sync.RWMutex
	items map[K]expiredValue[V]
}

// NewWriteHeavyCacheExpired creates a new instance of WriteHeavyCacheExpired
func NewWriteHeavyCacheExpired[K comparable, V any]() *WriteHeavyCacheExpired[K, V] {
	return &WriteHeavyCacheExpired[K, V]{items: make(map[K]expiredValue[V])}
}

// NewReadHeavyCacheExpired creates a new instance of ReadHeavyCacheExpired
func NewReadHeavyCacheExpired[K comparable, V any]() *ReadHeavyCacheExpired[K, V] {
	return &ReadHeavyCacheExpired[K, V]{items: make(map[K]expiredValue[V])}
}

// Set method for WriteHeavyCacheExpired with a specified expiration duration
func (c *WriteHeavyCacheExpired[K, V]) Set(key K, value V, duration time.Duration) {
	val := expiredValue[V]{
		value:  value,
		expire: time.Now().Add(duration),
	}
	c.Lock()
	defer c.Unlock()
	c.items[key] = val
}

// Get method for WriteHeavyCacheExpired
func (c *WriteHeavyCacheExpired[K, V]) Get(key K) (V, bool) {
	c.Lock()
	defer c.Unlock()
	v, found := c.items[key]
	if !found || time.Now().After(v.expire) {
		var zero V
		return zero, false
	}
	return v.value, true
}

// Set method for ReadHeavyCacheExpired with a specified expiration duration
func (c *ReadHeavyCacheExpired[K, V]) Set(key K, value V, duration time.Duration) {
	val := expiredValue[V]{
		value:  value,
		expire: time.Now().Add(duration),
	}
	c.Lock()
	defer c.Unlock()
	c.items[key] = val
}

// Get method for ReadHeavyCacheExpired
func (c *ReadHeavyCacheExpired[K, V]) Get(key K) (V, bool) {
	c.RLock()
	defer c.RUnlock()
	v, found := c.items[key]
	if !found || time.Now().After(v.expire) {
		var zero V
		return zero, false
	}
	return v.value, true
}

// WriteHeavyCacheInteger is a cache optimized for write-heavy operations for integer-like types.
// It uses a Mutex to synchronize access to the cache items.
type WriteHeavyCacheInteger[K comparable, V interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}] struct {
	sync.Mutex // WriteHeavyCacheInteger uses Mutex for write-heavy scenarios
	items      map[K]V
}

// ReadHeavyCacheInteger is a cache optimized for read-heavy operations for integer-like types.
// It uses an RWMutex to allow concurrent reads and synchronized writes.
type ReadHeavyCacheInteger[K comparable, V interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}] struct {
	sync.RWMutex // ReadHeavyCacheInteger uses RWMutex for read-heavy scenarios
	items        map[K]V
}

// NewWriteHeavyCacheInteger creates a new write-heavy cache for integer types
func NewWriteHeavyCacheInteger[K comparable, V interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}]() *WriteHeavyCacheInteger[K, V] {
	return &WriteHeavyCacheInteger[K, V]{
		items: make(map[K]V),
	}
}

// NewReadHeavyCacheInteger creates a new read-heavy cache for integer types
func NewReadHeavyCacheInteger[K comparable, V interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}]() *ReadHeavyCacheInteger[K, V] {
	return &ReadHeavyCacheInteger[K, V]{
		items: make(map[K]V),
	}
}

// Set sets a value in WriteHeavyCacheInteger, locking for the write operation
func (c *WriteHeavyCacheInteger[K, V]) Set(key K, value V) {
	c.Lock()
	c.items[key] = value
	c.Unlock()
}

// Get retrieves a value from WriteHeavyCacheInteger, locking for read as well
func (c *WriteHeavyCacheInteger[K, V]) Get(key K) (V, bool) {
	c.Lock()
	v, found := c.items[key]
	c.Unlock()
	return v, found
}

// Incr increments a value in WriteHeavyCacheInteger, locking for the operation
func (c *WriteHeavyCacheInteger[K, V]) Incr(key K, value V) {
	c.Lock()
	v, found := c.items[key]
	if found {
		c.items[key] = v + value
	} else {
		c.items[key] = value
	}
	c.Unlock()
}

// Clear removes all items from WriteHeavyCacheInteger.
func (c *WriteHeavyCacheInteger[K, V]) Clear() {
	c.Lock()
	c.items = make(map[K]V)
	c.Unlock()
}

// Set sets a value in ReadHeavyCacheInteger, locking for the write operation
func (c *ReadHeavyCacheInteger[K, V]) Set(key K, value V) {
	c.Lock()
	c.items[key] = value
	c.Unlock()
}

// Get retrieves a value from ReadHeavyCacheInteger, using a read lock
func (c *ReadHeavyCacheInteger[K, V]) Get(key K) (V, bool) {
	c.RLock()
	v, found := c.items[key]
	c.RUnlock()
	return v, found
}

// Incr increments a value in ReadHeavyCacheInteger, locking for the operation
func (c *ReadHeavyCacheInteger[K, V]) Incr(key K, value V) {
	c.Lock()
	v, found := c.items[key]
	if found {
		c.items[key] = v + value
	} else {
		c.items[key] = value
	}
	c.Unlock()
}

// Clear removes all items from ReadHeavyCacheExpired.
func (c *ReadHeavyCacheInteger[K, V]) Clear() {
	c.Lock()
	c.items = make(map[K]V)
	c.Unlock()
}

// LockManager manages a set of mutexes identified by keys of type K.
// It is designed to provide fine-grained locking for operations on individual keys.
type LockManager[K comparable] struct {
	mu    sync.Mutex
	locks map[K]*sync.Mutex
}

// NewLockManager creates a new instance of LockManager.
func NewLockManager[K comparable]() *LockManager[K] {
	return &LockManager[K]{
		locks: make(map[K]*sync.Mutex),
	}
}

// GetMutex retrieves the mutex associated with the given key, creating it if it doesn't exist.
// This function uses a double-check locking mechanism to minimize locking overhead
// for existing mutexes, improving performance in read-heavy scenarios.
func (lm *LockManager[K]) getMutex(id K) *sync.Mutex {
	// First check without locking to improve read performance.
	if lock, exists := lm.locks[id]; exists {
		return lock
	}

	lm.mu.Lock()
	defer lm.mu.Unlock()

	// Double-check to avoid race conditions and ensure only one mutex is created for the key.
	if lock, exists := lm.locks[id]; exists {
		return lock
	}

	lock := &sync.Mutex{}
	lm.locks[id] = lock
	return lock
}

// Lock locks the mutex associated with the given key.
func (lm *LockManager[K]) Lock(id K) {
	lm.getMutex(id).Lock()
}

// GetAndLock retrieves the mutex associated with the given key, locks it, and returns the locked mutex.
// This is useful for cases where you want to obtain and lock the mutex in a single line.
// For example, you can use it like this:
//
//	defer lm.GetAndLock(id).Unlock()
//
// This pattern allows you to ensure the mutex is unlocked when the surrounding function exits.
func (lm *LockManager[K]) GetAndLock(id K) *sync.Mutex {
	mu := lm.getMutex(id)
	mu.Lock()
	return mu
}

// Unlock unlocks the mutex associated with the given key.
func (lm *LockManager[K]) Unlock(id K) {
	lm.getMutex(id).Unlock()
}
