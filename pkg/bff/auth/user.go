package auth

import (
	"errors"

	"github.com/auth0/go-auth0/management"
	"github.com/trisacrypto/directory/pkg/bff/auth/cache"
	"github.com/trisacrypto/directory/pkg/bff/config"
)

type UserFetcher struct {
	client *management.UserManager
}

func NewUserFetcher(client *management.UserManager) *UserFetcher {
	return &UserFetcher{
		client: client,
	}
}

// UserDetails contains information about the user that can be safely cached on the BFF
// server.
type UserDetails struct {
	Name  string
	Roles []string
}

// UserFetcher implements the ResourceFetcher interface
var _ cache.ResourceFetcher = &UserFetcher{}

func (f *UserFetcher) Get(id string) (data interface{}, err error) {
	var user *management.User
	if user, err = f.client.User.Read(id); err != nil {
		return nil, err
	}

	profile := &UserDetails{}
	if user.Name != nil {
		profile.Name = *user.Name
	}

	var roles *management.RoleList
	if roles, err = f.client.User.Roles(id); err != nil {
		return nil, err
	}

	profile.Roles = make([]string, len(roles.Roles))
	for i, role := range roles.Roles {
		profile.Roles[i] = *role.Name
	}

	return profile, nil
}

// UserCache caches user details on the BFF server to reduce the number of calls to
// Auth0.
type UserCache struct {
	cache *cache.TTLCache
}

func NewUserCache(conf config.CacheConfig, client *management.UserManager) *UserCache {
	ttl := cache.NewTTLCache(conf)
	ttl.SetFetcher(NewUserFetcher(client))
	ttl.SetStore(cache.NewTTLStore(conf.TTLMean, conf.TTLSigma))
	return &UserCache{
		cache: ttl,
	}
}

// Get is a helper to fetch user data from a generic cache and assert it to the
// expected type.
func (c *UserCache) Get(id string) (details *UserDetails, err error) {
	var data interface{}
	if data, err = c.cache.Get(id); err != nil {
		return nil, err
	}

	if details, ok := data.(*UserDetails); ok {
		return details, nil
	}

	return nil, errors.New("could not assert user data to the UserDetails type")
}

// UserDisplayName is a helper to get the user's display name from the Auth0 user
// record. This should be used when the backend needs to retrieve a user-facing display
// name for the user and returns an error if no name is available.
func UserDisplayName(user *management.User) (string, error) {
	if user == nil {
		return "", errors.New("user record is nil")
	}

	// Prefer the user's actual name if available
	switch {
	case user.Name != nil && *user.Name != "":
		return *user.Name, nil
	case user.Email != nil && *user.Email != "":
		return *user.Email, nil
	default:
		return "", errors.New("user record has no name or email address")
	}
}
