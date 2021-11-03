package sectigo

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSectigo(t *testing.T) {
	suite.Run(t, new(SectigoTestSuite))
}

type SectigoTestSuite struct {
	suite.Suite
}

func (s *SectigoTestSuite) BeforeTest(suiteName, testName string) {
	mockCredentialsCache = Credentials{
		Username: "foo",
		Password: "supersecret",
	}
	mockBackend = &mockServer{
		users: map[string]*mockUser{
			"foo": {
				username:   "foo",
				password:   "supersecret",
				authorized: true,
			},
		},
		access:  map[string]string{},
		refresh: map[string]string{},
	}
}

func (s *SectigoTestSuite) AfterTest(suiteName, testName string) {
	mockCredentialsCache = Credentials{}
}

func (s *SectigoTestSuite) TestCredsCopy() {
	require := s.Require()
	api, err := New("foo", "supersecret", "CipherTrace EE")
	require.NoError(err)

	// Ensure that creds are copied and are not the same object
	creds := api.Creds()
	require.NotEqual(&api.creds, &creds)

	require.Equal(api.creds.Creds().Username, creds.Username)
	creds.Username = "superbunny"
	require.NotEqual(api.creds.Creds().Username, creds.Username)
	require.Equal(api.creds.Creds().Username, "foo")
}

func (s *SectigoTestSuite) TestAuthenticate() {
	require := s.Require()
	api, err := NewMock("foo", "supersecret", "CipherTrace EE")
	require.NoError(err)

	// Authenticate a user
	err = api.Authenticate()
	require.NoError(err)
}
