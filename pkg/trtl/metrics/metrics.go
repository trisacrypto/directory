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
	// Basic RPC Metrics
	PmPuts       *prometheus.CounterVec   // count of trtl Puts per namespace
	PmGets       *prometheus.CounterVec   // count of trtl Gets per namespace
	PmDels       *prometheus.CounterVec   // count of trtl Deletes per namespace
	PmIters      *prometheus.CounterVec   // count of trtl Iters per namespace
	PmRPCLatency *prometheus.HistogramVec // the time it is taking for successful RPC calls to complete, labeled by RPC type, success, and failure
	// PmObjects    *prometheus.CounterVec   // count of objects being managed by trtl, by namespace
	// PmTombstones *prometheus.CounterVec   // count of tombstones per namespace; increases on delete, decrease on overwrite of tombstone

	// Anti-Entropy Metrics
	PmAESyncs         *prometheus.CounterVec   // count of anti entropy sessions per peer, per region, and by perspective (initiator/remote)
	PmAESyncLatency   *prometheus.HistogramVec // duration of anti entropy sessions (initiator perspective), by peer and region
	PmAEPhase1Latency *prometheus.HistogramVec // time phase 1 of anti-entropy is taking from the perspective of the initiator, by peer
	PmAEPhase2Latency *prometheus.HistogramVec // time phase 2 of anti-entropy is taking from the perspective of the remote, by peer
	PmAEVersions      *prometheus.HistogramVec // count of all observed versions, per peer and region
	PmAEUpdates       *prometheus.HistogramVec // pushed objects during anti entropy, by peer and region
	PmAERepairs       *prometheus.HistogramVec // pulled objects during anti entropy, by peer and region
	PmAEStomps        *prometheus.CounterVec   // count of stomped versions, per peer and region
	PmAESkips         *prometheus.CounterVec   // count of skipped versions, per peer and region

)

// A MetricsService manages Prometheus metrics
type MetricsService struct {
	srv *http.Server
}

// New creates a metrics service and also initializes all of the prometheus metrics.
// The trtl server *must* create the metrics service by calling New before any
// metrics are logged to Prometheus. Even in the case of tests, the metrics service
// must be created before the tests can be run.
func New() (*MetricsService, error) {
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

	PmRPCLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: PmNamespace,
		Name:      "latency",
		Help:      "time to RPC call completion, labeled by RPC (Put, Get, Delete, Iter)",
	}, []string{"call"})

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

	PmAESyncs = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: PmNamespace,
		Name:      "syncs",
		Help:      "the count of anti-entropy sessions, labeled by peer, region, and perspective",
	}, []string{"peer", "region", "perspective"})

	PmAESyncLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: PmNamespace,
		Name:      "sync_latency",
		Help:      "total duration of anti-entropy (originator perspective), labeled by peer and region",
	}, []string{"peer", "region"})

	PmAEPhase1Latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: PmNamespace,
		Name:      "phase1_latency",
		Help:      "duration of anti-entropy phase 1 (originator perspective), labeled by peer",
	}, []string{"peer"})

	PmAEPhase2Latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: PmNamespace,
		Name:      "phase2_latency",
		Help:      "duration of anti-entropy phase 2 (remote perspective), labeled by peer",
	}, []string{"peer"})

	PmAEVersions = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: PmNamespace,
		Name:      "versions",
		Help:      "count of all observed versions, labeled by peer, region, and perspective",
	}, []string{"peer", "region", "perspective"})

	PmAERepairs = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: PmNamespace,
		Name:      "pulls",
		Help:      "pulled objects during anti entropy, labeled by peer, region, and perspective",
	}, []string{"peer", "region", "perspective"})

	PmAEUpdates = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: PmNamespace,
		Name:      "pushes",
		Help:      "pushed objects during anti entropy, labeled by peer, region and perspective",
	}, []string{"peer", "region", "perspective"})

	PmAEStomps = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: PmNamespace,
		Name:      "stomps",
		Help:      "count of stomped versions, labeled by peer and region",
	}, []string{"peer", "region"})

	PmAESkips = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: PmNamespace,
		Name:      "skips",
		Help:      "count of skipped versions, labeled by peer and region",
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
	if err := prometheus.Register(PmRPCLatency); err != nil {
		log.Debug().Err(err).Msg("unable to register PmLatency")
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

	if err := prometheus.Register(PmAESyncs); err != nil {
		log.Debug().Err(err).Msg("unable to register PmAESyncs")
	}
	if err := prometheus.Register(PmAESyncLatency); err != nil {
		log.Debug().Err(err).Msg("unable to register PmAESyncLatency")
	}
	if err := prometheus.Register(PmAEPhase1Latency); err != nil {
		log.Debug().Err(err).Msg("unable to register PmAEPhase1Latency")
	}
	if err := prometheus.Register(PmAEPhase2Latency); err != nil {
		log.Debug().Err(err).Msg("unable to register PmAEPhase2Latency")
	}
	if err := prometheus.Register(PmAEVersions); err != nil {
		log.Debug().Err(err).Msg("unable to register PmAEVersions")
		return err
	}
	if err := prometheus.Register(PmAERepairs); err != nil {
		log.Debug().Err(err).Msg("unable to register PmAEPulls")
		return err
	}
	if err := prometheus.Register(PmAEUpdates); err != nil {
		log.Debug().Err(err).Msg("unable to register PmAEPushes")
		return err
	}
	if err := prometheus.Register(PmAEStomps); err != nil {
		log.Debug().Err(err).Msg("unable to register PmAEStomps")
		return err
	}
	if err := prometheus.Register(PmAESkips); err != nil {
		log.Debug().Err(err).Msg("unable to register PmAESkips")
		return err
	}
	return nil
}
