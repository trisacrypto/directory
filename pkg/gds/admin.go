package gds

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/idtoken"

	"github.com/trisacrypto/directory/pkg"
	admin "github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/secrets"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
	"github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/store"
	storeerrors "github.com/trisacrypto/directory/pkg/store/errors"
	"github.com/trisacrypto/directory/pkg/utils"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	"github.com/trisacrypto/trisa/pkg/ivms101"
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
	if a.tokens, err = tokens.New(a.conf.TokenKeys, a.conf.Audience); err != nil {
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
	log.Debug().
		Strs("authorized_domains", s.conf.Oauth.AuthorizedEmailDomains).
		Strs("allowed_origins", s.conf.AllowOrigins).
		Msg("authorization context")

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

func (s *Admin) setupRoutes() (err error) {
	var (
		tags      gin.HandlerFunc
		tracing   gin.HandlerFunc
		adminTags map[string]string
	)

	if s.svc.conf.Sentry.UseSentry() {
		adminTags = map[string]string{"service": "admin"}
		tags = sentry.UseTags(adminTags)
	}

	if s.svc.conf.Sentry.UsePerformanceTracking() {
		tracing = sentry.TrackPerformance(adminTags)
	}

	// Application Middleware
	// NOTE: ordering is very important to how middleware is handled.
	middlewares := []gin.HandlerFunc{
		// Logging should be outside so we can record the complete latency of requests.
		// Note: logging panics will not recover.
		logger.GinLogger("gds_admin_v2"),

		// Panic recovery middleware; note: gin middleware needs to be added before sentry
		gin.Recovery(),
		sentrygin.New(sentrygin.Options{
			Repanic:         true,
			WaitForDelivery: false,
		}),

		// Add searchable tags to the sentry context.
		tags,

		// Tracing helps us with our performance metrics and should be as early in the
		// chain as possible. It is after recovery to ensure trace panics recover.
		tracing,

		// CORS configuration allows the front-end to make cross-origin requests.
		cors.New(cors.Config{
			AllowOrigins:     s.conf.AllowOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-CSRF-TOKEN", "sentry-trace", "baggage"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),

		// Maintenance mode handling - does not require authentication
		s.Available(),
	}

	// Add the middleware to the router
	for _, middleware := range middlewares {
		if middleware != nil {
			s.router.Use(middleware)
		}
	}

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
		v2.GET("/countries", authorize, s.ListCountries)

		// VASP routes all must be authenticated (some CSRF protection required)
		vasps := v2.Group("/vasps", authorize)
		{
			vasps.GET("", s.ListVASPs)
			vasps.GET("/:vaspID", s.RetrieveVASP)
			vasps.PATCH("/:vaspID", csrf, s.UpdateVASP)
			vasps.DELETE("/:vaspID", csrf, s.DeleteVASP)
			vasps.GET("/:vaspID/certificates", s.ListCertificates)
			vasps.GET("/:vaspID/review", s.ReviewToken)
			vasps.POST("/:vaspID/review", csrf, s.Review)
			vasps.POST("/:vaspID/resend", csrf, s.Resend)

			contacts := vasps.Group("/:vaspID/contacts")
			{
				contacts.PUT("/:kind", csrf, s.ReplaceContact)
				contacts.DELETE("/:kind", csrf, s.DeleteContact)
			}

			notes := vasps.Group("/:vaspID/notes")
			{
				notes.GET("", s.ListReviewNotes)
				notes.POST("", csrf, s.CreateReviewNote)
				notes.PUT("/:noteID", csrf, s.UpdateReviewNote)
				notes.DELETE("/:noteID", csrf, s.DeleteReviewNote)
			}
		}
	}

	// NotFound and NotAllowed requests
	s.router.NoRoute(admin.NotFound)
	s.router.NoMethod(admin.NotAllowed)
	return nil
}

// Retrieve user claims from the Context for access to provided user info.
func (s *Admin) getClaims(c *gin.Context) (claims *tokens.Claims, err error) {
	value, exists := c.Get(admin.UserClaims)
	if exists && value != nil {
		var ok bool
		if claims, ok = value.(*tokens.Claims); !ok {
			return nil, fmt.Errorf("claims is an incorrect type, expecting *tokens.Claims found %T", value)
		}
	} else {
		return nil, errors.New("no user claims in context")
	}
	return claims, nil
}

// Set the maximum age of authentication protection cookies.
const protectAuthenticateMaxAge = time.Minute * 10

// ProtectAuthenticate prepares the front-end for submitting a login token by setting
// the double cookie tokens for CSRF protection. The front-end should call this before
// posting credentials from Google.
func (s *Admin) ProtectAuthenticate(c *gin.Context) {
	expiresAt := time.Now().Add(protectAuthenticateMaxAge)
	if err := admin.SetDoubleCookieTokens(c, s.conf.CookieDomain, expiresAt); err != nil {
		sentry.Error(c).Err(err).Msg("could not set cookies")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not set cookies"))
		return
	}
	c.JSON(http.StatusOK, &admin.Reply{Success: true})
}

// Authenticate expects a Google OAuth JWT token that is verified by the server. Once
// verified, the JWT claims are authenticated against the server. Provided valid claims,
// the server will issue access and refresh tokens that the client should submit in the
// Authorization header for all future requests. This method also resets the CSRF double
// cookies to ensure that max-age matches the duration of the refresh tokens.
func (s *Admin) Authenticate(c *gin.Context) {
	var (
		err       error
		in        *admin.AuthRequest
		out       *admin.AuthReply
		claims    *idtoken.Payload
		expiresAt time.Time
	)

	// Parse incoming JSON data from the client request
	in = new(admin.AuthRequest)
	if err = c.ShouldBind(&in); err != nil {
		sentry.Warn(c).Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(err))
		return
	}

	// Check that a credential was posted
	if in.Credential == "" {
		c.JSON(http.StatusUnauthorized, admin.ErrorResponse("invalid credentials"))
		return
	}

	// Validate the credential with Google
	if claims, err = s.tokens.Validate(c.Request.Context(), in.Credential, s.conf.Oauth.GoogleAudience); err != nil {
		sentry.Warn(c).Err(err).Msg("invalid credentials used for authentication")
		c.JSON(http.StatusUnauthorized, admin.ErrorResponse("invalid credentials"))
		return
	}

	// Verify that the domain is one of our authorized domains
	if err = s.checkAuthorizedDomain(claims); err != nil {
		sentry.Warn(c).Err(err).Msg("access request from unauthorized domain")
		c.JSON(http.StatusUnauthorized, admin.ErrorResponse("invalid credentials"))
		return
	}

	// At this point request has been authenticated and authorized, create credentials.
	if out, expiresAt, err = s.createAuthReply(claims); err != nil {
		// NOTE: additional error logging happens in createAuthReply
		sentry.Error(c).Err(err).Msg("could not authenticate user")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not authenticate with credentials"))
		return
	}

	// Refresh the double cookies for CSRF protection while using the access/refresh tokens
	if err := admin.SetDoubleCookieTokens(c, s.conf.CookieDomain, expiresAt); err != nil {
		sentry.Error(c).Err(err).Msg("could not set double cookie tokens")
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
	for _, authorized := range s.conf.Oauth.AuthorizedEmailDomains {
		if domains == authorized {
			// Found an authorized domain!
			return nil
		}
	}

	return fmt.Errorf("%s is not in the configured authorized domains", domains)
}

func (s *Admin) createAuthReply(creds interface{}) (out *admin.AuthReply, expiresAt time.Time, err error) {
	var accessToken, refreshToken *jwt.Token

	// Create the access and refresh tokens from the claims
	if accessToken, err = s.tokens.CreateAccessToken(creds); err != nil {
		return nil, time.Time{}, fmt.Errorf("could not create access token: %w", err)
	}

	if refreshToken, err = s.tokens.CreateRefreshToken(accessToken); err != nil {
		return nil, time.Time{}, fmt.Errorf("could not create refresh token: %w", err)
	}

	// Sign the tokens and return the response
	out = new(admin.AuthReply)
	if out.AccessToken, err = s.tokens.Sign(accessToken); err != nil {
		return nil, time.Time{}, fmt.Errorf("could not sign access token: %w", err)
	}
	if out.RefreshToken, err = s.tokens.Sign(refreshToken); err != nil {
		return nil, time.Time{}, fmt.Errorf("could not sign refresh token: %w", err)
	}

	// Refresh the double cookies for CSRF protection while using the access/refresh tokens
	expiresAt = refreshToken.Claims.(*tokens.Claims).ExpiresAt.Time
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
		expiresAt     time.Time
		accessClaims  *tokens.Claims
		refreshClaims *tokens.Claims
	)

	// Get the Bearer token from the Authorization header (contains access token)
	if tks, err = admin.GetAccessToken(c); err != nil {
		sentry.Warn(c).Err(err).Msg("reauthenticate called without access token")
		c.JSON(http.StatusUnauthorized, admin.ErrorResponse("request is not authorized"))
		return
	}

	// Parse the access token from the Authorization header without validating the
	// claims, e.g. it doesn't matter if the access token is expired, but it should be
	// signed correctly by the token server.
	if accessClaims, err = s.tokens.Parse(tks); err != nil {
		sentry.Warn(c).Err(err).Msg("reauthenticate called with invalid access token")
		c.JSON(http.StatusUnauthorized, admin.ErrorResponse("request is not authorized"))
		return
	}

	// Parse incoming JSON data from the client request (contains refresh token)
	in = new(admin.AuthRequest)
	if err = c.ShouldBind(&in); err != nil {
		sentry.Warn(c).Err(err).Msg("could not bind request")
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
		sentry.Warn(c).Err(err).Msg("could not verify refresh token")
		c.JSON(http.StatusUnauthorized, admin.ErrorResponse("invalid credentials"))
		return
	}

	// Ensure the refresh token and admin token match
	// TODO: verify the in.Credential is a refresh token using the subject or audience
	if accessClaims.ID != refreshClaims.ID {
		sentry.Warn(c).Msg("mismatched access and refresh token pair")
		c.JSON(http.StatusUnauthorized, admin.ErrorResponse("invalid credentials"))
		return
	}

	// At this point we've validated the reauthentication and are ready to reissue tokens
	if out, expiresAt, err = s.createAuthReply(accessClaims); err != nil {
		// NOTE: additional error logging happens in createAuthReply
		sentry.Error(c).Err(err).Msg("could not reauthenticate user")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not reauthenticate with credentials"))
		return
	}

	if err := admin.SetDoubleCookieTokens(c, s.conf.CookieDomain, expiresAt); err != nil {
		sentry.Error(c).Err(err).Msg("could not set double cookie tokens")
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

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Query the list of VASPs from the data store to perform aggregation counts.
	iter := s.db.ListVASPs(ctx)
	for iter.Next() {
		// Fetch VASP from the database
		var vasp *pb.VASP
		var err error
		if vasp, err = iter.VASP(); err != nil {
			sentry.Error(c).Err(err).Msg("could not parse VASP from database")
			continue
		}

		// Count VASPs
		out.VASPsCount++

		// Count Statuses and any status that is "pending" -- awaiting action by a reviewer.
		out.Statuses[vasp.VerificationStatus.String()]++
		if int32(vasp.VerificationStatus) < int32(pb.VerificationState_VERIFIED) || vasp.VerificationStatus == pb.VerificationState_APPEALED {
			out.PendingRegistrations++
		}

	}

	if err := iter.Error(); err != nil {
		iter.Release()
		sentry.Error(c).Err(err).Msg("could not iterate over vasps in store")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse(err))
		return
	}
	iter.Release()

	// Loop over the email contact records next
	iter2 := s.db.ListEmails(ctx)
	for iter2.Next() {
		out.ContactsCount++
		if contact, err := iter2.Email(); err == nil && contact.Verified {
			out.VerifiedContacts++
		} else if err != nil {
			sentry.Error(c).Err(err).Msg("could not parse email from database")
		}
	}

	if err := iter2.Error(); err != nil {
		iter2.Release()
		sentry.Error(c).Err(err).Msg("could not iterate over emails in store")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse(err))
		return
	}
	iter2.Release()

	// Loop over certificate requests next
	iter3 := s.db.ListCertReqs(ctx)
	for iter3.Next() {
		// Fetch CertificateRequest from the database
		var certreq *models.CertificateRequest
		var err error
		if certreq, err = iter3.CertReq(); err != nil {
			sentry.Error(c).Err(err).Msg("could not parse CertificateRequest from database")
			continue
		}

		out.CertReqs[certreq.Status.String()]++
		if certreq.Status == models.CertificateRequestState_COMPLETED {
			out.CertificatesIssued++
		}
	}

	if err := iter3.Error(); err != nil {
		iter3.Release()
		sentry.Error(c).Err(err).Msg("could not iterate over certreqs in store")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse(err))
		return
	}
	iter3.Release()

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
	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()
	iter := s.db.ListVASPs(ctx)
	defer iter.Release()
	for iter.Next() {
		// Fetch VASP from the database
		var vasp *pb.VASP
		var err error
		if vasp, err = iter.VASP(); err != nil {
			sentry.Error(c).Err(err).Msg("could not parse VASP from database")
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
					// Since this is not a unique index and multiple certs have been
					// issued to organizations in the past, we will encounter name
					// collisions here. We want this to be at debug level instead of at
					// warning level to avoid alert spam.
					log.Debug().Str("name", name).Msg("duplicate name detected")
				}
			}
		}
	}

	if err := iter.Error(); err != nil {
		sentry.Error(c).Err(err).Msg("could not iterate over vasps in store")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse(err))
		return
	}

	// In case any of the names were empty string, delete it (no guard required)
	delete(out.Names, "")

	// Successful request, return the VASP list JSON data
	c.JSON(http.StatusOK, out)
}

