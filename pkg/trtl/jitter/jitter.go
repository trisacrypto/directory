/*
Package jitter provides a stochastic ticker that returns ticks at a random interval
specified by a normal distribution with a mean periodicity and a standard deviation,
sigma both of which are time.Durations.

This is a simplified version of https://github.com/lthibault/jitterbug.
*/
package jitter

import (
	"math/rand"
	"time"
)

// New returns a new stochastic timer with the specified interval and standard deviation.
func New(interval, sigma time.Duration) (j *Ticker) {
	c := make(chan time.Time)
	j = &Ticker{
		C:        c,
		Interval: interval,
		Sigma:    sigma,
		done:     make(chan struct{}),
	}
	go j.loop(c)
	return j
}

// Ticker behaves like time.Ticker (listen on Ticker.C for timestamps).
type Ticker struct {
	C        <-chan time.Time
	Interval time.Duration
	Sigma    time.Duration
	done     chan struct{}
}

// Stop the Ticker.
func (j *Ticker) Stop() {
	close(j.done)
}

func (t *Ticker) loop(c chan<- time.Time) {
	defer close(c)
	for {
		time.Sleep(t.calcDelay())

		select {
		case <-t.done:
			return
		case c <- time.Now():
		default: // there may be no recv
		}
	}
}

func (j *Ticker) calcDelay() time.Duration {
	s := rand.NormFloat64() * float64(j.Sigma)
	return j.Interval + time.Duration(s)
}
