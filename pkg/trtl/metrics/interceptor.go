package metrics

import (
	"context"
	"time"

	"github.com/trisacrypto/directory/pkg/utils/interceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryMonitoring is an interceptor that handles the generic gRPC monitoring prometheus
// metrics for trtl, tracking each RPC by service and method name. Unary RPC handlers
// should update the namespace on the context for this interceptor to correctly track
// what is happening on a per-namespace basis.
func UnaryMonitoring() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		// Parse the service and method from the RPC info
		service, method := interceptors.ParseMethod(info.FullMethod)
		PmRPCStarted.WithLabelValues("unary", service, method).Inc()

		// Update the context with the shared namespace
		ctx = WithSharedNamespace(ctx)

		// Execute the handler tracking how long it takes
		start := time.Now()
		rep, err := handler(ctx, req)
		duration := time.Since(start)

		// Extract the code and namespace from the request
		code := status.Code(err)
		namespace := GetNamespace(ctx)

		PmRPCHandled.WithLabelValues(namespace, "unary", service, method, code.String()).Inc()
		PmRPCUnaryLatency.WithLabelValues(namespace, "unary", service, method).Observe(duration.Seconds())
		return rep, err
	}
}

// StreamMonitoring is an interceptor that handels generic gRPC streaming monitoring,
// updating prometheus metrics for trtl and tracking each RPC by service and method.
// Streaming RPC handlers should update the namespace on the context for this
// interceptor to correctly track what is happening on a per-namespace basis.
func StreamMonitoring() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		// Parse the service and method from the RPC info
		service, method := interceptors.ParseMethod(info.FullMethod)
		stream_type := interceptors.StreamType(info)
		PmRPCStarted.WithLabelValues(stream_type, service, method).Inc()

		// Wrap stream with a monitored stream handler
		ctx := WithSharedNamespace(stream.Context())
		stream = &MonitoredStream{stream, 0, 0, ctx}

		// Handle the request and track how long it takes
		start := time.Now()
		err = handler(srv, stream)
		duration := time.Since(start)

		// Extract the code and namespace from the request
		code := status.Code(err)
		namespace := GetNamespace(ctx)

		PmRPCHandled.WithLabelValues(namespace, stream_type, service, method, code.String()).Inc()
		PmRPCStreamDuration.WithLabelValues(namespace, stream_type, service, method).Observe(duration.Seconds())
		PmMsgsPerStream.WithLabelValues(namespace, stream_type, service, method).Observe(stream.(*MonitoredStream).Msgs())
		return err
	}
}

// MonitoredStream wraps a grpc.ServerStream allowing it to increment Sent and Recv
// message counters when they are called by the application.
type MonitoredStream struct {
	grpc.ServerStream
	sent uint64
	recv uint64
	ctx  context.Context
}

// Increment the number of sent messages if there is no error on Send.
func (s *MonitoredStream) SendMsg(m interface{}) (err error) {
	if err = s.ServerStream.SendMsg(m); err == nil {
		s.sent++
	}
	return err
}

// Increment the number of received messages if there is no error on Recv.
func (s *MonitoredStream) RecvMsg(m interface{}) (err error) {
	if err = s.ServerStream.RecvMsg(m); err == nil {
		s.recv++
	}
	return err
}

func (s *MonitoredStream) Msgs() float64 {
	return float64(s.sent + s.recv)
}

func (s *MonitoredStream) Context() context.Context {
	return s.ctx
}

// SharedNamespace is a bit of a hack, using a pointer we allow the child handler to
// update the namespace so that the parent interceptor can use the namespace as a label
// for monitoring. This is not a correct or standard way to use contexts and it is not
// thread-safe even though contexts should be.
type SharedNamespace struct {
	Namespace string
}

type MetricsKey uint8

var NamespaceKey MetricsKey = 1

// WithSharedNamespace updates the context with a shared namespace value.
func WithSharedNamespace(ctx context.Context) context.Context {
	shared := &SharedNamespace{Namespace: "unknown"}
	return context.WithValue(ctx, NamespaceKey, shared)
}

// UpdateNamespace updates the context with the shared namespace value if available.
func UpdateNamespace(ctx context.Context, namespace string) {
	val := ctx.Value(NamespaceKey)
	if val == nil {
		return
	}

	shared, ok := val.(*SharedNamespace)
	if !ok {
		return
	}

	shared.Namespace = namespace
}

// GetNamespace returns the namespace from the shared namespace context value.
func GetNamespace(ctx context.Context) string {
	val := ctx.Value(NamespaceKey)
	if val == nil {
		return ""
	}

	shared, ok := val.(*SharedNamespace)
	if !ok {
		return ""
	}

	return shared.Namespace
}
