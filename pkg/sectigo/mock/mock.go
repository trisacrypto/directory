/*
 * Package mock provides an httptest.Server that allows mock interactions with the
 * Sectigo API for both unit tests of the Sectigo package and integration tests with GDS
 * Sectigo/Certificate operations.
 *
 * The basic use case is to create a *mock.Server using mock.New(). Requests to the
 * mock server endpoint (fetched using mock.Server.URL()) will return "happy path"
 * responses for a single fake certificate interaction. You can also modify the response
 * from the mock server by using mock.Handle or one of the handler methods.
 */
package mock

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/trisacrypto/directory/pkg/sectigo"
)

// Server helps verify that the success paths of the Sectigo API calls are working,
// while also providing a way for tests to inject handlers to test the errors that
// might be returned from the Sectigo API.
type Server struct {
	server   *httptest.Server
	router   *gin.Engine
	handlers map[string]gin.HandlerFunc
	calls    map[string]int
}

// New initializes a new server which mocks the Sectigo REST API. By default, it sets up
// HTTP handlers which return 200 OK responses with mocked data, but custom handlers can
// be passed in by tests to test specific error paths. Note that this function modifies
// the externally used baseURL for the Sectigo endpoint, and the caller must close the
// server when done.
func New() (s *Server, err error) {
	gin.SetMode(gin.TestMode)

	s = &Server{
		handlers: make(map[string]gin.HandlerFunc),
		router:   gin.New(),
		calls:    make(map[string]int),
	}

	s.setupHandlers()
	s.server = httptest.NewServer(s.router)
	sectigo.SetBaseURL(s.URL())
	return s, nil
}

// Close the test server to complete the tests and cleanup.
func (s *Server) Close() {
	s.server.Close()
	sectigo.ResetBaseURL()
}

// URL returns the URL of the test server.
func (s *Server) URL() *url.URL {
	u, err := url.Parse(s.server.URL)
	if err != nil {
		panic(err)
	}
	return u
}

// GetCalls returns the map of called endpoints.
func (s *Server) GetCalls() map[string]int {
	return s.calls
}

// Handle is a helper function that adds a handler to the mock server's handlers map and
// returns that handler function when the endpoint is called.
func (s *Server) Handle(endpoint string, handler gin.HandlerFunc) error {
	if _, ok := s.handlers[endpoint]; !ok {
		return fmt.Errorf("unhandled endpoint %s", endpoint)
	}
	s.handlers[endpoint] = handler
	return nil
}

func (s *Server) handle(endpoint, method string, handler gin.HandlerFunc) {
	// Get the path from the endpoint
	s.handlers[endpoint] = handler
	ep, err := sectigo.Endpoint(endpoint)
	if err != nil {
		panic(err)
	}

	// Modify any %d format directives to gin :params
	for i := 0; strings.Contains(ep.Path, "%d"); i++ {
		ep.Path = strings.Replace(ep.Path, "%d", fmt.Sprintf(":param%d", i), 1)
	}

	// If there are any other % directives, panic
	if strings.Contains(ep.Path, "%") {
		panic(fmt.Errorf("invalid path %q", ep.Path))
	}

	s.router.Handle(method, ep.Path, func(c *gin.Context) {
		if mockHandlerFunc, ok := s.handlers[endpoint]; ok {
			mockHandlerFunc(c)
		} else {
			c.JSON(http.StatusNotFound, "endpoint not found")
		}
	})
	s.calls[endpoint] = 0
}

