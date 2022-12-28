package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// The maintenance interceptor only allows Status endpoint to be queried and returns a
// service unavailable error otherwise.
func UnaryMaintenance(statusEndpoint string) grpc.UnaryServerInterceptor {
	// This interceptor will supercede all following interceptors.
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		// Allow the Status endpoint through, otherwise return unavailable
		if info.FullMethod == statusEndpoint {
			return handler(ctx, req)
		}
		return nil, status.Error(codes.Unavailable, "service is currently in maintenance mode")
	}
}

// The stream maintenance interceptor simply returns an unavailable error.
func StreamMaintenance(statusEndpoint string) grpc.StreamServerInterceptor {
	// This interceptor will supercede all following interceptors
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		return status.Error(codes.Unavailable, "service is currently in maintenance mode")
	}
}
