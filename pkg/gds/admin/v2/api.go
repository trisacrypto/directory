package admin

import (
	"context"
	"time"
)

//===========================================================================
// Service Interface
//===========================================================================

// DirectoryAdministrationClient defines client-side interactions with the API.
type DirectoryAdministrationClient interface {
	Login(ctx context.Context) (err error)
	Status(ctx context.Context) (out *StatusReply, err error)
	Authenticate(ctx context.Context, in *AuthRequest) (out *AuthReply, err error)
	Reauthenticate(ctx context.Context, in *AuthRequest) (out *AuthReply, err error)
	Summary(ctx context.Context) (out *SummaryReply, err error)
	Autocomplete(ctx context.Context) (out *AutocompleteReply, err error)
	ListVASPs(ctx context.Context, params *ListVASPsParams) (out *ListVASPsReply, err error)
	RetrieveVASP(ctx context.Context, id string) (out *RetrieveVASPReply, err error)
	Review(ctx context.Context, in *ReviewRequest) (out *ReviewReply, err error)
	Resend(ctx context.Context, in *ResendRequest) (out *ResendReply, err error)
	ReviewTimeline(ctx context.Context, params *ReviewTimelineParams) (out *ReviewTimelineReply, err error)
}

//===========================================================================
// Top Level Requests and Responses
//===========================================================================

// Reply contains standard fields that are used for generic API responses and errors
type Reply struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty" yaml:"error,omitempty"`
}

// StatusReply is returned on status requests. Note that no request is needed.
type StatusReply struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Version   string    `json:"version,omitempty"`
}

//===========================================================================
// Admin v2 API Requests and Responses
//===========================================================================

// AuthRequest is used by both the Authenticate and Reauthenticate API calls. In
// Authenticate, the credential should be the OAuth2 JWT token supplied by the Identity
// Service Provider. In Reauthenticate the credential should be the refresh token
// returned by the Authenticate request.
type AuthRequest struct {
	Credential string `json:"credential"`
}

// AuthReply returns access and refresh tokens. The access token should be used as a
// a Bearer token in the Authorization header for all authenticated requests including
// Reauthentication. The refresh token is used in the Reauthenticate method to retrieve
// new access tokens without requiring a new login.
type AuthReply struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// SummaryReply provides aggregate statistics that describe the state of the GDS.
type SummaryReply struct {
	VASPsCount           int            `json:"vasps_count"`           // the total number of VASPs in any state in GDS
	PendingRegistrations int            `json:"pending_registrations"` // the number of registrations pending (in any pre-review status)
	ContactsCount        int            `json:"contacts_count"`        // the number of contacts in the system
	VerifiedContacts     int            `json:"verified_contacts"`     // the number of verified contacts in the system
	CertificatesIssued   int            `json:"certificates_issued"`   // the number of certificates issued by the GDS
	Statuses             map[string]int `json:"statuses"`              // the counts of all statuses in the system
	CertReqs             map[string]int `json:"certreqs"`              // The counts of all certificate request statuses
}

// AutocompleteReply contains a mapping of name to VASP UUID for the search bar.
type AutocompleteReply struct {
	Names map[string]string `json:"names"`
}

// ListVASPsParams is a request-like struct that passes query params to the ListVASPs
// GET request. All query params are optional and modify how and what data is retrieved.
type ListVASPsParams struct {
	Status   string `url:"status,omitempty" form:"status"`
	Page     int    `url:"page,omitempty" form:"page" default:"1"`             // defaults to page 1 if not included
	PageSize int    `url:"page_size,omitempty" form:"page_size" default:"100"` // defaults to 100 if not included
}

