package emails

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/utils/emails"
	"github.com/trisacrypto/directory/pkg/utils/emails/mock"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// New email manager with the specified configuration.
func New(conf config.EmailConfig) (m *EmailManager, err error) {
	m = &EmailManager{conf: conf}
	if conf.Testing {
		log.Warn().Bool("testing", conf.Testing).Str("storage", conf.Storage).Msg("using mock sendgrid client")
		m.client = &mock.SendGridClient{
			Storage: conf.Storage,
		}
	} else {
		if conf.SendGridAPIKey == "" {
			return nil, errors.New("cannot create sendgrid client without API key")
		}
		m.client = sendgrid.NewSendClient(conf.SendGridAPIKey)
	}

	// Warn if email configuration isn't complete and will produce partial emails.
	if conf.VerifyContactBaseURL == "" || conf.AdminReviewBaseURL == "" {
		log.Warn().
			Bool("missing_verify_contact_base_url", conf.VerifyContactBaseURL == "").
			Bool("missing_admin_review_base_url", conf.AdminReviewBaseURL == "").
			Msg("partial email configuration, some emails may not include links")
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
	client       emails.EmailClient
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
// email address. Caller must update the VASP record on the data store after calling
// this function.
func (m *EmailManager) SendVerifyContacts(vasp *pb.VASP) (sent int, err error) {
	// Attempt at least one delivery, don't give up just because one email failed
	// Track how many emails and errors occurred during delivery.
	var nErrors int
	iter := models.NewContactIterator(vasp.Contacts, true, false)
	for iter.Next() {
		contact, kind := iter.Value()

		var verified bool
		if _, verified, err = models.GetContactVerification(contact); err != nil {
			log.Error().Err(err).Str("vasp", vasp.Id).Msg("failed to get contact verification")
			return sent, err
		}

		if !verified {
			if err := m.SendVerifyContact(vasp, contact); err != nil {
				nErrors++
				log.Error().Err(err).Str("vasp", vasp.Id).Str("contact", kind).Msg("failed to send verify contact email")
			} else {
				sent++
			}
		}
	}

	// Return an error if no emails were delivered
	if sent == 0 {
		return sent, fmt.Errorf("no verify contact emails were successfully sent (%d errors)", nErrors)
	}
	return sent, nil
}

// SendVerifyContact sends a verification email to a contact.
func (m *EmailManager) SendVerifyContact(vasp *pb.VASP, contact *pb.Contact) (err error) {
	ctx := VerifyContactData{
		Name:        contact.Name,
		VID:         vasp.Id,
		BaseURL:     m.conf.VerifyContactBaseURL,
		DirectoryID: m.conf.DirectoryID,
	}

	if ctx.Token, _, err = models.GetContactVerification(contact); err != nil {
		log.Error().Err(err).Str("vasp", vasp.Id).Msg("failed to get contact verification")
		return err
	}

	msg, err := VerifyContactEmail(
		m.serviceEmail.Name, m.serviceEmail.Address,
		contact.Name, contact.Email,
		ctx,
	)
	if err != nil {
		log.Error().Err(err).Msg("could not create verify contact email")
		return err
	}

	if err = m.Send(msg); err != nil {
		log.Error().Err(err).Msg("could not send verify contact email")
		return err
	}

	if err = models.AppendEmailLog(contact, string(admin.ResendVerifyContact), msg.Subject); err != nil {
		log.Error().Err(err).Msg("could not log verify contact email")
	}
	return nil
}

// SendReviewRequest is a shortcut for iComply verification in which we simply send
// an email to the TRISA admins and have them manually verify registrations.
func (m *EmailManager) SendReviewRequest(vasp *pb.VASP) (sent int, err error) {
	// Create the template context with the admin verification token
	ctx := ReviewRequestData{
		VID:                 vasp.Id,
		RegisteredDirectory: m.conf.DirectoryID,
		BaseURL:             m.conf.AdminReviewBaseURL,
	}
	if ctx.Token, err = models.GetAdminVerificationToken(vasp); err != nil {
		return 0, err
	}

	// Remove sensitive data so it's not sent in the form.
	clone := proto.Clone(vasp).(*pb.VASP)
	models.SetAdminVerificationToken(clone, "[REDACTED]")

	iter := models.NewContactIterator(clone.Contacts, false, false)
	for iter.Next() {
		contact, _ := iter.Value()
		_, verified, _ := models.GetContactVerification(contact)
		models.SetContactVerification(contact, "[REDACTED]", verified)
	}

	// Marshal the VASP struct for review in the email.
	jsonpb := protojson.MarshalOptions{
		Multiline:       false,
		AllowPartial:    true,
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
	}

	var data []byte
	if data, err = jsonpb.Marshal(clone); err != nil {
		return 0, err
	}

	// Convert JSON to YAML to make it more human readable
	// If the conversion fails, then the JSON data will be kept
	if yamlData, err := yaml.JSONToYAML(data); err == nil {
		data = yamlData
	}
	ctx.Request = string(data)

	// Attach the JSON data as an attachment
	ctx.Attachment = data

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
// Caller must update the VASP record on the data store after calling this function.
func (m *EmailManager) SendRejectRegistration(vasp *pb.VASP, reason string) (sent int, err error) {
	var errs *multierror.Error
	ctx := RejectRegistrationData{
		VID:    vasp.Id,
		Reason: reason,
	}

	// Attempt at least one delivery, don't give up just because one email failed
	// Track how many emails and errors occurred during delivery.
	iter := models.NewContactIterator(vasp.Contacts, true, true)
	for iter.Next() {
		var contact *pb.Contact
		var kind string
		contact, kind = iter.Value()
		ctx.Name = contact.Name
		msg, err := RejectRegistrationEmail(
			m.serviceEmail.Name, m.serviceEmail.Address,
			contact.Name, contact.Email,
			ctx,
		)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("could not create reject registration email for %s contact: %s", kind, err))
			continue
		}

		if err = m.Send(msg); err != nil {
			errs = multierror.Append(errs, fmt.Errorf("could not send reject registration email for %s contact: %s", kind, err))
			continue
		}

		sent++

		if err = models.AppendEmailLog(contact, string(admin.ResendRejection), msg.Subject); err != nil {
			errs = multierror.Append(errs, fmt.Errorf("could not log reject registration email for %s contact: %s", kind, err))
			continue
		}
	}

	if iterErrs := iter.Error(); iterErrs != nil {
		errs = multierror.Append(errs, iterErrs)
	}

	if sent == 0 {
		errs = multierror.Append(errs, fmt.Errorf("no reject registration emails were successfully sent"))
	}

	return sent, errs.ErrorOrNil()
}

// SendDeliverCertificates sends the PKCS12 encrypted certificate files to the VASP
// contacts as an attachment, completing the certificate issuance process. This method
// only sends the certificate attachment to one email (to limit the delivery of a secure
// email), ranking the contact emails by priority. Caller must update the VASP record on
// the data store after calling this function.
func (m *EmailManager) SendDeliverCertificates(vasp *pb.VASP, path string) (sent int, err error) {
	var errs *multierror.Error
	ctx := DeliverCertsData{
		VID:                 vasp.Id,
		CommonName:          vasp.CommonName,
		SerialNumber:        hex.EncodeToString(vasp.IdentityCertificate.SerialNumber),
		Endpoint:            vasp.TrisaEndpoint,
		RegisteredDirectory: m.conf.DirectoryID,
	}

	// Attempt at least one delivery, don't give up just because one email failed
	// Track how many emails and errors occurred during delivery.
	// Note: new contact iterator provides the contact email prioritization order.
	iter := models.NewContactIterator(vasp.Contacts, true, true)
	for iter.Next() {
		var contact *pb.Contact
		var kind string
		contact, kind = iter.Value()
		ctx.Name = contact.Name
		msg, err := DeliverCertsEmail(
			m.serviceEmail.Name, m.serviceEmail.Address,
			contact.Name, contact.Email,
			path, ctx,
		)

		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("could not create deliver certificates email for %s contact: %s", kind, err))
			continue
		}

		if err = m.Send(msg); err != nil {
			errs = multierror.Append(errs, fmt.Errorf("could not send deliver certificates email for %s contact: %s", kind, err))
			continue
		}

		sent++

		if err = models.AppendEmailLog(contact, string(admin.ResendDeliverCerts), msg.Subject); err != nil {
			errs = multierror.Append(errs, fmt.Errorf("could not log deliver certificates email for %s contact: %s", kind, err))
			continue
		}

		// If we've successfully sent one cert delivery message, then stop sending
		// the message so that we only send it a single time.
		break
	}

	if iterErrs := iter.Error(); iterErrs != nil {
		errs = multierror.Append(errs, iterErrs)
	}

	if sent == 0 {
		errs = multierror.Append(errs, fmt.Errorf("no deliver certificates emails were successfully sent"))
	}

	return sent, errs.ErrorOrNil()
}

