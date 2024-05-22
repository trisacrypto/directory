package emails

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
)

// Compile all email templates into a top-level global variable.
var templates map[string]*template.Template

func init() {
	templates = make(map[string]*template.Template)
	for _, name := range AssetNames() {
		data := MustAssetString(name)
		templates[name] = template.Must(template.New(name).Parse(data))
	}
}

//===========================================================================
// Template Contexts
//===========================================================================

const (
	UnknownDate             = "unknown date"
	DateFormat              = "Monday, January 2, 2006"
	UnspecifiedOrganization = "Unknown Legal Name"
)

// VerifyContactData to complete the verify contact email templates.
type VerifyContactData struct {
	Name        string // Used to address the email
	Token       string // The unique token needed to verify the email
	VID         string // The ID of the VASP/Registration
	BaseURL     string // The URL of the verify contact endpoint to build the VerifyContactURL
	DirectoryID string // The registered directory to build a URL accessible by the BFF
}

// VerifyContactURL composes the link to verify the contact from the context. If the
// link is not able to be composed, the function returns an empty string and logs an
// error because without the link the email is fairly useless.
func (d VerifyContactData) VerifyContactURL() *url.URL {
	var (
		link *url.URL
		err  error
	)
	if d.BaseURL != "" {
		if link, err = url.Parse(d.BaseURL); err != nil {
			sentry.Error(nil).Err(err).Msg("could not include verify contact link in email, could not parse verify contact base url")
			return nil
		}
	} else {
		sentry.Error(nil).Msg("could not include verify contact link in email, no verify contact base url")
		return nil
	}

	params := link.Query()
	params.Set("vaspID", d.VID)
	params.Set("token", d.Token)
	params.Set("registered_directory", d.DirectoryID)
	link.RawQuery = params.Encode()
	return link
}

func (d VerifyContactData) VerifyContactURLUnencoded() template.HTML {
	url := d.VerifyContactURL()
	return template.HTML(url.String())
}

// ReviewRequestData to complete review request email templates.
type ReviewRequestData struct {
	VID                 string // The ID of the VASP/Registration
	Token               string // The unique token needed to review the registration
	Request             string // The review request data as a nicely formatted JSON or YAML string
	RegisteredDirectory string // The directory name for the review request
	Attachment          []byte // Data to attach to the email
	BaseURL             string // The URL of the admin review endpoint to build the AdminReviewURL
}

// AdminReviewURL composes a link to the VASP detail in the admin UI. If the base url is
// missing or can't be parsed, it logs an warning and returns empty string.
// The AdminReviewURL is useful, but it is not a critical error.
func (d ReviewRequestData) AdminReviewURL() string {
	var (
		link *url.URL
		err  error
	)
	if d.BaseURL != "" {
		if link, err = url.Parse(d.BaseURL); err != nil {
			sentry.Warn(nil).Err(err).Msg("could not include admin review link in email, could not parse admin base url")
			return ""
		}
	} else {
		sentry.Warn(nil).Msg("could not include admin review link in email, no admin base url")
		return ""
	}
	return link.ResolveReference(&url.URL{Path: d.VID}).String()
}

// RejectRegistrationData to complete reject registration email templates.
type RejectRegistrationData struct {
	Name                string // Used to address the email
	VID                 string // The ID of the VASP/Registration
	Organization        string // The name of the organization (if it exists)
	CommonName          string // The common name assigned to the cert
	RegisteredDirectory string // The directory name for the certificates being issued
	Reason              string // A description of why the registration request was rejected
}

// DeliverCertsData to complete deliver certs email templates.
type DeliverCertsData struct {
	Name                string // Used to address the email
	VID                 string // The ID of the VASP/Registration
	Organization        string // The name of the organization (if it exists)
	CommonName          string // The common name assigned to the cert
	SerialNumber        string // The serial number of the certificate
	Endpoint            string // The expected endpoint for the TRISA service
	RegisteredDirectory string // The directory name for the certificates being issued
}

