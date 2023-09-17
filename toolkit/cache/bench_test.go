package cache_test

import (
	"github.com/unionj-cloud/go-doudou/v2/toolkit/cache"
	"strings"
	"testing"

)

func BenchmarkOnce(b *testing.B) {
	mycache := newCacheWithLocal(newRing())
	obj := &Object{
		Str: strings.Repeat("my very large string", 10),
		Num: 42,
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var dst Object
			err := mycache.Once(&cache.Item{
				Key:   "bench-once",
				Value: &dst,
				Do: func(*cache.Item) (interface{}, error) {
					return obj, nil
				},
			})
			if err != nil {
				b.Fatal(err)
			}
			if dst.Num != 42 {
				b.Fatalf("%d != 42", dst.Num)
			}
		}
	})
}

func BenchmarkSet(b *testing.B) {
	mycache := newCacheWithLocal(newRing())
	obj := &Object{
		Str: strings.Repeat("my very large string", 10),
		Num: 42,
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := mycache.Set(&cache.Item{
				Key:   "bench-set",
				Value: obj,
			}); err != nil {
				b.Fatal(err)
			}
		}
	})
}
