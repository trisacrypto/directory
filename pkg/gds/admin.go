package gds

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/idtoken"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/trisacrypto/directory/pkg"
	admin "github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/store"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

// NewAdmin creates a new GDS admin server derived from a parent Service.
func NewAdmin(svc *Service) (a *Admin, err error) {
	// Define the base admin server
	a = &Admin{
		svc:  svc,
		conf: &svc.conf.Admin,
		db:   svc.db,
	}

	// Create the token manager
	if a.tokens, err = tokens.New(a.conf.TokenKeys); err != nil {
		return nil, err
	}

	// Create the router
	gin.SetMode(a.conf.Mode)
	a.router = gin.New()
	if err = a.setupRoutes(); err != nil {
		return nil, err
	}

	// Create the http server
	a.srv = &http.Server{
		Addr:         a.conf.BindAddr,
		Handler:      a.router,
		ErrorLog:     nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Debug().Msg("created admin api http server with gin router")
	return a, nil
}

// Admin implements the DirectoryAdministrationServer as defined by the v2 JSON API.
// This service is the primary interaction point with authorized TRISA users that are
// performing secure commands with authentication.
type Admin struct {
	sync.RWMutex
	svc     *Service             // The parent Service the admin server uses to interact with other components
	srv     *http.Server         // The HTTP server that listens on its own independent port
	conf    *config.AdminConfig  // The admin server specific configuration (alias to s.svc.conf.Admin)
	tokens  *tokens.TokenManager // A token manager that signs JWT tokens with RSA keys
	db      store.Store          // Database connection for loading objects (alias to s.svc.db)
	router  *gin.Engine          // The HTTP handler and associated middleware
	healthy bool                 // application state of the server
}

// Serve GRPC requests on the specified address.
func (s *Admin) Serve() (err error) {
	// If not enabled, ignore the call to Serve and exit without error.
	if !s.conf.Enabled {
		log.Warn().Msg("directory administration service is not enabled")
		return nil
	}

	// This service should start in maintenance mode and return unavailable.
	s.SetHealth(!s.svc.conf.Maintenance)
	if s.svc.conf.Maintenance {
		log.Warn().Msg("directory administration service starting in maintenance mode")
	}

	// Note authorization context
	log.Debug().Strs("authorized_domains", s.conf.AuthorizedDomains).Strs("allowed_origins", s.conf.AllowOrigins).Msg("authorization context")

	// Listen for TCP requests on the specified address and port
	log.Info().
		Str("listen", s.conf.BindAddr).
		Str("version", pkg.Version()).
		Msg("directory administration server started")

	if err = s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Shutdown the Directory Administration Service gracefully
func (s *Admin) Shutdown() (err error) {
	log.Debug().Msg("gracefully shutting down directory administration server")

	// Gracefully shutdown admin API server
	s.SetHealth(false)
	s.srv.SetKeepAlivesEnabled(false)

	// Require shutdown in 30 seconds without blocking
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = s.srv.Shutdown(ctx); err != nil {
		return err
	}

	log.Debug().Msg("successful shutdown of admin api server")
	return nil
}

// Routes returns the Admin API router for testing purposes.
func (s *Admin) Routes() http.Handler {
	return s.router
}

func (s *Admin) setupRoutes() (err error) {
	// Application Middleware
	s.router.Use(ginzerolog.Logger("gin"))
	s.router.Use(gin.Recovery())
	s.router.Use(s.Available())

	// Add CORS configuration
	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     s.conf.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-CSRF-TOKEN"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Route-specific middleware
	authorize := admin.Authorization(s.tokens)
	csrf := admin.DoubleCookie()

	// Add the v2 API routes
	v2 := s.router.Group("/v2")
	{
		// Heartbeat route (no authentication required)
		v2.GET("/status", s.Status)

		// Authentication and user management routes (some CSRF protection required)
		v2.GET("/authenticate", s.ProtectAuthenticate)
		v2.POST("/authenticate", csrf, s.Authenticate)
		v2.POST("/reauthenticate", csrf, s.Reauthenticate)

		// Information routes (must be authenticated)
		v2.GET("/summary", authorize, s.Summary)
		v2.GET("/autocomplete", authorize, s.Autocomplete)
		v2.GET("/reviews", authorize, s.ReviewTimeline)

		// VASP routes all must be authenticated (some CSRF protection required)
		vasps := v2.Group("/vasps", authorize)
		{
			vasps.GET("", s.ListVASPs)
			vasps.GET("/:vaspID", s.RetrieveVASP)
			vasps.POST("/:vaspID/review", csrf, s.Review)
			vasps.POST("/:vaspID/resend", csrf, s.Resend)
		}
	}

	// NotFound and NotAllowed requests
	s.router.NoRoute(admin.NotFound)
	s.router.NoMethod(admin.NotAllowed)
	return nil
}

// Set the maximum age of authentication protection cookies.
const protectAuthenticateMaxAge = time.Minute * 10

// ProtectAuthenticate prepares the front-end for submitting a login token by setting
// the double cookie tokens for CSRF protection. The front-end should call this before
// posting credentials from Google.
func (s *Admin) ProtectAuthenticate(c *gin.Context) {
	expiresAt := time.Now().Add(protectAuthenticateMaxAge).Unix()
	if err := admin.SetDoubleCookieTokens(c, expiresAt); err != nil {
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not set cookies"))
		return
	}
	c.JSON(http.StatusOK, &admin.Reply{Success: true})
}

// Authenticate expects a Google OAuth JWT token that is verified by the server. Once
// verified, the JWT claims are authenticated against the server. Provided valid claims,
// the server will issue access and referesh tokens that the client should submit in the
// Authorization header for all future requests. This method also resets the CSRF double
// cookies to ensure that max-age matches the duration of the refresh tokens.
func (s *Admin) Authenticate(c *gin.Context) {
	var (
		err       error
		in        *admin.AuthRequest
		out       *admin.AuthReply
		claims    *idtoken.Payload
		expiresAt int64
	)

	// Parse incoming JSON data from the client request
	in = new(admin.AuthRequest)
	if err = c.ShouldBind(&in); err != nil {
		log.Warn().Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(err))
		return
	}

	// Check that a credential was posted
	if in.Credential == "" {
		c.JSON(http.StatusUnauthorized, admin.ErrorResponse("invalid credentials"))
		return
	}

	// Validate the credential with Google
	if claims, err = idtoken.Validate(c.Request.Context(), in.Credential, s.conf.Audience); err != nil {
		log.Warn().Err(err).Msg("invalid credentials used for authentication")
		c.JSON(http.StatusUnauthorized, admin.ErrorResponse("invalid credentials"))
		return
	}

	// Verify that the domain is one of our authorized domains
	if err = s.checkAuthorizedDomain(claims); err != nil {
		log.Warn().Err(err).Msg("access request from unauthorized domain")
		c.JSON(http.StatusUnauthorized, admin.ErrorResponse("invalid credentials"))
		return
	}

	// At this point request has been authenticated and authorized, create credentials.
	if out, expiresAt, err = s.createAuthReply(claims); err != nil {
		// NOTE: additional error logging happens in createAuthReply
		log.Error().Err(err).Msg("could not authenticate user")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not authenticate with credentials"))
		return
	}

	// Refresh the double cookies for CSRF protection while using the access/refresh tokens
	if err := admin.SetDoubleCookieTokens(c, expiresAt); err != nil {
		log.Error().Err(err).Msg("could not set double cookie tokens")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not set cookies"))
		return
	}

	// Return successful authentication!
	c.JSON(http.StatusOK, out)
}

func (s *Admin) checkAuthorizedDomain(claims *idtoken.Payload) error {
	// Fetch the claim from the idtoken payload
	domain, ok := claims.Claims["hd"]
	if !ok {
		return errors.New("no hd claim to verify authorized domain with")
	}

	// Convert the domain into a string for verification
	domains, ok := domain.(string)
	if !ok {
		return fmt.Errorf("claim type %T unparseable", domain)
	}

	// Process the HD domain for string comparison purposes
	domains = strings.ToLower(strings.TrimSpace(domains))

	// Search the authorized domains, if found return nil
	for _, authorized := range s.conf.AuthorizedDomains {
		if domains == authorized {
			// Found an authorized domain!
			return nil
		}
	}

	return fmt.Errorf("%s is not in the configured authorized domains", domains)
}

func (s *Admin) createAuthReply(creds interface{}) (out *admin.AuthReply, expiresAt int64, err error) {
	var accessToken, refreshToken *jwt.Token

	// Create the access and refresh tokens from the claims
	if accessToken, err = s.tokens.CreateAccessToken(creds); err != nil {
		log.Error().Err(err).Msg("could not create access token")
		return nil, 0, err
	}

	if refreshToken, err = s.tokens.CreateRefreshToken(accessToken); err != nil {
		log.Error().Err(err).Msg("could not create refresh token")
		return nil, 0, err
	}

	// Sign the tokens and return the response
	out = new(admin.AuthReply)
	if out.AccessToken, err = s.tokens.Sign(accessToken); err != nil {
		log.Error().Err(err).Msg("could not sign access token")
		return nil, 0, err
	}
	if out.RefreshToken, err = s.tokens.Sign(refreshToken); err != nil {
		log.Error().Err(err).Msg("could not sign refresh token")
		return nil, 0, err
	}

	// Refresh the double cookies for CSRF protection while using the access/refresh tokens
	expiresAt = refreshToken.Claims.(*tokens.Claims).ExpiresAt
	return out, expiresAt, nil
}

// Reauthenticate allows the submission of a refresh token to reauthenticate an expired
// or expiring access token and issues a new token pair. The access token must still be
// provided in the Authorization header as a Bearer token, even if it is expired since
// the access token contains the claims that need to be reissued. The refresh token is
// posted in the request body as the credential. This method also resets the CSRF double
// cookies to ensure that the max-age matches the duration of the refresh tokens.
func (s *Admin) Reauthenticate(c *gin.Context) {
	var (
		err           error
		tks           string
		in            *admin.AuthRequest
		out           *admin.AuthReply
		expiresAt     int64
		accessClaims  *tokens.Claims
		refreshClaims *tokens.Claims
	)

	// Get the Bearer token from the Authorization header (contains access token)
	if tks, err = admin.GetAccessToken(c); err != nil {
		log.Warn().Err(err).Msg("reauthenticate called without access token")
		c.JSON(http.StatusUnauthorized, admin.ErrorResponse("request is not authorized"))
		return
	}

	// Parse the access token from the Authorization header without validating the
	// claims, e.g. it doesn't matter if the access token is expired, but it should be
	// signed correctly by the token server.
	if accessClaims, err = s.tokens.Parse(tks); err != nil {
		log.Warn().Err(err).Msg("reauthenticate called with invalid access token")
		c.JSON(http.StatusUnauthorized, admin.ErrorResponse("request is not authorized"))
		return
	}

	// Parse incoming JSON data from the client request (contains refresh token)
	in = new(admin.AuthRequest)
	if err = c.ShouldBind(&in); err != nil {
		log.Warn().Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(err))
		return
	}

	// Check that a credential was posted
	if in.Credential == "" {
		c.JSON(http.StatusUnauthorized, admin.ErrorResponse("invalid credentials"))
		return
	}

	// Validate the refresh token
	if refreshClaims, err = s.tokens.Verify(in.Credential); err != nil {
		log.Warn().Err(err).Msg("could not verify refresh token")
		c.JSON(http.StatusUnauthorized, admin.ErrorResponse("invalid credentials"))
		return
	}

	// Ensure the refresh token and admin token match
	// TODO: verify the in.Credential is a refresh token using the subject or audience
	if accessClaims.Id != refreshClaims.Id {
		log.Warn().Msg("mismatched access and refresh token pair")
		c.JSON(http.StatusUnauthorized, admin.ErrorResponse("invalid credentials"))
		return
	}

	// At this point we've validated the reauthentication and are ready to reissue tokens
	if out, expiresAt, err = s.createAuthReply(accessClaims); err != nil {
		// NOTE: additional error logging happens in createAuthReply
		log.Error().Err(err).Msg("could not reauthenticate user")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not reauthenticate with credentials"))
		return
	}

	if err := admin.SetDoubleCookieTokens(c, expiresAt); err != nil {
		log.Error().Err(err).Msg("could not set double cookie tokens")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not set cookies"))
		return
	}

	// Return successful reauthentication!
	c.JSON(http.StatusOK, out)
}

