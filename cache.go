package cache

import "sync"

type WriteHeavyCache[K comparable, V any] struct {
	sync.Mutex // WriteHeavyCache uses Mutex for all operations
	items      map[K]V
}

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

// WriteHeavyCacheInteger is a write-heavy cache for integer-like types
type WriteHeavyCacheInteger[K comparable, V interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}] struct {
	sync.Mutex // WriteHeavyCacheInteger uses Mutex for write-heavy scenarios
	items      map[K]V
}

// ReadHeavyCacheInteger is a read-heavy cache for integer-like types
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