// ReviewTimeline returns a list of time series records containing registration state counts by week.
func (s *Admin) ReviewTimeline(c *gin.Context) {
	// Go needs this constant to determine the time format
	const timeFormat = "2006-01-02"
	var (
		err        error
		in         *admin.ReviewTimelineParams
		out        *admin.ReviewTimelineReply
		week       *utils.Week
		weekIter   *utils.WeekIterator
		startTime  time.Time
		endTime    time.Time
		vaspCounts []map[string]bool
	)

	// Get request parameters
	in = new(admin.ReviewTimelineParams)
	if err = c.ShouldBindQuery(&in); err != nil {
		sentry.Warn(c).Err(err).Msg("could not bind request with query params")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(err))
		return
	}

	if in.Start != "" {
		// Parse start date
		if startTime, err = time.Parse(timeFormat, in.Start); err != nil {
			sentry.Warn(c).Err(err).Msg("could not parse start date")
			c.JSON(http.StatusBadRequest, admin.ErrorResponse(fmt.Errorf("invalid start date: %s", in.Start)))
			return
		}
		// If the request is before the epoch then it's probably an error
		epoch := time.Unix(0, 0)
		if startTime.Before(epoch) {
			sentry.Warn(c).Err(err).Msg("start date is before epoch")
			c.JSON(http.StatusBadRequest, admin.ErrorResponse(fmt.Errorf("start date can't be before %s", epoch.Format(timeFormat))))
			return
		}
	} else {
		// Default to 1 year ago
		startTime = time.Now().AddDate(-1, 0, 0)
	}

	if in.End != "" {
		// Parse end date
		if endTime, err = time.Parse(timeFormat, in.End); err != nil {
			sentry.Warn(c).Err(err).Msg("could not parse end date")
			c.JSON(http.StatusBadRequest, admin.ErrorResponse(fmt.Errorf("invalid end date: %s", in.End)))
			return
		}
	} else {
		// Default value for end date
		endTime = time.Now()
	}

	if weekIter, err = utils.GetWeekIterator(startTime, endTime); err != nil {
		sentry.Warn(c).Err(err).Msg("invalid timeline request: start date can't be after current date")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(fmt.Errorf("start date must be before end date")))
		return
	}

	// Initialize required counting structs
	vaspCounts = make([]map[string]bool, 0, 1)
	out = &admin.ReviewTimelineReply{
		Weeks: make([]admin.ReviewTimelineRecord, 0, 1),
	}

	// Iterate over the weeks and record the week start dates
	for {
		var ok bool
		if week, ok = weekIter.Next(); !ok {
			break
		}

		record := admin.ReviewTimelineRecord{
			Week:          week.Date.Format(timeFormat),
			VASPsUpdated:  0,
			Registrations: make(map[string]int),
		}

		// Need to initialize the map entries so that all verification states show up in
		// the JSON output, even if the count is 0
		var s int32
		for s = 0; s <= int32(pb.VerificationState_ERRORED); s++ {
			record.Registrations[pb.VerificationState_name[s]] = 0
		}
		out.Weeks = append(out.Weeks, record)
		vaspCounts = append(vaspCounts, make(map[string]bool))
	}

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Iterate over the VASPs and count registrations
	iter := s.db.ListVASPs(ctx)
	defer iter.Release()
	for iter.Next() {
		// Fetch VASP from the database
		var vasp *pb.VASP
		var err error
		if vasp, err = iter.VASP(); err != nil {
			sentry.Error(c).Err(err).Msg("could not parse VASP from database")
			continue
		}

		// Get VASP audit log
		var auditLog []*models.AuditLogEntry
		if auditLog, err = models.GetAuditLog(vasp); err != nil {
			sentry.Warn(c).Err(err).Msg("could not retrieve audit log for vasp")
			continue
		}

		// Iterate over VASP audit log and count registrations
		for _, entry := range auditLog {
			var timestamp time.Time
			if timestamp, err = time.Parse(time.RFC3339, entry.Timestamp); err != nil {
				sentry.Warn(c).Err(err).Msg("could not parse timestamp in audit log entry")
				continue
			}

			// Determine which week number the recorded date falls under
			weekNum := utils.NewWeek(timestamp).Sub(weekIter.Start)
			if weekNum >= 0 && weekNum < len(out.Weeks) {
				// Count VASP if we haven't seen it before
				if _, exists := vaspCounts[weekNum][vasp.Id]; !exists {
					vaspCounts[weekNum][vasp.Id] = true
					out.Weeks[weekNum].VASPsUpdated++
				}

				// Count registration state if it changed
				if entry.PreviousState != entry.CurrentState {
					out.Weeks[weekNum].Registrations[entry.CurrentState.String()]++
				}
			}
		}
	}

	if err := iter.Error(); err != nil {
		sentry.Error(c).Err(err).Msg("could not iterate over vasps in store")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, out)
}

