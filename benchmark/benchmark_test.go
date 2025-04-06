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

	for i := range b.N {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := keys[i%10]
			sg.Do(key, func() (any, error) {
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

	for i := range b.N {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := keys[i%10]
			ii, _, _ := sg.Do(key, func() (any, error) {
				return i, nil
			})
			if _, ok := ii.(int); !ok {
				b.Errorf("unexpected type: %T", ii)
				return
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

	for i := range b.N {
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

	for i := range b.N {
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
	for i := range n {
		keys[i] = "key-" + strconv.Itoa(i)
	}
	return keys
}