// Summary provides aggregate statistics that describe the state of the GDS.
func (s *Admin) Summary(c *gin.Context) {
	// Prepare the output response
	out := &admin.SummaryReply{
		Statuses: make(map[string]int),
		CertReqs: make(map[string]int),
	}

	// Query the list of VASPs from the data store to perform aggregation counts.
	iter := s.db.ListVASPs()
	for iter.Next() {
		// Fetch VASP from the database
		var vasp *pb.VASP
		if vasp = iter.VASP(); vasp == nil {
			// VASP could not be parsed; error logged in VASP() method continue iteration
			continue
		}

		// Count VASPs
		out.VASPsCount++

		// Count contacts
		contacts := []*pb.Contact{
			vasp.Contacts.Administrative, vasp.Contacts.Legal,
			vasp.Contacts.Technical, vasp.Contacts.Billing,
		}

		for _, contact := range contacts {
			if contact != nil && contact.Email != "" {
				out.ContactsCount++
				if _, verified, _ := models.GetContactVerification(contact); verified {
					out.VerifiedContacts++
				}
			}
		}

		// Count Statuses and any status that is "pending" -- awaiting action by a reviewer.
		out.Statuses[vasp.VerificationStatus.String()]++
		if int32(vasp.VerificationStatus) < int32(pb.VerificationState_VERIFIED) || vasp.VerificationStatus == pb.VerificationState_APPEALED {
			out.PendingRegistrations++
		}

	}

	if err := iter.Error(); err != nil {
		iter.Release()
		log.Warn().Err(err).Msg("could not iterate over vasps in store")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse(err))
		return
	}
	iter.Release()

	// Loop over certificate requests next
	iter2 := s.db.ListCertReqs()
	for iter2.Next() {
		// Fetch CertificateRequest from the database
		var certreq *models.CertificateRequest
		if certreq = iter2.CertReq(); certreq == nil {
			// CertificateRequest could not be parsed; error logged in CertReq() method continue iteration
			continue
		}

		out.CertReqs[certreq.Status.String()]++
		if certreq.Status == models.CertificateRequestState_COMPLETED {
			out.CertificatesIssued++
		}
	}

	if err := iter2.Error(); err != nil {
		iter2.Release()
		log.Warn().Err(err).Msg("could not iterate over certreqs in store")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse(err))
		return
	}
	iter2.Release()

	// Successful request, return the VASP list JSON data
	c.JSON(http.StatusOK, out)
}

