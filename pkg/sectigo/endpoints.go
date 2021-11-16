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

// endpoints maps methods to URLs, which are full constructed with the baseURL. Some
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
	userProfilesEP:                  {Path: "/api/v1/users/profiles"},
	usersOrganizationsEP:            {Path: "/api/v1/users/organizations/%d"},
	removeUserEP:                    {Path: "/api/v1/users/remove_user"},
	currentUserEP:                   {Path: "/api/v1/users/current"},
	updateEmailEP:                   {Path: "/api/v1/users/update_email"},
	checkUserEP:                     {Path: "/api/v1/users/check_user"},
	userCredentialsEP:               {Path: "/api/v1/users/credentials"},
	ProfilesEP:                      {Path: "/api/v1/profiles"},
	ProfileDetailEP:                 {Path: "/api/v1/profiles/%d"},
	profileBalanceEP:                {Path: "/api/v1/profiles/%d/balance"},
	ProfileParametersEP:             {Path: "/api/v1/profiles/%d/parameters"},
	profileSubjectDNEP:              {Path: "/api/v1/profiles/%d/dn"},
	organizationsEP:                 {Path: "/api/v1/organizations"},
	organizationDetailEP:            {Path: "/api/v1/organizations/%d"},
	CurrentUserOrganizationEP:       {Path: "/api/v1/organizations/user"},
	checkOrganizationEP:             {Path: "/api/v1/organizations/check_organization"},
	updateAuthorityEP:               {Path: "/api/v1/organizations/%d/authorities/%d"},
	organizationAuthoritiesEP:       {Path: "/api/v1/organizations/%d/all_authorities"},
	AuthorityDetailEP:               {Path: "/api/v1/organizations/%d/authority/%d"},
	authoritiesEP:                   {Path: "/api/v1/organizations/%d/authority"},
	organizationListItemsEP:         {Path: "/api/v1/organizations/select/items"},
	organizationParametersEP:        {Path: "/api/v1/organizations/profile_parameters/%d"},
	ecosystemsEP:                    {Path: "/api/v1/ecosystems"},
	userEcosystemEP:                 {Path: "/api/v1/ecosystems/ecosystem"},
	ecosystemBalanceEP:              {Path: "/api/v1/ecosystems/ecosystem/balance"},
	ecosystemsStatisticsEP:          {Path: "/api/v1/ecosystems/statistics"},
	ecosystemAdminDetailEP:          {Path: "/api/v1/ecosystems/users/%d"},
	ecosystemAdminsEP:               {Path: "/api/v1/ecosystems/users"},
	FindCertificateEP:               {Path: "/api/v1/certificates/find"},
	revokeDeviceCertificateEP:       {Path: "/api/v1/certificates/revoke"},
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
	userProfilesEP                  = "userProfiles"
	usersOrganizationsEP            = "usersOrganizations"
	removeUserEP                    = "removeUser"
	currentUserEP                   = "currentUser"
	updateEmailEP                   = "updateEmail"
	checkUserEP                     = "checkUser"
	userCredentialsEP               = "userCredentials"
	ProfilesEP                      = "profiles"
	ProfileDetailEP                 = "profileDetail"
	profileBalanceEP                = "profileBalance"
	ProfileParametersEP             = "profileParameters"
	profileSubjectDNEP              = "profileSubjectDN"
	organizationsEP                 = "organizations"
	organizationDetailEP            = "organizationDetail"
	CurrentUserOrganizationEP       = "currentUserOrganization"
	checkOrganizationEP             = "checkOrganization"
	updateAuthorityEP               = "updateAuthority"
	organizationAuthoritiesEP       = "organizationAuthorities"
	AuthorityDetailEP               = "authorityDetail"
	authoritiesEP                   = "authorities"
	organizationListItemsEP         = "organizationListItems"
	organizationParametersEP        = "organizationParameters"
	ecosystemsEP                    = "ecosystems"
	userEcosystemEP                 = "userEcosystem"
	ecosystemBalanceEP              = "ecosystemBalance"
	ecosystemsStatisticsEP          = "ecosystemsStatistics"
	ecosystemAdminDetailEP          = "ecosystemAdminDetail"
	ecosystemAdminsEP               = "ecosystemAdmins"
	FindCertificateEP               = "findCertificate"
	revokeDeviceCertificateEP       = "revokeDeviceCertificate"
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

// SetBaseURL updates Sectigo to use a different scheme and host to determine endpoints
func SetBaseURL(u *url.URL) {
	baseURL = u
}

// ResetBaseURL to the default Sectigo API endpoint
func ResetBaseURL() {
	baseURL = defaultURL
}
