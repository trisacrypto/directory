package models

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"strconv"
	"strings"

	"net"

	"github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/utils"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	// Basic Details Fields
	FieldWebsite          = "website"
	FieldBusinessCategory = "business_category"
	FieldVASPCategories   = "vasp_categories"
	FieldEstablishedOn    = "established_on"
	FieldOrganizationName = "organization_name"

	// Legal Person Entity Fields
	FieldEntity                              = "entity"
	FieldEntityName                          = "entity.name"
	FieldEntityNameIdentifiers               = "entity.name.name_identifiers"
	FieldEntityLocalNameIdentifiers          = "entity.name.local_name_identifiers"
	FieldEntityPhoneticNameIdentifiers       = "entity.name.phonetic_name_identifiers"
	FieldEntityGeographicAddresses           = "entity.geographic_addresses"
	FieldEntityGeographicAddressLines        = "entity.geographic_addresses.address_line"
	FieldEntityGeographicAddressCountry      = "entity.geographic_addresses.country"
	FieldEntityCustomerNumber                = "entity.customer_number"
	FieldEntityNationalIdentification        = "entity.national_identification"
	FieldEntityNationalIdentificationID      = "entity.national_identification.national_identifier"
	FieldEntityNationalIdentificationType    = "entity.national_identification.national_identifier_type"
	FieldEntityNationalIdentificationCountry = "entity.national_identification.country_of_issue"
	FieldEntityNationalIdentificationRA      = "entity.national_identification.registration_authority"
	FieldEntityCountryOfRegistration         = "entity.country_of_registration"

	// Contacts Fields
	FieldContacts                    = "contacts"
	FieldContactsTechnical           = "contacts.technical"
	FieldContactsTechnicalName       = "contacts.technical.name"
	FieldContactsTechnicalEmail      = "contacts.technical.email"
	FieldContactsTechnicalPhone      = "contacts.technical.phone"
	FieldContactsAdministrative      = "contacts.administrative"
	FieldContactsAdministrativeName  = "contacts.administrative.name"
	FieldContactsAdministrativeEmail = "contacts.administrative.email"
	FieldContactsAdministrativePhone = "contacts.administrative.phone"
	FieldContactsLegal               = "contacts.legal"
	FieldContactsLegalName           = "contacts.legal.name"
	FieldContactsLegalEmail          = "contacts.legal.email"
	FieldContactsLegalPhone          = "contacts.legal.phone"
	FieldContactsBilling             = "contacts.billing"
	FieldContactsBillingName         = "contacts.billing.name"
	FieldContactsBillingEmail        = "contacts.billing.email"
	FieldContactsBillingPhone        = "contacts.billing.phone"

	// TRIXO fields
	FieldTRIXO                                = "trixo"
	FieldTRIXOPrimaryNationalJurisdiction     = "trixo.primary_national_jurisdiction"
	FieldTRIXOPrimaryRegulator                = "trixo.primary_regulator"
	FieldTRIXOFinancialTransfersPermitted     = "trixo.financial_transfers_permitted"
	FieldTRIXOOtherJurisdictions              = "trixo.other_jurisdictions"
	FieldTRIXOOtherJurisdictionsCountry       = "trixo.other_jurisdictions.country"
	FieldTRIXOOtherJurisdictionsRegulatorName = "trixo.other_jurisdictions.regulator_name"
	FieldTRIXOOtherJurisdictionsLicenseNumber = "trixo.other_jurisdictions.license_number"
	FieldTRIXOHasRequiredRegulatoryProgram    = "trixo.has_required_regulatory_program"
	FieldTRIXOConductsCustomerKYC             = "trixo.conducts_customer_kyc"
	FieldTRIXOKYCThreshold                    = "trixo.kyc_threshold"
	FieldTRIXOKYCThresholdCurrency            = "trixo.kyc_threshold_currency"
	FieldTRIXOMustComplyTravelRule            = "trixo.must_comply_travel_rule"
	FieldTRIXOApplicableRegulations           = "trixo.applicable_regulations"
	FieldTRIXOComplianceThreshold             = "trixo.compliance_threshold"
	FieldTRIXOComplianceThresholdCurrency     = "trixo.compliance_threshold_currency"
	FieldTRIXOMustSafeguardPII                = "trixo.must_safeguard_pii"
	FieldTRIXOSafeGuardsPII                   = "trixo.safeguards_pii"

	// TRISA Details Fields
	FieldTestNet           = "testnet"
	FieldTestNetCommonName = "testnet.common_name"
	FieldTestNetEndpoint   = "testnet.endpoint"
	FieldTestNetDNSNames   = "testnet.dns_names"
	FieldMainNet           = "mainnet"
	FieldMainNetCommonName = "mainnet.common_name"
	FieldMainNetEndpoint   = "mainnet.endpoint"
	FieldMainNetDNSNames   = "mainnet.dns_names"

	// Default fields
	FieldState = "state"
)

// StepType represents a collection of fields in the registration form that are handled
// together as a single step when the user is filling in the registration form.
type StepType string

