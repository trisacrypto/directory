package bff

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/db"
	"github.com/trisacrypto/directory/pkg/bff/db/models/v1"
)

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

	// Fetch the record from the database
	if org, err = s.db.Organizations().Retrieve(c.Request.Context(), claims.OrgID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
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

// AddOrganizationCollaborator adds a new collaborator record to an organization.
func AddOrganizationCollaborator(org *models.Organization, email string) (err error) {
	if email == "" {
		return errors.New("email address is required to add an organization collaborator")
	}

	if org.Collaborators == nil {
		org.Collaborators = make(map[string]*models.Collaborator)
	}

	// Don't overwrite an existing collaborator
	if _, ok := org.Collaborators[email]; ok {
		return fmt.Errorf("collaborator with email address %s already exists", email)
	}

	// Create a new collaborator record
	org.Collaborators[email] = &models.Collaborator{
		Email:     email,
		CreatedAt: time.Now().Format(time.RFC3339Nano),
	}
	return nil
}

// GetOrganizationCollaborator returns the collaborator record for the given email address.
func GetOrganizationCollaborator(org *models.Organization, email string) (collaborator *models.Collaborator, err error) {
	if email == "" {
		return nil, errors.New("email address is required to get an organization collaborator")
	}

	if org.Collaborators == nil {
		org.Collaborators = make(map[string]*models.Collaborator)
	}

	// Lookup the collaborator record
	var ok bool
	if collaborator, ok = org.Collaborators[email]; !ok {
		return nil, fmt.Errorf("collaborator with email address %s does not exist", email)
	}

	return collaborator, nil
}
