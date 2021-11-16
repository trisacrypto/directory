package sectigo_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
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
	s.api, err = New("foo", "supersecret", "CipherTrace EE")
	require.NoError(err)
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

func (s *SectigoTestSuite) TestAuthenticate() {
	require := s.Require()
	m, err := mock.New()
	require.NoError(err)
	defer m.Close()

	err = s.api.Authenticate()
	require.NoError(err)
}

func (s *SectigoTestSuite) TestAuthenticateInvalidCreds() {
	require := s.Require()
	m, err := mock.New()
	require.NoError(err)
	defer m.Close()

	m.Handle(AuthenticateEP, func(c *gin.Context) {
		var (
			in *AuthenticationRequest
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

	s.api, err = New("invalid", "invalid", "CipherTrace EE")
	require.NoError(err)
	err = s.api.Authenticate()
	require.Error(err)
}