const (
	StepNone         StepType = ""
	StepAll          StepType = "all"
	StepBasicDetails StepType = "basic"
	StepLegalPerson  StepType = "legal"
	StepContacts     StepType = "contacts"
	StepTRISA        StepType = "trisa"
	StepTRIXO        StepType = "trixo"
)

// Parse a string as a step type.
func ParseStepType(s string) (StepType, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case string(StepNone):
		return StepNone, nil
	case string(StepAll):
		return StepAll, nil
	case string(StepBasicDetails):
		return StepBasicDetails, nil
	case string(StepLegalPerson):
		return StepLegalPerson, nil
	case string(StepContacts):
		return StepContacts, nil
	case string(StepTRIXO):
		return StepTRIXO, nil
	case string(StepTRISA):
		return StepTRISA, nil
	default:
		return StepNone, fmt.Errorf("unknown registration form step %q", s)
	}
}

func (s StepType) String() string {
	return string(s)
}

// Validate the registration form returning all field errors as opposed to a single
// error that shortcircuits when the first validation error is found. If a step is
// specified then only that step's fields are validated.
func (r *RegistrationForm) Validate(step StepType) error {
	switch step {
	case StepBasicDetails:
		return r.ValidateBasicDetails()
	case StepLegalPerson:
		return r.ValidateLegalPerson()
	case StepContacts:
		return r.ValidateContacts()
	case StepTRIXO:
		return r.ValidateTRIXO()
	case StepTRISA:
		return r.ValidateTRISA()
	case StepNone, StepAll:
		// If empty string, or "all" validate the entire form.
		errs := make(ValidationErrors, 0)
		errs, _ = errs.Append(r.ValidateBasicDetails())
		errs, _ = errs.Append(r.ValidateLegalPerson())
		errs, _ = errs.Append(r.ValidateContacts())
		errs, _ = errs.Append(r.ValidateTRIXO())
		errs, _ = errs.Append(r.ValidateTRISA())

		if len(errs) == 0 {
			return nil
		}
		return errs
	default:
		return fmt.Errorf("unknown registration form step %q", step)
	}
}

// Validate only the fields in the basic details step.
func (r *RegistrationForm) ValidateBasicDetails() error {
	err := make(ValidationErrors, 0)

	// Validate website
	// TODO: Check if this is a valid URL?
	if strings.TrimSpace(r.Website) == "" {
		err = append(err, &ValidationError{
			Field: FieldWebsite,
			Err:   ErrMissingField.Error(),
		})
	}

	// Validate business category
	if r.BusinessCategory == pb.BusinessCategory_UNKNOWN_ENTITY {
		err = append(err, &ValidationError{
			Field: FieldBusinessCategory,
			Err:   ErrMissingField.Error(),
		})
	}

	// Validate VASP categories if provided
	// TODO: Check if these are valid categories?
	for i, cat := range r.VaspCategories {
		r.VaspCategories[i] = strings.TrimSpace(cat)
		if r.VaspCategories[i] == "" {
			err = append(err, &ValidationError{
				Field: FieldVASPCategories,
				Err:   ErrMissingField.Error(),
				Index: i,
			})
		}
	}

	// Validate established date
	// TODO: Check if this is a valid date?
	if strings.TrimSpace(r.EstablishedOn) == "" {
		err = append(err, &ValidationError{
			Field: FieldEstablishedOn,
			Err:   ErrMissingField.Error(),
		})
	}

	// Validate organization name
	if strings.TrimSpace(r.OrganizationName) == "" {
		err = append(err, &ValidationError{
			Field: FieldOrganizationName,
			Err:   ErrMissingField.Error(),
		})
	}

	if len(err) == 0 {
		return nil
	}
	return err
}

