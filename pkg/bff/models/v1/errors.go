package models

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrYesNo                   = errors.New("field must be either 'yes' or 'no'")
	ErrYesNoPartially          = errors.New("field must be either 'yes', 'no', or 'partially'")
	ErrNegativeValue           = errors.New("field cannot be negative")
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
	ErrNoLegalNameIdentifier   = errors.New("at least one legal name identifier is required")
	ErrNoGeographicAddress     = errors.New("at least one geographic address is required")
	ErrTooManyAddressLines     = errors.New("an address can have at most 7 address lines")
	ErrInvalidAddress          = errors.New("address must have either address lines or street name and building number/name")
	ErrNoAddressLines          = errors.New("geographic address requires at least one address line")
	ErrInvalidCountry          = errors.New("country must be an ISO-3166-1 alpha-2 code")
	ErrInvalidCustomerNumber   = errors.New("customer number is optional but can be at most 50 characters")
	ErrInvalidLegalNatID       = errors.New("national identifier type must be RAID, MISC, LEIX, or TXID")
	ErrInvalidLEI              = errors.New("LEI identifier must not be longer than 35 characters")
	ErrNoCountryNatID          = errors.New("country of issue must be empty for legal persons")
	ErrLegalNatIDRequired      = errors.New("national identification is required to verify legal person")
	ErrNoRAForLEIX             = errors.New("registration authority must be empty for identifier type LEI")
	ErrRARequired              = errors.New("registration authority must be specified unless the identifier type is LEI")
	ErrInvalidEndpoint         = errors.New("endpoint must have the format host:port")
	ErrDuplicateEndpoint       = errors.New("mainnet endpoint cannot be the same as testnet endpoint")
	ErrMissingHost             = errors.New("endpoint string must have a host")
	ErrMissingPort             = errors.New("endpoint string must have a port")
	ErrInvalidPort             = errors.New("port must be a number between 1 and 65535")
	ErrInvalidCommonName       = errors.New("common name cannot contain wildcard characters (*) and must not have a scheme or port")
	ErrCommonNameMismatch      = errors.New("common name must match the endpoint host")
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
