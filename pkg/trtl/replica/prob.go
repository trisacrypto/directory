package replica

import (
	crand "crypto/rand"
	"encoding/binary"
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	lambda = -0.004
)

var (
	random *rand.Rand
	mu     sync.Mutex
)

func init() {
	var b [8]byte
	if _, err := crand.Read(b[:]); err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number")
	}
	random = rand.New(rand.NewSource(int64(binary.LittleEndian.Uint64(b[:]))))
}

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
	mu.Lock()
	defer mu.Unlock()
	return random.Float64() < TimeProbability(ts)
}

// Seed the random number generator for testing purposes
func TestSeed(seed int64) {
	random = rand.New(rand.NewSource(seed))
}
