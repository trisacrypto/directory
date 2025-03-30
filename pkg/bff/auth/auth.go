package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/auth0/go-auth0/management"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth/authtest"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
)

const (
	ScopeAnonymous          = "anonymous"
	ContextUserInfo         = "auth0_user_info"
	ContextBFFClaims        = "auth0_bff_claims"
	ContextRegisteredClaims = "auth0_registered_claims"
)

// AnonymousClaims are used to identify unauthenticated requests that have no permissions.
var AnonymousClaims = Claims{Scope: ScopeAnonymous, Permissions: nil}

// Claims extracts custom data from the JWT token provided by Auth0
type Claims struct {
	Scope         string   `json:"scope"`
	Permissions   []string `json:"permissions"`
	OrgID         string   `json:"https://trisa.directory/orgid"`
	VASPs         VASPs    `json:"https://trisa.directory/vasps"`
	Organizations []string `json:"https://trisa.directory/organizations"`
	Email         string   `json:"https://trisa.directory/email"`
}

// Validate implements the validator.CustomClaims interface for Auth0 parsing.
// Claims can have empty scope (e.g. no permissions) and no associated VASP.
func (c Claims) Validate(ctx context.Context) error {
	return nil
}

// HasScope checks if the claims contain the specified scope.
func (c Claims) HasScope(requiredScope string) bool {
	scopes := strings.Split(c.Scope, " ")
	for _, scope := range scopes {
		if scope == requiredScope {
			return true
		}
	}
	return false
}

// HasPermission checks if the claims contain the specified permission.
func (c Claims) HasPermission(requiredPermission string) bool {
	for _, permission := range c.Permissions {
		if permission == requiredPermission {
			return true
		}
	}
	return false
}

// HasAllPermissions checks if all specified permissions are in the claims.
func (c Claims) HasAllPermissions(requiredPermissions ...string) bool {
	for _, requiredPermission := range requiredPermissions {
		if !c.HasPermission(requiredPermission) {
			return false
		}
	}
	return true
}

// IsAnonymous returns true if the claims refer to an anonymous user
func (c Claims) IsAnonymous() bool {
	return c.HasScope(ScopeAnonymous)
}

// NewManagementClient creates a new Auth0 management client from the configuration.
func NewManagementClient(conf config.AuthConfig) (manager *management.Management, err error) {
	var options []management.Option
	domain := conf.Domain
	if conf.Testing {
		var ts *authtest.Server
		if ts, err = authtest.Serve(); err != nil {
			return nil, err
		}
		options = append(options, management.WithClient(ts.Client()))
		options = append(options, management.WithStaticToken(conf.ClientSecret))
		domain = ts.URL.Host
	} else {
		options = append(options, conf.ClientCredentials())
	}

	// Connect to the management API
	if manager, err = management.New(domain, options...); err != nil {
		return nil, err
	}
	return manager, nil
}

// NewClaims implements the validator custom claims initializer interface.
func NewClaims() validator.CustomClaims {
	return &Claims{}
}

// Authenticate is a middleware that will parse and validate any Auth0 token provided
// in the header of the request and will add the claims to the request context for
// downstream processing. If no JWT token is present in the header, this middleware will
// mark the request as unauthenticated but it does not perform any authorization. If the
// JWT token is invalid this middleware will return a 403 Forbidden response.
func Authenticate(conf config.AuthConfig, options ...jwks.ProviderOption) (_ gin.HandlerFunc, err error) {
	// Parse the issuer url to ensure it is correctly configured.
	var issuerURL *url.URL
	if issuerURL, err = conf.IssuerURL(); err != nil {
		return nil, err
	}

	// If we're in testing mode and no other options have been provided, connect to the
	// default authtest server to validate local, test credentials
	if conf.Testing && len(options) == 0 {
		var ts *authtest.Server
		if ts, err = authtest.Serve(); err != nil {
			return nil, err
		}
		options = append(options, WithHTTPClient(ts.Client()))
	}

	// The caching provider fetches the JWKS (JSON Web Key Set) public keys used to
	// validate JWT signatures to prove that they were issued by auth0. The JWKS are
	// cached for the configured TTL (default 5 minutes) before being refetched.
	provider := jwks.NewCachingProvider(issuerURL, conf.ProviderCache, options...)

	// Create the JWT validator from the configuration. The validator parses the JWT
	// token, confirms it is not expired, configured for the correct audience, and has
	// been signed by auth0 -- this is the workhorse of the authentication middleware.
	var auth0 *validator.Validator
	if auth0, err = validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{conf.Audience},
		validator.WithCustomClaims(NewClaims),
		validator.WithAllowedClockSkew(500*time.Millisecond),
	); err != nil {
		return nil, fmt.Errorf("could not set up JWT validator: %s", err)
	}

	return func(c *gin.Context) {
		var (
			err    error
			tks    string
			claims interface{}
		)

		if tks, err = jwtmiddleware.AuthHeaderTokenExtractor(c.Request); err != nil || tks == "" {
			// The most common reason there is no token in the header is because it is
			// not provided -- add an unauthenticated, anonymous user to the context.
			// The second most common case is that the token is malformed or incorrectly
			// structured. If this is a recurring problem, we will have to add extra
			// checks to determine if an authorization header was provided or not.
			if err != nil {
				sentry.Warn(c).Err(err).Msg("could not extract token from authorization header")
			}

			// Add anonymous user and empty claims to context
			log.Debug().Msg("anonymous user")
			c.Set(ContextBFFClaims, &AnonymousClaims)
		} else {
			// If a token is provided in the authorization header, verify that it was
			// correctly signed using auth0's public keys and add user claims to the
			// context. If the token is not valid return a forbidden error.
			if claims, err = auth0.ValidateToken(c.Request.Context(), tks); err != nil {
				sentry.Warn(c).Err(err).Msg("invalid authorization token")
				c.AbortWithStatusJSON(http.StatusForbidden, api.ErrorResponse(ErrInvalidAuthToken))
				return
			}

			// Set the claims on the gin context for downstream processing
			// NOTE: invalid type assertions will cause panics which will be recovered
			claims := claims.(*validator.ValidatedClaims)
			c.Set(ContextBFFClaims, claims.CustomClaims.(*Claims))
			c.Set(ContextRegisteredClaims, &claims.RegisteredClaims)
		}

		// Continue handling the request with next middleware.
		c.Next()
	}, nil
}

