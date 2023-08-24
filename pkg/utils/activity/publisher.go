package activity

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rotationalio/go-ensign"
	"github.com/rotationalio/go-ensign/mock"
	"github.com/rs/zerolog/log"
)

var (
	start, stop sync.Once
	mu          sync.Mutex
	running     bool
	enabled     bool
	topic       string
	ticker      *time.Ticker
	recv        chan *Entry
	activity    *NetworkActivity
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
			var window time.Duration
			if window = conf.AggregationWindow; window <= 0 {
				err = ErrInvalidWindow
				return
			}

			topic = conf.Topic
			ticker = time.NewTicker(window)
			recv = make(chan *Entry, 1000)
			activity = New(conf.Network, window, time.Now())

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
		case entry, ok := <-recv:
			if !ok {
				return
			}

			// Add the entry to the aggregation
			if entry.vasp != uuid.Nil {
				activity.IncrVASP(entry.vasp, entry.activity)
			} else {
				activity.Incr(entry.activity)
			}
		case <-ticker.C:
			var (
				err   error
				event *ensign.Event
			)
			if event, err = activity.Event(); err != nil {
				log.Error().Err(err).Msg("could not create activity event")
				activity.Reset()
				continue
			}

			if err = client.Publish(topic, event); err != nil {
				log.Error().Err(err).Msg("could not publish activity event")
			}

			// Reset the activity counts for the next window
			activity.Reset()
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

// Add the activity entry to the publisher.
func (e *Entry) Add() {
	e.ts = time.Now()
	if enabled {
		select {
		case recv <- e:
		default:
		}
	}
}
