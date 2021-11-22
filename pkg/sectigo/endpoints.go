package sectigo

import (
	"fmt"
	"net/url"
)

const (
	contentType         = "application/json;charset=UTF-8"
	downloadContentType = "application/octet-stream"
	userAgent           = "TRISA GDS Sectigo Client v1.1"
)

// baseURL is used to construct API endpoints for Sectigo methods
var defaultURL = &url.URL{Scheme: "https", Host: "iot.sectigo.com"}
var baseURL = defaultURL

// endpoints maps methods to URLs, which are fully constructed using the baseURL. Some
// endpoint paths contain string format verbs intended for dynamic REST urls by making a
// copy of the url.URL and replacing the Path with the output of
// fmt.Sprintf(endpoints[method].Path, param). If the format verb is not replaced in
// this way, it will be encoded as %25<verb> which will likely return a 404 error. Note
// that value copying a url.URL is not a deep copy and the internal user info will not
// be copied.
//
// Many endpoints require URL parameters. See the API documentation for more details.
// Endpoint URLs specify paths only, copy the URL and add URL parameters for dynamic
// endpoints.
var endpoints = map[string]*url.URL{
	AuthenticateEP:                  {Path: "/auth/pwd"},
	RefreshEP:                       {Path: "/auth/refresh"},
	BatchesEP:                       {Path: "/api/v1/batches"},
	BatchDetailEP:                   {Path: "/api/v1/batches/%d"},
	BatchStatusEP:                   {Path: "/api/v1/batches/%d/status"},
	BatchAuditLogEP:                 {Path: "/api/v1/batches/%d/auditLog"},
	BatchProcessingInfoEP:           {Path: "/api/v1/batches/%d/processing_info"},
	BatchDevicesAuditLogEP:          {Path: "/api/v1/batches/%d/devices/auditLog"},
	BatchPreviewEP:                  {Path: "/api/v1/batches/preview"},
	CreateSingleCertBatchEP:         {Path: "/api/v1/batches/createSingleCertBatch"},
	UploadEP:                        {Path: "/api/v1/batches/upload"},
	UploadCSVEP:                     {Path: "/api/v2/organizations/%d/profiles/%d/batches/csv-upload"},
	GeneratorsEP:                    {Path: "/api/v1/generators"},
	DownloadEP:                      {Path: "/api/v1/batches/%d/download"},
	DevicesEP:                       {Path: "/api/v1/devices"},
	UserAuthoritiesEP:               {Path: "/api/v1/authorities/allowed"},
	AuthorityBalanceUsedEP:          {Path: "/api/v1/authorities/%d/balanceused/%d"},
	AuthorityBalanceAvailableEP:     {Path: "/api/v1/authorities/%d/balanceavailable/%d"},
	AuthorityUserBalanceAvailableEP: {Path: "/api/v1/authorities/%d/balanceavailable"},
	UsersEP:                         {Path: "/api/v1/users"},
	CheckEmailEP:                    {Path: "/api/v1/users/check_email"},
	UserDetailEP:                    {Path: "/api/v1/users/%d"},
	UserProfilesEP:                  {Path: "/api/v1/users/profiles"},
	UsersOrganizationsEP:            {Path: "/api/v1/users/organizations/%d"},
	RemoveUserEP:                    {Path: "/api/v1/users/remove_user"},
	CurrentUserEP:                   {Path: "/api/v1/users/current"},
	UpdateEmailEP:                   {Path: "/api/v1/users/update_email"},
	CheckUserEP:                     {Path: "/api/v1/users/check_user"},
	UserCredentialsEP:               {Path: "/api/v1/users/credentials"},
	ProfilesEP:                      {Path: "/api/v1/profiles"},
	ProfileDetailEP:                 {Path: "/api/v1/profiles/%d"},
	ProfileBalanceEP:                {Path: "/api/v1/profiles/%d/balance"},
	ProfileParametersEP:             {Path: "/api/v1/profiles/%d/parameters"},
	ProfileSubjectDNEP:              {Path: "/api/v1/profiles/%d/dn"},
	OrganizationsEP:                 {Path: "/api/v1/organizations"},
	OrganizationDetailEP:            {Path: "/api/v1/organizations/%d"},
	CurrentUserOrganizationEP:       {Path: "/api/v1/organizations/user"},
	CheckOrganizationEP:             {Path: "/api/v1/organizations/check_organization"},
	UpdateAuthorityEP:               {Path: "/api/v1/organizations/%d/authorities/%d"},
	OrganizationAuthoritiesEP:       {Path: "/api/v1/organizations/%d/all_authorities"},
	AuthorityDetailEP:               {Path: "/api/v1/organizations/%d/authority/%d"},
	AuthoritiesEP:                   {Path: "/api/v1/organizations/%d/authority"},
	OrganizationListItemsEP:         {Path: "/api/v1/organizations/select/items"},
	OrganizationParametersEP:        {Path: "/api/v1/organizations/profile_parameters/%d"},
	EcosystemsEP:                    {Path: "/api/v1/ecosystems"},
	UserEcosystemEP:                 {Path: "/api/v1/ecosystems/ecosystem"},
	EcosystemBalanceEP:              {Path: "/api/v1/ecosystems/ecosystem/balance"},
	EcosystemsStatisticsEP:          {Path: "/api/v1/ecosystems/statistics"},
	EcosystemAdminDetailEP:          {Path: "/api/v1/ecosystems/users/%d"},
	EcosystemAdminsEP:               {Path: "/api/v1/ecosystems/users"},
	FindCertificateEP:               {Path: "/api/v1/certificates/find"},
	RevokeDeviceCertificateEP:       {Path: "/api/v1/certificates/revoke"},
	RevokeCertificateEP:             {Path: "/api/v1/certificates/%d/revoke"},
}

