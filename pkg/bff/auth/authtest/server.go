/*
Package authtest provides a wrapped httptest.Server that will respond to auth0 requests.
The most common request is related to authentication and token verification, to
authenticate requests to the BFF server, use this package's token generation methods to
create a token that will be validated by the authentication middleware. Note that you
will have to configure the Authenticate middleware to use the correct TLS client.

This module also provides a singleton authtest.Server that can be used on demand from
both tests and live server code by calling the package level functions authtest.Serve()
and authtest.Close respectively. This ensures that tests do not require injection of
the authentication mechanism. The first time that authtest.Serve is called a new server
will be created; and the first time authtest.Close is called, the server will be closed.
Note however that a new server will not be created on subsequent calls, so it's
important to ensure that Close is not called before the tests are complete.
*/
package authtest

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"time"

	"github.com/auth0/go-auth0/management"
	"github.com/golang-jwt/jwt/v4"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"gopkg.in/square/go-jose.v2"
)

const (
	KeyID        = "StyqeY8Kl4Eam28KsUs"
	ClientID     = "a5laOSr0NOX1L53yBaNtumKOoExFxptc"
	ClientSecret = "me4JZSvBvPSnBaM0h0AoXgXPn1VBiBMz0bL7E/sV1isndP9lZ5ptm5NWA9IkKwEb"
	Audience     = "http://localhost"
	Name         = "Leopold Wentzel"
	Email        = "leopold.wentzel@gmail.com"
	UserID       = "test|abcdefg1234567890"
	UserRole     = "Organization Collaborator"
	OrgID        = "b1b9e9b1-9a44-4317-aefa-473971b4df42"
	MainNetVASP  = "87d92fd1-53cf-47d8-85b1-048e8a38ced9"
	TestNetVASP  = "d0082f55-d3ba-4726-a46d-85e3f5a2911f"
	Scope        = "openid profile email"
)

var (
	srv       *Server
	srvErr    error
	srvCreate sync.Once
	srvClose  sync.Once
)

// Serve creates the singleton authtest server if it does not already exist and returns
// it for use in tests and test dependency injection. If creating the server resulted in
// an error then the error is returned. Once Close is called, this method will return
// nil since the server is a singleton and can only be created once. Ensure that Close
// is not called until the tests are complete.
func Serve() (*Server, error) {
	srvCreate.Do(func() {
		srv, srvErr = New()
	})
	return srv, srvErr
}

// Close shuts down the single authtest server and cleans it up. This method should only
// be called once when tests are completed. When the singleton server is shutdown it can
// no longer be created a second time because of the use of sync.Once.
func Close() {
	srvClose.Do(func() {
		if srv != nil {
			srv.Close()
		}
		srv = nil
		srvErr = errors.New("the authtest server has been closed")
	})
}

// Server wraps an httptest.Server to provide a default handler for auth0 requests.
type Server struct {
	srv       *httptest.Server
	mux       *http.ServeMux
	URL       *url.URL
	keys      *rsa.PrivateKey
	users     map[string]*management.User
	userRoles map[string]*management.RoleList
	roles     *management.RoleList
}

// New starts and returns a new Auth0 server using TLS. The caller should call close
// when finished, to shut it down. The server can also issue tokens for authentication.
func New() (s *Server, err error) {
	s = &Server{}

	// Create RSA Private Keys to sign auth tokens with
	if s.keys, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
		return nil, err
	}

	// Create some default users without any associated app metadata
	s.users = NewUsers()

	// Create some default roles for the users
	s.userRoles = NewUserRoles()

	// Create a default role list
	s.roles = NewRoles()

	// Setup routes for the mux
	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/.well-known/openid-configuration", s.OpenIDConfiguration)
	s.mux.HandleFunc("/.well-known/jwks.json", s.JWKS)
	s.mux.HandleFunc("/api/v2/users/"+UserID, s.Users)
	s.mux.HandleFunc("/api/v2/users/"+UserID+"/roles", s.UserRoles)
	s.mux.HandleFunc("/api/v2/roles", s.Roles)

	s.srv = httptest.NewTLSServer(s.mux)
	s.URL, _ = url.Parse(s.srv.URL)
	return s, nil
}