// Validate only the fields in the legal person step.
func (r *RegistrationForm) ValidateLegalPerson() error {
	err := make(ValidationErrors, 0)
	if r.Entity == nil {
		return append(err, &ValidationError{
			Field: FieldEntity,
			Err:   ErrMissingField.Error(),
		})
	}

	// Validate name identifiers
	if r.Entity.Name == nil {
		err = append(err, &ValidationError{
			Field: FieldEntityName,
			Err:   ErrMissingField.Error(),
		})
	} else {
		// Ensure there is at least one legal name identifier
		var legalNames uint32
		for i, name := range r.Entity.Name.NameIdentifiers {
			if name.LegalPersonNameIdentifierType == ivms101.LegalPersonLegal {
				legalNames++
			}

			if verr := ValidateLegalPersonName(name); verr != nil {
				err = append(err, &ValidationError{
					Field: FieldEntityNameIdentifiers,
					Err:   verr.Error(),
					Index: i,
				})
			}
		}

		if legalNames == 0 {
			err = append(err, &ValidationError{
				Field: FieldEntityNameIdentifiers,
				Err:   ErrNoLegalNameIdentifier.Error(),
			})
		}

		// Validate local name identifiers
		for i, name := range r.Entity.Name.LocalNameIdentifiers {
			if verr := ValidateLegalPersonLocalName(name); verr != nil {
				err = append(err, &ValidationError{
					Field: FieldEntityLocalNameIdentifiers,
					Err:   verr.Error(),
					Index: i,
				})
			}
		}

		// Validate phonetic name identifiers
		for i, name := range r.Entity.Name.PhoneticNameIdentifiers {
			if verr := ValidateLegalPersonLocalName(name); verr != nil {
				err = append(err, &ValidationError{
					Field: FieldEntityPhoneticNameIdentifiers,
					Err:   verr.Error(),
					Index: i,
				})
			}
		}
	}

	// Validate Geographic Addresses
	if len(r.Entity.GeographicAddresses) == 0 {
		err = append(err, &ValidationError{
			Field: FieldEntityGeographicAddresses,
			Err:   ErrNoGeographicAddress.Error(),
		})
	} else {
		for i, addr := range r.Entity.GeographicAddresses {
			// TODO: do we need to validate address type code?

			// There can be at most 7 address lines
			if len(addr.AddressLine) > 7 {
				err = append(err, &ValidationError{
					Field: FieldEntityGeographicAddressLines,
					Err:   ErrTooManyAddressLines.Error(),
					Index: i,
				})
			}

			// Valid address is either address lines or street name + building number.
			if len(addr.AddressLine) == 0 && (addr.StreetName == "" && (addr.BuildingName == "" || addr.BuildingNumber == "")) {
				err = append(err, &ValidationError{
					Field: FieldEntityGeographicAddresses,
					Err:   ErrInvalidAddress.Error(),
					Index: i,
				})
			}

			// Address lines cannot all be blank
			var validAddrLines uint16
			for i, line := range addr.AddressLine {
				addr.AddressLine[i] = strings.TrimSpace(line)
				if strings.TrimSpace(line) != "" {
					validAddrLines++
				}
			}

			if validAddrLines == 0 {
				err = append(err, &ValidationError{
					Field: FieldEntityGeographicAddressLines,
					Err:   ErrNoAddressLines.Error(),
					Index: i,
				})
			}

			// Country must be an alpha-2 country code
			addr.Country = strings.TrimSpace(addr.Country)
			if addr.Country == "" {
				err = append(err, &ValidationError{
					Field: FieldEntityGeographicAddressCountry,
					Err:   ErrMissingField.Error(),
					Index: i,
				})
			} else if len(addr.Country) != 2 {
				err = append(err, &ValidationError{
					Field: FieldEntityGeographicAddressCountry,
					Err:   ErrInvalidCountry.Error(),
					Index: i,
				})
			}
		}
	}

	// Customer number must not be greater than 50 chars
	r.Entity.CustomerNumber = strings.TrimSpace(r.Entity.CustomerNumber)
	if r.Entity.CustomerNumber != "" && len(r.Entity.CustomerNumber) > 50 {
		err = append(err, &ValidationError{
			Field: FieldEntityCustomerNumber,
			Err:   ErrInvalidCustomerNumber.Error(),
		})
	}

	// Validate National Identification
	if r.Entity.NationalIdentification != nil {
		// Validate National Identification
		r.Entity.NationalIdentification.NationalIdentifier = strings.TrimSpace(r.Entity.NationalIdentification.NationalIdentifier)
		if r.Entity.NationalIdentification.NationalIdentifier == "" {
			err = append(err, &ValidationError{
				Field: FieldEntityNationalIdentificationID,
				Err:   ErrMissingField.Error(),
			})
		}

		// Validate National Identification Type Code
		if !(r.Entity.NationalIdentification.NationalIdentifierType == ivms101.NationalIdentifierRAID ||
			r.Entity.NationalIdentification.NationalIdentifierType == ivms101.NationalIdentifierMISC ||
			r.Entity.NationalIdentification.NationalIdentifierType == ivms101.NationalIdentifierLEIX ||
			r.Entity.NationalIdentification.NationalIdentifierType == ivms101.NationalIdentifierTXID) {
			err = append(err, &ValidationError{
				Field: FieldEntityNationalIdentificationType,
				Err:   ErrInvalidLegalNatID.Error(),
			})
		}

		// TODO: validate LEI with checksum
		if r.Entity.NationalIdentification.NationalIdentifierType == ivms101.NationalIdentifierLEIX {
			if len(r.Entity.NationalIdentification.NationalIdentifier) > 35 {
				err = append(err, &ValidationError{
					Field: FieldEntityNationalIdentificationID,
					Err:   ErrInvalidLEI.Error(),
				})
			}
		}

		// Country of issue is only used for natural persons
		r.Entity.NationalIdentification.CountryOfIssue = strings.TrimSpace(r.Entity.NationalIdentification.CountryOfIssue)
		if r.Entity.NationalIdentification.CountryOfIssue != "" {
			err = append(err, &ValidationError{
				Field: FieldEntityNationalIdentificationCountry,
				Err:   ErrNoCountryNatID.Error(),
			})
		}

		// If the ID is an LEIX then registration authority must be empty and vice-versa.
		r.Entity.NationalIdentification.RegistrationAuthority = strings.TrimSpace(r.Entity.NationalIdentification.RegistrationAuthority)
		if r.Entity.NationalIdentification.NationalIdentifierType != ivms101.NationalIdentifierLEIX {
			if r.Entity.NationalIdentification.RegistrationAuthority == "" {
				err = append(err, &ValidationError{
					Field: FieldEntityNationalIdentificationRA,
					Err:   ErrRARequired.Error(),
				})
			}
		} else {
			// If the ID is an LEIX, Registration Authority must be empty
			if r.Entity.NationalIdentification.RegistrationAuthority != "" {
				err = append(err, &ValidationError{
					Field: FieldEntityNationalIdentificationRA,
					Err:   ErrNoRAForLEIX.Error(),
				})
			}
		}
	} else {
		err = append(err, &ValidationError{
			Field: FieldEntityNationalIdentification,
			Err:   ErrLegalNatIDRequired.Error(),
		})
	}

	// Country Code Constratint
	r.Entity.CountryOfRegistration = strings.TrimSpace(r.Entity.CountryOfRegistration)
	if r.Entity.CountryOfRegistration != "" {
		// TODO: ensure the country code is valid?
		if len(r.Entity.CountryOfRegistration) != 2 {
			err = append(err, &ValidationError{
				Field: FieldEntityCountryOfRegistration,
				Err:   ErrInvalidCountry.Error(),
			})
		}
	} else {
		err = append(err, &ValidationError{
			Field: FieldEntityCountryOfRegistration,
			Err:   ErrMissingField.Error(),
		})
	}

	// Final validation just to check and make sure we didn't miss anything
	if verr := r.Entity.Validate(); verr != nil {
		// C9 constraint errors are ignored because they are also ignored in GDS
		if verr.Error() != ivms101.ErrCompleteNationalIdentifierCountry.Error() {
			err = append(err, &ValidationError{
				Field: FieldEntity,
				Err:   verr.Error(),
			})
		}
	}

	if len(err) == 0 {
		return nil
	}
	return err
}

