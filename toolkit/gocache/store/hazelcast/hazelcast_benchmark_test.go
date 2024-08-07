package hazelcast

import (
	"context"
	"fmt"
	"math"
	"testing"

	"github.com/hazelcast/hazelcast-go-client"
	lib_store "github.com/unionj-cloud/go-doudou/v2/toolkit/gocache/lib/store"
)

func BenchmarkHazelcastSet(b *testing.B) {
	ctx := context.Background()

	hzClient, err := hazelcast.StartNewClient(ctx)
	if err != nil {
		b.Fatalf("Failed to start client: %v", err)
	}

	store := NewHazelcast(hzClient, "gocache")

	for k := 0.; k <= 10; k++ {
		n := int(math.Pow(2, k))
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			for i := 0; i < b.N*n; i++ {
				key := fmt.Sprintf("test-%d", n)
				value := []byte(fmt.Sprintf("value-%d", n))
				store.Set(ctx, key, value, lib_store.WithTags([]string{fmt.Sprintf("tag-%d", n)}))
			}
		})
	}
}

func BenchmarkHazelcastGet(b *testing.B) {
	ctx := context.Background()

	hzClient, err := hazelcast.StartNewClient(ctx)
	if err != nil {
		b.Fatalf("Failed to start client: %v", err)
	}

	store := NewHazelcast(hzClient, "gocache")

	key := "test"
	value := []byte("value")

	store.Set(ctx, key, value)

	for k := 0.; k <= 10; k++ {
		n := int(math.Pow(2, k))
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			for i := 0; i < b.N*n; i++ {
				_, _ = store.Get(ctx, key)
			}
		})
	}
}