// Autocomplete returns a mapping of name to VASP UUID for the search bar.
func (s *Admin) Autocomplete(c *gin.Context) {
	// Prepare the output response
	out := &admin.AutocompleteReply{
		Names: make(map[string]string),
	}

	// Query the list of VASPs from the data store to perform aggregation counts.
	// NOTE: we could have just queried the names index, which would be a lot faster
	// than iterating over the VASPs; if the UI requires more complex information
	// storage then the VASP iteration is better (or a better index). If it doesn't,
	// then this should be refactored to simply fetch the index and return it.
	iter := s.db.ListVASPs()
	defer iter.Release()
	for iter.Next() {
		// Fetch VASP from the database
		var vasp *pb.VASP
		if vasp = iter.VASP(); vasp == nil {
			// VASP could not be parsed; error logged in VASP() method continue iteration
			continue
		}

		// Add top level names to the autocomplete
		out.Names[vasp.CommonName] = vasp.Id
		out.Names[vasp.Website] = vasp.Website

		// Add all legal person names
		if vasp.Entity != nil {
			for _, name := range vasp.Entity.Names() {
				if _, ok := out.Names[name]; !ok {
					out.Names[name] = vasp.Id
				} else {
					log.Warn().Str("name", name).Msg("duplicate name detected")
				}
			}
		}
	}

	if err := iter.Error(); err != nil {
		log.Warn().Err(err).Msg("could not iterate over vasps in store")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse(err))
		return
	}

	// In case any of the names were empty string, delete it (no guard required)
	delete(out.Names, "")

	// Successful request, return the VASP list JSON data
	c.JSON(http.StatusOK, out)
}

