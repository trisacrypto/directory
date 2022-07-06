package bff

import (
	"context"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

// RPC allows the BFF to issue arbitrary client methods in parallel to both the
// testnet and the mainnet. The combined client object, which contains separate
// sub-clients for the GDS and members services, and network name are passed into the
// function, allowing the RPC to make any directory service or members service RPC
// call and log with the associated network.
type RPC func(ctx context.Context, client GlobalDirectoryClient, network string) (proto.Message, error)

// ParallelGDSRequests makes concurrent requests to both the testnet and the mainnet,
// storing the results and errors in a slice of length 2 ([testnet, mainnet]). If the
// flatten bool is true, then nil values are removed from the slice (though this will
// make which network returned the result ambiguous).
func (s *Server) ParallelGDSRequests(ctx context.Context, rpc RPC, flatten bool) (results []proto.Message, errs []error) {
	// Create the results and errors slices
	results = make([]proto.Message, 2)
	errs = make([]error, 2)

	// Execute the request in parallel to both the testnet and the mainnet
	ctx, cancel := context.WithTimeout(ctx, 25*time.Second)
	var wg sync.WaitGroup
	defer cancel()
	wg.Add(2)

	// Create a closure to execute the rpc
	closure := func(client GlobalDirectoryClient, idx int, network string) {
		defer wg.Done()
		results[idx], errs[idx] = rpc(ctx, client, network)
	}

	// execute both requests
	go closure(s.testnetGDS, 0, testnet)
	go closure(s.mainnetGDS, 1, mainnet)
	wg.Wait()

	// flatten rpc and error if requested
	if flatten {
		return FlattenResults(results), FlattenErrs(errs)
	}
	return results, errs
}

// FlattenResults removes nil proto messages from the slice (exported for testing purposes).
func FlattenResults(in []proto.Message) (out []proto.Message) {
	out = make([]proto.Message, 0, len(in))
	for _, msg := range in {
		if msg != nil {
			out = append(out, msg)
		}
	}
	return out
}

// FlattenErrs removes nil errors from the slice (exported for testing purposes).
func FlattenErrs(in []error) (out []error) {
	out = make([]error, 0, len(in))
	for _, err := range in {
		if err != nil {
			out = append(out, err)
		}
	}
	return out
}
