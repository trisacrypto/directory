package emails

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

// New email manager with the specified configuration.
func New(conf config.EmailConfig) (m *EmailManager, err error) {
	m = &EmailManager{
		conf:   conf,
		client: sendgrid.NewSendClient(conf.SendGridAPIKey),
	}

	// Parse the admin and service emails from the configuration
	if m.serviceEmail, err = mail.ParseAddress(conf.ServiceEmail); err != nil {
		return nil, fmt.Errorf("could not parse service email %q: %s", conf.ServiceEmail, err)
	}

	if m.adminsEmail, err = mail.ParseAddress(conf.AdminEmail); err != nil {
		return nil, fmt.Errorf("could not parse admin email %q: %s", conf.AdminEmail, err)
	}

	return m, nil
}

// EmailManager allows the server to send rich emails using the SendGrid service.
type EmailManager struct {
	conf         config.EmailConfig
	client       *sendgrid.Client
	serviceEmail *mail.Address
	adminsEmail  *mail.Address
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

// SendVerifyContacts creates a verification token for each contact in the VASP contact
// list and sends them the verification email with instructions on how to verify their
// email address.
func (m *EmailManager) SendVerifyContacts(vasp *pb.VASP) (err error) {
	var contacts = []*pb.Contact{
		vasp.Contacts.Technical, vasp.Contacts.Administrative,
		vasp.Contacts.Billing, vasp.Contacts.Legal,
	}

	for _, contact := range contacts {
		if contact == nil || contact.Email == "" {
			continue
		}

		ctx := VerifyContactData{
			Name: contact.Name,
			VID:  vasp.Id,
		}
		if ctx.Token, _, err = models.GetContactVerification(contact); err != nil {
			return err
		}

		msg, err := VerifyContactEmail(
			m.serviceEmail.Name, m.serviceEmail.Address,
			contact.Name, contact.Email,
			ctx,
		)
		if err != nil {
			return err
		}

		if err = m.Send(msg); err != nil {
			return err
		}
	}

	return nil
}

// SendReviewRequest is a shortcut for iComply verification in which we simply send
// an email to the TRISA admins and have them manually verify registrations.
func (m *EmailManager) SendReviewRequest(vasp *pb.VASP) (err error) {
	var data []byte
	if data, err = json.MarshalIndent(vasp, "", "  "); err != nil {
		return err
	}

	ctx := ReviewRequestData{
		VID:     vasp.Id,
		Request: string(data),
	}
	if ctx.Token, err = models.GetAdminVerificationToken(vasp); err != nil {
		return err
	}

	msg, err := ReviewRequestEmail(
		m.serviceEmail.Name, m.serviceEmail.Address,
		m.adminsEmail.Name, m.adminsEmail.Address,
		ctx,
	)
	if err != nil {
		return err
	}

	if err = m.Send(msg); err != nil {
		return err
	}

	return nil
}

// SendRejectRegistration sends a notification to all VASP contacts that their
// registration status is rejected without certificate issuance and explains why.
func (m *EmailManager) SendRejectRegistration(vasp *pb.VASP, reason string) (err error) {
	ctx := RejectRegistrationData{
		VID:    vasp.Id,
		Reason: reason,
	}

	var contacts = []*pb.Contact{
		vasp.Contacts.Technical, vasp.Contacts.Administrative,
		vasp.Contacts.Billing, vasp.Contacts.Legal,
	}

	for _, contact := range contacts {
		var verified bool
		if _, verified, err = models.GetContactVerification(contact); err != nil {
			return err
		}

		if contact != nil && verified {
			ctx.Name = contact.Name
			msg, err := RejectRegistrationEmail(
				m.serviceEmail.Name, m.serviceEmail.Address,
				contact.Name, contact.Email,
				ctx,
			)
			if err != nil {
				return err
			}

			if err = m.Send(msg); err != nil {
				return err
			}
		}
	}

	return nil
}

// SendDeliverCertificates sends the PKCS12 encrypted certificate files to the VASP
// contacts as an attachment, completing the certificate issuance process.
func (m *EmailManager) SendDeliverCertificates(vasp *pb.VASP, path string) (err error) {
	ctx := DeliverCertsData{
		VID:          vasp.Id,
		CommonName:   vasp.CommonName,
		SerialNumber: hex.EncodeToString(vasp.IdentityCertificate.SerialNumber),
		Endpoint:     vasp.TrisaEndpoint,
	}

	var contacts = []*pb.Contact{
		vasp.Contacts.Technical, vasp.Contacts.Administrative,
		vasp.Contacts.Billing, vasp.Contacts.Legal,
	}

	for _, contact := range contacts {
		var verified bool
		if _, verified, err = models.GetContactVerification(contact); err != nil {
			return err
		}
		if contact != nil && verified {
			ctx.Name = contact.Name
			msg, err := DeliverCertsEmail(
				m.serviceEmail.Name, m.serviceEmail.Address,
				contact.Name, contact.Email,
				path, ctx,
			)

			if err != nil {
				return err
			}

			if err = m.Send(msg); err != nil {
				return err
			}
		}
	}

	return nil
}
