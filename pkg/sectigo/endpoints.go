package sectigo

import (
	"fmt"
	"net/url"
)

func init() {
	buildEndpoints()
}

const (
	contentType         = "application/json;charset=UTF-8"
	downloadContentType = "application/octet-stream"
	userAgent           = "TRISADS Sectigo Client v1.0"
)

// baseURL is used to construct API endpoints for Sectigo methods.
// var baseURL = &url.URL{Scheme: "http", Host: "localhost:8812"}
var baseURL = &url.URL{Scheme: "https", Host: "iot.sectigo.com"}

// endpoints maps methods to URLs, which are full constructed with the baseURL in the
// package init function, which calls buildEndpoints(). Some endpoint paths contain
// string format verbs intended for dynamic REST urls by making a copy of the url.URL
// and replacing the Path with the output of fmt.Sprintf(endpoints[method].Path, param).
// If the format verb is not replaced in this way, it will be encoded as %25<verb> which
// will likely return a 404 error. Note that value copying a url.URL is not a deep copy
// and the internal user info will not be copied.
//
// Many endpoints require URL parameters. See the API documentation for more details.
// Endpoint URLs specify paths only, copy the URL and add URL parameters for dynamic
// endpoints.
var endpoints = map[string]*url.URL{
	authenticateEP:                  {Path: "/auth/pwd"},
	refreshEP:                       {Path: "/auth/refresh"},
	batchesEP:                       {Path: "/api/v1/batches"},
	batchDetailEP:                   {Path: "/api/v1/batches/%d"},
	batchStatusEP:                   {Path: "/api/v1/batches/%d/status"},
	batchAuditLogEP:                 {Path: "/api/v1/batches/%d/auditLog"},
	batchProcessingInfoEP:           {Path: "/api/v1/batches/%d/processing_info"},
	batchDevicesAuditLogEP:          {Path: "/api/v1/batches/%d/devices/auditLog"},
	batchPreviewEP:                  {Path: "/api/v1/batches/preview"},
	createSingleCertBatchEP:         {Path: "/api/v1/batches/createSingleCertBatch"},
	uploadEP:                        {Path: "/api/v1/upload"},
	uploadCSVEP:                     {Path: "/api/v2/organizations/%d/profiles/%d/batches/csv-upload"},
	generatorsEP:                    {Path: "/api/v1/generators"},
	downloadEP:                      {Path: "/api/v1/batches/%d/download"},
	devicesEP:                       {Path: "/api/v1/devices"},
	userAuthoritiesEP:               {Path: "/api/v1/authorities/allowed"},
	authorityBalanceUsedEP:          {Path: "/api/v1/authorities/%d/balanceused/%d"},
	authorityBalanceAvailableEP:     {Path: "/api/v1/authorities/%d/balanceavailable/%d"},
	authorityUserBalanceAvailableEP: {Path: "/api/v1/authorities/%d/balanceavailable"},
	usersEP:                         {Path: "/api/v1/users"},
	checkEmailEP:                    {Path: "/api/v1/users/check_email"},
	userDetailEP:                    {Path: "/api/v1/users/%d"},
	userProfilesEP:                  {Path: "/api/v1/users/profiles"},
	usersOrganizationsEP:            {Path: "/api/v1/users/organizations/%d"},
	removeUserEP:                    {Path: "/api/v1/users/remove_user"},
	currentUserEP:                   {Path: "/api/v1/users/current"},
	updateEmailEP:                   {Path: "/api/v1/users/update_email"},
	checkUserEP:                     {Path: "/api/v1/users/check_user"},
	userCredentialsEP:               {Path: "/api/v1/users/credentials"},
	profilesEP:                      {Path: "/api/v1/profiles"},
	profileDetailEP:                 {Path: "/api/v1/profiles/%d"},
	profileBalanceEP:                {Path: "/api/v1/profiles/%d/balance"},
	profileParametersEP:             {Path: "/api/v1/profiles/%d/parameters"},
	profileSubjectDNEP:              {Path: "/api/v1/profiles/%d/dn"},
	organizationsEP:                 {Path: "/api/v1/organizations"},
	organizationDetailEP:            {Path: "/api/v1/organizations/%d"},
	currentUserOrganizationEP:       {Path: "/api/v1/organizations/user"},
	checkOrganizationEP:             {Path: "/api/v1/organizations/check_organization"},
	updateAuthorityEP:               {Path: "/api/v1/organizations/%d/authorities/%d"},
	organizationAuthoritiesEP:       {Path: "/api/v1/organizations/%d/all_authorities"},
	authorityDetailEP:               {Path: "/api/v1/organizations/%d/authority/%d"},
	authoritiesEP:                   {Path: "/api/v1/organizations/%d/authority"},
	organizationListItemsEP:         {Path: "/api/v1/organizations/select/items"},
	organizationParametersEP:        {Path: "/api/v1/organizations/profile_parameters/%d"},
	ecosystemsEP:                    {Path: "/api/v1/ecosystems"},
	userEcosystemEP:                 {Path: "/api/v1/ecosystems/ecosystem"},
	ecosystemBalanceEP:              {Path: "/api/v1/ecosystems/ecosystem/balance"},
	ecosystemsStatisticsEP:          {Path: "/api/v1/ecosystems/statistics"},
	ecosystemAdminDetailEP:          {Path: "/api/v1/ecosystems/users/%d"},
	ecosystemAdminsEP:               {Path: "/api/v1/ecosystems/users"},
	findCertificateEP:               {Path: "/api/v1/certificates/find"},
	revokeDeviceCertificateEP:       {Path: "/api/v1/certificates/revoke"},
	revokeCertificateEP:             {Path: "/api/v1/certificates/%d/revoke"},
}

