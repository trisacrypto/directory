package bff

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

const (
	testnet = "testnet"
	mainnet = "mainnet"
)

// ConnectGDS creates a gRPC client to the TRISA Directory Service specified in the
// configuration. This method is used to connect to both the TestNet and the MainNet and
// to connect to mock GDS services in testing using buffconn.
func ConnectGDS(conf config.DirectoryConfig) (_ gds.TRISADirectoryClient, err error) {
	// Create the Dial options with required credentials
	var opts []grpc.DialOption
	if conf.Insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	}

	ctx, cancel := context.WithTimeout(context.Background(), conf.Timeout)
	defer cancel()

	// Connect the directory client (non-blocking)
	var cc *grpc.ClientConn
	if cc, err = grpc.DialContext(ctx, conf.Endpoint, opts...); err != nil {
		return nil, err
	}
	return gds.NewTRISADirectoryClient(cc), nil
}

// Lookup makes a request on behalf of the user to both the TestNet and MainNet GDS
// servers, returning 1-2 results (e.g. either or both GDS responses). If no results
// are returned, Lookup returns a 404 not found error. If one of the GDS requests fails,
// the error is logged, but the valid response is returned. If both GDS requests fail,
// a 500 error is returned. This endpoint passes through the response from GDS as JSON,
// the result should contain a registered_directory field that identifies which network
// the record is associated with.
func (s *Server) Lookup(c *gin.Context) {
	// Bind the parameters associated with the lookup
	params := &api.LookupParams{}
	if err := c.ShouldBindQuery(&params); err != nil {
		log.Warn().Err(err).Msg("could not bind request with query params")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Ensure that we have either ID or CommonName
	if params.ID == "" && params.CommonName == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("must provide either uuid or common_name in query params"))
		return
	}

	// Create the LookupRequest, omit registered directory as it is assumed we're
	// looking up the value for the directory we're making the request to.
	req := &gds.LookupRequest{Id: params.ID, CommonName: params.CommonName}
	errs := make([]error, 2)
	results := make([]map[string]interface{}, 2)

	// Execute the request in parallel to both the testnet and the mainnet
	ctx, cancel := context.WithTimeout(c.Request.Context(), 25*time.Second)
	var wg sync.WaitGroup
	defer cancel()
	wg.Add(2)

	lookup := func(client gds.TRISADirectoryClient, idx int, name string) {
		defer wg.Done()
		var (
			err error
			rep *gds.LookupReply
		)

		if rep, err = client.Lookup(ctx, req); err != nil {
			serr, _ := status.FromError(err)
			if serr.Code() != codes.NotFound {
				log.Error().Err(err).Str("network", name).Msg("GDS lookup unsuccessful")
				errs[idx] = err
			}
			return
		}

		// There is currently an error message on the reply that is unused. This check
		// future proofs the use case where that field is returned and also ensures that
		// an empty reply is not returned.
		if rep.Error != nil && rep.Error.Code != 0 {
			rerr := fmt.Errorf("[%d] %s", rep.Error.Code, rep.Error.Message)
			log.Warn().Err(rerr).Msg("received error in response body with a gRPC status ok")
			if rep.Id == "" && rep.CommonName == "" {
				// If we don't have an ID or common name, don't return an empty result
				errs[idx] = rerr
				return
			}
		}

		if results[idx], err = wire.Rewire(rep); err != nil {
			log.Error().Err(err).Str("network", name).Msg("could not rewire LookupReply")
			errs[idx] = err
		}
	}

	go lookup(s.testnet, 0, testnet)
	go lookup(s.mainnet, 1, mainnet)
	wg.Wait()

	// If there were multiple errors, return a 500
	nErrs := 0
	for _, err := range errs {
		if err != nil {
			nErrs++
		}
	}
	if nErrs == 2 {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("unable to execute Lookup request"))
		return
	}

	// Flatten results
	out := &api.LookupReply{
		Results: make([]map[string]interface{}, 0, 2),
	}
	for _, result := range results {
		if len(result) > 0 {
			out.Results = append(out.Results, result)
		}
	}

	// Check if there are results to return
	if len(out.Results) == 0 {
		c.JSON(http.StatusNotFound, api.ErrorResponse("no results returned for query"))
		return
	}

	// Serialize the results and return a successful response
	c.JSON(http.StatusOK, out)
}