// ListVASPs returns a paginated, summary data structure of all VASPs managed by the
// directory service. This is an authenticated endpoint that is used to support the
// Admin UI and facilitate the review and registration process.
func (s *Admin) ListVASPs(c *gin.Context) {
	var (
		err error
		in  *admin.ListVASPsParams
		out *admin.ListVASPsReply
	)

	in = new(admin.ListVASPsParams)
	if err = c.ShouldBindQuery(&in); err != nil {
		log.Warn().Err(err).Msg("could not bind request with query params")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(err))
		return
	}

	// Determine status filter
	var status pb.VerificationState
	if in.Status != "" {
		in.Status = strings.ToUpper(strings.ReplaceAll(in.Status, " ", "_"))
		sn, ok := pb.VerificationState_value[in.Status]
		if !ok {
			log.Warn().Str("status", in.Status).Msg("unknown verification status")
			c.JSON(http.StatusBadRequest, admin.ErrorResponse(fmt.Errorf("unknown verification status %q", in.Status)))
			return
		}
		status = pb.VerificationState(sn)
	}

	// Set pagination defaults if not specified in query
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 100
	}

	// Determine pagination index range (indexed by 1)
	minIndex := (in.Page - 1) * in.PageSize
	maxIndex := minIndex + in.PageSize
	log.Debug().Int("page", in.Page).Int("page_size", in.PageSize).Int("min_index", minIndex).Int("max_index", maxIndex).Msg("paginating vasps")

	out = &admin.ListVASPsReply{
		VASPs:    make([]admin.VASPSnippet, 0),
		Page:     in.Page,
		PageSize: in.PageSize,
	}

	// Query the list of VASPs from the data store
	iter := s.db.ListVASPs()
	defer iter.Release()
	for iter.Next() {
		out.Count++
		if out.Count >= minIndex && out.Count < maxIndex {
			// In the page range so add to the list reply
			// Fetch VASP from the database
			var vasp *pb.VASP
			if vasp = iter.VASP(); vasp == nil {
				// VASP could not be parsed; error logged in VASP() method continue iteration
				continue
			}

			// Check the status before continuing
			if status != pb.VerificationState_NO_VERIFICATION && vasp.VerificationStatus != status {
				continue
			}

			// Build the snippet
			snippet := admin.VASPSnippet{
				ID:                 vasp.Id,
				CommonName:         vasp.CommonName,
				VerificationStatus: vasp.VerificationStatus.String(),
				LastUpdated:        vasp.LastUpdated,
				Traveler:           models.IsTraveler(vasp),
			}

			// Name is a computed value, ignore errors in finding the name.
			snippet.Name, _ = vasp.Name()

			// Add verified contacts to snippet
			snippet.VerifiedContacts = models.ContactVerifications(vasp)

			// Append to list in reply
			out.VASPs = append(out.VASPs, snippet)
		}
	}

	if err = iter.Error(); err != nil {
		log.Warn().Err(err).Msg("could not iterate over vasps in store")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse(err))
		return
	}

	// Successful request, return the VASP list JSON data
	c.JSON(http.StatusOK, out)
}