// ExpiresAdminNotificationData to complete expires admin notification email templates.
type ExpiresAdminNotificationData struct {
	VID                 string    // The ID of the VASP/Registration
	Organization        string    // The name of the organization (if it exists)
	CommonName          string    // The common name assigned to the cert
	SerialNumber        string    // The serial number of the certificate
	Endpoint            string    // The expected endpoint for the TRISA service
	RegisteredDirectory string    // The directory name for the certificates being issued
	Expiration          time.Time // The timestamp that the certificates expire
	Reissuance          time.Time // The timestamp the certificates will be reissued
	BaseURL             string    // The URL of the admin review endpoint to build the AdminReviewURL
}

// AdminReviewURL composes a link to the VASP detail in the admin UI. If the base url is
// missing or can't be parsed, it logs an warning and returns empty string.
// The AdminReviewURL is useful, but it is not a critical error.
func (d ExpiresAdminNotificationData) AdminReviewURL() string {
	var (
		link *url.URL
		err  error
	)
	if d.BaseURL != "" {
		if link, err = url.Parse(d.BaseURL); err != nil {
			sentry.Warn(nil).Err(err).Msg("could not include admin review link in email, could not parse admin base url")
			return ""
		}
	} else {
		sentry.Warn(nil).Msg("could not include admin review link in email, no admin base url")
		return ""
	}
	return link.ResolveReference(&url.URL{Path: d.VID}).String()
}

// ExpirationDate formats the expiration date for rendering in the email.
func (d ExpiresAdminNotificationData) ExpirationDate() string {
	if d.Expiration.IsZero() {
		return UnknownDate
	}
	return d.Expiration.Format(DateFormat)
}

// ReissueDate formats the reissuance date for rendering in the email.
func (d ExpiresAdminNotificationData) ReissueDate() string {
	if d.Reissuance.IsZero() {
		return UnknownDate
	}
	return d.Reissuance.Format(DateFormat)
}

// ReissuanceReminderData to complete reissue reminder email templates.
type ReissuanceReminderData struct {
	Name                string    // Used to address the email
	VID                 string    // The ID of the VASP/Registration
	Organization        string    // The name of the organization (if it exists)
	CommonName          string    // The common name assigned to the cert
	SerialNumber        string    // The serial number of the certificate
	Endpoint            string    // The expected endpoint for the TRISA service
	RegisteredDirectory string    // The directory name for the certificates being issued
	Expiration          time.Time // The timestamp that the certificates expire
	Reissuance          time.Time // The timestamp the certificates will be reissued
}

// ExpirationDate formats the expiration date for rendering in the email.
func (d ReissuanceReminderData) ExpirationDate() string {
	if d.Expiration.IsZero() {
		return UnknownDate
	}
	return d.Expiration.Format(DateFormat)
}

// ReissueDate formats the reissuance date for rendering in the email.
func (d ReissuanceReminderData) ReissueDate() string {
	if d.Reissuance.IsZero() {
		return UnknownDate
	}
	return d.Reissuance.Format(DateFormat)
}

// ReissuanceStartedData to complete reissue reminder email templates.
type ReissuanceStartedData struct {
	Name                string // Used to address the email
	VID                 string // The ID of the VASP/Registration
	Organization        string // The name of the organization (if it exists)
	CommonName          string // The common name assigned to the cert
	Endpoint            string // The expected endpoint for the TRISA service
	RegisteredDirectory string // The directory name for the certificates being issued
	WhisperURL          string // Secure one-time whisper link for password retrieval
}

// ReissuanceAdminNotificationData to complete reissuance admin notification email templates.
type ReissuanceAdminNotificationData struct {
	VID                 string    // The ID of the VASP/Registration
	Organization        string    // The name of the organization (if it exists)
	CommonName          string    // The common name assigned to the cert
	SerialNumber        string    // The serial number of the certificate
	Endpoint            string    // The expected endpoint for the TRISA service
	RegisteredDirectory string    // The directory name for the certificate that was issued
	Expiration          time.Time // The timestamp when the certificate expires
	Reissuance          time.Time // The timestamp when the certificate was reissued
	BaseURL             string    // The URL of the admin review endpoint to build the AdminReviewURL
}

// AdminReviewURL composes a link to the VASP detail in the admin UI. If the base url is
// missing or can't be parsed, it logs an warning and returns empty string.
// The AdminReviewURL is useful, but it is not a critical error.
func (d ReissuanceAdminNotificationData) AdminReviewURL() string {
	var (
		link *url.URL
		err  error
	)
	if d.BaseURL != "" {
		if link, err = url.Parse(d.BaseURL); err != nil {
			sentry.Warn(nil).Err(err).Msg("could not include admin review link in email, could not parse admin base url")
			return ""
		}
	} else {
		sentry.Warn(nil).Msg("could not include admin review link in email, no admin base url")
		return ""
	}
	return link.ResolveReference(&url.URL{Path: d.VID}).String()
}

