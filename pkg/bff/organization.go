package bff

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
	storeerrors "github.com/trisacrypto/directory/pkg/store/errors"
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
		log.Error().Err(err).Str("org_id", id).Msg("could not parse organization ID")
		return nil, err
	}

	// Fetch the record from the database
	if org, err = s.db.RetrieveOrganization(uuid); err != nil {
		return nil, err
	}

	return org, nil
}
