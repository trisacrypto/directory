package tokens

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/segmentio/ksuid"
)

// MockTokenManager creates a new TokenManager with a randomly generated RSA key to be
// used for testing external code that depends on the Token Manager.
func MockTokenManager() (tm *TokenManager, err error) {
	tm = &TokenManager{
		audience: "http://localhost",
		subject:  "testing",
		issuer:   "http://localhost",
		keys:     make(map[ksuid.KSUID]*rsa.PublicKey),
	}

	var key *rsa.PrivateKey
	if key, err = rsa.GenerateKey(rand.Reader, 1024); err != nil {
		return nil, err
	}

	tm.currentKeyID = ksuid.New()
	tm.currentKey = key
	tm.keys[tm.currentKeyID] = &key.PublicKey
	return tm, nil
}
