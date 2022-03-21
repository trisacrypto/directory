package api

import (
	"context"
)

//===========================================================================
// Service Interface
//===========================================================================

type BFFClient interface {
	Status(ctx context.Context, in *StatusParams) (out *StatusReply, err error)
	Lookup(ctx context.Context, in *LookupParams) (out *LookupReply, err error)
	Register(ctx context.Context, in *RegisterRequest) (out *RegisterReply, err error)
}

//===========================================================================
// Top Level Requests and Responses
//===========================================================================

// Reply contains standard fields that are used for generic API responses and errors
type Reply struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty" yaml:"error,omitempty"`
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

//===========================================================================
// BFF v1 API Requests and Responses
//===========================================================================

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
// will have either vaspdirectory.net or trisatest.net inside of it - which can be used
// to identify which network the record is associated with. The protocol buffers may
// also contain an "error" field - the BFF will handle this field by logging the error
// but will exclude it from any results returned.
type LookupReply struct {
	Results []map[string]interface{} `json:"results"`
}

// RegisterRequest is converted into a protocol buffer RegisterRequest to send to the
// specified GDS in v1.4. In v1.5 this will be expanded to include fields for both
// TestNet and MainNet registration (e.g. multiple fields for endpoint and domain).
// All generic fields (e.g. map[string]interface{}) should match their protobuf spec.
type RegisterRequest struct {
	// The network of the directory to register with, e.g. testnet or mainnet
	// NOTE: not required for front-end, used to determine endpoint in the client.
	Network string `json:"network,omitempty"`

	// RegisterRequest ProtocolBuffer fields
	Entity           map[string]interface{} `json:"entity"`
	Contacts         map[string]interface{} `json:"contacts"`
	TRISAEndpoint    string                 `json:"trisa_endpoint"`
	CommonName       string                 `json:"common_name"`
	Website          string                 `json:"website"`
	BusinessCategory string                 `json:"business_category"`
	VASPCategories   []string               `json:"vasp_categories"`
	EstablishedOn    string                 `json:"established_on"`
	TRIXO            map[string]interface{} `json:"trixo"`
}

// RegisterReply is converted from a protocol buffer RegisterReply.
type RegisterReply struct {
	Error               map[string]interface{} `json:"error"`
	Id                  string                 `json:"id"`
	RegisteredDirectory string                 `json:"registered_directory"`
	CommonName          string                 `json:"common_name"`
	Status              string                 `json:"status"`
	Message             string                 `json:"message"`
	PKCS12Password      string                 `json:"pkcs12password"`
}
