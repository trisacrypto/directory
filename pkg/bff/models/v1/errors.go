package models

import "errors"

var (
	ErrInvalidCollaborator = errors.New("collaborator record is invalid")
	ErrCollaboratorExists  = errors.New("collaborator already exists in organization")
	ErrMaxCollaborators    = errors.New("maximum number of collaborators reached")
)