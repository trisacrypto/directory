package gds_test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/trisacrypto/directory/pkg/gds"
	admin "github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
	"github.com/trisacrypto/directory/pkg/utils"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

func (s *gdsTestSuite) initAdmin(dbPath string) (admin *gds.Admin) {
	var err error
	require := s.Require()
	s.conf.Database = config.DatabaseConfig{
		URL: "leveldb:///" + dbPath,
	}
	admin, s.db, s.tokens, err = gds.NewMockedAdmin(*s.conf)
	require.NoError(err)
	return admin
}

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
	if request.params != nil {
		for k, v := range request.params {
			c.Params = append(c.Params, gin.Param{
				Key:   k,
				Value: v,
			})
		}
	}
	return c, w
}

// doRequest is a helper function for making an admin API request and retrieving
// the response.
func (s *gdsTestSuite) doRequest(fn gin.HandlerFunc, c *gin.Context, w *httptest.ResponseRecorder, reply interface{}) (res *http.Response) {
	require := s.Require()
	// Call the admin function and return the HTTP response
	fn(c)
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

// Test that we get a good response from ProtectAuthenticate.
func (s *gdsTestSuite) TestProtectAuthenticate() {
	require := s.Require()
	a := s.initAdmin(s.dbPath)

	request := &httpRequest{
		method: http.MethodPost,
		path:   "/v2/foo",
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
		require.Equal(s.conf.Admin.CookieDomain, cookie.Domain)
	}
}

// Test the Authenticate endpoint.
func (s *gdsTestSuite) TestAuthenticate() {
	require := s.Require()
	a := s.initAdmin(s.dbPath)

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

	// TODO: Test successful authentication path
}

// Test the Reauthenticate endpoint.
func (s *gdsTestSuite) TestReauthenticate() {
	require := s.Require()
	a := s.initAdmin(s.dbPath)

	claims := &tokens.Claims{
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.NewString(),
			Audience:  "http://localhost",
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	refreshToken, err := s.tokens.CreateRefreshToken(accessToken)
	require.NoError(err)
	access, err := s.tokens.Sign(accessToken)
	require.NoError(err)
	refresh, err := s.tokens.Sign(refreshToken)
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
	claims.Id = uuid.NewString()
	otherToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	other, err := s.tokens.Sign(otherToken)
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
		require.Equal(s.conf.Admin.CookieDomain, cookie.Domain)
	}
}

// Test that the Summary endpoint returns the correct response.
func (s *gdsTestSuite) TestSummary() {
	require := s.Require()
	a := s.initAdmin(s.dbPath)

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
		PendingRegistrations: 5,
		ContactsCount:        40,
		VerifiedContacts:     28,
		CertificatesIssued:   4,
		Statuses: map[string]int{
			pb.VerificationState_APPEALED.String():       1,
			pb.VerificationState_ERRORED.String():        1,
			pb.VerificationState_PENDING_REVIEW.String(): 2,
			pb.VerificationState_REJECTED.String():       2,
			pb.VerificationState_SUBMITTED.String():      2,
			pb.VerificationState_VERIFIED.String():       6,
		},
		CertReqs: map[string]int{
			models.CertificateRequestState_COMPLETED.String():   4,
			models.CertificateRequestState_CR_ERRORED.String():  1,
			models.CertificateRequestState_CR_REJECTED.String(): 2,
			models.CertificateRequestState_INITIALIZED.String(): 3,
		},
	}
	require.Equal(expected, actual)
}

// Test that the Autocomplete endpoint returns the correct response.
func (s *gdsTestSuite) TestAutocomplete() {
	require := s.Require()
	a := s.initAdmin(s.smallDBPath)

	request := &httpRequest{
		method: http.MethodGet,
		path:   "/v2/autocomplete",
	}
	actual := &admin.AutocompleteReply{}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.Autocomplete, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)

	// Construct the expected response
	expected := &admin.AutocompleteReply{
		Names: map[string]string{
			"trisa.charliebank.io":         getVASPIDFromKey(smallDBFixtures[0]),
			"https://trisa.charliebank.io": "https://trisa.charliebank.io",
			"CharlieBank":                  getVASPIDFromKey(smallDBFixtures[0]),
			"trisa.delta.io":               getVASPIDFromKey(smallDBFixtures[1]),
			"https://trisa.delta.io":       "https://trisa.delta.io",
			"Delta Assets":                 getVASPIDFromKey(smallDBFixtures[1]),
		},
	}
	require.Equal(expected, actual)
}

