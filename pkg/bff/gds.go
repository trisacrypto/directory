package bff

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	models "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

const (
	testnet       = "testnet"
	mainnet       = "mainnet"
	trisatest     = "trisatest.net"
	vaspdirectory = "vaspdirectory.net"
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

// Register makes a request on behalf of the user to either the TestNet or the MainNet
// GDS server based on the URL endpoint. This method is essentially a passthrough
// conversion of a JSON request to a gRPC request to GDS.
func (s *Server) Register(c *gin.Context) {
	// Get the network from the URL
	network := strings.ToLower(c.Param("network"))

	// Parse the incoming JSON data from the client request
	var err error
	in := &api.RegisterRequest{}
	if err = c.ShouldBind(&in); err != nil {
		log.Warn().Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Sanity check: validate the network with the request if supplied
	if in.Network != "" && in.Network != network {
		log.Warn().Str("data", in.Network).Str("param", network).Msg("mismatched request network and URL")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("the request network does not match the URL endpoint"))
		return
	}

	// Prevent panics by requiring data that will be unwired
	switch {
	case in.BusinessCategory == "":
		c.JSON(http.StatusBadRequest, api.ErrorResponse("business category is required"))
		return
	case len(in.Entity) == 0:
		c.JSON(http.StatusBadRequest, api.ErrorResponse("entity is required"))
		return
	case len(in.Contacts) == 0:
		c.JSON(http.StatusBadRequest, api.ErrorResponse("contacts are required"))
		return
	case len(in.TRIXO) == 0:
		c.JSON(http.StatusBadRequest, api.ErrorResponse("trixo is required"))
		return
	}

	// Create the RegisterRequest to send to GDS
	req := &gds.RegisterRequest{
		Entity:           &ivms101.LegalPerson{},
		Contacts:         &models.Contacts{},
		TrisaEndpoint:    in.TRISAEndpoint,
		CommonName:       in.CommonName,
		Website:          in.Website,
		BusinessCategory: models.BusinessCategoryUnknown,
		VaspCategories:   in.VASPCategories,
		EstablishedOn:    in.EstablishedOn,
		Trixo:            &models.TRIXOQuestionnaire{},
	}

	// Unwire the protocol buffers into the request
	if err = wire.Unwire(in.Entity, req.Entity); err != nil {
		log.Warn().Err(err).Msg("could not unwire legal person entity")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse legal person entity"))
		return
	}

	if err = wire.Unwire(in.Contacts, req.Contacts); err != nil {
		log.Warn().Err(err).Msg("could not unwire contacts")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse contacts"))
		return
	}

	if req.BusinessCategory, err = models.ParseBusinessCategory(in.BusinessCategory); err != nil {
		log.Warn().Err(err).Str("input", in.BusinessCategory).Msg("could not parse business category")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	if err = wire.Unwire(in.TRIXO, req.Trixo); err != nil {
		log.Warn().Err(err).Msg("could not unwire TRIXO form")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse TRIXO form"))
		return
	}

	// Make the GDS request
	log.Debug().Str("network", network).Msg("issuing GDS register request")
	ctx, cancel := context.WithTimeout(c.Request.Context(), 25*time.Second)
	var rep *gds.RegisterReply
	defer cancel()

	switch network {
	case testnet:
		rep, err = s.testnet.Register(ctx, req)
	case mainnet:
		rep, err = s.mainnet.Register(ctx, req)
	default:
		c.JSON(http.StatusNotFound, api.ErrorResponse("network should be either testnet or mainnet"))
		return
	}

	// Handle GDS errors
	if err != nil {
		serr, _ := status.FromError(err)
		switch serr.Code() {
		case codes.InvalidArgument, codes.AlreadyExists:
			c.JSON(http.StatusBadRequest, api.ErrorResponse(serr.Message()))
		case codes.Aborted:
			c.JSON(http.StatusConflict, api.ErrorResponse(serr.Message()))
		default:
			log.Error().Err(err).Str("code", serr.Code().String()).Str("network", network).Msg("could not register with directory service")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(fmt.Errorf("could not register with %s", network)))
		}
		return
	}

	// Create the response from the reply
	out := &api.RegisterReply{
		Id:                  rep.Id,
		RegisteredDirectory: rep.RegisteredDirectory,
		CommonName:          rep.CommonName,
		Status:              rep.Status.String(),
		Message:             rep.Message,
		PKCS12Password:      rep.Pkcs12Password,
	}

	if rep.Error != nil && rep.Error.Code != 0 {
		if out.Error, err = wire.Rewire(rep.Error); err != nil {
			log.Error().Err(err).Msg("could not rewire response error struct")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(fmt.Errorf("could not handle register response from %s", network)))
			return
		}
	}

	c.JSON(http.StatusOK, out)
}

// VerifyContact is currently a passthrough helper that forwards the verify contact
// request from the user interface to the GDS that needs contact verification.
func (s *Server) VerifyContact(c *gin.Context) {
	// Bind the parameters associated with the verify contact request
	params := &api.VerifyContactParams{}
	if err := c.ShouldBindQuery(&params); err != nil {
		log.Warn().Err(err).Msg("could not bind request with query params")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Ensure that we have all required parameters
	if params.ID == "" || params.Token == "" || params.Directory == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("must provide vaspID, token, and registered_directory in query parameters"))
		return
	}

	// Ensure the registered_directory is one we understand
	params.Directory = strings.ToLower(params.Directory)
	if params.Directory != trisatest && params.Directory != vaspdirectory {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("unknown registered directory"))
	}

	// Make the GDS request
	log.Debug().Str("registered_directory", params.Directory).Msg("issuing GDS verify contact request")
	req := &gds.VerifyContactRequest{Id: params.ID, Token: params.Token}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 25*time.Second)
	defer cancel()

	var (
		err error
		rep *gds.VerifyContactReply
	)

	switch params.Directory {
	case trisatest:
		rep, err = s.testnet.VerifyContact(ctx, req)
	case vaspdirectory:
		rep, err = s.mainnet.VerifyContact(ctx, req)
	default:
		log.Error().Str("registered_directory", params.Directory).Str("endpoint", "verify").Msg("unhandled directory")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not verify contact"))
		return
	}

	// Handle GDS errors
	if err != nil {
		serr, _ := status.FromError(err)
		switch serr.Code() {
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, api.ErrorResponse(serr.Message()))
		case codes.NotFound:
			c.JSON(http.StatusNotFound, api.ErrorResponse(serr.Message()))
		case codes.Aborted:
			c.JSON(http.StatusConflict, api.ErrorResponse(serr.Message()))
		default:
			log.Error().Err(err).Str("code", serr.Code().String()).Str("registered_directory", params.Directory).Msg("could not verify contact")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(serr.Message()))
		}
		return
	}

	// Create the response from the reply
	out := &api.VerifyContactReply{
		Status:  rep.Status.String(),
		Message: rep.Message,
	}

	if rep.Error != nil && rep.Error.Code != 0 {
		if out.Error, err = wire.Rewire(rep.Error); err != nil {
			log.Error().Err(err).Msg("could not rewire response error struct")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(fmt.Errorf("could not handle verify contact response from %s", params.Directory)))
			return
		}
	}

	c.JSON(http.StatusOK, out)
}