// ListCountries returns a list of countries with VASP registrations, sorted by
// registration count. This is an authenticated endpoint used to support visualizations
// on the Admin UI.
func (s *Admin) ListCountries(c *gin.Context) {
	var err error

	// Count registrations by country
	countries := make(map[string]int)

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	iter := s.db.ListVASPs(ctx)
	defer iter.Release()
	for iter.Next() {
		// Fetch VASP from the database
		var vasp *pb.VASP
		if vasp, err = iter.VASP(); err != nil {
			sentry.Error(c).Err(err).Str("vasp_id", vasp.Id).Msg("could not parse VASP from database")
			continue
		}

		// Only count verified VASPs
		if vasp.VerificationStatus != pb.VerificationState_VERIFIED {
			continue
		}

		// Prevent nil panics -- if this happens then the VASP is in an invalid state.
		if vasp.Entity == nil {
			sentry.Warn(c).Str("vasp_id", vasp.Id).Msg("vasp entity is nil")
			continue
		}

		// Counts are aggregated by country ISO code
		var code string
		if code = vasp.Entity.CountryOfRegistration; code == "" {
			// NOTE: Validation code does not require country of registration so
			// ignore VASP records without it but log the message for possible fixes.
			log.Debug().Str("vasp_id", vasp.Id).Msg("vasp country code is empty")
			continue
		}

		// Increment country count
		if _, ok := countries[code]; !ok {
			countries[code] = 0
		}
		countries[code]++
	}

	if err := iter.Error(); err != nil {
		sentry.Error(c).Err(err).Msg("could not iterate over vasps in store")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse(err))
		return
	}

	// Build the response with countries sorted by descending registration count
	out := make([]*admin.CountryRecord, 0)
	for code, count := range countries {
		i := sort.Search(len(out), func(i int) bool {
			return out[i].Registrations <= count
		})
		out = append(out, nil)
		copy(out[i+1:], out[i:])
		out[i] = &admin.CountryRecord{
			ISOCode:       code,
			Registrations: count,
		}
	}

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
		sentry.Warn(c).Err(err).Msg("could not bind request with query params")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(err))
		return
	}

	// Determine status filter
	filters := make(map[pb.VerificationState]struct{})
	if in.StatusFilters != nil {
		for i, s := range in.StatusFilters {
			in.StatusFilters[i] = strings.ToUpper(strings.ReplaceAll(s, " ", "_"))
			sn, ok := pb.VerificationState_value[in.StatusFilters[i]]
			if !ok {
				sentry.Warn(c).Str("status", in.StatusFilters[i]).Msg("unknown verification status")
				c.JSON(http.StatusBadRequest, admin.ErrorResponse(fmt.Errorf("unknown verification status %q", in.StatusFilters[i])))
				return
			}
			filters[pb.VerificationState(sn)] = struct{}{}
		}
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

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Query the list of VASPs from the data store
	iter := s.db.ListVASPs(ctx)
	defer iter.Release()
	for out.Count = 0; iter.Next(); out.Count++ {
		if out.Count >= minIndex && out.Count < maxIndex {
			// In the page range so add to the list reply
			// Fetch VASP from the database
			var vasp *pb.VASP
			var err error
			if vasp, err = iter.VASP(); err != nil {
				sentry.Error(c).Err(err).Msg("could not parse VASP from database")
				out.Count--
				continue
			}

			// Check against the status filters before continuing
			if _, ok := filters[vasp.VerificationStatus]; len(filters) > 0 && !ok {
				out.Count--
				continue
			}

			// Build the snippet
			snippet := admin.VASPSnippet{
				ID:                  vasp.Id,
				CommonName:          vasp.CommonName,
				RegisteredDirectory: vasp.RegisteredDirectory,
				VerificationStatus:  vasp.VerificationStatus.String(),
				LastUpdated:         vasp.LastUpdated,
				VerifiedOn:          vasp.VerifiedOn,
				Traveler:            models.IsTraveler(vasp),
			}

			// Add certificate serial number if it exists
			if vasp.IdentityCertificate != nil {
				snippet.CertificateSerial = models.GetCertID(vasp.IdentityCertificate)
				snippet.CertificateIssued = vasp.IdentityCertificate.NotBefore
				snippet.CertificateExpiration = vasp.IdentityCertificate.NotAfter
			}

			// Name is a computed value, ignore errors in finding the name.
			snippet.Name, _ = vasp.Name()

			// Add verified contacts to snippet
			var contacts *models.Contacts
			if contacts, err = s.db.VASPContacts(ctx, vasp); err != nil {
				sentry.Error(c).Err(err).Msg("could not get contact verifications")
			} else {
				snippet.VerifiedContacts = contacts.ContactVerifications()
			}

			// Append to list in reply
			out.VASPs = append(out.VASPs, snippet)
		}
	}

	if err = iter.Error(); err != nil {
		sentry.Error(c).Err(err).Msg("could not iterate over vasps in store")
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
	logctx := sentry.With(c).Str("id", vaspID)

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Attempt to fetch the VASP from the database
	if vasp, err = s.db.RetrieveVASP(ctx, vaspID); err != nil {
		logctx.Warn().Err(err).Msg("could not retrieve vasp")
		c.JSON(http.StatusNotFound, admin.ErrorResponse("could not retrieve VASP record by ID"))
		return
	}

	// Prepare VASP detail response (both retrieve and update use this method)
	// NOTE: VASP is modified in this step, must not save VASP after this!
	if out, err = s.prepareVASPDetail(ctx, vasp, logctx); err != nil {
		// NOTE: logging occurs in prepareVASPDetail
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not create VASP detail"))
		return
	}

	// Successful request, return the VASP detail JSON data
	c.JSON(http.StatusOK, out)
}

func (s *Admin) prepareVASPDetail(ctx context.Context, vasp *pb.VASP, logctx *sentry.Logger) (out *admin.RetrieveVASPReply, err error) {
	var contacts *models.Contacts
	if contacts, err = s.db.VASPContacts(ctx, vasp); err != nil {
		return nil, err
	}

	out = &admin.RetrieveVASPReply{
		Traveler:         models.IsTraveler(vasp),
		VerifiedContacts: contacts.VerifiedContacts(),
	}

	// Attempt to determine the VASP name from IVMS 101 data.
	if out.Name, err = vasp.Name(); err != nil {
		// This is a serious data validation error that needs to be addressed ASAP by
		// the operations team but should not block this API return.
		logctx.Error().Err(err).Msg("could not get VASP name")
	}

	// Add the audit log to the response, on error, create empty audit log response
	if auditLog, err := models.GetAuditLog(vasp); err != nil {
		logctx.Warn().Err(err).Msg("could not get audit log for VASP detail")
	} else {
		out.AuditLog = make([]map[string]interface{}, 0, len(auditLog))
		for i, entry := range auditLog {
			if rewiredEntry, err := wire.Rewire(entry); err != nil {
				// If we cannot rewire an audit log entry, do not serialize any audit
				// log entries to prevent confusion about what has happened in the log.
				logctx.Warn().Err(err).Int("index", i).Msg("could not rewire audit log entry for VASP detail")
				out.AuditLog = nil
				break
			} else {
				out.AuditLog = append(out.AuditLog, rewiredEntry)
			}
		}
	}

	// Add a unified email log to the response
	if emailLog, err := models.GetVASPEmailLog(vasp); err != nil {
		logctx.Warn().Err(err).Msg("could not get email log for VASP detail")
	} else {
		out.EmailLog = make([]map[string]interface{}, 0)
		for i, entry := range emailLog {
			if rewiredEntry, err := wire.Rewire(entry); err != nil {
				// If we cannot rewire an email log entry, do not serialize any email
				// log entries to prevent confusion about what has happened in the log.
				logctx.Warn().Err(err).Int("index", i).Msg("could not rewire email log entry for VASP detail")
				out.EmailLog = nil
				break
			} else {
				out.EmailLog = append(out.EmailLog, rewiredEntry)
			}
		}
	}

	// Remove extra data from the VASP
	// Must be done after verified contacts is computed
	// Note: This is safe because nothing is saved back to the database!
	vasp.Extra = nil
	for _, contact := range []*pb.Contact{vasp.Contacts.Technical, vasp.Contacts.Administrative, vasp.Contacts.Legal, vasp.Contacts.Billing} {
		if contact != nil {
			contact.Extra = nil
		}
	}

	// Rewire the VASP from protocol buffers to specific JSON serialization context
	if out.VASP, err = wire.Rewire(vasp); err != nil {
		logctx.Warn().Err(err).Msg("could rewire vasp json")
		return nil, err
	}

	// Convert the identity certificate serial number a capital hex encoded string
	if vasp.IdentityCertificate != nil {
		if _, ok := out.VASP["identity_certificate"].(map[string]interface{})["serial_number"]; !ok {
			logctx.Warn().Msg("could not parse identity certificate serial number from vasp json")
			return nil, errors.New("could not parse identity certificate serial number from vasp json")
		}
		out.VASP["identity_certificate"].(map[string]interface{})["serial_number"] = models.GetCertID(vasp.IdentityCertificate)
	}

	// Convert the signing certificate serial numbers to capital hex encoded strings
	for i, cert := range vasp.SigningCertificates {
		if _, ok := out.VASP["signing_certificates"].([]interface{})[i].(map[string]interface{})["serial_number"]; !ok {
			logctx.Warn().Msg("could not parse signing certificate serial number from vasp json")
			return nil, errors.New("could not parse signing certificate serial number from vasp json")
		}
		out.VASP["signing_certificates"].([]interface{})[i].(map[string]interface{})["serial_number"] = models.GetCertID(cert)
	}

	return out, nil
}

// UpdateVASP is a single entry point to a variety of different patches that can be made
// to the VASP object. In particular, the user may update the business details (website,
// categories, and established on), update the IVMS 101 Legal Person entity, change
// their responses to the TRIXO form, update the common name or endpoint, or manage
// contact details. Although technically, this endpoint would allow all those changes to
// be made simultaneously, the idea is that the PATCH only happens inside of those
// collections or groups of fields. Individual update methods define the logic for how
// each of those groups is updated together.
func (s *Admin) UpdateVASP(c *gin.Context) {
	var (
		err    error
		vaspID string
		vasp   *pb.VASP
		in     *admin.UpdateVASPRequest
		out    *admin.RetrieveVASPReply
		claims *tokens.Claims
	)
	// Get vaspID from the URL
	vaspID = c.Param("vaspID")

	// Parse incoming JSON data from the client request
	in = new(admin.UpdateVASPRequest)
	if err = c.ShouldBind(&in); err != nil {
		sentry.Warn(c).Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(err))
		return
	}

	// Sanity Check: Validate VASP ID
	if in.VASP != "" && in.VASP != vaspID {
		sentry.Warn(c).Str("id", in.VASP).Str("vasp_id", vaspID).Msg("mismatched request ID and URL")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse("the request ID does not match the URL endpoint"))
		return
	}

	// Create a log context for downstream logging
	logctx := sentry.With(c).Str("id", vaspID)

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Attempt to fetch the VASP from the database
	if vasp, err = s.db.RetrieveVASP(ctx, vaspID); err != nil {
		logctx.Debug().Err(err).Msg("could not retrieve vasp")
		c.JSON(http.StatusNotFound, admin.ErrorResponse("could not retrieve VASP record by ID"))
		return
	}

	// Get user claims for audit log tracing
	if claims, err = s.getClaims(c); err != nil {
		logctx.Error().Err(err).Msg("could not get user claims for audit log")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not update VASP record audit log"))
		return
	}

	// Apply changes and record if anything has changed. Note that all update methods
	// may change the VASP and must return if they created a modification requiring the
	// VASP to be saved back to the database.
	var (
		code     int
		updated  bool
		nChanges uint8
	)

	// Update business information
	if updated, code, err = s.updateVASPBusinessInfo(vasp, in, logctx); err != nil {
		// NOTE: logging happens in the update helper function
		c.JSON(code, admin.ErrorResponse(err))
		return
	} else if updated {
		// Log all of the updates that were made in one log message.
		logctx = logctx.Str("business_info", "VASP business information updated")
		nChanges++
	}

	// Update VASP entity information
	if updated, code, err = s.updateVASPEntity(vasp, in.Entity, logctx); err != nil {
		// NOTE: logging happens in the update helper function
		c.JSON(code, admin.ErrorResponse(err))
		return
	} else if updated {
		// Log all of the updates that were made in one log message.
		logctx = logctx.Str("vasp_entity", "VASP IVMS101 record updated")
		nChanges++
	}

	// Update VASP TRIXO form
	if updated, code, err = s.updateVASPTRIXO(vasp, in.TRIXO, logctx); err != nil {
		// NOTE: logging happens in the update helper function
		c.JSON(code, admin.ErrorResponse(err))
		return
	} else if updated {
		// Log all of the updates that were made in one log message.
		logctx = logctx.Str("vasp_trixo", "VASP TRIXO form updated")
		nChanges++
	}

	// Update common name and trisa endpoint - this will also update any certificate requests.
	// NOTE: if updated is true and any failure occurs after this point, the certificate requests
	// will be in an inconsistent state. Transactions would be very nice here, but instead this
	// function is performed as late as possible to minimize the chance of other errors.
	// The log messages that follow this line of code may be at a higher level than we might
	// expect so that we can debug the case where the database moves to an inconsistent state.
	if updated, code, err = s.updateVASPEndpoint(vasp, in.CommonName, in.TRISAEndpoint, claims.Email, logctx); err != nil {
		// NOTE: logging happens in the update helper function
		c.JSON(code, admin.ErrorResponse(err))
		return
	} else if updated {
		// Log all of the updates that were made in one log message.
		logctx = logctx.Str("trisa_endpoint", "trisa endpoint and common name updated")
		nChanges++
	}

	// Check if we've updated anything, if not return an error to indicate to the front
	// end that no work was performed.
	if nChanges == 0 {
		logctx.Debug().Msg("no updates on VASP occurred")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse("no updates made to VASP record"))
		return
	}

	// Validate that the VASP record is still correct after the changes.
	// Note: the updateVASPEntity and updateVASPEndpoint both make similar checks to
	// ensure that certificate requests are not saved when the VASP record is not valid.
	if err = models.ValidateVASP(vasp, true); err != nil {
		sentry.Warn(c).Err(err).Msg("invalid or incomplete VASP record on update")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(fmt.Errorf("validation error: %s", err)))
		return
	}

	// Add a record to the audit log
	// NOTE: if validation, adding a record to the audit log, or saving the VASP back to
	// disk errors, we could end up in an inconsistent state where we have changes to
	// the certificate request that do not appear in the VASP audit log.
	if err = models.UpdateVerificationStatus(vasp, vasp.VerificationStatus, "VASP record updated by admin", claims.Email); err != nil {
		logctx.Error().Err(err).Msg("could not add audit log entry by updating the verification status")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not update VASP audit log"))
		return
	}

	// Since updates have occurred, save the changes
	// TODO: transactions would be super nice here so we could rollback any certificate request changes
	if err = s.db.UpdateVASP(ctx, vasp); err != nil {
		logctx.Error().Err(err).Msg("could not save VASP after update")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not update VASP"))
		return
	}

	// Create the response to send back, ensuring extra fields are removed.
	// Prepare VASP detail response (both retrieve and update use this method)
	// NOTE: VASP is modified in this step, must not save VASP after this!
	if out, err = s.prepareVASPDetail(ctx, vasp, logctx); err != nil {
		// NOTE: logging occurs in prepareVASPDetail
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not create VASP detail"))
		return
	}

	// Successful request, return the VASP detail JSON data
	logctx.Info().Uint8("n_changes", nChanges).Msg("update VASP completed")
	c.JSON(http.StatusOK, admin.UpdateVASPReply(*out))
}

