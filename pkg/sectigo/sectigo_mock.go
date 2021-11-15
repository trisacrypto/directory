package sectigo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// This helps verify that the success paths of the Sectigo API calls are working,
// while also providing a way for tests to inject handlers to test the errors that
// might be returned from the Sectigo API.

type mockServer struct {
	server   *httptest.Server
	handlers map[string]*mockHandlerFunc
}

type mockHandlerFunc struct {
	method      string
	handlerFunc func(c *gin.Context)
}

// setupServer configures routing for the mock server, including adding the default and
// custom handlers.
func (s *mockServer) setupServer() (r *gin.Engine) {
	gin.SetMode(gin.TestMode)
	r = gin.Default()
	for path, h := range s.handlers {
		r.Handle(h.method, path, h.handlerFunc)
	}
	return r
}

// addHandler is a helper function that adds a handler to the mock server's handlers map.
func (s *mockServer) addHandler(path, method string, handlerFunc func(c *gin.Context)) {
	if s.handlers == nil {
		s.handlers = make(map[string]*mockHandlerFunc)
	}
	s.handlers[path] = &mockHandlerFunc{
		method:      method,
		handlerFunc: handlerFunc,
	}
}

// setupHandlers is a helper function which instantiates the handlers map with the
// default handlers defined in this file, and then overrides them with any custom
// handlers passed into the function.
func (s *mockServer) setupHandlers(handlers map[string]*mockHandlerFunc) {
	// Default handlers
	s.addHandler(endpoints[authenticateEP].Path, http.MethodPost, s.authenticate)
	s.addHandler(endpoints[refreshEP].Path, http.MethodPost, s.refresh)
	s.addHandler(endpoints[createSingleCertBatchEP].Path, http.MethodPut, s.createSingleCertBatch)
	s.addHandler(endpoints[uploadCSREP].Path, http.MethodPost, s.uploadCSR)
	s.addHandler(endpoints[batchDetailEP].Path, http.MethodGet, s.batchDetail)
	s.addHandler(endpoints[batchProcessingInfoEP].Path, http.MethodGet, s.batchProcessingInfo)
	s.addHandler(endpoints[downloadEP].Path, http.MethodGet, s.download)
	s.addHandler(endpoints[devicesEP].Path, http.MethodGet, s.devices)
	s.addHandler(endpoints[userAuthoritiesEP].Path, http.MethodGet, s.userAuthorities)
	s.addHandler(endpoints[authorityUserBalanceAvailableEP].Path, http.MethodGet, s.authorityUserBalanceAvailable)
	s.addHandler(endpoints[profilesEP].Path, http.MethodGet, s.profiles)
	s.addHandler(endpoints[profileParametersEP].Path, http.MethodGet, s.profileParameters)
	s.addHandler(endpoints[profileDetailEP].Path, http.MethodGet, s.profileDetail)
	s.addHandler(endpoints[currentUserOrganizationEP].Path, http.MethodGet, s.currentUserOrganization)
	s.addHandler(endpoints[findCertificateEP].Path, http.MethodPost, s.findCertificate)
	s.addHandler(endpoints[revokeCertificateEP].Path, http.MethodGet, s.revokeCertificate)

	// Custom handlers
	for path, h := range handlers {
		s.addHandler(path, h.method, h.handlerFunc)
	}
}

// NewMockServer initializes a new server which mocks the Sectigo REST API. By default,
// it sets up HTTP handlers which return 200 OK responses with mocked data, but custom
// handlers can be passed in by tests to test specific error paths. Note that this
// function modifies the externally used baseURL for the Sectigo endpoint, and the
// caller must close the server when done.
func NewMockServer(handlers map[string]*mockHandlerFunc) (s *mockServer, err error) {
	s = &mockServer{}
	s.setupHandlers(handlers)
	s.server = httptest.NewServer(s.setupServer())
	baseURL, _ = url.Parse(s.server.URL)
	return s, nil
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

func (s *mockServer) authenticate(c *gin.Context) {
	var (
		in      *AuthenticationRequest
		access  string
		refresh string
		err     error
	)
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
	c.JSON(http.StatusOK, &AuthenticationReply{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

func (s *mockServer) refresh(c *gin.Context) {
	var (
		access  string
		refresh string
		err     error
	)
	if access, err = generateToken(); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if refresh, err = generateToken(); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, &AuthenticationReply{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

func (s *mockServer) createSingleCertBatch(c *gin.Context) {
	var in *CreateSingleCertBatchRequest
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, &BatchResponse{
		BatchID:      42,
		CreationDate: time.Now().Format(time.RFC3339),
		Status:       "completed",
		Active:       false,
		BatchName:    in.BatchName,
	})
}

func (s *mockServer) uploadCSR(c *gin.Context) {
	var in *UploadCSRBatchRequest
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, &BatchResponse{
		BatchID:      42,
		CreationDate: time.Now().Format(time.RFC3339),
		Status:       "completed",
		Active:       false,
	})
}

func (s *mockServer) batchDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, "missing batch id")
		return
	}

	c.JSON(http.StatusOK, &BatchResponse{
		BatchID:      42,
		CreationDate: time.Now().Format(time.RFC3339),
		Status:       "completed",
		Active:       false,
	})
}

func (s *mockServer) batchProcessingInfo(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, "missing batch id")
		return
	}

	c.JSON(http.StatusOK, &ProcessingInfoResponse{
		Active:  0,
		Success: 1,
		Failed:  0,
	})
}

func (s *mockServer) download(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, "missing batch id")
		return
	}

	c.JSON(http.StatusOK, &ProcessingInfoResponse{
		Active:  0,
		Success: 1,
		Failed:  0,
	})
}

func (s *mockServer) devices(c *gin.Context) {
	c.JSON(http.StatusOK, &LicensesUsedResponse{
		Ordered: 1,
		Issued:  1,
	})
}

func (s *mockServer) userAuthorities(c *gin.Context) {
	c.JSON(http.StatusOK, []*AuthorityResponse{
		{
			ID:      42,
			Balance: 100,
			Enabled: true,
		},
	})
}

func (s *mockServer) authorityUserBalanceAvailable(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, "missing authority id")
		return
	}
	c.JSON(http.StatusOK, 100)
}

func (s *mockServer) profiles(c *gin.Context) {
	c.JSON(http.StatusOK, []*ProfileResponse{
		{
			ProfileID: 42,
			CA:        "sectigo",
		},
	})
}

func (s *mockServer) profileParameters(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, "missing profile id")
		return
	}
	c.JSON(http.StatusOK, []*ProfileParamsResponse{
		{
			Name:    "foo",
			Message: "bar",
		},
	})
}

func (s *mockServer) profileDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, "missing profile id")
		return
	}
	c.JSON(http.StatusOK, &ProfileDetailResponse{
		ProfileName: "foo",
		ProfileID:   42,
	})
}

func (s *mockServer) currentUserOrganization(c *gin.Context) {
	c.JSON(http.StatusOK, &OrganizationResponse{
		OrganizationID:   42,
		OrganizationName: "foo.io",
	})
}

func (s *mockServer) findCertificate(c *gin.Context) {
	var in *FindCertificateRequest
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, &FindCertificateResponse{
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

func (s *mockServer) revokeCertificate(c *gin.Context) {
	var in *RevokeCertificateRequest
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}
