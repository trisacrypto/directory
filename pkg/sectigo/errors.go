package sectigo

import "errors"

// Standard errors issued by the Sectigo client.
var (
	ErrNotAuthenticated     = errors.New("not authenticated")
	ErrCredentialsMismatch  = errors.New("requires both username and password")
	ErrTokensMismatch       = errors.New("both access and refresh tokens required")
	ErrNoCredentials        = errors.New("no API access credentials")
	ErrInvalidCredentials   = errors.New("could not authenticate credentials")
	ErrNotAuthorized        = errors.New("user is not authorized for this endpoint")
	ErrTokensExpired        = errors.New("access and refresh tokens have expired")
	ErrInvalidClaims        = errors.New("jwt claims do not have required timestamps")
	ErrMustUseTLSAuth       = errors.New("account requires TLS client authentication")
	ErrPKCSPasswordRequired = errors.New("pkcs12 password required for cert params")
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
	return e.Message
}
