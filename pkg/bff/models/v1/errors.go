package models

import "errors"

var (
	ErrInvalidCollaborator = errors.New("collaborator record is invalid")
	ErrCollaboratorExists  = errors.New("collaborator already exists in organization")
)