// Test the ListVASPs endpoint.
func (s *gdsTestSuite) TestListVASPs() {
	require := s.Require()
	a := s.initAdmin(s.smallDBPath)

	snippets := []admin.VASPSnippet{
		{
			ID:                  getVASPIDFromKey(smallDBFixtures[0]),
			Name:                "CharlieBank",
			CommonName:          "trisa.charliebank.io",
			RegisteredDirectory: "trisatest.net",
			VerificationStatus:  pb.VerificationState_SUBMITTED.String(),
			LastUpdated:         "2021-09-27T04:12:23Z",
			VerifiedContacts: map[string]bool{
				"administrative": false,
				"billing":        true,
				"legal":          true,
				"technical":      false,
			},
		},
		{
			ID:                  getVASPIDFromKey(smallDBFixtures[1]),
			Name:                "Delta Assets",
			CommonName:          "trisa.delta.io",
			RegisteredDirectory: "trisatest.net",
			VerificationStatus:  pb.VerificationState_APPEALED.String(),
			LastUpdated:         "2021-09-18T10:58:22Z",
			VerifiedContacts: map[string]bool{
				"administrative": false,
				"billing":        true,
				"legal":          true,
				"technical":      false,
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
	sort.Slice(actual.VASPs, func(i, j int) bool {
		return actual.VASPs[i].ID < actual.VASPs[j].ID
	})
	require.Equal(snippets, actual.VASPs)

	// List VASPs with an invalid status
	request.path = "/v2/vasps?status=invalid"
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ListVASPs, c, w, nil)
	require.Equal(http.StatusBadRequest, rep.StatusCode)

	// List VASPs with the specified status
	request.path = "/v2/vasps?status=" + snippets[0].VerificationStatus
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
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ListVASPs, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	require.Equal(2, actual.Count)
	require.Equal(1, actual.Page)
	require.Equal(100, actual.PageSize)
	require.Len(actual.VASPs, 2)
	require.Equal(snippets[0], actual.VASPs[0])
	require.Equal(snippets[1], actual.VASPs[1])

	// List VASPs on multiple pages
	request.path = "/v2/vasps?page=1&page_size=1"
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ListVASPs, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	require.Equal(len(snippets), actual.Count)
	require.Equal(1, actual.Page)
	require.Equal(1, actual.PageSize)
	require.Len(actual.VASPs, 1)
	require.Equal(snippets[0], actual.VASPs[0])

	request.path = "/v2/vasps?page=2&page_size=1"
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ListVASPs, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	require.Equal(len(snippets), actual.Count)
	require.Equal(2, actual.Page)
	require.Equal(1, actual.PageSize)
	require.Len(actual.VASPs, 1)
	require.Equal(snippets[1], actual.VASPs[0])

	request.path = "/v2/vasps?page=3&page_size=1"
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
	require := s.Require()
	a := s.initAdmin(s.smallDBPath)

	// Retrieve a VASP that doesn't exist
	request := &httpRequest{
		method: http.MethodGet,
		path:   "/v2/vasps/invalid",
	}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.RetrieveVASP, c, w, nil)
	require.Equal(http.StatusNotFound, rep.StatusCode)

	// Retrieve a VASP that exists
	request.path = "/v2/vasps/" + getVASPIDFromKey(smallDBFixtures[0])
	request.params = map[string]string{
		"vaspID": getVASPIDFromKey(smallDBFixtures[0]),
	}
	actual := &admin.RetrieveVASPReply{}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.RetrieveVASP, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	expected := &admin.RetrieveVASPReply{
		Name: "CharlieBank",
		VerifiedContacts: map[string]string{
			"billing": "anthony@charliebank.net",
			"legal":   "sonia@charliebank.com",
		},
		Traveler: false,
		AuditLog: []map[string]interface{}{
			{
				"current_state":  pb.VerificationState_SUBMITTED.String(),
				"description":    "register request received",
				"previous_state": pb.VerificationState_NO_VERIFICATION.String(),
				"source":         "automated",
				"timestamp":      "2021-09-27T04:12:23Z",
			},
		},
	}
	// RetrieveVASP removes the extra data from the VASP before returning it
	obj, ok := s.fixtures[smallDBFixtures[0]]
	require.True(ok)
	v := obj.(*pb.VASP)
	v.Extra = nil
	v.Contacts.Administrative.Extra = nil
	v.Contacts.Legal.Extra = nil
	v.Contacts.Technical.Extra = nil
	v.Contacts.Billing.Extra = nil
	var err error
	expected.VASP, err = utils.Rewire(v)
	require.NoError(err)
	require.Equal(expected, actual)
}

// Test the CreateReviewNote endpoint.
func (s *gdsTestSuite) TestCreateReviewNote() {
	require := s.Require()
	a := s.initAdmin(s.dbPath)

	vasps := []string{
		"d9da630e-41aa-11ec-9d29-acde48001122",
	}

	// Supplying an invalid note ID
	request := &httpRequest{
		method: http.MethodPost,
		path:   "/v2/vasps/" + vasps[0] + "/notes",
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
		VASP:   vasps[0],
		NoteID: "89bceb0e-41aa-11ec-9d29-acde48001122",
		Text:   "foo",
	}
	request.params = map[string]string{
		"vaspID": vasps[0],
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
	v, err := s.db.RetrieveVASP(vasps[0])
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
	require := s.Require()
	a := s.initAdmin(s.dbPath)

	vasps := []string{
		"d9da630e-41aa-11ec-9d29-acde48001122",
	}

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
		"vaspID": vasps[0],
	}
	actual := &admin.ListReviewNotesReply{}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ListReviewNotes, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	expected := &admin.ListReviewNotesReply{
		Notes: []admin.ReviewNote{
			{
				Author:   "admin@trisa.io",
				Created:  "2021-07-26T11:25:44Z",
				Editor:   "juanreyes@example.net",
				ID:       "d9ebcfa4-41aa-11ec-9d29-acde48001122",
				Modified: "2021-08-30T22:40:24Z",
				Text:     "Porro magnam amet ut.",
			},
		},
	}
	require.Equal(expected, actual)
}

// Test the UpdateReviewNote endpoint.
func (s *gdsTestSuite) TestUpdateReviewNote() {
	require := s.Require()
	a := s.initAdmin(s.dbPath)

	vasps := []string{
		"d9da630e-41aa-11ec-9d29-acde48001122",
	}

	noteID := "d9ebcfa4-41aa-11ec-9d29-acde48001122"

	// Supplying an invalid note ID
	request := &httpRequest{
		method: http.MethodPut,
		path:   "/v2/vasps/" + vasps[0] + "/notes/invalid",
		in: &admin.ModifyReviewNoteRequest{
			VASP:   vasps[0],
			NoteID: "invalid slug",
		},
		params: map[string]string{
			"vaspID": vasps[0],
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
		VASP:   vasps[0],
		NoteID: noteID,
		Text:   "bar",
	}
	request.params = map[string]string{
		"vaspID": vasps[0],
		"noteID": noteID,
	}
	actual := &admin.ReviewNote{}
	modified := time.Now()
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.UpdateReviewNote, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	// Validate returned note
	require.Equal(noteID, actual.ID)
	require.Equal("2021-07-26T11:25:44Z", actual.Created)
	ts, err := time.Parse(time.RFC3339, actual.Modified)
	require.NoError(err)
	require.True(ts.Sub(modified) < time.Minute)
	require.Equal("admin@trisa.io", actual.Author)
	require.Equal(request.claims.Email, actual.Editor)
	require.Equal("bar", actual.Text)
	// Record on the database should be updated
	v, err := s.db.RetrieveVASP(vasps[0])
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
	require := s.Require()
	a := s.initAdmin(s.dbPath)

	vasps := []string{
		"d9da630e-41aa-11ec-9d29-acde48001122",
	}

	noteID := "d9ebcfa4-41aa-11ec-9d29-acde48001122"

	// Supplying an invalid note ID
	request := &httpRequest{
		method: http.MethodDelete,
		path:   "/v2/vasps/" + vasps[0] + "/notes/invalid",
		params: map[string]string{
			"vaspID": vasps[0],
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
		"vaspID": vasps[0],
		"noteID": noteID,
	}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.DeleteReviewNote, c, w, nil)
	require.Equal(http.StatusOK, rep.StatusCode)
	// Record on the database should be deleted
	v, err := s.db.RetrieveVASP(vasps[0])
	require.NoError(err)
	notes, err := models.GetReviewNotes(v)
	require.NoError(err)
	require.Len(notes, 0)
}

// Test the Review endpoint.
func (s *gdsTestSuite) TestReview() {
	require := s.Require()
	a := s.initAdmin(s.dbPath)

	vasps := []string{
		"d9da630e-41aa-11ec-9d29-acde48001122",
	}

	// Supplying an invalid VASP ID
	request := &httpRequest{
		method: http.MethodPost,
		path:   "/v2/vasps/invalid/review",
		in: &admin.ReviewRequest{
			ID:                     "invalid",
			AdminVerificationToken: "foo",
			Accept:                 true,
		},
		params: map[string]string{
			"vaspID": "invalid",
		},
	}
	c, w := s.makeRequest(request)
	rep := s.doRequest(a.Review, c, w, nil)
	require.Equal(http.StatusNotFound, rep.StatusCode)

	// No verification token supplied
	request.in = &admin.ReviewRequest{
		ID:     vasps[0],
		Accept: true,
	}
	request.params = map[string]string{
		"vaspID": vasps[0],
	}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.Review, c, w, nil)
	require.Equal(http.StatusBadRequest, rep.StatusCode)

	// No rejection reason supplied
	request.in = &admin.ReviewRequest{
		ID:                     vasps[0],
		AdminVerificationToken: "foo",
		Accept:                 false,
	}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.Review, c, w, nil)
	require.Equal(http.StatusBadRequest, rep.StatusCode)

	// TODO: Add tests for accept and reject paths. This requires VASP fixtures with
	// a non-empty adminVerificationToken.
}

