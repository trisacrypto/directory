package bff

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/auth0/go-auth0/management"
	"github.com/gin-gonic/gin"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/utils"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
)

const orgIDParam = "orgid"

// AddCollaborator creates a new collaborator with the email address in the request.
// The endpoint adds the collaborator to the organization record associated with the
// user and sends a verification email to the provided email address.
//
// @Summary Add collaborator [update:collaborators]
// @Description Invite a new collaborator to the user's organization.
// @Tags collaborators
// @Accept json
// @Produce json
// @Param collaborator body models.Collaborator true "Collaborator to add"
// @Success 200 {object} models.Collaborator
// @Failure 400 {object} api.Reply "Invalid collaborator, email address is required"
// @Failure 401 {object} api.Reply
// @Failure 403 {object} api.Reply "Maximum number of collaborators reached"
// @Failure 409 {object} api.Reply "Collaborator already exists"
// @Failure 500 {object} api.Reply
// @Router /collaborators [post]
func (s *Server) AddCollaborator(c *gin.Context) {
	var (
		err          error
		collaborator *models.Collaborator
		inviter      *management.User
		org          *models.Organization
	)

	// Fetch the organization from the claims
	// NOTE: This method handles the error logging and response
	if org, err = s.OrganizationFromClaims(c); err != nil {
		return
	}

	// The invoking user is the inviter
	if inviter, err = auth.GetUserInfo(c); err != nil {
		sentry.Error(c).Err(err).Msg("add collaborator handler requires user info; expected middleware to return 401")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not identify user"))
		return
	}

	// Unmarshal the collaborator from the POST request
	collaborator = &models.Collaborator{}
	if err = c.ShouldBind(collaborator); err != nil {
		sentry.Warn(c).Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Collaborator invites are valid for seven days
	collaborator.ExpiresAt = time.Now().AddDate(0, 0, 7).Format(time.RFC3339Nano)

	// Add the collaborator to the organization
	if err = org.AddCollaborator(collaborator); err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidCollaborator):
			sentry.Warn(c).Err(err).Msg("invalid collaborator")
			c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		case errors.Is(err, models.ErrCollaboratorExists):
			sentry.Warn(c).Err(err).Str("email", collaborator.Email).Msg("collaborator already exists")
			c.JSON(http.StatusConflict, api.ErrorResponse(err))
		case errors.Is(err, models.ErrMaxCollaborators):
			sentry.Warn(c).Err(err).Int("maximum", models.MaxCollaborators).Msg("maximum number of collaborators reached")
			c.JSON(http.StatusForbidden, api.ErrorResponse(err))
		default:
			sentry.Error(c).Err(err).Msg("could not add collaborator")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		}
		return
	}

	// Handle both users who have already registered and completely new users
	var user *management.User
	var baseURL string
	if user, err = s.FindUserByEmail(c.Request.Context(), collaborator.Email); err != nil {
		if !errors.Is(err, ErrUserEmailNotFound) {
			sentry.Error(c).Err(err).Str("email", collaborator.Email).Msg("error finding user by email in Auth0")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not add collaborator"))
			return
		}

		// If the user doesn't exist in Auth0 then direct them to the signup page
		baseURL = s.conf.RegisterURL
		user = &management.User{
			Email: &collaborator.Email,
		}
	} else {
		// If the user already exists then direct them to the login page
		baseURL = s.conf.LoginURL
	}

	// Include the organization ID in the invite URL
	var inviteURL *url.URL
	if inviteURL, err = url.Parse(baseURL); err != nil {
		sentry.Error(c).Err(err).Str("url", baseURL).Msg("could not parse invite URL")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not add collaborator"))
		return
	}
	q := inviteURL.Query()
	q.Set(orgIDParam, org.Id)
	inviteURL.RawQuery = q.Encode()

	// Send the invite email to the user
	if err = s.email.SendUserInvite(user, inviter, org, inviteURL); err != nil {
		sentry.Error(c).Err(err).Str("email", *user.Email).Msg("error sending user invite email")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not add collaborator"))
		return
	}

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Save the updated organization
	if err = s.db.UpdateOrganization(ctx, org); err != nil {
		sentry.Error(c).Err(err).Msg("could not save organization with new collaborator")
		c.JSON(http.StatusInternalServerError, "could not add collaborator")
		return
	}

	c.JSON(http.StatusOK, collaborator)
}

