package tokens

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"

	"github.com/segmentio/ksuid"
	"google.golang.org/api/idtoken"
)

// MockTokenManager creates a new TokenManager with a randomly generated RSA key to be
// used for testing external code that depends on the Token Manager.
func MockTokenManager() (tm *TokenManager, err error) {
	tm = &TokenManager{
		audience: "http://localhost",
		keys:     make(map[ksuid.KSUID]*rsa.PublicKey),
	}
	tm.validate = tm.mockValidate

	var key *rsa.PrivateKey
	if key, err = rsa.GenerateKey(rand.Reader, 1024); err != nil {
		return nil, err
	}

	tm.currentKeyID = ksuid.New()
	tm.currentKey = key
	tm.keys[tm.currentKeyID] = &key.PublicKey
	return tm, nil
}

// mockValidate is a mock for Google's token validation which enables authentication
// testing by internally parsing the given token and injecting the claims into the
// returned payload. This will only validate badly constructed tokens and tokens where
// the supplied audience does not match the audience used to sign the token.
func (tm *TokenManager) mockValidate(ctx context.Context, idToken string, audience string) (payload *idtoken.Payload, err error) {
	var claims *Claims
	if claims, err = tm.Parse(idToken); err != nil {
		return nil, err
	}

	if claims.Audience != audience {
		return nil, errors.New("audience in token does not match given audience")
	}

	payload = &idtoken.Payload{
		Issuer:   claims.Issuer,
		Audience: claims.Audience,
		Expires:  claims.ExpiresAt,
		IssuedAt: claims.IssuedAt,
		Subject:  claims.Subject,
		Claims:   make(map[string]interface{}),
	}

	// Marshal the claims to JSON
	var data []byte
	if data, err = json.Marshal(claims); err != nil {
		return nil, err
	}

	// Unmarshal the claims back into the payload
	if err = json.Unmarshal(data, &payload.Claims); err != nil {
		return nil, err
	}

	return payload, nil
}