// SendExpiresAdminNotification sends the admins a notice that an identity certificate
// will be expiring soon. This allows the admins to determine if a new review of the
// TRISA member is necessary before the reissuance process begins.
func (m *EmailManager) SendExpiresAdminNotification(vasp *pb.VASP, timeWindow int, reissueDate time.Time) (sent int, err error) {
	// Make sure the email has not already been sent recently
	if emailCount, err := models.GetSentAdminEmailCount(vasp, string(admin.ReissuanceReminder), timeWindow); err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("error retrieving admin email log for %s's reissuance reminder", vasp.Id))
		return 0, err
	} else if emailCount > 0 {
		return 0, nil
	}

	// Create the template context
	ctx := ExpiresAdminNotificationData{
		VID:                 vasp.Id,
		CommonName:          vasp.CommonName,
		Endpoint:            vasp.TrisaEndpoint,
		RegisteredDirectory: m.conf.DirectoryID,
		Reissuance:          reissueDate,
		BaseURL:             m.conf.AdminReviewBaseURL,
	}

	if vasp.IdentityCertificate != nil {
		// TODO: ensure the timestamp format is correct
		ctx.SerialNumber = strings.ToUpper(hex.EncodeToString(vasp.IdentityCertificate.SerialNumber))
		ctx.Expiration, _ = time.Parse(time.RFC3339, vasp.IdentityCertificate.NotAfter)
	}

	msg, err := ExpiresAdminNotificationEmail(
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
	sent++

	if err = models.AppendAdminEmailLog(vasp, string(admin.ReissuanceReminder), msg.Subject); err != nil {
		return 0, err
	}

	return sent, nil
}