// endpoint name constants to prevent typos at compile time rather than at runtime.
const (
	authenticateEP                  = "authenticate"
	refreshEP                       = "refresh"
	batchesEP                       = "batches"
	batchDetailEP                   = "batchDetail"
	batchStatusEP                   = "batchStatus"
	batchAuditLogEP                 = "batchAuditLog"
	batchProcessingInfoEP           = "batchProcessingInfo"
	batchDevicesAuditLogEP          = "batchDevicesAuditLog"
	batchPreviewEP                  = "batchPreview"
	createSingleCertBatchEP         = "createSingleCertBatch"
	uploadEP                        = "upload"
	uploadCSVEP                     = "uploadCSV"
	generatorsEP                    = "generators"
	downloadEP                      = "download"
	devicesEP                       = "devices"
	userAuthoritiesEP               = "userAuthorities"
	authorityBalanceUsedEP          = "authorityBalanceUsed"
	authorityBalanceAvailableEP     = "authorityBalanceAvailable"
	authorityUserBalanceAvailableEP = "authorityUserBalanceAvailable"
	usersEP                         = "users"
	checkEmailEP                    = "checkEmail"
	userDetailEP                    = "userDetail"
	userProfilesEP                  = "userProfiles"
	usersOrganizationsEP            = "usersOrganizations"
	removeUserEP                    = "removeUser"
	currentUserEP                   = "currentUser"
	updateEmailEP                   = "updateEmail"
	checkUserEP                     = "checkUser"
	userCredentialsEP               = "userCredentials"
	profilesEP                      = "profiles"
	profileDetailEP                 = "profileDetail"
	profileBalanceEP                = "profileBalance"
	profileParametersEP             = "profileParameters"
	profileSubjectDNEP              = "profileSubjectDN"
	organizationsEP                 = "organizations"
	organizationDetailEP            = "organizationDetail"
	currentUserOrganizationEP       = "currentUserOrganization"
	checkOrganizationEP             = "checkOrganization"
	updateAuthorityEP               = "updateAuthority"
	organizationAuthoritiesEP       = "organizationAuthorities"
	authorityDetailEP               = "authorityDetail"
	authoritiesEP                   = "authorities"
	organizationListItemsEP         = "organizationListItems"
	organizationParametersEP        = "organizationParameters"
	ecosystemsEP                    = "ecosystems"
	userEcosystemEP                 = "userEcosystem"
	ecosystemBalanceEP              = "ecosystemBalance"
	ecosystemsStatisticsEP          = "ecosystemsStatistics"
	ecosystemAdminDetailEP          = "ecosystemAdminDetail"
	ecosystemAdminsEP               = "ecosystemAdmins"
	findCertificateEP               = "findCertificate"
	revokeDeviceCertificateEP       = "revokeDeviceCertificate"
	revokeCertificateEP             = "revokeCertificate"
)

// Convert the endpoints into absolute URLs by resolving them with the base URL.
func buildEndpoints() {
	for key, endpoint := range endpoints {
		endpoints[key] = baseURL.ResolveReference(endpoint)
	}
}

// Get a URL for the specified endpoint with the given parameters.
func urlFor(endpoint string, params ...interface{}) string {
	u, ok := endpoints[endpoint]
	if !ok {
		// this is a developer error, so panic
		panic(fmt.Sprintf("no endpoint named %q", endpoint))
	}

	// Copy the URL so the original URL isn't modified
	if len(params) > 0 {
		v := *u
		v.Path = fmt.Sprintf(u.Path, params...)
		return v.String()
	}
	return u.String()
}
