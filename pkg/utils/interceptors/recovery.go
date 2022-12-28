package interceptors

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Recovers a unary handler from a panic and converts the error into an internal error.
func UnaryRecovery() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		panicked := true
		defer func() {
			// NOTE: recover only works for the current go routine so panics in any
			// go routine launched by the handler will not be recovered by this function
			if r := recover(); r != nil || panicked {
				log.WithLevel(zerolog.PanicLevel).
					Err(fmt.Errorf("%v", r)).
					Bool("panicked", panicked).
					Str("stack_trace", string(debug.Stack())).
					Msg("grpc server has recovered from a panic")
				err = status.Error(codes.Internal, "an unhandled exception occurred")
			}
		}()

		rep, err := handler(ctx, req)
		panicked = false
		return rep, err
	}
}

// Recovers a stream handler from a panic and converts the error into an internal error.
func StreamRecovery() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		panicked := true

		defer func() {
			// NOTE: recover only works for the current go routine so panics in any
			// go routine launched by the handler will not be recovered by this function
			if r := recover(); r != nil || panicked {
				log.WithLevel(zerolog.PanicLevel).
					Err(fmt.Errorf("%v", r)).
					Bool("panicked", panicked).
					Str("stack_trace", string(debug.Stack())).
					Msg("grpc server has recovered from a panic")
				err = status.Error(codes.Internal, "an unhandled exception occurred")
			}
		}()

		err = handler(srv, stream)
		panicked = false
		return err
	}
}