// Helper function for HandleCertificateReissuance that sends reissuance reminder emails
// a vasp's contacts, ensuring at least one of the contact's receives the reminder or a
// critical alert is raised.
func (m *EmailManager) SendContactReissuanceReminder(vasp *pb.VASP, timeWindow int, reissuanceDate time.Time) (err error) {
	ReissuanceData := ReissuanceReminderData{
		VID:                 vasp.Id,
		CommonName:          vasp.CommonName,
		Endpoint:            vasp.TrisaEndpoint,
		RegisteredDirectory: m.conf.ServiceEmail,
		Reissuance:          reissuanceDate,
	}

	if vasp.IdentityCertificate != nil {
		ReissuanceData.SerialNumber = strings.ToUpper(hex.EncodeToString(vasp.IdentityCertificate.SerialNumber))
		if ReissuanceData.Expiration, err = time.Parse(time.RFC3339, vasp.IdentityCertificate.NotAfter); err != nil {
			return fmt.Errorf("could not parse vasp certificate expiration date for %s", vasp.Id)
		}
	}

	// Iterate through the VASP's verified contacts and send the reissuance reminder email.
	contactsToNotify, err := getContactsToNotify(vasp.Contacts)
	if err != nil {
		return err
	}
	if len(contactsToNotify) == 0 {
		log.WithLevel(zerolog.FatalLevel).Str("vasp_id", vasp.Id).Msg("cert-manager could not find a verified contact for vasp")
		return nil
	}

	for _, contact := range contactsToNotify {
		// Make sure that the reminder email hasn't already been sent to this contact.
		reissuanceReminder := string(admin.ReissuanceReminder)
		emailCount, err := models.GetSentEmailCount(contact, reissuanceReminder, timeWindow)
		if err != nil {
			log.Error().Err(err).Str("vasp_id", vasp.Id).Str("contact", contact.Name).Msg("could not retrieve email count from email log")
			continue
		}
		if emailCount > 0 {
			continue
		}

		// Create the reissuance reminder email.
		ReissuanceData.Name = contact.Name
		msg, err := ReissuanceReminderEmail(
			m.serviceEmail.Name, m.serviceEmail.Address,
			contact.Name, contact.Email,
			ReissuanceData,
		)
		if err != nil {
			log.Error().Err(err).Str("vasp_id", vasp.Id).Str("contact", contact.Name).Msg("could not create reissuance reminder email")
			continue
		}

		if err = m.Send(msg); err != nil {
			log.Error().Err(err).Str("vasp_id", vasp.Id).Str("contact", contact.Name).Msg("error sending reissuance reminder email")
			continue
		}
		if err = models.AppendEmailLog(contact, reissuanceReminder, msg.Subject); err != nil {
			log.Error().Err(err).Str("vasp_id", vasp.Id).Str("contact", contact.Name).Msg("error appending to email log")
		}
	}
	return nil
}

