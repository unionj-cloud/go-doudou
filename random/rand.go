package random

import (
	"math/rand"
	"time"
)

// RandInt generates random int between min and max
func RandInt(min int, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min) + min
}
