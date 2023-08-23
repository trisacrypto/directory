package ensign

import "errors"

var (
	ErrMissingClientID     = errors.New("missing client id")
	ErrMissingClientSecret = errors.New("missing client secret")
	ErrMissingEndpoint     = errors.New("missing endpoint")
	ErrMissingAuthURL      = errors.New("missing auth url")
)
