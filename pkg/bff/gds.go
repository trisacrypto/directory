package bff

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	models "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const (
	testnet       = "testnet"
	mainnet       = "mainnet"
	trisatest     = "trisatest.net"
	vaspdirectory = "vaspdirectory.net"
)

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

	// Create an RPC func for making a parallel GDS request
	lookup := func(ctx context.Context, client GlobalDirectoryClient, network string) (_ proto.Message, err error) {
		var rep *gds.LookupReply
		if rep, err = client.Lookup(ctx, req); err != nil {
			// If the code is not found then do not return an error, just no result.
			serr, _ := status.FromError(err)
			if serr.Code() != codes.NotFound {
				log.Error().Err(err).Str("network", network).Msg("GDS lookup unsuccessful")
				return nil, err
			}
			return nil, nil
		}

		// There is currently an error message on the reply that is unused. This check
		// future proofs the use case where that field is returned and also ensures that
		// an empty reply is not returned.
		if rep.Error != nil && rep.Error.Code != 0 {
			rerr := fmt.Errorf("[%d] %s", rep.Error.Code, rep.Error.Message)
			log.Warn().Err(rerr).Msg("received error in response body with a gRPC status ok")
			if rep.Id == "" && rep.CommonName == "" {
				// If we don't have an ID or common name, don't return an empty result
				return nil, rerr
			}
		}
		return rep, nil
	}

	// Execute the parallel GDS lookup request, ensuring that flatten is false with the
	// expectation that TestNet will be in the 0 index and MainNet in the 1 index.
	results, errs := s.ParallelGDSRequests(c.Request.Context(), lookup, false)

	// If there were multiple errors, return a 500
	// Because the results cannot be flattened we have to check each err individually.
	if errs[0] != nil && errs[1] != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("unable to execute Lookup request"))
		return
	}

	// Check if there are results to return
	// Because the results cannot be flattened we have to check each result individually.
	if results[0] == nil && results[1] == nil {
		c.JSON(http.StatusNotFound, api.ErrorResponse("no results returned for query"))
		return
	}

	// Rewire the results into a JSON response
	out := &api.LookupReply{}
	for idx, result := range results {
		// Skip over nil results
		if result == nil {
			continue
		}

		data, err := wire.Rewire(result)
		if err != nil {
			log.Error().Err(err).Msg("could not rewire LookupReply")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not process lookup reply"))
			return
		}

		switch idx {
		case 0:
			out.TestNet = data
		case 1:
			out.MainNet = data
		}
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
		rep, err = s.testnetGDS.Register(ctx, req)
	case mainnet:
		rep, err = s.mainnetGDS.Register(ctx, req)
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

	// TODO: save the response with the organization details in the Organization document.
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
		rep, err = s.testnetGDS.VerifyContact(ctx, req)
	case vaspdirectory:
		rep, err = s.mainnetGDS.VerifyContact(ctx, req)
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
