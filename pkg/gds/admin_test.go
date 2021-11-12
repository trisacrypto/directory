package gds

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	admin "github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/emails"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/store"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
	"github.com/trisacrypto/directory/pkg/utils"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

func initAdmin(s *gdsTestSuite) (admin *Admin) {
	require := s.Require()
	db, err := store.Open(config.DatabaseConfig{
		URL: "leveldb:///" + s.db,
	})
	require.NoError(err)
	manager, err := emails.New(config.EmailConfig{
		ServiceEmail: "service@example.com",
		AdminEmail:   "admin@example.com",
		Testing:      true,
	})
	require.NoError(err)
	return &Admin{
		svc: &Service{
			email: manager,
		},
		db:   db,
		conf: &config.AdminConfig{},
	}
}

// generateToken generates a fake JWT token to send back to the Sectigo client.
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

// doAdminRequest is a helper function for making an admin API request and retrieving
// the response.
func (s *gdsTestSuite) doRequest(fn func(c *gin.Context), method, path string, headers, params map[string]string, request interface{}, reply interface{}, claims *tokens.Claims) (status int) {
	var body io.ReadWriter
	var err error
	require := s.Require()

	// Encode the JSON request
	if request != nil {
		body = &bytes.Buffer{}
		err = json.NewEncoder(body).Encode(request)
		require.NoError(err)
	} else {
		body = nil
	}

	// Construct the HTTP request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	if claims != nil {
		c.Set(admin.UserClaims, claims)
	}
	c.Request = httptest.NewRequest(method, path, body)
	c.Request.Header.Add("Content-Type", "application/json")
	if headers != nil {
		for k, v := range headers {
			c.Request.Header.Add(k, v)
		}
	}
	if params != nil {
		for k, v := range params {
			c.Params = append(c.Params, gin.Param{
				Key:   k,
				Value: v,
			})
		}
	}

	// Call the admin function and return the response
	fn(c)
	res := w.Result()
	defer res.Body.Close()
	if reply != nil {
		bytes, err := ioutil.ReadAll(res.Body)
		require.NoError(err)
		err = json.Unmarshal(bytes, reply)
		require.NoError(err)
	}
	return res.StatusCode
}

// Test that we get a good response from ProtectAuthenticate.
func (s *gdsTestSuite) TestProtectAuthenticate() {
	require := s.Require()
	a := initAdmin(s)

	actual := &admin.Reply{}
	status := s.doRequest(a.ProtectAuthenticate, http.MethodPost, "/v2/protect/authenticate", nil, nil, nil, actual, nil)
	require.Equal(http.StatusOK, status)
	expected := &admin.Reply{Success: true}
	require.Equal(expected, actual)
}

// Test the Authenticate endpoint.
func (s *gdsTestSuite) TestAuthenticate() {
	require := s.Require()
	a := initAdmin(s)

	// Missing credential
	req := &admin.AuthRequest{}
	status := s.doRequest(a.Authenticate, http.MethodPost, "/v2/authenticate", nil, nil, req, nil, nil)
	require.Equal(http.StatusUnauthorized, status)

	// Invalid credential
	req = &admin.AuthRequest{
		Credential: "invalid",
	}
	status = s.doRequest(a.Authenticate, http.MethodPost, "/v2/authenticate", nil, nil, req, nil, nil)
	require.Equal(http.StatusUnauthorized, status)

	// TODO: Test the successful authentication path
}

// Test the Reauthenticate endpoint.
func (s *gdsTestSuite) TestReauthenticate() {
	require := s.Require()
	a := initAdmin(s)

	// Missing access token
	req := &admin.AuthRequest{}
	status := s.doRequest(a.Reauthenticate, http.MethodPost, "/v2/reauthenticate", nil, nil, req, nil, nil)
	require.Equal(http.StatusUnauthorized, status)

	// Invalid access token
	headers := map[string]string{
		"Authorization": "Bearer invalid",
	}
	status = s.doRequest(a.Reauthenticate, http.MethodPost, "/v2/reauthenticate", headers, nil, req, nil, nil)
	require.Equal(http.StatusUnauthorized, status)

	// TODO: Mock token manager to test the successful reauthentication path
}