// Update the VASP business information such as website, business and VASP categories,
// and established on date. This information can be modified at anytime but cannot be
// set to an empty value, otherwise the PATCH update will not take place.
func (s *Admin) updateVASPBusinessInfo(vasp *pb.VASP, in *admin.UpdateVASPRequest, logctx *sentry.Logger) (updated bool, _ int, err error) {
	if in.Website != "" {
		vasp.Website = in.Website
		updated = true
	}

	if in.BusinessCategory != "" {
		category, ok := pb.BusinessCategory_value[in.BusinessCategory]
		if !ok {
			return false, http.StatusBadRequest, errors.New("could not parse business category")
		}

		vasp.BusinessCategory = pb.BusinessCategory(category)
		updated = true
	}

	if len(in.VASPCategories) > 0 {
		vasp.VaspCategories = in.VASPCategories
		updated = true
	}

	if in.EstablishedOn != "" {
		vasp.EstablishedOn = in.EstablishedOn
		updated = true
	}

	return updated, http.StatusOK, nil
}

// Update the VASP IVMS101 Legal Person entity; the LegalPerson entity must be valid.
// This method completely overwrites the previous LegalPerson entity, no field-level
// patching is available.
func (s *Admin) updateVASPEntity(vasp *pb.VASP, data map[string]interface{}, logctx *sentry.Logger) (_ bool, _ int, err error) {
	// Check if entity data has been supplied, otherwise do not update.
	if len(data) == 0 {
		return false, http.StatusOK, nil
	}

	// Remarshal the JSON IVMS 101 entity
	entity := &ivms101.LegalPerson{}
	if err = wire.Unwire(data, entity); err != nil {
		logctx.Warn().Err(err).Msg("could not unwire JSON data into an IVMS 101 LegalPerson")
		return false, http.StatusBadRequest, errors.New("could not parse IVMS 101 LegalPerson entity")
	}

	// Validation here is an extra guard, even though validate is also called in the
	// primary RPC function. This is to ensure that an invalid VASP doesn't have
	// certificate requests updated inappropriately.
	// NOTE: other methods ignore ErrCompleteNationalIdentifierLegalPerson, but it is not
	// ignored here, requiring the admin to determine how best to accurately update the entity.
	if err = entity.Validate(); err != nil {
		logctx.Debug().Err(err).Msg("invalid IVMS 101 LegalPerson struct")
		return false, http.StatusBadRequest, err
	}

	vasp.Entity = entity
	return true, http.StatusOK, nil
}

// Update the VASP TRIXO form; the TRIXO form really has no internal validation.
// This method completely overwrites the previous LegalPerson entity, no field-level
// patching is available.
func (s *Admin) updateVASPTRIXO(vasp *pb.VASP, data map[string]interface{}, logctx *sentry.Logger) (_ bool, _ int, err error) {
	// Check if trixo data has been supplied, otherwise do not update.
	if len(data) == 0 {
		return false, http.StatusOK, nil
	}

	// Remarshal the JSON TRIXO questionnaire
	trixo := &pb.TRIXOQuestionnaire{}
	if err = wire.Unwire(data, trixo); err != nil {
		logctx.Warn().Err(err).Msg("could not unwire JSON data into an valid TRIXO Questionnaire")
		return false, http.StatusBadRequest, errors.New("could not parse TRIXO questionnaire")
	}

	vasp.Trixo = trixo
	return true, http.StatusOK, nil
}

