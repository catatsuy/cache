package cache

import "sync"

type WriteHeavyCache[K comparable, V any] struct {
	sync.Mutex // For Write Heavy, use Mutex for all operations
	items      map[K]V
}

type ReadHeavyCache[K comparable, V any] struct {
	sync.RWMutex // For Read Heavy, allow concurrent read access
	items        map[K]V
}

// Write Heavy methods
func (c *WriteHeavyCache[K, V]) Set(key K, value V) {
	c.Lock() // Lock for writing
	c.items[key] = value
	c.Unlock()
}

func (c *WriteHeavyCache[K, V]) Get(key K) (V, bool) {
	c.Lock() // Lock for reading as well, single-thread access
	v, found := c.items[key]
	c.Unlock()
	return v, found
}

// Read Heavy methods
func (c *ReadHeavyCache[K, V]) Set(key K, value V) {
	c.Lock() // Lock for writing
	c.items[key] = value
	c.Unlock()
}

func (c *ReadHeavyCache[K, V]) Get(key K) (V, bool) {
	c.RLock() // RLock for reading, allows multiple concurrent reads
	v, found := c.items[key]
	c.RUnlock()
	return v, found
}

// Constructor for Write Heavy cache
func NewWriteHeavyCache[K comparable, V any]() *WriteHeavyCache[K, V] {
	return &WriteHeavyCache[K, V]{
		items: make(map[K]V),
	}
}

// Constructor for Read Heavy cache
func NewReadHeavyCache[K comparable, V any]() *ReadHeavyCache[K, V] {
	return &ReadHeavyCache[K, V]{
		items: make(map[K]V),
	}
}

// cacheInteger with manual type constraints for integer-like types.
type WriteHeavyCacheInteger[K comparable, V interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}] struct {
	sync.Mutex // Use Mutex for write-heavy scenarios
	items      map[K]V
}

type ReadHeavyCacheInteger[K comparable, V interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}] struct {
	sync.RWMutex // Use RWMutex for read-heavy scenarios
	items        map[K]V
}

// NewWriteHeavyCacheInteger constructor for creating a new write-heavy cache.
func NewWriteHeavyCacheInteger[K comparable, V interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}]() *WriteHeavyCacheInteger[K, V] {
	return &WriteHeavyCacheInteger[K, V]{
		items: make(map[K]V),
	}
}

// NewReadHeavyCacheInteger constructor for creating a new read-heavy cache.
func NewReadHeavyCacheInteger[K comparable, V interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}]() *ReadHeavyCacheInteger[K, V] {
	return &ReadHeavyCacheInteger[K, V]{
		items: make(map[K]V),
	}
}

// Write-heavy Set method using Mutex
func (c *WriteHeavyCacheInteger[K, V]) Set(key K, value V) {
	c.Lock() // Lock for writing
	c.items[key] = value
	c.Unlock() // Unlock after write is complete
}

// Write-heavy Get method using Mutex
func (c *WriteHeavyCacheInteger[K, V]) Get(key K) (V, bool) {
	c.Lock() // Lock for reading (Mutex is used for both read and write)
	v, found := c.items[key]
	c.Unlock() // Unlock after read is complete
	return v, found
}

// Write-heavy Incr method using Mutex
func (c *WriteHeavyCacheInteger[K, V]) Incr(key K, value V) {
	c.Lock() // Lock for incrementing
	v, found := c.items[key]
	if found {
		c.items[key] = v + value
	} else {
		c.items[key] = value
	}
	c.Unlock() // Unlock after modification
}

// Read-heavy Set method using RWMutex
func (c *ReadHeavyCacheInteger[K, V]) Set(key K, value V) {
	c.Lock() // Lock for writing
	c.items[key] = value
	c.Unlock() // Unlock after write is complete
}

// Read-heavy Get method using RWMutex
func (c *ReadHeavyCacheInteger[K, V]) Get(key K) (V, bool) {
	c.RLock() // RLock for reading, allows multiple concurrent reads
	v, found := c.items[key]
	c.RUnlock() // Unlock after read is complete
	return v, found
}

// Read-heavy Incr method using RWMutex
func (c *ReadHeavyCacheInteger[K, V]) Incr(key K, value V) {
	c.Lock() // Lock for incrementing (write operation)
	v, found := c.items[key]
	if found {
		c.items[key] = v + value
	} else {
		c.items[key] = value
	}
	c.Unlock() // Unlock after modification
}