// Test that the Summary endpoint returns the correct response.
func (s *gdsTestSuite) TestSummary() {
	require := s.Require()
	a := initAdmin(s)

	actual := &admin.SummaryReply{}
	status := s.doRequest(a.Summary, http.MethodGet, "/v2/summary", nil, nil, nil, actual, nil)
	require.Equal(http.StatusOK, status)

	// Test against the expected response
	expected := &admin.SummaryReply{
		VASPsCount:           len(s.dbVASPs),
		PendingRegistrations: 5,
		ContactsCount:        40,
		VerifiedContacts:     28,
		CertificatesIssued:   0,
		Statuses: map[string]int{
			pb.VerificationState_APPEALED.String():       1,
			pb.VerificationState_ERRORED.String():        1,
			pb.VerificationState_PENDING_REVIEW.String(): 2,
			pb.VerificationState_REJECTED.String():       2,
			pb.VerificationState_SUBMITTED.String():      2,
			pb.VerificationState_VERIFIED.String():       6,
		},
		CertReqs: map[string]int{},
	}
	require.Equal(expected, actual)
}

// Test that the Autocomplete endpoint returns the correct response.
func (s *gdsTestSuite) TestAutocomplete() {
	require := s.Require()
	vasps := []string{
		"d9da630e-41aa-11ec-9d29-acde48001122",
		"d9efca14-41aa-11ec-9d29-acde48001122",
	}
	WriteVASPs(s, vasps)
	a := initAdmin(s)

	actual := &admin.AutocompleteReply{}
	status := s.doRequest(a.Autocomplete, http.MethodGet, "/v2/autocomplete", nil, nil, nil, actual, nil)
	require.Equal(http.StatusOK, status)

	// Construct the expected response
	expected := &admin.AutocompleteReply{
		Names: map[string]string{
			"trisa.charliebank.io":         vasps[0],
			"https://trisa.charliebank.io": "https://trisa.charliebank.io",
			"CharlieBank":                  vasps[0],
			"trisa.delta.io":               vasps[1],
			"https://trisa.delta.io":       "https://trisa.delta.io",
			"Delta Assets":                 vasps[1],
		},
	}
	require.Equal(expected, actual)
}

