package bff

import "errors"

var (
	ErrNotFound             = errors.New("key not found in database")
	ErrUnsuccessfulPut      = errors.New("unable to successfully make Put request to trtl")
	ErrUnsuccessfulDelete   = errors.New("unable to successfully make Delete request to trtl")
	ErrEmptyAnnouncement    = errors.New("cannot post a zero-valued announcement")
	ErrUnboundedRecent      = errors.New("cannot specify zero-valued not before otherwise announcements fetch is unbounded")
	ErrInvalidUserRole      = errors.New("invalid user role specified")
	ErrUserEmailNotFound    = errors.New("could not find user by email address")
	ErrMultipleEmailUsers   = errors.New("multiple users found by email address")
	ErrDomainAlreadyExists  = errors.New("the specified domain already exists")
	ErrNoVerificationStatus = errors.New("no verification status supplied in context")
)