// setupHandlers is a helper function which instantiates the handlers map with the
// default handlers defined in this file. Note that a Sectigo endpoint needs to be
// configured with a default method in order for it to be available for mocking during
// an external test, otherwise it will return a 404 error.
func (s *Server) setupHandlers() {
	// Default handlers
	s.handle(sectigo.AuthenticateEP, http.MethodPost, s.authenticate)
	s.handle(sectigo.RefreshEP, http.MethodPost, s.refresh)
	s.handle(sectigo.CreateSingleCertBatchEP, http.MethodPut, s.createSingleCertBatch)
	s.handle(sectigo.UploadCSREP, http.MethodPost, s.uploadCSR)
	s.handle(sectigo.BatchDetailEP, http.MethodGet, s.batchDetail)
	s.handle(sectigo.BatchProcessingInfoEP, http.MethodGet, s.batchProcessingInfo)
	s.handle(sectigo.DownloadEP, http.MethodGet, s.download)
	s.handle(sectigo.DevicesEP, http.MethodGet, s.devices)
	s.handle(sectigo.UserAuthoritiesEP, http.MethodGet, s.userAuthorities)
	s.handle(sectigo.AuthorityUserBalanceAvailableEP, http.MethodGet, s.authorityUserBalanceAvailable)
	s.handle(sectigo.ProfilesEP, http.MethodGet, s.profiles)
	s.handle(sectigo.ProfileParametersEP, http.MethodGet, s.profileParameters)
	s.handle(sectigo.ProfileDetailEP, http.MethodGet, s.profileDetail)
	s.handle(sectigo.CurrentUserOrganizationEP, http.MethodGet, s.currentUserOrganization)
	s.handle(sectigo.FindCertificateEP, http.MethodPost, s.findCertificate)
	s.handle(sectigo.RevokeCertificateEP, http.MethodPost, s.revokeCertificate)
}

func generateToken() (string, error) {
	var token *jwt.Token
	claims := jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
	}
	if token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims); token == nil {
		return "", fmt.Errorf("could not generate jwt token")
	}
	var signed string
	var err error
	if signed, err = token.SignedString([]byte("foo")); err != nil {
		return "", err
	}
	return signed, nil
}

