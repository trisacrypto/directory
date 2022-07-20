package gds_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	admin "github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/emails"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

// httpRequest is a helper struct to make it easier to organize all the different
// parameters required for making an in-code API request.
type httpRequest struct {
	method  string
	path    string
	headers map[string]string
	params  map[string]string
	in      interface{}
	claims  *tokens.Claims
}

// makeRequest creates a new HTTP request and returns the gin context along with the
// http.ResponseRecorder.
func (s *gdsTestSuite) makeRequest(request *httpRequest) (*gin.Context, *httptest.ResponseRecorder) {
	var body io.ReadWriter
	var err error
	require := s.Require()

	// Encode the JSON request
	if request.in != nil {
		body = &bytes.Buffer{}
		err = json.NewEncoder(body).Encode(request.in)
		require.NoError(err)
	} else {
		body = nil
	}

	// Construct the HTTP request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	if request.claims != nil {
		c.Set(admin.UserClaims, request.claims)
	}
	c.Request = httptest.NewRequest(request.method, request.path, body)
	c.Request.Header.Add("Content-Type", "application/json")
	if request.headers != nil {
		for k, v := range request.headers {
			c.Request.Header.Add(k, v)
		}
	}

	for k, v := range request.params {
		c.Params = append(c.Params, gin.Param{
			Key:   k,
			Value: v,
		})
	}
	return c, w
}

// doRequest is a helper function for making an admin API request and retrieving
// the response.
func (s *gdsTestSuite) doRequest(handle gin.HandlerFunc, c *gin.Context, w *httptest.ResponseRecorder, reply interface{}) (res *http.Response) {
	require := s.Require()
	// Call the admin function and return the HTTP response
	handle(c)
	res = w.Result()
	defer res.Body.Close()
	if reply != nil {
		bytes, err := ioutil.ReadAll(res.Body)
		require.NoError(err)
		err = json.Unmarshal(bytes, reply)
		require.NoError(err)
	}
	return res
}

// createAccessCredential is a helper function for generating JWT credential strings
// using the mocked token manager for authentication tests.
func (s *gdsTestSuite) createAccessString(creds map[string]interface{}) string {
	require := s.Require()
	tm := s.svc.GetAdmin().GetTokenManager()
	accessToken, err := tm.CreateAccessToken(creds)
	require.NoError(err)
	access, err := tm.Sign(accessToken)
	require.NoError(err)
	return access
}

// APIError is a helper function for asserting that an expected API error is returned.
func (s *gdsTestSuite) APIError(expectedCode int, expectedMessage string, rep *http.Response) {
	require := s.Require()
	require.NotNil(rep, "no HTTP response returned")
	require.Equal(expectedCode, rep.StatusCode, "expected status code does not match response")

	defer rep.Body.Close()
	data := &admin.Reply{}
	require.NoError(json.NewDecoder(rep.Body).Decode(data), "could not decode admin.Reply JSON")
	require.NotNil(data, "no data was returned")
	require.False(data.Success, "API returned a success response")
	require.Equal(expectedMessage, data.Error, "error message mismatch")
}

// Test that the middleware returns the corect error when making unauthenticated
// requests to protected endpoints.
func (s *gdsTestSuite) TestMiddleware() {
	endpoints := []struct {
		name      string
		method    string
		path      string
		authorize bool
		csrf      bool
	}{
		// CSRF protected endpoints
		{"authenticate", http.MethodPost, "/v2/authenticate", false, true},
		{"reauthenticate", http.MethodPost, "/v2/reauthenticate", false, true},
		// Authenticated endpoints
		{"summary", http.MethodGet, "/v2/summary", true, false},
		{"autocomplete", http.MethodGet, "/v2/autocomplete", true, false},
		{"reviews", http.MethodGet, "/v2/reviews", true, false},
		{"listVASPs", http.MethodGet, "/v2/vasps", true, false},
		{"retrieveVASP", http.MethodGet, "/v2/vasps/42", true, false},
		{"listReviewNotes", http.MethodGet, "/v2/vasps/42/notes", true, false},
		{"listCertificates", http.MethodGet, "/v2/vasps/42/certificates", true, false},
		// Authenticated and CSRF protected endpoints
		{"updateVASP", http.MethodPatch, "/v2/vasps/42", true, true},
		{"deleteVASP", http.MethodDelete, "/v2/vasps/42", true, true},
		{"replaceContact", http.MethodPut, "/v2/vasps/42/contacts/kind", true, true},
		{"deleteContact", http.MethodDelete, "/v2/vasps/42/contacts/kind", true, true},
		{"review", http.MethodPost, "/v2/vasps/42/review", true, true},
		{"resend", http.MethodPost, "/v2/vasps/42/resend", true, true},
		{"createReviewNote", http.MethodPost, "/v2/vasps/42/notes", true, true},
		{"updateReviewNote", http.MethodPut, "/v2/vasps/42/notes/1", true, true},
		{"deleteReviewNote", http.MethodDelete, "/v2/vasps/42/notes/1", true, true},
	}
	serv := httptest.NewServer(s.svc.GetAdmin().GetRouter())
	defer serv.Close()

	// Endpoints should return unavailable when in maintenance mode/unhealthy
	s.svc.GetAdmin().SetHealth(false)
	for _, endpoint := range endpoints {
		s.T().Run(endpoint.name, func(t *testing.T) {
			r, err := http.NewRequest(endpoint.method, serv.URL+endpoint.path, nil)
			require.NoError(t, err)
			res, err := http.DefaultClient.Do(r)
			require.NoError(t, err)
			defer res.Body.Close()
			require.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
		})
	}

	// Endpoints that are authenticated or CSRF protected
	s.svc.GetAdmin().SetHealth(true)
	for _, endpoint := range endpoints {
		switch {
		case endpoint.authorize && endpoint.csrf:
			s.T().Run(endpoint.name, func(t *testing.T) {
				// Request is not authenticated
				r, err := http.NewRequest(endpoint.method, serv.URL+endpoint.path, nil)
				require.NoError(t, err)
				res, err := http.DefaultClient.Do(r)
				require.NoError(t, err)
				require.Equal(t, http.StatusUnauthorized, res.StatusCode)
				// Request is authenticated but CSRF token is missing
				r, err = http.NewRequest(endpoint.method, serv.URL+endpoint.path, nil)
				require.NoError(t, err)
				creds := map[string]interface{}{
					"sub":     "102374163855881761273",
					"hd":      "example.com",
					"email":   "jon@example.com",
					"name":    "Jon Doe",
					"picture": "https://foo.googleusercontent.com/test!/Aoh14gJceTrUA",
				}
				access := s.createAccessString(creds)
				r.Header.Add("Authorization", "Bearer "+access)
				res, err = http.DefaultClient.Do(r)
				require.NoError(t, err)
				defer res.Body.Close()
				require.Equal(t, http.StatusForbidden, res.StatusCode)
			})
		case endpoint.authorize || endpoint.csrf:
			var status int
			if endpoint.authorize {
				status = http.StatusUnauthorized
			} else {
				status = http.StatusForbidden
			}
			s.T().Run(endpoint.name, func(t *testing.T) {
				r, err := http.NewRequest(endpoint.method, serv.URL+endpoint.path, nil)
				require.NoError(t, err)
				res, err := http.DefaultClient.Do(r)
				require.NoError(t, err)
				defer res.Body.Close()
				require.Equal(t, status, res.StatusCode)
			})
		default:
			s.Require().Fail(fmt.Sprintf("misconfigured test: %s, authorize or csrf must be true", endpoint.name))
		}
	}
}

// Test that we get a good response from ProtectAuthenticate.
func (s *gdsTestSuite) TestProtectAuthenticate() {
	a := s.svc.GetAdmin()
	require := s.Require()
	request := &httpRequest{
		method: http.MethodGet,
		path:   "/v2/authenticate",
	}

	actual := &admin.Reply{}
	c, w := s.makeRequest(request)
	res := s.doRequest(a.ProtectAuthenticate, c, w, actual)
	require.Equal(http.StatusOK, res.StatusCode)
	expected := &admin.Reply{Success: true}
	require.Equal(expected, actual)

	// Double cookie tokens should be set
	cookies := res.Cookies()
	require.Len(cookies, 2)
	for _, cookie := range cookies {
		require.Equal(s.svc.GetConf().Admin.CookieDomain, cookie.Domain)
	}
}

