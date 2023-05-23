package models

import "errors"

var (
	ErrInvalidField        = errors.New("invalid field")
	ErrNegativeValue       = errors.New("field cannot be negative")
	ErrMissingField        = errors.New("missing required field")
	ErrInvalidCollaborator = errors.New("collaborator record is invalid")
	ErrCollaboratorExists  = errors.New("collaborator already exists in organization")
	ErrMaxCollaborators    = errors.New("maximum number of collaborators reached")
)
