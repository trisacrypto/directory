package sectigo

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

func NewMock(username, password, profile string) (client *Sectigo, err error) {
	client = &Sectigo{
		creds: &mockCredentialsManager{
			creds: Credentials{},
		},
		client: &mockHTTPClient{
			CheckRedirect: certificateAuthRedirectPolicy,
		},
		profile: profile,
	}

	if err = client.creds.Load(username, password); err != nil {
		return nil, err
	}

	return client, nil
}

type mockCredentialsManager struct {
	creds Credentials
}

type mockHTTPClient struct {
	CheckRedirect func(req *http.Request, via []*http.Request) error
}

type mockCredentials struct {
	creds Credentials
}

var mockCredentialsCache Credentials

type mockUser struct {
	username   string
	password   string
	authorized bool
	access     string
	refresh    string
}

type mockBatch struct {
	processing ProcessingInfoResponse
	filename   string
}

type mockProfile struct {
	params  []*ProfileParamsResponse
	details *ProfileDetailResponse
}

type mockCertificate struct {
	DeviceID     int
	CommonName   string
	SerialNumber string
	CreationDate string
	Status       string
}

type mockServer struct {
	certs        bool
	users        map[string]*mockUser
	access       map[string]string
	refresh      map[string]string
	ids          map[int]string
	batches      map[int]*mockBatch
	numBatches   int
	licenses     *LicensesUsedResponse
	authorities  map[int]*AuthorityResponse
	profiles     map[int]*mockProfile
	org          *OrganizationResponse
	certificates []*mockCertificate
}