// ListVASPsReply contains a summary data structure of all VASPs managed by the directory.
// The list reply contains standard pagination information, including the count of all
// VASPs and the links to previous and next.
type ListVASPsReply struct {
	VASPs    []VASPSnippet `json:"vasps"`
	Count    int           `json:"count"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

// VASPSnippet provides summary information about a VASP.
type VASPSnippet struct {
	ID                 string          `json:"id"`
	Name               string          `json:"name"`
	CommonName         string          `json:"common_name"`
	VerificationStatus string          `json:"verification_status,omitempty"`
	LastUpdated        string          `json:"last_updated,omitempty"`
	Traveler           bool            `json:"traveler"`
	VerifiedContacts   map[string]bool `json:"verified_contacts"`
}

// RetrieveVASPReply returns a pb.VASP record that has been marshaled by protojson and
// includes extra information such as verified contacts, whether or not the VASP is a
// Traveler node, and other pre-computed data to facilitate administrative actions. The
// serialized pb.VASP record is returned to make sure that the Admin API keeps up with
// changes in the TRISA library. Admin API developers should reference the
// trisacrypto/trisa library to ensure they have all of the requried data that is
// returned. Go developers should unmarshal the data into a *pb.VASP struct.
type RetrieveVASPReply struct {
	Name             string                 `json:"name"`
	VASP             map[string]interface{} `json:"vasp"`
	VerifiedContacts map[string]string      `json:"verified_contacts"`
	Traveler         bool                   `json:"traveler"`
}

// ReviewTimelineParams contains the start and end date for the requested timeline.
type ReviewTimelineParams struct {
	Start string `url:"start,omitempty" form:"start"`
	End   string `url:"end,omitempty" form:"end"`
}

// ReviewTimelineRecord contains counts of VASP registration states over a single week.
type ReviewTimelineRecord struct {
	Week          string         `json:"week"`
	VASPsUpdated  int            `json:"vasps_updated"`
	Registrations map[string]int `json:"registrations"`
}

// ReviewTimelineReply returns a list of time series records containing registration counts.
type ReviewTimelineReply struct {
	Weeks []ReviewTimelineRecord `json:"weeks"`
}

//===========================================================================
// VASP Action RPCs
//===========================================================================

// Registration review requests are sent via email to the TRISA admin email address with
// a lightweight token for review. This endpoint allows administrators to submit a review
// determination back to the directory server.
//
// TODO: adapt this so there is no admin verification token and authentication is used.
type ReviewRequest struct {
	// The ID of the VASP to perform the review for (optional - is part of the URL)
	ID string `json:"vasp_id,omitempty"`

	// The verification token sent in the review email.
	// Lightweight authentication should be replaced by authentication and audit logs.
	AdminVerificationToken string `json:"admin_verification_token"`

	// If accept is false the request will be rejected and a reject reason must be
	// specfied. If it is true, then the certificate issuance process will begin.
	Accept       bool   `json:"accept"`
	RejectReason string `json:"reject_reason,omitempty"`
}

// ReviewReply returns verification status of the VASP Registration.
type ReviewReply struct {
	// Status must be a valid trisa.gds.models.v1beta1.VerificationState
	Status  string `json:"status"`
	Message string `json:"message"`
}

// ResendActions to use in ResendRequests
type ResendAction string

const (
	ResendVerifyContact ResendAction = "verify_contact"
	ResendReview        ResendAction = "review"
	ResendDeliverCerts  ResendAction = "deliver_certs"
	ResendRejection     ResendAction = "rejection"
)

// ResendRequest allows extra attempts to resend emails to be made if they were not
// delivered or recieved the first time. This is a routine action that may need to be
// carried out from time to time.
type ResendRequest struct {
	// The ID of the VASP to resend emails for (optional - is part of the URL)
	ID string `json:"vasp_id,omitempty"`

	// The resend action type, must parse to a ResendAction enumeration. If the action
	// is "rejection" then a reason must be supplied for the rejection as well.
	Action ResendAction `json:"action"`
	Reason string       `json:"reason,omitempty"`
}

// ResendReply returns the number of emails sent and a status message from the server.
type ResendReply struct {
	Sent    int    `json:"sent"`
	Message string `json:"message"`
}
