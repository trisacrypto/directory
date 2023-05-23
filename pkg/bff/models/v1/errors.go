package models

import (
	"errors"
	"fmt"
	"strings"
)

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
	ErrLegalPersonNameLength   = errors.New("legal person name must be less than 100 characters")
)

type ValidationError struct {
	Field string
	Err   string
	Index int
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("invalid field %s: %s", v.Field, v.Err)
}

type ValidationErrors []*ValidationError

func (v ValidationErrors) Error() string {
	errs := make([]string, 0, len(v))
	for _, e := range v {
		errs = append(errs, e.Error())
	}
	return fmt.Sprintf("%d validation errors occurred:\n%s", len(v), strings.Join(errs, "\n"))
}

// If err is a ValidationErrors then append them to this list of validation errors and
// return true, otherwise return false since we can't append random errors.
func (v ValidationErrors) Append(err error) (ValidationErrors, bool) {
	if err == nil {
		return v, true
	}

	var e *ValidationError
	if errors.As(err, &e) {
		return append(v, e), true
	}

	var es ValidationErrors
	if errors.As(err, &es) {
		return append(v, es...), true
	}
	return v, false
}
