package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	wr "github.com/mroth/weightedrand"
	"github.com/trisacrypto/directory/pkg/trtl/jitter"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/trisacrypto/trisa/pkg/trust"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

const (
	interval  = 5 * time.Second                                        // the ticker interval, default to 5 seconds
	sigma     = 100 * time.Millisecond                                 // the amount of jitter, default to 100 ms
	accesses  = 15                                                     // desired accesses per interval, default to 15
	regions   = 7                                                      // number of regions simultaneously using the accessor
	endpoint  = "localhost:4436"                                       // the endpoint of the running trtl server
	insecure  = true                                                   // connect without mTLS
	certPath  = "fixtures/certs/mtls.client.dev/client.pem"            // path to file on disk with certificate key-pair if using mTLS
	poolPath  = "fixtures/certs/mtls.client.dev/certificate.chain.pem" // path to file on disk with trust pool if using mTLS
	keyspace  = 1000                                                   // the number of keys the simulator operates on
	chunkSize = 512                                                    // size of each write
)

// dummy namespaces for simulation
var namespaces = []string{"catchfireBarons", "falselightCutters", "fullCrowns", ""}

// probabilities for reads/writes/deletes
var probabilities = map[string]uint{
	"read":   60,
	"write":  38,
	"delete": 2,
}

func main() {

	// assumes trtl is already being served (e.g. from the trtl cli)
	// sim needs the endpoint (e.g. localhost:port) + certs (just stubs for now)
	// run trtl in insecure mode
	sim := new(endpoint, insecure)
	simClient, err := sim.connect()
	if err != nil {
		panic(err)
	}

	ticker := jitter.New(interval, sigma)

	// start on the first tick
	for ; true; <-ticker.C {

		// multiple regions access the data store concurrently; 1 routine per region
		for r := 1; r <= regions; r++ {
			go sim.accessor(simClient)
		}
	}
}

type Simulator struct {
	Endpoint string      `yaml:"endpoint"`           // the replica endpoint to connect to
	Insecure bool        `yaml:"insecure,omitempty"` // do not connect with TLS
	Selector *wr.Chooser `yaml:"chooser,omitempty"`  // random selection helper
	CertPath string      `yaml:cert_path,omitempty"` // path to certificate key pair for client side mTLS
	PoolPath string      `yaml:poo_path,omitempty"`  // path to certificate trust chain for client side mTLS
}

func new(endpoint string, insecure bool) *Simulator {
	// initialize weighted probability selector
	selector := initialize()
	return &Simulator{
		Endpoint: endpoint,
		Insecure: insecure,
		Selector: selector,
		CertPath: certPath,
		PoolPath: poolPath,
	}
}

// Connect to the trtl server and return a gRPC client
func (s *Simulator) connect() (_ pb.TrtlClient, err error) {
	var opts []grpc.DialOption
	if s.Insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		var sz *trust.Serializer
		if sz, err = trust.NewSerializer(false); err != nil {
			return nil, err
		}

		var pool trust.ProviderPool
		if pool, err = sz.ReadPoolFile(s.PoolPath); err != nil {
			return nil, err
		}

		var provider *trust.Provider
		if provider, err = sz.ReadFile(s.CertPath); err != nil {
			return nil, err
		}

		var cert tls.Certificate
		if cert, err = provider.GetKeyPair(); err != nil {
			return nil, err
		}

		var certPool *x509.CertPool
		if certPool, err = pool.GetCertPool(false); err != nil {
			return nil, err
		}
		var u *url.URL
		if u, err = url.Parse(s.Endpoint); err != nil {
			return nil, err
		}
		conf := &tls.Config{
			ServerName:   u.Host,
			Certificates: []tls.Certificate{cert},
			RootCAs:      certPool,
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(conf)))
	}

	// Connect the replica client
	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(s.Endpoint, opts...); err != nil {
		return nil, err
	}
	return pb.NewTrtlClient(cc), nil
}

// Create a fixed number of random accesses across the namespace and keyspace options
func (s *Simulator) accessor(client pb.TrtlClient) {
	// Make as many accesses as directed by the global variable/config
	for a := 1; a <= accesses; a++ {

		// randomly select a namespace from namespaces
		ns := namespaces[rand.Intn(len(namespaces))]

		// randomly select a key from keyspace and cast to bytes
		key := []byte(strconv.Itoa(rand.Intn(keyspace)))

		// select from read, write, delete using probabilities
		switch s.Selector.Pick().(string) {
		case "read":
			// execute Get
			req := &pb.GetRequest{
				Key:       key,
				Namespace: ns,
				Options: &pb.Options{
					ReturnMeta: false,
				},
			}
			if _, err := client.Get(context.TODO(), req); err != nil {
				// Skip if the object doesn't exist
				if serr, ok := status.FromError(err); ok && serr.Code() == codes.NotFound {
					continue
				}
				panic(fmt.Errorf("could not read from database: %v", err))
			}

		case "write":
			// create random data
			val := make([]byte, chunkSize)

			// execute Put
			req := &pb.PutRequest{
				Key:       key,
				Value:     val,
				Namespace: ns,
				Options: &pb.Options{
					ReturnMeta: false,
				},
			}
			if _, err := client.Put(context.TODO(), req); err != nil {
				fmt.Print(err)
				panic("could not write to database")
			}
		case "delete":
			// First check that the object exists
			req := &pb.GetRequest{
				Key:       key,
				Namespace: ns,
				Options: &pb.Options{
					ReturnMeta: false,
				},
			}
			// Object doesn't exist so skip this one
			if _, err := client.Get(context.TODO(), req); err != nil {
				if serr, ok := status.FromError(err); ok && serr.Code() == codes.NotFound {
					continue
				}
			} else {
				// Object does exist so execute Delete
				req := &pb.DeleteRequest{
					Key:       key,
					Namespace: ns,
					Options: &pb.Options{
						ReturnMeta: false,
					},
				}
				if _, err := client.Delete(context.TODO(), req); err != nil {
					if serr, ok := status.FromError(err); ok && serr.Code() == codes.NotFound {
						continue
					}
					panic(fmt.Errorf("could not delete from database: %v", err))
				}
			}
		default:
			panic(errors.New("unknown database operation"))
		}
	}
}

//===========================================================================
// Helper for weighted probability selection
//===========================================================================

func initialize() *wr.Chooser {
	rand.Seed(time.Now().UTC().UnixNano())
	choices := make([]wr.Choice, len(namespaces))
	for i, w := range probabilities {
		choices = append(choices, wr.Choice{Item: i, Weight: w})
	}
	chooser, err := wr.NewChooser(choices...)
	if err != nil {
		panic("error in chooser creation")
	}
	return chooser
}