func (s *Admin) RetrieveVASP(c *gin.Context) {
	var (
		err    error
		vaspID string
		vasp   *pb.VASP
		out    *admin.RetrieveVASPReply
	)
	// Get vaspID from the URL
	vaspID = c.Param("vaspID")

	// Attempt to fetch the VASP from the database
	if vasp, err = s.db.RetrieveVASP(vaspID); err != nil {
		log.Warn().Err(err).Str("id", vaspID).Msg("could not retrieve vasp")
		c.JSON(http.StatusNotFound, admin.ErrorResponse("could not retrieve VASP record by ID"))
		return
	}

	// Create the response to send back
	out = &admin.RetrieveVASPReply{
		VerifiedContacts: models.VerifiedContacts(vasp),
		Traveler:         models.IsTraveler(vasp),
	}
	if out.Name, err = vasp.Name(); err != nil {
		// This is a serious data validation error that needs to be addressed ASAP by
		// the operations team but should not block this API return.
		log.Error().Err(err).Msg("could not get VASP name")
	}

	// Remove extra data from the VASP
	// Must be done after verified contacts is computed
	// This is safe because nothing is saved back to the database
	vasp.Extra = nil
	if vasp.Contacts.Administrative != nil {
		vasp.Contacts.Administrative.Extra = nil
	}
	if vasp.Contacts.Legal != nil {
		vasp.Contacts.Legal.Extra = nil
	}
	if vasp.Contacts.Technical != nil {
		vasp.Contacts.Technical.Extra = nil
	}
	if vasp.Contacts.Billing != nil {
		vasp.Contacts.Billing.Extra = nil
	}

	// Serialize the VASP from protojson
	jsonpb := protojson.MarshalOptions{
		Multiline:       false,
		AllowPartial:    true,
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
	}

	var data []byte
	if data, err = jsonpb.Marshal(vasp); err != nil {
		log.Warn().Err(err).Str("id", vaspID).Msg("could marshal vasp json")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not create VASP json detail"))
		return
	}

	// Remarshal the JSON (unnecessary work, but done to make things easier)
	if err = json.Unmarshal(data, &out.VASP); err != nil {
		log.Warn().Err(err).Str("id", vaspID).Msg("could unmarshal vasp json")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not create VASP json detail"))
		return
	}

	// Successful request, return the VASP detail JSON data
	c.JSON(http.StatusOK, out)
}

// Review a registration request and either accept or reject it. On accept, the
// certificate request that was created on verify is used to send a Sectigo request and
// the certificate manager process watches it until the certificate has been issued. On
// reject, the VASP and certificate request records are deleted and the reject reason is
// sent to the technical contact.
func (s *Admin) Review(c *gin.Context) {
	var (
		err    error
		in     *admin.ReviewRequest
		out    *admin.ReviewReply
		vasp   *pb.VASP
		vaspID string
	)

	// Get vaspID from the URL
	vaspID = c.Param("vaspID")

	// Parse incoming JSON data from the client request
	in = new(admin.ReviewRequest)
	if err := c.ShouldBind(&in); err != nil {
		log.Warn().Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(err))
		return
	}

	// Validate review request
	if in.ID != "" && in.ID != vaspID {
		log.Warn().Msg("mismatched request ID and URL")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse("the request ID does not match the URL endpoint"))
		return
	}

	if in.AdminVerificationToken == "" {
		log.Warn().Msg("no verification token specified")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse("the admin verification token is required"))
		return
	}

	if !in.Accept && in.RejectReason == "" {
		log.Warn().Msg("missing reject reason")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse("if rejecting the request, a reason must be supplied"))
		return
	}

	// Lookup the VASP record associated with the request
	if vasp, err = s.db.RetrieveVASP(vaspID); err != nil {
		log.Warn().Err(err).Str("id", vaspID).Msg("could not retrieve vasp")
		c.JSON(http.StatusNotFound, admin.ErrorResponse("could not retrieve VASP record by ID"))
		return
	}

	// Check that the administration verification token is correct
	var adminVerificationToken string
	if adminVerificationToken, err = models.GetAdminVerificationToken(vasp); err != nil {
		log.Error().Err(err).Msg("could not retrieve admin token from extra data field on VASP")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not retrieve admin token from data"))
		return
	}
	if in.AdminVerificationToken != adminVerificationToken {
		log.Warn().Err(err).Str("vasp", vaspID).Msg("incorrect admin verification token")
		c.JSON(http.StatusUnauthorized, admin.ErrorResponse("admin verification token not accepted"))
		return
	}

	// Retrieve user claims for access to provided user info
	var claims *tokens.Claims
	value, exists := c.Get(admin.UserClaims)
	if exists && value != nil {
		var ok bool
		claims, ok = value.(*tokens.Claims)
		if !ok {
			err = fmt.Errorf("claims is an incorrect type, expecting *tokens.Claims found %T", value)
			log.Error().Err(err).Msg("could not retrieve user claims")
			c.JSON(http.StatusInternalServerError, admin.ErrorResponse("unable to retrieve user info"))
			return
		}
	} else {
		log.Error().Err(fmt.Errorf("no user claims in context")).Msg("could not retrieve user claims")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("unable to retrieve user info"))
		return
	}

	// Accept or reject the request
	out = &admin.ReviewReply{}
	if in.Accept {
		if out.Message, err = s.acceptRegistration(vasp, claims); err != nil {
			log.Error().Err(err).Msg("could not accept VASP registration")
			c.JSON(http.StatusInternalServerError, admin.ErrorResponse("unable to accept VASP registration request"))
			return
		}
	} else {
		if out.Message, err = s.rejectRegistration(vasp, in.RejectReason, claims); err != nil {
			log.Error().Err(err).Msg("could not reject VASP registration")
			c.JSON(http.StatusInternalServerError, admin.ErrorResponse("unable to reject VASP registration request"))
			return
		}
	}

	name, _ := vasp.Name()
	out.Status = vasp.VerificationStatus.String()
	log.Info().Str("vasp", vasp.Id).Str("name", name).Bool("accepted", in.Accept).Msg("registration reviewed")
	c.JSON(http.StatusOK, out)
}

