package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
	members "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
)

//===========================================================================
// Service Interface
//===========================================================================

type BFFClient interface {
	// Unauthenticated Endpoints
	Status(context.Context, *StatusParams) (*StatusReply, error)
	Lookup(context.Context, *LookupParams) (*LookupReply, error)
	LookupAutocomplete(context.Context) (map[string]string, error)
	VerifyContact(context.Context, *VerifyContactParams) (*VerifyContactReply, error)
	NetworkActivity(context.Context) (*NetworkActivityReply, error)

	// User Management Endpoints
	Login(context.Context, *LoginParams) error
	ListUserRoles(context.Context) ([]string, error)

	// Authenticated Endpoints
	UpdateUser(context.Context, *UpdateUserParams) error
	UserOrganization(context.Context) (*OrganizationReply, error)

	// Organization management
	CreateOrganization(context.Context, *OrganizationParams) (*OrganizationReply, error)
	DeleteOrganization(_ context.Context, id string) error
	PatchOrganization(_ context.Context, id string, request *OrganizationParams) (*OrganizationReply, error)
	ListOrganizations(context.Context, *ListOrganizationsParams) (*ListOrganizationsReply, error)

	// Collaborators endpoint
	AddCollaborator(context.Context, *models.Collaborator) (*models.Collaborator, error)
	ListCollaborators(context.Context) (*ListCollaboratorsReply, error)
	UpdateCollaboratorRoles(_ context.Context, id string, request *UpdateRolesParams) (*models.Collaborator, error)
	DeleteCollaborator(_ context.Context, id string) error

	MemberList(context.Context, *MemberPageInfo) (*MemberListReply, error)
	MemberDetails(context.Context, *MemberDetailsParams) (*MemberDetailsReply, error)

	// Registration form
	LoadRegistrationForm(context.Context, *RegistrationFormParams) (*RegistrationForm, error)
	SaveRegistrationForm(context.Context, *RegistrationForm) (*RegistrationForm, error)
	ResetRegistrationForm(context.Context, *RegistrationFormParams) (*RegistrationForm, error)
	SubmitRegistration(_ context.Context, network string) (*RegisterReply, error)
	RegistrationStatus(context.Context) (*RegistrationStatus, error)

	// Overview and announcements
	Overview(context.Context) (*OverviewReply, error)
	Announcements(context.Context) (*AnnouncementsReply, error)
	MakeAnnouncement(context.Context, *models.Announcement) error
	Attention(context.Context) (*AttentionReply, error)

	// Certificate management
	Certificates(context.Context) (*CertificatesReply, error)
}

//===========================================================================
// Top Level Requests and Responses
//===========================================================================

// Reply contains standard fields that are used for generic API responses and errors
type Reply struct {
	Success      bool   `json:"success"`
	Error        string `json:"error,omitempty" yaml:"error,omitempty"`
	RefreshToken bool   `json:"refresh_token,omitempty" yaml:"refresh_token,omitempty"`
}

// StatusParams is parsed from the query parameters of the GET request
type StatusParams struct {
	NoGDS bool `url:"nogds,omitempty" form:"nogds" default:"false"`
}

// StatusReply is returned on status requests. Note that no request is needed.
type StatusReply struct {
	Status  string `json:"status"`
	Uptime  string `json:"uptime,omitempty"`
	Version string `json:"version,omitempty"`
	TestNet string `json:"testnet,omitempty"`
	MainNet string `json:"mainnet,omitempty"`
}

// A per-field validation error that is intended for human consumption - if the field is
// not valid (e.g. empty when required, doesn't match regular expression, etc.) then
// this struct is meant to be sent back so the front-end can render the message to the
// user in a help-box or similar. If the field is an array element, then the index field
// will contain the index of the erroring element.
type FieldValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
	Index int    `json:"index"`
}

func (f *FieldValidationError) String() string {
	return fmt.Sprintf("%s: %s", f.Field, f.Error)
}

func NewFieldValidationError(err error) *FieldValidationError {
	var verr *models.ValidationError
	if errors.As(err, &verr) {
		return &FieldValidationError{Field: verr.Field, Error: verr.Err, Index: verr.Index}
	}
	return &FieldValidationError{Error: err.Error()}
}

func FromValidationErrors(err error) []*FieldValidationError {
	var verrs models.ValidationErrors
	if errors.As(err, &verrs) {
		out := make([]*FieldValidationError, 0, len(verrs))
		for _, verr := range verrs {
			out = append(out, NewFieldValidationError(verr))
		}
		return out
	}

	out := make([]*FieldValidationError, 0, 1)
	return append(out, NewFieldValidationError(err))
}

//===========================================================================
// BFF v1 API Requests and Responses
//===========================================================================

// UpdateUserParams is used to update the user's profile information.
type UpdateUserParams struct {
	Name string `json:"name,omitempty"`
}

// OrganizationParams is used to create and update organizations.
type OrganizationParams struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

