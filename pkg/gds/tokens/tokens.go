// Package tokens handles the creation and verification of JWT tokens for authentication.
package tokens

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/segmentio/ksuid"
	"google.golang.org/api/idtoken"
)

// Token time constraint constants.
const (
	accessTokenDuration  = 1 * time.Hour
	refreshTokenDuration = 2 * time.Hour
	accessRefreshOverlap = -15 * time.Minute
)

// Global variables that should really not be changed except between major versions.
var (
	signingMethod = jwt.SigningMethodRS256
)

// TokenManager handles the creation and verification of RSA signed JWT tokens. To
// facilitate signing key rollover, TokenManager can accept multiple keys identified by
// a ksuid. JWT tokens generated by token managers include a kid in the header that
// allows the token manager to verify the key with the specified signature. To sign keys
// the token manager will always use the latest private key by ksuid.
//
// When the TokenManager creates tokens it will use JWT standard claims as well as
// extended claims based on Oauth credentials. The standard claims included are exp, nbf
// aud, and sub. The iss claim is optional and would duplicate aud, so it is omitted.
// On token verification, the exp, nbf, and aud claims are validated.
type TokenManager struct {
	audience     string
	currentKeyID ksuid.KSUID
	currentKey   *rsa.PrivateKey
	keys         map[ksuid.KSUID]*rsa.PublicKey
	validate     func(ctx context.Context, idToken, audience string) (*idtoken.Payload, error)
}

// Claims implements custom claims for the GDS application to hold user data provided
// from external openid sources. It also embeds the standard JWT claims.
type Claims struct {
	jwt.StandardClaims
	Domain  string `json:"hd,omitempty"`
	Email   string `json:"email,omitempty"`
	Name    string `json:"name,omitempty"`
	Picture string `json:"picture,omitempty"`
}

// New creates a TokenManager with the specified keys which should be a mapping of KSUID
// strings to paths to files that contain PEM encoded RSA private keys. This input is
// specifically designed for the config environment variable so that keys can be loaded
// from k8s or vault secrets that are mounted as files on disk.
func New(keys map[string]string, audience string) (tm *TokenManager, err error) {
	tm = &TokenManager{
		keys:     make(map[ksuid.KSUID]*rsa.PublicKey),
		audience: audience,
		validate: idtoken.Validate,
	}

	for kid, path := range keys {
		// Parse the key id
		var keyID ksuid.KSUID
		if keyID, err = ksuid.Parse(kid); err != nil {
			return nil, fmt.Errorf("could not parse kid %q for path %s: %s", kid, path, err)
		}

		// Load the keys from disk
		var data []byte
		if data, err = ioutil.ReadFile(path); err != nil {
			return nil, fmt.Errorf("could not read kid %s from %s: %s", kid, path, err)
		}

		var key *rsa.PrivateKey
		if key, err = jwt.ParseRSAPrivateKeyFromPEM(data); err != nil {
			return nil, fmt.Errorf("could not parse RSA private key kid %s from %s: %s", kid, path, err)
		}

		// Add the key to the key map
		tm.keys[keyID] = &key.PublicKey

		// Set the current key if it is the latest key
		if tm.currentKey == nil || keyID.Time().After(tm.currentKeyID.Time()) {
			tm.currentKey = key
			tm.currentKeyID = keyID
		}
	}

	return tm, nil
}

// Verify an access or a refresh token after parsing and return its claims.
func (tm *TokenManager) Verify(tks string) (claims *Claims, err error) {
	var token *jwt.Token
	if token, err = jwt.ParseWithClaims(tks, &Claims{}, tm.keyFunc); err != nil {
		return nil, err
	}

	var ok bool
	if claims, ok = token.Claims.(*Claims); ok && token.Valid {
		if !claims.VerifyAudience(tm.audience, true) {
			return nil, fmt.Errorf("invalid audience %q", claims.Audience)
		}

		return claims, nil
	}

	return nil, fmt.Errorf("could not parse or verify GDS claims from %T", token.Claims)
}

// Parse an access or refresh token verifying its signature but without verifying its
// claims. This ensures that valid JWT tokens are still accepted but claims can be
// handled on a case-by-case basis; for example by validating an expired access token
// during reauthentication.
func (tm *TokenManager) Parse(tks string) (claims *Claims, err error) {
	parser := &jwt.Parser{SkipClaimsValidation: true}
	claims = &Claims{}
	if _, err = parser.ParseWithClaims(tks, claims, tm.keyFunc); err != nil {
		return nil, err
	}
	return claims, nil
}

// Sign an access or refresh token and return the token string.
func (tm *TokenManager) Sign(token *jwt.Token) (tks string, err error) {
	// Sanity check to prevent nil panics.
	if tm.currentKey == nil || tm.currentKeyID.IsNil() {
		return "", errors.New("token manager not initialized with signing keys")
	}

	// Add the kid (key id - this is the standard 3 letter JWT name) to the header.
	token.Header["kid"] = tm.currentKeyID.String()

	// Return the signed string
	return token.SignedString(tm.currentKey)
}