// Test the Authenticate endpoint.
func (s *gdsTestSuite) TestAuthenticate() {
	require := s.Require()
	s.LoadFullFixtures()
	a := s.svc.GetAdmin()

	// Missing credential
	request := &httpRequest{
		method: http.MethodPost,
		path:   "/v2/authenticate",
		in:     &admin.AuthRequest{},
	}
	c, w := s.makeRequest(request)
	res := s.doRequest(a.Authenticate, c, w, nil)
	require.Equal(http.StatusUnauthorized, res.StatusCode)

	// Invalid credential
	request.in = &admin.AuthRequest{
		Credential: "invalid",
	}
	c, w = s.makeRequest(request)
	res = s.doRequest(a.Authenticate, c, w, nil)
	require.Equal(http.StatusUnauthorized, res.StatusCode)

	// Unauthorized domain
	creds := map[string]interface{}{
		"sub":     "102374163855881761273",
		"hd":      "unauthorized.dev",
		"email":   "jon@gds.dev",
		"name":    "Jon Doe",
		"picture": "https://foo.googleusercontent.com/test!/Aoh14gJceTrUA",
	}
	access := s.createAccessString(creds)
	request.in = &admin.AuthRequest{
		Credential: access,
	}
	c, w = s.makeRequest(request)
	res = s.doRequest(a.Authenticate, c, w, nil)
	require.Equal(http.StatusUnauthorized, res.StatusCode)

	// Successful authentication
	creds["hd"] = "gds.dev"
	access = s.createAccessString(creds)
	request.in = &admin.AuthRequest{
		Credential: access,
	}
	c, w = s.makeRequest(request)
	res = s.doRequest(a.Authenticate, c, w, nil)
	require.Equal(http.StatusOK, res.StatusCode)
	// Double cookie tokens should be set
	cookies := res.Cookies()
	require.Len(cookies, 2)
	for _, cookie := range cookies {
		require.Equal(s.svc.GetConf().Admin.CookieDomain, cookie.Domain)
	}
}

// Test the Reauthenticate endpoint.
func (s *gdsTestSuite) TestReauthenticate() {
	s.LoadFullFixtures()
	require := s.Require()
	a := s.svc.GetAdmin()
	tm := a.GetTokenManager()

	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Audience:  jwt.ClaimStrings{"http://localhost"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	refreshToken, err := tm.CreateRefreshToken(accessToken)
	require.NoError(err)
	access, err := tm.Sign(accessToken)
	require.NoError(err)
	refresh, err := tm.Sign(refreshToken)
	require.NoError(err)

	// Missing access token
	request := &httpRequest{
		method: http.MethodPost,
		path:   "/v2/reauthenticate",
		in: &admin.AuthRequest{
			Credential: refresh,
		},
	}
	c, w := s.makeRequest(request)
	res := s.doRequest(a.Reauthenticate, c, w, nil)
	require.Equal(http.StatusUnauthorized, res.StatusCode)

	// Invalid access token
	request.headers = map[string]string{
		"Authorization": "Bearer invalid",
	}
	c, w = s.makeRequest(request)
	res = s.doRequest(a.Reauthenticate, c, w, nil)
	require.Equal(http.StatusUnauthorized, res.StatusCode)

	// Missing refresh token
	request.in = &admin.AuthRequest{}
	request.headers = map[string]string{
		"Authorization": "Bearer " + access,
	}
	c, w = s.makeRequest(request)
	res = s.doRequest(a.Reauthenticate, c, w, nil)
	require.Equal(http.StatusUnauthorized, res.StatusCode)

	// Invalid refresh token
	request.in = &admin.AuthRequest{
		Credential: "invalid",
	}
	c, w = s.makeRequest(request)
	res = s.doRequest(a.Reauthenticate, c, w, nil)
	require.Equal(http.StatusUnauthorized, res.StatusCode)

	// Mismatched access and refresh tokens
	claims.ID = uuid.NewString()
	otherToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	other, err := tm.Sign(otherToken)
	require.NoError(err)
	request.in = &admin.AuthRequest{
		Credential: refresh,
	}
	request.headers = map[string]string{
		"Authorization": "Bearer " + other,
	}
	c, w = s.makeRequest(request)
	res = s.doRequest(a.Reauthenticate, c, w, nil)
	require.Equal(http.StatusUnauthorized, res.StatusCode)

	// Successful reauthentication
	request.in = &admin.AuthRequest{
		Credential: refresh,
	}
	request.headers = map[string]string{
		"Authorization": "Bearer " + access,
	}
	c, w = s.makeRequest(request)
	res = s.doRequest(a.Reauthenticate, c, w, nil)
	require.Equal(http.StatusOK, res.StatusCode)
	// Double cookie tokens should be set
	cookies := res.Cookies()
	require.Len(cookies, 2)
	for _, cookie := range cookies {
		require.Equal(s.svc.GetConf().Admin.CookieDomain, cookie.Domain)
	}
}

// Test that the Summary endpoint returns the correct response.
func (s *gdsTestSuite) TestSummary() {
	s.LoadFullFixtures()
	require := s.Require()
	a := s.svc.GetAdmin()

	request := &httpRequest{
		method: http.MethodGet,
		path:   "/v2/summary",
	}
	actual := &admin.SummaryReply{}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.Summary, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)

	// Test against the expected response
	expected := &admin.SummaryReply{
		VASPsCount:           14,
		PendingRegistrations: 6,
		ContactsCount:        39,
		VerifiedContacts:     26,
		CertificatesIssued:   3,
		Statuses: map[string]int{
			pb.VerificationState_APPEALED.String():            1,
			pb.VerificationState_ERRORED.String():             1,
			pb.VerificationState_ISSUING_CERTIFICATE.String(): 1,
			pb.VerificationState_PENDING_REVIEW.String():      2,
			pb.VerificationState_REVIEWED.String():            1,
			pb.VerificationState_REJECTED.String():            2,
			pb.VerificationState_SUBMITTED.String():           1,
			pb.VerificationState_VERIFIED.String():            5,
		},
		CertReqs: map[string]int{
			models.CertificateRequestState_COMPLETED.String():       3,
			models.CertificateRequestState_PROCESSING.String():      1,
			models.CertificateRequestState_CR_ERRORED.String():      1,
			models.CertificateRequestState_CR_REJECTED.String():     1,
			models.CertificateRequestState_INITIALIZED.String():     3,
			models.CertificateRequestState_READY_TO_SUBMIT.String(): 1,
		},
	}
	require.Equal(expected, actual, "unexpected summary reply, have the fixtures changed?")
}

// Test that the Autocomplete endpoint returns the correct response.
func (s *gdsTestSuite) TestAutocomplete() {
	s.LoadSmallFixtures()
	require := s.Require()
	a := s.svc.GetAdmin()

	request := &httpRequest{
		method: http.MethodGet,
		path:   "/v2/autocomplete",
	}
	actual := &admin.AutocompleteReply{}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.Autocomplete, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)

	// Construct the expected response
	charlieID := s.fixtures[vasps]["charliebank"].(*pb.VASP).Id
	deltaID := s.fixtures[vasps]["delta"].(*pb.VASP).Id
	expected := &admin.AutocompleteReply{
		Names: map[string]string{
			"trisa.charliebank.io":         charlieID,
			"https://trisa.charliebank.io": "https://trisa.charliebank.io",
			"CharlieBank":                  charlieID,
			"trisa.delta.io":               deltaID,
			"https://trisa.delta.io":       "https://trisa.delta.io",
			"Delta Assets":                 deltaID,
		},
	}
	require.Equal(expected, actual)
}

