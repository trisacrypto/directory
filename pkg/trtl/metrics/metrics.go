package metrics

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/trtl/config"
)

// Prometheus namespaces for the collectors defined in this package.
const (
	PmNamespaceTrtl = "trtl"
	PmNamespaceGRPC = "grpc"
)

// All trtl specific collectors for observability are defined here.
var (
	// Basic RPC Metrics
	PmRPCStarted        *prometheus.CounterVec   // RPCs started by method, namespace
	PmRPCHandled        *prometheus.CounterVec   // RPCs completed by method, namespace, and code
	PmRPCUnaryLatency   *prometheus.HistogramVec // the time it is taking for successful unary RPC calls to complete, labeled by RPC type, namespace, and code
	PmRPCStreamDuration *prometheus.HistogramVec // the time it is taking for successful streaming RPC calls to complete, labeled by type, namespace, and code
	PmMsgsPerStream     *prometheus.HistogramVec // the number of messages sent and recv per streaming RPC, labeled by type, namespace, and code

	// Storage Metrics
	// TODO: add version metrics
	PmTrtlReads         *prometheus.CounterVec   // number of reads, e.g. Get and Iter to the embedded database, by namespace
	PmTrtlBytesRead     *prometheus.CounterVec   // number of bytes read by trtl operations by namespace
	PmTrtlWrites        *prometheus.CounterVec   // number of writes, e.g. Puts and Deletes to the embedded database, by namespace
	PmTrtlBytesWritten  *prometheus.CounterVec   // number of bytes written by trtl operations by namespace
	PmObjectSize        *prometheus.HistogramVec // average size in bytes of objects stored in trtl, by namespace
	PmDatabaseSize      *prometheus.GaugeVec     // current size in bytes of all objects in the database, by namespace
	PmCurrentObjects    *prometheus.GaugeVec     // current number of objects in the database, by namespace
	PmCurrentTombstones *prometheus.GaugeVec     // current number of tombstones int he database, by namespace

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

// Ensure that the collectors are only registered once even if multiple metrics servers
// are created and initialized from external packages.
var (
	register    sync.Once
	registerErr error
)

// A MetricsService manages Prometheus metrics
type MetricsService struct {
	srv *http.Server
	cfg config.MetricsConfig
}

// New creates a metrics service and also initializes all of the prometheus metrics.
// The trtl server *must* create the metrics service by calling New before any
// metrics are logged to Prometheus.
func New(conf config.MetricsConfig) (*MetricsService, error) {
	if err := RegisterMetrics(); err != nil {
		return nil, err
	}

	// Setup the prometheus handler and collectors server
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return &MetricsService{
		cfg: conf,
		srv: &http.Server{
			Addr:         conf.Addr,
			Handler:      mux,
			ErrorLog:     nil,
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}, nil
}

// Serve serves the Prometheus metrics
func (m *MetricsService) Serve() error {
	// If metrics are not enabled return without starting the server.
	if !m.cfg.Enabled {
		return nil
	}

	// Serve the metrics server in its own go routine
	go func() {
		if err := m.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Err(err).Msg("metrics server shutdown prematurely")
		}
	}()

	// Log the metrics service starting up
	log.Info().Str("addr", fmt.Sprintf("http://%s/metrics", m.cfg.Addr)).Msg("metrics server started and ready for prometheus collector")
	return nil
}

// Gracefully shutdown the Prometheus metrics service
func (m *MetricsService) Shutdown(ctx context.Context) error {
	// If metrics are not enabled, return without shutting down
	if !m.cfg.Enabled {
		return nil
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	if err := m.srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

// Initializes and registers the metrics collectors in Prometheus. This function can
// safely be called multiple times and the collectors will only be registered once. This
// method can be used prior to tests to ensure that there are no nil panics for handlers
// that make use of the collectors; otherwise it will be called when creating a new
// metrics server in preparation for exposing application metrics.
func RegisterMetrics() error {
	register.Do(func() {
		registerErr = registerMetrics()
	})
	return registerErr
}

func registerMetrics() error {
	// Track all collectors to make it easier to register them after initialization
	collectors := make([]prometheus.Collector, 0, 22)

	// Basic RPC Metrics
	PmRPCStarted = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: PmNamespaceGRPC,
		Name:      "server_started",
		Help:      "count the number of RPCs started on the server",
	}, []string{"type", "service", "method"})
	collectors = append(collectors, PmRPCStarted)

	PmRPCHandled = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: PmNamespaceGRPC,
		Name:      "server_handled",
		Help:      "count the number of RPCs completed on the server",
	}, []string{"namespace", "type", "service", "method", "code"})
	collectors = append(collectors, PmRPCHandled)

	PmRPCUnaryLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: PmNamespaceGRPC,
		Name:      "unary_latency",
		Help:      "response latency (in seconds) of the application handler for unary rpcs",
		Buckets:   prometheus.ExponentialBuckets(5, 2, 12),
	}, []string{"namespace", "type", "service", "method"})
	collectors = append(collectors, PmRPCUnaryLatency)

	PmRPCStreamDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: PmNamespaceGRPC,
		Name:      "streaming_duration",
		Help:      "durations of streaming application handslers",
		Buckets:   prometheus.ExponentialBuckets(5, 2, 12),
	}, []string{"namespace", "type", "service", "method"})
	collectors = append(collectors, PmRPCStreamDuration)

	PmMsgsPerStream = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: PmNamespaceGRPC,
		Name:      "messages_per_stream",
		Help:      "number of messages sent and recv per streaming rpc handler",
	}, []string{"namespace", "type", "service", "method"})
	collectors = append(collectors, PmMsgsPerStream)

	// Storage Metrics
	PmTrtlReads = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: PmNamespaceTrtl,
		Name:      "reads",
		Help:      "the number of reads to the embedded database (e.g. Get and Iter)",
	}, []string{"namespace"})
	collectors = append(collectors, PmTrtlReads)

	PmTrtlBytesRead = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: PmNamespaceTrtl,
		Name:      "bytes_read",
		Help:      "the number of bytes read by trtl operations",
	}, []string{"namespace"})
	collectors = append(collectors, PmTrtlBytesRead)

	PmTrtlWrites = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: PmNamespaceTrtl,
		Name:      "writes",
		Help:      "the number of writes to the embedded database (e.g. Put and Delete)",
	}, []string{"namespace"})
	collectors = append(collectors, PmTrtlWrites)

	PmTrtlBytesWritten = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: PmNamespaceTrtl,
		Name:      "bytes_written",
		Help:      "the number of bytes written by trtl operations",
	}, []string{"namespace"})
	collectors = append(collectors, PmTrtlBytesWritten)

	PmObjectSize = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: PmNamespaceTrtl,
		Name:      "object_size",
		Help:      "size in bytes of each version saved in the trtl database",
	}, []string{"namespace"})
	collectors = append(collectors, PmObjectSize)

	PmDatabaseSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: PmNamespaceTrtl,
		Name:      "database_size",
		Help:      "current size in bytes of all objects in the database",
	}, []string{"namespace"})
	collectors = append(collectors, PmDatabaseSize)

	PmCurrentObjects = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: PmNamespaceTrtl,
		Name:      "objects",
		Help:      "current number of objects in the database",
	}, []string{"namespace"})
	collectors = append(collectors, PmCurrentObjects)

	PmCurrentTombstones = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: PmNamespaceTrtl,
		Name:      "tombstones",
		Help:      "current number of tombstones in the database",
	}, []string{"namespace"})
	collectors = append(collectors, PmCurrentTombstones)

	// Anti-Entropy Metrics
	PmAESyncs = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: PmNamespaceTrtl,
		Name:      "syncs",
		Help:      "the count of anti-entropy sessions, labeled by peer, region, and perspective",
	}, []string{"peer", "region", "perspective"})
	collectors = append(collectors, PmAESyncs)

	PmAESyncLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: PmNamespaceTrtl,
		Name:      "sync_latency",
		Help:      "total duration of anti-entropy (originator perspective), labeled by peer and region",
		Buckets:   prometheus.LinearBuckets(10, 10, 50),
	}, []string{"peer", "region"})
	collectors = append(collectors, PmAESyncLatency)

	PmAEPhase1Latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: PmNamespaceTrtl,
		Name:      "phase1_latency",
		Help:      "duration of anti-entropy phase 1 (originator perspective), labeled by peer",
		Buckets:   prometheus.LinearBuckets(1, 10, 50),
	}, []string{"peer"})
	collectors = append(collectors, PmAEPhase1Latency)

	PmAEPhase2Latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: PmNamespaceTrtl,
		Name:      "phase2_latency",
		Help:      "duration of anti-entropy phase 2 (remote perspective), labeled by peer",
		Buckets:   prometheus.LinearBuckets(1, 10, 50),
	}, []string{"peer"})
	collectors = append(collectors, PmAEPhase2Latency)

	PmAEVersions = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: PmNamespaceTrtl,
		Name:      "versions",
		Help:      "count of all observed versions, labeled by peer, region, and perspective",
		Buckets:   prometheus.LinearBuckets(10, 2000, 1000),
	}, []string{"peer", "region", "perspective"})
	collectors = append(collectors, PmAEVersions)

	PmAERepairs = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: PmNamespaceTrtl,
		Name:      "pulls",
		Help:      "pulled objects during anti entropy, labeled by peer, region, and perspective",
		Buckets:   prometheus.LinearBuckets(10, 100, 1000),
	}, []string{"peer", "region", "perspective"})
	collectors = append(collectors, PmAERepairs)

	PmAEUpdates = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: PmNamespaceTrtl,
		Name:      "pushes",
		Help:      "pushed objects during anti entropy, labeled by peer, region and perspective",
		Buckets:   prometheus.LinearBuckets(10, 100, 1000),
	}, []string{"peer", "region", "perspective"})
	collectors = append(collectors, PmAEUpdates)

	PmAEStomps = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: PmNamespaceTrtl,
		Name:      "stomps",
		Help:      "count of stomped versions, labeled by peer and region",
	}, []string{"peer", "region"})
	collectors = append(collectors, PmAEStomps)

	PmAESkips = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: PmNamespaceTrtl,
		Name:      "skips",
		Help:      "count of skipped versions, labeled by peer and region",
	}, []string{"peer", "region"})
	collectors = append(collectors, PmAESkips)

	// Register all collectors
	for _, collector := range collectors {
		if err := prometheus.Register(collector); err != nil {
			log.Debug().Err(err).Msg("could not register collector")
			return err
		}
	}
	return nil
}
