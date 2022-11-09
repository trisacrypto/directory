package emails

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"

	"github.com/auth0/go-auth0/management"
	"github.com/rs/zerolog/log"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/utils/emails"
	"github.com/trisacrypto/directory/pkg/utils/emails/mock"
)

func New(conf config.EmailConfig) (m *EmailManager, err error) {
	m = &EmailManager{
		conf: conf,
	}
	if conf.Testing {
		log.Warn().Bool("testing", conf.Testing).Str("storage", conf.Storage).Msg("using mock sendgrid client")
		m.client = &mock.SendGridClient{
			Storage: conf.Storage,
		}
	} else {
		m.client = sendgrid.NewSendClient(conf.SendGridAPIKey)
	}

	// Parse the service email from the configuration
	if m.serviceEmail, err = mail.ParseAddress(conf.ServiceEmail); err != nil {
		return nil, fmt.Errorf("could not parse service email %q: %s", conf.ServiceEmail, err)
	}

	return m, nil
}

// EmailManager allows the BFF to send rich emails using the SendGrid service.
type EmailManager struct {
	conf         config.EmailConfig
	client       emails.EmailClient
	serviceEmail *mail.Address
}

func (m *EmailManager) Send(message *sgmail.SGMailV3) (err error) {
	var rep *rest.Response
	if rep, err = m.client.Send(message); err != nil {
		return err
	}

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		return errors.New(rep.Body)
	}

	return nil
}

// SendUserInvite sends an email to a user inviting them to join an organization.
func (m *EmailManager) SendUserInvite(user *management.User, inviter *management.User, org *models.Organization, inviteURL *url.URL) (err error) {
	var (
		userName, userEmail, inviterName, orgName string
	)

	if user.Email == nil || *user.Email == "" {
		return errors.New("user has no email address")
	}
	userEmail = *user.Email

	if user.Name != nil && *user.Name != "" {
		userName = *user.Name
	}

	switch {
	case inviter.Name != nil && *inviter.Name != "":
		inviterName = *inviter.Name
	case inviter.Email != nil && *inviter.Email != "":
		inviterName = *inviter.Email
	default:
		return errors.New("inviter user has no name or email address")
	}

	if org.Name != "" {
		orgName = org.Name
	} else {
		orgName = "their organization"
	}

	ctx := InviteUserData{
		Inviter:      inviterName,
		Organization: orgName,
		InviteURL:    inviteURL.String(),
	}

	if userName != "" {
		ctx.User = userName
	} else {
		ctx.User = userEmail
	}

	msg, err := InviteUserEmail(m.serviceEmail.Name, m.serviceEmail.Address, userName, userEmail, ctx)
	if err != nil {
		log.Error().Err(err).Msg("could not create user invite email")
		return err
	}

	if err = m.Send(msg); err != nil {
		log.Error().Err(err).Msg("could not send user invite email")
		return err
	}
	return nil
}