// Test the ListVASPs endpoint.
func (s *gdsTestSuite) TestListVASPs() {
	require := s.Require()
	vasps := []string{
		"d9da630e-41aa-11ec-9d29-acde48001122",
		"d9efca14-41aa-11ec-9d29-acde48001122",
	}
	WriteVASPs(s, vasps)
	a := initAdmin(s)

	snippets := []admin.VASPSnippet{
		{
			ID:                  vasps[0],
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
			ID:                  vasps[1],
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
	actual := &admin.ListVASPsReply{}
	status := s.doRequest(a.ListVASPs, http.MethodGet, "/v2/vasps", nil, nil, nil, actual, nil)
	require.Equal(http.StatusOK, status)
	require.Equal(len(vasps), actual.Count)
	require.Equal(1, actual.Page)
	require.Equal(100, actual.PageSize)
	require.Len(actual.VASPs, len(snippets))
	sort.Slice(actual.VASPs, func(i, j int) bool {
		return actual.VASPs[i].ID < actual.VASPs[j].ID
	})
	require.Equal(snippets, actual.VASPs)

	// List VASPs with an invalid status
	status = s.doRequest(a.ListVASPs, http.MethodGet, "/v2/vasps?status=invalid", nil, nil, nil, nil, nil)
	require.Equal(http.StatusBadRequest, status)

	// List VASPs with the specified status
	status = s.doRequest(a.ListVASPs, http.MethodGet, "/v2/vasps?status="+snippets[0].VerificationStatus, nil, nil, nil, actual, nil)
	require.Equal(http.StatusOK, status)
	require.Equal(len(vasps), actual.Count)
	require.Equal(1, actual.Page)
	require.Equal(100, actual.PageSize)
	require.Len(actual.VASPs, 1)
	require.Equal(snippets[0], actual.VASPs[0])

	// List VASPs on multiple pages
	status = s.doRequest(a.ListVASPs, http.MethodGet, "/v2/vasps?page=1&page_size=1", nil, nil, nil, actual, nil)
	require.Equal(http.StatusOK, status)
	require.Equal(len(vasps), actual.Count)
	require.Equal(1, actual.Page)
	require.Equal(1, actual.PageSize)
	require.Len(actual.VASPs, 1)
	require.Equal(snippets[0], actual.VASPs[0])

	status = s.doRequest(a.ListVASPs, http.MethodGet, "/v2/vasps?page=2&page_size=1", nil, nil, nil, actual, nil)
	require.Equal(http.StatusOK, status)
	require.Equal(len(vasps), actual.Count)
	require.Equal(2, actual.Page)
	require.Equal(1, actual.PageSize)
	require.Len(actual.VASPs, 1)
	require.Equal(snippets[1], actual.VASPs[0])

	status = s.doRequest(a.ListVASPs, http.MethodGet, "/v2/vasps?page=3&page_size=1", nil, nil, nil, actual, nil)
	require.Equal(http.StatusOK, status)
	require.Equal(len(vasps), actual.Count)
	require.Equal(3, actual.Page)
	require.Equal(1, actual.PageSize)
	require.Len(actual.VASPs, 0)
}

// Test the RetrieveVASP endpoint.
func (s *gdsTestSuite) TestRetrieveVASP() {
	require := s.Require()
	vasps := []string{
		"d9da630e-41aa-11ec-9d29-acde48001122",
	}
	WriteVASPs(s, vasps)
	a := initAdmin(s)

	// Retrieve a VASP that doesn't exist
	status := s.doRequest(a.RetrieveVASP, http.MethodGet, "/v2/vasps/invalid", nil, nil, nil, nil, nil)
	require.Equal(http.StatusNotFound, status)

	// Retrieve a VASP that exists
	actual := &admin.RetrieveVASPReply{}
	params := map[string]string{
		"vaspID": vasps[0],
	}
	status = s.doRequest(a.RetrieveVASP, http.MethodGet, "/v2/vasps/"+vasps[0], nil, params, nil, actual, nil)
	require.Equal(http.StatusOK, status)
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
	v := *s.dbVASPs[vasps[0]]
	v.Extra = nil
	v.Contacts.Administrative.Extra = nil
	v.Contacts.Legal.Extra = nil
	v.Contacts.Technical.Extra = nil
	v.Contacts.Billing.Extra = nil
	var err error
	expected.VASP, err = utils.Rewire(&v)
	require.NoError(err)
	require.Equal(expected, actual)
}

// Test the CreateReviewNote endpoint.
func (s *gdsTestSuite) TestCreateReviewNote() {
	require := s.Require()
	a := initAdmin(s)

	vasps := []string{
		"d9da630e-41aa-11ec-9d29-acde48001122",
	}

	// Supplying an invalid note ID
	req := &admin.ModifyReviewNoteRequest{
		NoteID: "invalid slug",
	}
	claims := &tokens.Claims{
		Email: "admin@example.com",
	}
	status := s.doRequest(a.CreateReviewNote, http.MethodPost, "/v2/vasps/"+vasps[0]+"/notes", nil, nil, req, nil, claims)
	require.Equal(http.StatusBadRequest, status)

	// Supplying an invalid VASP ID
	req = &admin.ModifyReviewNoteRequest{
		VASP: "invalid",
	}
	params := map[string]string{
		"vaspID": "invalid",
	}
	status = s.doRequest(a.CreateReviewNote, http.MethodPost, "/v2/vasps/invalid/notes", nil, params, req, nil, claims)
	require.Equal(http.StatusNotFound, status)

	// Successfully creating a review note
	req = &admin.ModifyReviewNoteRequest{
		VASP:   vasps[0],
		NoteID: "89bceb0e-41aa-11ec-9d29-acde48001122",
		Text:   "foo",
	}
	params = map[string]string{
		"vaspID": vasps[0],
	}
	actual := &admin.ReviewNote{}
	created := time.Now()
	status = s.doRequest(a.CreateReviewNote, http.MethodPost, "/v2/vasps/"+vasps[0]+"/notes", nil, params, req, actual, claims)
	require.Equal(http.StatusCreated, status)
	// Validate returned note
	require.Equal(req.NoteID, actual.ID)
	ts, err := time.Parse(time.RFC3339, actual.Created)
	require.NoError(err)
	require.True(ts.Sub(created) < time.Minute)
	require.Empty(actual.Modified)
	require.Equal(claims.Email, actual.Author)
	require.Empty(actual.Editor)
	require.Equal(req.Text, actual.Text)
	// Record on the database should be updated
	v, err := a.db.RetrieveVASP(vasps[0])
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
	status = s.doRequest(a.CreateReviewNote, http.MethodPost, "/v2/vasps/"+vasps[0]+"/notes", nil, params, req, nil, claims)
	require.Equal(http.StatusBadRequest, status)
}

// Test the ListReviewNotes endpoint.
func (s *gdsTestSuite) TestListReviewNotes() {
	require := s.Require()
	a := initAdmin(s)

	vasps := []string{
		"d9da630e-41aa-11ec-9d29-acde48001122",
	}

	// Supplying an invalid VASP ID
	params := map[string]string{
		"vaspID": "invalid",
	}
	status := s.doRequest(a.ListReviewNotes, http.MethodGet, "/v2/vasps/invalid/notes", nil, params, nil, nil, nil)
	require.Equal(http.StatusNotFound, status)

	// Successfully listing review notes
	params = map[string]string{
		"vaspID": vasps[0],
	}
	actual := &admin.ListReviewNotesReply{}
	status = s.doRequest(a.ListReviewNotes, http.MethodGet, "/v2/vasps/"+vasps[0]+"/notes", nil, params, nil, actual, nil)
	require.Equal(http.StatusOK, status)
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
	a := initAdmin(s)

	vasps := []string{
		"d9da630e-41aa-11ec-9d29-acde48001122",
	}

	noteID := "d9ebcfa4-41aa-11ec-9d29-acde48001122"

	// Supplying an invalid note ID
	req := &admin.ModifyReviewNoteRequest{
		VASP:   vasps[0],
		NoteID: "invalid slug",
	}
	params := map[string]string{
		"vaspID": vasps[0],
		"noteID": "invalid slug",
	}
	claims := &tokens.Claims{
		Email: "admin@example.com",
	}
	status := s.doRequest(a.UpdateReviewNote, http.MethodPut, "/v2/vasps/"+vasps[0]+"/notes/invalid", nil, params, req, nil, claims)
	require.Equal(http.StatusNotFound, status)

	// Supplying an invalid VASP ID
	req = &admin.ModifyReviewNoteRequest{
		VASP:   "invalid",
		NoteID: noteID,
	}
	params = map[string]string{
		"vaspID": "invalid",
		"noteID": noteID,
	}
	status = s.doRequest(a.UpdateReviewNote, http.MethodPut, "/v2/vasps/invalid/notes/"+noteID, nil, params, req, nil, claims)
	require.Equal(http.StatusNotFound, status)

	// Successfully updating a review note
	req = &admin.ModifyReviewNoteRequest{
		VASP:   vasps[0],
		NoteID: noteID,
		Text:   "bar",
	}
	params = map[string]string{
		"vaspID": vasps[0],
		"noteID": noteID,
	}
	actual := &admin.ReviewNote{}
	modified := time.Now()
	status = s.doRequest(a.UpdateReviewNote, http.MethodPut, "/v2/vasps/"+vasps[0]+"/notes/"+noteID, nil, params, req, actual, claims)
	require.Equal(http.StatusOK, status)
	// Validate returned note
	require.Equal(req.NoteID, actual.ID)
	require.Equal("2021-07-26T11:25:44Z", actual.Created)
	ts, err := time.Parse(time.RFC3339, actual.Modified)
	require.NoError(err)
	require.True(ts.Sub(modified) < time.Minute)
	require.Equal("admin@trisa.io", actual.Author)
	require.Equal(claims.Email, actual.Editor)
	require.Equal(req.Text, actual.Text)
	// Record on the database should be updated
	v, err := a.db.RetrieveVASP(vasps[0])
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
	a := initAdmin(s)

	vasps := []string{
		"d9da630e-41aa-11ec-9d29-acde48001122",
	}

	noteID := "d9ebcfa4-41aa-11ec-9d29-acde48001122"

	// Supplying an invalid note ID
	params := map[string]string{
		"vaspID": vasps[0],
		"noteID": "invalid slug",
	}
	status := s.doRequest(a.DeleteReviewNote, http.MethodDelete, "/v2/vasps/"+vasps[0]+"/notes/invalid", nil, params, nil, nil, nil)
	require.Equal(http.StatusNotFound, status)

	// Supplying an invalid VASP ID
	params = map[string]string{
		"vaspID": "invalid",
		"noteID": noteID,
	}
	status = s.doRequest(a.DeleteReviewNote, http.MethodDelete, "/v2/vasps/invalid/notes/"+noteID, nil, params, nil, nil, nil)
	require.Equal(http.StatusNotFound, status)

	// Successfully deleting a review note
	params = map[string]string{
		"vaspID": vasps[0],
		"noteID": noteID,
	}
	status = s.doRequest(a.DeleteReviewNote, http.MethodDelete, "/v2/vasps/"+vasps[0]+"/notes/"+noteID, nil, params, nil, nil, nil)
	require.Equal(http.StatusOK, status)
	// Record on the database should be deleted
	v, err := a.db.RetrieveVASP(vasps[0])
	require.NoError(err)
	notes, err := models.GetReviewNotes(v)
	require.NoError(err)
	require.Len(notes, 0)
}

// Test the Review endpoint.
func (s *gdsTestSuite) TestReview() {
	require := s.Require()
	a := initAdmin(s)

	vasps := []string{
		"d9da630e-41aa-11ec-9d29-acde48001122",
	}

	// Supplying an invalid VASP ID
	req := &admin.ReviewRequest{
		ID:                     "invalid",
		AdminVerificationToken: "foo",
		Accept:                 true,
	}
	params := map[string]string{
		"vaspID": "invalid",
	}
	status := s.doRequest(a.Review, http.MethodPost, "/v2/vasps/invalid/review", nil, params, req, nil, nil)
	require.Equal(http.StatusNotFound, status)

	// No verification token supplied
	req = &admin.ReviewRequest{
		ID:     vasps[0],
		Accept: true,
	}
	params = map[string]string{
		"vaspID": vasps[0],
	}
	status = s.doRequest(a.Review, http.MethodPost, "/v2/vasps/"+vasps[0]+"/review", nil, params, req, nil, nil)
	require.Equal(http.StatusBadRequest, status)

	// No rejection reason supplied
	req = &admin.ReviewRequest{
		ID:                     vasps[0],
		AdminVerificationToken: "foo",
		Accept:                 false,
	}
	status = s.doRequest(a.Review, http.MethodPost, "/v2/vasps/"+vasps[0]+"/review", nil, params, req, nil, nil)
	require.Equal(http.StatusBadRequest, status)

	// TODO: Test the accept and reject paths
}

// Test the Resend endpoint.
func (s *gdsTestSuite) TestResend() {
	require := s.Require()
	a := initAdmin(s)

	vaspErrored := "da2b165a-41aa-11ec-9d29-acde48001122"
	vaspRejected := "da8bd0e4-41aa-11ec-9d29-acde48001122"

	// Supplying an invalid VASP ID
	req := &admin.ResendRequest{
		ID:     "invalid",
		Action: admin.ResendVerifyContact,
	}
	params := map[string]string{
		"vaspID": "invalid",
	}
	status := s.doRequest(a.Resend, http.MethodPost, "/v2/vasps/invalid/resend", nil, params, req, nil, nil)
	require.Equal(http.StatusNotFound, status)

	// ResendVerifyContact email
	req = &admin.ResendRequest{
		ID:     vaspErrored,
		Action: admin.ResendVerifyContact,
		Reason: "verify",
	}
	params = map[string]string{
		"vaspID": vaspErrored,
	}
	actual := &admin.ResendReply{}
	sent := time.Now()
	status = s.doRequest(a.Resend, http.MethodPost, "/v2/vasps/"+vaspErrored+"/resend", nil, params, req, actual, nil)
	require.Equal(http.StatusOK, status)
	require.Equal(1, actual.Sent)
	require.Contains(actual.Message, "contact verification emails resent")
	// Email audit log should be updated
	v, err := a.db.RetrieveVASP(vaspErrored)
	require.NoError(err)
	emails, err := models.GetEmailLog(v.Contacts.Billing)
	require.NoError(err)
	require.Len(emails, 1)
	ts, err := time.Parse(time.RFC3339, emails[0].Timestamp)
	require.NoError(err)
	require.True(ts.Sub(sent) < time.Minute)
	require.Equal("verify_contact", emails[0].Reason)

	// ResendReview email
	req = &admin.ResendRequest{
		ID:     vaspErrored,
		Action: admin.ResendReview,
		Reason: "review",
	}
	params = map[string]string{
		"vaspID": vaspErrored,
	}
	actual = &admin.ResendReply{}
	sent = time.Now()
	status = s.doRequest(a.Resend, http.MethodPost, "/v2/vasps/"+vaspErrored+"/resend", nil, params, req, actual, nil)
	require.Equal(http.StatusOK, status)
	require.Equal(1, actual.Sent)
	require.Contains(actual.Message, "review request resent")

	// ResendRejection email
	req = &admin.ResendRequest{
		ID:     vaspRejected,
		Action: admin.ResendRejection,
		Reason: "reject",
	}
	params = map[string]string{
		"vaspID": vaspRejected,
	}
	actual = &admin.ResendReply{}
	sent = time.Now()
	status = s.doRequest(a.Resend, http.MethodPost, "/v2/vasps/"+vaspRejected+"/resend", nil, params, req, actual, nil)
	require.Equal(http.StatusOK, status)
	require.Equal(2, actual.Sent)
	require.Contains(actual.Message, "rejection emails resent")
	// Email audit logs should be updated
	v, err = a.db.RetrieveVASP(vaspRejected)
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
	vasps := []string{
		"d9da630e-41aa-11ec-9d29-acde48001122",
	}
	WriteVASPs(s, vasps)
	a := initAdmin(s)

	// Invalid start date
	status := s.doRequest(a.ReviewTimeline, http.MethodPost, "/v2/reviews?start=09-01-2021", nil, nil, nil, nil, nil)
	require.Equal(http.StatusBadRequest, status)

	// Start date is before epoch
	status = s.doRequest(a.ReviewTimeline, http.MethodPost, "/v2/reviews?start=1968-01-01", nil, nil, nil, nil, nil)
	require.Equal(http.StatusBadRequest, status)

	// Invalid end date
	status = s.doRequest(a.ReviewTimeline, http.MethodPost, "/v2/reviews?end=09-01-2021", nil, nil, nil, nil, nil)
	require.Equal(http.StatusBadRequest, status)

	// Start date is after end date
	status = s.doRequest(a.ReviewTimeline, http.MethodPost, "/v2/reviews?start=2021-01-01&end=2020-01-01", nil, nil, nil, nil, nil)
	require.Equal(http.StatusBadRequest, status)

	// Successful retrieval of review timeline
	actual := &admin.ReviewTimelineReply{}
	status = s.doRequest(a.ReviewTimeline, http.MethodPost, "/v2/reviews?start=2021-09-20&end=2021-09-30", nil, nil, nil, actual, nil)
	require.Equal(http.StatusOK, status)
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
