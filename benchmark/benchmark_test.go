package benchmark_test

import (
	"strconv"
	"testing"

	"github.com/catatsuy/cache"
	cs "github.com/catatsuy/sync/singleflight"
	"golang.org/x/sync/singleflight"
)

func BenchmarkSingleflight(b *testing.B) {
	for _, keys := range []int{1, 10} { // 衝突度
		b.Run("std/keys="+strconv.Itoa(keys), func(b *testing.B) {
			var sg singleflight.Group
			runStd(b, &sg, keys)
		})
		b.Run("std-cast/keys="+strconv.Itoa(keys), func(b *testing.B) {
			var sg singleflight.Group
			runStdWithCast(b, &sg, keys)
		})
		b.Run("generics/keys="+strconv.Itoa(keys), func(b *testing.B) {
			var sg cs.Group[int]
			runGenerics(b, &sg, keys)
		})
		b.Run("custom/keys="+strconv.Itoa(keys), func(b *testing.B) {
			sf := cache.NewSingleflightGroup[int]()
			runCustom(b, sf, keys)
		})
	}
}

func runStd(b *testing.B, sg *singleflight.Group, keyCount int) {
	b.ReportAllocs()
	keys := genKeys(keyCount)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := keys[i%keyCount]
			i++
			v, _, _ := sg.Do(key, func() (any, error) { return i, nil })

			_ = v
		}
	})
}

func runStdWithCast(b *testing.B, sg *singleflight.Group, keyCount int) {
	b.ReportAllocs()
	keys := genKeys(keyCount)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := keys[i%keyCount]
			i++
			v, _, _ := sg.Do(key, func() (any, error) { return i, nil })

			_ = v.(int)
		}
	})
}

func runGenerics(b *testing.B, sg *cs.Group[int], keyCount int) {
	b.ReportAllocs()
	keys := genKeys(keyCount)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := keys[i%keyCount]
			i++
			v, _, _ := sg.Do(key, func() (int, error) { return i, nil })

			_ = v
		}
	})
}

func runCustom(b *testing.B, sf *cache.SingleflightGroup[int], keyCount int) {
	b.ReportAllocs()
	keys := genKeys(keyCount)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := keys[i%keyCount]
			i++
			v, _ := sf.Do(key, func() (int, error) { return i, nil })

			_ = v
		}
	})
}

func genKeys(n int) []string {
	keys := make([]string, n)
	for i := range n {
		keys[i] = "key-" + strconv.Itoa(i)
	}
	return keys
}
