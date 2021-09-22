package secrets

import "errors"

var (
	ErrSecretNotFound    = errors.New("could not add secret version - not found")
	ErrFileSizeLimit     = errors.New("could not add secret version - file size exceeds limit")
	ErrPermissionsDenied = errors.New("could not add secret version - permissions denied at project level")
)
