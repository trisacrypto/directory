package trtl

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"net"
	"runtime/debug"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	status "google.golang.org/grpc/status"
)

type ContextKey string

// PeerInfo stores information about the identity of a remote peer.
type PeerInfo struct {
	Name        *pkix.Name
	DNSNames    []string
	IPAddresses []net.IP
}

// The interceptor intercepts incoming gRPC requests and adds remote peer information to the context.
func (t *Server) interceptor(ctx context.Context, in interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (out interface{}, err error) {
	// Track how long the method takes to execute.
	start := time.Now()
	panicked := true

	// Set the service tag
	if t.conf.Sentry.UseSentry() {
		sentry.CurrentHub().Scope().SetTag("service", "trtl")
	}

	// Recover from panics in the handler.
	defer func() {
		if r := recover(); r != nil || panicked {
			if t.conf.Sentry.UseSentry() {
				sentry.CurrentHub().Recover(r)
			}
			log.WithLevel(zerolog.PanicLevel).
				Err(fmt.Errorf("%v", r)).
				Str("stack_trace", string(debug.Stack())).
				Msg("trtl server has recovered from a panic")
			err = status.Error(codes.Internal, "an unhandled exception occurred")
		}
	}()

	// Check if we're in maintenance mode - status method should still return a full response
	if t.conf.Maintenance && info.FullMethod != "/trtl.v1.Trtl/Status" {
		err = status.Error(codes.Unavailable, "the trtl service is currently in maintenance mode")
		log.Trace().Err(err).Str("method", info.FullMethod).Msg("trtl service unavailable during maintenance mode")

		panicked = false
		return nil, err
	}

	// Fetch peer information from the TLS info if we're not in insecure mode.
	if !t.conf.MTLS.Insecure {
		var peer *PeerInfo
		if peer, err = PeerFromTLS(ctx); err != nil {
			err = status.Error(codes.Unauthenticated, "unable to retrieve authenticated peer information")
			log.Warn().Err(err).Str("method", info.FullMethod).Msg("unauthenticated access detected")

			panicked = false
			return nil, err
		}

		// Add peer information to the context.
		ctx = context.WithValue(ctx, ContextKey("peer"), peer)
	}

	// Call the handler to finalize the request and get the response.
	var span *sentry.Span
	if t.conf.Sentry.UsePerformanceTracking() {
		span = sentry.StartSpan(ctx, "grpc handler", sentry.TransactionName(info.FullMethod))
	}
	out, err = handler(ctx, in)
	if t.conf.Sentry.UsePerformanceTracking() {
		span.Finish()
	}

	// Log with zerolog - checkout grpclog.LoggerV2 for default logging.
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
func (t *Server) streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	// Track how long the method takes to execute
	start := time.Now()
	panicked := true

	// Set the service tag for sentry
	if t.conf.Sentry.UseSentry() {
		sentry.CurrentHub().Scope().SetTag("service", "trtl")
	}

	// Recover from panics in the handler
	defer func() {
		if r := recover(); r != nil || panicked {
			if t.conf.Sentry.UseSentry() {
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
	if t.conf.Maintenance {
		err = status.Error(codes.Unavailable, "the trtl service is currently in maintenance mode")
		log.Trace().Err(err).Str("method", info.FullMethod).Msg("trtl service unavailable during maintenance mode")

		panicked = false
		return err
	}

	// Fetch peer information from the TLS info if we're not in insecure mode.
	if !t.conf.MTLS.Insecure {
		if _, err = PeerFromTLS(ss.Context()); err != nil {
			err = status.Error(codes.Unauthenticated, "unable to retrieve authenticated peer information")
			log.Warn().Err(err).Str("method", info.FullMethod).Msg("unauthenticated access detected")

			panicked = false
			return err
		}

		// TODO: add peer information to the context, this requires wrapping the ServerStream
		// See: https://stackoverflow.com/questions/60982406/how-to-safely-add-values-to-grpc-serverstream-in-interceptor
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

// PeerFromTLS looks up the TLSInfo from the incoming gRPC connection to retrieve
// information about the remote peer from the certificate.
func PeerFromTLS(ctx context.Context) (info *PeerInfo, err error) {
	var (
		ok      bool
		gp      *peer.Peer
		tlsAuth credentials.TLSInfo
		chains  [][]*x509.Certificate
	)

	if gp, ok = peer.FromContext(ctx); !ok {
		return nil, errors.New("no peer found in context")
	}

	if tlsAuth, ok = gp.AuthInfo.(credentials.TLSInfo); !ok {
		// If there is no mTLS information return nil peer info.
		if gp.AuthInfo == nil {
			return nil, nil
		}
		return nil, fmt.Errorf("unexpected peer transport credentials type: %T", gp.AuthInfo)
	}

	chains = tlsAuth.State.VerifiedChains
	if len(chains) == 0 {
		return nil, errors.New("could not find certificate chain")
	}

	// Search for a valid peer certificate.
	for _, c := range chains {
		// Certificate chain should contain at least the peer certificate and the CA.
		if len(c) >= 2 {
			// The peer certificate is always first.
			info = &PeerInfo{
				Name:        &c[0].Subject,
				DNSNames:    c[0].DNSNames,
				IPAddresses: []net.IP{},
			}
			var addr net.IP
			if addr = net.ParseIP(gp.Addr.String()); addr != nil {
				// Only add the net.Addr if it's parseable
				info.IPAddresses = append(info.IPAddresses, addr)
			}

			info.IPAddresses = append(info.IPAddresses, c[0].IPAddresses...)
			return info, nil
		}
	}

	return nil, errors.New("could not find peer certificate info")
}