func ValidateLegalPersonName(name *ivms101.LegalPersonNameId) error {
	// Validate the name identifier
	name.LegalPersonName = strings.TrimSpace(name.LegalPersonName)
	if name.LegalPersonName == "" {
		return ErrMissingField
	}

	if len(name.LegalPersonName) > 100 {
		return ErrLegalPersonNameLength
	}

	// TODO: does the legal person name type code need to be validated?
	return nil
}

func ValidateLegalPersonLocalName(name *ivms101.LocalLegalPersonNameId) error {
	// Validate the name identifier
	name.LegalPersonName = strings.TrimSpace(name.LegalPersonName)
	if name.LegalPersonName == "" {
		return ErrMissingField
	}

	if len(name.LegalPersonName) > 100 {
		return ErrLegalPersonNameLength
	}

	// TODO: does the legal person name type code need to be validated?
	return nil
}

// Validate only the fields in the contacts step.
func (r *RegistrationForm) ValidateContacts() error {
	err := make(ValidationErrors, 0)
	if r.Contacts == nil {
		err = append(err, &ValidationError{Field: FieldContacts, Err: ErrMissingField.Error()})
		return err
	}

	// Create a contacts data structure for validation that does not have associated
	// email records. The email records are created after validation during the GDS
	// registration process and not on the BFF.
	contacts := &models.Contacts{Contacts: r.Contacts}

	// Validate each non-zero contact

	iter := contacts.NewIterator()
	for iter.Next() {
		contact := iter.Contact()
		err, _ = err.Append(ValidateContact(contact.Contact, FieldContacts+"."+contact.Kind))
	}

	// Check that all required contacts are present (special rules)
	switch contacts.Length() {
	case 0:
		// At least one contact is required
		err = append(err, &ValidationError{Field: FieldContacts, Err: ErrNoContacts.Error()})
	case 1:
		// If there is only one contact, it must be the admin; if not highlight the missing fields
		if !contacts.Has(models.AdministrativeContact) {
			// Global contact error
			err = append(err, &ValidationError{Field: FieldContacts, Err: ErrMissingContact.Error()})
			switch {
			case contacts.Has(models.TechnicalContact):
				// If the technical contact is filled in then nominate the legal/admin contact to be populated
				err = append(err, &ValidationError{Field: FieldContactsAdministrative, Err: ErrMissingAdminOrLegal.Error()})
				err = append(err, &ValidationError{Field: FieldContactsLegal, Err: ErrMissingAdminOrLegal.Error()})
			case contacts.Has(models.LegalContact):
				// If the legal contact is filled in then nominate the technical/admin contact to be populated
				err = append(err, &ValidationError{Field: FieldContactsAdministrative, Err: ErrMissingAdminOrTechnical.Error()})
				err = append(err, &ValidationError{Field: FieldContactsTechnical, Err: ErrMissingAdminOrTechnical.Error()})
			default:
				// Otherwise say that one of the other fields is required
				err = append(err, &ValidationError{Field: FieldContactsAdministrative, Err: ErrMissingContact.Error()})
				err = append(err, &ValidationError{Field: FieldContactsTechnical, Err: ErrMissingContact.Error()})
				err = append(err, &ValidationError{Field: FieldContactsLegal, Err: ErrMissingContact.Error()})
			}
		}
	default:
		// If there are at least two contacts, either admin or technical must be present
		if !contacts.Has(models.AdministrativeContact) && !contacts.Has(models.TechnicalContact) {
			err = append(err, &ValidationError{Field: FieldContacts, Err: ErrMissingContact.Error()})
			err = append(err, &ValidationError{Field: FieldContactsAdministrative, Err: ErrMissingAdminOrTechnical.Error()})
			err = append(err, &ValidationError{Field: FieldContactsTechnical, Err: ErrMissingAdminOrTechnical.Error()})
		}
		// Admin or legal must be present
		if !contacts.Has(models.AdministrativeContact) && !contacts.Has(models.LegalContact) {
			err = append(err, &ValidationError{Field: FieldContacts, Err: ErrMissingContact.Error()})
			err = append(err, &ValidationError{Field: FieldContactsAdministrative, Err: ErrMissingAdminOrLegal.Error()})
			err = append(err, &ValidationError{Field: FieldContactsLegal, Err: ErrMissingAdminOrLegal.Error()})
		}
	}

	if len(err) == 0 {
		return nil
	}
	return err
}