// ExpirationDate formats the expiration date for rendering in the email.
func (d ReissuanceAdminNotificationData) ExpirationDate() string {
	if d.Expiration.IsZero() {
		return UnknownDate
	}
	return d.Expiration.Format(DateFormat)
}

// ReissueDate formats the reissuance date for rendering in the email.
func (d ReissuanceAdminNotificationData) ReissueDate() string {
	if d.Reissuance.IsZero() {
		return UnknownDate
	}
	return d.Reissuance.Format(DateFormat)
}

//===========================================================================
// Email Builders
//===========================================================================

// VerifyContactEmail creates a new verify contact email, ready for sending by rendering
// the text and html templates with the supplied data then constructing a sendgrid email.
func VerifyContactEmail(sender, senderEmail, recipient, recipientEmail string, data VerifyContactData) (message *mail.SGMailV3, err error) {
	var text, html string
	if text, html, err = Render("verify_contact", data); err != nil {
		return nil, err
	}

	return mail.NewSingleEmail(
		mail.NewEmail(sender, senderEmail),
		VerifyContactRE,
		mail.NewEmail(recipient, recipientEmail),
		text,
		html,
	), nil
}

// ReviewRequestEmail creates a new review request email, ready for sending by rendering
// the text and html templates with the supplied data then constructing a sendgrid email.
func ReviewRequestEmail(sender, senderEmail, recipient, recipientEmail string, data ReviewRequestData) (message *mail.SGMailV3, err error) {
	var text, html string
	if text, html, err = Render("review_request", data); err != nil {
		return nil, err
	}

	message = mail.NewSingleEmail(
		mail.NewEmail(sender, senderEmail),
		ReviewRequestRE,
		mail.NewEmail(recipient, recipientEmail),
		text,
		html,
	)

	if err = AttachJSON(message, data.Attachment, fmt.Sprintf("%s.json", data.VID)); err != nil {
		// Log the error but do not stop sending the message
		sentry.Error(nil).Err(err).Msg("could not attach JSON data to review request email")
	}
	return message, nil
}

// RejectRegistrationEmail creates a new reject registration email, ready for sending by
// rendering the text and html templates with the supplied data then constructing a
// sendgrid email.
func RejectRegistrationEmail(sender, senderEmail, recipient, recipientEmail string, data RejectRegistrationData) (message *mail.SGMailV3, err error) {
	var text, html string
	if text, html, err = Render("reject_registration", data); err != nil {
		return nil, err
	}

	return mail.NewSingleEmail(
		mail.NewEmail(sender, senderEmail),
		RejectRegistrationRE,
		mail.NewEmail(recipient, recipientEmail),
		text,
		html,
	), nil
}

// DeliverCertsEmail creates a new deliver certs email, ready for sending by rendering
// the text and html templates with the supplied data, loading the attachment from disk
// then constructing a sendgrid email.
func DeliverCertsEmail(sender, senderEmail, recipient, recipientEmail, attachmentPath string, data DeliverCertsData) (message *mail.SGMailV3, err error) {
	var text, html string
	if text, html, err = Render("deliver_certs", data); err != nil {
		return nil, err
	}

	message = mail.NewSingleEmail(
		mail.NewEmail(sender, senderEmail),
		DeliverCertsRE,
		mail.NewEmail(recipient, recipientEmail),
		text,
		html,
	)

	// Add attachment from a file on disk.
	if err = LoadAttachment(message, attachmentPath); err != nil {
		return nil, err
	}

	return message, nil
}

// ExpiresAdminNotificationEmail creates a new certs expired admin notification email,
// ready for sending by rendering the text and html templates with the supplied data.
func ExpiresAdminNotificationEmail(sender, senderEmail, recipient, recipientEmail string, data ExpiresAdminNotificationData) (message *mail.SGMailV3, err error) {
	var text, html string
	if text, html, err = Render("expires_admin_notification", data); err != nil {
		return nil, err
	}

	message = mail.NewSingleEmail(
		mail.NewEmail(sender, senderEmail),
		ExpiresAdminNotificationRE,
		mail.NewEmail(recipient, recipientEmail),
		text,
		html,
	)

	return message, nil
}