func (s *Admin) updateVASPEndpoint(vasp *pb.VASP, commonName, endpoint, source string, logctx *sentry.Logger) (_ bool, _ int, err error) {
	if commonName == "" && endpoint == "" {
		return false, http.StatusOK, nil
	}

	// Compute the common name from the TRISA endpoint if not specified
	if commonName == "" && endpoint != "" {
		if commonName, _, err = net.SplitHostPort(endpoint); err != nil {
			logctx.Warn().Err(err).Str("endpoint", endpoint).Msg("could not parse common name from endpoint")
			return false, http.StatusBadRequest, errors.New("no common name supplied, could not parse common name from endpoint")
		}
	}

	// Check if any changes are required
	if vasp.CommonName == commonName && vasp.TrisaEndpoint == endpoint {
		return false, http.StatusOK, nil
	}

	// Check if this is just an endpoint change
	if vasp.CommonName == commonName {
		vasp.TrisaEndpoint = endpoint
		logctx.Info().Msg("trisa endpoint updated without change to common name")
		return true, http.StatusOK, nil
	}

	if vasp.VerificationStatus >= pb.VerificationState_REVIEWED {
		// Cannot change common name after certificates have been issued
		logctx.Warn().Str("status", vasp.VerificationStatus.String()).Str("common_name", commonName).Msg("could not update VASP common name")
		return false, http.StatusBadRequest, errors.New("cannot update common name in current state")
	}

	// Make changes to both the VASP and the CertificateRequest
	vasp.CommonName = commonName
	vasp.TrisaEndpoint = endpoint

	// Get the Certificate Request IDs from the VASP model
	var certreqs []string
	if certreqs, err = models.GetCertReqIDs(vasp); err != nil {
		logctx.Error().Err(err).Msg("could not get certificate requests for VASP")
		return false, http.StatusInternalServerError, errors.New("could not update certificate request with common name")
	}

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Loop through all of the certificate requests and check if they can be updated
	ncertreqs := 0
	for _, certreqID := range certreqs {
		var certreq *models.CertificateRequest
		if certreq, err = s.db.RetrieveCertReq(ctx, certreqID); err != nil {
			logctx.Error().Err(err).Str("certreq_id", certreqID).Msg("could not fetch certificate request for VASP")
			return false, http.StatusInternalServerError, errors.New("could not update certificate request with common name")
		}

		// If the certificate request has already been submitted, we cannot change its common name
		if certreq.Status > models.CertificateRequestState_READY_TO_SUBMIT {
			logctx.Debug().Str("status", certreq.Status.String()).Str("certreq_id", certreqID).Msg("could not update certificate request")
			continue
		}

		// Update certificate request and add an audit log entry
		certreq.CommonName = commonName
		if err = models.UpdateCertificateRequestStatus(certreq, certreq.Status, "common name changed", source); err != nil {
			logctx.Error().Err(err).Str("certreq_id", certreqID).Msg("could not update certificate request status to add audit log entry")
			continue
		}

		// Store the certificate request back to disk
		if err = s.db.UpdateCertReq(ctx, certreq); err != nil {
			logctx.Error().Err(err).Str("certreq_id", certreqID).Msg("could not update certificate request for VASP")
			continue
		}

		logctx.Info().Str("certreq_id", certreqID).Msg("certificate request updated")
		ncertreqs++
	}

	if ncertreqs == 0 {
		logctx.Error().Msg("no certificate requests updated with common name")
		return false, http.StatusInternalServerError, errors.New("could not update certificate request with common name")
	}

	// NOTE: from this point on it's possible that we have an unsaved VASP that has had
	// modifications to its certificate requests. If the VASP is not saved after this
	// method, it could lead to an inconsistency that needs to be repaired manually.
	return true, http.StatusOK, nil
}

// DeleteVASP removes a VASP and its associated certificate requests if and only if the
// VASP verification status is in PENDING_REVIEW or earlier or ERRORED.
func (s *Admin) DeleteVASP(c *gin.Context) {
	var (
		vaspID     string
		vasp       *pb.VASP
		certReqIDs []string
		err        error
	)

	vaspID = c.Param("vaspID")

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Retrieve the VASP from the database
	if vasp, err = s.db.RetrieveVASP(ctx, vaspID); err != nil {
		sentry.Warn(c).Err(err).Msg("could not retrieve VASP from database")
		c.JSON(http.StatusNotFound, admin.ErrorResponse("could not retrieve VASP record by ID"))
		return
	}

	// Only allow deletions if the VASP has not been reviewed yet
	if vasp.VerificationStatus > pb.VerificationState_PENDING_REVIEW && vasp.VerificationStatus < pb.VerificationState_ERRORED {
		sentry.Warn(c).Str("status", vasp.VerificationStatus.String()).Msg("VASP is in invalid state for deletion")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse("cannot delete VASP in its current state"))
		return
	}

	// Retrieve the associated certificate requests
	if certReqIDs, err = models.GetCertReqIDs(vasp); err != nil {
		sentry.Error(c).Err(err).Msg("could not retrieve certificate request IDs for VASP")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not retrieve certificate requests for VASP"))
		return
	}

	// Delete the certificate requests
	for _, id := range certReqIDs {
		if err = s.db.DeleteCertReq(ctx, id); err != nil {
			sentry.Error(c).Err(err).Str("certreq_id", id).Msg("could not delete certificate request")
			c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not delete associated certificate request"))
			return
		}
	}

	// Delete the VASP object to finalize the VASP deletion
	if err = s.db.DeleteVASP(ctx, vaspID); err != nil {
		sentry.Error(c).Err(err).Msg("could not delete VASP from database")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not delete VASP record by ID"))
		return
	}

	c.JSON(http.StatusOK, admin.Reply{Success: true})
}

// ListCertificates returns a list of certificates for the VASP.
func (s *Admin) ListCertificates(c *gin.Context) {
	var err error

	// Get vaspID from the URL
	vaspID := c.Param("vaspID")

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Retrieve the VASP from the database
	var vasp *pb.VASP
	if vasp, err = s.db.RetrieveVASP(ctx, vaspID); err != nil {
		sentry.Warn(c).Err(err).Str("vasp_id", vaspID).Msg("could not retrieve VASP from database")
		c.JSON(http.StatusNotFound, admin.ErrorResponse("could not retrieve VASP record by ID"))
		return
	}

	// Retrieve the Certificate IDs from the VASP
	var ids []string
	if ids, err = models.GetCertIDs(vasp); err != nil {
		sentry.Error(c).Err(err).Str("vasp_id", vaspID).Msg("could not retrieve certificate IDs for VASP")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not retrieve certificate IDs for VASP"))
		return
	}

	// Construct the reply
	out := &admin.ListCertificatesReply{
		Certificates: make([]admin.Certificate, 0),
	}

	for _, id := range ids {
		// Retrieve the Certificate from the database
		var cert *models.Certificate
		if cert, err = s.db.RetrieveCert(ctx, id); err != nil {
			sentry.Error(c).Err(err).Str("cert_id", id).Msg("could not retrieve certificate from database")
			c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not retrieve certificate by ID"))
			return
		}

		// Construct an entry for the reply
		entry := admin.Certificate{
			SerialNumber: id,
			IssuedAt:     cert.Details.NotBefore,
			ExpiresAt:    cert.Details.NotAfter,
			Status:       cert.Status.String(),
		}
		if entry.Details, err = wire.Rewire(cert.Details); err != nil {
			sentry.Error(c).Err(err).Str("cert_id", id).Msg("could not serialize certificate details")
			c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not serialize certificate details"))
			return
		}

		out.Certificates = append(out.Certificates, entry)
	}

	c.JSON(http.StatusOK, out)
}

