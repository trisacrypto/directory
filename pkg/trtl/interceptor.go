package trtl

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"net"
	"time"

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

	// Check if we're in maintenance mode
	// TODO: update to consider the Status endpoint, once it has been added
	if t.conf.Maintenance {
		return nil, status.Error(codes.Unavailable, "the trtl service is currently in maintenance mode")
	}

	// Fetch peer information from the TLS info.
	var peer *PeerInfo
	if peer, err = peerFromTLS(ctx); err != nil {
		log.Error().Err(err).Msg("unable to retrieve remote peer info")
		// TODO: Uncomment when mTLS is added
		//return nil, status.Error(codes.Unauthenticated, "unable to retrieve authenticated peer information")
	}

	// Add peer information to the context.
	ctx = context.WithValue(ctx, ContextKey("peer"), peer)

	// Call the handler to finalize the request and get the response.
	out, err = handler(ctx, in)

	// Log with zerolog - checkout grpclog.LoggerV2 for default logging.
	log.Debug().Str("method", info.FullMethod).Str("latency", time.Since(start).String()).Err(err).Msg("gRPC request complete")
	return out, err
}

// peerFromTLS looks up the TLSInfo from the incoming gRPC connection to retrieve
// information about the remote peer from the certificate.
func peerFromTLS(ctx context.Context) (info *PeerInfo, err error) {
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
				IPAddresses: []net.IP{net.ParseIP(gp.Addr.String())},
			}
			info.IPAddresses = append(info.IPAddresses, c[0].IPAddresses...)
			return info, nil
		}
	}

	return nil, errors.New("could not find peer certificate info")
}