// Helper function for SendContactReissuanceReminder that builds the list of verified contacts
// to send reissuance reminder emails to based on the following logic:
// 		1. Send to the Technical contact if verified, else
// 		2. Send to the Administrative contact if verified, else
// 		3. Send to all other verified contacts
func getContactsToNotify(contacts *pb.Contacts) (contactsToNotify []*pb.Contact, err error) {
	if verified, err := models.ContactIsVerified(contacts.Technical); err != nil {
		return nil, err
	} else if verified {
		return []*pb.Contact{contacts.Technical}, nil
	}

	if verified, err := models.ContactIsVerified(contacts.Administrative); err != nil {
		return nil, err
	} else if verified {
		return []*pb.Contact{contacts.Administrative}, nil
	}

	if verified, err := models.ContactIsVerified(contacts.Legal); err != nil {
		return nil, err
	} else if verified {
		contactsToNotify = append(contactsToNotify, contacts.Legal)
	}

	if verified, err := models.ContactIsVerified(contacts.Billing); err != nil {
		return nil, err
	} else if verified {
		contactsToNotify = append(contactsToNotify, contacts.Billing)
	}

	return contactsToNotify, nil
}

// SendReissuanceReminder sends a reminder to all verified contacts that their identity
// certificates will be expiring soon and that the system will automatically reissue the
// certs on a particular date.
func (m *EmailManager) SendReissuanceReminder(vasp *pb.VASP, reissueDate time.Time) (sent int, err error) {
	var errs *multierror.Error
	ctx := ReissuanceReminderData{
		VID:                 vasp.Id,
		CommonName:          vasp.CommonName,
		Endpoint:            vasp.TrisaEndpoint,
		RegisteredDirectory: m.conf.DirectoryID,
		Reissuance:          reissueDate,
	}

	if vasp.IdentityCertificate != nil {
		// TODO: ensure the timestamp format is correct
		ctx.SerialNumber = strings.ToUpper(hex.EncodeToString(vasp.IdentityCertificate.SerialNumber))
		ctx.Expiration, _ = time.Parse(time.RFC3339, vasp.IdentityCertificate.NotAfter)
	}

	// Attempt at least one delivery, don't give up just because one email failed.
	// Track how many emails and errors occurred during delivery.
	iter := models.NewContactIterator(vasp.Contacts, true, true)
	for iter.Next() {
		contact, kind := iter.Value()
		ctx.Name = contact.Name

		msg, err := ReissuanceReminderEmail(
			m.serviceEmail.Name, m.serviceEmail.Address,
			contact.Name, contact.Email,
			ctx,
		)

		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("could not create reissuance reminder email for %s contact: %s", kind, err))
			continue
		}

		if err = m.Send(msg); err != nil {
			errs = multierror.Append(errs, fmt.Errorf("could not send reissuance reminder email for %s contact: %s", kind, err))
			continue
		}

		sent++

		if err = models.AppendEmailLog(contact, string(admin.ReissuanceReminder), msg.Subject); err != nil {
			errs = multierror.Append(errs, fmt.Errorf("could not log reissuance reminder email for %s contact: %s", kind, err))
			continue
		}
	}

	if iterErrs := iter.Error(); iterErrs != nil {
		errs = multierror.Append(errs, iterErrs)
	}

	if sent == 0 {
		errs = multierror.Append(errs, fmt.Errorf("no reissuance reminder emails were successfully sent"))
	}

	return sent, errs.ErrorOrNil()
}

