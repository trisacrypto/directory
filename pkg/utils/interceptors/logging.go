package interceptors

import (
	"context"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryLogging handles generic logging of RPC methods to zerolog.
func UnaryLogging() grpc.UnaryServerInterceptor {
	version := pkg.Version()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		// Parse the method for tags
		service, method := ParseMethod(info.FullMethod)

		// Handle the request, tracing how long it takes.
		start := time.Now()
		rep, err := handler(ctx, req)
		duration := time.Since(start)
		code := status.Code(err)

		// Prepare log context for logging
		logctx := log.With().
			Str("type", "unary").
			Str("service", service).
			Str("method", method).
			Str("version", version).
			Uint32("code", uint32(code)).
			Dur("duration", duration).
			Logger()

		// Log based on the error code
		switch code {
		case codes.OK:
			logctx.Info().Msg(info.FullMethod)
		case codes.Unknown:
			logctx.Error().Err(err).Msgf("unknown error handling %s", info.FullMethod)
		case codes.DeadlineExceeded, codes.Canceled, codes.Aborted, codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.Unauthenticated:
			logctx.Info().Err(err).Msg(info.FullMethod)
		case codes.PermissionDenied, codes.ResourceExhausted, codes.FailedPrecondition, codes.OutOfRange, codes.Unavailable:
			logctx.Warn().Err(err).Msg(info.FullMethod)
		case codes.Unimplemented, codes.Internal, codes.DataLoss:
			logctx.Error().Err(err).Str("full_method", info.FullMethod).Msg(err.Error())
		default:
			logctx.Error().Err(err).Msgf("unhandled error code %s: %s", code, info.FullMethod)
		}

		return rep, err
	}
}

// StreamLogging handles generic logging of RPC methods to zerolog.
func StreamLogging() grpc.StreamServerInterceptor {
	version := pkg.Version()
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		// Parse the method for tags
		service, method := ParseMethod(info.FullMethod)

		// Handle the request and trace how long the request takes.
		start := time.Now()
		err = handler(srv, stream)
		duration := time.Since(start)
		code := status.Code(err)

		// Prepare log context for logging
		logctx := log.With().
			Str("type", "stream").
			Str("stream_type", StreamType(info)).
			Str("service", service).
			Str("method", method).
			Str("version", version).
			Uint32("code", uint32(code)).
			Dur("duration", duration).
			Logger()

		switch code {
		case codes.OK:
			logctx.Info().Msg(info.FullMethod)
		case codes.Unknown:
			logctx.Error().Err(err).Msgf("unknown error handling %s", info.FullMethod)
		case codes.Canceled, codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.Unauthenticated:
			logctx.Info().Err(err).Msg(info.FullMethod)
		case codes.DeadlineExceeded, codes.PermissionDenied, codes.ResourceExhausted, codes.FailedPrecondition, codes.Aborted, codes.OutOfRange, codes.Unavailable:
			logctx.Warn().Err(err).Msg(info.FullMethod)
		case codes.Unimplemented, codes.Internal, codes.DataLoss:
			logctx.Error().Err(err).Str("full_method", info.FullMethod).Msg(err.Error())
		default:
			logctx.Error().Err(err).Msgf("unhandled error code %s: %s", code, info.FullMethod)
		}
		return err
	}
}

func ParseMethod(method string) (string, string) {
	method = strings.TrimPrefix(method, "/") // remove leading slash
	if i := strings.Index(method, "/"); i >= 0 {
		return method[:i], method[i+1:]
	}
	return "unknown", "unknown"
}

func StreamType(info *grpc.StreamServerInfo) string {
	if !info.IsClientStream && !info.IsServerStream {
		return "unary"
	}
	if info.IsClientStream && !info.IsServerStream {
		return "client_stream"
	}
	if !info.IsClientStream && info.IsServerStream {
		return "server_stream"
	}
	return "bidirectional"
}
