package models

import (
	"errors"
	"fmt"
	"strings"

	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/encoding/protojson"
)

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
func (r *RegistrationForm) ValidateBasicDetails() (v ValidationErrors) {
	// Validate website
	// TODO: Check if this is a valid URL?
	if strings.TrimSpace(r.Website) == "" {
		v = append(v, &ValidationError{
			Field: "website",
			Err:   ErrMissingField.Error(),
		})
	}

	// Validate business category
	if r.BusinessCategory == pb.BusinessCategory_UNKNOWN_ENTITY {
		v = append(v, &ValidationError{
			Field: "business_category",
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
		v = append(v, &ValidationError{
			Field: "vasp_categories",
			Err:   ErrMissingField.Error(),
		})
	}

	// Validate established date
	// TODO: Check if this is a valid date?
	if strings.TrimSpace(r.EstablishedOn) == "" {
		v = append(v, &ValidationError{
			Field: "established_on",
			Err:   ErrMissingField.Error(),
		})
	}

	// Validate organization name
	if strings.TrimSpace(r.OrganizationName) == "" {
		v = append(v, &ValidationError{
			Field: "organization_name",
			Err:   ErrMissingField.Error(),
		})
	}

	return v
}

// Validate only the fields in the legal person step.
func (r *RegistrationForm) ValidateLegalPerson() ValidationErrors {
	// TODO: implement
	return nil
}

// Validate only the fields in the contacts step.
func (r *RegistrationForm) ValidateContacts() ValidationErrors {
	// TODO: implement
	return nil
}

// Validate only the fields in the trixo step.
func (r *RegistrationForm) ValidateTRIXO() (v ValidationErrors) {
	if r.Trixo == nil {
		v = append(v, &ValidationError{
			Field: "trixo",
			Err:   ErrMissingField.Error(),
		})
		return v
	}

	// TODO: Validate country name or ISO-3166-1 code
	if strings.TrimSpace(r.Trixo.PrimaryNationalJurisdiction) == "" {
		v = append(v, &ValidationError{
			Field: "trixo.primary_national_jurisdiction",
			Err:   ErrMissingField.Error(),
		})
	}

	if strings.TrimSpace(r.Trixo.PrimaryRegulator) == "" {
		v = append(v, &ValidationError{
			Field: "trixo.primary_regulator",
			Err:   ErrMissingField.Error(),
		})
	}

	finTransfers := strings.ToLower(strings.TrimSpace(r.Trixo.FinancialTransfersPermitted))
	if finTransfers == "" {
		v = append(v, &ValidationError{
			Field: "trixo.financial_transfers_permitted",
			Err:   ErrMissingField.Error(),
		})
	} else if finTransfers != "yes" && finTransfers != "no" && finTransfers != "partially" {
		v = append(v, &ValidationError{
			Field: "trixo.financial_transfers_permitted",
			Err:   ErrInvalidField.Error(),
		})
	}

	for i, juris := range r.Trixo.OtherJurisdictions {
		// TODO: Validate country name or ISO-3166-1 code
		if strings.TrimSpace(juris.Country) == "" {
			v = append(v, &ValidationError{
				Field: "trixo.other_jurisdictions.country",
				Err:   ErrMissingField.Error(),
				Index: i,
			})
		}

		if strings.TrimSpace(juris.RegulatorName) == "" {
			v = append(v, &ValidationError{
				Field: "trixo.other_jurisdictions.regulator_name",
				Err:   ErrMissingField.Error(),
				Index: i,
			})
		}

		if strings.TrimSpace(juris.LicenseNumber) == "" {
			v = append(v, &ValidationError{
				Field: "trixo.other_jurisdictions.license_number",
				Err:   ErrMissingField.Error(),
				Index: i,
			})
		}
	}

	hasReg := strings.ToLower(strings.TrimSpace(r.Trixo.HasRequiredRegulatoryProgram))
	if hasReg == "" {
		v = append(v, &ValidationError{
			Field: "trixo.has_required_regulatory_program",
			Err:   ErrMissingField.Error(),
		})
	} else if hasReg != "yes" && hasReg != "no" {
		v = append(v, &ValidationError{
			Field: "trixo.has_required_regulatory_program",
			Err:   ErrInvalidField.Error(),
		})
	}

	if r.Trixo.ConductsCustomerKyc {
		if r.Trixo.KycThreshold < 0 {
			v = append(v, &ValidationError{
				Field: "trixo.kyc_threshold",
				Err:   ErrNegativeValue.Error(),
			})
		}

		// TODO: Validate currency code
		if strings.TrimSpace(r.Trixo.KycThresholdCurrency) == "" {
			v = append(v, &ValidationError{
				Field: "trixo.kyc_threshold_currency",
				Err:   ErrMissingField.Error(),
			})
		}
	}

	if r.Trixo.MustComplyTravelRule {
		if len(r.Trixo.ApplicableRegulations) == 0 {
			v = append(v, &ValidationError{
				Field: "trixo.applicable_regulations",
				Err:   ErrMissingField.Error(),
			})
		}

		for i, reg := range r.Trixo.ApplicableRegulations {
			if strings.TrimSpace(reg) == "" {
				v = append(v, &ValidationError{
					Field: "trixo.applicable_regulations",
					Err:   ErrMissingField.Error(),
					Index: i,
				})
			}
		}

		if r.Trixo.ComplianceThreshold < 0 {
			v = append(v, &ValidationError{
				Field: "trixo.compliance_threshold",
				Err:   ErrNegativeValue.Error(),
			})
		}

		if strings.TrimSpace(r.Trixo.ComplianceThresholdCurrency) == "" {
			v = append(v, &ValidationError{
				Field: "trixo.compliance_threshold_currency",
				Err:   ErrMissingField.Error(),
			})
		}
	}

	return v
}

// Validate only the fields in the trisa step.
func (r *RegistrationForm) ValidateTRISA() ValidationErrors {
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
