package sectigo

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// This helps verify that the success paths of the Sectigo API calls are working,
// while also providing a way for tests to inject handlers to test the errors that
// might be returned from the Sectigo API.

type mockServer struct {
	server *httptest.Server
}

type mockHandlerFunc struct {
	method      string
	handlerFunc func(c *gin.Context)
}

func (s *mockServer) setupServer(handlers map[string]*mockHandlerFunc) (r *gin.Engine) {
	gin.SetMode(gin.TestMode)
	r = gin.Default()

	r.POST(endpoints[authenticateEP].Path, s.authenticate)
	r.POST(endpoints[refreshEP].Path, s.refresh)
	r.PUT(endpoints[createSingleCertBatchEP].Path, s.createSingleCertBatch)
	r.POST(endpoints[uploadCSREP].Path, s.uploadCSR)
	r.GET(endpoints[batchDetailEP].Path, s.batchDetail)
	r.GET(endpoints[batchProcessingInfoEP].Path, s.batchProcessingInfo)
	r.GET(endpoints[downloadEP].Path, s.download)
	r.GET(endpoints[devicesEP].Path, s.devices)
	r.GET(endpoints[userAuthoritiesEP].Path, s.userAuthorities)
	r.GET(endpoints[authorityUserBalanceAvailableEP].Path, s.authorityUserBalanceAvailable)
	r.GET(endpoints[profilesEP].Path, s.profiles)
	r.GET(endpoints[profileParametersEP].Path, s.profileParameters)
	r.GET(endpoints[profileDetailEP].Path, s.profileDetail)
	r.GET(endpoints[currentUserOrganizationEP].Path, s.currentUserOrganization)
	r.POST(endpoints[findCertificateEP].Path, s.findCertificate)
	r.GET(endpoints[revokeCertificateEP].Path, s.revokeCertificate)

	if handlers != nil {
		for path, h := range handlers {
			r.Handle(h.method, path, h.handlerFunc)
		}
	}
	return r
}

func NewMockServer(handlers map[string]*mockHandlerFunc) (s *mockServer, err error) {
	s = &mockServer{}
	s.server = httptest.NewUnstartedServer(s.setupServer(handlers))
	cert, err := tls.LoadX509KeyPair(filepath.Join("testdata", "certs", "server.crt"), filepath.Join("testdata", "certs", "server.key"))
	if err != nil {
		return nil, err
	}
	s.server.TLS = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	s.server.StartTLS()
	setBaseURL(&url.URL{Scheme: "http", Host: s.server.Listener.Addr().String()})
	fmt.Println(urlFor(authenticateEP))
	return s, nil
}

func (s *mockServer) authenticate(c *gin.Context) {
	fmt.Println("got here")

	var in *AuthenticationRequest
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, &AuthenticationReply{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	})
}

func (s *mockServer) refresh(c *gin.Context) {
	c.JSON(http.StatusOK, &AuthenticationReply{
		AccessToken:  "new-access",
		RefreshToken: "new-refresh",
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
