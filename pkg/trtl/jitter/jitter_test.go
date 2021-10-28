package jitter_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/trtl/jitter"
)

// NOTE: because this is a stochastic test it might fail; run again to make sure the
// failure wasn't a fluke of some extremely unlikely event in the normal distribution.
func TestTicker(t *testing.T) {
	ticker := jitter.New(100*time.Millisecond, 5*time.Millisecond)
	prev := time.Now()
	prevDelta := time.Duration(0)
	for i := 0; i < 10; i++ {
		now := <-ticker.C
		delta := now.Sub(prev)
		require.NotEqual(t, prevDelta, delta)
		require.Less(t, time.Duration(10*time.Millisecond), delta)
		require.Greater(t, time.Duration(200*time.Millisecond), delta)
		prev = now
		prevDelta = delta
	}

	ticker.Stop()
}
