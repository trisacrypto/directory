package replica

import (
	"math"
	"math/rand"
	"time"
)

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

const (
	lambda = -0.004
)

var (
	random *rand.Rand
)

// TimeProbability returns an exponentially decaying probability between 0 and 1 that
// smoothly decreases to zero over the course of a day (24 hours). E.g. the longer the
// time since the timestamp, the lower the probability that is returned.
func TimeProbability(ts time.Time) float64 {
	since := time.Since(ts).Minutes()
	return math.Exp(since * lambda)
}

// ReplicateObjectRoulette performs a roulette roll to see if the object should be
// replicated basesd on its TimeProbability.
func ReplicateObjectRoulette(ts time.Time) bool {
	return random.Float64() < TimeProbability(ts)
}

func Seed(seed int64) {
	random = rand.New(rand.NewSource(seed))
}
