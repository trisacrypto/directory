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
	Status(ctx context.Context) (out *StatusReply, err error)
	ListVASPs(ctx context.Context, params *ListVASPsParams) (out *ListVASPsReply, err error)
	Review(ctx context.Context, in *ReviewRequest) (out *ReviewReply, err error)
	Resend(ctx context.Context, in *ResendRequest) (out *ResendReply, err error)
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

// ListVASPsParams is a request-like struct that passes query params to the ListVASPs
// GET request. All query params are optional and modify how and what data is retrieved.
type ListVASPsParams struct {
	Page     int `url:"page,omitempty"`      // defaults to page 1 if not included
	PageSize int `url:"page_size,omitempty"` // defaults to 100 if not included
}

// ListVASPsReply contains a summary data structure of all VASPs managed by the directory.
// The list reply contains standard pagination information, including the count of all
// VASPs and the links to previous and next.
type ListVASPsReply struct {
	VASPs    []VASPSnippet `json:"vasps"`
	Count    int           `json:"count"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
	Next     string        `json:"next,omitempty"`
	Previous string        `json:"previous,omitempty"`
}

type VASPSnippet struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	CommonName       string   `json:"common_name"`
	Status           string   `json:"status,omitempty"`
	LastUpdated      string   `json:"last_updated,omitempty"`
	Traveler         bool     `json:"traveler"`
	VerifiedContacts []string `json:"verified_contacts"`
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