// Accept the VASP registration and begin the certificate issuance process.
func (s *Admin) acceptRegistration(vasp *pb.VASP, claims *tokens.Claims) (msg string, err error) {
	// Change the VASP verification status
	if err = models.SetAdminVerificationToken(vasp, ""); err != nil {
		return "", err
	}
	vasp.VerifiedOn = time.Now().Format(time.RFC3339)
	if err := models.UpdateVerificationStatus(vasp, pb.VerificationState_REVIEWED, "registration request received", claims.Email); err != nil {
		return "", err
	}
	if err = s.db.UpdateVASP(vasp); err != nil {
		return "", err
	}

	// Mark any initialized certificate requests for this VASP as ready to submit
	// NOTE: there should only be one certificate request per VASP, but no errors occur
	// if there are more than one (other than a logged warning).
	var ncertreqs int
	careqs := s.db.ListCertReqs()

	for careqs.Next() {
		req := careqs.CertReq()
		if req != nil && req.Vasp == vasp.Id && req.Status == models.CertificateRequestState_INITIALIZED {
			// TODO: Replace "email" in the source parameter with user email address.
			if err = models.UpdateCertificateRequestStatus(req, models.CertificateRequestState_READY_TO_SUBMIT, "registration request received", "email"); err != nil {
				return "", err
			}
			if err = s.db.UpdateCertReq(req); err != nil {
				return "", err
			}
			ncertreqs++
		}
	}

	if err = careqs.Error(); err != nil {
		careqs.Release()
		return "", err
	}
	careqs.Release()

	switch ncertreqs {
	case 0:
		return "", errors.New("no certificate requests found for VASP registration")
	case 1:
		log.Debug().Str("vasp", vasp.Id).Msg("certificate request marked as ready to submit")
	default:
		log.Warn().Str("vasp", vasp.Id).Int("requests", ncertreqs).Msg("multiple certificate requests marked as ready to submit")
	}

	// Send successful response
	var name string
	if name, err = vasp.Name(); err != nil {
		name = vasp.Id
	}
	return fmt.Sprintf("registration request for %s has been approved and a Sectigo certificate will be requested", name), nil
}

// Reject the VASP registration and notify the contacts of the result.
func (s *Admin) rejectRegistration(vasp *pb.VASP, reason string, claims *tokens.Claims) (msg string, err error) {
	// Change the VASP verification status
	if err = models.SetAdminVerificationToken(vasp, ""); err != nil {
		return "", err
	}
	if err := models.UpdateVerificationStatus(vasp, pb.VerificationState_REJECTED, "registration rejected", claims.Email); err != nil {
		return "", err
	}
	if err = s.db.UpdateVASP(vasp); err != nil {
		return "", err
	}

	// Delete all pending certificate requests
	var ncertreqs int
	careqs := s.db.ListCertReqs()

	for careqs.Next() {
		req := careqs.CertReq()
		if req != nil && req.Vasp == vasp.Id {
			if err = s.db.DeleteCertReq(req.Id); err != nil {
				log.Error().Err(err).Str("id", req.Id).Msg("could not delete certificate request")
			}
			ncertreqs++
		}
	}

	if err = careqs.Error(); err != nil {
		careqs.Release()
		return "", err
	}
	careqs.Release()

	// Log deletion of certificate requests
	switch ncertreqs {
	case 0:
		log.Warn().Str("vasp", vasp.Id).Msg("no certificate requests deleted")
	case 1:
		log.Debug().Str("vasp", vasp.Id).Msg("certificate request deleted")
	default:
		log.Warn().Str("vasp", vasp.Id).Msg("multiple certificate requests deleted")
	}

	// Notify the VASP contacts that the registration request has been rejected.
	if _, err = s.svc.email.SendRejectRegistration(vasp, reason); err != nil {
		return "", err
	}

	// Send successful response
	var name string
	if name, err = vasp.Name(); err != nil {
		name = vasp.Id
	}
	return fmt.Sprintf("registration request for %s has been rejected and its contacts notified", name), nil
}

