package db

import "errors"

var (
	ErrNotFound           = errors.New("key not found in database")
	ErrUnsuccessfulPut    = errors.New("unable to successfully make Put request to trtl")
	ErrUnsuccessfulDelete = errors.New("unable to successfully make Delete request to trtl")
	ErrEmptyAnnouncement  = errors.New("cannot post a zero-valued announcement")
)
