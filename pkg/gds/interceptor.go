package gds

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	members "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
	api "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) unaryInterceptor(ctx context.Context, in interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (out interface{}, err error) {
	// Track how long the method takes to execute.
	start := time.Now()
	panicked := true

	// Set the service tag
	if s.conf.Sentry.UseSentry() {
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
		sentry.CurrentHub().Scope().SetTag("service", service)
	}

	// Recover from panics in the handler.
	// See: https://github.com/grpc-ecosystem/go-grpc-middleware/blob/4705cb37b9857ad51b4c96ff5a2f3c60afe442cf/recovery/interceptors.go#L21-L37
	defer func() {
		if r := recover(); r != nil || panicked {
			if s.conf.Sentry.UseSentry() {
				sentry.CurrentHub().Recover(r)
			}
			log.WithLevel(zerolog.PanicLevel).
				Err(fmt.Errorf("%v", r)).
				Str("stack_trace", string(debug.Stack())).
				Msg("grpc server has recovered from a panic")
			err = status.Error(codes.Internal, "an unhandled exception occurred")
		}
	}()

	// Check if we're in maintenance mode - status method should still return a full response.
	if s.conf.Maintenance && info.FullMethod != "/trisa.gds.api.v1beta1.TRISADirectory/Status" {
		err = status.Error(codes.Unavailable, "the GDS service is currently in maintenance mode")
		log.Trace().Err(err).Str("method", info.FullMethod).Msg("gds service unavailable during maintenance mode")

		panicked = false
		return nil, err
	}

	// Call the handler to finalize the request and get the response.
	var span *sentry.Span
	if s.conf.Sentry.UsePerformanceTracking() {
		span = sentry.StartSpan(ctx, "grpc", sentry.TransactionName(info.FullMethod))
	}
	out, err = handler(ctx, in)
	if s.conf.Sentry.UsePerformanceTracking() {
		span.Finish()
	}

	// Log with zerolog - checkout grpclog.LoggerV2 for default logging.
	// TODO: add remote peer information if using mTLS
	log.Debug().
		Err(err).
		Str("method", info.FullMethod).
		Str("latency", time.Since(start).String()).
		Msg("gRPC request complete")

	panicked = false
	return out, err
}

// The streamInterceptor intercepts incoming gRPC streaming requests and adds remote
// peer information to the context, performing maintenance mode checks and panic recovery.
func (s *Service) streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	// Track how long the method takes to execute.
	start := time.Now()
	panicked := true

	// Set the service tag
	// TODO: set the type of the stream server to determine if the tag is gds or members
	// Currently there are no streaming RPCs in our defined services (unreachable code)
	if s.conf.Sentry.UseSentry() {
		var service string
		switch srv.(type) {
		default:
			log.WithLevel(zerolog.PanicLevel).Err(fmt.Errorf("unknown service type: %T", srv))
			return status.Error(codes.Unimplemented, "unknown service type for request")
		}
		sentry.CurrentHub().Scope().SetTag("service", service)
	}

	// Recover from panics in the handler
	defer func() {
		if r := recover(); r != nil || panicked {
			if s.conf.Sentry.UseSentry() {
				sentry.CurrentHub().Recover(r)
			}
			log.WithLevel(zerolog.PanicLevel).
				Err(fmt.Errorf("%v", r)).
				Str("stack_trace", string(debug.Stack())).
				Msg("trtl server has recovered from a panic")
			err = status.Error(codes.Internal, "an unhandled exception occurred")
		}
	}()

	// Check if we're in maintenance mode -- no streaming method should be available
	if s.conf.Maintenance {
		err = status.Error(codes.Unavailable, "the GDS service is currently in maintenance mode")
		log.Trace().Err(err).Str("method", info.FullMethod).Msg("gds service unavailable during maintenance mode")

		panicked = false
		return err
	}

	// Call the handler to execute the stream RPC
	// NOTE: sentry performance tracking is not valid here since streams can take an
	// arbitrarily long time to complete and minimizing latency is not necessarily desirable.
	err = handler(srv, ss)

	// Log with zerolog - check grpclog.LoggerV2 for default logging
	log.Debug().
		Err(err).
		Str("method", info.FullMethod).
		Str("duration", time.Since(start).String()).
		Msg("grpc stream request complete")

	panicked = false
	return err
}