// Validate a single contact, using the field name to construct errors.
func ValidateContact(contact *pb.Contact, fieldName string) error {
	err := make(ValidationErrors, 0)

	name := strings.TrimSpace(contact.Name)
	if name == "" {
		err = append(err, &ValidationError{Field: fieldName + ".name", Err: ErrMissingField.Error()})
	} else if len(name) < 2 {
		// Check that the name is long enough to match GDS validation
		err = append(err, &ValidationError{Field: fieldName + ".name", Err: ErrTooShort.Error()})
	}

	email := strings.TrimSpace(contact.Email)
	if email == "" {
		err = append(err, &ValidationError{Field: fieldName + ".email", Err: ErrMissingField.Error()})
	} else {
		// Check that the email is parseable by RFC 5322.
		if _, verr := mail.ParseAddress(email); verr != nil {
			err = append(err, &ValidationError{Field: fieldName + ".email", Err: ErrInvalidEmail.Error()})
		}
	}

	// Phone number is optional for technical contacts
	phone := strings.TrimSpace(contact.Phone)
	if phone == "" && fieldName != FieldContactsTechnical {
		err = append(err, &ValidationError{Field: fieldName + ".phone", Err: ErrMissingField.Error()})
	}

	if len(err) == 0 {
		return nil
	}
	return err
}

