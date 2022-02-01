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
	PmPuts  *prometheus.CounterVec // count of trtl Puts per namespace
	PmGets  *prometheus.CounterVec // count of trtl Gets per namespace
	PmDels  *prometheus.CounterVec // count of trtl Deletes per namespace
	PmIters *prometheus.CounterVec // count of trtl Iters per namespace
	// PmObjects    *prometheus.CounterVec   // count of objects being managed by trtl, by namespace
	// PmTombstones *prometheus.CounterVec   // count of tombstones per namespace; increases on delete, decrease on overwrite of tombstone
	PmLatency       *prometheus.HistogramVec // the time it is taking for successful RPC calls to complete, labeled by RPC type, success, and failure
	PmAESyncs       *prometheus.CounterVec   // count of anti entropy sessions per peer and per region
	PmAESyncLatency *prometheus.HistogramVec // the time it is taking for anti entropy sessions to complete, by peer
	PmAEPushes      *prometheus.HistogramVec // pushed objects during anti entropy, by peer and region
	PmAEPulls       *prometheus.HistogramVec // pulled objects during anti entropy, by peer and region
	PmAEPushVSPull  prometheus.Gauge         // a gauge of objects pushed vs pulled
	PmAEStomps      *prometheus.CounterVec   // count of stomped versions, per peer and region
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
	PmNamespace = "trtl"
)

func initMetrics() {
	PmPuts = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: PmNamespace,
		Name:      "puts",
		Help:      "the count of puts, labeled by namespace",
	}, []string{"namespace"})

	PmGets = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: PmNamespace,
		Name:      "gets",
		Help:      "the count of gets, labeled by namespace",
	}, []string{"namespace"})

	PmDels = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: PmNamespace,
		Name:      "deletes",
		Help:      "the count of deletes, labeled by namespace",
	}, []string{"namespace"})

	PmIters = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: PmNamespace,
		Name:      "iters",
		Help:      "the count of iters, labeled by namespace",
	}, []string{"namespace"})

	// PmObjects = prometheus.NewCounterVec(prometheus.CounterOpts{
	// 	Namespace: PmNamespace,
	// 	Name:      "objects",
	// 	Help:      "the count of trtl objects, labeled by namespace",
	// }, []string{"namespace"})

	// PmTombstones = prometheus.NewCounterVec(prometheus.CounterOpts{
	// 	Namespace: PmNamespace,
	// 	Name:      "tombstones",
	// 	Help:      "the count of tombstones, labeled by namespace",
	// }, []string{"namespace"})

	PmLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: PmNamespace,
		Name:      "latency",
		Help:      "time to RPC call completion, labeled by RPC (Put, Get, Delete, Iter)",
	}, []string{"call"})

	PmAESyncs = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: PmNamespace,
		Name:      "syncs",
		Help:      "the count of anti-entropy sessions, labeled by peer and region",
	}, []string{"peer", "region"})

	PmAESyncLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: PmNamespace,
		Name:      "sync_latency",
		Help:      "time to anti-entropy session completion, labeled by peer",
	}, []string{"peer"})

	PmAEPulls = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: PmNamespace,
		Name:      "pulls",
		Help:      "pulled objects during anti entropy, labeled by peer and region",
	}, []string{"peer", "region"})

	PmAEPushes = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: PmNamespace,
		Name:      "pushes",
		Help:      "pushed objects during anti entropy, labeled by peer and region",
	}, []string{"peer", "region"})

	PmAEPushVSPull = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: PmNamespace,
		Name:      "push_vs_pull",
		Help:      "objects pushed vs pulled",
	})

	PmAEStomps = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: PmNamespace,
		Name:      "stomps",
		Help:      "count of stomped versions, labeled by peer and region",
	}, []string{"peer", "region"})
}

func registerMetrics() error {
	if err := prometheus.Register(PmPuts); err != nil {
		log.Debug().Err(err).Msg("unable to register PmPuts")
		return err
	}
	if err := prometheus.Register(PmGets); err != nil {
		log.Debug().Err(err).Msg("unable to register PmGets")
		return err
	}
	if err := prometheus.Register(PmDels); err != nil {
		log.Debug().Err(err).Msg("unable to register PmDels")
		return err
	}
	if err := prometheus.Register(PmIters); err != nil {
		log.Debug().Err(err).Msg("unable to register PmIters")
		return err
	}
	// if err := prometheus.Register(PmObjects); err != nil {
	// 	log.Debug().Err(err).Msg("unable to register PmObjects")
	// 	return err
	// }
	// if err := prometheus.Register(PmTombstones); err != nil {
	// 	log.Debug().Err(err).Msg("unable to register PmTombstones")
	// 	return err
	// }
	if err := prometheus.Register(PmLatency); err != nil {
		log.Debug().Err(err).Msg("unable to register PmLatency")
		return err
	}
	if err := prometheus.Register(PmAESyncs); err != nil {
		log.Debug().Err(err).Msg("unable to register PmAESyncs")
	}
	if err := prometheus.Register(PmAESyncLatency); err != nil {
		log.Debug().Err(err).Msg("unable to register PmAESyncLatency")
	}
	if err := prometheus.Register(PmAEPulls); err != nil {
		log.Debug().Err(err).Msg("unable to register PmAEPulls")
		return err
	}
	if err := prometheus.Register(PmAEPushes); err != nil {
		log.Debug().Err(err).Msg("unable to register PmAEPushes")
		return err
	}
	if err := prometheus.Register(PmAEPushVSPull); err != nil {
		log.Debug().Err(err).Msg("unable to register PmAEPushVSPull")
		return err
	}
	if err := prometheus.Register(PmAEStomps); err != nil {
		log.Debug().Err(err).Msg("unable to register PmAEStomps")
		return err
	}
	return nil
}