// CreateAccessToken from the verified Google credential payload or from an previous
// token if the access token is being reauthorized from previous credentials. Note that
// the returned token only contains the claims and is unsigned.
func (tm *TokenManager) CreateAccessToken(creds interface{}) (_ *jwt.Token, err error) {
	// Create the claims for the access token, using access token defaults.
	now := time.Now()
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.NewString(), // ID is randomly generated and shared between access and refresh tokens.
			Audience:  tm.audience,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(accessTokenDuration).Unix(),
		},
	}

	// Populate the claims from the credentials based on type.
	switch t := creds.(type) {
	case *idtoken.Payload:
		// This case is extracting the claims from a verified idtoken payload from Google,
		// e.g. this is extracting the claims from from the Google OAuth credential.
		if err = claims.extractClaims(t.Claims); err != nil {
			return nil, err
		}
	case jwt.MapClaims, map[string]interface{}:
		// This is a generic case that exercises the extractClaims function, it is currently
		// only used for testing, but might also be used if we have different credential sources.
		if err = claims.extractClaims(t.(map[string]interface{})); err != nil {
			return nil, err
		}
	case *Claims:
		// This case extracts the claims from a previously issued access token and is used
		// on reauthenticate to issue a new access token from a soon-to-expire access token.
		claims.Domain = t.Domain
		claims.Email = t.Email
		claims.Name = t.Name
		claims.Picture = t.Picture
		claims.Subject = t.Subject
	default:
		return nil, fmt.Errorf("cannot create access token from %T", t)
	}

	return jwt.NewWithClaims(signingMethod, claims), nil
}

// CreateRefreshToken from the Access token claims with predefined expiration. Note that
// the returned token only contains the claims and is unsigned.
func (tm *TokenManager) CreateRefreshToken(accessToken *jwt.Token) (refreshToken *jwt.Token, err error) {
	accessClaims, ok := accessToken.Claims.(*Claims)
	if !ok {
		return nil, errors.New("could not retrieve GDS claims from access token")
	}

	// Create claims for the refresh token from the access token defaults.
	// Note that refresh token claims are GDS claims but do not have credentials, which
	// means the refresh token can also be parsed with standard claims.
	// TODO: should we make this a refresh-specific audience or subject?
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			Id:        accessClaims.Id, // ID is randomly generated and shared between access and refresh tokens.
			Audience:  accessClaims.Audience,
			Issuer:    accessClaims.Issuer,
			Subject:   accessClaims.Subject,
			IssuedAt:  accessClaims.IssuedAt,
			NotBefore: time.Unix(accessClaims.ExpiresAt, 0).Add(accessRefreshOverlap).Unix(),
			ExpiresAt: time.Unix(accessClaims.IssuedAt, 0).Add(refreshTokenDuration).Unix(),
		},
	}

	return jwt.NewWithClaims(signingMethod, claims), nil
}

// Keys returns the map of ksuid to public key for use externally.
func (tm *TokenManager) Keys() map[ksuid.KSUID]*rsa.PublicKey {
	return tm.keys
}

// CurrentKey returns the ksuid of the current key being used to sign tokens.
func (tm *TokenManager) CurrentKey() ksuid.KSUID {
	return tm.currentKeyID
}

// Validate the given token using the provided audience and return the token's payload.
// This method provides a convenient way for tests to circumvent Google's specific
// validation logic in order to test successful authentication.
func (tm *TokenManager) Validate(ctx context.Context, idToken, audience string) (*idtoken.Payload, error) {
	if tm.validate == nil {
		return nil, errors.New("no token validation function configured")
	}
	return tm.validate(ctx, idToken, audience)
}

// keyFunc is an jwt.KeyFunc that selects the RSA public key from the list of managed
// internal keys based on the kid in the token header. If the kid does not exist an
// error is returned and the token will not be able to be verified.
func (tm *TokenManager) keyFunc(token *jwt.Token) (key interface{}, err error) {
	// Per JWT security notice: do not forget to validate alg is expected
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	// Fetch the kid from the header
	kid, ok := token.Header["kid"]
	if !ok {
		return nil, errors.New("token does not have kid in header")
	}

	// Parse the kid
	var keyID ksuid.KSUID
	if keyID, err = ksuid.Parse(kid.(string)); err != nil {
		return nil, fmt.Errorf("could not parse kid: %s", err)
	}

	// Fetch the key from the list of managed keys
	if key, ok = tm.keys[keyID]; !ok {
		return nil, errors.New("unknown signing key")
	}
	return key, nil
}

func (c *Claims) extractClaims(o map[string]interface{}) (err error) {
	// Extract required claims
	for _, key := range []string{"hd", "email"} {
		if err = c.extractClaim(key, o); err != nil {
			return fmt.Errorf("missing required claim: %s", err)
		}
	}

	// Extract optional claims (do not return an error)
	for _, key := range []string{"name", "picture", "sub"} {
		c.extractClaim(key, o)
	}

	return nil
}

func (c *Claims) extractClaim(key string, o map[string]interface{}) error {
	// Fetch the claim from the map
	val, ok := o[key]
	if !ok || val == nil {
		return fmt.Errorf("no claim %q found in token", key)
	}

	var vals string
	if vals, ok = val.(string); !ok {
		return fmt.Errorf("claim %q is not a string", key)
	}

	if vals == "" {
		return fmt.Errorf("claim %q is empty", key)
	}

	switch key {
	case "sub":
		c.Subject = vals
	case "hd":
		c.Domain = vals
	case "email":
		c.Email = vals
	case "name":
		c.Name = vals
	case "picture":
		c.Picture = vals
	default:
		return fmt.Errorf("unhandled claim %q", key)
	}

	return nil
}