// Test the ListVASPs endpoint.
func (s *gdsTestSuite) TestListVASPs() {
	s.LoadSmallFixtures()
	require := s.Require()
	a := s.svc.GetAdmin()

	// Impose an ordering so we can verify the results.
	sortByName := func(s []admin.VASPSnippet) {
		sort.Slice(s, func(i, j int) bool {
			return s[i].Name < s[j].Name
		})
	}

	snippets := []admin.VASPSnippet{
		{
			ID:                  s.fixtures[vasps]["charliebank"].(*pb.VASP).Id,
			Name:                "CharlieBank",
			CommonName:          "trisa.charliebank.io",
			RegisteredDirectory: "trisatest.net",
			VerificationStatus:  pb.VerificationState_SUBMITTED.String(),
			VerifiedContacts: map[string]bool{
				"administrative": false,
				"billing":        false,
				"legal":          false,
			},
		},
		{
			ID:                  s.fixtures[vasps]["delta"].(*pb.VASP).Id,
			Name:                "Delta Assets",
			CommonName:          "trisa.delta.io",
			RegisteredDirectory: "trisatest.net",
			VerificationStatus:  pb.VerificationState_APPEALED.String(),
			VerifiedContacts: map[string]bool{
				"billing": true,
				"legal":   true,
			},
		},
	}

	// List all VASPs on the same page
	request := &httpRequest{
		method: http.MethodGet,
		path:   "/v2/vasps",
	}
	actual := &admin.ListVASPsReply{}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.ListVASPs, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	require.Equal(len(snippets), actual.Count)
	require.Equal(1, actual.Page)
	require.Equal(100, actual.PageSize)
	require.Len(actual.VASPs, len(snippets))
	sortByName(actual.VASPs)

	// Make sure the last modified timestamps are the same since the fixture will be
	// updated when it is inserted into the database.
	for i, vasp := range actual.VASPs {
		snippets[i].LastUpdated = vasp.LastUpdated
	}

	require.Equal(snippets, actual.VASPs)

	// List VASPs with an invalid status
	request.path = "/v2/vasps?status=invalid"
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ListVASPs, c, w, nil)
	require.Equal(http.StatusBadRequest, rep.StatusCode)

	// List VASPs with the specified status
	request.path = "/v2/vasps?status=" + snippets[0].VerificationStatus
	actual = &admin.ListVASPsReply{}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ListVASPs, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	require.Equal(1, actual.Count)
	require.Equal(1, actual.Page)
	require.Equal(100, actual.PageSize)
	require.Len(actual.VASPs, 1)
	require.Equal(snippets[0], actual.VASPs[0])

	// List VASPs with multiple status filters
	request.path = "/v2/vasps?status=" + snippets[0].VerificationStatus + "&status=" + snippets[1].VerificationStatus
	actual = &admin.ListVASPsReply{}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ListVASPs, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	require.Equal(2, actual.Count)
	require.Equal(1, actual.Page)
	require.Equal(100, actual.PageSize)
	require.Len(actual.VASPs, 2)
	sortByName(actual.VASPs)
	require.Equal(snippets, actual.VASPs)

	// List VASPs on multiple pages
	request.path = "/v2/vasps?page=1&page_size=1"
	actual = &admin.ListVASPsReply{}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ListVASPs, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	require.Equal(len(snippets), actual.Count)
	require.Equal(1, actual.Page)
	require.Equal(1, actual.PageSize)
	require.Len(actual.VASPs, 1)
	pageResults := []admin.VASPSnippet{snippets[0]}

	request.path = "/v2/vasps?page=2&page_size=1"
	actual = &admin.ListVASPsReply{}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ListVASPs, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	require.Equal(len(snippets), actual.Count)
	require.Equal(2, actual.Page)
	require.Equal(1, actual.PageSize)
	require.Len(actual.VASPs, 1)
	pageResults = append(pageResults, snippets[1])
	sortByName(pageResults)
	require.Equal(snippets, pageResults)

	request.path = "/v2/vasps?page=3&page_size=1"
	actual = &admin.ListVASPsReply{}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ListVASPs, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	require.Equal(len(snippets), actual.Count)
	require.Equal(3, actual.Page)
	require.Equal(1, actual.PageSize)
	require.Len(actual.VASPs, 0)
}

// Test the RetrieveVASP endpoint.
func (s *gdsTestSuite) TestRetrieveVASP() {
	s.LoadFullFixtures()
	require := s.Require()
	a := s.svc.GetAdmin()

	// Retrieve a VASP that doesn't exist
	request := &httpRequest{
		method: http.MethodGet,
		path:   "/v2/vasps/invalid",
	}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.RetrieveVASP, c, w, nil)
	require.Equal(http.StatusNotFound, rep.StatusCode)

	hotel := s.fixtures[vasps]["hotel"].(*pb.VASP)

	// Retrieve a VASP that exists
	request.path = "/v2/vasps/" + hotel.Id
	request.params = map[string]string{
		"vaspID": hotel.Id,
	}
	actual := &admin.RetrieveVASPReply{}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.RetrieveVASP, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	expected := &admin.RetrieveVASPReply{
		Name:     "Hotel Corp",
		Traveler: false,
		AuditLog: []map[string]interface{}{
			{
				"current_state":  pb.VerificationState_SUBMITTED.String(),
				"description":    "register request received",
				"previous_state": pb.VerificationState_NO_VERIFICATION.String(),
				"source":         "automated",
				"timestamp":      "2021-06-17T11:12:23Z",
			},
			{
				"current_state":  pb.VerificationState_EMAIL_VERIFIED.String(),
				"description":    "completed email verification",
				"previous_state": pb.VerificationState_SUBMITTED.String(),
				"source":         "automated",
				"timestamp":      "2021-06-21T14:34:49Z",
			},
			{
				"current_state":  pb.VerificationState_PENDING_REVIEW.String(),
				"description":    "review email sent",
				"previous_state": pb.VerificationState_EMAIL_VERIFIED.String(),
				"source":         "automated",
				"timestamp":      "2021-07-01T20:59:04Z",
			},
			{
				"current_state":  pb.VerificationState_REVIEWED.String(),
				"description":    "registration request received",
				"previous_state": pb.VerificationState_PENDING_REVIEW.String(),
				"source":         "admin@rotational.io",
				"timestamp":      "2021-08-10T21:37:14Z",
			},
			{
				"current_state":  pb.VerificationState_ISSUING_CERTIFICATE.String(),
				"description":    "issuing certificate",
				"previous_state": pb.VerificationState_REVIEWED.String(),
				"source":         "automated",
				"timestamp":      "2021-08-25T18:03:15Z",
			},
			{
				"current_state":  pb.VerificationState_VERIFIED.String(),
				"description":    "certificate issued",
				"previous_state": pb.VerificationState_ISSUING_CERTIFICATE.String(),
				"source":         "automated",
				"timestamp":      "2021-10-21T15:52:08Z",
			},
		},
	}

	actualVASP, err := remarshalProto(vasps, actual.VASP)
	require.NoError(err, "could not remarshal actual VASP")

	// RetrieveVASP removes the extra data and modifies the serial numbers before
	// returning the VASP
	s.CompareFixture(vasps, hotel.Id, actualVASP, true, true)

	// Check the verified contacts
	require.Len(actual.VerifiedContacts, 2)

	// Verify that the identity certificate serial number was converted to a capital hex encoded string
	expectedSerial := fmt.Sprintf("%X", hotel.IdentityCertificate.SerialNumber)
	actualSerial, ok := actual.VASP["identity_certificate"].(map[string]interface{})["serial_number"]
	require.True(ok, "identity_certificate.serial_number not found in VASP json")
	require.Equal(expectedSerial, actualSerial)

	// Verify that the signing cerificate serial numbers were converted to capital hex encoded strings
	actualCerts, ok := actual.VASP["signing_certificates"].([]interface{})
	require.True(ok, "signing_certificates not found in VASP json")
	require.Len(actualCerts, len(hotel.SigningCertificates))
	for i, cert := range actualCerts {
		actualSerial, ok := cert.(map[string]interface{})["serial_number"]
		require.True(ok, "signing certificate serial number not found in VASP json")
		expectedSerial := fmt.Sprintf("%X", hotel.SigningCertificates[i].SerialNumber)
		require.Equal(expectedSerial, actualSerial)
	}

	// Compare the rest of the non-vasp results
	actual.VASP = nil
	actual.VerifiedContacts = nil
	require.Equal(expected, actual)
}