// Test the Resend endpoint.
func (s *gdsTestSuite) TestResend() {
	require := s.Require()
	a := s.initAdmin(s.dbPath)

	vaspErrored := "da2b165a-41aa-11ec-9d29-acde48001122"
	vaspRejected := "da8bd0e4-41aa-11ec-9d29-acde48001122"

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
	sent := time.Now()
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.Resend, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	require.Equal(1, actual.Sent)
	require.Contains(actual.Message, "contact verification emails resent")
	// Email audit log should be updated
	v, err := s.db.RetrieveVASP(vaspErrored)
	require.NoError(err)
	emails, err := models.GetEmailLog(v.Contacts.Billing)
	require.NoError(err)
	require.Len(emails, 1)
	ts, err := time.Parse(time.RFC3339, emails[0].Timestamp)
	require.NoError(err)
	require.True(ts.Sub(sent) < time.Minute)
	require.Equal("verify_contact", emails[0].Reason)

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
	sent = time.Now()
	c, w = s.makeRequest(request)
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
	sent = time.Now()
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.Resend, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	require.Equal(2, actual.Sent)
	require.Contains(actual.Message, "rejection emails resent")
	// Email audit logs should be updated
	v, err = s.db.RetrieveVASP(vaspRejected)
	require.NoError(err)
	emails, err = models.GetEmailLog(v.Contacts.Administrative)
	require.NoError(err)
	require.Len(emails, 1)
	ts, err = time.Parse(time.RFC3339, emails[0].Timestamp)
	require.NoError(err)
	require.True(ts.Sub(sent) < time.Minute)
	require.Equal("rejection", emails[0].Reason)
	emails, err = models.GetEmailLog(v.Contacts.Legal)
	require.NoError(err)
	require.Len(emails, 1)
	ts, err = time.Parse(time.RFC3339, emails[0].Timestamp)
	require.NoError(err)
	require.True(ts.Sub(sent) < time.Minute)
	require.Equal("rejection", emails[0].Reason)
}

