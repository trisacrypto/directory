package emails

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"

	"github.com/ghodss/yaml"
	"github.com/rs/zerolog/log"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/encoding/protojson"
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
func (m *EmailManager) SendVerifyContacts(vasp *pb.VASP) (sent int, err error) {
	var contacts = []*pb.Contact{
		vasp.Contacts.Technical, vasp.Contacts.Administrative,
		vasp.Contacts.Billing, vasp.Contacts.Legal,
	}

	// Attempt at least one delivery, don't give up just because one email failed
	// Track how many emails and errors occurred during delivery.
	var nErrors int
	for idx, contact := range contacts {
		// Skip any null contacts or contacts without email addresses
		if contact == nil || contact.Email == "" {
			continue
		}

		var verified bool
		ctx := VerifyContactData{
			Name: contact.Name,
			VID:  vasp.Id,
		}
		if ctx.Token, verified, err = models.GetContactVerification(contact); err != nil {
			// If we can't get the verification token, this is a fatal error
			return sent, err
		}

		// If the contact has already been verified, then do not send a verify contact request
		if !verified {
			msg, err := VerifyContactEmail(
				m.serviceEmail.Name, m.serviceEmail.Address,
				contact.Name, contact.Email,
				ctx,
			)
			if err != nil {
				log.Error().Err(err).Str("vasp", vasp.Id).Int("contact", idx).Msg("could not create verify contact email")
				nErrors++
				continue
			}

			if err = m.Send(msg); err != nil {
				log.Error().Err(err).Str("vasp", vasp.Id).Int("contact", idx).Msg("could not send verify contact email")
				nErrors++
				continue
			}

			sent++
		}
	}

	// Return an error if no emails were delivered
	if sent == 0 {
		return sent, fmt.Errorf("no verify contact emails were successfully sent (%d errors)", nErrors)
	}
	return sent, nil
}

// SendReviewRequest is a shortcut for iComply verification in which we simply send
// an email to the TRISA admins and have them manually verify registrations.
func (m *EmailManager) SendReviewRequest(vasp *pb.VASP) (sent int, err error) {
	jsonpb := protojson.MarshalOptions{
		Multiline:       false,
		AllowPartial:    true,
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
	}

	var data []byte
	if data, err = jsonpb.Marshal(vasp); err != nil {
		return 0, err
	}

	// Convert JSON to YAML to make it more human readable
	// If the conversion fails, then the JSON data will be kept
	if yamlData, err := yaml.JSONToYAML(data); err == nil {
		data = yamlData
	}

	ctx := ReviewRequestData{
		VID:     vasp.Id,
		Request: string(data),
	}
	if ctx.Token, err = models.GetAdminVerificationToken(vasp); err != nil {
		return 0, err
	}

	msg, err := ReviewRequestEmail(
		m.serviceEmail.Name, m.serviceEmail.Address,
		m.adminsEmail.Name, m.adminsEmail.Address,
		ctx,
	)
	if err != nil {
		return 0, err
	}

	if err = m.Send(msg); err != nil {
		return 0, err
	}

	return 1, nil
}

// SendRejectRegistration sends a notification to all VASP contacts that their
// registration status is rejected without certificate issuance and explains why.
func (m *EmailManager) SendRejectRegistration(vasp *pb.VASP, reason string) (sent int, err error) {
	ctx := RejectRegistrationData{
		VID:    vasp.Id,
		Reason: reason,
	}

	var contacts = []*pb.Contact{
		vasp.Contacts.Technical, vasp.Contacts.Administrative,
		vasp.Contacts.Billing, vasp.Contacts.Legal,
	}

	// Attempt at least one delivery, don't give up just because one email failed
	// Track how many emails and errors occurred during delivery.
	var nErrors uint8
	for idx, contact := range contacts {
		// Skip any contacts that we can't send emails to
		if contact == nil || contact.Email == "" {
			continue
		}

		var verified bool
		if _, verified, err = models.GetContactVerification(contact); err != nil {
			// If we can't get the verification, this is a fatal error
			return sent, err
		}

		if verified {
			ctx.Name = contact.Name
			msg, err := RejectRegistrationEmail(
				m.serviceEmail.Name, m.serviceEmail.Address,
				contact.Name, contact.Email,
				ctx,
			)
			if err != nil {
				log.Error().Err(err).Str("vasp", vasp.Id).Int("contact", idx).Msg("could not create reject registration email")
				nErrors++
				continue
			}

			if err = m.Send(msg); err != nil {
				log.Error().Err(err).Str("vasp", vasp.Id).Int("contact", idx).Msg("could not send reject registration email")
				nErrors++
				continue
			}

			sent++
		}
	}

	// Return an error if no emails were delivered
	if sent == 0 {
		return sent, fmt.Errorf("no registration rejection emails were successfully sent (%d errors)", nErrors)
	}
	return sent, nil
}

// SendDeliverCertificates sends the PKCS12 encrypted certificate files to the VASP
// contacts as an attachment, completing the certificate issuance process. This method
// only sends the certificate attachment to one email (to limit the delivery of a secure
// email), ranking the contact emails by priority.
func (m *EmailManager) SendDeliverCertificates(vasp *pb.VASP, path string) (sent int, err error) {
	ctx := DeliverCertsData{
		VID:          vasp.Id,
		CommonName:   vasp.CommonName,
		SerialNumber: hex.EncodeToString(vasp.IdentityCertificate.SerialNumber),
		Endpoint:     vasp.TrisaEndpoint,
	}

	// These contacts are ordered by priority, e.g. first try to send to the technical
	// contact, then the administrative, etc.
	var contacts = []*pb.Contact{
		vasp.Contacts.Technical, vasp.Contacts.Administrative,
		vasp.Contacts.Legal, vasp.Contacts.Billing,
	}

	// Attempt at least one delivery, don't give up just because one email failed
	// Track how many emails and errors occurred during delivery.
	var nErrors uint8
	for idx, contact := range contacts {
		// Skip any null contacts or contacts without email addresses
		if contact == nil || contact.Email == "" {
			continue
		}

		var verified bool
		if _, verified, err = models.GetContactVerification(contact); err != nil {
			// If we can't get the verification this is a fatal error
			return sent, err
		}

		if verified {
			ctx.Name = contact.Name
			msg, err := DeliverCertsEmail(
				m.serviceEmail.Name, m.serviceEmail.Address,
				contact.Name, contact.Email,
				path, ctx,
			)

			if err != nil {
				log.Error().Err(err).Str("vasp", vasp.Id).Int("contact", idx).Msg("could not create deliver certs email")
				nErrors++
				continue
			}

			if err = m.Send(msg); err != nil {
				log.Error().Err(err).Str("vasp", vasp.Id).Int("contact", idx).Msg("could not send deliver certs email")
				nErrors++
				continue
			}

			// If we've successfully sent one cert delivery message, then stop sending
			// the message so that we only send it a single time.
			sent++
			break
		}
	}

	// Return an error if no emails were delivered
	if sent == 0 {
		return sent, fmt.Errorf("no certificate delivery emails were successfully sent (%d errors)", nErrors)
	}
	return sent, nil
}
