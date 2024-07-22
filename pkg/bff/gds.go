package bff

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/auth0/go-auth0/management"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/config"
	records "github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/store"
	"github.com/trisacrypto/directory/pkg/utils"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
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
//
// @Summary Lookup a VASP record by name or ID
// @Description Lookup a VASP record in both TestNet and MainNet, returning either or both results.
// @Tags GDS
// @Accept json
// @Produce json
// @Param params body api.LookupParams true "Lookup parameters"
// @Success 200 {object} api.LookupReply
// @Failure 400 {object} api.Reply "Either ID or CommonName must be provided"
// @Failure 404 {object} api.Reply "No results returned for query"
// @Failure 500 {object} api.Reply "Internal server error"
// @Router /lookup [get]
func (s *Server) Lookup(c *gin.Context) {
	// Bind the parameters associated with the lookup
	params := &api.LookupParams{}
	if err := c.ShouldBindQuery(&params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not bind request with query params")
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
				sentry.Error(ctx).Err(err).Str("network", network).Msg("GDS lookup unsuccessful")
				return nil, err
			}
			return nil, nil
		}

		// There is currently an error message on the reply that is unused. This check
		// future proofs the use case where that field is returned and also ensures that
		// an empty reply is not returned.
		if rep.Error != nil && rep.Error.Code != 0 {
			rerr := fmt.Errorf("[%d] %s", rep.Error.Code, rep.Error.Message)
			sentry.Warn(ctx).Err(rerr).Msg("received error in response body with a gRPC status ok")
			if rep.Id == "" && rep.CommonName == "" {
				// If we don't have an ID or common name, don't return an empty result
				return nil, rerr
			}
		}
		return rep, nil
	}

	// Execute the parallel GDS lookup request, ensuring that flatten is false with the
	// expectation that TestNet will be in the 0 index and MainNet in the 1 index.
	results, errs := s.ParallelGDSRequests(sentry.RequestContext(c), lookup, false)

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
			sentry.Error(c).Err(err).Msg("unexpected result type")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
			return
		}

		data, err := wire.Rewire(result.(proto.Message))
		if err != nil {
			sentry.Error(c).Err(err).Msg("could not rewire LookupReply")
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

// LookupAutocomplete makes a request on behalf of the user to both the TestNet and
// MainNet GDS databases and returns the complete deduplicated list of verified VASP
// names to the user, to facilitate client-side autocomplete functionality. This is an
// unauthenticated endpoint so it's publicly accessible.
//
// @Summary Get the names of verified VASPs for autocomplete
// @Description Get the names of all the verified VASPs in both TestNet and MainNet.
// @Tags GDS
// @Produce json
// @Success 200 {list} string "List of VASP names"
// @Failure 500 {object} api.Reply
// @Router /lookup/autocommplete [get]
func (s *Server) LookupAutocomplete(c *gin.Context) {
	// Create an RPC func for making a parallel GDS request for the list of VASPs
	rpc := func(ctx context.Context, db store.Store, network string) (rep interface{}, err error) {
		iter := db.ListVASPs(ctx)
		defer iter.Release()

		names := make(map[string]string)
		for iter.Next() {
			var vasp *pb.VASP
			if vasp, err = iter.VASP(); err != nil {
				sentry.Error(c).Err(err).Str("network", network).Msg("could not get VASP from ListVASP iterator")
				continue
			}

			// Filter any VASPs that are not verified
			if vasp.VerificationStatus != pb.VerificationState_VERIFIED {
				continue
			}

			var name string
			if name, err = vasp.Name(); err != nil || name == "" {
				sentry.Warn(c).Err(err).Str("network", network).Str("id", iter.Id()).Msg("could not resolve VASP name from VASP")
				continue
			}

			names[name] = vasp.CommonName
			names[vasp.CommonName] = vasp.CommonName
		}

		if err = iter.Error(); err != nil {
			sentry.Error(c).Err(err).Str("network", network).Msg("could not iterate over VASPs")
			return nil, err
		}

		return names, nil
	}

	// Execute the parallel GDS list request
	results, errs := s.ParallelDBRequests(sentry.RequestContext(c), rpc, true)
	if len(errs) > 0 {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve VASP names"))
		return
	}

	// Deduplicate the results for the response
	out := make(map[string]string, 0)
	for _, vasps := range results {
		var (
			vaspNames map[string]string
			ok        bool
		)
		if vaspNames, ok = vasps.(map[string]string); !ok {
			sentry.Error(c).Msg("unexpected result type")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve VASP names"))
			return
		}

		// Add the name to the list
		for name, commonName := range vaspNames {
			out[name] = commonName
		}
	}

	c.JSON(http.StatusOK, out)
}

// VerifyContact is currently a passthrough helper that forwards the verify contact
// request from the user interface to the GDS that needs contact verification.
//
// @Summary Verify a VASP contact
// @Description Verify a VASP contact using a TestNet or MainNet GDS.
// @Tags GDS
// @Accept json
// @Produce json
// @Param params body api.VerifyContactParams true "Verify contact parameters"
// @Success 200 {object} api.VerifyContactReply
// @Failure 400 {object} api.Reply
// @Failure 404 {object} api.Reply
// @Failure 409 {object} api.Reply
// @Failure 500 {object} api.Reply
// @Router /verify [get]
func (s *Server) VerifyContact(c *gin.Context) {
	// Bind the parameters associated with the verify contact request
	params := &api.VerifyContactParams{}
	if err := c.ShouldBindQuery(&params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not bind request with query params")
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
	ctx, cancel := context.WithTimeout(sentry.RequestContext(c), 25*time.Second)
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
		sentry.Error(c).Str("registered_directory", params.Directory).Str("endpoint", "verify").Msg("unhandled directory")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not verify contact"))
		return
	}

	// Handle GDS errors
	if err != nil {
		log.Warn().Err(err).Str("directory", params.Directory).Msg("could not verify contact")
		serr, _ := status.FromError(err)
		switch serr.Code() {
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, api.ErrorResponse(serr.Message()))
		case codes.NotFound:
			c.JSON(http.StatusNotFound, api.ErrorResponse(serr.Message()))
		case codes.Aborted:
			c.JSON(http.StatusConflict, api.ErrorResponse(serr.Message()))
		default:
			sentry.Error(c).Err(err).Str("code", serr.Code().String()).Str("registered_directory", params.Directory).Msg("could not verify contact")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(serr.Message()))
		}
		return
	}

	log.Info().
		Str("status", rep.Status.String()).
		Str("reply", rep.Message).
		Str("directory", params.Directory).
		Msg("gds contact verification completed")

	// Create the response from the reply
	out := &api.VerifyContactReply{
		Status:  rep.Status.String(),
		Message: rep.Message,
	}

	if rep.Error != nil && rep.Error.Code != 0 {
		if out.Error, err = wire.Rewire(rep.Error); err != nil {
			sentry.Error(c).Err(err).Msg("could not rewire response error struct")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(fmt.Errorf("could not handle verify contact response from %s", params.Directory)))
			return
		}
	}

	c.JSON(http.StatusOK, out)
}