func (s *gdsTestSuite) TestUpdateVASP() {
	s.LoadSmallFixtures()
	defer s.ResetFixtures()

	a := s.svc.GetAdmin()

	// Attempt to update a VASP that doesn't exist
	request := &httpRequest{
		method: http.MethodPatch,
		path:   "/v2/vasps/invalid",
		params: map[string]string{
			"vaspID": "invalid",
		},
		in: &admin.UpdateVASPRequest{},
		claims: &tokens.Claims{
			Email: "admin@example.com",
		},
	}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.UpdateVASP, c, w, nil)
	s.APIError(http.StatusNotFound, "could not retrieve VASP record by ID", rep)

	// Update a VASP that exists
	charlieVASP := s.fixtures[vasps]["charliebank"].(*pb.VASP)
	charlieID := charlieVASP.Id
	request.path = "/v2/vasps/" + charlieID
	request.params["vaspID"] = charlieID

	// Test request VASP and URL do not match returns a 400 error
	request.in = &admin.UpdateVASPRequest{VASP: "bce77e90-82e0-4685-8139-6ec5d4b83615"}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.UpdateVASP, c, w, nil)
	s.APIError(http.StatusBadRequest, "the request ID does not match the URL endpoint", rep)

	// Test an update with no changes returns a 400 error
	request.in = &admin.UpdateVASPRequest{}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.UpdateVASP, c, w, nil)
	s.APIError(http.StatusBadRequest, "no updates made to VASP record", rep)

	// TODO: Test bad business category (not parseable) returns 400 error
	// TODO: Test updating website, business category, vasp categories, and established on
	// TODO: Test invalid IVMS 101 returns 400 error
	// TODO: Test update VASP entity
	// TODO: Test update TRIXO form
	// TODO: Test compute common name from endpoint retuns an error if endpoint is "foo"
	// TODO: Test endpoint-only change with no change to common name is successful
	// TODO: Test an update to common name for reviewed VASP returns a 400 error
	// TODO: Test no certificate requests updated returns an error
	// TODO: Test common name change with incorrect endpoint returns an error
	// TODO: Test common name-only change with correct endpoint is successful
	// TODO: Test endpoint-only change with change to common name is successful
	// TODO: Test common name and endpoint change is successful
}

func (s *gdsTestSuite) TestDeleteVASP() {
	s.LoadFullFixtures()
	defer s.ResetFixtures()
	defer s.loadReferenceFixtures()

	require := s.Require()
	a := s.svc.GetAdmin()

	deltaID := s.fixtures[vasps]["delta"].(*pb.VASP).Id
	julietID := s.fixtures[vasps]["juliet"].(*pb.VASP).Id
	xrayID := s.fixtures[certreqs]["xray"].(*models.CertificateRequest).Id

	golf := s.fixtures[vasps]["golfbucks"].(*pb.VASP)

	// Attempt to delete a VASP that doesn't exist
	request := &httpRequest{
		method: http.MethodDelete,
		path:   "/v2/vasps/invalid",
		params: map[string]string{
			"vaspID": "invalid",
		},
		claims: &tokens.Claims{
			Email: "admin@example.com",
		},
	}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.DeleteVASP, c, w, nil)
	s.APIError(http.StatusNotFound, "could not retrieve VASP record by ID", rep)

	// VASP is in an invalid state for deletion
	request.path = "/v2/vasps/" + deltaID
	request.params["vaspID"] = deltaID
	msg := "cannot delete VASP in its current state"
	for status := pb.VerificationState_REVIEWED; status < pb.VerificationState_ERRORED; status++ {
		s.SetVerificationStatus(deltaID, status)
		c, w = s.makeRequest(request)
		rep = s.doRequest(a.DeleteVASP, c, w, nil)
		s.APIError(http.StatusBadRequest, msg, rep)
	}

	// Successfully deleting a VASP and its certificate requests
	id := golf.Id
	for status := pb.VerificationState_NO_VERIFICATION; status < pb.VerificationState_REVIEWED; status++ {
		s.SetVerificationStatus(id, status)
		request.path = "/v2/vasps/" + id
		request.params["vaspID"] = id
		c, w = s.makeRequest(request)
		rep = s.doRequest(a.DeleteVASP, c, w, nil)
		require.Equal(http.StatusOK, rep.StatusCode)
		_, err := s.svc.GetStore().RetrieveVASP(golf.Id)
		require.Error(err)

		// Recreate the VASP
		golf.Id = ""
		id, err = s.svc.GetStore().CreateVASP(golf)
		require.NoError(err)
	}

	// Make sure we can also delete a VASP in the ERRORED state
	s.SetVerificationStatus(julietID, pb.VerificationState_ERRORED)
	request.path = "/v2/vasps/" + julietID
	request.params["vaspID"] = julietID
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.DeleteVASP, c, w, nil)
	require.Equal(http.StatusOK, rep.StatusCode)
	_, err := s.svc.GetStore().RetrieveVASP(julietID)
	require.Error(err)
	_, err = s.svc.GetStore().RetrieveCertReq(xrayID)
	require.Error(err)
}

// Test the ListCertificates endpoint
func (s *gdsTestSuite) TestListCertificates() {
	s.LoadFullFixtures()
	defer s.ResetFixtures()

	require := s.Require()
	a := s.svc.GetAdmin()

	// CharlieBank has no certificates
	charlieID := s.fixtures[vasps]["charliebank"].(*pb.VASP).Id

	// HotelCorp has a few certificates
	hotelID := s.fixtures[vasps]["hotel"].(*pb.VASP).Id
	uniform := s.fixtures[certs]["uniform"].(*models.Certificate)
	uniformDetails, err := wire.Rewire(uniform.Details)
	require.NoError(err)
	victor := s.fixtures[certs]["victor"].(*models.Certificate)
	victorDetails, err := wire.Rewire(victor.Details)
	require.NoError(err)
	zulu := s.fixtures[certs]["zulu"].(*models.Certificate)
	zuluDetails, err := wire.Rewire(zulu.Details)
	require.NoError(err)
	certificates := []admin.Certificate{
		{
			SerialNumber: uniform.Id,
			IssuedAt:     uniform.Details.NotBefore,
			ExpiresAt:    uniform.Details.NotAfter,
			Status:       "ISSUED",
			Details:      uniformDetails,
		},
		{
			SerialNumber: victor.Id,
			IssuedAt:     victor.Details.NotBefore,
			ExpiresAt:    victor.Details.NotAfter,
			Status:       "EXPIRED",
			Details:      victorDetails,
		},
		{
			SerialNumber: zulu.Id,
			IssuedAt:     zulu.Details.NotBefore,
			ExpiresAt:    zulu.Details.NotAfter,
			Status:       "REVOKED",
			Details:      zuluDetails,
		},
	}

	// Attempt to retrieve certificates for a VASP that doesn't exist
	request := &httpRequest{
		method: http.MethodGet,
		path:   "/v2/vasps/invalid/certificates",
		params: map[string]string{
			"vaspID": "invalid",
		},
		claims: &tokens.Claims{
			Email: "admin@example.com",
		},
	}
	actual := &admin.ListCertificatesReply{}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.ListCertificates, c, w, nil)
	s.APIError(http.StatusNotFound, "could not retrieve VASP record by ID", rep)

	// No certificates exist for the VASP
	request.path = "/v2/vasps/" + charlieID + "/certificates"
	request.params["vaspID"] = charlieID
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ListCertificates, c, w, nil)
	require.Equal(http.StatusOK, rep.StatusCode)
	require.Empty(actual.Certificates)

	// The VASP contains a few certificates
	request.path = "/v2/vasps/" + hotelID + "/certificates"
	request.params["vaspID"] = hotelID
	actual = &admin.ListCertificatesReply{}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ListCertificates, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	require.Len(actual.Certificates, len(certificates))
	require.ElementsMatch(certificates, actual.Certificates)
}

