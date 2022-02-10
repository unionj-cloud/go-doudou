package loadbalance

import (
	"math"
	"math/rand"
)

func Subset(backends []string, clientId int, subsetSize int) []string {
	subsetCount := int(math.Ceil(float64(len(backends)) / float64(subsetSize)))
	round := int64(clientId / subsetCount)
	rand.Seed(round)
	rand.Shuffle(len(backends), func(i, j int) {
		backends[i], backends[j] = backends[j], backends[i]
	})
	subsetId := clientId % subsetCount
	start := subsetId * subsetSize
	return backends[start : start+subsetSize]
}