// ReissuanceReminderEmail creates a new reissuance reminder email, ready for sending by
// rendering the text and html templates with the supplied data.
func ReissuanceReminderEmail(sender, senderEmail, recipient, recipientEmail string, data ReissuanceReminderData) (message *mail.SGMailV3, err error) {
	var text, html string
	if text, html, err = Render("reissuance_reminder", data); err != nil {
		return nil, err
	}

	message = mail.NewSingleEmail(
		mail.NewEmail(sender, senderEmail),
		ReissuanceReminderRE,
		mail.NewEmail(recipient, recipientEmail),
		text,
		html,
	)

	return message, nil
}

// ReissuanceStartedEmail creates a new reissuance started email, ready for sending by
// rendering the text and html templates with the supplied data.
func ReissuanceStartedEmail(sender, senderEmail, recipient, recipientEmail string, data ReissuanceStartedData) (message *mail.SGMailV3, err error) {
	var text, html string
	if text, html, err = Render("reissuance_started", data); err != nil {
		return nil, err
	}

	message = mail.NewSingleEmail(
		mail.NewEmail(sender, senderEmail),
		ReissuanceStartedRE,
		mail.NewEmail(recipient, recipientEmail),
		text,
		html,
	)

	return message, nil
}

// ReissuanceAdminNotificationEmail creates a new certs reissuance admin notification email,
// ready for sending by rendering the text and html templates with the supplied data.
func ReissuanceAdminNotificationEmail(sender, senderEmail, recipient, recipientEmail string, data ReissuanceAdminNotificationData) (message *mail.SGMailV3, err error) {
	var text, html string
	if text, html, err = Render("reissuance_admin_notification", data); err != nil {
		return nil, err
	}

	message = mail.NewSingleEmail(
		mail.NewEmail(sender, senderEmail),
		ReissuanceAdminNotificationRE,
		mail.NewEmail(recipient, recipientEmail),
		text,
		html,
	)

	return message, nil
}

//===========================================================================
// Template Builders
//===========================================================================

// Render returns the text and html executed templates for the specified name and data.
// Ensure that the extension is not supplied to the render method.
func Render(name string, data interface{}) (text, html string, err error) {
	if text, err = render(name+".txt", data); err != nil {
		return "", "", err
	}

	if html, err = render(name+".html", data); err != nil {
		return "", "", err
	}

	return text, html, nil
}

func render(name string, data interface{}) (_ string, err error) {
	var (
		ok bool
		t  *template.Template
	)

	if t, ok = templates[name]; !ok {
		return "", fmt.Errorf("could not find %q in templates", name)
	}

	buf := &strings.Builder{}
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// LoadAttachment onto email from a file on disk.
func LoadAttachment(message *mail.SGMailV3, attachmentPath string) (err error) {
	// Read and encode the attachment data
	var data []byte
	if data, err = os.ReadFile(attachmentPath); err != nil {
		return err
	}
	encoded := base64.StdEncoding.EncodeToString(data)

	// Create the attachment
	// TODO: detect mimetype rather than assuming zip
	attach := mail.NewAttachment()
	attach.SetContent(encoded)
	attach.SetType("application/zip")
	attach.SetFilename(filepath.Base(attachmentPath))
	attach.SetDisposition("attachment")
	message.AddAttachment(attach)
	return nil
}

// AttachJSON by marshaling the specified data into human-readable data and encode and
// attach it to the email as a file.
func AttachJSON(message *mail.SGMailV3, data []byte, filename string) (err error) {
	// Encode the data to attach to the email
	encoded := base64.StdEncoding.EncodeToString(data)

	// Create the attachment
	attach := mail.NewAttachment()
	attach.SetContent(encoded)
	attach.SetType("application/json")
	attach.SetFilename(filename)
	attach.SetDisposition("attachment")
	message.AddAttachment(attach)
	return nil

}

//===========================================================================
// Testing Interfaces
//===========================================================================

type AdminReview interface {
	AdminReviewURL() string
}

type FutureReissuer interface {
	ExpirationDate() string
	ReissueDate() string
}