// ListCollaborators lists all the collaborators on the user's organization.
//
// @Summary List collaborators [read:collaborators]
// @Description Returns all collaborators on the user's organization sorted by email address.
// @Tags collaborators
// @Produce json
// @Success 200 {object} api.ListCollaboratorsReply
// @Failure 401 {object} api.Reply
// @Failure 500 {object} api.Reply
// @Router /collaborators [get]
func (s *Server) ListCollaborators(c *gin.Context) {
	var (
		err error
		org *models.Organization
	)

	// Fetch the organization from the claims
	// NOTE: This method handles the error logging and response
	if org, err = s.OrganizationFromClaims(c); err != nil {
		return
	}

	// Build the response from the internal map
	out := &api.ListCollaboratorsReply{
		Collaborators: make([]*models.Collaborator, 0),
	}

	for _, collab := range org.Collaborators {
		if err = s.LoadCollaboratorDetails(c.Request.Context(), collab); err != nil {
			sentry.Error(c).Err(err).Str("collabID", collab.Key()).Msg("could not load collaborator details")
		}

		// Enforce consistent ordering by email address
		out.Collaborators = InsortCollaborator(out.Collaborators, collab, func(a, b *models.Collaborator) bool {
			return a.Email < b.Email
		})
	}

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Collaborators exist on the organization record so we must persist the updated
	// organization record to the database
	if err = s.db.UpdateOrganization(ctx, org); err != nil {
		sentry.Error(c).Err(err).Msg("could not save organization with updated collaborators")
	}

	c.JSON(http.StatusOK, out)
}

// UpdateCollaboratorRoles updates the roles of the collaborator ID in the request,
// ensuring that the roles are updated both on the organization record and in Auth0.
// The user must have the update:collaborators permission to make this request.
//
// @Summary Update collaborator roles [update:collaborators]
// @Description Replace the roles of the collaborator with the given ID.
// @Tags collaborators
// @Accept json
// @Produce json
// @Param id path string true "Collaborator ID"
// @Param roles body api.UpdateRolesParams true "New roles for the collaborator"
// @Success 200 {object} api.UpdateRolesParams
// @Failure 400 {object} api.Reply
// @Failure 401 {object} api.Reply
// @Failure 404 {object} api.Reply
// @Failure 500 {object} api.Reply
// @Router /collaborators/{id} [post]
func (s *Server) UpdateCollaboratorRoles(c *gin.Context) {
	var (
		err          error
		collaborator *models.Collaborator
		org          *models.Organization
	)

	// Fetch the organization from the claims
	// NOTE: This method handles the error logging and response
	if org, err = s.OrganizationFromClaims(c); err != nil {
		return
	}

	// Get the collabID from the URL
	collabID := c.Param("collabID")

	// Unmarshal the roles from the POST request
	params := &api.UpdateRolesParams{}
	if err = c.ShouldBind(params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Make sure a 404 is returned if the collaborator does not exist
	var ok bool
	if collaborator, ok = org.Collaborators[collabID]; !ok {
		sentry.Warn(c).Str("collabID", collabID).Msg("collaborator does not exist")
		c.JSON(http.StatusNotFound, api.ErrorResponse("collaborator does not exist"))
		return
	}

	// Collaborator needs to be a verified user in Auth0
	if !collaborator.Verified || collaborator.UserId == "" {
		sentry.Warn(c).Str("collabID", collabID).Msg("cannot update roles for unverified collaborator")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("cannot update roles for unverified collaborator"))
		return
	}

	// Update the users's roles in Auth0
	if err = s.AssignRoles(c.Request.Context(), collaborator.UserId, params.Roles); err != nil {
		if errors.Is(err, ErrInvalidUserRole) {
			c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		} else {
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update collaborator roles"))
		}
		return
	}

	// Update the collaborator record with the info in Auth0
	if err = s.LoadCollaboratorDetails(c.Request.Context(), collaborator); err != nil {
		sentry.Error(c).Err(err).Str("collabID", collabID).Str("auth0_id", collaborator.UserId).Msg("could not update collaborator record")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update collaborator record"))
		return
	}

	// Make sure the collaborator record has updated timestamps
	collaborator.ModifiedAt = time.Now().UTC().Format(time.RFC3339Nano)
	if collaborator.CreatedAt == "" {
		collaborator.CreatedAt = collaborator.ModifiedAt
	}

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Save the updated organization
	if err = s.db.UpdateOrganization(ctx, org); err != nil {
		sentry.Error(c).Err(err).Msg("could not save organization with new collaborator")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not replace collaborator"))
		return
	}

	c.JSON(http.StatusOK, collaborator)
}