// Authorize is a middleware that requires specific permissions in an authenticated
// user's claims. If those permissions do not match or the request is unauthenticated
// the middleware returns a 401 Unauthorized response. The Authorize middleware must
// follow the Authenticate middleware.
func Authorize(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := GetClaims(c)
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(ErrNoAuthorization))
			return
		}

		if claims.IsAnonymous() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(ErrAuthRequired))
			return
		}

		if !claims.HasAllPermissions(permissions...) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(ErrNoPermission))
			return
		}

		c.Next()
	}
}

// UserInfo is a middleware that requires an authenticated user's claims, it then
// fetches the user profile including app_data from Auth0 and adds them to the Gin
// context. This middleware is primarily used for endpoints that manage the user state,
// not for endpoints that simply need access to resources or permissions (those should
// be added to the claims to prevent calls to Auth0 on every RPC). If the user is not
// authenticated before this step, a 401 is returned.
func UserInfo(conf config.AuthConfig) (_ gin.HandlerFunc, err error) {
	var manager *management.Management
	if manager, err = NewManagementClient(conf); err != nil {
		return nil, err
	}

	return func(c *gin.Context) {
		claims, err := GetRegisteredClaims(c)
		if err != nil || claims.Subject == "" {
			// c.Error will panic if err is nil so wrap in a guard
			if err != nil {
				c.Error(err)
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(ErrNoAuthUser))
			return
		}

		user, err := manager.User.Read(c.Request.Context(), claims.Subject)
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadGateway, api.ErrorResponse(ErrNoAuthUserData))
			return
		}

		// User must have complete identity data for downstream processing
		if user.ID == nil || *user.ID == "" || user.Email == nil || *user.Email == "" {
			c.Error(ErrIncompleteUser)
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(ErrIncompleteUser))
			return
		}

		// User must be email verified to access the API
		if user.EmailVerified == nil || !*user.EmailVerified {
			c.Error(fmt.Errorf("user %s is not email verified", *user.ID))
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(ErrUnverifiedUser))
			return
		}

		c.Set(ContextUserInfo, user)
		c.Next()
	}, nil
}

// GetClaims fetches and parses the BFF claims from the gin context. Returns an error if
// no claims exist on the context rather than returning anonymous claims. Panics if the
// claims are an incorrect type, but the panic should be recovered by middleware.
func GetClaims(c *gin.Context) (*Claims, error) {
	claims, exists := c.Get(ContextBFFClaims)
	if !exists {
		return nil, ErrNoClaims
	}
	return claims.(*Claims), nil
}

// GetRegisteredClaims fetches and parses the access token claims from the gin context.
// Returns an error if no claims exist on the context rather than returning zero-valued
// claims. Panics if the claims are an incorrect type, but should be recovered.
func GetRegisteredClaims(c *gin.Context) (*validator.RegisteredClaims, error) {
	claims, exists := c.Get(ContextRegisteredClaims)
	if !exists {
		return nil, ErrNoClaims
	}
	rclaims := claims.(*validator.RegisteredClaims)
	return rclaims, nil
}

// GetUserInfo fetches the user info from the gin context. Returns an error if no user
// exists on the context or if the user value is nil. Panics if user is incorrect type.
func GetUserInfo(c *gin.Context) (*management.User, error) {
	user, exists := c.Get(ContextUserInfo)
	if !exists || user == nil {
		return nil, ErrNoUserInfo
	}
	return user.(*management.User), nil
}
