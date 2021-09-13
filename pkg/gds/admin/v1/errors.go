package admin

import "errors"

var (
	ErrInvalidResendAction = errors.New("invalid resend action")
	ErrIDRequred           = errors.New("request requires a valid ID to determine endpoint")
)