// OrganizationReply contains high level information about an organization.
type OrganizationReply struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Domain       string `json:"domain"`
	CreatedAt    string `json:"created_at"`
	LastLogin    string `json:"last_login"`
	RefreshToken bool   `json:"refresh_token,omitempty"`
}

// ListOrganizationsParams contains query parameters for listing organizations.
type ListOrganizationsParams struct {
	Name     string `url:"name,omitempty" form:"name"`
	Page     int    `url:"page,omitempty" form:"page" default:"1"`
	PageSize int    `url:"page_size,omitempty" form:"page_size" default:"8"`
}

// ListOrganizationsReply contains a page of organizations.
type ListOrganizationsReply struct {
	Organizations []*OrganizationReply `json:"organizations"`
	Count         int                  `json:"count"`
	Page          int                  `json:"page"`
	PageSize      int                  `json:"page_size"`
}

// LoginParams contains additional information needed for post-authentication checks
// during user login.
type LoginParams struct {
	OrgID string `json:"orgid"`
}

// UpdateRolesParams contains a list of new roles for a collaborator.
type UpdateRolesParams struct {
	Roles []string `json:"roles"`
}

// ListCollaboratorsReply contains a list of collaborators.
type ListCollaboratorsReply struct {
	Collaborators []*models.Collaborator `json:"collaborators"`
}

// LookupParams is converted into a GDS LookupRequest.
type LookupParams struct {
	ID         string `url:"uuid,omitempty" form:"uuid"`
	CommonName string `url:"common_name,omitempty" form:"common_name"`
}

// LookupReply can return 1-2 results either one result found from one directory
// service or results found from both TestNet and MainNet. If no results are found, the
// Lookup endpoint returns a 404 error (not found). The result is the simplest case,
// just a JSON serialization of the protocol buffers returned from GDS to help long term
// maintainability. The protocol buffers contain a "registered_directory" field that
// will have either trisa.directory or testnet.directory inside of it - which can be used
// to identify which network the record is associated with. The protocol buffers may
// also contain an "error" field - the BFF will handle this field by logging the error
// but will exclude it from any results returned.
type LookupReply struct {
	TestNet map[string]interface{} `json:"testnet"`
	MainNet map[string]interface{} `json:"mainnet"`
}

// VerifyContactParams is converted into a GDS VerifyContactRequest.
type VerifyContactParams struct {
	ID        string `url:"vaspID,omitempty" form:"vaspID"`
	Token     string `url:"token,omitempty" form:"token"`
	Directory string `url:"registered_directory,omitempty" form:"registered_directory"`
}

// VerifyContactReply
type VerifyContactReply struct {
	Error   map[string]interface{} `json:"error,omitempty"`
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
}

type RegistrationFormStep string

const (
	StepBasicDetails RegistrationFormStep = "basic"
	StepLegalPerson  RegistrationFormStep = "legal"
	StepContacts     RegistrationFormStep = "contacts"
	StepTRISA        RegistrationFormStep = "trisa"
	StepTRIXO        RegistrationFormStep = "trixo"
)

// Allows the front-end to specify which part of the registration form they want to
// fetch or delete.
// GET /v1/registration will return the entire registration form, while
// GET /v1/registration?step=trixo would return just the TRIXO form
// DELETE /v1/registration will reset the entire registration form, while
// DELETE /v1/registration?step=trixo would reset just the TRIXO form
type RegistrationFormParams struct {
	Step RegistrationFormStep `url:"step,omitempty" form:"step"`
}

// RegistrationForm is a wrapper around the models.RegistrationForm that includes API-
// specific details such as the step and field validation errors.
type RegistrationForm struct {
	Step   RegistrationFormStep     `json:"step,omitempty"`
	Form   *models.RegistrationForm `json:"form"`
	Errors []*FieldValidationError  `json:"errors,omitempty"`
}

// MarshalStepJSON removes any unnecessary fields from the registration form.
func (r *RegistrationForm) MarshalStepJSON() (_ gin.H, err error) {
	// Marshal everything but the registration form
	form := r.Form
	r.Form = nil
	defer func() {
		// Reset the form
		r.Form = form
	}()

	var data []byte
	if data, err = json.Marshal(r); err != nil {
		return nil, err
	}

	var intermediate gin.H
	if err = json.Unmarshal(data, &intermediate); err != nil {
		return nil, err
	}

	var step models.StepType
	if step, err = models.ParseStepType(string(r.Step)); err != nil {
		return nil, err
	}

	// Marshal the registration form with the step
	if intermediate["form"], err = form.MarshalStep(step); err != nil {
		return nil, err
	}
	return intermediate, nil
}

