package bff

import (
	"context"
	"sync"
	"time"

	members "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	"google.golang.org/protobuf/proto"
)

// GDSRPC allows the BFF to issue arbitrary GDS client methods in parallel to both the
// testnet and the mainnet. The GDS client and network are passed into the function,
// allowing the RPC to make any GDS RPC call and log with the associated network.
type GDSRPC func(ctx context.Context, client gds.TRISADirectoryClient, network string) (proto.Message, error)

// ParallelGDSRequest makes concurrent requests to both the testnet and the mainnet,
// storing the results and errors in a slice of length 2 ([testnet, mainnet]). If the
// flatten bool is true, then nil values are removed from the slice (though this will
// make which network returned the result ambiguous).
func (s *Server) ParallelGDSRequests(ctx context.Context, rpc GDSRPC, flatten bool) (results []proto.Message, errs []error) {
	// Create the results and errors slices
	results = make([]proto.Message, 2)
	errs = make([]error, 2)

	// Execute the request in parallel to both the testnet and the mainnet
	ctx, cancel := context.WithTimeout(ctx, 25*time.Second)
	var wg sync.WaitGroup
	defer cancel()
	wg.Add(2)

	// Create a closure to execute the rpc
	closure := func(client gds.TRISADirectoryClient, idx int, network string) {
		defer wg.Done()
		results[idx], errs[idx] = rpc(ctx, client, network)
	}

	// execute both requests
	go closure(s.testnet.gds, 0, testnet)
	go closure(s.mainnet.gds, 1, mainnet)
	wg.Wait()

	// flatten rpc and error if requested
	if flatten {
		return FlattenResults(results), FlattenErrs(errs)
	}
	return results, errs
}

// MembersRPC allows the BFF to issue arbitrary members client methods in parallel to both the
// testnet and the mainnet. The members client and network are passed into the function,
// allowing the RPC to make any GDS RPC call and log with the associated network.
type MembersRPC func(ctx context.Context, client members.TRISAMembersClient, network string) (proto.Message, error)

// ParallelMembersRequest makes concurrent requests to both the testnet and the mainnet,
// storing the results and errors in a slice of length 2 ([testnet, mainnet]). If the
// flatten bool is true, then nil values are removed from the slice (though this will
// make which network returned the result ambiguous).
func (s *Server) ParallelMembersRequests(ctx context.Context, rpc MembersRPC, flatten bool) (results []proto.Message, errs []error) {
	// Create the results and errors slices
	results = make([]proto.Message, 2)
	errs = make([]error, 2)

	// Execute the request in parallel to both the testnet and the mainnet
	ctx, cancel := context.WithTimeout(ctx, 25*time.Second)
	var wg sync.WaitGroup
	defer cancel()
	wg.Add(2)

	// Create the closures to execute the rpcs
	closure := func(client members.TRISAMembersClient, idx int, network string) {
		defer wg.Done()
		results[idx], errs[idx] = rpc(ctx, client, network)
	}

	// execute both requests
	go closure(s.testnet.members, 0, testnet)
	go closure(s.mainnet.members, 1, mainnet)
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
