package errors

import "errors"

var (
	ErrCorruptedIndex    = errors.New("search indices are invalid")
	ErrCorruptedSequence = errors.New("primary key sequence is invalid")
	ErrDuplicateEntity   = errors.New("entity unique constraints violated")
	ErrEntityNotFound    = errors.New("entity not found")
	ErrIDAlreadySet      = errors.New("record must not have an ID (use update instead)")
	ErrIncompleteRecord  = errors.New("record is missing required fields")
	ErrProtocol          = errors.New("unexpected protocol error")
)
