package trtl

import (
	"github.com/trisacrypto/directory/pkg/utils/interceptors"
	sentryutil "github.com/trisacrypto/directory/pkg/utils/sentry"
	"google.golang.org/grpc"
)

const statusMethod = "/trtl.v1.Trtl/Status"

// Prepares the interceptors (middleware) for the unary RPC endpoints of the server.
// The first interceptor will be the outer-most handler and the last will be the
// inner-most wrapper around the final handler. All unary interceptors returned by this
// method should be changed using grpc.ChaingUnaryInterceptor().
func (t *Server) UnaryInterceptors() []grpc.UnaryServerInterceptor {
	// Prepare Sentry configuration
	// NOTE: this will override any user-configured settings
	t.conf.Sentry.Service = "trtl"
	t.conf.Sentry.Repanic = true

	// If we're in maintenance mode, only return maintenance and recovery
	if t.conf.Maintenance {
		return []grpc.UnaryServerInterceptor{
			interceptors.UnaryRecovery(),
			interceptors.UnaryMaintenance(statusMethod),
		}
	}

	// Return Unary interceptors
	return []grpc.UnaryServerInterceptor{
		interceptors.UnaryLogging(),
		interceptors.UnaryRecovery(),
		sentryutil.UnaryInterceptor(t.conf.Sentry),
		interceptors.UnaryMTLS(),
	}
}

// Prepares the interceptors (middleware) for the stream RPC endpoints of the server.
// The first interceptor will be the outer-most handler and the last will be the
// inner-most wrapper around the final handler. All stream interceptors returned by this
// method should be changed using grpc.ChaingStreamInterceptor().
func (t *Server) StreamInterceptors() []grpc.StreamServerInterceptor {
	// Prepare Sentry configuration
	// NOTE: this will override any user-configured settings
	t.conf.Sentry.Service = "trtl"
	t.conf.Sentry.Repanic = true

	// If we're in maintenance mode, only return maintenance and recovery
	if t.conf.Maintenance {
		return []grpc.StreamServerInterceptor{
			interceptors.StreamRecovery(),
			interceptors.StreamMaintenance(statusMethod),
		}
	}

	// Return Unary interceptors
	return []grpc.StreamServerInterceptor{
		interceptors.StreamLogging(),
		interceptors.StreamRecovery(),
		sentryutil.StreamInterceptor(t.conf.Sentry),
		interceptors.StreamMTLS(),
	}
}
