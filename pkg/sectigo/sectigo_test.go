package sectigo_test

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	. "github.com/trisacrypto/directory/pkg/sectigo"
	"github.com/trisacrypto/directory/pkg/sectigo/mock"
)

func TestSectigo(t *testing.T) {
	suite.Run(t, new(SectigoTestSuite))
}

type SectigoTestSuite struct {
	suite.Suite
	api *Sectigo
}

func (s *SectigoTestSuite) BeforeTest(suiteName, testName string) {
	var err error
	require := s.Require()

	// Note: the username and password will be set to sectigo.MockUsername and
	// sectigo.MockPassword because config.Testing==true.
	conf := Config{
		Profile:     "CipherTrace EE",
		Environment: "testing",
	}
	require.NoError(mock.Start(conf.Profile))
	s.api, err = New(conf)
	require.NoError(err)
}

func (s *SectigoTestSuite) AfterTest(suiteName, testName string) {
	creds := s.api.Creds()
	if path := creds.CacheFile(); path != "" {
		os.RemoveAll(path)
	}
	mock.Stop()
}

func (s *SectigoTestSuite) TestSortedParams() {
	require := s.Require()
	for key, params := range Profiles {
		require.True(sort.StringsAreSorted(params), "the %s params are not sorted", key)
	}
}

func (s *SectigoTestSuite) TestCredsCopy() {
	require := s.Require()

	// Test the internal Sectigo credentials
	creds := s.api.Creds()
	require.Equal(MockUsername, creds.Username)
	require.Equal(MockPassword, creds.Password)

	// Ensure that creds are copied and are not the same object
	creds.Username = "superbunny"
	creds.Password = "knockknock"
	require.NotEqual(&creds, s.api.Creds())

	orig := s.api.Creds()
	require.NotEqual(creds.Username, orig.Username)
	require.NotEqual(creds.Password, orig.Password)
}

func (s *SectigoTestSuite) TestSuccessfulCalls() {
	require := s.Require()
	tests := []struct {
		name string
		f    func(t *testing.T)
	}{
		{name: AuthenticateEP, f: s.authenticate},
		{name: RefreshEP, f: s.refresh}, // must be called after authenticate
		{name: CreateSingleCertBatchEP, f: s.createSingleCertBatch},
		{name: UploadCSREP, f: s.uploadCSRBatch},
		{name: BatchDetailEP, f: s.batchDetail},
		{name: BatchProcessingInfoEP, f: s.processingInfo},
		{name: BatchStatusEP, f: s.batchStatus},
		{name: DownloadEP, f: s.download},
		{name: DevicesEP, f: s.licensesUsed},
		{name: UserAuthoritiesEP, f: s.userAuthorities},
		{name: AuthorityBalanceAvailableEP, f: s.authorityAvailableBalance},
		{name: CurrentUserOrganizationEP, f: s.organization},
		{name: ProfilesEP, f: s.profiles},
		{name: ProfileParametersEP, f: s.profileParams},
		{name: ProfileDetailEP, f: s.profileDetail},
		{name: FindCertificateEP, f: s.findCertificate},
		{name: RevokeCertificateEP, f: s.revokeCertificate},
	}
	for _, t := range tests {
		s.T().Run(t.name, t.f)
	}
	calls := mock.Get().GetCalls()
	calls.Range(func(endpoint, count interface{}) bool {
		num, ok := count.(int)
		require.True(ok, "unexpected type in call map, expected int")
		require.Equal(1, num, fmt.Errorf("wrong number of calls to endpoint %s", endpoint))
		return true
	})
}

func (s *SectigoTestSuite) authenticate(t *testing.T) {
	require.NoError(t, s.api.Authenticate())
}

// Note: refresh has an external dependency on the refresh token being set in the
// sectigo client Credentials object.
func (s *SectigoTestSuite) refresh(t *testing.T) {
	require.NoError(t, s.api.Refresh())
}

func (s *SectigoTestSuite) createSingleCertBatch(t *testing.T) {
	rep, err := s.api.CreateSingleCertBatch(1, "foo", map[string]string{
		"commonName":     "foo",
		"dNSName":        "foo.example.com",
		"pkcs12Password": "bar",
	})
	require.NoError(t, err)
	require.NotNil(t, rep)
}

