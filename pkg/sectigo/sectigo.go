/*
Package sectigo provides API access to the Sectigo IoT Manager 20.7, which is used to
sign certificate requests for directory service certificate issuance.
*/
package sectigo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Valid Sectigo Certificate Profile Names and IDs
// TODO: do not hardcode this, but get programatically from Sectigo API
const (
	ProfileCipherTraceEE                     = "CipherTrace EE"
	ProfileIDCipherTraceEE                   = "17"
	ProfileCipherTraceEndEntityCertificate   = "CipherTrace End Entity Certificate"
	ProfileIDCipherTraceEndEntityCertificate = "85"
)

var AllProfiles = [4]string{
	ProfileCipherTraceEE, ProfileIDCipherTraceEE,
	ProfileCipherTraceEndEntityCertificate, ProfileIDCipherTraceEndEntityCertificate,
}

// Sectigo provides authenticated http requests to the Sectigo IoT Manager 20.7 REST API.
// See documentation at: https://support.sectigo.com/Com_KnowledgeDetailPage?Id=kA01N000000bvCJ
//
// Most endpoints require an JWT access token set in an Authorization: Bearer header to
// provide information about an authenticated user. The authenticate method will request
// access and refresh tokens based on user credentials. Each access token has a validity
// of 600 seconds, when the access token expires, the refresh token should be used to
// request a new access token without requiring the user to resupply credentials.
//
// The client handles authentication by checking if the tokens are valid before every
// request, and if not either refreshes the token or reauthenticates using its
// credentials.
type Sectigo struct {
	sync.RWMutex
	client  http.Client
	creds   *Credentials
	profile string
	testing bool
}

// New creates a Sectigo client ready to make HTTP requests, but unauthenticated. The
// username and password will be loaded from the environment if not given - from
// $SECTIGO_USERNAME and $SECTIGO_PASSWORD respectively; alternatively if not given and
// not stored in the environment, as long as valid access credentials are cached the
// credentials will be loaded.
func New(conf Config) (client *Sectigo, err error) {
	client = &Sectigo{
		creds: &Credentials{},
		client: http.Client{
			CheckRedirect: certificateAuthRedirectPolicy,
		},
		profile: conf.Profile,
		testing: conf.Testing,
	}

	if conf.Testing {
		// Add mock credentials to the client if we're in testing mode
		if conf.Username == "" {
			conf.Username = MockUsername
		}

		if conf.Password == "" {
			conf.Password = MockPassword
		}
	}

	if err = client.creds.Load(conf.Username, conf.Password); err != nil {
		return nil, err
	}

	return client, nil
}

// Authenticate the user with the specified credentials to get new access and refresh tokens.
// This method will replace the access tokens even if already present and valid. If
// certificate authentication is enabled then the response will be a 307 status code,
// if wrong user name and password a 401 status code and if a correct user name and
// password but the user does not have authority, a 403 status code.
func (s *Sectigo) Authenticate() (err error) {
	data := AuthenticationRequest{
		Username: s.creds.Username,
		Password: s.creds.Password,
	}

	body := new(bytes.Buffer)
	if err = json.NewEncoder(body).Encode(data); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, urlFor(AuthenticateEP), body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", contentType)

	rep, err := s.Do(req)
	if err != nil {
		return err
	}
	defer rep.Body.Close()

	// Handle error states
	switch rep.StatusCode {
	case http.StatusUnauthorized:
		return ErrInvalidCredentials
	case http.StatusForbidden:
		return ErrNotAuthorized
	case http.StatusTemporaryRedirect:
		return ErrMustUseTLSAuth
	}

	if rep.StatusCode != http.StatusOK {
		return fmt.Errorf("unhandled status code: %d", rep.StatusCode)
	}

	// We've got a successful response - deserialize request body
	tokens := &AuthenticationReply{}
	if err = json.NewDecoder(rep.Body).Decode(&tokens); err != nil {
		return err
	}

	if err = s.creds.Update(tokens.AccessToken, tokens.RefreshToken); err != nil {
		return err
	}
	return nil
}