// Test the ReplaceContact endpoint
func (s *gdsTestSuite) TestReplaceContact() {
	s.LoadSmallFixtures()
	defer s.ResetFixtures()
	defer emails.PurgeMockEmails()
	defer s.loadReferenceFixtures()

	require := s.Require()
	a := s.svc.GetAdmin()

	charlieVASP := s.fixtures[vasps]["charliebank"].(*pb.VASP)
	charlieID := charlieVASP.Id

	// Attempt to update a VASP that doesn't exist
	contact := charlieVASP.Contacts.Administrative
	contactRequest, err := wire.Rewire(contact)
	require.NoError(err, "could not rewire contact")
	request := &httpRequest{
		method: http.MethodPut,
		path:   "/v2/vasps/invalid/contacts/administrative",
		params: map[string]string{
			"vaspID": "invalid",
			"kind":   "administrative",
		},
		in: &admin.ReplaceContactRequest{
			VASP:    "invalid",
			Kind:    "administrative",
			Contact: contactRequest,
		},
		claims: &tokens.Claims{
			Email: "admin@example.com",
		},
	}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.ReplaceContact, c, w, nil)
	s.APIError(http.StatusNotFound, "could not retrieve VASP record by ID", rep)

	// Replacing a contact kind that doesn't exist
	request.path = "/v2/vasps/" + charlieID + "/contacts/invalid"
	request.params["vaspID"] = charlieID
	request.params["kind"] = "invalid"
	request.in = &admin.ReplaceContactRequest{
		VASP:    charlieID,
		Kind:    "invalid",
		Contact: contactRequest,
	}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ReplaceContact, c, w, nil)
	s.APIError(http.StatusBadRequest, "invalid contact kind provided", rep)

	// Replacing a contact with no data
	request.path = "/v2/vasps/" + charlieID + "/contacts/administrative"
	request.params["kind"] = "administrative"
	request.in = &admin.ReplaceContactRequest{
		VASP:    charlieID,
		Kind:    "administrative",
		Contact: map[string]interface{}{},
	}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ReplaceContact, c, w, nil)
	s.APIError(http.StatusBadRequest, "contact data is required for ReplaceContact request", rep)

	// Test removing a contact email is not allowed
	contact = &pb.Contact{
		Email: "",
	}
	contactRequest, err = wire.Rewire(contact)
	require.NoError(err, "could not rewire contact")
	request.in = &admin.ReplaceContactRequest{
		VASP:    charlieID,
		Kind:    "administrative",
		Contact: contactRequest,
	}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ReplaceContact, c, w, nil)
	require.Equal(http.StatusBadRequest, rep.StatusCode)

	// Successfully replacing a contact name
	contact = charlieVASP.Contacts.Administrative
	contact.Name = "Clark Kent"
	contactRequest, err = wire.Rewire(contact)
	require.NoError(err, "could not rewire contact")
	request.in = &admin.ReplaceContactRequest{
		VASP:    charlieID,
		Kind:    "administrative",
		Contact: contactRequest,
	}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ReplaceContact, c, w, nil)
	require.Equal(http.StatusOK, rep.StatusCode)
	vasp, err := s.svc.GetStore().RetrieveVASP(charlieID)
	require.NoError(err, "could not retrieve VASP record")
	require.Equal(contact.Name, vasp.Contacts.Administrative.Name)
	require.Equal(contact.Email, vasp.Contacts.Administrative.Email)
	require.Equal(contact.Phone, vasp.Contacts.Administrative.Phone)
	require.Nil(vasp.Contacts.Technical)

	// Successfully replacing a contact email
	contact = charlieVASP.Contacts.Administrative
	contact.Email = "clark.kent@charliebank.com"
	contactRequest, err = wire.Rewire(contact)
	require.NoError(err, "could not rewire contact")
	request.in = &admin.ReplaceContactRequest{
		VASP:    charlieID,
		Kind:    "administrative",
		Contact: contactRequest,
	}
	c, w = s.makeRequest(request)
	adminSent := time.Now()
	rep = s.doRequest(a.ReplaceContact, c, w, nil)
	require.Equal(http.StatusOK, rep.StatusCode)
	vasp, err = s.svc.GetStore().RetrieveVASP(charlieID)
	require.NoError(err, "could not retrieve VASP record")
	require.Equal(contact.Name, vasp.Contacts.Administrative.Name)
	require.Equal(contact.Email, vasp.Contacts.Administrative.Email)
	require.Equal(contact.Phone, vasp.Contacts.Administrative.Phone)
	// Should no longer be verified
	token, verified, err := models.GetContactVerification(vasp.Contacts.Administrative)
	require.NoError(err, "could not retrieve contact verification")
	require.NotEmpty(token)
	require.False(verified)

	// Successfully adding a new contact
	contact = &pb.Contact{
		Name:  "Lois Lane",
		Email: "lois.lane@charliebank.com",
	}
	contactRequest, err = wire.Rewire(contact)
	require.NoError(err, "could not rewire contact")
	request.params["kind"] = "technical"
	request.in = &admin.ReplaceContactRequest{
		VASP:    charlieID,
		Kind:    "technical",
		Contact: contactRequest,
	}
	c, w = s.makeRequest(request)
	technicalSent := time.Now()
	rep = s.doRequest(a.ReplaceContact, c, w, nil)
	require.Equal(http.StatusOK, rep.StatusCode)
	vasp, err = s.svc.GetStore().RetrieveVASP(charlieID)
	require.NoError(err, "could not retrieve VASP record")
	require.NotNil(vasp.Contacts.Technical)
	require.Equal(contact.Name, vasp.Contacts.Technical.Name)
	require.Equal(contact.Email, vasp.Contacts.Technical.Email)
	// Should not be verified
	token, verified, err = models.GetContactVerification(vasp.Contacts.Technical)
	require.NoError(err, "could not retrieve contact verification")
	require.NotEmpty(token)
	require.False(verified)

	// Should send verification emails to the two contacts
	messages := []*emailMeta{
		{
			contact:   vasp.Contacts.Administrative,
			to:        "clark.kent@charliebank.com",
			from:      s.svc.GetConf().Email.ServiceEmail,
			subject:   emails.VerifyContactRE,
			reason:    "verify_contact",
			timestamp: adminSent,
		},
		{
			contact:   vasp.Contacts.Technical,
			to:        "lois.lane@charliebank.com",
			from:      s.svc.GetConf().Email.ServiceEmail,
			subject:   emails.VerifyContactRE,
			reason:    "verify_contact",
			timestamp: technicalSent,
		},
	}
	s.CheckEmails(messages)
}

// Test the DeleteContact endpoint
func (s *gdsTestSuite) TestDeleteContact() {
	s.LoadSmallFixtures()
	defer s.ResetFixtures()
	defer s.loadReferenceFixtures()

	require := s.Require()
	a := s.svc.GetAdmin()

	charlieVASP := s.fixtures[vasps]["charliebank"].(*pb.VASP)
	charlieID := charlieVASP.Id

	// Attempt to delete a VASP that doesn't exist
	request := &httpRequest{
		method: http.MethodDelete,
		path:   "/v2/vasps/invalid/contacts/administrative",
		params: map[string]string{
			"vaspID": "invalid",
			"kind":   "administrative",
		},
		claims: &tokens.Claims{
			Email: "admin@example.com",
		},
	}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.DeleteContact, c, w, nil)
	s.APIError(http.StatusNotFound, "could not retrieve VASP record by ID", rep)

	// Attempt to delete a contact that doesn't exist
	request.path = "/v2/vasps/charliebank/contacts/invalid"
	request.params["vaspID"] = charlieID
	request.params["kind"] = "invalid"
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.DeleteContact, c, w, nil)
	s.APIError(http.StatusBadRequest, "invalid contact kind provided", rep)

	// Test deleting a contact
	request.path = "/v2/vasps/charliebank/contacts/administrative"
	request.params["kind"] = "administrative"
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.DeleteContact, c, w, nil)
	require.Equal(http.StatusOK, rep.StatusCode)
	vasp, err := s.svc.GetStore().RetrieveVASP(charlieID)
	require.NoError(err, "could not retrieve VASP record")
	require.Nil(vasp.Contacts.Administrative)
	require.NotNil(vasp.Contacts.Billing)

	// Test deleting another contact
	request.path = "/v2/vasps/charliebank/contacts/billing"
	request.params["kind"] = "billing"
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.DeleteContact, c, w, nil)
	require.Equal(http.StatusOK, rep.StatusCode)
	vasp, err = s.svc.GetStore().RetrieveVASP(charlieID)
	require.NoError(err, "could not retrieve VASP record")
	require.Nil(vasp.Contacts.Billing)

	// Test deleting the remaining contact is not allowed
	request.path = "/v2/vasps/charliebank/contacts/legal"
	request.params["kind"] = "legal"
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.DeleteContact, c, w, nil)
	require.Equal(http.StatusBadRequest, rep.StatusCode)
}

