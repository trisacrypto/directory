package interceptors

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// UnaryMTLS adds authenticated peer info to the context and returns an unauthenticated
// error if that peer information is not available.
func UnaryMTLS() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		var peer *PeerInfo
		if peer, err = PeerFromTLS(ctx); err != nil {
			err = status.Error(codes.Unauthenticated, "unable to retrieve authenticated peer information")
			log.Debug().Err(err).Str("method", info.FullMethod).Msg("unauthenticated access detected")
			return nil, err
		}

		// Add peer information to the context
		ctx = context.WithValue(ctx, ContextKey("peer"), peer)
		return handler(ctx, req)
	}
}

// StreamMTLS adds authenticated peer info to the context and returns an unauthenticated
// error if that peer information is not available.
func StreamMTLS() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		// TODO: add peer information to the context, this requires wrapping the ServerStream
		// See: https://stackoverflow.com/questions/60982406/how-to-safely-add-values-to-grpc-serverstream-in-interceptor
		if _, err = PeerFromTLS(stream.Context()); err != nil {
			err = status.Error(codes.Unauthenticated, "unable to retrieve authenticated peer information")
			log.Debug().Err(err).Str("method", info.FullMethod).Msg("unauthenticated access detected")
			return err
		}
		return handler(srv, stream)
	}
}

type ContextKey string

// PeerInfo stores information about the identity of a remote peer.
type PeerInfo struct {
	Name        *pkix.Name
	DNSNames    []string
	IPAddresses []net.IP
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
