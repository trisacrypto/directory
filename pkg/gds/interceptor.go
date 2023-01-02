package gds

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	members "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
	"github.com/trisacrypto/directory/pkg/utils/interceptors"
	sentryutil "github.com/trisacrypto/directory/pkg/utils/sentry"
	api "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const statusMethod = "/trisa.gds.api.v1beta1.TRISADirectory/Status"

// Prepares the interceptors (middleware) for the unary RPC endpoints of the server.
// The first interceptor will be the outer-most handler and the last will be the
// inner-most wrapper around the final handler. All unary interceptors returned by this
// method should be chained using grpc.ChainUnaryInterceptor().
func (s *Service) UnaryInterceptors() []grpc.UnaryServerInterceptor {
	// Prepare Sentry configuration
	s.conf.Sentry.Repanic = true

	// If we're in maintenance mode, only return maintenance and recovery
	if s.conf.Maintenance {
		return []grpc.UnaryServerInterceptor{
			interceptors.UnaryRecovery(),
			interceptors.UnaryMaintenance(statusMethod),
		}
	}

	// Return Unary interceptors
	opts := []grpc.UnaryServerInterceptor{
		interceptors.UnaryLogging(),
		interceptors.UnaryRecovery(),
		sentryutil.UnaryInterceptor(s.conf.Sentry),
		ServiceTag(),
	}

	return opts
}

// Prepares the interceptors (middleware) for the stream RPC endpoints of the server.
// The first interceptor will be the outer-most handler and the last will be the
// inner-most wrapper around the final handler. All stream interceptors returned by this
// method should be chained using grpc.ChainStreamInterceptor().
func (s *Service) StreamInterceptors() []grpc.StreamServerInterceptor {
	// Prepare Sentry configuration
	s.conf.Sentry.Repanic = true

	// If we're in maintenance mode, only return maintenance and recovery
	if s.conf.Maintenance {
		return []grpc.StreamServerInterceptor{
			interceptors.StreamRecovery(),
			interceptors.StreamMaintenance(statusMethod),
		}
	}

	// Return Stream interceptors
	opts := []grpc.StreamServerInterceptor{
		interceptors.StreamLogging(),
		interceptors.StreamRecovery(),
		sentryutil.StreamInterceptor(s.conf.Sentry),
	}
	return opts
}

func ServiceTag() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, in interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (out interface{}, err error) {
		hub := sentry.GetHubFromContext(ctx)
		if hub == nil {
			hub = sentry.CurrentHub().Clone()
			ctx = sentry.SetHubOnContext(ctx, hub)
		}

		var service string
		switch info.Server.(type) {
		case api.TRISADirectoryServer:
			service = "gds"
		case members.TRISAMembersServer:
			service = "members"
		default:
			log.WithLevel(zerolog.PanicLevel).Err(fmt.Errorf("unknown service type: %T", info.Server))
			return nil, status.Error(codes.Unimplemented, "unknown service type for request")
		}

		hub.Scope().SetTag("service", service)
		return handler(ctx, in)
	}
}
