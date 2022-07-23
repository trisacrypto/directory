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
				Msg("grpc server has recovered from a panic")
			err = status.Error(codes.Internal, "an unhandled exception occurred")
		}
	}()

	// Check if we're in maintenance mode
	if t.conf.Maintenance {
		// The only RPC we allow in maintenance mode is Status
		if info.FullMethod == "/trtl.v1.Trtl/Status" {
			return handler(ctx, in)
		}

		// Otherwise we stop processing here and return unavailable
		return nil, status.Error(codes.Unavailable, "the trtl service is currently in maintenance mode")
	}

	// Fetch peer information from the TLS info if we're not in insecure mode.
	if !t.conf.MTLS.Insecure {
		var peer *PeerInfo
		if peer, err = PeerFromTLS(ctx); err != nil {
			return nil, status.Error(codes.Unauthenticated, "unable to retrieve authenticated peer information")
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
	panicked = false

	// Log with zerolog - checkout grpclog.LoggerV2 for default logging.
	log.Debug().
		Err(err).
		Str("method", info.FullMethod).
		Str("latency", time.Since(start).String()).
		Msg("gRPC request complete")
	return out, err
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