// ReplaceContact completely replaces a contact on a VASP with a new contact.
func (s *Admin) ReplaceContact(c *gin.Context) {
	var (
		in           *admin.ReplaceContactRequest
		contact      *models.ContactRecord
		vasp         *pb.VASP
		email        *models.Email
		emailUpdated bool
		err          error
	)

	// Get vaspID from the URL
	vaspID := c.Param("vaspID")
	kind := c.Param("kind")

	// Parse incoming JSON data from the client request
	in = new(admin.ReplaceContactRequest)
	if err = c.ShouldBind(&in); err != nil {
		sentry.Warn(c).Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(err))
		return
	}

	// Sanity check: validate VASP ID
	if in.VASP != "" && in.VASP != vaspID {
		sentry.Warn(c).Str("id", in.VASP).Str("vasp_id", vaspID).Msg("mismatched request ID and URL")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse("the request ID does not match the URL endpoint"))
		return
	}

	// Sanity check: validate contact kind
	if in.Kind != "" && in.Kind != kind {
		sentry.Warn(c).Str("kind", in.Kind).Str("kind", kind).Msg("mismatched contact kind and URL")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse("the contact kind does not match the URL endpoint"))
		return
	}

	// Kind must be one of the accepted values
	if !models.ContactKindIsValid(kind) {
		sentry.Warn(c).Str("kind", kind).Msg("invalid contact kind")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse("invalid contact kind provided"))
		return
	}

	// Contact data must be provided
	if len(in.Contact) == 0 {
		sentry.Warn(c).Msg("missing contact data on ReplaceContact request")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse("contact data is required for ReplaceContact request"))
		return
	}

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Retrieve the VASP from the database
	if vasp, err = s.db.RetrieveVASP(ctx, vaspID); err != nil {
		sentry.Warn(c).Err(err).Msg("could not retrieve VASP from database")
		c.JSON(http.StatusNotFound, admin.ErrorResponse("could not retrieve VASP record by ID"))
		return
	}

	// Retrieve the VASP contacts from the database
	var contacts *models.Contacts
	if contacts, err = s.db.VASPContacts(ctx, vasp); err != nil {
		sentry.Error(c).Err(err).Msg("could not retrieve contacts for VASP")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not retrieve VASP contacts"))
		return
	}

	// Remarshal the JSON contact data
	update := &pb.Contact{}
	if err = wire.Unwire(in.Contact, update); err != nil {
		sentry.Warn(c).Err(err).Msg("could not unmarshal contact data")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(err))
		return
	}

	// Verify that updated contact has data
	if update.IsZero() {
		c.JSON(http.StatusBadRequest, admin.ErrorResponse("invalid contact data: missing required fields"))
		return
	}

	if contact = contacts.Get(kind); contact == nil {
		// If the contact does not exist on the VASP then create it
		email, emailUpdated = contacts.Add(kind, update)
	} else {
		// Otherwise replace the existing contact info
		email, emailUpdated = contact.Update(update)
	}

	// New VASP record must be valid
	if err = models.ValidateVASP(vasp, true); err != nil {
		sentry.Warn(c).Err(err).Msg("invalid VASP record after update")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(fmt.Errorf("validation error: %s", err)))
		return
	}

	if emailUpdated {
		// Attempt to retrieve the email from the database if it exists
		var emailRecord *models.Email
		if emailRecord, err = s.db.RetrieveEmail(ctx, email.Email); err != nil {
			// If this is not a not found error then something else went wrong
			if !errors.Is(err, storeerrors.ErrEntityNotFound) {
				sentry.Error(c).Err(err).Msg("could not retrieve email from database")
				c.JSON(http.StatusInternalServerError, admin.ErrorResponse(err))
				return
			}

			// If the email record does not exist in the database create it
			email.Token = secrets.CreateToken(models.VerificationTokenLength)
			email.Verified = false
			if _, err = s.db.CreateEmail(ctx, email); err != nil {
				sentry.Error(c).Err(err).Msg("could not create email")
				c.JSON(http.StatusInternalServerError, admin.ErrorResponse(err))
				return
			}

			// Send the email verification request
			if err = s.svc.email.SendVerifyContact(vasp.Id, contact); err != nil {
				sentry.Error(c).Err(err).Msg("could not send verification email")
				c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not send verification email to the new contact"))
				return
			}

			// Update the email record with the send log after sending contact verification.
			if err = s.db.UpdateEmail(ctx, emailRecord); err != nil {
				sentry.Error(c).Err(err).Msg("could not save updated email")
			}
		} else {
			// Update the email record with the vaspID and new name if it exists.
			if email.Name != "" {
				emailRecord.Name = email.Name
			}
			emailRecord.AddVASP(vasp.Id)
			if err = s.db.UpdateEmail(ctx, emailRecord); err != nil {
				sentry.Error(c).Err(err).Msg("could not save updated email")
				c.JSON(http.StatusInternalServerError, admin.ErrorResponse(err))
				return
			}
		}
	}

	// Commit the contact changes to the database
	if err = s.db.UpdateVASP(ctx, vasp); err != nil {
		sentry.Error(c).Err(err).Msg("could not update VASP in database")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not update VASP record by ID"))
		return
	}

	c.JSON(http.StatusOK, admin.Reply{Success: true})
}

// DeleteContact deletes a contact on a VASP.
func (s *Admin) DeleteContact(c *gin.Context) {
	var (
		vasp *pb.VASP
		err  error
	)

	// Get vaspID from the URL
	vaspID := c.Param("vaspID")
	kind := c.Param("kind")

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Retrieve the VASP from the database
	if vasp, err = s.db.RetrieveVASP(ctx, vaspID); err != nil {
		sentry.Warn(c).Err(err).Msg("could not retrieve VASP from database")
		c.JSON(http.StatusNotFound, admin.ErrorResponse("could not retrieve VASP record by ID"))
		return
	}

	// Kind must be one of the accepted values
	if !models.ContactKindIsValid(kind) {
		c.JSON(http.StatusBadRequest, admin.ErrorResponse("invalid contact kind provided"))
		return
	}

	var contacts *models.Contacts
	if contacts, err = s.db.VASPContacts(ctx, vasp); err != nil {
		sentry.Error(c).Err(err).Msg("could not retrieve vasp contacts")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse(err))
		return
	}

	if email, emailUpdated := contacts.Delete(kind); emailUpdated {
		if err = s.db.UpdateEmail(ctx, email); err != nil {
			sentry.Error(c).Err(err).Msg("could not update contact to remove vaspID")
		}
	}

	// New VASP record must be valid
	if err = models.ValidateVASP(vasp, true); err != nil {
		sentry.Warn(c).Err(err).Msg("invalid VASP record after update")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(fmt.Errorf("validation error: %s", err)))
		return
	}

	// Commit the contact changes to the database
	if err = s.db.UpdateVASP(ctx, vasp); err != nil {
		sentry.Error(c).Err(err).Msg("could not update VASP in database")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not update VASP record by ID"))
		return
	}

	c.JSON(http.StatusOK, admin.Reply{Success: true})
}

// CreateReviewNote creates a new review note given the vaspID param and a CreateReviewNoteRequest.
func (s *Admin) CreateReviewNote(c *gin.Context) {
	var (
		err    error
		in     *admin.ModifyReviewNoteRequest
		note   *models.ReviewNote
		vasp   *pb.VASP
		claims *tokens.Claims
		vaspID string
		noteID string
	)

	// Get vaspID from the URL
	vaspID = c.Param("vaspID")

	// Parse incoming JSON data from the client request
	in = new(admin.ModifyReviewNoteRequest)
	if err = c.ShouldBind(&in); err != nil {
		sentry.Warn(c).Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(err))
		return
	}

	// Validate VASP ID
	if in.VASP != "" && in.VASP != vaspID {
		sentry.Warn(c).Str("id", in.VASP).Str("vasp_id", vaspID).Msg("mismatched request ID and URL")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse("the request ID does not match the URL endpoint"))
		return
	}

	// Retrieve author email
	if claims, err = s.getClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not retrieve user claims")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("unable to retrieve user info"))
		return
	}

	if in.NoteID == "" {
		// Create note ID if not provided
		noteID = uuid.New().String()
	} else {
		noteID = in.NoteID
		// Only allow reasonable length note IDs (generated IDs are also 36 characters)
		if len(noteID) > 36 {
			sentry.Warn(c).Err(err).Msg("invalid note ID")
			c.JSON(http.StatusBadRequest, admin.ErrorResponse("note ID cannot be longer than 36 characters"))
			return
		}

		// Only allow note IDs that can be used in request URLs
		if escaped := url.QueryEscape(noteID); noteID != escaped {
			sentry.Warn(c).Err(err).Msg("invalid note ID")
			c.JSON(http.StatusBadRequest, admin.ErrorResponse(fmt.Errorf("note ID contains unescaped characters: %s", noteID)))
			return
		}
	}

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Lookup the VASP record associated with the request
	if vasp, err = s.db.RetrieveVASP(ctx, vaspID); err != nil {
		sentry.Warn(c).Err(err).Str("id", vaspID).Msg("could not retrieve vasp")
		c.JSON(http.StatusNotFound, admin.ErrorResponse("could not retrieve VASP record by ID"))
		return
	}

	// Create the note
	if note, err = models.CreateReviewNote(vasp, noteID, claims.Email, in.Text); err != nil {
		sentry.Warn(c).Err(err).Msg("error creating review note")
		if err == models.ErrorAlreadyExists {
			c.JSON(http.StatusBadRequest, admin.ErrorResponse("note already exists"))
		} else {
			c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not create review note"))
		}
		return
	}

	// Persist the VASP record to the database
	if err = s.db.UpdateVASP(ctx, vasp); err != nil {
		sentry.Error(c).Err(err).Msg("error updating VASP record")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not update VASP record"))
		return
	}

	c.JSON(http.StatusCreated, &admin.ReviewNote{
		ID:       note.Id,
		Created:  note.Created,
		Modified: note.Modified,
		Author:   note.Author,
		Editor:   note.Editor,
		Text:     note.Text,
	})
}

// ListReviewNotes returns a list of review notes given a vaspID param.
func (s *Admin) ListReviewNotes(c *gin.Context) {
	var (
		err    error
		out    *admin.ListReviewNotesReply
		vasp   *pb.VASP
		vaspID string
		notes  map[string]*models.ReviewNote
	)

	// Get vaspID from the URL
	vaspID = c.Param("vaspID")

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Lookup the VASP record associated with the request
	if vasp, err = s.db.RetrieveVASP(ctx, vaspID); err != nil {
		sentry.Warn(c).Err(err).Str("id", vaspID).Msg("could not retrieve vasp")
		c.JSON(http.StatusNotFound, admin.ErrorResponse("could not retrieve VASP record by ID"))
		return
	}

	// Retrieve the slice of notes
	if notes, err = models.GetReviewNotes(vasp); err != nil {
		sentry.Error(c).Err(err).Msg("error retrieving review notes")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not retrieve review notes"))
		return
	}

	// Compose the JSON response
	out = &admin.ListReviewNotesReply{
		Notes: []admin.ReviewNote{},
	}
	for _, n := range notes {
		out.Notes = append(out.Notes, admin.ReviewNote{
			ID:       n.Id,
			Created:  n.Created,
			Modified: n.Modified,
			Author:   n.Author,
			Editor:   n.Editor,
			Text:     n.Text,
		})
	}

	c.JSON(http.StatusOK, out)
}

