package cache_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/unionj-cloud/go-doudou/v2/toolkit/cache"
)

func TestTinyLFU_Get_CorruptionOnExpiry(t *testing.T) {
	strFor := func(i int) string {
		return fmt.Sprintf("a string %d", i)
	}
	keyName := func(i int) string {
		return fmt.Sprintf("key-%00000d", i)
	}

	mycache := cache.NewTinyLFU(1000, 1*time.Second)
	size := 50000
	// Put a bunch of stuff in the cache with a TTL of 1 second
	for i := 0; i < size; i++ {
		key := keyName(i)
		mycache.Set(key, []byte(strFor(i)))
	}

	// Read stuff for a bit longer than the TTL - that's when the corruption occurs
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	done := ctx.Done()
loop:
	for {
		select {
		case <-done:
			// this is expected
			break loop
		default:
			i := rand.Intn(size)
			key := keyName(i)

			b, ok := mycache.Get(key)
			if !ok {
				continue loop
			}

			got := string(b)
			expected := strFor(i)
			if got != expected {
				t.Fatalf("expected=%q got=%q key=%q", expected, got, key)
			}
		}
	}
}