// Validate only the fields in the trixo step.
func (r *RegistrationForm) ValidateTRIXO() error {
	err := make(ValidationErrors, 0)

	if r.Trixo == nil {
		err = append(err, &ValidationError{
			Field: FieldTRIXO,
			Err:   ErrMissingField.Error(),
		})
		return err
	}

	r.Trixo.PrimaryNationalJurisdiction = strings.TrimSpace(r.Trixo.PrimaryNationalJurisdiction)
	if r.Trixo.PrimaryNationalJurisdiction == "" {
		err = append(err, &ValidationError{
			Field: FieldTRIXOPrimaryNationalJurisdiction,
			Err:   ErrMissingField.Error(),
		})
	} else if len(r.Trixo.PrimaryNationalJurisdiction) != 2 {
		err = append(err, &ValidationError{
			Field: FieldTRIXOPrimaryNationalJurisdiction,
			Err:   ErrInvalidCountry.Error(),
		})
	}

	r.Trixo.PrimaryRegulator = strings.TrimSpace(r.Trixo.PrimaryRegulator)
	if r.Trixo.PrimaryRegulator == "" {
		err = append(err, &ValidationError{
			Field: FieldTRIXOPrimaryRegulator,
			Err:   ErrMissingField.Error(),
		})
	}

	r.Trixo.FinancialTransfersPermitted = strings.ToLower(strings.TrimSpace(r.Trixo.FinancialTransfersPermitted))
	if r.Trixo.FinancialTransfersPermitted == "" {
		err = append(err, &ValidationError{
			Field: FieldTRIXOFinancialTransfersPermitted,
			Err:   ErrMissingField.Error(),
		})
	} else if r.Trixo.FinancialTransfersPermitted != "yes" && r.Trixo.FinancialTransfersPermitted != "no" && r.Trixo.FinancialTransfersPermitted != "partially" {
		err = append(err, &ValidationError{
			Field: FieldTRIXOFinancialTransfersPermitted,
			Err:   ErrYesNoPartially.Error(),
		})
	}

	for i, juris := range r.Trixo.OtherJurisdictions {
		r.Trixo.OtherJurisdictions[i].Country = strings.TrimSpace(juris.Country)
		if r.Trixo.OtherJurisdictions[i].Country == "" {
			err = append(err, &ValidationError{
				Field: FieldTRIXOOtherJurisdictionsCountry,
				Err:   ErrMissingField.Error(),
				Index: i,
			})
		} else if len(r.Trixo.OtherJurisdictions[i].Country) != 2 {
			err = append(err, &ValidationError{
				Field: FieldTRIXOOtherJurisdictionsCountry,
				Err:   ErrInvalidCountry.Error(),
				Index: i,
			})
		}

		r.Trixo.OtherJurisdictions[i].RegulatorName = strings.TrimSpace(juris.RegulatorName)
		if r.Trixo.OtherJurisdictions[i].RegulatorName == "" {
			err = append(err, &ValidationError{
				Field: FieldTRIXOOtherJurisdictionsRegulatorName,
				Err:   ErrMissingField.Error(),
				Index: i,
			})
		}
	}

	r.Trixo.HasRequiredRegulatoryProgram = strings.ToLower(strings.TrimSpace(r.Trixo.HasRequiredRegulatoryProgram))
	if r.Trixo.HasRequiredRegulatoryProgram == "" {
		err = append(err, &ValidationError{
			Field: FieldTRIXOHasRequiredRegulatoryProgram,
			Err:   ErrMissingField.Error(),
		})
	} else if r.Trixo.HasRequiredRegulatoryProgram != "yes" && r.Trixo.HasRequiredRegulatoryProgram != "no" {
		err = append(err, &ValidationError{
			Field: FieldTRIXOHasRequiredRegulatoryProgram,
			Err:   ErrYesNo.Error(),
		})
	}

	if r.Trixo.ConductsCustomerKyc {
		if r.Trixo.KycThreshold < 0 {
			err = append(err, &ValidationError{
				Field: FieldTRIXOKYCThreshold,
				Err:   ErrNegativeValue.Error(),
			})
		}

		// TODO: Validate currency code
		r.Trixo.KycThresholdCurrency = strings.TrimSpace(r.Trixo.KycThresholdCurrency)
		if r.Trixo.KycThresholdCurrency == "" {
			err = append(err, &ValidationError{
				Field: FieldTRIXOKYCThresholdCurrency,
				Err:   ErrMissingField.Error(),
			})
		}
	}

	if r.Trixo.MustComplyTravelRule {
		for i, reg := range r.Trixo.ApplicableRegulations {
			r.Trixo.ApplicableRegulations[i] = strings.TrimSpace(reg)
			if r.Trixo.ApplicableRegulations[i] == "" {
				err = append(err, &ValidationError{
					Field: FieldTRIXOApplicableRegulations,
					Err:   ErrMissingField.Error(),
					Index: i,
				})
			}
		}

		if r.Trixo.ComplianceThreshold < 0 {
			err = append(err, &ValidationError{
				Field: FieldTRIXOComplianceThreshold,
				Err:   ErrNegativeValue.Error(),
			})
		}

		r.Trixo.ComplianceThresholdCurrency = strings.TrimSpace(r.Trixo.ComplianceThresholdCurrency)
		if r.Trixo.ComplianceThresholdCurrency == "" {
			err = append(err, &ValidationError{
				Field: FieldTRIXOComplianceThresholdCurrency,
				Err:   ErrMissingField.Error(),
			})
		}
	}

	if len(err) == 0 {
		return nil
	}

	return err
}

// Validate only the fields in the trisa step.
func (r *RegistrationForm) ValidateTRISA() error {
	err := make(ValidationErrors, 0)

	// At least one network must be specified
	if NetworkDetailsIsZero(r.Testnet) && NetworkDetailsIsZero(r.Mainnet) {
		err = append(err, &ValidationError{Field: FieldTestNet, Err: ErrMissingTestNetOrMainNet.Error()})
		err = append(err, &ValidationError{Field: FieldMainNet, Err: ErrMissingTestNetOrMainNet.Error()})
		return err
	}

	// Only validate testnet if it has data on it
	if !NetworkDetailsIsZero(r.Testnet) {
		err, _ = err.Append(validateNetwork(r.Testnet, FieldTestNet))
	}

	// Only validate mainnet if it has data on it
	if !NetworkDetailsIsZero(r.Mainnet) {
		err, _ = err.Append(validateNetwork(r.Mainnet, FieldMainNet))
		if !NetworkDetailsIsZero(r.Testnet) && r.Mainnet.Endpoint != "" && r.Mainnet.Endpoint == r.Testnet.Endpoint {
			err = append(err, &ValidationError{Field: FieldMainNetEndpoint, Err: ErrDuplicateEndpoint.Error()})
		}
	}

	if len(err) == 0 {
		return nil
	}

	return err
}

// NetworkDetailsIsZero returns true if the network details are nil or zero.
func NetworkDetailsIsZero(details *NetworkDetails) bool {
	if details == nil {
		return true
	}

	return details.Endpoint == "" && details.CommonName == "" && details.DnsNames == nil
}

