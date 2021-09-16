package gds

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/trisacrypto/directory/pkg"
	admin "github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/store"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

// NewAdmin creates a new GDS admin server derived from a parent Service.
func NewAdmin(svc *Service) (a *Admin, err error) {
	a = &Admin{
		svc:  svc,
		conf: &svc.conf.Admin,
		db:   svc.db,
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
	svc     *Service            // The parent Service the admin server uses to interact with other components
	srv     *http.Server        // The HTTP server that listens on its own independent port
	conf    *config.AdminConfig // The admin server specific configuration (alias to s.svc.conf.Admin)
	db      store.Store         // Database connection for loading objects (alias to s.svc.db)
	router  *gin.Engine         // The HTTP handler and associated middleware
	healthy bool                // application state of the server
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

	// Add CORS configuration
	// TODO: configure origins from the environment rather than hard-coding
	s.router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	s.router.Use(s.Available())

	// Add the v2 API routes
	v2 := s.router.Group("/v2")
	v2.GET("/status", s.Status)
	v2.GET("/vasps", s.ListVASPs)
	v2.POST("/vasps/:vaspID/review", s.Review)
	v2.POST("/vasps/:vaspID/resend", s.Resend)

	return nil
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
			vasp := iter.VASP()

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
			contacts := models.VerifiedContacts(vasp)
			snippet.VerifiedContacts = make([]string, 0, len(contacts))
			for key := range contacts {
				snippet.VerifiedContacts = append(snippet.VerifiedContacts, key)
			}

			// Append to list in reply
			out.VASPs = append(out.VASPs, snippet)
		}
	}

	if err = iter.Error(); err != nil {
		log.Warn().Err(err).Msg("could iterate over vasps in store")
		c.JSON(http.StatusInternalServerError, admin.ErrorResponse(err))
		return
	}

	// Successful request, return the VASP list JSON data
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

	// Accept or reject the request
	out = &admin.ReviewReply{}
	if in.Accept {
		if out.Message, err = s.acceptRegistration(vasp); err != nil {
			log.Error().Err(err).Msg("could not accept VASP registration")
			c.JSON(http.StatusInternalServerError, admin.ErrorResponse("unable to accept VASP registration request"))
			return
		}
	} else {
		if out.Message, err = s.rejectRegistration(vasp, in.RejectReason); err != nil {
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
func (s *Admin) acceptRegistration(vasp *pb.VASP) (msg string, err error) {
	// Change the VASP verification status
	if err = models.SetAdminVerificationToken(vasp, ""); err != nil {
		return "", err
	}
	vasp.VerifiedOn = time.Now().Format(time.RFC3339)
	vasp.VerificationStatus = pb.VerificationState_REVIEWED
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
			req.Status = models.CertificateRequestState_READY_TO_SUBMIT
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
func (s *Admin) rejectRegistration(vasp *pb.VASP, reason string) (msg string, err error) {
	// Change the VASP verification status
	if err = models.SetAdminVerificationToken(vasp, ""); err != nil {
		return "", err
	}
	vasp.VerificationStatus = pb.VerificationState_REJECTED
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