// UpdateReviewNote updates the text of a review note given vaspIP and noteID params
// and an UpdateReviewNoteRequest.
func (s *Admin) UpdateReviewNote(c *gin.Context) {
	var (
		err    error
		in     *admin.ModifyReviewNoteRequest
		note   *models.ReviewNote
		vasp   *pb.VASP
		claims *tokens.Claims
		vaspID string
		noteID string
	)

	// Get vaspID and noteID from the URL
	vaspID = c.Param("vaspID")
	noteID = c.Param("noteID")

	// Parse incoming JSON data from the client request
	in = new(admin.ModifyReviewNoteRequest)
	if err = c.ShouldBind(&in); err != nil {
		sentry.Warn(c).Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(err))
		return
	}

	// Validate VASP ID
	if in.VASP != "" && in.VASP != vaspID {
		sentry.Warn(c).Str("id", in.VASP).Str("vasp_id", vaspID).Msg("mismatched request ID and URL")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse("the request VASP ID does not match the URL endpoint"))
		return
	}

	// Validate note ID
	if in.NoteID != "" && in.NoteID != noteID {
		sentry.Warn(c).Str("id", in.NoteID).Str("note_id", noteID).Msg("mismatched request ID and URL")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse("the request Note ID does not match the URL endpoint"))
		return
	}

	// Retrieve author email
	if claims, err = s.getClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not retrieve user claims")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("unable to retrieve user info"))
		return
	}

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Lookup the VASP record associated with the request
	if vasp, err = s.db.RetrieveVASP(ctx, vaspID); err != nil {
		sentry.Warn(c).Err(err).Str("id", vaspID).Msg("could not retrieve vasp")
		c.JSON(http.StatusNotFound, admin.ErrorResponse("could not retrieve VASP record by ID"))
		return
	}

	// Update the note
	if note, err = models.UpdateReviewNote(vasp, noteID, claims.Email, in.Text); err != nil {
		sentry.Error(c).Err(err).Msg("error updating review note")
		if err == models.ErrorNotFound {
			c.JSON(http.StatusNotFound, admin.ErrorResponse("review note not found"))
		} else {
			c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not update review note"))
		}
		return
	}

	// Persist the VASP record to the database
	if err = s.db.UpdateVASP(ctx, vasp); err != nil {
		sentry.Error(c).Err(err).Msg("error updating VASP record")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not update VASP record"))
		return
	}

	c.JSON(http.StatusOK, &admin.ReviewNote{
		ID:       note.Id,
		Created:  note.Created,
		Modified: note.Modified,
		Author:   note.Author,
		Editor:   note.Editor,
		Text:     note.Text,
	})
}

// DeleteReviewNote deletes a review note given vaspID and noteID params.
func (s *Admin) DeleteReviewNote(c *gin.Context) {
	var (
		err    error
		vasp   *pb.VASP
		vaspID string
		noteID string
	)

	// Get vaspID and noteID from the URL
	vaspID = c.Param("vaspID")
	noteID = c.Param("noteID")

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Lookup the VASP record associated with the request
	if vasp, err = s.db.RetrieveVASP(ctx, vaspID); err != nil {
		sentry.Warn(c).Err(err).Str("id", vaspID).Msg("could not retrieve vasp")
		c.JSON(http.StatusNotFound, admin.ErrorResponse("could not retrieve VASP record by ID"))
		return
	}

	// Delete the note
	if err = models.DeleteReviewNote(vasp, noteID); err != nil {
		sentry.Warn(c).Err(err).Msg("error deleting review note")
		if err == models.ErrorNotFound {
			c.JSON(http.StatusNotFound, admin.ErrorResponse("review note not found"))
		} else {
			c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not delete review note"))
		}
		return
	}

	// Persist the VASP record to the database
	if err = s.db.UpdateVASP(ctx, vasp); err != nil {
		sentry.Error(c).Err(err).Msg("error updating VASP record")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not update VASP record"))
		return
	}

	c.JSON(http.StatusOK, &admin.Reply{Success: true})
}

// ReviewToken returns the admin verification token of the VASP if the VASP is in a
// state that it can be reviewed in, e.g. PENDING_REVIEW, otherwise a 404 is returned.
func (s *Admin) ReviewToken(c *gin.Context) {
	var (
		err    error
		out    *admin.ReviewTokenReply
		vasp   *pb.VASP
		vaspID string
	)

	// Get vaspID from the URL
	vaspID = c.Param("vaspID")

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Lookup the VASP record associated with the request
	if vasp, err = s.db.RetrieveVASP(ctx, vaspID); err != nil {
		sentry.Warn(c).Err(err).Str("id", vaspID).Msg("could not retrieve vasp")
		c.JSON(http.StatusNotFound, admin.ErrorResponse("could not retrieve VASP record by ID"))
		return
	}

	// Check if the VASP is in a state where it can be reviewed
	if vasp.VerificationStatus != pb.VerificationState_PENDING_REVIEW {
		log.Debug().Str("id", vaspID).Str("status", vasp.VerificationStatus.String()).Msg("could not retrieve admin verification token in current state")
		c.JSON(http.StatusNotFound, admin.ErrorResponse("admin verification token not available if VASP is not pending review"))
		return
	}

	// Construct the reply
	out = &admin.ReviewTokenReply{}
	if out.AdminVerificationToken, err = models.GetAdminVerificationToken(vasp); err != nil {
		sentry.Error(c).Err(err).Str("id", vaspID).Msg("could not retrieve admin verification token")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not retrieve admin verification token"))
		return
	}

	// Check that an admin verification token will be returned
	if out.AdminVerificationToken == "" {
		sentry.Error(c).Str("id", vaspID).Str("status", vasp.VerificationStatus.String()).Msg("admin verification token not available to review VASP")
		c.JSON(http.StatusNotFound, admin.ErrorResponse("could not retrieve admin verification token"))
		return
	}

	// Return the request with the admin verification token
	c.JSON(http.StatusOK, out)
}

// Review a registration request and either accept or reject it. On accept, the
// certificate request that was created on verify is used to send a Sectigo request and
// the certificate manager process watches it until the certificate has been issued. On
// reject, the VASP and certificate request records are deleted and the reject reason is
// sent to the technical contact.
func (s *Admin) Review(c *gin.Context) {
	var (
		err      error
		in       *admin.ReviewRequest
		out      *admin.ReviewReply
		vasp     *pb.VASP
		vaspID   string
		contacts *models.Contacts
	)

	// Get vaspID from the URL
	vaspID = c.Param("vaspID")

	// Parse incoming JSON data from the client request
	in = new(admin.ReviewRequest)
	if err := c.ShouldBind(&in); err != nil {
		sentry.Warn(c).Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(err))
		return
	}

	// Validate review request
	if in.ID != "" && in.ID != vaspID {
		sentry.Warn(c).Msg("mismatched request ID and URL")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse("the request ID does not match the URL endpoint"))
		return
	}

	if in.AdminVerificationToken == "" {
		sentry.Warn(c).Msg("no verification token specified")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse("the admin verification token is required"))
		return
	}

	if !in.Accept && in.RejectReason == "" {
		sentry.Warn(c).Msg("missing reject reason")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse("if rejecting the request, a reason must be supplied"))
		return
	}

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Lookup the VASP record associated with the request
	if vasp, err = s.db.RetrieveVASP(ctx, vaspID); err != nil {
		sentry.Warn(c).Err(err).Str("id", vaspID).Msg("could not retrieve vasp")
		c.JSON(http.StatusNotFound, admin.ErrorResponse("could not retrieve VASP record by ID"))
		return
	}

	// Lookup the contacts records associated with the VASP
	if contacts, err = s.db.VASPContacts(ctx, vasp); err != nil {
		sentry.Warn(c).Err(err).Str("id", vaspID).Msg("could not retrieve vasp contacts")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not retrieve vasp contacts"))
		return
	}

	// Check that the administration verification token is correct
	var adminVerificationToken string
	if adminVerificationToken, err = models.GetAdminVerificationToken(vasp); err != nil {
		sentry.Error(c).Err(err).Str("id", vaspID).Msg("could not retrieve admin verification token")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not retrieve admin verification token"))
		return
	}
	if in.AdminVerificationToken != adminVerificationToken {
		sentry.Warn(c).Err(err).Str("vasp", vaspID).Msg("incorrect admin verification token")
		c.JSON(http.StatusUnauthorized, admin.ErrorResponse("admin verification token not accepted"))
		return
	}

	// Retrieve user claims for access to provided user info
	var claims *tokens.Claims
	if claims, err = s.getClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not retrieve user claims")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("unable to retrieve user info"))
		return
	}

	// Accept or reject the request
	out = &admin.ReviewReply{}
	logctx := sentry.With(c).Str("vaspID", vasp.Id)

	if in.Accept {
		if out.Message, err = s.acceptRegistration(vasp, claims, logctx); err != nil {
			sentry.Error(c).Err(err).Msg("could not accept VASP registration")
			c.JSON(http.StatusInternalServerError, admin.ErrorResponse("unable to accept VASP registration request"))
			return
		}
	} else {
		if out.Message, err = s.rejectRegistration(vasp, contacts, in.RejectReason, claims, logctx); err != nil {
			sentry.Error(c).Err(err).Msg("could not reject VASP registration")
			c.JSON(http.StatusInternalServerError, admin.ErrorResponse("unable to reject VASP registration request"))
			return
		}
	}

	// Persist the VASP record to the database
	if err = s.db.UpdateVASP(ctx, vasp); err != nil {
		sentry.Error(c).Err(err).Msg("error updating VASP record")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not update VASP record"))
		return
	}

	name, _ := vasp.Name()
	out.Status = vasp.VerificationStatus.String()
	log.Info().Str("vasp", vasp.Id).Str("name", name).Bool("accepted", in.Accept).Msg("registration reviewed")
	c.JSON(http.StatusOK, out)
}