// Test the ReviewTimeline endpoint.
func (s *gdsTestSuite) TestReviewTimeline() {
	require := s.Require()
	a := s.initAdmin(s.smallDBPath)

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
	request.path = "/v2/reviews?start=2021-09-20&end=2021-09-30"
	actual := &admin.ReviewTimelineReply{}
	c, w = s.makeRequest(request)
	rep = s.doRequest(a.ReviewTimeline, c, w, actual)
	require.Equal(http.StatusOK, rep.StatusCode)
	expected := &admin.ReviewTimelineReply{
		Weeks: []admin.ReviewTimelineRecord{
			{
				Week:         "2021-09-20",
				VASPsUpdated: 0,
				Registrations: map[string]int{
					pb.VerificationState_NO_VERIFICATION.String():     0,
					pb.VerificationState_SUBMITTED.String():           0,
					pb.VerificationState_EMAIL_VERIFIED.String():      0,
					pb.VerificationState_PENDING_REVIEW.String():      0,
					pb.VerificationState_REVIEWED.String():            0,
					pb.VerificationState_ISSUING_CERTIFICATE.String(): 0,
					pb.VerificationState_VERIFIED.String():            0,
					pb.VerificationState_REJECTED.String():            0,
					pb.VerificationState_APPEALED.String():            0,
					pb.VerificationState_ERRORED.String():             0,
				},
			},
			{
				Week:         "2021-09-27",
				VASPsUpdated: 1,
				Registrations: map[string]int{
					pb.VerificationState_NO_VERIFICATION.String():     0,
					pb.VerificationState_SUBMITTED.String():           1,
					pb.VerificationState_EMAIL_VERIFIED.String():      0,
					pb.VerificationState_PENDING_REVIEW.String():      0,
					pb.VerificationState_REVIEWED.String():            0,
					pb.VerificationState_ISSUING_CERTIFICATE.String(): 0,
					pb.VerificationState_VERIFIED.String():            0,
					pb.VerificationState_REJECTED.String():            0,
					pb.VerificationState_APPEALED.String():            0,
					pb.VerificationState_ERRORED.String():             0,
				},
			},
		},
	}
	require.Equal(expected, actual)
}
