package auth

import "errors"

var (
	ErrUnauthenticated = errors.New("request is unauthenticated")
	ErrNoClaims        = errors.New("no claims found on the request context")
)
