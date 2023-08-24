package activity

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rotationalio/go-ensign"
	"github.com/rotationalio/go-ensign/mock"
)

var (
	start, stop sync.Once
	mu          sync.Mutex
	running     bool
	enabled     bool
	ticker      *time.Ticker
	recv        chan *Entry
	wg          *sync.WaitGroup
	client      *ensign.Client
	emock       *mock.Ensign
)

// Start the global activity publisher from the configuration.
func Start(conf Config) (err error) {
	start.Do(func() {
		// Validate the configuration
		if err = conf.Validate(); err != nil {
			return
		}

		enabled = conf.Enabled
		if enabled {
			ticker = time.NewTicker(conf.AggregationWindow)
			recv = make(chan *Entry, 1000)

			if conf.Testing {
				// In testing mode, create the Ensign client using a mock server
				emock = mock.New(nil)
				if client, err = ensign.New(ensign.WithMock(emock)); err != nil {
					return
				}
			} else if client, err = conf.Ensign.Client(); err != nil {
				return
			}

			wg = &sync.WaitGroup{}
			wg.Add(1)
			go Publish()
		}
	})

	return err
}

// Global goroutine that publishes activity entries from the receiver channel to the
// Ensign topic as events.
func Publish() {
	mu.Lock()
	running = true
	mu.Unlock()
	defer wg.Done()
	for {
		select {
		case _, ok := <-recv:
			if !ok {
				return
			}

			// TODO: Add the activity event to the aggregation
		case <-ticker.C:
			// TODO: Publish the aggregated events to Ensign and reset the aggregation
			client.Status(context.Background())
		}
	}
}

// Stop the publisher and wait for the goroutine to exit.
func Stop() {
	stop.Do(func() {
		mu.Lock()
		defer mu.Unlock()
		// Only stop the publisher if it is running
		if running {
			close(recv)
			wg.Wait()
		}

		running = false
	})
}

// Reset the publisher to allow NewPublisher() to be called again, this method should
// only be used for testing.
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	start = sync.Once{}
	stop = sync.Once{}

	if running {
		close(recv)
		wg.Wait()
	}
	running = false
	enabled = false
	ticker = nil
	recv = nil
	wg = nil
	client = nil
	emock = nil
}

// Expose the ensign server mock to the tests.
func GetEnsignMock() *mock.Ensign {
	return emock
}

// Entries are created from external go routines and are eventually published as Events
// to Ensign by the activity publisher.
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

// Lookup creates a new activity entry for a lookup. Must call Add() to commit the
// entry.
func Lookup() *Entry {
	return newEvent(LookupActivity)
}

// Search creates a new activity entry for a search. Must call Add() to commit the
// entry.
func Search() *Entry {
	return newEvent(SearchActivity)
}

// Search creates a new activity entry for registering a Vasp. Must call Add() to
// commit the entry.
func Register() *Entry {
	return newEvent(RegisterActivity)
}

// Add the activity entry to the publisher.
func (e *Entry) Add() {
	e.ts = time.Now()
	if enabled {
		recv <- e
	}
}
