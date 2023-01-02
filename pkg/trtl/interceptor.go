package trtl

import (
	"github.com/trisacrypto/directory/pkg/trtl/metrics"
	"github.com/trisacrypto/directory/pkg/utils/interceptors"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
	"google.golang.org/grpc"
)

const statusMethod = "/trtl.v1.Trtl/Status"

// Prepares the interceptors (middleware) for the unary RPC endpoints of the server.
// The first interceptor will be the outer-most handler and the last will be the
// inner-most wrapper around the final handler. All unary interceptors returned by this
// method should be chained using grpc.ChainUnaryInterceptor().
func (t *Server) UnaryInterceptors() []grpc.UnaryServerInterceptor {
	// Prepare Sentry configuration
	t.conf.Sentry.Repanic = true

	// If we're in maintenance mode, only return maintenance and recovery
	if t.conf.Maintenance {
		return []grpc.UnaryServerInterceptor{
			interceptors.UnaryRecovery(),
			interceptors.UnaryMaintenance(statusMethod),
		}
	}

	// Return Unary interceptors
	opts := []grpc.UnaryServerInterceptor{
		interceptors.UnaryLogging(),
		interceptors.UnaryRecovery(),
		sentry.UnaryInterceptor(t.conf.Sentry),
	}

	// If monitoring is enabled prepend it to the list of interceptors
	if t.conf.Metrics.Enabled {
		opts = append([]grpc.UnaryServerInterceptor{metrics.UnaryMonitoring()}, opts...)
	}

	// If mTLS is enabled append it to the list of interceptors
	if !t.conf.MTLS.Insecure {
		opts = append(opts, interceptors.UnaryMTLS())
	}
	return opts
}

// Prepares the interceptors (middleware) for the stream RPC endpoints of the server.
// The first interceptor will be the outer-most handler and the last will be the
// inner-most wrapper around the final handler. All stream interceptors returned by this
// method should be chained using grpc.ChainStreamInterceptor().
func (t *Server) StreamInterceptors() []grpc.StreamServerInterceptor {
	// Prepare Sentry configuration
	t.conf.Sentry.Repanic = true

	// If we're in maintenance mode, only return maintenance and recovery
	if t.conf.Maintenance {
		return []grpc.StreamServerInterceptor{
			interceptors.StreamRecovery(),
			interceptors.StreamMaintenance(statusMethod),
		}
	}

	// Return Stream interceptors
	opts := []grpc.StreamServerInterceptor{
		interceptors.StreamLogging(),
		interceptors.StreamRecovery(),
		sentry.StreamInterceptor(t.conf.Sentry),
	}

	// If monitoring is enabled prepend it to the list of interceptors
	if t.conf.Metrics.Enabled {
		opts = append([]grpc.StreamServerInterceptor{metrics.StreamMonitoring()}, opts...)
	}

	// If mTLS is enabled append it to the list of interceptors
	if !t.conf.MTLS.Insecure {
		opts = append(opts, interceptors.StreamMTLS())
	}
	return opts
}