// Validate a network details field.
func validateNetwork(details *NetworkDetails, fieldName string) error {
	verr := make(ValidationErrors, 0)

	// Validate the endpoint
	var host string
	details.Endpoint = strings.TrimSpace(details.Endpoint)
	if details.Endpoint == "" {
		verr = append(verr, &ValidationError{Field: fieldName + ".endpoint", Err: ErrMissingField.Error()})
	} else {
		var (
			port string
			err  error
		)
		if host, port, err = net.SplitHostPort(details.Endpoint); err != nil {
			verr = append(verr, &ValidationError{Field: fieldName + ".endpoint", Err: ErrInvalidEndpoint.Error()})
		} else {
			if host == "" {
				verr = append(verr, &ValidationError{Field: fieldName + ".endpoint", Err: ErrMissingHost.Error()})
			}

			if port == "" {
				verr = append(verr, &ValidationError{Field: fieldName + ".endpoint", Err: ErrMissingPort.Error()})
			} else if _, err := strconv.Atoi(port); err != nil {
				verr = append(verr, &ValidationError{Field: fieldName + ".endpoint", Err: ErrInvalidPort.Error()})
			}
		}
	}

	// Validate the common name
	details.CommonName = strings.TrimSpace(details.CommonName)
	if details.CommonName == "" {
		verr = append(verr, &ValidationError{Field: fieldName + ".common_name", Err: ErrMissingField.Error()})
	} else {
		if utils.ValidateCommonName(details.CommonName) != nil {
			verr = append(verr, &ValidationError{Field: fieldName + ".common_name", Err: ErrInvalidCommonName.Error()})
		}

		// Common name must match the endpoint host if host is not an IP address
		if net.ParseIP(host) == nil && host != details.CommonName {
			verr = append(verr, &ValidationError{Field: fieldName + ".common_name", Err: ErrCommonNameMismatch.Error()})
		}
	}

	// Validate the DNS names
	for i, dns := range details.DnsNames {
		details.DnsNames[i] = strings.TrimSpace(dns)
		if details.DnsNames[i] == "" {
			verr = append(verr, &ValidationError{Field: fieldName + ".dns_names", Err: ErrMissingField.Error(), Index: i})
		} else {
			if utils.ValidateCommonName(details.DnsNames[i]) != nil {
				verr = append(verr, &ValidationError{Field: fieldName + ".dns_names", Err: ErrInvalidCommonName.Error(), Index: i})
			}
		}
	}

	if len(verr) == 0 {
		return nil
	}
	return verr
}

// Update the registration form from another registration form model. If a step is
// specified then only that step from the other registration form is copied to this form
// otherwise the entire registration form is updated. If there is an update error it is
// returned, otherwise validation is performed and ValidationErrors are returned.
func (r *RegistrationForm) Update(o *RegistrationForm, step StepType) error {
	// No matter what, update the form state
	if o.State != nil {
		r.State = o.State
	}

	switch step {
	case StepBasicDetails:
		return r.UpdateBasicDetails(o)
	case StepLegalPerson:
		return r.UpdateLegalPerson(o)
	case StepContacts:
		return r.UpdateContacts(o)
	case StepTRIXO:
		return r.UpdateTRIXO(o)
	case StepTRISA:
		return r.UpdateTRISA(o)
	case StepNone, StepAll:
		var ok bool
		errs := make(ValidationErrors, 0)

		err := r.UpdateBasicDetails(o)
		if errs, ok = errs.Append(err); !ok {
			return err
		}

		err = r.UpdateLegalPerson(o)
		if errs, ok = errs.Append(err); !ok {
			return err
		}

		err = r.UpdateContacts(o)
		if errs, ok = errs.Append(err); !ok {
			return err
		}

		err = r.UpdateTRIXO(o)
		if errs, ok = errs.Append(err); !ok {
			return err
		}

		err = r.UpdateTRISA(o)
		if errs, ok = errs.Append(err); !ok {
			return err
		}

		if len(errs) == 0 {
			return nil
		}
		return errs
	default:
		return fmt.Errorf("unknown step %q", step)
	}
}

// Update only the fields from the basic details step.
func (r *RegistrationForm) UpdateBasicDetails(o *RegistrationForm) error {
	// TODO: make this functionality "PATCH" right now it is "PUT"
	r.Website = o.Website
	r.BusinessCategory = o.BusinessCategory
	r.VaspCategories = o.VaspCategories
	r.EstablishedOn = o.EstablishedOn
	r.OrganizationName = o.OrganizationName
	return r.ValidateBasicDetails()
}

// Update only the fields from the legal person step.
func (r *RegistrationForm) UpdateLegalPerson(o *RegistrationForm) error {
	// TODO: make this functionality "PATCH" right now it is "PUT"
	r.Entity = o.Entity
	return r.ValidateLegalPerson()
}

// Update only the fields from the contacts step.
func (r *RegistrationForm) UpdateContacts(o *RegistrationForm) error {
	// TODO: make this functionality "PATCH" right now it is "PUT"
	r.Contacts = o.Contacts
	return r.ValidateContacts()
}

// Update only the fields from the TRIXO step.
func (r *RegistrationForm) UpdateTRIXO(o *RegistrationForm) error {
	// TODO: make this functionality "PATCH" right now it is "PUT"
	r.Trixo = o.Trixo
	return r.ValidateTRIXO()
}