// Test the CreateReviewNote endpoint.
func (s *gdsTestSuite) TestCreateReviewNote() {
	s.LoadFullFixtures()
	defer s.ResetFixtures()

	require := s.Require()
	a := s.svc.GetAdmin()

	charlieID := s.fixtures[vasps]["charliebank"].(*pb.VASP).Id

	// Supplying an invalid note ID
	request := &httpRequest{
		method: http.MethodPost,
		path:   "/v2/vasps/" + charlieID + "/notes",
		in: &admin.ModifyReviewNoteRequest{
			NoteID: "invalid slug",
		},
		claims: &tokens.Claims{
			Email: "admin@example.com",
		},
	}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.CreateReviewNote, c, w, nil)
	require.Equal(http.StatusBadRequest, rep.StatusCode)

	// Supplying an invalid VASP ID
	request.in = &admin.ModifyReviewNoteRequest{
		VASP: "invalid",
	}
	request.params = map[string]string{
		"vaspID": "invalid",
	}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.CreateReviewNote, c, w, nil)
	require.Equal(http.StatusNotFound, rep.StatusCode)

	// Successfully creating a review note
	request.in = &admin.ModifyReviewNoteRequest{
		VASP:   charlieID,
		NoteID: "89bceb0e-41aa-11ec-9d29-acde48001122",
		Text:   "foo",
	}
	request.params = map[string]string{
		"vaspID": charlieID,
	}
	actual := &admin.ReviewNote{}
	created := time.Now()
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.CreateReviewNote, c, w, actual)
	require.Equal(http.StatusCreated, rep.StatusCode)
	// Validate returned note
	require.Equal("89bceb0e-41aa-11ec-9d29-acde48001122", actual.ID)
	ts, err := time.Parse(time.RFC3339, actual.Created)
	require.NoError(err)
	require.True(ts.Sub(created) < time.Minute)
	require.Empty(actual.Modified)
	require.Equal(request.claims.Email, actual.Author)
	require.Empty(actual.Editor)
	require.Equal("foo", actual.Text)
	// Record on the database should be updated
	v, err := s.svc.GetStore().RetrieveVASP(charlieID)
	require.NoError(err)
	notes, err := models.GetReviewNotes(v)
	require.NoError(err)
	require.Len(notes, 2)
	require.Contains(notes, actual.ID)
	require.Equal(actual.ID, notes[actual.ID].Id)
	require.Equal(actual.Text, notes[actual.ID].Text)
	require.Equal(actual.Author, notes[actual.ID].Author)
	require.Equal(actual.Created, notes[actual.ID].Created)
	require.Empty(notes[actual.ID].Modified)
	require.Empty(notes[actual.ID].Editor)

	// Should not be able to create a review note if it already exists
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.CreateReviewNote, c, w, nil)
	require.Equal(http.StatusBadRequest, rep.StatusCode)
}

// Test the ListReviewNotes endpoint.
func (s *gdsTestSuite) TestListReviewNotes() {
	s.LoadFullFixtures()
	require := s.Require()
	a := s.svc.GetAdmin()

	charlieID := s.fixtures[vasps]["charliebank"].(*pb.VASP).Id

	// Supplying an invalid VASP ID
	request := &httpRequest{
		method: http.MethodGet,
		path:   "/v2/vasps/invalid/notes",
		params: map[string]string{
			"vaspID": "invalid",
		},
	}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.ListReviewNotes, c, w, nil)
	require.Equal(http.StatusNotFound, rep.StatusCode)

	// Successfully listing review notes
	request.params = map[string]string{
		"vaspID": charlieID,
	}
	actual := &admin.ListReviewNotesReply{}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ListReviewNotes, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	require.Len(actual.Notes, 1)
	require.Equal("admin@trisa.io", actual.Notes[0].Author)
	require.NotEmpty(actual.Notes[0].Created)
	require.NotEmpty(actual.Notes[0].Editor)
	require.NotEmpty(actual.Notes[0].ID)
	require.NotEmpty(actual.Notes[0].Modified)
	require.NotEmpty(actual.Notes[0].Text)
}

// Test the UpdateReviewNote endpoint.
func (s *gdsTestSuite) TestUpdateReviewNote() {
	s.LoadFullFixtures()
	defer s.ResetFixtures()

	require := s.Require()
	a := s.svc.GetAdmin()

	charlieID := s.fixtures[vasps]["charliebank"].(*pb.VASP).Id
	noteID := "5daa4ff0-9011-4b61-a8b3-9b0ff1ec4927"

	// Supplying an invalid note ID
	request := &httpRequest{
		method: http.MethodPut,
		path:   "/v2/vasps/" + charlieID + "/notes/invalid",
		in: &admin.ModifyReviewNoteRequest{
			VASP:   charlieID,
			NoteID: "invalid slug",
		},
		params: map[string]string{
			"vaspID": charlieID,
			"noteID": "invalid slug",
		},
		claims: &tokens.Claims{
			Email: "admin@example.com",
		},
	}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.UpdateReviewNote, c, w, nil)
	require.Equal(http.StatusNotFound, rep.StatusCode)

	// Supplying an invalid VASP ID
	request.in = &admin.ModifyReviewNoteRequest{
		VASP:   "invalid",
		NoteID: noteID,
	}
	request.params = map[string]string{
		"vaspID": "invalid",
		"noteID": noteID,
	}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.UpdateReviewNote, c, w, nil)
	require.Equal(http.StatusNotFound, rep.StatusCode)

	// Successfully updating a review note
	request.in = &admin.ModifyReviewNoteRequest{
		VASP:   charlieID,
		NoteID: noteID,
		Text:   "bar",
	}
	request.params = map[string]string{
		"vaspID": charlieID,
		"noteID": noteID,
	}
	actual := &admin.ReviewNote{}
	modified := time.Now()
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.UpdateReviewNote, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	// Validate returned note
	require.Equal(noteID, actual.ID)
	require.Equal("2021-07-12T20:12:02Z", actual.Created)
	ts, err := time.Parse(time.RFC3339, actual.Modified)
	require.NoError(err)
	require.True(ts.Sub(modified) < time.Minute)
	require.Equal("admin@trisa.io", actual.Author)
	require.Equal(request.claims.Email, actual.Editor)
	require.Equal("bar", actual.Text)
	// Record on the database should be updated
	v, err := s.svc.GetStore().RetrieveVASP(charlieID)
	require.NoError(err)
	notes, err := models.GetReviewNotes(v)
	require.NoError(err)
	require.Len(notes, 1)
	require.Contains(notes, actual.ID)
	require.Equal(actual.ID, notes[actual.ID].Id)
	require.Equal(actual.Text, notes[actual.ID].Text)
	require.Equal(actual.Author, notes[actual.ID].Author)
	require.Equal(actual.Created, notes[actual.ID].Created)
	require.Equal(actual.Modified, notes[actual.ID].Modified)
	require.Equal(actual.Editor, notes[actual.ID].Editor)
}

// Test the DeleteReviewNote endpoint.
func (s *gdsTestSuite) TestDeleteReviewNote() {
	s.LoadFullFixtures()
	defer s.ResetFixtures()

	require := s.Require()
	a := s.svc.GetAdmin()

	charlieID := s.fixtures[vasps]["charliebank"].(*pb.VASP).Id
	noteID := "5daa4ff0-9011-4b61-a8b3-9b0ff1ec4927"

	// Supplying an invalid note ID
	request := &httpRequest{
		method: http.MethodDelete,
		path:   "/v2/vasps/" + charlieID + "/notes/invalid",
		params: map[string]string{
			"vaspID": charlieID,
			"noteID": "invalid slug",
		},
	}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.DeleteReviewNote, c, w, nil)
	require.Equal(http.StatusNotFound, rep.StatusCode)

	// Supplying an invalid VASP ID
	request.params = map[string]string{
		"vaspID": "invalid",
		"noteID": noteID,
	}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.DeleteReviewNote, c, w, nil)
	require.Equal(http.StatusNotFound, rep.StatusCode)

	// Successfully deleting a review note
	request.params = map[string]string{
		"vaspID": charlieID,
		"noteID": noteID,
	}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.DeleteReviewNote, c, w, nil)
	require.Equal(http.StatusOK, rep.StatusCode)

	// Record on the database should be deleted
	v, err := s.svc.GetStore().RetrieveVASP(charlieID)
	require.NoError(err)
	notes, err := models.GetReviewNotes(v)
	require.NoError(err)
	require.Len(notes, 0)
}

// Test the Review Token endpoint
func (s *gdsTestSuite) TestReviewToken() {
	s.LoadFullFixtures()

	require := s.Require()
	a := s.svc.GetAdmin()

	echo := s.fixtures[vasps]["echo"].(*pb.VASP)
	juliet := s.fixtures[vasps]["juliet"].(*pb.VASP)

	require.NotEqual(pb.VerificationState_PENDING_REVIEW, echo.VerificationStatus, "echo must not be in PENDING_REVIEW for this test to pass")
	require.Equal(pb.VerificationState_PENDING_REVIEW, juliet.VerificationStatus, "juliet must be in PENDING_REVIEW for this test to pass")

	// Test not found with invalid ID
	request := &httpRequest{
		method: http.MethodGet,
		path:   "/v2/vasps/invalid/review",
		params: map[string]string{
			"vaspID": "invalid",
		},
		claims: &tokens.Claims{
			Email: "admin@example.com",
		},
	}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.ReviewToken, c, w, nil)
	require.Equal(http.StatusNotFound, rep.StatusCode)

	// Test not found with VASP not in a PENDING_REVIEW state
	request.path = fmt.Sprintf("/v2/vasps/%s/review", echo.Id)
	request.params["vaspID"] = echo.Id
	c2, w2 := s.makeRequest(request)
	rep2 := s.doRequest(a.ReviewToken, c2, w2, nil)
	require.Equal(http.StatusNotFound, rep2.StatusCode)

	// Ensure Juliet has an admin verification token
	avt, err := models.GetAdminVerificationToken(juliet)
	require.NoError(err, "could not get admin verification token from juliet")
	require.NotEmpty(avt, "juliet fixture does not have an admin verification token")

	// Test valid response returned when VASP is in a PENDING_REVIEW state
	out := &admin.ReviewTokenReply{}
	request.path = fmt.Sprintf("/v2/vasps/%s/review", juliet.Id)
	request.params["vaspID"] = juliet.Id
	c3, w3 := s.makeRequest(request)
	rep3 := s.doRequest(a.ReviewToken, c3, w3, out)
	require.Equal(http.StatusOK, rep3.StatusCode)
	require.Equal(avt, out.AdminVerificationToken)
}

