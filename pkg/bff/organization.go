package bff

import (
	"errors"
	"net/http"
	"strings"

	"github.com/auth0/go-auth0/management"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
	storeerrors "github.com/trisacrypto/directory/pkg/store/errors"
)

// CreateOrganization creates a new organization in the database. This endpoint returns
// an error if the organization already exists and the user is assigned to it. The user
// must have the create:organizations permission to perform this action.
func (s *Server) CreateOrganization(c *gin.Context) {
	var (
		err  error
		user *management.User
	)

	// Fetch the user from the context
	if user, err = auth.GetUserInfo(c); err != nil {
		log.Error().Err(err).Msg("create organization handler requires user info; expected middleware to return 401")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not identify user to create organization"))
		return
	}

	// Load the user app metadata to check their organization assignments
	appdata := &auth.AppMetadata{}
	if err = appdata.Load(user.AppMetadata); err != nil {
		log.Error().Err(err).Msg("could not parse user app metadata")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not parse user app metadata"))
		return
	}

	// Unmarshal the params from the POST request
	params := &api.OrganizationParams{}
	if err := c.ShouldBind(params); err != nil {
		log.Warn().Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, &api.Reply{Error: err.Error()})
		return
	}

	// Name is a required paramater
	if params.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("must provide name in request params"))
		return
	}

	// Domain is a required paramater
	domain := normalizeDomain(params.Domain)
	if domain == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("must provide domain in request params"))
		return
	}

	// Don't allow the user to create duplicate organizations
	// TODO: Should we do a universal check against the database using an index?
	for _, id := range appdata.GetOrganizations() {
		if org, err := s.OrganizationFromID(id); err != nil {
			log.Error().Err(err).Str("org_id", id).Msg("could not retrieve organization from database")
		} else if org.Domain == domain {
			c.JSON(http.StatusConflict, api.ErrorResponse("organization with domain name already exists"))
			return
		}
	}

	// Create a new organization in the database
	var org *models.Organization
	if org, err = s.db.CreateOrganization(); err != nil {
		log.Error().Err(err).Msg("could not create organization in database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create organization"))
		return
	}

	// Update the organization with the request params
	org.Name = params.Name
	org.Domain = domain
	if err = s.db.UpdateOrganization(org); err != nil {
		log.Error().Err(err).Msg("could not update organization in database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create organization"))
		return
	}

	// Assign the user to the organization
	appdata.AddOrganization(org.Id)
	if err = s.SaveAuth0AppMetadata(*user.ID, *appdata); err != nil {
		log.Error().Err(err).Msg("could not update user app metadata")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update user app metadata"))
		return
	}

	// Build the response
	out := &api.OrganizationReply{
		ID:           org.Id,
		Name:         org.Name,
		Domain:       org.Domain,
		CreatedAt:    org.Created,
		RefreshToken: true,
	}
	c.JSON(http.StatusOK, out)
}

// ListOrganizations returns a list of organizations that the user is a member of. The
// user must have the read:organizations permission to perform this action.
func (s *Server) ListOrganizations(c *gin.Context) {
	var (
		err  error
		user *management.User
	)

	// Fetch the user from the context
	if user, err = auth.GetUserInfo(c); err != nil {
		log.Error().Err(err).Msg("list organizations handler requires user info; expected middleware to return 401")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not identify user to list organizations"))
		return
	}

	// Load the user app metadata to check their organization assignments
	appdata := &auth.AppMetadata{}
	if err = appdata.Load(user.AppMetadata); err != nil {
		log.Error().Err(err).Msg("could not parse user app metadata")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not parse user app metadata"))
		return
	}

	// Build the response
	// TODO: Return a last login timestamp so the frontend can order by last used
	out := make([]*api.OrganizationReply, 0)
	for _, id := range appdata.GetOrganizations() {
		if org, err := s.OrganizationFromID(id); err != nil {
			log.Error().Err(err).Str("org_id", id).Msg("could not retrieve organization from database")
		} else {
			out = append(out, &api.OrganizationReply{
				ID:        org.Id,
				Name:      org.Name,
				Domain:    org.Domain,
				CreatedAt: org.Created,
			})
		}
	}

	c.JSON(http.StatusOK, out)
}

// Normalize a domain name for matching purposes.
func normalizeDomain(domain string) string {
	return strings.ToLower(strings.TrimRight(strings.TrimSpace(domain), "."))
}

// OrganizationFromClaims is a helper method to retrieve the organization for a
// particular request by fetching the orgID from the claims and querying the database.
// If there is an error fetching the organization, the appropriate error response is
// made on the gin writer and logged. The caller should check for error and return.
func (s *Server) OrganizationFromClaims(c *gin.Context) (org *models.Organization, err error) {
	// Retrieve the organization ID from the claims
	var claims *auth.Claims
	if claims, err = auth.GetClaims(c); err != nil {
		log.Error().Err(err).Msg("could not retrieve claims to fetch orgID")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not identify organization"))
		return nil, err
	}

	// If there is no organization ID, something went wrong
	if claims.OrgID == "" {
		log.Warn().Msg("claims do not contain an orgID")
		api.MustRefreshToken(c, "missing claims info, try logging out and logging back in")
		return nil, errors.New("missing organization ID in claims")
	}

	// Fetch the organization from the database
	if org, err = s.OrganizationFromID(claims.OrgID); err != nil {
		if errors.Is(err, storeerrors.ErrEntityNotFound) {
			log.Warn().Err(err).Msg("could not find organization in database from orgID in claims")
			api.MustRefreshToken(c, "no organization found, try logging out and logging back in")
			return nil, err
		}

		log.Error().Err(err).Msg("could not retrieve organization")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not identify organization"))
		return nil, err
	}

	return org, nil
}

func (s *Server) OrganizationFromID(id string) (org *models.Organization, err error) {
	// Organizations are stored by UUID in the database
	var uuid uuid.UUID
	if uuid, err = models.ParseOrgID(id); err != nil {
		return nil, err
	}

	// Fetch the record from the database
	if org, err = s.db.RetrieveOrganization(uuid); err != nil {
		return nil, err
	}

	return org, nil
}
