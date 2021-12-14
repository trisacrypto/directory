package trtl

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

var (
	pmPuts       *prometheus.CounterVec   // count of trtl Puts per namespace
	pmGets       *prometheus.CounterVec   // count of trtl Gets per namespace
	pmDels       *prometheus.CounterVec   // count of trtl Deletes per namespace
	pmIters      *prometheus.CounterVec   // count of trtl Iters per namespace
	pmObjects    *prometheus.CounterVec   // count of objects being managed by trtl, by namespace
	pmTombstones *prometheus.CounterVec   // count of tombstones per namespace; increases on delete, decrease on overwrite of tombstone
	pmLatency    *prometheus.HistogramVec // the time it is taking for RPC calls to complete, labeled by RPC type, success, and failure
)

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

	pmObjects = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: pmNamespace,
		Name:      "objects",
		Help:      "the count of trtl objects, labeled by namespace",
	}, []string{"namespace"})

	pmTombstones = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: pmNamespace,
		Name:      "tombstones",
		Help:      "the count of tombstones, labeled by namespace",
	}, []string{"namespace"})

	pmLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: pmNamespace,
		Name:      "latency",
		Help:      "time to RPC call completion, labeled by RPC (Put, Get, Delete, Iter)",
	}, []string{"call"})
}

func serveMetrics(metricsAddr string) {
	log.Info().Msg(fmt.Sprintf("serving prometheus metrics at http://%s/metrics", metricsAddr))
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(metricsAddr, nil); err != nil {
		log.Error().Err(err).Str("metricsAddr", metricsAddr).Msg("unable to serve prometheus metrics")
	}
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
	if err := prometheus.Register(pmObjects); err != nil {
		log.Debug().Err(err).Msg("unable to register pmObjects")
		return err
	}
	if err := prometheus.Register(pmTombstones); err != nil {
		log.Debug().Err(err).Msg("unable to register pmTombstones")
		return err
	}
	if err := prometheus.Register(pmLatency); err != nil {
		log.Debug().Err(err).Msg("unable to register pmLatency")
		return err
	}

	return nil
}
