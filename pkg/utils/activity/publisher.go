package activity

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/rotationalio/go-ensign"
)

var (
	mu          sync.Mutex
	enabled     bool
	topic       string
	activity    *NetworkActivity
	ticker      *time.Ticker
	recv        chan *Entry
	wg          *sync.WaitGroup
	client      *ensign.Client
	start, stop sync.Once
	running     uint32
)

// Create the global activity publisher from the configuration.
func NewPublisher(conf Config) (err error) {
	start.Do(func() {

		// TODO: Validate the configuration
		/*
			if err = conf.Validate(); err != nil {
				return err
			}*/

		// Set variables for the singleton
		enabled = conf.Enabled
		topic = conf.Topic
		ticker = time.NewTicker(conf.AggregationWindow)
		recv = make(chan *Entry, 1000)
		wg = &sync.WaitGroup{}

		// TODO: Create Ensign client from configuration

		if enabled {
			wg.Add(1)
			go Publish()
		}
	})

	return err
}

// Publish events from the receiver channel to the Ensign topic.
func Publish() {
	defer wg.Done()
	atomic.AddUint32(&running, 1)
	for {
		select {
		case entry, ok := <-recv:
			if !ok {
				return
			}
		// TODO: Add cases for the ticker and recv channel. If the ticker goes off,
		// publish aggregated events to Ensign. When an event is received on the
		// channel, add it to the aggregation.
		case <-done:
			return
		}
	}
}

// Stop the publisher and wait for the goroutine to exit.
func Stop() {
	stop.Do(func() {
		if atomic.LoadUint32(&running) > 0 {
			close(recv)
			wg.Wait()
		}

		atomic.AddUint32(&running, ^uint32(0))
	})
}

// Entries are created from external go routines and are eventually published by the
// activity publisher.
type Entry struct {
	ts       time.Time
	vasp     uuid.UUID
	activity Activity
}

// Create a new event from an activity type.
func newEvent(activity Activity) *Entry {
	return &Entry{
		activity: activity,
	}
}

// VASP adds a VASP UUID to the event.
func (e *Entry) VASP(id uuid.UUID) *Entry {
	e.vasp = id
	return e
}

// Lookup creates a new activity event for a lookup. Must call Add() to commit the
// event.
func Lookup() *Entry {
	return newEvent(LookupActivity)
}

// Search creates a new activity event for a search. Must call Add() to commit the
// event.
func Search() *Entry {
	return newEvent(SearchActivity)
}

// Send an activity event to the publisher.
func (e *Entry) Add() {
	e.ts = time.Now()
	if enabled {
		recv <- e
	}
}
