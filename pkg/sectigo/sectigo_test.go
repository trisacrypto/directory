package sectigo_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/sectigo"
	"github.com/trisacrypto/directory/pkg/sectigo/mock"
)

func TestSectigo(t *testing.T) {
	suite.Run(t, new(SectigoTestSuite))
}

type SectigoTestSuite struct {
	suite.Suite
	api *sectigo.Sectigo
}

func (s *SectigoTestSuite) BeforeTest(suiteName, testName string) {
	var err error
	require := s.Require()
	conf := sectigo.Config{
		Username: "foo",
		Password: "supersecret",
		Profile:  "CipherTrace EE",
		Testing:  true,
	}
	require.NoError(mock.Start())
	s.api, err = sectigo.New(conf)
	require.NoError(err)
}

func (s *SectigoTestSuite) AfterTest(suiteName, testName string) {
	creds := s.api.Creds()
	if path := creds.CacheFile(); path != "" {
		os.RemoveAll(path)
	}
	mock.Stop()
}

func (s *SectigoTestSuite) TestCredsCopy() {
	require := s.Require()

	// Test the internal Sectigo credentials
	creds := s.api.Creds()
	require.Equal("foo", creds.Username)
	require.Equal("supersecret", creds.Password)

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
		{name: sectigo.AuthenticateEP, f: s.authenticate},
		{name: sectigo.RefreshEP, f: s.refresh}, // must be called after authenticate
		{name: sectigo.CreateSingleCertBatchEP, f: s.createSingleCertBatch},
		{name: sectigo.UploadCSREP, f: s.uploadCSRBatch},
		{name: sectigo.BatchDetailEP, f: s.batchDetail},
		{name: sectigo.BatchProcessingInfoEP, f: s.processingInfo},
		{name: sectigo.BatchStatusEP, f: s.batchStatus},
		{name: sectigo.DownloadEP, f: s.download},
		{name: sectigo.DevicesEP, f: s.licensesUsed},
		{name: sectigo.UserAuthoritiesEP, f: s.userAuthorities},
		{name: sectigo.AuthorityBalanceAvailableEP, f: s.authorityAvailableBalance},
		{name: sectigo.CurrentUserOrganizationEP, f: s.organization},
		{name: sectigo.ProfilesEP, f: s.profiles},
		{name: sectigo.ProfileParametersEP, f: s.profileParams},
		{name: sectigo.ProfileDetailEP, f: s.profileDetail},
		{name: sectigo.FindCertificateEP, f: s.findCertificate},
		{name: sectigo.RevokeCertificateEP, f: s.revokeCertificate},
	}
	for _, t := range tests {
		s.T().Run(t.name, t.f)
	}
	calls := mock.Get().GetCalls()
	for endpoint, count := range calls {
		require.Equal(1, count, fmt.Errorf("wrong number of calls to endpoint %s", endpoint))
	}
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
	rep, err := s.api.CreateSingleCertBatch(42, "foo", map[string]string{"foo": "bar"})
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
	dir, err := ioutil.TempDir("", "download-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	rep, err := s.api.Download(42, dir)
	require.NoError(t, err)
	require.Equal(t, filepath.Join(dir, "42.zip"), rep)
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
	rep, err := s.api.AuthorityAvailableBalance(42)
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
	rep, err := s.api.ProfileParams(42)
	require.NoError(t, err)
	require.NotNil(t, rep)
}

func (s *SectigoTestSuite) profileDetail(t *testing.T) {
	rep, err := s.api.ProfileDetail(42)
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

	mock.Handle(sectigo.AuthenticateEP, func(c *gin.Context) {
		var (
			in *sectigo.AuthenticationRequest
		)
		if err := c.BindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		if in.Username != "foo" || in.Password != "supersecret" {
			c.JSON(http.StatusUnauthorized, "invalid credentials")
			return
		}

		c.JSON(http.StatusInternalServerError, "how did we get here?")
	})

	conf := sectigo.Config{
		Username: "invalid",
		Password: "invalid",
		Testing:  true,
	}
	var err error
	s.api, err = sectigo.New(conf)
	require.NoError(err)
	require.EqualError(s.api.Authenticate(), sectigo.ErrInvalidCredentials.Error())
}
