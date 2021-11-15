package sectigo

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
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
	// Ensure that creds are copied and are not the same object
	creds := s.api.Creds()
	require.NotEqual(&s.api.creds, &creds)

	require.Equal(s.api.creds.Username, creds.Username)
	creds.Username = "superbunny"
	require.NotEqual(s.api.creds.Username, creds.Username)
	require.Equal(s.api.creds.Username, "foo")
}

func (s *SectigoTestSuite) TestAuthenticate() {
	require := s.Require()
	m, err := NewMockServer(nil)
	require.NoError(err)
	defer m.server.Close()

	err = s.api.Authenticate()
	require.NoError(err)
}

func (s *SectigoTestSuite) TestAuthenticateInvalidCreds() {
	require := s.Require()
	m, err := NewMockServer(map[string]*mockHandlerFunc{
		endpoints[authenticateEP].Path: {
			method: http.MethodPost,
			handlerFunc: func(c *gin.Context) {
				var (
					in      *AuthenticationRequest
					access  string
					refresh string
					err     error
				)
				if err := c.BindJSON(&in); err != nil {
					c.JSON(http.StatusBadRequest, err)
					return
				}
				if in.Username != "foo" || in.Password != "supersecret" {
					c.JSON(http.StatusUnauthorized, "invalid credentials")
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
			},
		},
	})
	require.NoError(err)
	defer m.server.Close()

	s.api.creds = &Credentials{
		Username: "invalid",
		Password: "invalid",
	}
	err = s.api.Authenticate()
	require.Error(err)
}
