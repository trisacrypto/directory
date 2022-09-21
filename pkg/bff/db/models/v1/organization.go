package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	models "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
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

// Add a new collaborator to an organization record. The given collaborator record is
// validated before being added to the organization.
// Note: The caller is responsible for saving the updated organization record to the
// database.
func (org *Organization) AddCollaborator(collab *Collaborator) (err error) {
	// TODO: More comprehensive validation of the collaborator record
	if collab.Email == "" {
		return errors.New("email address is required to add an organization collaborator")
	}

	if org.Collaborators == nil {
		org.Collaborators = make(map[string]*Collaborator)
	}

	// Don't overwrite an existing collaborator
	if _, ok := org.Collaborators[collab.Email]; ok {
		return fmt.Errorf("collaborator with email address %s already exists", collab.Email)
	}

	// Make sure the record has a created timestamp
	if collab.CreatedAt == "" {
		collab.CreatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	}

	// Add the collaborator to the organization
	org.Collaborators[collab.Email] = collab
	return nil
}

// Return the collaborator record for the given email address.
func (org *Organization) GetCollaborator(email string) (collaborator *Collaborator, err error) {
	if email == "" {
		return nil, errors.New("email address is required to get an organization collaborator")
	}

	// Lookup the collaborator record
	var ok bool
	if collaborator, ok = org.Collaborators[email]; !ok {
		return nil, fmt.Errorf("collaborator with email address %s does not exist", email)
	}

	return collaborator, nil
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
		Entity: &ivms101.LegalPerson{
			Name: &ivms101.LegalPersonName{
				NameIdentifiers: []*ivms101.LegalPersonNameId{
					{
						LegalPersonNameIdentifierType: ivms101.LegalPersonNameTypeCode_LEGAL_PERSON_NAME_TYPE_CODE_LEGL,
					},
				},
			},
			GeographicAddresses: []*ivms101.Address{
				{
					AddressType: ivms101.AddressTypeCode_ADDRESS_TYPE_CODE_BIZZ,
				},
			},
			NationalIdentification: &ivms101.NationalIdentification{
				NationalIdentifierType: ivms101.NationalIdentifierTypeCode_NATIONAL_IDENTIFIER_TYPE_CODE_LEIX,
				RegistrationAuthority:  "RA777777",
			},
		},
		Contacts: &models.Contacts{
			Technical:      &models.Contact{},
			Administrative: &models.Contact{},
			Legal:          &models.Contact{},
			Billing:        &models.Contact{},
		},
		Trixo: &models.TRIXOQuestionnaire{
			FinancialTransfersPermitted:  "no",
			HasRequiredRegulatoryProgram: "no",
			KycThreshold:                 10,
			KycThresholdCurrency:         "USD",
			ApplicableRegulations: []string{
				"FATF Recommendation 16",
			},
			ComplianceThreshold:         3000,
			ComplianceThresholdCurrency: "USD",
		},
		Testnet: &NetworkDetails{},
		Mainnet: &NetworkDetails{},
		State: &FormState{
			Current: 1,
			Steps: []*FormStep{
				{
					Key:    1,
					Status: "progress",
				},
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
	case config.TestNet:
		if r.Testnet == nil {
			return false
		}
	case config.MainNet:
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
