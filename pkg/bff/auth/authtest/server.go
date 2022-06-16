/*
Package authtest provides a wrapped httptest.Server that will respond to auth0 requests.
The most common request is related to authentication and token verification, to
authenticate requests to the BFF server, use this package's token generation methods to
create a token that will be validated by the authentication middleware. Note that you
will have to configure the Authenticate middleware to use the correct TLS client.
*/
package authtest

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"gopkg.in/square/go-jose.v2"
)

const (
	KeyID        = "StyqeY8Kl4Eam28KsUs"
	ClientID     = "a5laOSr0NOX1L53yBaNtumKOoExFxptc"
	ClientSecret = "me4JZSvBvPSnBaM0h0AoXgXPn1VBiBMz0bL7E/sV1isndP9lZ5ptm5NWA9IkKwEb"
	Audience     = "http://localhost:4437/"
	Email        = "leopold.wentzel@gmail.com"
	UserID       = "test|abcdefg1234567890"
	OrgID        = "b1b9e9b1-9a44-4317-aefa-473971b4df42"
	MainNetVASP  = "87d92fd1-53cf-47d8-85b1-048e8a38ced9"
	TestNetVASP  = "d0082f55-d3ba-4726-a46d-85e3f5a2911f"
	Scope        = "openid profile email"
)

// Server wraps an httptest.Server to provide a default handler for auth0 requests.
type Server struct {
	srv  *httptest.Server
	mux  *http.ServeMux
	URL  *url.URL
	keys *rsa.PrivateKey
}

// New starts and returns a new Auth0 server using TLS. The caller should call close
// when finished, to shut it down. The server can also issue tokens for authentication.
func New() (s *Server, err error) {
	s = &Server{}

	// Create RSA Private Keys to sign auth tokens with
	if s.keys, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
		return nil, err
	}

	// Setup routes for the mux
	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/.well-known/openid-configuration", s.OpenIDConfiguration)
	s.mux.HandleFunc("/.well-known/jwks.json", s.JWKS)

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
			"testnet": TestNetVASP,
			"mainnet": MainNetVASP,
		},
		Scope:       Scope,
		Permissions: permissions,
	}
	return s.NewTokenWithClaims(claims)
}

// NewTokenWithClaims allows test user to specifically configure their claims.
func (s *Server) NewTokenWithClaims(claims *Claims) (tks string, err error) {
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