// Resend emails in case they went to spam or the initial email send failed.
func (s *Admin) Resend(c *gin.Context) {
	var (
		err    error
		in     *admin.ResendRequest
		out    *admin.ResendReply
		vasp   *pb.VASP
		vaspID string
	)

	// Get vaspID from the URL
	vaspID = c.Param("vaspID")

	// Parse incoming JSON data from the client request
	in = new(admin.ResendRequest)
	if err := c.ShouldBind(&in); err != nil {
		log.Warn().Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(err))
		return
	}

	// Validate resend request
	if in.ID != "" && in.ID != vaspID {
		log.Warn().Str("id", in.ID).Str("vasp_id", vaspID).Msg("mismatched request ID and URL")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse("the request ID does not match the URL endpoint"))
		return
	}

	// Lookup the VASP record associated with the resend request
	if vasp, err = s.db.RetrieveVASP(vaspID); err != nil {
		log.Warn().Err(err).Str("id", vaspID).Msg("could not retrieve vasp")
		c.JSON(http.StatusNotFound, admin.ErrorResponse("could not retrieve VASP record by ID"))
		return
	}

	// Handle different resend request types
	out = &admin.ResendReply{}
	switch in.Action {
	case admin.ResendVerifyContact:
		if out.Sent, err = s.svc.email.SendVerifyContacts(vasp); err != nil {
			log.Warn().Err(err).Int("sent", out.Sent).Msg("could not resend verify contacts emails")
			c.JSON(http.StatusInternalServerError, admin.ErrorResponse(fmt.Errorf("could not resend contact verification emails: %s", err)))
			return
		}
		out.Message = "contact verification emails resent to all unverified contacts"

	case admin.ResendReview:
		if out.Sent, err = s.svc.email.SendReviewRequest(vasp); err != nil {
			log.Warn().Err(err).Int("sent", out.Sent).Msg("could not resend review request")
			c.JSON(http.StatusInternalServerError, admin.ErrorResponse(fmt.Errorf("could not resend review request: %s", err)))
			return
		}
		out.Message = "review request resent to TRISA admins"

	case admin.ResendDeliverCerts:
		// TODO: check verification state and cert request state
		// TODO: in order to implement this, we'd have to fetch the certs from Google Secrets
		// TODO: if implemented, log which contact was sent the certs (e.g. technical, admin, etc.)
		// TODO: when above implemented, also log which contact was sent certs in acceptRegistration
		log.Warn().Msg("resend cert delivery not yet implemented")
		c.JSON(http.StatusNotImplemented, admin.ErrorResponse("resend cert delivery not yet implemented"))
		return

	case admin.ResendRejection:
		// Only send a rejection email if we're in the rejected state
		if vasp.VerificationStatus != pb.VerificationState_REJECTED {
			log.Warn().Err(err).Str("status", vasp.VerificationStatus.String()).Msg("cannot resend rejection emails in current state")
			c.JSON(http.StatusBadRequest, admin.ErrorResponse("VASP record verification status cannot send rejection email"))
			return
		}

		// A reason must be specified to send a rejection email (it's not stored)
		if in.Reason == "" {
			log.Warn().Str("resend_type", string(in.Action)).Msg("invalid resend request: missing reason argument")
			c.JSON(http.StatusBadRequest, admin.ErrorResponse("must specify reason for rejection to resend email"))
			return
		}
		if out.Sent, err = s.svc.email.SendRejectRegistration(vasp, in.Reason); err != nil {
			log.Warn().Err(err).Int("sent", out.Sent).Msg("could not resend rejection emails")
			c.JSON(http.StatusInternalServerError, admin.ErrorResponse(fmt.Errorf("could not resend rejection emails: %s", err)))
			return
		}
		out.Message = "rejection emails resent to all verified contacts"

	default:
		log.Warn().Str("resend_type", string(in.Action)).Msg("invalid resend request: unhandled resend request type")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(fmt.Errorf("unknown resend request type %q", in.Action)))
		return
	}

	log.Info().Str("id", vasp.Id).Int("sent", out.Sent).Str("resend_type", string(in.Action)).Msg("resend request complete")
	c.JSON(http.StatusOK, out)
}