// Test the Review endpoint with invalid parameters.
func (s *gdsTestSuite) TestReviewInvalid() {
	s.LoadFullFixtures()

	require := s.Require()
	a := s.svc.GetAdmin()

	julietVASP := s.fixtures[vasps]["juliet"].(*pb.VASP)

	// Supplying an invalid VASP ID
	request := &httpRequest{
		method: http.MethodPost,
		path:   "/v2/vasps/invalid/review",
		in: &admin.ReviewRequest{
			ID:                     "invalid",
			AdminVerificationToken: "supersecrettoken",
			Accept:                 true,
		},
		params: map[string]string{
			"vaspID": "invalid",
		},
		claims: &tokens.Claims{
			Email: "admin@example.com",
		},
	}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.Review, c, w, nil)
	require.Equal(http.StatusNotFound, rep.StatusCode)

	// No verification token supplied
	request.in = &admin.ReviewRequest{
		ID:     julietVASP.Id,
		Accept: true,
	}
	request.params = map[string]string{
		"vaspID": julietVASP.Id,
	}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.Review, c, w, nil)
	require.Equal(http.StatusBadRequest, rep.StatusCode)

	// Wrong verification token supplied
	request.in = &admin.ReviewRequest{
		ID:                     julietVASP.Id,
		AdminVerificationToken: "invalid",
		Accept:                 true,
	}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.Review, c, w, nil)
	require.Equal(http.StatusUnauthorized, rep.StatusCode)
}

// Test the Review endpoint for the accept case.
func (s *gdsTestSuite) TestReviewAccept() {
	s.LoadFullFixtures()
	defer s.ResetFixtures()

	require := s.Require()
	a := s.svc.GetAdmin()

	var err error
	charlieID := s.fixtures[vasps]["charliebank"].(*pb.VASP).Id
	julietVASP := s.fixtures[vasps]["juliet"].(*pb.VASP)
	xrayID := s.fixtures[certreqs]["xray"].(*models.CertificateRequest).Id

	// VASP does not have an admin verification token
	request := &httpRequest{
		method: http.MethodPost,
		path:   "/v2/vasps/invalid/review",
		in: &admin.ReviewRequest{
			ID:                     charlieID,
			AdminVerificationToken: "supersecrettoken",
			Accept:                 true,
		},
		params: map[string]string{
			"vaspID": charlieID,
		},
		claims: &tokens.Claims{
			Email: "admin@example.com",
		},
	}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.Review, c, w, nil)
	require.Equal(http.StatusUnauthorized, rep.StatusCode)

	// Test incorrect admin verification token
	request.in = &admin.ReviewRequest{
		ID:                     julietVASP.Id,
		AdminVerificationToken: "supersecrettoken",
		Accept:                 true,
	}
	request.params = map[string]string{
		"vaspID": julietVASP.Id,
	}
	actual := &admin.ReviewReply{}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.Review, c, w, actual)
	require.Equal(http.StatusUnauthorized, rep.StatusCode)

	// Successfully accepting a registration request
	reviewRequest := &admin.ReviewRequest{
		ID:     julietVASP.Id,
		Accept: true,
	}
	reviewRequest.AdminVerificationToken, err = models.GetAdminVerificationToken(julietVASP)
	require.NoError(err, "could not get required admin verification token for juliet")

	request.in = reviewRequest
	actual = &admin.ReviewReply{}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.Review, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	require.Equal(pb.VerificationState_REVIEWED.String(), actual.Status)
	require.Contains(actual.Message, "has been approved")

	// VASP state should be changed to REVIEWED
	v, err := s.svc.GetStore().RetrieveVASP(julietVASP.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_REVIEWED, v.VerificationStatus)

	// VASP audit log should contain the new entry
	log, err := models.GetAuditLog(v)
	require.NoError(err)
	require.Len(log, 4)
	require.Equal(pb.VerificationState_SUBMITTED, log[0].CurrentState)
	require.Equal(pb.VerificationState_EMAIL_VERIFIED, log[1].CurrentState)
	require.Equal(pb.VerificationState_PENDING_REVIEW, log[2].CurrentState)
	require.Equal(pb.VerificationState_PENDING_REVIEW, log[3].PreviousState)
	require.Equal(pb.VerificationState_REVIEWED, log[3].CurrentState)
	require.Equal(request.claims.Email, log[3].Source)

	// Certificate request should be changed to READY_TO_SUBMIT
	cert, err := s.svc.GetStore().RetrieveCertReq(xrayID)
	require.NoError(err)
	require.Equal(models.CertificateRequestState_READY_TO_SUBMIT, cert.Status)

	// Certificate request audit log should contain the new entry
	require.Len(cert.AuditLog, 2)
	require.Equal(models.CertificateRequestState_INITIALIZED, cert.AuditLog[0].CurrentState)
	require.Equal(models.CertificateRequestState_INITIALIZED, cert.AuditLog[1].PreviousState)
	require.Equal(models.CertificateRequestState_READY_TO_SUBMIT, cert.AuditLog[1].CurrentState)
	require.Equal(request.claims.Email, cert.AuditLog[1].Source)
}

// Test the Review endpoint for the reject case.
func (s *gdsTestSuite) TestReviewReject() {
	s.LoadFullFixtures()
	defer s.ResetFixtures()
	defer emails.PurgeMockEmails()

	require := s.Require()
	a := s.svc.GetAdmin()

	charlieID := s.fixtures[vasps]["charliebank"].(*pb.VASP).Id
	julietVASP := s.fixtures[vasps]["juliet"].(*pb.VASP)
	xrayID := s.fixtures[certreqs]["xray"].(*models.CertificateRequest).Id

	// Test when VASP does not have admin verification token
	request := &httpRequest{
		method: http.MethodPost,
		path:   "/v2/vasps/invalid/review",
		in: &admin.ReviewRequest{
			ID:                     charlieID,
			AdminVerificationToken: "supersecrettoken",
			Accept:                 false,
			RejectReason:           "some reason",
		},
		params: map[string]string{
			"vaspID": charlieID,
		},
		claims: &tokens.Claims{
			Email: "admin@example.com",
		},
	}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.Review, c, w, nil)
	require.Equal(http.StatusUnauthorized, rep.StatusCode)

	// Incorrect admin verification token
	request.in = &admin.ReviewRequest{
		ID:                     julietVASP.Id,
		AdminVerificationToken: "supersecrettoken",
		Accept:                 false,
		RejectReason:           "just joking around",
	}
	request.params = map[string]string{
		"vaspID": julietVASP.Id,
	}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.Review, c, w, nil)
	require.Equal(http.StatusUnauthorized, rep.StatusCode)

	avt, err := models.GetAdminVerificationToken(julietVASP)
	require.NoError(err, "could not fetch admin verification token for juliet")
	require.NotEmpty(avt, "juliet does not have an admin verification token")

	// No rejection reason supplied
	request.in = &admin.ReviewRequest{
		ID:                     julietVASP.Id,
		AdminVerificationToken: avt,
		Accept:                 false,
	}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.Review, c, w, nil)
	require.Equal(http.StatusBadRequest, rep.StatusCode)

	// Successfully rejecting a registration request
	request.in = &admin.ReviewRequest{
		ID:                     julietVASP.Id,
		AdminVerificationToken: avt,
		Accept:                 false,
		RejectReason:           "some reason",
	}
	actual := &admin.ReviewReply{}
	c, w = s.makeRequest(request)
	sent := time.Now()
	rep = s.doRequest(a.Review, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	require.Equal(pb.VerificationState_REJECTED.String(), actual.Status)
	require.Contains(actual.Message, "has been rejected")

	// VASP state should be changed to REJECTED
	v, err := s.svc.GetStore().RetrieveVASP(julietVASP.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_REJECTED, v.VerificationStatus)

	// VASP audit log should contain the new entry
	log, err := models.GetAuditLog(v)
	require.NoError(err)
	require.Len(log, 4)
	require.Equal(pb.VerificationState_SUBMITTED, log[0].CurrentState)
	require.Equal(pb.VerificationState_EMAIL_VERIFIED, log[1].CurrentState)
	require.Equal(pb.VerificationState_PENDING_REVIEW, log[2].CurrentState)
	require.Equal(pb.VerificationState_PENDING_REVIEW, log[3].PreviousState)
	require.Equal(pb.VerificationState_REJECTED, log[3].CurrentState)
	require.Equal(request.claims.Email, log[3].Source)

	// Certificate request should be deleted
	_, err = s.svc.GetStore().RetrieveCertReq(xrayID)
	require.Error(err)

	emailLog, err := models.GetEmailLog(v.Contacts.Administrative)
	require.NoError(err)
	require.Len(emailLog, 1)

	// Rejection emails should be sent to the verified contacts
	messages := []*emailMeta{
		{
			contact:   v.Contacts.Administrative,
			to:        v.Contacts.Administrative.Email,
			from:      s.svc.GetConf().Email.ServiceEmail,
			subject:   emails.RejectRegistrationRE,
			reason:    string(admin.ResendRejection),
			timestamp: sent,
		},
		{
			contact:   v.Contacts.Legal,
			to:        v.Contacts.Legal.Email,
			from:      s.svc.GetConf().Email.ServiceEmail,
			subject:   emails.RejectRegistrationRE,
			reason:    string(admin.ResendRejection),
			timestamp: sent,
		},
	}
	s.CheckEmails(messages)
}