// Refresh the access token using the refresh token. Note that this method does not
// check if the credentials are refreshable, it only issues the refresh request with
// the refresh access token if it exists. If the refresh token does not exist, then an
// error is returned.
func (s *Sectigo) Refresh() (err error) {
	if s.creds.RefreshToken == "" {
		return ErrNotAuthenticated
	}

	body := new(bytes.Buffer)
	fmt.Fprintf(body, "%s", s.creds.RefreshToken)

	req, err := http.NewRequest(http.MethodPost, urlFor(RefreshEP), body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", contentType)

	rep, err := s.Do(req)
	if err != nil {
		return err
	}
	defer rep.Body.Close()

	// Handle error states
	switch rep.StatusCode {
	case http.StatusUnauthorized:
		return ErrInvalidCredentials
	case http.StatusForbidden:
		return ErrNotAuthorized
	}

	if rep.StatusCode != http.StatusOK {
		return fmt.Errorf("unhandled status code: %d", rep.StatusCode)
	}

	// We've got a successful response - deserialize request body
	tokens := &AuthenticationReply{}
	if err = json.NewDecoder(rep.Body).Decode(&tokens); err != nil {
		return err
	}

	// It appears that sectigo reuses the refresh token.
	// TODO: verify refresh behavior to ensure that it's used correctly
	if tokens.RefreshToken == "" {
		tokens.RefreshToken = s.creds.RefreshToken
	}

	if err = s.creds.Update(tokens.AccessToken, tokens.RefreshToken); err != nil {
		return err
	}
	return nil
}

// CreateSingleCertBatch issues a new single certificate batch.
// User must be authenticated with role 'USER' and has permission to create request.
// You may get http code 400 if supplied values in profileParams fails to validate over
// rules specified in "profile".
func (s *Sectigo) CreateSingleCertBatch(authority int, name string, params map[string]string) (batch *BatchResponse, err error) {
	// perform preflight check for authenticated endpoint
	if err = s.preflight(); err != nil {
		return nil, err
	}

	batchInfo := &CreateSingleCertBatchRequest{
		AuthorityID:   authority,
		BatchName:     name,
		ProfileParams: params,
	}

	// create request
	var req *http.Request
	if req, err = s.newRequest(http.MethodPut, urlFor(CreateSingleCertBatchEP), batchInfo); err != nil {
		return nil, err
	}

	var rep *http.Response
	if rep, err = s.Do(req); err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	if err = s.checkStatus(rep); err != nil {
		return nil, err
	}

	if err = json.NewDecoder(rep.Body).Decode(&batch); err != nil {
		return nil, err
	}
	return batch, nil
}

// UploadCSRBatch CSR or bulk ZIP file.
// User must be authenticated with role 'USER' and must has permission to read this profile.
// That part contains the CSRs with a Content-Dispostion of form-data and a name
// parameter of files. The filename parameter is not used. The Content-Type of the part
// should match the uploaded file. The uploaded CSRs can be a single text file with
// multiple CSRs in PEM form using standard BEGIN/END separators or a zip file
// containing multiple CSRs files. When uploading a single text file the Content-Type
// can be text/plain, application/octet-stream or application/x-x509-ca-cert. When
// uploading a zip file the Content-Type must be application/zip. The zip file must
// contain each CSR in a file with the extension .csr or .pem.
func (s *Sectigo) UploadCSRBatch(profileId int, filename string, csrData []byte, params map[string]string) (batch *BatchResponse, err error) {
	// perform preflight check for authenticated endpoint
	if err = s.preflight(); err != nil {
		return nil, err
	}

	// do not continue if we're not authenticated
	if !s.creds.Valid() {
		return nil, ErrNotAuthenticated
	}

	// create multipart body for uploading file and form data
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// write the profileID as the first part
	var part io.Writer
	if part, err = writer.CreateFormField("profileId"); err != nil {
		return nil, fmt.Errorf("could not write profileId part: %s", err)
	}
	if _, err = fmt.Fprintf(part, "%d", profileId); err != nil {
		return nil, fmt.Errorf("could not write profileId: %s", err)
	}

	// write the csr batch request JSON serialized into the multipart request
	batchInfo := &UploadCSRBatchRequest{
		ProfileID:     profileId,
		BatchType:     "BATCH",
		ProfileParams: params,
	}
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="csrBatchRequest"; filename="blob"`)
	header.Set("Content-Type", "application/json")
	if part, err = writer.CreatePart(header); err != nil {
		return nil, fmt.Errorf("could not write csrBatchRequest part: %s", err)
	}
	if err = json.NewEncoder(part).Encode(batchInfo); err != nil {
		return nil, fmt.Errorf("could not encode csrBatchRequest: %s", err)
	}

	// write the csr data itself into the multipart request
	// NOTE: according to the documentation the filename is not used
	if part, err = writer.CreateFormFile("files", filename); err != nil {
		return nil, fmt.Errorf("could not write files part: %s", err)
	}
	if _, err = part.Write(csrData); err != nil {
		return nil, fmt.Errorf("could not write csrData: %s", err)
	}

	// make sure the writer is closed before creating the request to ensure that the
	// content length header is accurate (otherwise there will be a 500 error)
	writer.Close()

	// create multipart request (cannot use newRequest to ensure multipart is constructed correctly)
	var req *http.Request
	if req, err = http.NewRequest(http.MethodPost, urlFor(UploadCSREP), body); err != nil {
		return nil, err
	}

	// set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.creds.AccessToken))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", contentType)
	req.Header.Set("User-Agent", userAgent)

	// execute the request
	var rep *http.Response
	if rep, err = s.Do(req); err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	// handle http errors
	if err = s.checkStatus(rep); err != nil {
		return nil, err
	}

	// return the response
	if err = json.NewDecoder(rep.Body).Decode(&batch); err != nil {
		return nil, err
	}
	return batch, nil
}

// BatchDetail returns batch information by batch id.
// User must be authenticated with role 'USER' and has permission to read this batch.
func (s *Sectigo) BatchDetail(id int) (batch *BatchResponse, err error) {
	// perform preflight check for authenticated endpoint
	if err = s.preflight(); err != nil {
		return nil, err
	}

	// create request
	var req *http.Request
	if req, err = s.newRequest(http.MethodGet, urlFor(BatchDetailEP, id), nil); err != nil {
		return nil, err
	}

	var rep *http.Response
	if rep, err = s.Do(req); err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	if err = s.checkStatus(rep); err != nil {
		return nil, err
	}

	if err = json.NewDecoder(rep.Body).Decode(&batch); err != nil {
		return nil, err
	}
	return batch, nil
}

// BatchStatus views batch status by batchId.
// User must be authenticated with role 'USER' and has permission to read this batch.
// Statuses can be: "Ready for download", "Processing", "Failed" or "Not Acceptable. Status: ".
func (s *Sectigo) BatchStatus(batch int) (status string, err error) {
	// perform preflight check for authenticated endpoint
	if err = s.preflight(); err != nil {
		return "", err
	}

	// create request
	var req *http.Request
	if req, err = s.newRequest(http.MethodGet, urlFor(BatchStatusEP, batch), nil); err != nil {
		return "", err
	}

	var rep *http.Response
	if rep, err = s.Do(req); err != nil {
		return "", err
	}
	defer rep.Body.Close()

	if err = s.checkStatus(rep); err != nil {
		return "", err
	}

	var data []byte
	if data, err = ioutil.ReadAll(rep.Body); err != nil {
		return "", err
	}

	// Ensure the status appears as though it is one of our status constants
	status = strings.TrimSpace(string(data))
	status = strings.ToUpper(status)
	status = strings.Replace(status, " ", "_", -1)
	return status, nil
}

// ProcessingInfo returns batch processing status by batch id.
// User must be authenticated with role 'USER' and has permission to read this batch.
func (s *Sectigo) ProcessingInfo(batch int) (status *ProcessingInfoResponse, err error) {
	// perform preflight check for authenticated endpoint
	if err = s.preflight(); err != nil {
		return nil, err
	}

	// create request
	var req *http.Request
	if req, err = s.newRequest(http.MethodGet, urlFor(BatchProcessingInfoEP, batch), nil); err != nil {
		return nil, err
	}

	var rep *http.Response
	if rep, err = s.Do(req); err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	if err = s.checkStatus(rep); err != nil {
		return nil, err
	}

	if err = json.NewDecoder(rep.Body).Decode(&status); err != nil {
		return nil, err
	}
	return status, nil
}

// Download batch as a ZIP file.
// Dir should be a directory, filename is detected from content-disposition.
// User must be authenticated with role 'USER' and batch must be readable.
func (s *Sectigo) Download(batch int, dir string) (path string, err error) {
	// Verify download location
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("directory %q does not exist", dir)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("path %q is not a directory", dir)
	}

	// perform preflight check for authenticated endpoint
	if err = s.preflight(); err != nil {
		return "", err
	}

	// create request
	var req *http.Request
	if req, err = s.newRequest(http.MethodGet, urlFor(DownloadEP, batch), nil); err != nil {
		return "", err
	}

	// Set different content-type and accept headers
	req.Header.Set("Content-Type", downloadContentType)
	req.Header.Set("Accept", downloadContentType)

	var rep *http.Response
	if rep, err = s.Do(req); err != nil {
		return "", err
	}
	defer rep.Body.Close()

	// Parse the Content-Disposition header to get the download filename
	var filename string
	contentDisposition := rep.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err == nil {
			filename = params["filename"]
		}
	}

	if filename == "" {
		filename = fmt.Sprintf("%d.zip", batch)
	}
	path = filepath.Join(dir, filename)

	// Create the file to write the download into
	// TODO: get the filename from the headers and treat path as a dir
	var out *os.File
	if out, err = os.Create(path); err != nil {
		return "", err
	}
	defer out.Close()

	if _, err = io.Copy(out, rep.Body); err != nil {
		return "", err
	}
	return path, nil
}

// LicensesUsed returns statistic for Ordered/Issued certificates (licenses used)
// User must be authenticated with role 'USER'
func (s *Sectigo) LicensesUsed() (stats *LicensesUsedResponse, err error) {
	// perform preflight check for authenticated endpoint
	if err = s.preflight(); err != nil {
		return nil, err
	}

	// create request
	var req *http.Request
	if req, err = s.newRequest(http.MethodGet, urlFor(DevicesEP), nil); err != nil {
		return nil, err
	}

	var rep *http.Response
	if rep, err = s.Do(req); err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	if err = s.checkStatus(rep); err != nil {
		return nil, err
	}

	if err = json.NewDecoder(rep.Body).Decode(&stats); err != nil {
		return nil, err
	}
	return stats, nil
}

// UserAuthorities returns a list of all Authorities by Ecosystem and Current User
// User must be authenticated.
func (s *Sectigo) UserAuthorities() (authorities []*AuthorityResponse, err error) {
	// perform preflight check for authenticated endpoint
	if err = s.preflight(); err != nil {
		return nil, err
	}

	// create request
	var req *http.Request
	if req, err = s.newRequest(http.MethodGet, urlFor(UserAuthoritiesEP), nil); err != nil {
		return nil, err
	}

	var rep *http.Response
	if rep, err = s.Do(req); err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	if err = s.checkStatus(rep); err != nil {
		return nil, err
	}

	if err = json.NewDecoder(rep.Body).Decode(&authorities); err != nil {
		return nil, err
	}
	return authorities, nil
}

// AuthorityAvailableBalance returns balance available for the specified user/authority
// User must be authenticated.
func (s *Sectigo) AuthorityAvailableBalance(id int) (balance int, err error) {
	// perform preflight check for authenticated endpoint
	if err = s.preflight(); err != nil {
		return 0, err
	}

	// create request
	var req *http.Request
	if req, err = s.newRequest(http.MethodGet, urlFor(AuthorityUserBalanceAvailableEP, id), nil); err != nil {
		return 0, err
	}

	var rep *http.Response
	if rep, err = s.Do(req); err != nil {
		return 0, err
	}
	defer rep.Body.Close()

	if err = s.checkStatus(rep); err != nil {
		return 0, err
	}

	if err = json.NewDecoder(rep.Body).Decode(&balance); err != nil {
		return 0, err
	}
	return balance, nil

}

// Profiles returns a list of all profiles available to the user.
// User must be authenticated.
func (s *Sectigo) Profiles() (profiles []*ProfileResponse, err error) {
	// perform preflight check for authenticated endpoint
	if err = s.preflight(); err != nil {
		return nil, err
	}

	// create request
	var req *http.Request
	if req, err = s.newRequest(http.MethodGet, urlFor(ProfilesEP), nil); err != nil {
		return nil, err
	}

	var rep *http.Response
	if rep, err = s.Do(req); err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	if err = s.checkStatus(rep); err != nil {
		return nil, err
	}

	if err = json.NewDecoder(rep.Body).Decode(&profiles); err != nil {
		return nil, err
	}
	return profiles, nil
}

// ProfileParams lists the parameters acceptable and required by profileId
// User must be authenticated with role 'ADMIN' or 'USER' and permission to read this profile
func (s *Sectigo) ProfileParams(id int) (params []*ProfileParamsResponse, err error) {
	// perform preflight check for authenticated endpoint
	if err = s.preflight(); err != nil {
		return nil, err
	}

	// create request
	var req *http.Request
	if req, err = s.newRequest(http.MethodGet, urlFor(ProfileParametersEP, id), nil); err != nil {
		return nil, err
	}

	var rep *http.Response
	if rep, err = s.Do(req); err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	if err = s.checkStatus(rep); err != nil {
		return nil, err
	}

	if err = json.NewDecoder(rep.Body).Decode(&params); err != nil {
		return nil, err
	}
	return params, nil
}

// ProfileDetail gets extended profile information.
// User must be authenticated with role 'ADMIN' or 'USER' and permission to read this profile.
func (s *Sectigo) ProfileDetail(id int) (profile *ProfileDetailResponse, err error) {
	// perform preflight check for authenticated endpoint
	if err = s.preflight(); err != nil {
		return nil, err
	}

	// create request
	var req *http.Request
	if req, err = s.newRequest(http.MethodGet, urlFor(ProfileDetailEP, id), nil); err != nil {
		return nil, err
	}

	var rep *http.Response
	if rep, err = s.Do(req); err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	if err = s.checkStatus(rep); err != nil {
		return nil, err
	}

	if err = json.NewDecoder(rep.Body).Decode(&profile); err != nil {
		return nil, err
	}
	return profile, nil
}

// Organization returns the organization for the current user. User must be authorized.
func (s *Sectigo) Organization() (org *OrganizationResponse, err error) {
	// perform preflight check for authenticated endpoint
	if err = s.preflight(); err != nil {
		return nil, err
	}

	// create request
	var req *http.Request
	if req, err = s.newRequest(http.MethodGet, urlFor(CurrentUserOrganizationEP), nil); err != nil {
		return nil, err
	}

	var rep *http.Response
	if rep, err = s.Do(req); err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	if err = s.checkStatus(rep); err != nil {
		return nil, err
	}

	if err = json.NewDecoder(rep.Body).Decode(&org); err != nil {
		return nil, err
	}
	return org, nil
}

// FindCertificate searches for certificates by common name and serial number.
func (s *Sectigo) FindCertificate(commonName, serialNumber string) (certs *FindCertificateResponse, err error) {
	// perform preflight check for authenticated endpoint
	if err = s.preflight(); err != nil {
		return nil, err
	}

	query := &FindCertificateRequest{
		CommonName:   commonName,
		SerialNumber: serialNumber,
	}

	// create request
	var req *http.Request
	if req, err = s.newRequest(http.MethodPost, urlFor(FindCertificateEP), query); err != nil {
		return nil, err
	}

	var rep *http.Response
	if rep, err = s.Do(req); err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	if err = s.checkStatus(rep); err != nil {
		return nil, err
	}

	if err = json.NewDecoder(rep.Body).Decode(&certs); err != nil {
		return nil, err
	}
	return certs, nil
}

// RevokeCertificate by serial number if the certificate was signed by the given authority.
// A reason code from RFC 5280 must be given. This method revokes single certificates
// unlike the RevokeDeviceCertificates method which can revoke multiple certificates by
// their assignment to specific Device IDs. If no error is returned, the certificate
// revocation was successful.
// User must be authenticated and has permission to update profile.
func (s *Sectigo) RevokeCertificate(profileID, reasonCode int, serialNumber string) (err error) {
	// perform preflight check for authenticated endpoint
	if err = s.preflight(); err != nil {
		return err
	}

	query := &RevokeCertificateRequest{
		ReasonCode:   reasonCode,
		SerialNumber: serialNumber,
	}

	// create request
	var req *http.Request
	if req, err = s.newRequest(http.MethodPost, urlFor(RevokeCertificateEP, profileID), query); err != nil {
		return err
	}

	var rep *http.Response
	if rep, err = s.Do(req); err != nil {
		return err
	}
	defer rep.Body.Close()

	if err = s.checkStatus(rep); err != nil {
		return err
	}
	return nil
}

// Creds returns a copy of the underlying credentials object.
func (s *Sectigo) Creds() Credentials {
	return *s.creds
}

// Returns a request with default headers set along with the authentication header.
// If the client has not been authenticated, then an error is returned.
func (s *Sectigo) newRequest(method, url string, data interface{}) (req *http.Request, err error) {
	s.RLock()
	defer s.RUnlock()
	if !s.creds.Valid() {
		return nil, ErrNotAuthenticated
	}

	if data != nil {
		// JSON serialize the data being sent to the request
		body := new(bytes.Buffer)
		if err = json.NewEncoder(body).Encode(data); err != nil {
			return nil, err
		}

		if req, err = http.NewRequest(method, url, body); err != nil {
			return nil, err
		}
	} else {
		// Create a request with an empty body
		if req, err = http.NewRequest(method, url, nil); err != nil {
			return nil, err
		}
	}

	// Set Headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.creds.AccessToken))
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", contentType)
	req.Header.Set("User-Agent", userAgent)

	return req, nil
}

// Preflight prepares to send a request that needs to be authenticated by checking the
// credentials and sending any authentication or refresh requests required.
func (s *Sectigo) preflight() (err error) {
	s.Lock()
	defer s.Unlock()
	if !s.creds.Valid() {
		if s.creds.Refreshable() {
			// Attempt to refresh the credentials, if there is no error, then continue.
			// However, the refresh endpoint may not be working so if it errors, attempt
			// to reauthenticate with username and password instead.
			// TODO: add logging so the server knows what's going on.
			if err = s.Refresh(); err == nil {
				return err
			}
		}

		// If we could not refresh, attempt to reauthenticate
		if err = s.Authenticate(); err != nil {
			return err
		}
	}

	// Check the credentials and if they're good, dump them to disk
	if err = s.creds.Check(); err != nil {
		s.creds.Clear()
		s.creds.Dump()
		return err
	}

	// Ignore any cache errors on reauthentication
	s.creds.Dump()

	// Good to go!
	return nil
}

// Do performs a sectigo client request and returns the response.
func (s *Sectigo) Do(req *http.Request) (*http.Response, error) {
	if s.testing {
		// Ensure that we're sending the requests to a test server
		host := baseURL.Hostname()
		if host != "localhost" && host != "127.0.0.1" {
			return nil, fmt.Errorf("sectigo hostname must be set to localhost in testing mode, is %s", host)
		}
	}
	return s.client.Do(req)
}

// Helper function to convert a non-200 HTTP status into an error, reading JSON error
// data if it's available, otherwise returning a simple error. Note that this method
// will attempt to read the body on error, so do not use it for error handling that
// requires knowledge of the body.
func (s *Sectigo) checkStatus(rep *http.Response) (err error) {
	// Check if status code is a good status.
	if rep.StatusCode >= 200 && rep.StatusCode < 300 {
		return nil
	}

	// Try to unmarshall the error from the response
	var e *APIError
	if err = json.NewDecoder(rep.Body).Decode(&e); err != nil {
		switch rep.StatusCode {
		case http.StatusUnauthorized:
			return ErrNotAuthenticated
		case http.StatusForbidden:
			return ErrNotAuthorized
		}

		// Return a simple error since the JSON could not be decoded.
		e = &APIError{
			Status:  rep.StatusCode,
			Message: rep.Status,
		}
	}
	return e
}

// The Sectigo API has a special authentication policy when certificate authentication
// is enabled. In this case, normal password authentication requests return a 307 with
// a URL to POST to the certificate auth location, meaning that the URL requires TLS
// client authentication. TLS client authentication is only required to obtain an access
// token. This function prevents multiple redirects in the case of a 307 by returning
// the redirect request. Other redirect status codes are followed.
func certificateAuthRedirectPolicy(req *http.Request, via []*http.Request) error {
	if req.Response.StatusCode == http.StatusTemporaryRedirect {
		return http.ErrUseLastResponse
	}
	return nil
}

func (s *Sectigo) Profile() string {
	return s.profile
}
