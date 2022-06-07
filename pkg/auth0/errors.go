package auth0

import (
	"errors"
	"fmt"
)

var (
	ErrNotAuthenticated = errors.New("auth0 client not authenticated")
)

// APIError is the JSON value returned from the server where applicable.
type APIError struct {
	StatusCode int    `json:"statusCode"`
	Status     string `json:"error"`
	Message    string `json:"message"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("[%d] %s: %s", e.StatusCode, e.Status, e.Message)
	}
	return fmt.Sprintf("[%d] %s", e.StatusCode, e.Status)
}
