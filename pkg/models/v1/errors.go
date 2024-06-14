package models

import "errors"

var (
	ErrNoEmailAddress    = errors.New("email record requires an email address")
	ErrVerifiedInvalid   = errors.New("a verified email must have a verified on timestamp and no token")
	ErrUnverifiedInvalid = errors.New("an unverified email must have a token and no verified on timestamp")
	ErrNoLogReason       = errors.New("cannot match on empty reason string")
	ErrInvalidWindow     = errors.New("time window must be a positive number of days")
)
