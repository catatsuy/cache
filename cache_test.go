package cache_test

import (
	"runtime"
	"sync"
	"testing"

	"github.com/catatsuy/cache"
)

func TestWriteHeavyCache_SetAndGet(t *testing.T) {
	cache := cache.NewWriteHeavyCache[string, int]()
	cache.Set("key1", 100)

	if value, found := cache.Get("key1"); !found {
		t.Errorf("Expected key1 to be found")
	} else if value != 100 {
		t.Errorf("Expected value 100 for key1, but got %d", value)
	}
}

func TestReadHeavyCache_SetAndGet(t *testing.T) {
	cache := cache.NewReadHeavyCache[string, int]()
	cache.Set("key1", 200)

	if value, found := cache.Get("key1"); !found {
		t.Errorf("Expected key1 to be found")
	} else if value != 200 {
		t.Errorf("Expected value 200 for key1, but got %d", value)
	}
}

func TestWriteHeavyCacheInteger_SetAndGet(t *testing.T) {
	cache := cache.NewWriteHeavyCacheInteger[int, int]()
	cache.Set(1, 300)

	if value, found := cache.Get(1); !found {
		t.Errorf("Expected key 1 to be found")
	} else if value != 300 {
		t.Errorf("Expected value 300 for key 1, but got %d", value)
	}
}

func TestReadHeavyCacheInteger_SetAndGet(t *testing.T) {
	cache := cache.NewReadHeavyCacheInteger[int, int]()
	cache.Set(1, 400)

	if value, found := cache.Get(1); !found {
		t.Errorf("Expected key 1 to be found")
	} else if value != 400 {
		t.Errorf("Expected value 400 for key 1, but got %d", value)
	}
}

func TestWriteHeavyCacheInteger_Incr(t *testing.T) {
	cache := cache.NewWriteHeavyCacheInteger[int, int]()
	cache.Incr(1, 10) // Key doesn't exist, should set to 10

	if value, found := cache.Get(1); !found {
		t.Errorf("Expected key 1 to be found")
	} else if value != 10 {
		t.Errorf("Expected value 10 for key 1, but got %d", value)
	}

	cache.Incr(1, 5) // Increment by 5, should be 15

	if value, found := cache.Get(1); !found {
		t.Errorf("Expected key 1 to be found")
	} else if value != 15 {
		t.Errorf("Expected value 15 for key 1, but got %d", value)
	}
}

func TestReadHeavyCacheInteger_Incr(t *testing.T) {
	cache := cache.NewReadHeavyCacheInteger[int, int]()
	cache.Incr(1, 20) // Key doesn't exist, should set to 20

	if value, found := cache.Get(1); !found {
		t.Errorf("Expected key 1 to be found")
	} else if value != 20 {
		t.Errorf("Expected value 20 for key 1, but got %d", value)
	}

	cache.Incr(1, 10) // Increment by 10, should be 30

	if value, found := cache.Get(1); !found {
		t.Errorf("Expected key 1 to be found")
	} else if value != 30 {
		t.Errorf("Expected value 30 for key 1, but got %d", value)
	}
}

func TestWriteHeavyCache_ParallelWrite(t *testing.T) {
	cache := cache.NewWriteHeavyCache[int, int]()
	var wg sync.WaitGroup

	numProcs := runtime.GOMAXPROCS(0) // Get the number of available processors

	// Write to the cache concurrently from GOMAXPROCS goroutines
	for p := range numProcs {
		wg.Add(1)
		go func(procID int) {
			defer wg.Done()
			for i := range 1000 {
				cache.Set(procID*1000+i, i) // Unique keys per goroutine to avoid race conditions
			}
		}(p)
	}

	wg.Wait()

	// Verify if all values are set correctly
	for p := 0; p < numProcs; p++ {
		for i := range 1000 {
			if value, found := cache.Get(p*1000 + i); !found || value != i {
				t.Errorf("Expected value %d for key %d, but got %d (found: %v)", i, p*1000+i, value, found)
			}
		}
	}
}

