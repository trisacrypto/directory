package authtest

import (
	"net/url"

	"github.com/golang-jwt/jwt/v4"
)

// Claims must be defined here both to ensure we can use jwt and to ensure there are no
// recursive imports. That means this claims struct MUST be kept up to date with the
// auth.Claims struct that uses this package for testing.
type Claims struct {
	jwt.RegisteredClaims
	Email       string            `json:"https://vaspdirectory.net/email"`
	OrgID       string            `json:"https://vaspdirectory.net/orgid"`
	VASPs       map[string]string `json:"https://vaspdirectory.net/vasps"`
	Scope       string            `json:"scope"`
	Permissions []string          `json:"permissions"`
}

type OpenIDConfiguration struct {
	Issuer                        string   `json:"issuer"`
	AuthorizationEP               string   `json:"authorization_endpoint"`
	TokenEP                       string   `json:"token_endpoint"`
	DeviceAuthorizationEP         string   `json:"device_authorization_endpoint"`
	UserInfoEP                    string   `json:"userinfo_endpoint"`
	MFAChallengeEP                string   `json:"mfa_challenge_endpoint"`
	JWKSURI                       string   `json:"jwks_uri"`
	RegistrationEP                string   `json:"registration_endpoint"`
	RevocationEP                  string   `json:"revocation_endpoint"`
	ScopesSupported               []string `json:"scopes_supported"`
	ResponseTypesSupported        []string `json:"response_types_supported"`
	CodeChallengeMethodsSupported []string `json:"code_challenge_methods_supported"`
	ResponseModesSupported        []string `json:"response_modes_supported"`
	SubjectTypesSupported         []string `json:"subject_types_supported"`
	IDTokenSigningAlgValues       []string `json:"id_token_signing_alg_values_supported"`
	TokenEndpointAuthMethods      []string `json:"token_endpoint_auth_methods_supported"`
	ClaimsSupported               []string `json:"claims_supported"`
	RequestURIPArameterSupported  bool     `json:"request_uri_parameter_supported"`
}

func NewOpenIDConfiguration(u *url.URL) *OpenIDConfiguration {
	return &OpenIDConfiguration{
		Issuer:                        u.ResolveReference(&url.URL{Path: "/"}).String(),
		AuthorizationEP:               u.ResolveReference(&url.URL{Path: "/authorize"}).String(),
		TokenEP:                       u.ResolveReference(&url.URL{Path: "/oauth/token"}).String(),
		DeviceAuthorizationEP:         u.ResolveReference(&url.URL{Path: "/oauth/device/code"}).String(),
		UserInfoEP:                    u.ResolveReference(&url.URL{Path: "/userinfo"}).String(),
		MFAChallengeEP:                u.ResolveReference(&url.URL{Path: "/mfa/challenge"}).String(),
		JWKSURI:                       u.ResolveReference(&url.URL{Path: "/.well-known/jwks.json"}).String(),
		RegistrationEP:                u.ResolveReference(&url.URL{Path: "/oidc/register"}).String(),
		RevocationEP:                  u.ResolveReference(&url.URL{Path: "/oauth/revoke"}).String(),
		ScopesSupported:               []string{"openid", "profile", "email"},
		ResponseTypesSupported:        []string{"token", "id_token"},
		CodeChallengeMethodsSupported: []string{"S256", "plain"},
		ResponseModesSupported:        []string{"query", "fragment", "form_post"},
		SubjectTypesSupported:         []string{"public"},
		IDTokenSigningAlgValues:       []string{"HS256", "RS256"},
		TokenEndpointAuthMethods:      []string{"client_secret_basic", "client_secret_post"},
		ClaimsSupported:               []string{"aud", "email", "exp", "iat", "iss", "sub"},
		RequestURIPArameterSupported:  false,
	}
}