// Returns the user's current registration form if it's available
//
// @Summary Get the user's current registration form [read:vasp]
// @Description Get the registration form associated with the user's organization.
// @Tags registration
// @Produce json
// @Param params body api.RegistrationFormParams false "Load registration form parameters"
// @Success 200 {object} object "Registration form"
// @Failure 400 {object} api.Reply
// @Failure 401 {object} api.Reply
// @Failure 500 {object} api.Reply
// @Router /register [get]
func (s *Server) LoadRegisterForm(c *gin.Context) {
	var (
		err    error
		org    *records.Organization
		step   records.StepType
		params *api.RegistrationFormParams
	)

	// Load the organization from the claims
	// NOTE: this method will handle the error logging and response.
	if org, err = s.OrganizationFromClaims(c); err != nil {
		return
	}

	// Bind the parameters associated with the load registration request
	// NOTE: the step is optional and does not need to be specified
	params = &api.RegistrationFormParams{}
	if err = c.ShouldBindQuery(&params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not bind request with query params")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Convert the step into a StepType
	if step, err = records.ParseStepType(string(params.Step)); err != nil {
		sentry.Warn(c).Err(err).Msg("user requested invalid form step type")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Return the registration form, ensuring nil is not serialized.
	if org.Registration == nil {
		org.Registration = records.NewRegisterForm()
	}

	// Prepare to return the requested output.
	out := &api.RegistrationForm{
		Step: api.RegistrationFormStep(string(step)),
	}

	// If necessary, truncate the form to the specified step
	if out.Form, err = org.Registration.Truncate(step); err != nil {
		sentry.Warn(c).Err(err).Str("step", string(step)).Msg("could not truncate registration form")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	// Get any field validation errors and attach them to the form
	if verrs := org.Registration.Validate(step); verrs != nil {
		var fields records.ValidationErrors
		if errors.As(verrs, &fields) {
			out.Errors = api.FromValidationErrors(fields)
		} else {
			sentry.Warn(c).Err(err).Msg("could not validate registration form")
			c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
			return
		}
	}

	var cleaned gin.H
	if cleaned, err = out.MarshalStepJSON(); err != nil {
		sentry.Warn(c).Err(err).Str("step", string(step)).Msg("could not marshal registration form for the requested step")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, cleaned)
}

// Saves the registration form on the BFF to allow multiple users to edit the
// registration form before it is submitted to the directory service.
//
// @Summary Save a registration form to the database [update:vasp]
// @Description Save a registration form to the user's organization in the database.
// @Tags registration
// @Accept json
// @Produce json
// @Param form body object true "Registration form"
// @Success 200 {object} object "Registration form"
// @Success 204 "Empty form was provided"
// @Failure 400 {object} api.Reply
// @Failure 401 {object} api.Reply
// @Failure 500 {object} api.Reply
// @Router /register [put]
func (s *Server) SaveRegisterForm(c *gin.Context) {
	// Parse the incoming JSON data from the client request
	var (
		err  error
		step records.StepType
		form *api.RegistrationForm
		org  *records.Organization
	)

	// Unmarshal the registration form from the POST request
	form = &api.RegistrationForm{}
	if err = c.ShouldBind(form); err != nil {
		sentry.Warn(c).Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Convert the step into a StepType
	if step, err = records.ParseStepType(string(form.Step)); err != nil {
		sentry.Warn(c).Err(err).Msg("user requested invalid form step type")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Load the organization from the claims
	// NOTE: this method will handle the error logging and response.
	if org, err = s.OrganizationFromClaims(c); err != nil {
		return
	}

	// If the organization form does not exist; create a new registration form
	if org.Registration == nil {
		org.Registration = records.NewRegisterForm()
	}

	// A form should always be provided to this endpoint, the delete endpoint must be
	// used to reset a form.
	if form.Form == nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("no form was provided"))
		return
	}

	// Mark the form as started, the BFF relies on this state so the frontend should
	// capture the updated form returned from this endpoint to avoid overwriting the
	// state metadata.
	// NOTE: If an empty form was passed in, the form will not be marked as started.
	if form.Form.State != nil && form.Form.State.Started == "" {
		form.Form.State.Started = time.Now().Format(time.RFC3339)
	}

	// Prepare the response
	out := &api.RegistrationForm{
		Step: api.RegistrationFormStep(step),
	}

	// Update the registration form step that has been POSTED.
	if err = org.Registration.Update(form.Form, step); err != nil {
		// If there were validation errors, attach them to the output
		var fields records.ValidationErrors
		if errors.As(err, &fields) {
			out.Errors = api.FromValidationErrors(fields)
		} else {
			sentry.Warn(c).Err(err).Msg("could not update registration form")
			c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
			return
		}
	}

	// Update the organizations form
	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	if err = s.db.UpdateOrganization(ctx, org); err != nil {
		sentry.Error(c).Err(err).Msg("could not update organization")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not save registration form"))
		return
	}

	// Return the updated form in a 200 OK response, truncated if necessary.
	if out.Form, err = org.Registration.Truncate(step); err != nil {
		sentry.Warn(c).Err(err).Str("step", string(step)).Msg("could not truncate registration form")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	var cleaned gin.H
	if cleaned, err = out.MarshalStepJSON(); err != nil {
		sentry.Warn(c).Err(err).Str("step", string(step)).Msg("could not marshal registration form for the requested step")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, cleaned)
}

// Resets the user's current registration form to the defaults.
//
// @Summary Reset the user's current registration form [update:vasp]
// @Description Reset the registration form associated with the user's organization for the requested step.
// @Tags registration
// @Produce json
// @Param params body api.RegistrationFormParams false "Reset registration form parameters"
// @Success 200 {object} object "Registration form"
// @Failure 400 {object} api.Reply
// @Failure 401 {object} api.Reply
// @Failure 500 {object} api.Reply
// @Router /register [delete]
func (s *Server) ResetRegisterForm(c *gin.Context) {
	var (
		err    error
		org    *records.Organization
		step   records.StepType
		params *api.RegistrationFormParams
	)

	// Load the organization from the claims
	// NOTE: this method will handle the error logging and response.
	if org, err = s.OrganizationFromClaims(c); err != nil {
		return
	}

	// Bind the parameters associated with the delete registration request
	// NOTE: the step is optional and does not need to be specified
	params = &api.RegistrationFormParams{}
	if err = c.ShouldBindQuery(&params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not bind request with query params")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Convert the step into a StepType
	if step, err = records.ParseStepType(string(params.Step)); err != nil {
		sentry.Warn(c).Err(err).Msg("user requested invalid form step type")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// If the organization form does not exist; create a new registration form
	if org.Registration == nil {
		org.Registration = records.NewRegisterForm()
	}

	// Prepare the response
	out := &api.RegistrationForm{
		Step: api.RegistrationFormStep(step),
	}

	// Delete the form step by doing an update with a default form
	if err = org.Registration.Update(records.NewRegisterForm(), step); err != nil {
		// Ignore validation errors on delete
		var fields records.ValidationErrors
		if !errors.As(err, &fields) {
			sentry.Warn(c).Err(err).Msg("could not reset registration form")
			c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
			return
		}
	}

	// Update the form on the organization
	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	if err = s.db.UpdateOrganization(ctx, org); err != nil {
		sentry.Error(c).Err(err).Msg("could not update organization with reset registration form")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not reset registration form"))
		return
	}

	// Return the updated form in a 200 OK response, truncated if necessary.
	if out.Form, err = org.Registration.Truncate(step); err != nil {
		sentry.Warn(c).Err(err).Str("step", string(step)).Msg("could not truncate registration form")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	var cleaned gin.H
	if cleaned, err = out.MarshalStepJSON(); err != nil {
		sentry.Warn(c).Err(err).Str("step", string(step)).Msg("could not marshal registration form for the requested step")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, cleaned)
}

// SubmitRegistration makes a request on behalf of the user to either the TestNet or the
// MainNet GDS server based on the URL endpoint. The endpoint will first load the saved
// registration form from the front-end and will parse it for some basic validity
// constraints - it will then submit the form and return any response from the directory.
//
// @Summary Submit a registration form to a directory service [update:vasp]
// @Description Submit a registration form to the TestNet or MainNet directory service.
// @Tags registration
// @Produce json
// @Param directory path string true "Directory service to submit the registration form to (testnet or mainnet)"
// @Success 200 {object} api.RegisterReply
// @Failure 400 {object} api.Reply
// @Failure 401 {object} api.Reply
// @Failure 409 {object} api.Reply
// @Failure 500 {object} api.Reply
// @Router /register/{directory} [post]
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

	// Fetch the user from the context
	var user *management.User
	if user, err = auth.GetUserInfo(c); err != nil {
		sentry.Error(c).Err(err).Msg("submit registration requires user info; expected middleware to return 401")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch user info"))
		return
	}

	// Load the app metadata for the user for updating
	appdata := &auth.AppMetadata{}
	if err = appdata.Load(user.AppMetadata); err != nil {
		sentry.Error(c).Err(err).Msg("could not parse user app metadata")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not parse user app metadata"))
		return
	}

	// Do not allow a registration form to be submitted twice
	switch network {
	case config.TestNet:
		if org.Testnet != nil && org.Testnet.Submitted != "" {
			err = fmt.Errorf("registration form has already been submitted to the %s", network)
			sentry.Warn(c).Err(err).Str("network", network).Str("orgID", org.Id).Msg("cannot resubmit registration")
			c.JSON(http.StatusConflict, api.ErrorResponse(err))
			return
		}
	case config.MainNet:
		if org.Mainnet != nil && org.Mainnet.Submitted != "" {
			err = fmt.Errorf("registration form has already been submitted to the %s", network)
			sentry.Warn(c).Err(err).Str("network", network).Str("orgID", org.Id).Msg("cannot resubmit registration")
			c.JSON(http.StatusConflict, api.ErrorResponse(err))
			return
		}
	}

	// Validate that a registration form exists on the organization
	if org.Registration == nil || !org.Registration.ReadyToSubmit(network) {
		sentry.Warn(c).Str("orgID", org.Id).Msg("cannot submit empty or partial registration form")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("registration form is not ready to submit"))
		return
	}

	// Create the RegisterRequest to send to GDS
	req := &gds.RegisterRequest{
		Entity:           org.Registration.Entity,
		Contacts:         org.Registration.Contacts,
		Website:          org.Registration.Website,
		BusinessCategory: org.Registration.BusinessCategory,
		VaspCategories:   org.Registration.VaspCategories,
		EstablishedOn:    org.Registration.EstablishedOn,
		Trixo:            org.Registration.Trixo,
	}

	// Make the GDS request
	var rep *gds.RegisterReply
	log.Debug().Str("network", network).Msg("issuing GDS register request")
	ctx, cancel := utils.WithDeadline(context.Background())
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
			sentry.Error(c).Err(err).Str("code", serr.Code().String()).Str("network", network).Msg("could not register with directory service")
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
		RefreshToken:        true,
	}

	if rep.Error != nil && rep.Error.Code != 0 {
		if out.Error, err = wire.Rewire(rep.Error); err != nil {
			sentry.Error(c).Err(err).Str("network", network).Msg("could not rewire response error struct")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(fmt.Errorf("could not handle register response from %s", network)))
			return
		}
	}

	// If there is only an error and no vaspID, return a 409 and the response
	if out.Error != nil && out.Id == "" {
		sentry.Error(c).Err(rep.Error).Str("network", network).Msg("received unexpected GDS OK with an error in the response")
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
		appdata.VASPs.TestNet = rep.Id
	case config.MainNet:
		org.Mainnet = directoryRecord
		appdata.VASPs.MainNet = rep.Id
	}

	if err = s.db.UpdateOrganization(ctx, org); err != nil {
		sentry.Error(c).Err(err).Str("network", network).Msg("could not update organization with directory record")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete registration submission"))
		return
	}

	// Commit the user metadata updates to auth0
	if err = s.SaveAuth0AppMetadata(c.Request.Context(), *user.ID, *appdata); err != nil {
		sentry.Error(c).Err(err).Str("user_id", *user.ID).Msg("could not save user app metadata")
	}

	c.JSON(http.StatusOK, out)
}

// GetCertificates makes parallel calls to the databases to get the certificate
// information for both testnet and mainnet. If testnetID or mainnetID are empty
// strings, this will simply return a nil response for the corresponding network so
// the caller can distinguish between a non registration and an error.
func (s *Server) GetCertificates(ctx context.Context, testnetID, mainnetID string) (testnetCerts, mainnetCerts []*models.Certificate, testnetErr, mainnetErr error) {
	// Create the RPC which can do both testnet and mainnet calls
	rpc := func(ctx context.Context, db store.Store, network string) (rep interface{}, err error) {
		var vaspID string
		switch network {
		case config.TestNet:
			vaspID = testnetID
		case config.MainNet:
			vaspID = mainnetID
		default:
			return nil, fmt.Errorf("unknown network: %s", network)
		}

		if vaspID == "" {
			// The VASP is not registered for this network, so do not error and return
			// nil
			return nil, nil
		}

		// Retrieve the VASP record from the database
		var vasp *pb.VASP
		if vasp, err = db.RetrieveVASP(ctx, vaspID); err != nil {
			return nil, fmt.Errorf("could not retrieve VASP %s from database: %s", vaspID, err)
		}

		// Retrieve the certificate IDs from the VASP record
		var ids []string
		if ids, err = models.GetCertIDs(vasp); err != nil {
			return nil, fmt.Errorf("could not retrieve certificate IDs from VASP %s: %s", vaspID, err)
		}

		// Construct the list of certificate records
		certs := make([]*models.Certificate, 0)

		for _, id := range ids {
			// Retrieve the certificate record from the database
			var cert *models.Certificate
			if cert, err = db.RetrieveCert(ctx, id); err != nil {
				return nil, fmt.Errorf("could not retrieve certificate %s from database: %s", id, err)
			}
			certs = append(certs, cert)
		}

		return certs, nil
	}

	// Perform the parallel requests
	results, errs := s.ParallelDBRequests(ctx, rpc, false)
	if len(errs) != 2 || len(results) != 2 {
		err := fmt.Errorf("unexpected number of results from parallel requests: %d", len(results))
		return nil, nil, err, err
	}

	// Parse the results
	var ok bool
	if errs[0] != nil {
		testnetErr = errs[0]
	} else if results[0] != nil {
		if testnetCerts, ok = results[0].([]*models.Certificate); !ok {
			testnetErr = fmt.Errorf("unexpected testnet result type returned from parallel certificate requests: %T", results[0])
		}
	}

	if errs[1] != nil {
		mainnetErr = errs[1]
	} else if results[1] != nil {
		if mainnetCerts, ok = results[1].([]*models.Certificate); !ok {
			mainnetErr = fmt.Errorf("unexpected mainnet result type returned from parallel certificate requests: %T", results[1])
		}
	}

	return testnetCerts, mainnetCerts, testnetErr, mainnetErr
}

// Certificates returns the list of certificates for the authenticated user.
//
// @Summary List certificates for the user [read:vasp]
// @Description Returns the certificates associated with the user's organization.
// @Tags certificates
// @Produce json
// @Success 200 {object} api.CertificatesReply
// @Failure 401 {object} api.Reply
// @Failure 500 {object} api.Reply
// @Router /certificates [get]
func (s *Server) Certificates(c *gin.Context) {
	var err error

	// Get the bff claims from the context
	var claims *auth.Claims
	if claims, err = auth.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("unable to retrieve bff claims from context")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	// Extract the VASP IDs from the claims
	// Note that if testnet or mainnet are absent from the VASPs struct, the ID will
	// default to an empty string, and GetCertificates will return nil for that network
	// instead of an error.
	testnetID := claims.VASPs.TestNet
	mainnetID := claims.VASPs.MainNet

	// Get the certificate replies from the admin APIs
	testnet, mainnet, testnetErr, mainnetErr := s.GetCertificates(sentry.RequestContext(c), testnetID, mainnetID)

	// Construct the response
	out := &api.CertificatesReply{
		Error:   api.NetworkError{},
		TestNet: make([]api.Certificate, 0),
		MainNet: make([]api.Certificate, 0),
	}

	// Populate the testnet response
	if testnetErr != nil {
		out.Error.TestNet = testnetErr.Error()
	} else {
		for _, cert := range testnet {
			entry := api.Certificate{
				SerialNumber: cert.Id,
				IssuedAt:     cert.Details.NotBefore,
				ExpiresAt:    cert.Details.NotAfter,
				Revoked:      cert.Status == models.CertificateState_REVOKED,
			}

			if entry.Details, err = wire.Rewire(cert.Details); err != nil {
				sentry.Error(c).Str("cert_id", cert.Id).Err(err).Msg("could not rewire testnet certificate details")
				c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
				return
			}

			out.TestNet = append(out.TestNet, entry)
		}
	}

	// Populate the mainnet response
	if mainnetErr != nil {
		out.Error.MainNet = mainnetErr.Error()
	} else {
		for _, cert := range mainnet {
			entry := api.Certificate{
				SerialNumber: cert.Id,
				IssuedAt:     cert.Details.NotBefore,
				ExpiresAt:    cert.Details.NotAfter,
				Revoked:      cert.Status == models.CertificateState_REVOKED,
			}

			if entry.Details, err = wire.Rewire(cert.Details); err != nil {
				sentry.Error(c).Str("cert_id", cert.Id).Err(err).Msg("could not rewire mainnet certificate details")
				c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
				return
			}

			out.MainNet = append(out.MainNet, entry)
		}
	}

	c.JSON(http.StatusOK, out)
}

const (
	testnetName          = "TestNet"
	mainnetName          = "MainNet"
	supportEmail         = "support@rotational.io"
	StartRegistration    = "Start the registration and verification process for your organization to receive an X.509 Identity Certificate and become a trusted member of the TRISA network."
	CompleteRegistration = "Complete the registration process and verification process for your organization to receive an X.509 Identity Certificate and become a trusted member of the TRISA network."
	SubmitTestnet        = "Review and submit your " + testnetName + " registration."
	SubmitMainnet        = "Review and submit your " + mainnetName + " registration."
	VerifyEmails         = "Your organization's %s registration has been submitted and verification emails have been sent to the contacts specified in the form. Contacts and email addresses must be verified as the first step in the approval process. Please request that contacts verify their email addresses promptly so that the TRISA Validation Team can proceed with the validation process. Please contact TRISA support at " + supportEmail + " if contacts have not received the verification email and link."
	RegistrationPending  = "Your organization's %s registration has been received and is pending approval. The TRISA Validation Team will notify you about the outcome."
	RegistrationRejected = "Your organization's %s registration has been rejected by the TRISA Validation Team. This means your organization is not a verified member of the TRISA network and cannot communicate with other members. Please contact TRISA support at " + supportEmail + " for additional details and next steps."
	RegistrationApproved = "Your organization's %s registration has been approved by the TRISA Validation Team. Take the next steps to integrate, test, and begin sending compliance messages with TRISA-verified counterparties."
	RenewCertificate     = "Your organization's %s X.509 Identity Certificate will expire on %s. Start the renewal process to receive a new X.509 Identity Certificate and remain a trusted member of the TRISA network."
	CertificateRevoked   = "Your organization's %s X.509 Identity Certificate has been revoked by TRISA. This means your organization is no longer a verified member of the TRISA network and can no longer communicate with other members. Please contact TRISA support at " + supportEmail + " for additional details and next steps."
)

// GetVASPs makes parallel calls to the databases to retrieve VASP records from
// testnet and mainnet. If testnet or mainnet are empty strings, this will simply
// return a nil response for the corresponding network so the caller can distinguish
// between a non registration and an error.
func (s *Server) GetVASPs(ctx context.Context, testnetID, mainnetID string) (testnetVASP, mainnetVASP *pb.VASP, testnetErr, mainnetErr error) {
	// Create the RPC which can do both testnet and mainnet calls
	rpc := func(ctx context.Context, db store.Store, network string) (rep interface{}, err error) {
		var vaspID string
		switch network {
		case config.TestNet:
			vaspID = testnetID
		case config.MainNet:
			vaspID = mainnetID
		default:
			return nil, fmt.Errorf("unknown network: %s", network)
		}

		if vaspID == "" {
			// The VASP is not registered for this network, so do not error and return
			// nil
			return nil, nil
		}
		return db.RetrieveVASP(ctx, vaspID)
	}

	// Perform the parallel requests
	results, errs := s.ParallelDBRequests(ctx, rpc, false)
	if len(errs) != 2 || len(results) != 2 {
		err := fmt.Errorf("unexpected number of results from parallel requests: %d", len(results))
		return nil, nil, err, err
	}

	// Parse the results
	if errs[0] != nil {
		testnetErr = errs[0]
	} else if results[0] != nil {
		var ok bool
		if testnetVASP, ok = results[0].(*pb.VASP); !ok {
			testnetErr = fmt.Errorf("unexpected testnet result type returned from parallel RetrieveVASP requests: %T", results[0])
		}
	}

	if errs[1] != nil {
		mainnetErr = errs[1]
	} else if results[1] != nil {
		var ok bool
		if mainnetVASP, ok = results[1].(*pb.VASP); !ok {
			mainnetErr = fmt.Errorf("unexpected mainnet result type returned from parallel RetrieveVASP requests: %T", results[1])
		}
	}

	return testnetVASP, mainnetVASP, testnetErr, mainnetErr
}

// registrationMessage returns a message corresponding to the registration state of the
// VASP. These states are distinct, so only one message is returned.
func registrationMessage(vasp *pb.VASP, network string) (msg *api.AttentionMessage, err error) {
	const expireLayout = "January 2, 2006"

	if vasp == nil {
		return nil, nil
	}

	switch {
	case vasp.VerificationStatus == pb.VerificationState_SUBMITTED:
		// Verify contact emails have been sent and are pending verification
		return &api.AttentionMessage{
			Message:  fmt.Sprintf(VerifyEmails, network),
			Severity: records.AttentionSeverity_INFO.String(),
			Action:   records.AttentionAction_VERIFY_EMAILS.String(),
		}, nil
	case vasp.VerificationStatus > pb.VerificationState_SUBMITTED && vasp.VerificationStatus < pb.VerificationState_VERIFIED:
		// The VASP is pending review and certificate issuance
		return &api.AttentionMessage{
			Message:  fmt.Sprintf(RegistrationPending, network),
			Severity: records.AttentionSeverity_INFO.String(),
			Action:   records.AttentionAction_NO_ACTION.String(),
		}, nil
	case vasp.IdentityCertificate != nil && vasp.IdentityCertificate.Revoked:
		// The VASP's certificate has been revoked
		return &api.AttentionMessage{
			Message:  fmt.Sprintf(CertificateRevoked, network),
			Severity: records.AttentionSeverity_ALERT.String(),
			Action:   records.AttentionAction_CONTACT_SUPPORT.String(),
		}, nil
	case vasp.IdentityCertificate != nil:
		// Certificate has been issued, check if it is about to expire
		var expiresAt time.Time
		if expiresAt, err = time.Parse(time.RFC3339, vasp.IdentityCertificate.NotAfter); err != nil {
			return nil, err
		}

		// Warn if less than 30 days before expiration
		if time.Until(expiresAt) < 30*24*time.Hour {
			return &api.AttentionMessage{
				Message:  fmt.Sprintf(RenewCertificate, network, expiresAt.Format(expireLayout)),
				Severity: records.AttentionSeverity_WARNING.String(),
				Action:   records.AttentionAction_RENEW_CERTIFICATE.String(),
			}, nil
		}
	case vasp.VerificationStatus == pb.VerificationState_VERIFIED:
		// The VASP is verified and the certificate has been issued
		return &api.AttentionMessage{
			Message:  fmt.Sprintf(RegistrationApproved, network),
			Severity: records.AttentionSeverity_SUCCESS.String(),
			Action:   records.AttentionAction_NO_ACTION.String(),
		}, nil
	case vasp.VerificationStatus == pb.VerificationState_REJECTED:
		// The VASP has been rejected, so no certificate was issued
		return &api.AttentionMessage{
			Message:  fmt.Sprintf(RegistrationRejected, network),
			Severity: records.AttentionSeverity_ALERT.String(),
			Action:   records.AttentionAction_CONTACT_SUPPORT.String(),
		}, nil
	default:
	}
	return nil, nil
}

// Attention returns the current attention messages for the authenticated user.
//
// @Summary Get attention alerts for the user [read:vasp]
// @Description Get attention alerts for the user regarding their organization's VASP registration status.
// @Tags registration
// @Produce json
// @Success 200 {object} api.AttentionReply
// @Success 204 "No attention messages"
// @Failure 401 {object} api.Reply
// @Failure 500 {object} api.Reply
// @Router /attention [get]
func (s *Server) Attention(c *gin.Context) {
	var err error

	// Get the bff claims from the context
	var claims *auth.Claims
	if claims, err = auth.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("unable to retrieve bff claims from context")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	// Retrieve the organization from the claims
	// NOTE: This method handles the error logging and response.
	var org *records.Organization
	if org, err = s.OrganizationFromClaims(c); err != nil {
		return
	}

	// Attention messages to return
	messages := make([]*api.AttentionMessage, 0)

	// Check the registration state, at most one of these messages will be returned
	testnetSubmitted := (org.Testnet != nil && org.Testnet.Submitted != "")
	mainnetSubmitted := (org.Mainnet != nil && org.Mainnet.Submitted != "")
	switch {
	case org.Registration == nil || org.Registration.State == nil || org.Registration.State.Started == "":
		// Registration has not started
		messages = append(messages, &api.AttentionMessage{
			Message:  StartRegistration,
			Severity: records.AttentionSeverity_INFO.String(),
			Action:   records.AttentionAction_START_REGISTRATION.String(),
		})
	case !testnetSubmitted && !mainnetSubmitted:
		// Registration has started but has not been completed
		messages = append(messages, &api.AttentionMessage{
			Message:  CompleteRegistration,
			Severity: records.AttentionSeverity_INFO.String(),
			Action:   records.AttentionAction_COMPLETE_REGISTRATION.String(),
		})
	case testnetSubmitted && !mainnetSubmitted:
		// Registration is submitted for testnet but not for mainnet
		messages = append(messages, &api.AttentionMessage{
			Message:  SubmitMainnet,
			Severity: records.AttentionSeverity_INFO.String(),
			Action:   records.AttentionAction_SUBMIT_MAINNET.String(),
		})
	case !testnetSubmitted && mainnetSubmitted:
		// Registration is submitted for mainnet but not for testnet
		messages = append(messages, &api.AttentionMessage{
			Message:  SubmitTestnet,
			Severity: records.AttentionSeverity_INFO.String(),
			Action:   records.AttentionAction_SUBMIT_TESTNET.String(),
		})
	default:
	}

	// Get the VASP records from the admin APIs
	// NOTE: This will not attempt to retrieve the VASP records if the VASP ID is not
	// set in the claims for a network and a nil result will be returned instead for
	// that network.
	testnetVASP, mainnetVASP, testnetErr, mainnetErr := s.GetVASPs(sentry.RequestContext(c), claims.VASPs.TestNet, claims.VASPs.MainNet)

	if testnetErr != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(testnetErr))
		return
	}

	if mainnetErr != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(mainnetErr))
		return
	}

	// Get attention messages relating to certificates
	var testnetMsg *api.AttentionMessage
	if testnetMsg, err = registrationMessage(testnetVASP, testnetName); err != nil {
		sentry.Error(c).Err(err).Msg("could not get testnet certificate attention message")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}
	if testnetMsg != nil {
		messages = append(messages, testnetMsg)
	}

	var mainnetMsg *api.AttentionMessage
	if mainnetMsg, err = registrationMessage(mainnetVASP, mainnetName); err != nil {
		sentry.Error(c).Err(err).Msg("could not get mainnet certificate attention message")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}
	if mainnetMsg != nil {
		messages = append(messages, mainnetMsg)
	}

	// Build the response
	if len(messages) == 0 {
		c.JSON(http.StatusNoContent, nil)
	} else {
		c.JSON(http.StatusOK, &api.AttentionReply{
			Messages: messages,
		})
	}
}

// RegistrationStatus returns the registration status for both testnet and mainnet for
// the user.
//
// @Summary Get current registration status for the user [read:vasp]
// @Description Returns timestamps indicating when the user has submitted their TestNet and MainNet registrations.
// @Tags registration
// @Produce json
// @Success 200 {object} api.RegistrationStatus
// @Failure 401 {object} api.Reply
// @Failure 500 {object} api.Reply
func (s *Server) RegistrationStatus(c *gin.Context) {
	var err error

	// Retrieve the organization from the claims
	// NOTE: This method handles the error logging and response.
	var org *records.Organization
	if org, err = s.OrganizationFromClaims(c); err != nil {
		return
	}

	// Build the response
	// TODO: We should be querying the VASP record instead to allow for re-registration
	out := &api.RegistrationStatus{}
	if org.Testnet != nil && org.Testnet.Submitted != "" {
		out.TestNetSubmitted = org.Testnet.Submitted
	}
	if org.Mainnet != nil && org.Mainnet.Submitted != "" {
		out.MainNetSubmitted = org.Mainnet.Submitted
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
