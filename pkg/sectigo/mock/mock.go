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
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/trisacrypto/directory/pkg/sectigo"
)

var mockServer *Server

type Profile struct {
	ID      string
	Name    string
	Balance int
	Params  []string
}

var subjectParams = [4]string{
	"organizationName",
	"localityName",
	"stateOrProvinceName",
	"countryName",
}

var nameParams = [3]string{
	"commonName",
	"dNSName",
	"pkcs12Password",
}

// All the profiles that the mock knows about.
var validProfiles = []Profile{
	{
		ID:      sectigo.ProfileIDCipherTraceEE,
		Name:    sectigo.ProfileCipherTraceEE,
		Balance: 50,
		Params:  nameParams[:],
	},
	{
		ID:      sectigo.ProfileIDCipherTraceEndEntityCertificate,
		Name:    sectigo.ProfileCipherTraceEndEntityCertificate,
		Balance: 100,
		Params:  append(nameParams[:], subjectParams[:]...),
	},
}

// Server helps verify that the success paths of the Sectigo API calls are working,
// while also providing a way for tests to inject handlers to test the errors that
// might be returned from the Sectigo API.
type Server struct {
	server          *httptest.Server
	router          *gin.Engine
	handlers        sync.Map
	calls           sync.Map
	profileLookup   map[string]Profile
	authorityLookup map[int]Profile
}

// Get returns the current mock server.
func Get() *Server {
	return mockServer
}

// Start() initializes a new server which mocks the Sectigo REST API. By default, it sets up
// HTTP handlers which return 200 OK responses with mocked data, but custom handlers can
// be passed in by tests to test specific error paths. Note that this function modifies
// the externally used baseURL for the Sectigo endpoint, and the caller must call Stop()
// to close the server and undo the mock when done. Start() takes one paramter, which is the
// name of the sectigo profile that the mock server uses to verify request parameters in
// createSingleCertBatch.
func Start(profile string) error {
	gin.SetMode(gin.TestMode)

	mockServer = &Server{
		router:          gin.New(),
		profileLookup:   make(map[string]Profile),
		authorityLookup: make(map[int]Profile),
	}

	// TODO: We are only loading one authority so that we can get consistent behavior
	// from the tests.
	authority := 1
	for _, p := range validProfiles {
		mockServer.profileLookup[p.ID] = p
		if p.Name == profile {
			mockServer.authorityLookup[authority] = p
		}
		authority++
	}

	mockServer.setupHandlers()
	mockServer.server = httptest.NewServer(mockServer.router)
	sectigo.SetBaseURL(mockServer.URL())
	return nil
}

// Stop the test server and reset the Sectigo server URL to complete the tests and cleanup.
func Stop() {
	if mockServer != nil {
		mockServer.server.Close()
	}
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
func (s *Server) GetCalls() sync.Map {
	return s.calls
}

func (s *Server) incrementCall(endpoint string) {
	var (
		v  interface{}
		ok bool
	)
	if v, ok = s.calls.Load(endpoint); ok {
		i := v.(int)
		s.calls.Store(endpoint, i+1)
	} else {
		s.calls.Store(endpoint, 0)
	}
}

func (s *Server) getHandler(endpoint string) (handler gin.HandlerFunc, err error) {
	var (
		h  interface{}
		ok bool
	)
	if h, ok = s.handlers.Load(endpoint); !ok {
		return nil, fmt.Errorf("endpoint not found in handler map: %s", endpoint)
	}
	if handler, ok = h.(gin.HandlerFunc); !ok {
		return nil, fmt.Errorf("unexpected type in handler map: %v", h)
	}
	return handler, nil
}

// Handle is a helper function that adds a handler to the mock server's handlers map and
// returns that handler function when the endpoint is called.
func Handle(endpoint string, handler gin.HandlerFunc) error {
	if _, ok := mockServer.handlers.Load(endpoint); !ok {
		return fmt.Errorf("unhandled endpoint %s", endpoint)
	}
	mockServer.handlers.Store(endpoint, handler)
	return nil
}

func (s *Server) handle(endpoint, method string, handler gin.HandlerFunc) {
	// Get the path from the endpoint
	s.handlers.Store(endpoint, handler)
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
		if mockHandlerFunc, err := s.getHandler(endpoint); err != nil {
			c.JSON(http.StatusNotFound, err)
		} else {
			mockHandlerFunc(c)
		}
	})
	s.calls.Store(endpoint, 0)
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
	s.handle(sectigo.BatchStatusEP, http.MethodGet, s.batchStatus)
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
	s.incrementCall(sectigo.AuthenticateEP)
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
	s.incrementCall(sectigo.RefreshEP)
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
	s.incrementCall(sectigo.CreateSingleCertBatchEP)
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var profile Profile
	var ok bool
	if profile, ok = s.authorityLookup[in.AuthorityID]; !ok {
		c.JSON(http.StatusNotFound, "unrecognized authority id")
		return
	}

	for k, v := range in.ProfileParams {
		if v == "" {
			c.JSON(http.StatusBadRequest, fmt.Errorf("empty profile parameter %s", k))
			return
		}

		found := false
		for _, p := range profile.Params {
			if p == k {
				found = true
				break
			}
		}
		if !found {
			c.JSON(http.StatusBadRequest, fmt.Errorf("unrecognized profile parameter %s", k))
			return
		}
	}

	for _, p := range profile.Params {
		if _, ok := in.ProfileParams[p]; !ok {
			c.JSON(http.StatusBadRequest, fmt.Errorf("missing profile parameter %s", p))
			return
		}
	}

	c.JSON(http.StatusOK, &sectigo.BatchResponse{
		BatchID:      42,
		CreationDate: time.Now().Format(time.RFC3339),
		Status:       sectigo.BatchStatusReadyForDownload,
		Active:       false,
		BatchName:    in.BatchName,
		OrderNumber:  23,
		Profile:      profile.Name,
	})
}

