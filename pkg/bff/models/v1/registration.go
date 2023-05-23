package models

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/trisacrypto/directory/pkg/models/v1"
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

	// Validate there is at least one VASP category provided
	// TODO: Check if these are valid categories?
	var hasCategory bool
	for _, cat := range r.VaspCategories {
		if strings.TrimSpace(cat) != "" {
			hasCategory = true
			break
		}
	}

	if !hasCategory {
		err = append(err, &ValidationError{
			Field: FieldVASPCategories,
			Err:   ErrMissingField.Error(),
		})
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

			if verr := r.validateLegalPersonName(name); verr != nil {
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
			if verr := r.validateLegalPersonLocalName(name); verr != nil {
				err = append(err, &ValidationError{
					Field: FieldEntityLocalNameIdentifiers,
					Err:   verr.Error(),
					Index: i,
				})
			}
		}

		// Validate phonetic name identifiers
		for i, name := range r.Entity.Name.PhoneticNameIdentifiers {
			if verr := r.validateLegalPersonLocalName(name); verr != nil {
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
			for _, line := range addr.AddressLine {
				if line != "" {
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
	if r.Entity.CustomerNumber != "" && len(r.Entity.CustomerNumber) > 50 {
		err = append(err, &ValidationError{
			Field: FieldEntityCustomerNumber,
			Err:   ErrInvalidCustomerNumber.Error(),
		})
	}

	// Validate National Identification
	if r.Entity.NationalIdentification != nil {
		// Validate National Identification
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
		if r.Entity.NationalIdentification.CountryOfIssue != "" {
			err = append(err, &ValidationError{
				Field: FieldEntityNationalIdentificationCountry,
				Err:   ErrNoCountryNatID.Error(),
			})
		}

		// If the ID is an LEIX then registration authority must be empty and vice-versa.
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
		err = append(err, &ValidationError{
			Field: FieldEntity,
			Err:   verr.Error(),
		})
	}

	if len(err) == 0 {
		return nil
	}
	return err
}

func (r *RegistrationForm) validateLegalPersonName(name *ivms101.LegalPersonNameId) error {
	// Validate the name identifier
	if name.LegalPersonName == "" {
		return ErrMissingField
	}

	if len(name.LegalPersonName) > 100 {
		return ErrLegalPersonNameLength
	}

	// TODO: does the legal person name type code need to be validated?
	return nil
}

func (r *RegistrationForm) validateLegalPersonLocalName(name *ivms101.LocalLegalPersonNameId) error {
	// Validate the name identifier
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

	// Validate each non-zero contact
	contacts := 0
	iter := models.NewContactIterator(r.Contacts, false, false)
	for iter.Next() {
		contacts++
		contact, field := iter.Value()
		if !models.ContactIsZero(contact) {
			err, _ = err.Append(ValidateContact(contact, FieldContacts+"."+field))
		}
	}

	// Check that all required contacts are present (special rules)
	switch contacts {
	case 0:
		// At least one contact is required
		err = append(err, &ValidationError{Field: FieldContacts, Err: ErrNoContacts.Error()})
	case 1:
		// If there is only one contact, it must be the admin; if not highlight the missing fields
		if models.ContactIsZero(r.Contacts.Administrative) {
			// Global contact error
			err = append(err, &ValidationError{Field: FieldContacts, Err: ErrMissingContact.Error()})
			switch {
			case !models.ContactIsZero(r.Contacts.Technical):
				// If the technical contact is filled in then nominate the legal/admin contact to be populated
				err = append(err, &ValidationError{Field: FieldContactsAdministrative, Err: ErrMissingAdminOrLegal.Error()})
				err = append(err, &ValidationError{Field: FieldContactsLegal, Err: ErrMissingAdminOrLegal.Error()})
			case !models.ContactIsZero(r.Contacts.Legal):
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
		if models.ContactIsZero(r.Contacts.Administrative) && models.ContactIsZero(r.Contacts.Technical) {
			err = append(err, &ValidationError{Field: FieldContacts, Err: ErrMissingContact.Error()})
			err = append(err, &ValidationError{Field: FieldContactsAdministrative, Err: ErrMissingAdminOrTechnical.Error()})
			err = append(err, &ValidationError{Field: FieldContactsTechnical, Err: ErrMissingAdminOrTechnical.Error()})
		}
		// Admin or legal must be present
		if models.ContactIsZero(r.Contacts.Administrative) && models.ContactIsZero(r.Contacts.Legal) {
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

	phone := strings.TrimSpace(contact.Phone)
	if phone == "" {
		err = append(err, &ValidationError{Field: fieldName + ".phone", Err: ErrMissingField.Error()})
	}

	// TODO: Ensure this is a valid phone number

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

	// TODO: Validate country name or ISO-3166-1 code
	if strings.TrimSpace(r.Trixo.PrimaryNationalJurisdiction) == "" {
		err = append(err, &ValidationError{
			Field: FieldTRIXOPrimaryNationalJurisdiction,
			Err:   ErrMissingField.Error(),
		})
	}

	if strings.TrimSpace(r.Trixo.PrimaryRegulator) == "" {
		err = append(err, &ValidationError{
			Field: FieldTRIXOPrimaryRegulator,
			Err:   ErrMissingField.Error(),
		})
	}

	finTransfers := strings.ToLower(strings.TrimSpace(r.Trixo.FinancialTransfersPermitted))
	if finTransfers == "" {
		err = append(err, &ValidationError{
			Field: FieldTRIXOFinancialTransfersPermitted,
			Err:   ErrMissingField.Error(),
		})
	} else if finTransfers != "yes" && finTransfers != "no" && finTransfers != "partially" {
		err = append(err, &ValidationError{
			Field: FieldTRIXOFinancialTransfersPermitted,
			Err:   ErrInvalidField.Error(),
		})
	}

	for i, juris := range r.Trixo.OtherJurisdictions {
		// TODO: Validate country name or ISO-3166-1 code
		if strings.TrimSpace(juris.Country) == "" {
			err = append(err, &ValidationError{
				Field: FieldTRIXOOtherJurisdictionsCountry,
				Err:   ErrMissingField.Error(),
				Index: i,
			})
		}

		if strings.TrimSpace(juris.RegulatorName) == "" {
			err = append(err, &ValidationError{
				Field: FieldTRIXOOtherJurisdictionsRegulatorName,
				Err:   ErrMissingField.Error(),
				Index: i,
			})
		}

		if strings.TrimSpace(juris.LicenseNumber) == "" {
			err = append(err, &ValidationError{
				Field: FieldTRIXOOtherJurisdictionsLicenseNumber,
				Err:   ErrMissingField.Error(),
				Index: i,
			})
		}
	}

	hasReg := strings.ToLower(strings.TrimSpace(r.Trixo.HasRequiredRegulatoryProgram))
	if hasReg == "" {
		err = append(err, &ValidationError{
			Field: FieldTRIXOHasRequiredRegulatoryProgram,
			Err:   ErrMissingField.Error(),
		})
	} else if hasReg != "yes" && hasReg != "no" {
		err = append(err, &ValidationError{
			Field: FieldTRIXOHasRequiredRegulatoryProgram,
			Err:   ErrInvalidField.Error(),
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
		if strings.TrimSpace(r.Trixo.KycThresholdCurrency) == "" {
			err = append(err, &ValidationError{
				Field: FieldTRIXOKYCThresholdCurrency,
				Err:   ErrMissingField.Error(),
			})
		}
	}

	if r.Trixo.MustComplyTravelRule {
		if len(r.Trixo.ApplicableRegulations) == 0 {
			err = append(err, &ValidationError{
				Field: FieldTRIXOApplicableRegulations,
				Err:   ErrMissingField.Error(),
			})
		}

		for i, reg := range r.Trixo.ApplicableRegulations {
			if strings.TrimSpace(reg) == "" {
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

		if strings.TrimSpace(r.Trixo.ComplianceThresholdCurrency) == "" {
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
	// TODO: implement
	return nil
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
