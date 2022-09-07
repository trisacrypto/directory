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
	"github.com/trisacrypto/directory/pkg/bff/config"
	records "github.com/trisacrypto/directory/pkg/bff/db/models/v1"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var (
	trisatest = map[string]struct{}{
		"trisatest.net": {},
		"trisatest.dev": {},
	}
	vaspdirectory = map[string]struct{}{
		"vaspdirectory.net": {},
		"vaspdirectory.dev": {},
	}
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

		if _, ok := result.(*gds.LookupReply); !ok {
			err := fmt.Errorf("unexpected result type: %T", result)
			log.Error().Err(err).Msg("unexpected result type")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
			return
		}

		data, err := wire.Rewire(result.(proto.Message))
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
	if !validRegisteredDirectory(params.Directory) {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("unknown registered directory"))
		return
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

	switch registeredDirectoryType(params.Directory) {
	case config.TestNet:
		rep, err = s.testnetGDS.VerifyContact(ctx, req)
	case config.MainNet:
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

// Returns the user's current registration form if it's available
func (s *Server) LoadRegisterForm(c *gin.Context) {
	// Load the organization from the claims
	// NOTE: this method will handle the error logging and response.
	org, err := s.OrganizationFromClaims(c)
	if err != nil {
		return
	}

	// Return the registration form, ensuring nil is not serialized.
	if org.Registration == nil {
		org.Registration = records.NewRegisterForm()
	}
	c.JSON(http.StatusOK, org.Registration)
}

// Saves the registration form on the BFF to allow multiple users to edit the
// registration form before it is submitted to the directory service.
func (s *Server) SaveRegisterForm(c *gin.Context) {
	// Parse the incoming JSON data from the client request
	var (
		err  error
		form *records.RegistrationForm
		org  *records.Organization
	)

	// Unmarshal the registration form from the POST request
	form = &records.RegistrationForm{}
	if err = c.ShouldBind(form); err != nil {
		log.Warn().Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Load the organization from the claims
	// NOTE: this method will handle the error logging and response.
	if org, err = s.OrganizationFromClaims(c); err != nil {
		return
	}

	// Mark the form as started, the BFF relies on this state so the frontend should
	// capture the updated form returned from this endpoint to avoid overwriting the
	// state metadata.
	// NOTE: If an empty form was passed in, the form will not be marked as started.
	if form.State != nil && form.State.Started == "" {
		form.State.Started = time.Now().Format(time.RFC3339)
	}

	// Update the organizations form
	org.Registration = form
	if err = s.db.Organizations().Update(c.Request.Context(), org); err != nil {
		log.Error().Err(err).Msg("could not update organization")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not save registration form"))
		return
	}

	if org.Registration.State == nil || org.Registration.State.Started == "" {
		// If an empty form was passed in, return a 204 No Content response
		c.Status(http.StatusNoContent)
	} else {
		// Otherwise, return the form in a 200 OK response
		c.JSON(http.StatusOK, org.Registration)
	}
}

// SubmitRegistration makes a request on behalf of the user to either the TestNet or the
// MainNet GDS server based on the URL endpoint. The endpoint will first load the saved
// registration form from the front-end and will parse it for some basic validity
// constraints - it will then submit the form and return any response from the directory.
func (s *Server) SubmitRegistration(c *gin.Context) {
	// Get the network from the URL
	var err error
	network := strings.ToLower(c.Param("network"))
	if network != config.TestNet && network != config.MainNet {
		c.JSON(http.StatusNotFound, api.ErrorResponse("network should be either testnet or mainnet"))
		return
	}

	// Load the organization from the claims
	// NOTE: this method will handle the error logging and response.
	var org *records.Organization
	if org, err = s.OrganizationFromClaims(c); err != nil {
		return
	}

	// Do not allow a registration form to be submitted twice
	switch network {
	case config.TestNet:
		if org.Testnet != nil && org.Testnet.Submitted != "" {
			err = fmt.Errorf("registration form has already been submitted to the %s", network)
			log.Warn().Err(err).Str("network", network).Str("orgID", org.Id).Msg("cannot resubmit registration")
			c.JSON(http.StatusConflict, api.ErrorResponse(err))
			return
		}
	case config.MainNet:
		if org.Mainnet != nil && org.Mainnet.Submitted != "" {
			err = fmt.Errorf("registration form has already been submitted to the %s", network)
			log.Warn().Err(err).Str("network", network).Str("orgID", org.Id).Msg("cannot resubmit registration")
			c.JSON(http.StatusConflict, api.ErrorResponse(err))
			return
		}
	}

	// Validate that a registration form exists on the organization
	if org.Registration == nil || !org.Registration.ReadyToSubmit(network) {
		log.Debug().Str("orgID", org.Id).Msg("cannot submit empty or partial registration form")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("registration form is not ready to submit"))
		return
	}

	// Create the RegisterRequest to send to GDS
	req := &gds.RegisterRequest{
		Entity:         org.Registration.Entity,
		Contacts:       org.Registration.Contacts,
		Website:        org.Registration.Website,
		VaspCategories: org.Registration.VaspCategories,
		EstablishedOn:  org.Registration.EstablishedOn,
		Trixo:          org.Registration.Trixo,
	}

	// Make the GDS request
	var rep *gds.RegisterReply
	log.Debug().Str("network", network).Msg("issuing GDS register request")
	ctx, cancel := context.WithTimeout(c.Request.Context(), 25*time.Second)
	defer cancel()

	switch network {
	case config.TestNet:
		req.TrisaEndpoint = org.Registration.Testnet.Endpoint
		req.CommonName = org.Registration.Testnet.CommonName
		rep, err = s.testnetGDS.Register(ctx, req)
	case config.MainNet:
		req.TrisaEndpoint = org.Registration.Mainnet.Endpoint
		req.CommonName = org.Registration.Mainnet.CommonName
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
			log.Error().Err(err).Str("network", network).Msg("could not rewire response error struct")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(fmt.Errorf("could not handle register response from %s", network)))
			return
		}
	}

	// If there is only an error and no vaspID, return a 409 and the response
	if out.Error != nil && out.Id == "" {
		log.Error().Err(rep.Error).Str("network", network).Msg("received unexpected GDS OK with an error in the response")
		c.JSON(http.StatusConflict, out)
		return
	}

	// Save the response with the organization details in the Organization document.
	directoryRecord := &records.DirectoryRecord{
		Id:                  rep.Id,
		RegisteredDirectory: rep.RegisteredDirectory,
		CommonName:          rep.CommonName,
		Submitted:           time.Now().Format(time.RFC3339),
	}

	switch network {
	case config.TestNet:
		org.Testnet = directoryRecord
	case config.MainNet:
		org.Mainnet = directoryRecord
	}

	if err = s.db.Organizations().Update(c.Request.Context(), org); err != nil {
		log.Error().Err(err).Str("network", network).Msg("could not update organization with directory record")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete registration submission"))
		return
	}

	c.JSON(http.StatusOK, out)
}

// Checks if the user supplied registered directory is one of the known directories that
// maps to either the testnet or to the mainnet (by domain).
func validRegisteredDirectory(r string) bool {
	if _, ok := trisatest[r]; ok {
		return true
	}

	if _, ok := vaspdirectory[r]; ok {
		return true
	}

	return false
}

// Returns either testnet or mainnet depending on the user supplied registered directory.
func registeredDirectoryType(r string) string {
	if _, ok := trisatest[r]; ok {
		return config.TestNet
	}

	if _, ok := vaspdirectory[r]; ok {
		return config.MainNet
	}

	return ""
}
