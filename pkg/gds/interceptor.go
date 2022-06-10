package gds

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) serverInterceptor(ctx context.Context, in interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (out interface{}, err error) {
	// Track how long the method takes to execute.
	start := time.Now()
	panicked := true

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
		return nil, status.Error(codes.Unavailable, "the GDS service is currently in maintenance mode")
	}

	// Call the handler to finalize the request and get the response.
	var span *sentry.Span
	if s.conf.Sentry.UsePerformanceTracking() {
		span = sentry.StartSpan(ctx, info.FullMethod)
	}
	out, err = handler(ctx, in)
	if s.conf.Sentry.UsePerformanceTracking() {
		span.Finish()
	}
	panicked = false

	// Log with zerolog - checkout grpclog.LoggerV2 for default logging.
	// TODO: add remote peer information if using mTLS
	log.Debug().Str("method", info.FullMethod).Str("latency", time.Since(start).String()).Err(err).Msg("gRPC request complete")
	return out, err
}
