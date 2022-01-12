package trtl

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

var (
	pmPuts  *prometheus.CounterVec // count of trtl Puts per namespace
	pmGets  *prometheus.CounterVec // count of trtl Gets per namespace
	pmDels  *prometheus.CounterVec // count of trtl Deletes per namespace
	pmIters *prometheus.CounterVec // count of trtl Iters per namespace
	// pmObjects    *prometheus.CounterVec   // count of objects being managed by trtl, by namespace
	// pmTombstones *prometheus.CounterVec   // count of tombstones per namespace; increases on delete, decrease on overwrite of tombstone
	pmLatency     *prometheus.HistogramVec // the time it is taking for successful RPC calls to complete, labeled by RPC type, success, and failure
	pmSyncs       *prometheus.CounterVec   // count of anti entropy sessions per peer and per region
	pmSyncLatency *prometheus.HistogramVec // the time it is taking for anti entropy sessions to complete, by peer
)

// A MetricsService manages Prometheus metrics
type MetricsService struct {
	srv *http.Server
}

func NewMetricsService() (*MetricsService, error) {
	initMetrics()
	return &MetricsService{srv: &http.Server{}}, nil
}

// Serve serves the Prometheus metrics
func (m *MetricsService) Serve(addr string) error {
	if err := registerMetrics(); err != nil {
		return err
	}
	m.srv.Addr = addr
	log.Info().Msg(fmt.Sprintf("serving prometheus metrics at http://%s/metrics", addr))
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := m.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Err(err).Msg("metrics server shutdown prematurely")
		}
	}()
	return nil
}

// Gracefully shutdown the Prometheus metrics service
func (m *MetricsService) Shutdown() error {
	// Might want to share context from Trtl more globally?
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := m.srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("unable to gracefully shutdown prometheus metrics")
	}
	return nil
}

const (
	pmNamespace = "trtl"
)

func initMetrics() {
	pmPuts = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: pmNamespace,
		Name:      "puts",
		Help:      "the count of puts, labeled by namespace",
	}, []string{"namespace"})

	pmGets = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: pmNamespace,
		Name:      "gets",
		Help:      "the count of gets, labeled by namespace",
	}, []string{"namespace"})

	pmDels = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: pmNamespace,
		Name:      "deletes",
		Help:      "the count of deletes, labeled by namespace",
	}, []string{"namespace"})

	pmIters = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: pmNamespace,
		Name:      "iters",
		Help:      "the count of iters, labeled by namespace",
	}, []string{"namespace"})

	// pmObjects = prometheus.NewCounterVec(prometheus.CounterOpts{
	// 	Namespace: pmNamespace,
	// 	Name:      "objects",
	// 	Help:      "the count of trtl objects, labeled by namespace",
	// }, []string{"namespace"})

	// pmTombstones = prometheus.NewCounterVec(prometheus.CounterOpts{
	// 	Namespace: pmNamespace,
	// 	Name:      "tombstones",
	// 	Help:      "the count of tombstones, labeled by namespace",
	// }, []string{"namespace"})

	pmLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: pmNamespace,
		Name:      "latency",
		Help:      "time to RPC call completion, labeled by RPC (Put, Get, Delete, Iter)",
	}, []string{"call"})

	pmSyncs = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: pmNamespace,
		Name:      "syncs",
		Help:      "the count of anti-entropy sessions, labeled by peer and region",
	}, []string{"peer", "region"})

	pmSyncLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: pmNamespace,
		Name:      "sync_latency",
		Help:      "time to anti-entropy session completion, labeled by peer",
	}, []string{"peer"})
}

func registerMetrics() error {
	if err := prometheus.Register(pmPuts); err != nil {
		log.Debug().Err(err).Msg("unable to register pmPuts")
		return err
	}
	if err := prometheus.Register(pmGets); err != nil {
		log.Debug().Err(err).Msg("unable to register pmGets")
		return err
	}
	if err := prometheus.Register(pmDels); err != nil {
		log.Debug().Err(err).Msg("unable to register pmDels")
		return err
	}
	if err := prometheus.Register(pmIters); err != nil {
		log.Debug().Err(err).Msg("unable to register pmIters")
		return err
	}
	// if err := prometheus.Register(pmObjects); err != nil {
	// 	log.Debug().Err(err).Msg("unable to register pmObjects")
	// 	return err
	// }
	// if err := prometheus.Register(pmTombstones); err != nil {
	// 	log.Debug().Err(err).Msg("unable to register pmTombstones")
	// 	return err
	// }
	if err := prometheus.Register(pmLatency); err != nil {
		log.Debug().Err(err).Msg("unable to register pmLatency")
		return err
	}
	if err := prometheus.Register(pmSyncs); err != nil {
		log.Debug().Err(err).Msg("unable to register pmSyncs")
	}
	if err := prometheus.Register(pmSyncLatency); err != nil {
		log.Debug().Err(err).Msg("unable to register pmSyncLatency")
	}
	return nil
}