func TestReadHeavyCache_ParallelWrite(t *testing.T) {
	cache := cache.NewReadHeavyCache[int, int]()
	var wg sync.WaitGroup

	numProcs := runtime.GOMAXPROCS(0) // Get the number of available processors

	// Write to the cache concurrently from GOMAXPROCS goroutines
	for p := range numProcs {
		wg.Add(1)
		go func(procID int) {
			defer wg.Done()
			for i := range 1000 {
				cache.Set(procID*1000+i, i) // Unique keys per goroutine to avoid race conditions
			}
		}(p)
	}

	wg.Wait()

	// Verify if all values are set correctly
	for p := 0; p < numProcs; p++ {
		for i := range 1000 {
			if value, found := cache.Get(p*1000 + i); !found || value != i {
				t.Errorf("Expected value %d for key %d, but got %d (found: %v)", i, p*1000+i, value, found)
			}
		}
	}
}

// Benchmark for WriteHeavyCache's Set method
func BenchmarkWriteHeavyCache_Set(b *testing.B) {
	cache := cache.NewWriteHeavyCache[int, int]()
	b.ResetTimer() // Reset the timer to ignore the setup time

	for i := range b.N {
		cache.Set(i, i)
	}
}

// Benchmark for WriteHeavyCache's Get method
func BenchmarkWriteHeavyCache_Get(b *testing.B) {
	cache := cache.NewWriteHeavyCache[int, int]()
	for i := range 1000 {
		cache.Set(i, i)
	}
	b.ResetTimer()

	for i := range b.N {
		cache.Get(i % 1000) // Access the existing keys
	}
}

// Benchmark for ReadHeavyCache's Set method
func BenchmarkReadHeavyCache_Set(b *testing.B) {
	cache := cache.NewReadHeavyCache[int, int]()
	b.ResetTimer()

	for i := range b.N {
		cache.Set(i, i)
	}
}

// Benchmark for ReadHeavyCache's Get method
func BenchmarkReadHeavyCache_Get(b *testing.B) {
	cache := cache.NewReadHeavyCache[int, int]()
	for i := range 1000 {
		cache.Set(i, i)
	}
	b.ResetTimer()

	for i := range b.N {
		cache.Get(i % 1000)
	}
}

// Benchmark for WriteHeavyCacheInteger's Incr method
func BenchmarkWriteHeavyCacheInteger_Incr(b *testing.B) {
	cache := cache.NewWriteHeavyCacheInteger[int, int]()
	for i := range 1000 {
		cache.Set(i, i)
	}
	b.ResetTimer()

	for i := range b.N {
		cache.Incr(i%1000, 1)
	}
}

// Benchmark for ReadHeavyCacheInteger's Incr method
func BenchmarkReadHeavyCacheInteger_Incr(b *testing.B) {
	cache := cache.NewReadHeavyCacheInteger[int, int]()
	for i := range 1000 {
		cache.Set(i, i)
	}
	b.ResetTimer()

	for i := range b.N {
		cache.Incr(i%1000, 1)
	}
}

func BenchmarkWriteHeavyCache_ParallelWrite(b *testing.B) {
	cache := cache.NewWriteHeavyCache[int, int]()
	var wg sync.WaitGroup

	numProcs := runtime.GOMAXPROCS(0) // Get the number of available processors

	b.ResetTimer() // Reset the timer to ignore setup time

	for range b.N {
		// Parallel writes using GOMAXPROCS goroutines
		for p := range numProcs {
			wg.Add(1)
			go func(procID int) {
				defer wg.Done()
				for j := 0; j < 1000; j++ {
					cache.Set(procID*1000+j, j)
				}
			}(p)
		}
		wg.Wait()
	}
}

