package trtl

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/rotationalio/honu"
	"github.com/rotationalio/honu/options"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/trtl/metrics"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
)

// Monitor is an independent service which periodically scans the trtl storage and
// determines how many objects, tombstones, size, etc. is utilized by the internal db.
type Monitor struct {
	conf config.MetricsConfig
	db   *honu.DB
	stop chan struct{}
}

func NewMonitor(conf config.MetricsConfig, db *honu.DB) (*Monitor, error) {
	return &Monitor{
		conf: conf,
		db:   db,
		stop: make(chan struct{}),
	}, nil
}

// Run the monitor which periodically wakes up and measures the metrics in the database.
func (m *Monitor) Run() {
	if !m.conf.Enabled {
		log.Warn().Msg("trtl monitor disabled")
		return
	}

	ticker := time.NewTicker(m.conf.Interval)
	log.Info().Dur("interval", m.conf.Interval).Msg("trtl monitor started")

monitor:
	for {
		// Wait for next tick or a stop message
		select {
		case <-m.stop:
			log.Info().Msg("trtl monitor stopping")
			return
		case <-ticker.C:
		}

		// Collect the metrics for the database
		log.Debug().Msg("starting metrics analysis of trtl database")
		start := time.Now()

		if err := m.Measure(); err != nil {
			sentry.Error(nil).Err(err).Msg("could not measure trtl database metrics")
			continue monitor
		}
		log.Debug().Dur("duration", time.Since(start)).Msg("trtl metrics analysis complete")
	}
}

func (m *Monitor) Shutdown() error {
	if m.stop != nil {
		// Will block until the current measurement is complete
		m.stop <- struct{}{}

		// Close the channel and set to nil so that multiple shutdown calls don't block
		close(m.stop)
		m.stop = nil
	}
	return nil
}

func (m *Monitor) Measure() (err error) {
	for _, namespace := range measuredNamespaces {
		if merr := m.MeasureNamespace(namespace); merr != nil {
			err = multierror.Append(err, fmt.Errorf("could not measure namespace %s: %w", namespace, merr))
		}
	}
	return err
}

func (m *Monitor) MeasureNamespace(namespace string) error {
	iter, err := m.db.Iter(nil, options.WithNamespace(namespace), options.WithTombstones())
	if err != nil {
		return err
	}
	defer iter.Release()
	log.Trace().Str("namespace", namespace).Msg("measuring namespace")

	// Measurements
	var objects, tombstones, bytes uint64

objects:
	for iter.Next() {
		// Load the object metadata to check if it is a tombstone
		obj, err := iter.Object()
		if err != nil {
			sentry.Error(nil).Err(err).
				Str("namespace", namespace).
				Bytes("key", iter.Key()).
				Msg("could not unmarshal honu metadata")
			continue objects
		}

		// Count tombstones vs objects
		if obj.Tombstone() {
			tombstones++
		} else {
			objects++
		}

		// Measure the size of the value in bytes
		bytes += uint64(len(iter.Value()))
	}

	if err := iter.Error(); err != nil {
		return err
	}

	metrics.PmDatabaseSize.WithLabelValues(namespace).Set(float64(bytes))
	metrics.PmCurrentObjects.WithLabelValues(namespace).Set(float64(objects))
	metrics.PmCurrentTombstones.WithLabelValues(namespace).Set(float64(tombstones))
	return nil
}