func (s *Server) uploadCSR(c *gin.Context) {
	var (
		form  *multipart.Form
		files []*multipart.FileHeader
		err   error
		ok    bool
	)
	s.incrementCall(sectigo.UploadCSREP)
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
	s.incrementCall(sectigo.BatchDetailEP)
	id := c.Param("param0")
	if id == "" {
		c.JSON(http.StatusBadRequest, "missing batch id")
		return
	}

	c.JSON(http.StatusOK, &sectigo.BatchResponse{
		BatchID:      42,
		CreationDate: time.Now().Format(time.RFC3339),
		Status:       sectigo.BatchStatusReadyForDownload,
		Active:       false,
	})
}

func (s *Server) batchStatus(c *gin.Context) {
	s.incrementCall(sectigo.BatchStatusEP)
	id := c.Param("param0")
	if id == "" {
		c.JSON(http.StatusBadRequest, "missing batch id")
		return
	}

	c.String(http.StatusOK, "Ready for download")
}

func (s *Server) batchProcessingInfo(c *gin.Context) {
	s.incrementCall(sectigo.BatchProcessingInfoEP)
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
	s.incrementCall(sectigo.DownloadEP)
	id := c.Param("param0")
	if id == "" {
		c.JSON(http.StatusBadRequest, "missing batch id")
		return
	}

	// Using runtime.Caller allows us to load the fixture using a path relative to this
	// file, since the mock can be invoked from a few different packages.
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("could not get caller context"))
		return
	}
	f, err := os.Open(filepath.Join(filepath.Dir(thisFile), "testdata", "certificate.zip"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.Header("Content-Disposition", "attachment; filename=certificate.zip")
	c.DataFromReader(http.StatusOK, info.Size(), "application/zip", f, nil)
}

func (s *Server) devices(c *gin.Context) {
	s.incrementCall(sectigo.DevicesEP)
	c.JSON(http.StatusOK, &sectigo.LicensesUsedResponse{
		Ordered: 1,
		Issued:  1,
	})
}

func (s *Server) userAuthorities(c *gin.Context) {
	s.incrementCall(sectigo.UserAuthoritiesEP)
	authorities := make([]*sectigo.AuthorityResponse, 0)
	for id, profile := range s.authorityLookup {
		authorities = append(authorities, &sectigo.AuthorityResponse{
			ID:          id,
			Balance:     profile.Balance,
			Enabled:     true,
			ProfileName: profile.Name,
		})
	}
	c.JSON(http.StatusOK, authorities)
}

func (s *Server) authorityUserBalanceAvailable(c *gin.Context) {
	s.incrementCall(sectigo.AuthorityUserBalanceAvailableEP)
	authority := c.Param("param0")
	if authority == "" {
		c.JSON(http.StatusBadRequest, "missing authority id")
		return
	}
	id, err := strconv.Atoi(authority)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid authority id")
		return
	}
	var auth Profile
	var ok bool
	if auth, ok = s.authorityLookup[id]; !ok {
		c.JSON(http.StatusNotFound, "unrecognized authority id")
		return
	}
	c.JSON(http.StatusOK, auth.Balance)
}

func (s *Server) profiles(c *gin.Context) {
	s.incrementCall(sectigo.ProfilesEP)
	profiles := make([]*sectigo.ProfileResponse, 0)
	for _, profile := range s.authorityLookup {
		id, err := strconv.Atoi(profile.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		profiles = append(profiles, &sectigo.ProfileResponse{
			ProfileID: id,
			CA:        "sectigo",
		})
	}
	c.JSON(http.StatusOK, profiles)
}

func (s *Server) profileParameters(c *gin.Context) {
	s.incrementCall(sectigo.ProfileParametersEP)
	id := c.Param("param0")
	if id == "" {
		c.JSON(http.StatusBadRequest, "missing profile id")
		return
	}
	var profile Profile
	var ok bool
	if profile, ok = s.profileLookup[id]; !ok {
		c.JSON(http.StatusNotFound, "unrecognized profile id")
		return
	}
	params := make([]*sectigo.ProfileParamsResponse, 0)
	for _, param := range profile.Params {
		params = append(params, &sectigo.ProfileParamsResponse{
			Name:     param,
			Required: true,
		})
	}

	c.JSON(http.StatusOK, params)
}

func (s *Server) profileDetail(c *gin.Context) {
	s.incrementCall(sectigo.ProfileDetailEP)
	id := c.Param("param0")
	if id == "" {
		c.JSON(http.StatusBadRequest, "missing profile id")
		return
	}
	var profile Profile
	var ok bool
	if profile, ok = s.profileLookup[id]; !ok {
		c.JSON(http.StatusNotFound, "unrecognized profile id")
		return
	}

	profileId, err := strconv.Atoi(profile.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &sectigo.ProfileDetailResponse{
		ProfileName: profile.Name,
		ProfileID:   profileId,
	})
}

func (s *Server) currentUserOrganization(c *gin.Context) {
	s.incrementCall(sectigo.CurrentUserOrganizationEP)
	c.JSON(http.StatusOK, &sectigo.OrganizationResponse{
		OrganizationID:   42,
		OrganizationName: "foo.io",
	})
}

func (s *Server) findCertificate(c *gin.Context) {
	s.incrementCall(sectigo.FindCertificateEP)
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
	s.incrementCall(sectigo.RevokeCertificateEP)
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