// SendReissuanceStarted sends the PKCS12 password via a secure one time link. This
// method only sends the PKCS12 password to one email (to limit the delivery of secure
// emails), ranking the contact emails by priority.
func (m *EmailManager) SendReissuanceStarted(vasp *pb.VASP, whisperLink string) (sent int, err error) {
	var errs *multierror.Error
	ctx := ReissuanceStartedData{
		VID:                 vasp.Id,
		CommonName:          vasp.CommonName,
		Endpoint:            vasp.TrisaEndpoint,
		RegisteredDirectory: m.conf.DirectoryID,
		WhisperURL:          whisperLink,
	}

	// Attempt at least one delivery, don't give up just because one email failed.
	// Track how many emails and errors are occurring during delivery.
	// Note: new contact iterator provides the contact email prioritization order.
	iter := models.NewContactIterator(vasp.Contacts, true, true)
	for iter.Next() {
		contact, kind := iter.Value()
		ctx.Name = contact.Name

		msg, err := ReissuanceStartedEmail(
			m.serviceEmail.Name, m.serviceEmail.Address,
			contact.Name, contact.Email,
			ctx,
		)

		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("could not create reissuance started email for %s contact: %s", kind, err))
			continue
		}

		if err = m.Send(msg); err != nil {
			errs = multierror.Append(errs, fmt.Errorf("could not send reissuance started email for %s contact: %s", kind, err))
			continue
		}

		sent++

		if err = models.AppendEmailLog(contact, string(admin.ReissuanceStarted), msg.Subject); err != nil {
			errs = multierror.Append(errs, fmt.Errorf("could not log reissuance started email for %s contact: %s", kind, err))
			continue
		}

		// If we've successfully send one reissuance started message, ten stop sending
		// messages to minimize how many contacts receive the secure one-time link.
		break
	}

	if iterErrs := iter.Error(); iterErrs != nil {
		errs = multierror.Append(errs, iterErrs)
	}

	if sent == 0 {
		errs = multierror.Append(errs, fmt.Errorf("no reissuance started emails were successfully sent"))
	}

	return sent, errs.ErrorOrNil()
}

// SendReissuanceAdminNotification sends the admins a notice that an identity certificate
// has been reissued. This allows the admins to know that the reissuance has been done automatically
func (m *EmailManager) SendReissuanceAdminNotification(vasp *pb.VASP, timeWindow int, reissueDate time.Time) (sent int, err error) {
	// Make sure the email has not already been sent recently
	if emailCount, err := models.GetSentAdminEmailCount(vasp, string(admin.ReissuanceStarted), timeWindow); err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("error retrieving admin email log for %s's reissuance notification", vasp.Id))
		return 0, err
	} else if emailCount > 0 {
		return 0, nil
	}

	// Create the template context
	ctx := ReissuanceAdminNotificationData{
		VID:                 vasp.Id,
		CommonName:          vasp.CommonName,
		Endpoint:            vasp.TrisaEndpoint,
		RegisteredDirectory: m.conf.DirectoryID,
		Reissuance:          reissueDate,
		BaseURL:             m.conf.AdminReviewBaseURL,
	}

	if vasp.IdentityCertificate != nil {
		// TODO: ensure the timestamp format is correct
		ctx.SerialNumber = strings.ToUpper(hex.EncodeToString(vasp.IdentityCertificate.SerialNumber))
		ctx.Expiration, err = time.Parse(time.RFC3339, vasp.IdentityCertificate.NotAfter)
		if err != nil {
			return 0, fmt.Errorf("error parsing timestamp: %v", err)
		}
	} else if vasp.IdentityCertificate == nil {
		return 0, fmt.Errorf("no identity certificate for vasp %s", vasp.Id)
	}

	// Create reissuance admin notifications email.
	msg, err := ReissuanceAdminNotificationEmail(
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
	sent++

	if err = models.AppendAdminEmailLog(vasp, string(admin.ReissuanceStarted), msg.Subject); err != nil {
		return sent, err
	}

	return sent, nil
}
