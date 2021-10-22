package admin

// Credentials is a generic interface that can be used to provide an access token to
// the admin API client for authentication. The credentials can be stored on disk,
// generated from private keys, or fetched via an Oauth2 workflow. The methods of the
// interface are passed a client so that the Credentials can use the API to make
// requests to the server, particularly the authenticate and reauthenticate endpoints.
type Credentials interface {
	// Login is called when the client does not have access or refresh tokens or when
	// both the access and refresh tokens are expired. The only API methods available to
	// the Credentials object are the endpoints that are unauthenticated.
	Login(api DirectoryAdministrationClient) (accessToken, refreshToken string, err error)

	// Refresh is called when the client has a valid refresh token but the access token is expired.
	Refresh(api DirectoryAdministrationClient) (accessToken, refreshToken string, err error)

	// Logout is not called explicitly by the client but is available on the client to
	// be called manually by the user. The intent is for any cached tokens to be
	// destroyed, requiring authentication before another request can be made.
	Logout(api DirectoryAdministrationClient) (err error)
}