func (s *Server) authenticate(c *gin.Context) {
	var (
		in      *sectigo.AuthenticationRequest
		access  string
		refresh string
		err     error
	)
	s.calls[sectigo.AuthenticateEP]++
	if err = c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if access, err = generateToken(); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if refresh, err = generateToken(); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, &sectigo.AuthenticationReply{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

func (s *Server) refresh(c *gin.Context) {
	var (
		access  string
		refresh string
		err     error
	)
	s.calls[sectigo.RefreshEP]++
	if access, err = generateToken(); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if refresh, err = generateToken(); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, &sectigo.AuthenticationReply{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

func (s *Server) createSingleCertBatch(c *gin.Context) {
	var in *sectigo.CreateSingleCertBatchRequest
	s.calls[sectigo.CreateSingleCertBatchEP]++
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, &sectigo.BatchResponse{
		BatchID:      42,
		CreationDate: time.Now().Format(time.RFC3339),
		Status:       "completed",
		Active:       false,
		BatchName:    in.BatchName,
	})
}

func (s *Server) uploadCSR(c *gin.Context) {
	var (
		form  *multipart.Form
		files []*multipart.FileHeader
		err   error
		ok    bool
	)
	s.calls[sectigo.UploadCSREP]++
	if form, err = c.MultipartForm(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if _, ok := form.Value["profileId"]; !ok {
		c.JSON(http.StatusBadRequest, fmt.Errorf("missing profile id"))
		return
	}
	if files, ok = form.File["files"]; !ok {
		c.JSON(http.StatusBadRequest, fmt.Errorf("multipart form missing files part"))
		return
	}
	if len(files) != 1 {
		c.JSON(http.StatusBadRequest, fmt.Errorf("expected 1 file, got %d", len(files)))
		return
	}
	c.JSON(http.StatusOK, &sectigo.BatchResponse{
		BatchID:      42,
		CreationDate: time.Now().Format(time.RFC3339),
		Status:       "completed",
		Active:       false,
	})
}

func (s *Server) batchDetail(c *gin.Context) {
	s.calls[sectigo.BatchDetailEP]++
	id := c.Param("param0")
	if id == "" {
		c.JSON(http.StatusBadRequest, "missing batch id")
		return
	}

	c.JSON(http.StatusOK, &sectigo.BatchResponse{
		BatchID:      42,
		CreationDate: time.Now().Format(time.RFC3339),
		Status:       "completed",
		Active:       false,
	})
}

func (s *Server) batchProcessingInfo(c *gin.Context) {
	s.calls[sectigo.BatchProcessingInfoEP]++
	id := c.Param("param0")
	if id == "" {
		c.JSON(http.StatusBadRequest, "missing batch id")
		return
	}

	c.JSON(http.StatusOK, &sectigo.ProcessingInfoResponse{
		Active:  0,
		Success: 1,
		Failed:  0,
	})
}

func (s *Server) download(c *gin.Context) {
	s.calls[sectigo.DownloadEP]++
	id := c.Param("param0")
	if id == "" {
		c.JSON(http.StatusBadRequest, "missing batch id")
		return
	}

	someJSON := struct {
		Field string `json:"field"`
	}{
		Field: "foo",
	}

	c.JSON(http.StatusOK, someJSON)
}

func (s *Server) devices(c *gin.Context) {
	s.calls[sectigo.DevicesEP]++
	c.JSON(http.StatusOK, &sectigo.LicensesUsedResponse{
		Ordered: 1,
		Issued:  1,
	})
}

func (s *Server) userAuthorities(c *gin.Context) {
	s.calls[sectigo.UserAuthoritiesEP]++
	c.JSON(http.StatusOK, []*sectigo.AuthorityResponse{
		{
			ID:      42,
			Balance: 100,
			Enabled: true,
		},
	})
}

func (s *Server) authorityUserBalanceAvailable(c *gin.Context) {
	s.calls[sectigo.AuthorityUserBalanceAvailableEP]++
	id := c.Param("param0")
	if id == "" {
		c.JSON(http.StatusBadRequest, "missing authority id")
		return
	}
	c.JSON(http.StatusOK, 100)
}

func (s *Server) profiles(c *gin.Context) {
	s.calls[sectigo.ProfilesEP]++
	c.JSON(http.StatusOK, []*sectigo.ProfileResponse{
		{
			ProfileID: 42,
			CA:        "sectigo",
		},
	})
}

func (s *Server) profileParameters(c *gin.Context) {
	s.calls[sectigo.ProfileParametersEP]++
	id := c.Param("param0")
	if id == "" {
		c.JSON(http.StatusBadRequest, "missing profile id")
		return
	}
	c.JSON(http.StatusOK, []*sectigo.ProfileParamsResponse{
		{
			Name:    "foo",
			Message: "bar",
		},
	})
}

func (s *Server) profileDetail(c *gin.Context) {
	s.calls[sectigo.ProfileDetailEP]++
	id := c.Param("param0")
	if id == "" {
		c.JSON(http.StatusBadRequest, "missing profile id")
		return
	}
	c.JSON(http.StatusOK, &sectigo.ProfileDetailResponse{
		ProfileName: "foo",
		ProfileID:   42,
	})
}

func (s *Server) currentUserOrganization(c *gin.Context) {
	s.calls[sectigo.CurrentUserOrganizationEP]++
	c.JSON(http.StatusOK, &sectigo.OrganizationResponse{
		OrganizationID:   42,
		OrganizationName: "foo.io",
	})
}

func (s *Server) findCertificate(c *gin.Context) {
	s.calls[sectigo.FindCertificateEP]++
	var in *sectigo.FindCertificateRequest
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, &sectigo.FindCertificateResponse{
		TotalCount: 1,
		Items: []struct {
			DeviceID     int    `json:"deviceId"`
			CommonName   string `json:"commonName"`
			SerialNumber string `json:"serialNumber"`
			CreationDate string `json:"creationDate"`
			Status       string `json:"status"`
		}{
			{
				DeviceID:     42,
				CommonName:   in.CommonName,
				SerialNumber: in.SerialNumber,
				CreationDate: time.Now().Format(time.RFC3339),
				Status:       "valid",
			},
		},
	})
}

func (s *Server) revokeCertificate(c *gin.Context) {
	s.calls[sectigo.RevokeCertificateEP]++
	id := c.Param("param0")
	if id == "" {
		c.JSON(http.StatusBadRequest, "missing certificate id")
		return
	}
	var in *sectigo.RevokeCertificateRequest
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}
