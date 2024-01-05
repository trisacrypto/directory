package replica_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/trtl/replica"
)

func TestTimeProbability(t *testing.T) {
	cases := []struct {
		before      time.Duration
		minExpected float64
		maxExpected float64
	}{
		{-1 * time.Second, 0.99, 1.0},
		{-1 * time.Minute, 0.99, 1.0},
		{-10 * time.Minute, 0.96, .97},
		{-30 * time.Minute, 0.88, .89},
		{-1 * time.Hour, 0.78, .79},
		{-2 * time.Hour, 0.61, .62},
		{-6 * time.Hour, 0.23, .24},
		{-12 * time.Hour, .05, .06},
		{-24 * time.Hour, .003, .004},
		{-48 * time.Hour, 0.0, .0001},
	}

	for i, tc := range cases {
		prob := replica.TimeProbability(time.Now().Add(tc.before))
		require.Greater(t, prob, tc.minExpected, "test case %d prob is less than min expected", i)
		require.Less(t, prob, tc.maxExpected, "test case %d prob is greater than max expected", i)
	}

}

func TestReplicateObjectRoulette(t *testing.T) {
	cases := []struct {
		rolls       int
		minExpected int
		maxExpected int
		before      time.Duration
	}{
		{100, 95, 98, -10 * time.Minute},
		{100, 65, 75, -2 * time.Hour},
		{100, 0, 10, -12 * time.Hour},
	}

	for i, tc := range cases {
		count := 0
		ts := time.Now().Add(tc.before)
		replica.Seed(24)

		for i := 0; i < tc.rolls; i++ {
			if replica.ReplicateObjectRoulette(ts) {
				count++
			}
		}
		require.Greater(t, count, tc.minExpected, "test case %d prob is less than min expected", i)
		require.Less(t, count, tc.maxExpected, "test case %d prob is greater than max expected", i)
	}
}
