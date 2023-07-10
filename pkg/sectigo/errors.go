package sectigo

import (
	"errors"
	"fmt"
)

// Standard errors issued by the Sectigo client.
var (
	ErrNotAuthenticated     = errors.New("sectigo client error: not authenticated")
	ErrCredentialsMismatch  = errors.New("sectigo client error: requires both username and password")
	ErrTokensMismatch       = errors.New("sectigo client error: both access and refresh tokens required")
	ErrNoCredentials        = errors.New("sectigo client error: no API access credentials")
	ErrInvalidCredentials   = errors.New("sectigo client error: could not authenticate credentials")
	ErrNotAuthorized        = errors.New("sectigo client error: user is not authorized for this endpoint")
	ErrTokensExpired        = errors.New("sectigo client error: access and refresh tokens have expired")
	ErrInvalidClaims        = errors.New("sectigo client error: jwt claims do not have required timestamps")
	ErrMustUseTLSAuth       = errors.New("sectigo client error: account requires TLS client authentication")
	ErrPKCSPasswordRequired = errors.New("sectigo client error: pkcs12 password required for cert params")
)

// APIError is unmarshalled from the JSON response of the Sectigo API and implements
// the error interface to correctly return error messages.
type APIError struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	ErrorCode int    `json:"errorCode"`
	Timestamp int    `json:"timestamp"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("%d %s %d %d", e.Status, e.Message, e.ErrorCode, e.Timestamp)
}