// RegisterReply is converted from a protocol buffer RegisterReply.
type RegisterReply struct {
	Error               map[string]interface{} `json:"error,omitempty"`
	Id                  string                 `json:"id"`
	RegisteredDirectory string                 `json:"registered_directory"`
	CommonName          string                 `json:"common_name"`
	Status              string                 `json:"status"`
	Message             string                 `json:"message"`
	PKCS12Password      string                 `json:"pkcs12password"`
	RefreshToken        bool                   `json:"refresh_token,omitempty"`
}

// RegistrationStatus is returned on registration status requests. This will contain
// RFC3339 formatted timestamps indicating when the registration was submitted for
// testnet and mainnet.
type RegistrationStatus struct {
	TestNetSubmitted string `json:"testnet_submitted,omitempty"`
	MainNetSubmitted string `json:"mainnet_submitted,omitempty"`
}

// OverviewReply is returned on overview requests.
type OverviewReply struct {
	Error   NetworkError    `json:"error,omitempty"`
	OrgID   string          `json:"org_id"`
	TestNet NetworkOverview `json:"testnet"`
	MainNet NetworkOverview `json:"mainnet"`
}

// NetworkOverview contains network-specific information.
type NetworkOverview struct {
	Status             string        `json:"status"`
	Vasps              int           `json:"vasps"`
	CertificatesIssued int           `json:"certificates_issued"`
	NewMembers         int           `json:"new_members"`
	MemberDetails      MemberDetails `json:"member_details"`
}

// MemberDetails contains VASP-specific information.
type MemberDetails struct {
	ID          string                 `json:"id"`
	Status      string                 `json:"status"`
	CountryCode string                 `json:"country_code"`
	FirstListed string                 `json:"first_listed"`
	VerifiedOn  string                 `json:"verified_on"`
	LastUpdated string                 `json:"last_updated"`
	Certificate map[string]interface{} `json:"certificate"`
}

// AnnouncementsReply contains up to the last 10 network announcements that were made in
// the past month. It does not require pagination since only relevant results are returned.
type AnnouncementsReply struct {
	Announcements []*models.Announcement `json:"announcements"`
	LastUpdated   string                 `json:"last_updated,omitempty"`
}

// CertificatesReply is returned on certificates requests.
type CertificatesReply struct {
	Error   NetworkError  `json:"network_error,omitempty"`
	TestNet []Certificate `json:"testnet"`
	MainNet []Certificate `json:"mainnet"`
}

// Certificate contains details about a certificate issued to a VASP.
type Certificate struct {
	SerialNumber string                 `json:"serial_number"`
	IssuedAt     string                 `json:"issued_at"`
	ExpiresAt    string                 `json:"expires_at"`
	Revoked      bool                   `json:"revoked"`
	Details      map[string]interface{} `json:"details"`
}

// MembersPageInfo enables paginated requests to the TRISAMembers/List RPC for the
// specified directory. Pagination is not stateful and requires a token.
type MemberPageInfo struct {
	Directory string `url:"registered_directory,omitempty" form:"registered_directory"`
	PageSize  int32  `url:"page_size,omitempty" form:"page_size"`
	PageToken string `url:"page_token,omitempty" form:"page_token"`
}

type MemberListReply struct {
	VASPs         []*members.VASPMember `json:"vasps"`
	NextPageToken string                `json:"next_page_token,omitempty"`
}

// MemberDetailsParams contains details required to identify a VASP member in a specific
// registered directory (e.g. testnet.directory or trisa.directory).
type MemberDetailsParams struct {
	ID        string `url:"-" form:"-"`
	Directory string `url:"registered_directory,omitempty" form:"registered_directory"`
}

// MemberDetailsReply contains sensitive details about a VASP member.
type MemberDetailsReply struct {
	Summary     map[string]interface{} `json:"summary"`
	LegalPerson map[string]interface{} `json:"legal_person"`
	Contacts    map[string]interface{} `json:"contacts"`
	Trixo       map[string]interface{} `json:"trixo"`
}

// AttentionReply contains all the current attention messages relevant to an
// organization.
type AttentionReply struct {
	Messages []*AttentionMessage `json:"messages"`
}

// AttentionMessage contains details about a single attention message.
type AttentionMessage struct {
	Message  string `json:"message"`
	Severity string `json:"severity"`
	Action   string `json:"action"`
}

// NetworkError is populated when the BFF receives an error from a network endpoint,
// containing an error string for each network that errored. This allows the client to
// distinguish between network errors and BFF errors and determine which network the
// errors originated from.
type NetworkError struct {
	TestNet string `json:"testnet,omitempty"`
	MainNet string `json:"mainnet,omitempty"`
}

// Activity is a time-aggregated collection of events (Search, Lookup, etc)
type Activity struct {
	Date   string `json:"date"`
	Events uint64 `json:"events"`
}

// NetworkActivityReply is a map of the network (TestNet or MainNet) to Activity,
// which is a time-aggregated collection of events (Search, Lookup, etc)
// that occurred on the network
type NetworkActivityReply struct {
	TestNet []Activity `json:"testnet"`
	MainNet []Activity `json:"mainnet"`
}
