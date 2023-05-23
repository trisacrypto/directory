package models

import "errors"

var (
	ErrMissingField            = errors.New("missing required field")
	ErrTooShort                = errors.New("field is too short")
	ErrInvalidEmail            = errors.New("field is not an email address")
	ErrNoContacts              = errors.New("at least one contact is required")
	ErrMissingContact          = errors.New("administrative contact or technical and legal contacts required")
	ErrMissingAdminOrLegal     = errors.New("administrative or legal contact required")
	ErrMissingAdminOrTechnical = errors.New("administrative or technical contact required")
	ErrInvalidCollaborator     = errors.New("collaborator record is invalid")
	ErrCollaboratorExists      = errors.New("collaborator already exists in organization")
	ErrMaxCollaborators        = errors.New("maximum number of collaborators reached")
)