// DeleteCollaborator deletes the collaborator in the request from the user's
// organization. The user must have the update:collaborators permission.
// Note: This does not return an error if the collaborator does not exist on the
// organization and instead returns a 200 OK response.
//
// @Summary Delete collaborator [update:collaborators]
// @Description Delete the collaborator with the given ID from the organization.
// @Tags collaborators
// @Produce json
// @Param id path string true "Collaborator ID"
// @Success 200 {object} api.Reply
// @Failure 401 {object} api.Reply
// @Failure 404 {object} api.Reply
// @Failure 500 {object} api.Reply
// @Router /collaborators/{id} [delete]
func (s *Server) DeleteCollaborator(c *gin.Context) {
	var (
		err          error
		collaborator *models.Collaborator
		org          *models.Organization
	)

	// Fetch the organization from the claims
	// NOTE: This method handles the error logging and response
	if org, err = s.OrganizationFromClaims(c); err != nil {
		return
	}

	// Get the collabID from the URL
	collabID := c.Param("collabID")

	// Make sure a 404 is returned if the collaborator does not exist
	var ok bool
	if collaborator, ok = org.Collaborators[collabID]; !ok {
		sentry.Warn(c).Str("collabID", collabID).Msg("collaborator does not exist")
		c.JSON(http.StatusNotFound, api.ErrorResponse("collaborator not found"))
		return
	}

	// If the collaborator is already verified in Auth0, then remove them from the
	// organization
	if collaborator.Verified && collaborator.UserId != "" {
		// Fetch the user from Auth0
		var user *management.User
		if user, err = s.auth0.User.Read(c.Request.Context(), collaborator.UserId); err != nil {
			sentry.Error(c).Err(err).Str("collabID", collabID).Str("auth0_id", collaborator.UserId).Msg("could not fetch user from Auth0")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch user from Auth0"))
			return
		}

		// Fetch the user app metadata
		appdata := &auth.AppMetadata{}
		if err = appdata.Load(user.AppMetadata); err != nil {
			sentry.Error(c).Err(err).Str("collabID", collabID).Str("auth0_id", collaborator.UserId).Msg("could not parse user app metadata")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not parse user app metadata"))
			return
		}

		// Update the app metadata with the removed organization
		appdata.ClearOrganization()
		if err = s.SaveAuth0AppMetadata(c.Request.Context(), *user.ID, *appdata); err != nil {
			sentry.Error(c).Err(err).Str("collabID", collabID).Str("auth0_id", collaborator.UserId).Msg("could not save user app metadata")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not save user app metadata"))
			return
		}
	}

	// Delete the collaborator from the organization record
	delete(org.Collaborators, collabID)

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Save the updated organization
	if err = s.db.UpdateOrganization(ctx, org); err != nil {
		sentry.Error(c).Err(err).Msg("could not save organization without collaborator")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not delete collaborator"))
		return
	}

	c.Status(http.StatusOK)
}

// LoadCollaboratorDetails updates a collaborator record with the user details in
// Auth0. The collaborator must have a user ID on it and the data in Auth0 will
// overwrite the data on the collaborator record.
func (s *Server) LoadCollaboratorDetails(ctx context.Context, collab *models.Collaborator) (err error) {
	// If the user is not verified in Auth0 then we can't retrieve the details
	if !collab.Verified {
		return nil
	}

	// Record must have a user ID to query Auth0
	if collab.UserId == "" {
		return errors.New("collaborator does not have a user ID")
	}

	// Fetch the user profile from Auth0
	var profile *auth.UserProfile
	if profile, err = s.FetchUserProfile(ctx, collab.UserId); err != nil {
		return err
	}

	// Refresh the collaborator record with the details
	collab.Name = profile.Name
	collab.Roles = profile.Roles

	return nil
}

// FetchUserProfile fetches a user profile by ID from the user cache or Auth0 if
// necessary.
func (s *Server) FetchUserProfile(ctx context.Context, id string) (profile *auth.UserProfile, err error) {
	if data, ok := s.users.Get(id); !ok {
		// If we can't get the profile from the cache then fetch it from Auth0
		var user *management.User
		if user, err = s.auth0.User.Read(ctx, id); err != nil {
			return nil, err
		}

		profile = &auth.UserProfile{}
		if user.Name != nil {
			profile.Name = *user.Name
		}

		var roles *management.RoleList
		if roles, err = s.auth0.User.Roles(ctx, id); err != nil {
			return nil, err
		}

		profile.Roles = make([]string, len(roles.Roles))
		for i, role := range roles.Roles {
			profile.Roles[i] = *role.Name
		}

		// Add the fetched profile to the cache
		s.users.Add(id, profile)
	} else if profile, ok = data.(*auth.UserProfile); !ok {
		return nil, errors.New("invalid user profile cache entry")
	}

	return profile, nil
}

// InsortCollaborator is a helper function to insert a collaborator into a sorted slice
// using a custom sort function.
func InsortCollaborator(collabs []*models.Collaborator, value *models.Collaborator, f func(a, b *models.Collaborator) bool) []*models.Collaborator {
	if collabs == nil || value == nil || f == nil {
		return nil
	}

	i := sort.Search(len(collabs), func(i int) bool { return f(value, collabs[i]) })
	collabs = append(collabs, nil)
	copy(collabs[i+1:], collabs[i:])
	collabs[i] = value
	return collabs
}