func (s *SectigoTestSuite) uploadCSRBatch(t *testing.T) {
	rep, err := s.api.UploadCSRBatch(42, "foo", []byte("foo"), map[string]string{"foo": "bar"})
	require.NoError(t, err)
	require.NotNil(t, rep)
}

func (s *SectigoTestSuite) batchDetail(t *testing.T) {
	rep, err := s.api.BatchDetail(42)
	require.NoError(t, err)
	require.NotNil(t, rep)
}

func (s *SectigoTestSuite) batchStatus(t *testing.T) {
	rep, err := s.api.BatchStatus(42)
	require.NoError(t, err)
	require.NotEmpty(t, rep)
	require.Equal(t, "READY_FOR_DOWNLOAD", rep)
}

func (s *SectigoTestSuite) processingInfo(t *testing.T) {
	rep, err := s.api.ProcessingInfo(42)
	require.NoError(t, err)
	require.NotNil(t, rep)
}

func (s *SectigoTestSuite) download(t *testing.T) {
	dir := t.TempDir()
	rep, err := s.api.Download(42, dir)
	require.NoError(t, err)
	require.Equal(t, filepath.Join(dir, "certificate.zip"), rep)
	require.FileExists(t, rep)
}

func (s *SectigoTestSuite) licensesUsed(t *testing.T) {
	rep, err := s.api.LicensesUsed()
	require.NoError(t, err)
	require.NotNil(t, rep)
}

func (s *SectigoTestSuite) userAuthorities(t *testing.T) {
	rep, err := s.api.UserAuthorities()
	require.NoError(t, err)
	require.NotEmpty(t, rep)
	for _, r := range rep {
		require.NotNil(t, r)
	}
}

func (s *SectigoTestSuite) authorityAvailableBalance(t *testing.T) {
	rep, err := s.api.AuthorityAvailableBalance(1)
	require.NoError(t, err)
	require.Greater(t, rep, 0)
}

func (s *SectigoTestSuite) profiles(t *testing.T) {
	rep, err := s.api.Profiles()
	require.NoError(t, err)
	require.NotEmpty(t, rep)
	for _, r := range rep {
		require.NotNil(t, r)
	}
}

func (s *SectigoTestSuite) profileParams(t *testing.T) {
	id, err := strconv.Atoi(ProfileIDCipherTraceEE)
	require.NoError(t, err)
	rep, err := s.api.ProfileParams(id)
	require.NoError(t, err)
	require.NotNil(t, rep)
}

func (s *SectigoTestSuite) profileDetail(t *testing.T) {
	id, err := strconv.Atoi(ProfileIDCipherTraceEE)
	require.NoError(t, err)
	rep, err := s.api.ProfileDetail(id)
	require.NoError(t, err)
	require.NotNil(t, rep)
}

func (s *SectigoTestSuite) organization(t *testing.T) {
	rep, err := s.api.Organization()
	require.NoError(t, err)
	require.NotNil(t, rep)
}

func (s *SectigoTestSuite) findCertificate(t *testing.T) {
	rep, err := s.api.FindCertificate("foo", "12345")
	require.NoError(t, err)
	require.NotNil(t, rep)
}

func (s *SectigoTestSuite) revokeCertificate(t *testing.T) {
	err := s.api.RevokeCertificate(42, 0, "12345")
	require.NoError(t, err)
}

func (s *SectigoTestSuite) TestAuthenticateInvalidCreds() {
	require := s.Require()

	mock.Handle(AuthenticateEP, func(c *gin.Context) {
		var (
			in *AuthenticationRequest
		)
		if err := c.BindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		if in.Username != MockUsername || in.Password != MockPassword {
			c.JSON(http.StatusUnauthorized, "invalid credentials")
			return
		}

		c.JSON(http.StatusInternalServerError, "how did we get here?")
	})

	conf := Config{
		Username:    "invalid",
		Password:    "invalid",
		Profile:     ProfileCipherTraceEE,
		Environment: "testing",
	}
	var err error
	s.api, err = New(conf)
	require.NoError(err)
	require.EqualError(s.api.Authenticate(), ErrInvalidCredentials.Error())
}