var mockEnv map[string]string
var mockBackend *mockServer

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (c *mockHTTPClient) Do(req *http.Request) (rep *http.Response, err error) {
	switch req.URL.Path {
	case urlFor(authenticateEP):
		var request AuthenticationRequest
		if err = json.NewDecoder(req.Body).Decode(&request); err != nil {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
			}, err
		}

		if mockBackend.certs {
			return &http.Response{
				StatusCode: http.StatusTemporaryRedirect,
			}, fmt.Errorf("certificate authentication is enabled")
		}

		var user *mockUser
		if user, ok := mockBackend.users[request.Username]; !ok || user.password != request.Password {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
			}, fmt.Errorf("invalid username or password")
		}

		if !user.authorized {
			return &http.Response{
				StatusCode: http.StatusForbidden,
			}, fmt.Errorf("user is not authorized")
		}

		delete(mockBackend.access, user.access)
		delete(mockBackend.refresh, user.refresh)
		response := AuthenticationReply{}
		if response.AccessToken, err = generateToken(); err != nil {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
		if response.RefreshToken, err = generateToken(); err != nil {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		mockBackend.access[response.AccessToken] = request.Username
		mockBackend.refresh[response.RefreshToken] = request.Username
		user.refresh = response.RefreshToken

		repBody := bytes.NewBuffer(nil)
		if err = json.NewEncoder(repBody).Encode(response); err != nil {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(repBody),
		}, nil
	case urlFor(refreshEP):
		var reqBody []byte
		if reqBody, err = ioutil.ReadAll(req.Body); err != nil {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		refresh := string(reqBody)
		var name string
		var user *mockUser
		var ok bool
		if name, ok = mockBackend.refresh[refresh]; !ok {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
			}, fmt.Errorf("invalid refresh token")
		}

		if user, ok = mockBackend.users[name]; !ok {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}, fmt.Errorf("user not found")
		}

		if !user.authorized {
			return &http.Response{
				StatusCode: http.StatusForbidden,
			}, fmt.Errorf("user is not authorized")
		}

		delete(mockBackend.access, user.access)
		delete(mockBackend.refresh, user.refresh)
		response := AuthenticationReply{}
		if response.AccessToken, err = generateToken(); err != nil {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
		if response.RefreshToken, err = generateToken(); err != nil {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		repBody := bytes.NewBuffer(nil)
		if err = json.NewEncoder(repBody).Encode(response); err != nil {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(repBody),
		}, nil
	case urlFor(createSingleCertBatchEP):
		var request CreateSingleCertBatchRequest
		if err = json.NewDecoder(req.Body).Decode(&request); err != nil {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
			}, err
		}

		response := &BatchResponse{}
		mockBackend.batches[mockBackend.numBatches] = &mockBatch{
			filename: "foo",
			processing: ProcessingInfoResponse{
				Active:  0,
				Success: 1,
				Failed:  0,
			},
		}
		mockBackend.numBatches++
		repBody := bytes.NewBuffer(nil)
		if err = json.NewEncoder(repBody).Encode(response); err != nil {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(repBody),
		}, nil
	case urlFor(uploadCSREP):
		var request UploadCSRBatchRequest
		if err = json.NewDecoder(req.Body).Decode(&request); err != nil {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
			}, err
		}

		response := &BatchResponse{}
		mockBackend.batches[mockBackend.numBatches] = &mockBatch{
			filename: "foo",
			processing: ProcessingInfoResponse{
				Active:  0,
				Success: 1,
				Failed:  0,
			},
		}
		mockBackend.numBatches++
		repBody := bytes.NewBuffer(nil)
		if err = json.NewEncoder(repBody).Encode(response); err != nil {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(repBody),
		}, nil
	case urlFor(batchDetailEP):
		param := req.URL.Query().Get("id")
		var id int
		if id, err = strconv.Atoi(param); err != nil {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
			}, err
		}
		var batch *mockBatch
		var ok bool
		if batch, ok = mockBackend.batches[id]; !ok {
			return &http.Response{
				StatusCode: http.StatusNotFound,
			}, fmt.Errorf("batch not found")
		}

		response := &BatchResponse{
			Active: batch.processing.Active > 0,
		}
		repBody := bytes.NewBuffer(nil)
		if err = json.NewEncoder(repBody).Encode(response); err != nil {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(repBody),
		}, nil
	case urlFor(batchProcessingInfoEP):
		param := req.URL.Query().Get("id")
		var id int
		if id, err = strconv.Atoi(param); err != nil {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
			}, err
		}
		var batch *mockBatch
		var ok bool
		if batch, ok = mockBackend.batches[id]; !ok {
			return &http.Response{
				StatusCode: http.StatusNotFound,
			}, fmt.Errorf("batch not found")
		}

		repBody := bytes.NewBuffer(nil)
		if err = json.NewEncoder(repBody).Encode(batch.processing); err != nil {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(repBody),
		}, nil
	case urlFor(downloadEP):
		param := req.URL.Query().Get("id")
		var id int
		if id, err = strconv.Atoi(param); err != nil {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
			}, err
		}
		var batch *mockBatch
		var ok bool
		if batch, ok = mockBackend.batches[id]; !ok {
			return &http.Response{
				StatusCode: http.StatusNotFound,
			}, fmt.Errorf("batch not found")
		}

		repBody := bytes.NewBuffer(nil)
		if err = json.NewEncoder(repBody).Encode(batch.processing); err != nil {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		response := &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(repBody),
		}
		response.Header.Set("Content-Disposition", fmt.Sprintf("filename=%s", batch.filename))
	case urlFor(devicesEP):
		repBody := bytes.NewBuffer(nil)
		if err = json.NewEncoder(repBody).Encode(mockBackend.licenses); err != nil {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(repBody),
		}, nil
	case urlFor(userAuthoritiesEP):
		repBody := bytes.NewBuffer(nil)
		auth := make([]*AuthorityResponse, 0, len(mockBackend.authorities))
		for _, value := range mockBackend.authorities {
			auth = append(auth, value)
		}
		if err = json.NewEncoder(repBody).Encode(auth); err != nil {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(repBody),
		}, nil
	case urlFor(authorityUserBalanceAvailableEP):
		param := req.URL.Query().Get("id")
		var id int
		if id, err = strconv.Atoi(param); err != nil {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
			}, err
		}
		var auth *AuthorityResponse
		var ok bool
		if auth, ok = mockBackend.authorities[id]; !ok {
			return &http.Response{
				StatusCode: http.StatusNotFound,
			}, fmt.Errorf("authority not found")
		}

		repBody := bytes.NewBuffer(nil)
		if err = json.NewEncoder(repBody).Encode(auth.Balance); err != nil {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(repBody),
		}, nil
	case urlFor(profilesEP):
		repBody := bytes.NewBuffer(nil)
		profiles := make([]*ProfileResponse, 0, len(mockBackend.profiles))
		for _, value := range mockBackend.profiles {
			profiles = append(profiles, &ProfileResponse{
				ProfileID: value.details.ProfileID,
			})
		}
		if err = json.NewEncoder(repBody).Encode(profiles); err != nil {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(repBody),
		}, nil
	case urlFor(profileParametersEP):
		param := req.URL.Query().Get("id")
		var id int
		if id, err = strconv.Atoi(param); err != nil {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
			}, err
		}
		var profile *mockProfile
		var ok bool
		if profile, ok = mockBackend.profiles[id]; !ok {
			return &http.Response{
				StatusCode: http.StatusNotFound,
			}, fmt.Errorf("profile not found")
		}

		repBody := bytes.NewBuffer(nil)
		if err = json.NewEncoder(repBody).Encode(profile.params); err != nil {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(repBody),
		}, nil
	case urlFor(profileDetailEP):
		param := req.URL.Query().Get("id")
		var id int
		if id, err = strconv.Atoi(param); err != nil {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
			}, err
		}
		var profile *mockProfile
		var ok bool
		if profile, ok = mockBackend.profiles[id]; !ok {
			return &http.Response{
				StatusCode: http.StatusNotFound,
			}, fmt.Errorf("profile not found")
		}

		repBody := bytes.NewBuffer(nil)
		if err = json.NewEncoder(repBody).Encode(profile.details); err != nil {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(repBody),
		}, nil
	case urlFor(currentUserOrganizationEP):
		repBody := bytes.NewBuffer(nil)
		if err = json.NewEncoder(repBody).Encode(mockBackend.org); err != nil {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(repBody),
		}, nil
	case urlFor(findCertificateEP):
		var request *FindCertificateRequest
		if err = json.NewDecoder(req.Body).Decode(request); err != nil {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
			}, err
		}

		response := &FindCertificateResponse{
			TotalCount: 0,
		}
		for _, c := range mockBackend.certificates {
			if c.SerialNumber == request.SerialNumber || c.CommonName == request.CommonName {
				response.TotalCount++
				response.Items = append(response.Items, struct {
					DeviceID     int    `json:"deviceId"`
					CommonName   string `json:"commonName"`
					SerialNumber string `json:"serialNumber"`
					CreationDate string `json:"creationDate"`
					Status       string `json:"status"`
				}{
					DeviceID:     c.DeviceID,
					CommonName:   c.CommonName,
					SerialNumber: c.SerialNumber,
					CreationDate: c.CreationDate,
					Status:       c.Status,
				})
			}
		}

		repBody := bytes.NewBuffer(nil)
		if err = json.NewEncoder(repBody).Encode(response); err != nil {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(repBody),
		}, nil
	case urlFor(revokeCertificateEP):
		var request *RevokeCertificateRequest
		if err = json.NewDecoder(req.Body).Decode(request); err != nil {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
			}, err
		}

		var found bool
		for _, c := range mockBackend.certificates {
			if c.SerialNumber == request.SerialNumber {
				found = true
				break
			}
		}

		if !found {
			return &http.Response{
				StatusCode: http.StatusNotFound,
			}, fmt.Errorf("certificate not found")
		}

		return &http.Response{
			StatusCode: http.StatusOK,
		}, nil
	default:
		return &http.Response{
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("unsupported request")
	}
	return &http.Response{
		StatusCode: http.StatusInternalServerError,
	}, fmt.Errorf("error parsing URL path")
}