// endpoint name constants to prevent typos at compile time rather than at runtime.
const (
	AuthenticateEP                  = "authenticate"
	RefreshEP                       = "refresh"
	BatchesEP                       = "batches"
	BatchDetailEP                   = "batchDetail"
	BatchStatusEP                   = "batchStatus"
	BatchAuditLogEP                 = "batchAuditLog"
	BatchProcessingInfoEP           = "batchProcessingInfo"
	BatchDevicesAuditLogEP          = "batchDevicesAuditLog"
	BatchPreviewEP                  = "batchPreview"
	CreateSingleCertBatchEP         = "createSingleCertBatch"
	UploadEP                        = "upload"
	UploadCSREP                     = UploadEP
	UploadCSVEP                     = "uploadCSV"
	GeneratorsEP                    = "generators"
	DownloadEP                      = "download"
	DevicesEP                       = "devices"
	UserAuthoritiesEP               = "userAuthorities"
	AuthorityBalanceUsedEP          = "authorityBalanceUsed"
	AuthorityBalanceAvailableEP     = "authorityBalanceAvailable"
	AuthorityUserBalanceAvailableEP = "authorityUserBalanceAvailable"
	UsersEP                         = "users"
	CheckEmailEP                    = "checkEmail"
	UserDetailEP                    = "userDetail"
	UserProfilesEP                  = "userProfiles"
	UsersOrganizationsEP            = "usersOrganizations"
	RemoveUserEP                    = "removeUser"
	CurrentUserEP                   = "currentUser"
	UpdateEmailEP                   = "updateEmail"
	CheckUserEP                     = "checkUser"
	UserCredentialsEP               = "userCredentials"
	ProfilesEP                      = "profiles"
	ProfileDetailEP                 = "profileDetail"
	ProfileBalanceEP                = "profileBalance"
	ProfileParametersEP             = "profileParameters"
	ProfileSubjectDNEP              = "profileSubjectDN"
	OrganizationsEP                 = "organizations"
	OrganizationDetailEP            = "organizationDetail"
	CurrentUserOrganizationEP       = "currentUserOrganization"
	CheckOrganizationEP             = "checkOrganization"
	UpdateAuthorityEP               = "updateAuthority"
	OrganizationAuthoritiesEP       = "organizationAuthorities"
	AuthorityDetailEP               = "authorityDetail"
	AuthoritiesEP                   = "authorities"
	OrganizationListItemsEP         = "organizationListItems"
	OrganizationParametersEP        = "organizationParameters"
	EcosystemsEP                    = "ecosystems"
	UserEcosystemEP                 = "userEcosystem"
	EcosystemBalanceEP              = "ecosystemBalance"
	EcosystemsStatisticsEP          = "ecosystemsStatistics"
	EcosystemAdminDetailEP          = "ecosystemAdminDetail"
	EcosystemAdminsEP               = "ecosystemAdmins"
	FindCertificateEP               = "findCertificate"
	RevokeDeviceCertificateEP       = "revokeDeviceCertificate"
	RevokeCertificateEP             = "revokeCertificate"
)

// Endpoint returns the resolved URL for the named endpoint
func Endpoint(endpoint string, params ...interface{}) (u *url.URL, err error) {
	var ok bool
	if u, ok = endpoints[endpoint]; !ok {
		return nil, fmt.Errorf("no endpoint named %q", endpoint)
	}

	// Construct the full URL from the base URL and endpoint path.
	v := baseURL.ResolveReference(u)

	if len(params) > 0 {
		v.Path = fmt.Sprintf(v.Path, params...)
	}
	return v, nil
}

// urlFor returns the string of the resolved URL for the named endpoint for constructing
// http requests inside of client methods. If the endpoint doesn't exist, it panics
// instead of returning an error because this is a developer, not a user error.
func urlFor(endpoint string, params ...interface{}) string {
	ep, err := Endpoint(endpoint, params...)
	if err != nil {
		// This is a developer error, so we should panic (original functionality)
		panic(err)
	}
	return ep.String()
}

// SetBaseURL updates Sectigo to use a different scheme and host to determine endpoints
func SetBaseURL(u *url.URL) {
	baseURL = u
}

// ResetBaseURL to the default Sectigo API endpoint
func ResetBaseURL() {
	baseURL = defaultURL
}