// ReviewTimeline returns a list of time series records containing registration state counts by week.
func (s *Admin) ReviewTimeline(c *gin.Context) {
	const timeFormat = "YYYY-MM-DD"
	var (
		err          error
		in           *admin.ReviewTimelineParams
		out          *admin.ReviewTimelineReply
		startTime    time.Time
		endTime      time.Time
		earliestTime time.Time
		weekTime     time.Time
		numWeeks     int
		vaspCounts   []map[string]bool
	)

	// Get request parameters
	in = new(admin.ReviewTimelineParams)
	if err = c.ShouldBindQuery(&in); err != nil {
		log.Warn().Err(err).Msg("could not bind request with query params")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(err))
		return
	}

	// Default value for start date
	earliestTime = time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Local)
	if in.Start != "" {
		// Parse start date
		if startTime, err = time.Parse(timeFormat, in.Start); err != nil {
			log.Warn().Err(err).Msg("could not parse start date")
			c.JSON(http.StatusBadRequest, admin.ErrorResponse(fmt.Errorf("invalid start date: %s", in.Start)))
		}
		// Hard limit on the earliest date to avoid making the server do unnecessary work
		if startTime.Before(earliestTime) {
			startTime = earliestTime
		}
	} else {
		startTime = earliestTime
	}

	if in.End != "" {
		// Parse end date
		if endTime, err = time.Parse(timeFormat, in.End); err != nil {
			log.Warn().Err(err).Msg("could not parse end date")
			c.JSON(http.StatusBadRequest, admin.ErrorResponse(fmt.Errorf("invalid end date: %s", in.End)))
		}
	} else {
		// Default value for end date
		endTime = time.Now()
	}
	if startTime.After(endTime) {
		log.Warn().Err(fmt.Errorf("start date after end date")).Msg("invalid timeline request: start date can't be after current date")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(fmt.Errorf("start date must be before end date")))
	}

	// Initialize required counting structs
	numWeeks = int(endTime.Sub(startTime).Hours()/24/7) + 1
	vaspCounts = make([]map[string]bool, numWeeks)
	out = &admin.ReviewTimelineReply{
		Weeks: make([]admin.ReviewTimelineRecord, numWeeks),
	}
	weekTime = startTime
	for i := 0; i < numWeeks; i++ {
		weekDate := weekTime.Format(timeFormat)
		record := admin.ReviewTimelineRecord{
			Week:          weekDate,
			VASPsCount:    0,
			Registrations: make(map[string]int),
		}
		var s int32
		for s = 0; s <= int32(pb.VerificationState_ERRORED); s++ {
			record.Registrations[pb.VerificationState_name[s]] = 0
		}
		out.Weeks[i] = record
		vaspCounts[i] = make(map[string]bool)
		weekTime = weekTime.Add(time.Hour * 24 * 7)
	}

	// Iterate over the VASPs and count registrations
	iter := s.db.ListVASPs()
	defer iter.Release()
	for iter.Next() {
		// Fetch VASP from the database
		var vasp *pb.VASP
		if vasp = iter.VASP(); vasp == nil {
			// VASP could not be parsed; error logged in VASP() method continue iteration
			continue
		}

		// Get VASP audit log
		var auditLog []*models.AuditLogEntry
		if auditLog, err = models.GetAuditLog(vasp); err != nil {
			log.Warn().Err(err).Msg("could not retrieve audit log for vasp")
			continue
		}

		// Iterate over VASP audit log and count registrations
		for _, entry := range auditLog {
			var timestamp time.Time
			if timestamp, err = time.Parse(time.RFC3339, entry.Timestamp); err != nil {
				log.Warn().Err(err).Msg("could not parse timestamp in audit log entry")
				continue
			}
			weekNum := int(timestamp.Sub(startTime).Hours() / 24 / 7)

			// Count VASP if we haven't seen it before
			if _, exists := vaspCounts[weekNum][vasp.Id]; !exists {
				vaspCounts[weekNum][vasp.Id] = true
				out.Weeks[weekNum].VASPsCount++
			}

			// Count registration state if it changed
			if entry.PreviousState != entry.CurrentState {
				out.Weeks[weekNum].Registrations[entry.CurrentState.String()]++
			}
		}
	}

	if err := iter.Error(); err != nil {
		log.Warn().Err(err).Msg("could not iterate over vasps in store")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, out)
}

const (
	serverStatusOK          = "ok"
	serverStatusMaintenance = "maintenance"
)

// Get current counts of registration statuses and certificate requests.
func (s *Admin) Status(c *gin.Context) {
	c.JSON(http.StatusOK, admin.StatusReply{
		Status:    serverStatusOK,
		Timestamp: time.Now(),
		Version:   pkg.Version(),
	})
}

// SetHealth sets the health status on the API server, putting it into unavailable mode
// if health is false, and removing maintenance mode if health is true.
func (s *Admin) SetHealth(health bool) {
	s.Lock()
	s.healthy = health
	s.Unlock()
	log.Debug().Bool("health", health).Msg("admin api server health set")
}

// Available is middleware that uses the healthy boolean to return a service unavailable
// http status code if the server is shutting down. It does this before all routes to
// ensure that complex handling doesn't bog down the server.
func (s *Admin) Available() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check health status (if unhealthy, assume maintenance mode)
		s.RLock()
		if !s.healthy {
			c.JSON(http.StatusServiceUnavailable, admin.StatusReply{
				Status:    serverStatusMaintenance,
				Timestamp: time.Now(),
				Version:   pkg.Version(),
			})
			c.Abort()
			s.RUnlock()
			return
		}
		s.RUnlock()
		c.Next()
	}
}