func BenchmarkReadHeavyCache_ParallelWrite(b *testing.B) {
	cache := cache.NewReadHeavyCache[int, int]()
	var wg sync.WaitGroup

	numProcs := runtime.GOMAXPROCS(0) // Get the number of available processors

	b.ResetTimer() // Reset the timer to ignore setup time

	for range b.N {
		// Parallel writes using GOMAXPROCS goroutines
		for p := range numProcs {
			wg.Add(1)
			go func(procID int) {
				defer wg.Done()
				for j := range 1000 {
					cache.Set(procID*1000+j, j)
				}
			}(p)
		}
		wg.Wait()
	}
}

// Benchmark for parallel writes in WriteHeavyCacheInteger
func BenchmarkWriteHeavyCacheInteger_ParallelWrite(b *testing.B) {
	cache := cache.NewWriteHeavyCacheInteger[int, int]()
	var wg sync.WaitGroup

	numProcs := runtime.GOMAXPROCS(0) // Get the number of available processors

	b.ResetTimer() // Reset the timer to ignore setup time

	for range b.N {
		// Parallel writes using GOMAXPROCS goroutines
		for p := 0; p < numProcs; p++ {
			wg.Add(1)
			go func(procID int) {
				defer wg.Done()
				for j := 0; j < 1000; j++ {
					cache.Set(procID*1000+j, j)
				}
			}(p)
		}
		wg.Wait()
	}
}

// Benchmark for parallel writes in ReadHeavyCacheInteger
func BenchmarkReadHeavyCacheInteger_ParallelWrite(b *testing.B) {
	cache := cache.NewReadHeavyCacheInteger[int, int]()
	var wg sync.WaitGroup

	numProcs := runtime.GOMAXPROCS(0) // Get the number of available processors

	b.ResetTimer() // Reset the timer to ignore setup time

	for range b.N {
		// Parallel writes using GOMAXPROCS goroutines
		for p := 0; p < numProcs; p++ {
			wg.Add(1)
			go func(procID int) {
				defer wg.Done()
				for j := 0; j < 1000; j++ {
					cache.Set(procID*1000+j, j)
				}
			}(p)
		}
		wg.Wait()
	}
}

// Benchmark for parallel increments in WriteHeavyCacheInteger
func BenchmarkWriteHeavyCacheInteger_ParallelIncr(b *testing.B) {
	cache := cache.NewWriteHeavyCacheInteger[int, int]()
	var wg sync.WaitGroup

	numProcs := runtime.GOMAXPROCS(0) // Get the number of available processors

	// Initialize cache with some values
	for i := range 1000 * numProcs {
		cache.Set(i, i)
	}

	b.ResetTimer() // Reset the timer to ignore setup time

	for range b.N {
		// Parallel increments using GOMAXPROCS goroutines
		for p := 0; p < numProcs; p++ {
			wg.Add(1)
			go func(procID int) {
				defer wg.Done()
				for j := 0; j < 1000; j++ {
					cache.Incr(procID*1000+j, 1)
				}
			}(p)
		}
		wg.Wait()
	}
}

// Benchmark for parallel increments in ReadHeavyCacheInteger
func BenchmarkReadHeavyCacheInteger_ParallelIncr(b *testing.B) {
	cache := cache.NewReadHeavyCacheInteger[int, int]()
	var wg sync.WaitGroup

	numProcs := runtime.GOMAXPROCS(0) // Get the number of available processors

	// Initialize cache with some values
	for i := range 1000 * numProcs {
		cache.Set(i, i)
	}

	b.ResetTimer() // Reset the timer to ignore setup time

	for range b.N {
		// Parallel increments using GOMAXPROCS goroutines
		for p := 0; p < numProcs; p++ {
			wg.Add(1)
			go func(procID int) {
				defer wg.Done()
				for j := 0; j < 1000; j++ {
					cache.Incr(procID*1000+j, 1)
				}
			}(p)
		}
		wg.Wait()
	}
}