// Test the Resend endpoint.
func (s *gdsTestSuite) TestResend() {
	s.LoadFullFixtures()
	defer s.ResetFixtures()
	defer emails.PurgeMockEmails()

	require := s.Require()
	a := s.svc.GetAdmin()

	vaspErrored := s.fixtures[vasps]["golfbucks"].(*pb.VASP).Id
	vaspRejected := s.fixtures[vasps]["lima"].(*pb.VASP).Id

	// Supplying an invalid VASP ID
	request := &httpRequest{
		method: http.MethodPost,
		path:   "/v2/vasps/invalid/resend",
		in: &admin.ResendRequest{
			ID:     "invalid",
			Action: admin.ResendVerifyContact,
		},
		params: map[string]string{
			"vaspID": "invalid",
		},
	}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.Resend, c, w, nil)
	require.Equal(http.StatusNotFound, rep.StatusCode)

	// ResendVerifyContact email
	request.in = &admin.ResendRequest{
		ID:     vaspErrored,
		Action: admin.ResendVerifyContact,
		Reason: "verify",
	}
	request.params = map[string]string{
		"vaspID": vaspErrored,
	}
	actual := &admin.ResendReply{}
	c, w = s.makeRequest(request)
	firstSend := time.Now()
	rep = s.doRequest(a.Resend, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	require.Equal(1, actual.Sent)
	require.Contains(actual.Message, "contact verification emails resent")

	// ResendReview email
	request.in = &admin.ResendRequest{
		ID:     vaspErrored,
		Action: admin.ResendReview,
		Reason: "review",
	}
	request.params = map[string]string{
		"vaspID": vaspErrored,
	}
	actual = &admin.ResendReply{}
	c, w = s.makeRequest(request)
	secondSend := time.Now()
	rep = s.doRequest(a.Resend, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	require.Equal(1, actual.Sent)
	require.Contains(actual.Message, "review request resent")

	// ResendRejection email
	request.in = &admin.ResendRequest{
		ID:     vaspRejected,
		Action: admin.ResendRejection,
		Reason: "reject",
	}
	request.params = map[string]string{
		"vaspID": vaspRejected,
	}
	actual = &admin.ResendReply{}
	c, w = s.makeRequest(request)
	thirdSend := time.Now()
	rep = s.doRequest(a.Resend, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	require.Equal(2, actual.Sent)
	require.Contains(actual.Message, "rejection emails resent")

	// Verify that all emails were sent
	errored, err := s.svc.GetStore().RetrieveVASP(vaspErrored)
	require.NoError(err)
	rejected, err := s.svc.GetStore().RetrieveVASP(vaspRejected)
	require.NoError(err)

	messages := []*emailMeta{
		{
			contact:   errored.Contacts.Billing,
			to:        errored.Contacts.Billing.Email,
			from:      s.svc.GetConf().Email.ServiceEmail,
			subject:   emails.VerifyContactRE,
			reason:    "verify_contact",
			timestamp: firstSend,
		},
		{
			to:        s.svc.GetConf().Email.AdminEmail,
			from:      s.svc.GetConf().Email.ServiceEmail,
			subject:   emails.ReviewRequestRE,
			timestamp: secondSend,
		},
		{
			contact:   rejected.Contacts.Administrative,
			to:        rejected.Contacts.Administrative.Email,
			from:      s.svc.GetConf().Email.ServiceEmail,
			subject:   emails.RejectRegistrationRE,
			reason:    "rejection",
			timestamp: thirdSend,
		},
		{
			contact:   rejected.Contacts.Legal,
			to:        rejected.Contacts.Legal.Email,
			from:      s.svc.GetConf().Email.ServiceEmail,
			subject:   emails.RejectRegistrationRE,
			reason:    "rejection",
			timestamp: thirdSend,
		},
	}
	s.CheckEmails(messages)
}

// Test the ReviewTimeline endpoint.
func (s *gdsTestSuite) TestReviewTimeline() {
	s.LoadSmallFixtures()
	require := s.Require()
	a := s.svc.GetAdmin()

	// Invalid start date
	request := &httpRequest{
		method: http.MethodGet,
		path:   "/v2/reviews?start=09-01-2021",
	}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.ReviewTimeline, c, w, nil)
	require.Equal(http.StatusBadRequest, rep.StatusCode)

	// Start date is before epoch
	request.path = "/v2/reviews?start=1968-01-01"
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ReviewTimeline, c, w, nil)
	require.Equal(http.StatusBadRequest, rep.StatusCode)

	// Invalid end date
	request.path = "/v2/reviews?end=09-01-2021"
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ReviewTimeline, c, w, nil)
	require.Equal(http.StatusBadRequest, rep.StatusCode)

	// Start date is after end date
	request.path = "/v2/reviews?start=2021-01-01&end=2020-01-01"
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ReviewTimeline, c, w, nil)
	require.Equal(http.StatusBadRequest, rep.StatusCode)

	// Successful retrieval of review timeline
	request.path = "/v2/reviews?start=2021-08-23&end=2021-09-01"
	actual := &admin.ReviewTimelineReply{}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ReviewTimeline, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	expected := &admin.ReviewTimelineReply{
		Weeks: []admin.ReviewTimelineRecord{
			{
				Week:         "2021-08-23",
				VASPsUpdated: 1,
				Registrations: map[string]int{
					pb.VerificationState_NO_VERIFICATION.String():     0,
					pb.VerificationState_SUBMITTED.String():           0,
					pb.VerificationState_EMAIL_VERIFIED.String():      0,
					pb.VerificationState_PENDING_REVIEW.String():      1,
					pb.VerificationState_REVIEWED.String():            0,
					pb.VerificationState_ISSUING_CERTIFICATE.String(): 0,
					pb.VerificationState_VERIFIED.String():            0,
					pb.VerificationState_REJECTED.String():            0,
					pb.VerificationState_APPEALED.String():            0,
					pb.VerificationState_ERRORED.String():             0,
				},
			},
			{
				Week:         "2021-08-30",
				VASPsUpdated: 2,
				Registrations: map[string]int{
					pb.VerificationState_NO_VERIFICATION.String():     0,
					pb.VerificationState_SUBMITTED.String():           1,
					pb.VerificationState_EMAIL_VERIFIED.String():      0,
					pb.VerificationState_PENDING_REVIEW.String():      0,
					pb.VerificationState_REVIEWED.String():            0,
					pb.VerificationState_ISSUING_CERTIFICATE.String(): 0,
					pb.VerificationState_VERIFIED.String():            0,
					pb.VerificationState_REJECTED.String():            1,
					pb.VerificationState_APPEALED.String():            0,
					pb.VerificationState_ERRORED.String():             0,
				},
			},
		},
	}
	require.Equal(expected, actual)
}