// Update only the fields from the TRISA step.
func (r *RegistrationForm) UpdateTRISA(o *RegistrationForm) error {
	// TODO: make this functionality "PATCH" right now it is "PUT"
	r.Testnet = o.Testnet
	r.Mainnet = o.Mainnet
	return r.ValidateTRISA()
}

// Truncate reutrns a new registration form with only the specified step's data. If none
// or all is specified then the original registration form is returned without error.
func (r *RegistrationForm) Truncate(step StepType) (*RegistrationForm, error) {
	switch step {
	case StepBasicDetails:
		return r.TruncateBasicDetails(), nil
	case StepLegalPerson:
		return r.TruncateLegalPerson(), nil
	case StepContacts:
		return r.TruncateContacts(), nil
	case StepTRIXO:
		return r.TruncateTRIXO(), nil
	case StepTRISA:
		return r.TruncateTRISA(), nil
	case StepNone, StepAll:
		return r, nil
	default:
		return nil, fmt.Errorf("unknown registration form step %q", step)
	}
}

// Returns a registration form with only the original details.
func (r *RegistrationForm) TruncateBasicDetails() *RegistrationForm {
	return &RegistrationForm{
		Website:          r.Website,
		BusinessCategory: r.BusinessCategory,
		VaspCategories:   r.VaspCategories,
		EstablishedOn:    r.EstablishedOn,
		OrganizationName: r.OrganizationName,
		State:            r.State,
	}
}

// Returns a registration form with only the IVMS101 legal person entity (same pointer).
func (r *RegistrationForm) TruncateLegalPerson() *RegistrationForm {
	return &RegistrationForm{
		Entity: r.Entity,
		State:  r.State,
	}
}

// Returns a registration form with only the contacts (same pointer).
func (r *RegistrationForm) TruncateContacts() *RegistrationForm {
	return &RegistrationForm{
		Contacts: r.Contacts,
		State:    r.State,
	}
}

// Returns a registration form with only the TRIXO form (same pointer).
func (r *RegistrationForm) TruncateTRIXO() *RegistrationForm {
	return &RegistrationForm{
		Trixo: r.Trixo,
		State: r.State,
	}
}

// Returns a registration form with only the network details (same pointers).
func (r *RegistrationForm) TruncateTRISA() *RegistrationForm {
	return &RegistrationForm{
		Testnet: r.Testnet,
		Mainnet: r.Mainnet,
		State:   r.State,
	}
}

// ProtocolBuffer JSON marshaling and unmarshaling ensures that the BFF JSON API works
// as expected with protocol buffer models that are stored in the database.
var (
	pbencoder = protojson.MarshalOptions{
		Multiline:       false,
		Indent:          "",
		AllowPartial:    true,
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
	}
	pbdecoder = protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: false,
	}
)

// MarshalJSON uses protojson with default marshaling options.
func (r *RegistrationForm) MarshalJSON() ([]byte, error) {
	return pbencoder.Marshal(r)
}

// UnmarshalJSON uses protojson with default unmarshaling options.
func (r *RegistrationForm) UnmarshalJSON(data []byte) error {
	return pbdecoder.Unmarshal(data, r)
}

// MarshalStepJSON excludes all fields not related to the specified step when returning
// JSON data, ensuring the front-end only gets the data it needs.
func (r *RegistrationForm) MarshalStepJSON(step StepType) (_ []byte, err error) {
	var intermediate map[string]interface{}
	if intermediate, err = r.MarshalStep(step); err != nil {
		return nil, err
	}
	return json.Marshal(intermediate)
}

func (r *RegistrationForm) MarshalStep(step StepType) (intermediate map[string]interface{}, err error) {
	// Create an intermediate JSON representation from the protojson.
	var data []byte
	if data, err = r.MarshalJSON(); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &intermediate); err != nil {
		return nil, err
	}

	// Delete all keys except for the specified step and state
	toKeep := map[string]struct{}{
		FieldState: {},
	}

	switch step {
	case StepAll, StepNone:
		return intermediate, nil
	case StepBasicDetails:
		toKeep[FieldWebsite] = struct{}{}
		toKeep[FieldBusinessCategory] = struct{}{}
		toKeep[FieldVASPCategories] = struct{}{}
		toKeep[FieldEstablishedOn] = struct{}{}
		toKeep[FieldOrganizationName] = struct{}{}
	case StepLegalPerson:
		toKeep[FieldEntity] = struct{}{}
	case StepContacts:
		toKeep[FieldContacts] = struct{}{}
	case StepTRIXO:
		toKeep[FieldTRIXO] = struct{}{}
	case StepTRISA:
		toKeep[FieldTestNet] = struct{}{}
		toKeep[FieldMainNet] = struct{}{}
	default:
		return nil, fmt.Errorf("unknown registration form step %q", step)
	}

	for key := range intermediate {
		if _, ok := toKeep[key]; !ok {
			delete(intermediate, key)
		}
	}

	return intermediate, nil
}