// Config returns an AuthConfig that can be used to setup middleware.
func (s *Server) Config() config.AuthConfig {
	return config.AuthConfig{
		Domain:        s.URL.Host,
		Audience:      Audience,
		ProviderCache: 30 * time.Second,
		ClientID:      ClientID,
		ClientSecret:  ClientSecret,
		Testing:       true,
	}
}

// Client returns the https configured client that can connect to this server.
func (s *Server) Client() *http.Client {
	return s.srv.Client()
}

// Close the server when you're done with your tests!
func (s *Server) Close() {
	s.srv.Close()
}

// NewToken returns a valid token with the specified permissions.
func (s *Server) NewToken(permissions ...string) (tks string, err error) {
	issuer := s.URL.ResolveReference(&url.URL{Path: "/"}).String()

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  issuer,
			Subject: UserID,
			Audience: jwt.ClaimStrings{
				issuer,
				Audience,
			},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
		Email: Email,
		OrgID: OrgID,
		VASPs: map[string]string{
			config.TestNet: TestNetVASP,
			config.MainNet: MainNetVASP,
		},
		Scope:       Scope,
		Permissions: permissions,
	}
	return s.NewTokenWithClaims(claims)
}

// NewTokenWithClaims allows test user to specifically configure their claims.
func (s *Server) NewTokenWithClaims(claims *Claims) (tks string, err error) {
	// Set required claims strings if they are not on the struct
	if claims.Issuer == "" {
		claims.Issuer = s.URL.ResolveReference(&url.URL{Path: "/"}).String()
	}

	if len(claims.Audience) == 0 {
		claims.Audience = jwt.ClaimStrings{
			claims.Issuer, Audience,
		}
	}

	if claims.Subject == "" {
		claims.Subject = UserID
	}

	if claims.IssuedAt == nil && claims.ExpiresAt == nil {
		claims.IssuedAt = jwt.NewNumericDate(time.Now())
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(15 * time.Minute))
	}

	if claims.Scope == "" {
		claims.Scope = Scope
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = KeyID
	return token.SignedString(s.keys)
}

func (s *Server) OpenIDConfiguration(w http.ResponseWriter, r *http.Request) {
	// Create data response to return
	oic := NewOpenIDConfiguration(s.URL)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oic)
}

func (s *Server) JWKS(w http.ResponseWriter, r *http.Request) {
	webkeys := &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{
				Key:       &s.keys.PublicKey,
				KeyID:     KeyID,
				Algorithm: jwt.SigningMethodRS256.Alg(),
			},
		},
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(webkeys)
}

func (s *Server) Users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.GetUser(w, r)
	case http.MethodPatch:
		s.PatchUserAppMetadata(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	// Return the user object from the map
	// TODO: Parse the user id from the request
	if user, ok := s.users[UserID]; ok {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *Server) UserRoles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.ListUserRoles(w, r)
	case http.MethodPost:
		s.AssignUserRoles(w, r)
	case http.MethodDelete:
		s.RemoveUserRoles(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) ListUserRoles(w http.ResponseWriter, r *http.Request) {
	// Return the roles object from the map
	// TODO: Parse the user id from the request
	if roles, ok := s.userRoles[UserID]; ok {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(roles)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *Server) AssignUserRoles(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.userRoles[UserID]; ok {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *Server) RemoveUserRoles(w http.ResponseWriter, r *http.Request) {
	// Note: This does not actually change the state on the server
	if _, ok := s.userRoles[UserID]; ok {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *Server) Roles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.GetRoles(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) GetRoles(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.roles)
}

func (s *Server) PatchUserAppMetadata(w http.ResponseWriter, r *http.Request) {
	// Get the user object from the request
	user := &management.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Replace the user app metadata
	s.users[UserID].AppMetadata = user.AppMetadata

	// Return the user object in the response
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Expose the test user's app metadata to the tests.
func (s *Server) GetUserAppMetadata() map[string]interface{} {
	return s.users[UserID].AppMetadata
}

// Update the test user with unstructured app metadata.
func (s *Server) UseAppMetadata(appdata map[string]interface{}) {
	s.users[UserID].AppMetadata = appdata
}