// Accept the VASP registration and begin the certificate issuance process.
func (s *Admin) acceptRegistration(vasp *pb.VASP, claims *tokens.Claims, logctx *sentry.Logger) (msg string, err error) {
	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Change the VASP verification status
	if err = models.SetAdminVerificationToken(vasp, ""); err != nil {
		return "", err
	}
	vasp.VerifiedOn = time.Now().Format(time.RFC3339)
	if err := models.UpdateVerificationStatus(vasp, pb.VerificationState_REVIEWED, "registration request received", claims.Email); err != nil {
		return "", err
	}
	if err = s.db.UpdateVASP(ctx, vasp); err != nil {
		return "", err
	}

	// Mark any initialized certificate requests for this VASP as ready to submit
	// NOTE: there should only be one certificate request per VASP, but no errors occur
	// if there are more than one (other than a logged warning).
	var (
		ncertreqs int
		careqs    []string
	)

	if careqs, err = models.GetCertReqIDs(vasp); err != nil {
		return "", err
	}

	for _, careqID := range careqs {
		var careq *models.CertificateRequest
		if careq, err = s.db.RetrieveCertReq(ctx, careqID); err != nil {
			logctx.Error().Err(err).Str("certreq", careqID).Msg("could not retrieve certificate request for VASP")
			continue
		}

		// Sanity check
		if careq.Vasp != vasp.Id {
			logctx.Warn().Str("certreq", careqID).Msg("vasp associated with unrelated certificate request")
			continue
		}

		if careq.Status == models.CertificateRequestState_INITIALIZED {
			if err = models.UpdateCertificateRequestStatus(careq, models.CertificateRequestState_READY_TO_SUBMIT, "registration request received", claims.Email); err != nil {
				return "", err
			}
			if err = s.db.UpdateCertReq(ctx, careq); err != nil {
				return "", err
			}
			ncertreqs++
		}
	}

	switch ncertreqs {
	case 0:
		return "", errors.New("no certificate requests found for VASP registration")
	case 1:
		logctx.Debug().Msg("certificate request marked as ready to submit")
	default:
		logctx.Warn().Int("requests", ncertreqs).Msg("multiple certificate requests marked as ready to submit")
	}

	// Send successful response
	var name string
	if name, err = vasp.Name(); err != nil {
		name = vasp.Id
	}
	return fmt.Sprintf("registration request for %s has been approved and a Sectigo certificate will be requested", name), nil
}

// Reject the VASP registration and notify the contacts of the result.
func (s *Admin) rejectRegistration(vasp *pb.VASP, contacts *models.Contacts, reason string, claims *tokens.Claims, logctx *sentry.Logger) (msg string, err error) {
	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Change the VASP verification status
	if err = models.SetAdminVerificationToken(vasp, ""); err != nil {
		return "", err
	}
	if err := models.UpdateVerificationStatus(vasp, pb.VerificationState_REJECTED, "registration rejected", claims.Email); err != nil {
		return "", err
	}
	if err = s.db.UpdateVASP(ctx, vasp); err != nil {
		return "", err
	}

	// Delete all pending certificate requests
	var (
		ncertreqs int
		careqs    []string
	)

	if careqs, err = models.GetCertReqIDs(vasp); err != nil {
		return "", err
	}

	for _, careqID := range careqs {
		var careq *models.CertificateRequest
		if careq, err = s.db.RetrieveCertReq(ctx, careqID); err != nil {
			logctx.Error().Err(err).Str("certreq", careqID).Msg("could not retrieve certificate request for VASP")
			continue
		}

		// Sanity check
		if careq.Vasp != vasp.Id {
			logctx.Warn().Str("certreq", careqID).Msg("vasp associated with unrelated certificate request")
			continue
		}

		// Delete the certificate request
		if err = s.db.DeleteCertReq(ctx, careq.Id); err != nil {
			logctx.Error().Err(err).Str("id", careq.Id).Msg("could not delete certificate request")
			continue
		}

		// Delete the VASP reference to the certificate request
		if err = models.DeleteCertReqID(vasp, careq.Id); err != nil {
			logctx.Error().Err(err).Str("certreq", careq.Id).Msg("could not delete certificate request ID from VASP")
		} else if err = s.db.UpdateVASP(ctx, vasp); err != nil {
			logctx.Error().Err(err).Msg("could not update VASP with deleted certificate request ID")
		}
		ncertreqs++
	}

	// Log deletion of certificate requests
	switch ncertreqs {
	case 0:
		logctx.Warn().Msg("no certificate requests deleted")
	case 1:
		logctx.Debug().Msg("certificate request deleted")
	default:
		logctx.Warn().Msg("multiple certificate requests deleted")
	}

	// Notify the VASP contacts that the registration request has been rejected.
	if _, err = s.svc.email.SendRejectRegistration(vasp, contacts, reason); err != nil {
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
		err      error
		in       *admin.ResendRequest
		out      *admin.ResendReply
		vasp     *pb.VASP
		vaspID   string
		contacts *models.Contacts
	)

	// Get vaspID from the URL
	vaspID = c.Param("vaspID")

	// Parse incoming JSON data from the client request
	in = new(admin.ResendRequest)
	if err := c.ShouldBind(&in); err != nil {
		sentry.Warn(c).Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(err))
		return
	}

	// Validate resend request
	if in.ID != "" && in.ID != vaspID {
		sentry.Warn(c).Str("id", in.ID).Str("vasp_id", vaspID).Msg("mismatched request ID and URL")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse("the request ID does not match the URL endpoint"))
		return
	}

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Lookup the VASP record associated with the resend request
	if vasp, err = s.db.RetrieveVASP(ctx, vaspID); err != nil {
		sentry.Warn(c).Err(err).Str("id", vaspID).Msg("could not retrieve vasp")
		c.JSON(http.StatusNotFound, admin.ErrorResponse("could not retrieve VASP record by ID"))
		return
	}

	// Lookup the contacts records associated with the VASP
	if contacts, err = s.db.VASPContacts(ctx, vasp); err != nil {
		sentry.Error(c).Err(err).Str("id", vaspID).Msg("could not retrieve vasp contacts")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse("could not retrieve vasp contacts"))
		return
	}

	// Handle different resend request types
	out = &admin.ResendReply{}
	switch in.Action {
	case admin.ResendVerifyContact:
		// Send Verify Contacts needs to include not just the VASPs but also the contacts from the database
		if out.Sent, err = s.svc.email.SendVerifyContacts(contacts); err != nil {
			sentry.Error(c).Err(err).Int("sent", out.Sent).Msg("could not resend verify contacts emails")
			c.JSON(http.StatusInternalServerError, admin.ErrorResponse(fmt.Errorf("could not resend contact verification emails: %s", err)))
			return
		}
		out.Message = "contact verification emails resent to all unverified contacts"

	case admin.ResendReview:
		if out.Sent, err = s.svc.email.SendReviewRequest(vasp); err != nil {
			sentry.Error(c).Err(err).Int("sent", out.Sent).Msg("could not resend review request")
			c.JSON(http.StatusInternalServerError, admin.ErrorResponse(fmt.Errorf("could not resend review request: %s", err)))
			return
		}
		out.Message = "review request resent to TRISA admins"

	case admin.ResendDeliverCerts:
		// TODO: check verification state and cert request state
		// TODO: in order to implement this, we'd have to fetch the certs from Google Secrets
		// TODO: if implemented, log which contact was sent the certs (e.g. technical, admin, etc.)
		// TODO: when above implemented, also log which contact was sent certs in acceptRegistration
		sentry.Warn(c).Msg("resend cert delivery not yet implemented")
		c.JSON(http.StatusNotImplemented, admin.ErrorResponse("resend cert delivery not yet implemented"))
		return

	case admin.ResendRejection:
		// Only send a rejection email if we're in the rejected state
		if vasp.VerificationStatus != pb.VerificationState_REJECTED {
			sentry.Warn(c).Err(err).Str("status", vasp.VerificationStatus.String()).Msg("cannot resend rejection emails in current state")
			c.JSON(http.StatusBadRequest, admin.ErrorResponse("VASP record verification status cannot send rejection email"))
			return
		}

		// A reason must be specified to send a rejection email (it's not stored)
		if in.Reason == "" {
			sentry.Warn(c).Str("resend_type", string(in.Action)).Msg("invalid resend request: missing reason argument")
			c.JSON(http.StatusBadRequest, admin.ErrorResponse("must specify reason for rejection to resend email"))
			return
		}
		if out.Sent, err = s.svc.email.SendRejectRegistration(vasp, contacts, in.Reason); err != nil {
			sentry.Error(c).Err(err).Int("sent", out.Sent).Msg("could not resend rejection emails")
			c.JSON(http.StatusInternalServerError, admin.ErrorResponse(fmt.Errorf("could not resend rejection emails: %s", err)))
			return
		}
		out.Message = "rejection emails resent to all verified contacts"

	default:
		sentry.Warn(c).Str("resend_type", string(in.Action)).Msg("invalid resend request: unhandled resend request type")
		c.JSON(http.StatusBadRequest, admin.ErrorResponse(fmt.Errorf("unknown resend request type %q", in.Action)))
		return
	}

	if err = s.db.UpdateVASP(ctx, vasp); err != nil {
		sentry.Error(c).Str("id", vasp.Id).Msg("error updating email logs on VASP")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse(fmt.Errorf("could not update VASP record: %s", err)))
		return
	}

	log.Info().Str("id", vasp.Id).Int("sent", out.Sent).Str("resend_type", string(in.Action)).Msg("resend request complete")
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

//===========================================================================
// Accessors - used primarily for testing
//===========================================================================

// GetTokenManager returns the underlying token manager for testing.
func (s *Admin) GetTokenManager() *tokens.TokenManager {
	return s.tokens
}

// GetRouter returns the Admin API router for testing purposes.
func (s *Admin) GetRouter() http.Handler {
	return s.router
}
