package models

import (
	"errors"

	"github.com/google/uuid"
	"github.com/trisacrypto/directory/pkg/bff/config"
)

var (
	ErrInvalidOrgID = errors.New("invalid organization id")
)

func (org *Organization) Key() []byte {
	if org.Id == "" {
		return uuid.Nil[:]
	}

	key := uuid.MustParse(org.Id)
	return key[:]
}

func ParseOrgID(orgID interface{}) (uuid.UUID, error) {
	switch t := orgID.(type) {
	case string:
		return uuid.Parse(t)
	case uuid.UUID:
		return t, nil
	case []byte:
		return uuid.FromBytes(t)
	default:
		return uuid.Nil, ErrInvalidOrgID
	}
}

// NewRegisterForm returns a new registration form with default values.
func NewRegisterForm() *RegistrationForm {
	// Make sure default values are populated for the frontend
	return &RegistrationForm{
		State: NewFormState(),
	}
}

// NewFormState returns a new form state with default values.
func NewFormState() *FormState {
	return &FormState{
		Current: 1,
		Steps: []*FormStep{
			{
				Key:    1,
				Status: "progress",
			},
		},
	}
}

// ReadyToSubmit performs very lightweight validation, ensuring that there are non-nil
// values on the nested data structures so that the request to the GDS does not fail.
// For data validation (required fields, types, etc.), we should rely on the GDS
// response to ensure that we're able to submit valid forms and that validation only
// occurs in one place in the code.
func (r *RegistrationForm) ReadyToSubmit(network string) bool {
	if r.VaspCategories == nil || r.Entity == nil || r.Contacts == nil || r.Trixo == nil {
		return false
	}

	switch network {
	case config.TestNetKey:
		if r.Testnet == nil {
			return false
		}
	case config.MainNetKey:
		if r.Mainnet == nil {
			return false
		}
	default:
		// If the network is not specified or a string like "all" or "both" is passed
		// in then the default behavior is to validate that both networks are ready.
		if r.Testnet == nil || r.Mainnet == nil {
			return false
		}
	}

	return true
}