func (c *mockCredentialsManager) Creds() *Credentials {
	return &c.creds
}

func (c *mockCredentialsManager) Load(username, password string) (err error) {
	var ok bool
	if username == "" {
		if username, ok = mockEnv[UsernameEnv]; ok {
			c.creds.Username = username
		}
	} else {
		c.creds.Username = username
	}

	if password == "" {
		if password, ok = mockEnv[PasswordEnv]; ok {
			c.creds.Password = password
		}
	} else {
		c.creds.Password = password
	}

	c.creds = mockCredentialsCache

	if err = c.Check(); err != nil {
		c.Clear()
		c.Dump()
	}

	if (c.creds.Username != "" || c.creds.Password != "") && (c.creds.Username == "" || c.creds.Password == "") {
		return ErrCredentialsMismatch
	}
	if c.creds.Username == "" && c.creds.Password == "" && c.creds.AccessToken == "" && c.creds.RefreshToken == "" {
		return ErrNoCredentials
	}

	return nil
}

func (c *mockCredentialsManager) Dump() (path string, err error) {
	mockCredentialsCache = c.creds
	return "", nil
}

func (c *mockCredentialsManager) Update(accessToken, refreshToken string) (err error) {
	var atc, rtc *jwt.StandardClaims
	if atc, err = parseToken(accessToken); err != nil {
		return fmt.Errorf("could not parse access token: %s", err)
	}

	if rtc, err = parseToken(refreshToken); err != nil {
		return fmt.Errorf("could not parse refresh token: %s", err)
	}

	c.creds.AccessToken = accessToken
	c.creds.RefreshToken = refreshToken
	c.creds.Subject = atc.Subject
	c.creds.IssuedAt = time.Unix(atc.IssuedAt, 0)
	c.creds.ExpiresAt = time.Unix(atc.ExpiresAt, 0)
	c.creds.RefreshBy = time.Unix(rtc.ExpiresAt, 0)

	if rtc.NotBefore > 0 {
		c.creds.NotBefore = time.Unix(rtc.NotBefore, 0)
	} else {
		c.creds.NotBefore = time.Unix(rtc.IssuedAt, 0)
	}

	if err = c.Check(); err != nil {
		c.Clear()
		c.Dump()
		return err
	}

	// If cache dump errors, do nothing - just keep going without cache
	c.Dump()
	return nil
}

func (c *mockCredentialsManager) Check() (err error) {
	creds := CredentialsManager{creds: c.creds}
	return creds.Check()
}

func (c *mockCredentialsManager) Valid() bool {
	creds := CredentialsManager{creds: c.creds}
	return creds.Valid()
}

func (c *mockCredentialsManager) Current() bool {
	creds := CredentialsManager{creds: c.creds}
	return creds.Current()
}

func (c *mockCredentialsManager) Refreshable() bool {
	creds := CredentialsManager{creds: c.creds}
	return creds.Refreshable()
}

func (c *mockCredentialsManager) Clear() {
	zeroTime := time.Time{}

	c.creds.AccessToken = ""
	c.creds.RefreshToken = ""
	c.creds.Subject = ""
	c.creds.IssuedAt = zeroTime
	c.creds.ExpiresAt = zeroTime
	c.creds.NotBefore = zeroTime
	c.creds.RefreshBy = zeroTime
}

func (c *mockCredentialsManager) CacheFile() string {
	return ""
}
