package emails

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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

// VerifyContactData to complete the verify contact email templates.
type VerifyContactData struct {
	Name    string // Used to address the email
	Token   string // The unique token needed to verify the email
	VID     string // The ID of the VASP/Registration
	BaseURL string // The URL of the verify contact endpoint to build the VerifyContactURL
}

// TODO: how to return errors instead of panic inside of template execution?
func (d VerifyContactData) VerifyContactURL() string {
	var (
		link *url.URL
		err  error
	)
	if d.BaseURL != "" {
		if link, err = url.Parse(d.BaseURL); err != nil {
			panic(err)
		}
	} else {
		panic("verify contact base url missing")
	}

	params := link.Query()
	params.Set("vaspID", d.VID)
	params.Set("token", d.Token)
	link.RawQuery = params.Encode()
	return link.String()
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

// TODO: how to return errors instead of panic inside of template execution?
func (d ReviewRequestData) AdminReviewURL() string {
	var (
		link *url.URL
		err  error
	)
	if d.BaseURL != "" {
		if link, err = url.Parse(d.BaseURL); err != nil {
			panic(err)
		}
	} else {
		panic("admin review base url missing")
	}

	return link.ResolveReference(&url.URL{Path: d.VID}).String()
}

// RejectRegistrationData to complete reject registration email templates.
type RejectRegistrationData struct {
	Name   string // Used to address the email
	VID    string // The ID of the VASP/Registration
	Reason string // A description of why the registration request was rejected
}

// DeliverCertsData to complete deliver certs email templates.
type DeliverCertsData struct {
	Name                string // Used to address the email
	VID                 string // The ID of the VASP/Registration
	CommonName          string // The common name assigned to the cert
	SerialNumber        string // The serial number of the certificate
	Endpoint            string // The expected endpoint for the TRISA service
	RegisteredDirectory string // The directory name for the certificates being issued
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
		log.Error().Err(err).Msg("could not attach JSON data to review request email")
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
	if data, err = ioutil.ReadFile(attachmentPath); err != nil {
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
