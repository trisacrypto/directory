package gds

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) serverInterceptor(ctx context.Context, in interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (out interface{}, err error) {
	// Track how long the method takes to execute.
	start := time.Now()

	// Check if we're in maintenance mode - status method should still return a full response.
	if s.conf.Maintenance && info.FullMethod != "/trisa.gds.api.v1beta1.TRISADirectory/Status" {
		return nil, status.Error(codes.Unavailable, "the GDS service is currently in maintenance mode")
	}

	// Call the handler to finalize the request and get the response.
	out, err = handler(ctx, in)

	// Log with zerolog - checkout grpclog.LoggerV2 for default logging.
	log.Debug().Str("method", info.FullMethod).Str("latency", time.Since(start).String()).Err(err)
	return out, err
}
