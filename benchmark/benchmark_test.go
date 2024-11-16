package benchmark_test

import (
	"strconv"
	"sync"
	"testing"

	"github.com/catatsuy/cache"
	cs "github.com/catatsuy/sync/singleflight"
	"golang.org/x/sync/singleflight"
)

func BenchmarkStandardSingleflight(b *testing.B) {
	var sg singleflight.Group
	var wg sync.WaitGroup

	keys := generateKeys(10)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := keys[i%10]
			sg.Do(key, func() (interface{}, error) {
				return i, nil
			})
		}(i)
	}
	wg.Wait()
}

func BenchmarkStandardSingleflightCast(b *testing.B) {
	var sg singleflight.Group
	var wg sync.WaitGroup

	keys := generateKeys(10)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := keys[i%10]
			ii, _, _ := sg.Do(key, func() (interface{}, error) {
				return i, nil
			})
			if _, ok := ii.(int); !ok {
				b.Fatalf("unexpected type: %T", ii)
			}
		}(i)
	}
	wg.Wait()
}

func BenchmarkGenericsSingleflight(b *testing.B) {
	var sg cs.Group[int]
	var wg sync.WaitGroup

	keys := generateKeys(10)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := keys[i%10]
			sg.Do(key, func() (int, error) {
				return i, nil
			})
		}(i)
	}
	wg.Wait()
}

func BenchmarkCustomSingleflight(b *testing.B) {
	sf := cache.NewSingleflightGroup[int]()
	var wg sync.WaitGroup

	keys := generateKeys(10)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := keys[i%10]
			sf.Do(key, func() (int, error) {
				return i, nil
			})
		}(i)
	}
	wg.Wait()
}

func generateKeys(n int) []string {
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		keys[i] = "key-" + strconv.Itoa(i)
	}
	return keys
}
