package trisads

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

// VerifyContactEmail creates a verification token for each contact in the VASP contact
// list and sends them the verification email with instructions on how to verify their
// email address.
func (s *Server) VerifyContactEmail(vasp *pb.VASP) (err error) {
	// Create the verification tokens and save the VASP back to the database
	var contacts = []*pb.Contact{
		vasp.Contacts.Technical, vasp.Contacts.Administrative,
		vasp.Contacts.Billing, vasp.Contacts.Legal,
	}

	for _, contact := range contacts {
		if contact != nil && contact.Email != "" {
			contact.Token = CreateToken(48)
			contact.Verified = false
		}
	}

	if err = s.db.Update(vasp); err != nil {
		log.Error().Msg("could not update vasp")
		return err
	}

	for _, contact := range contacts {
		if contact == nil || contact.Email == "" {
			continue
		}

		ctx := verifyContactContext{
			Name:  contact.Name,
			Token: contact.Token,
			VID:   vasp.Id,
		}

		var text, html string
		if text, err = execTemplateString(verifyContactPlainText, ctx); err != nil {
			return err
		}
		if html, err = execTemplateString(verifyContactHTML, ctx); err != nil {
			return err
		}

		if err = s.sendEmail(contact.Name, contact.Email, verifyContactSubject, text, html); err != nil {
			return err
		}

	}

	return nil
}

// ReviewRequestEmail is a shortcut for iComply verification in which we simply send
// an email to the TRISA admins and have them manually verify registrations.
func (s *Server) ReviewRequestEmail(vasp *pb.VASP) (err error) {
	// Create verification token for admin and update database
	// TODO: replace with actual authentication
	vasp.AdminVerificationToken = CreateToken(48)
	if err = s.db.Update(vasp); err != nil {
		return fmt.Errorf("could not save admin verification token: %s", err)
	}

	var data []byte
	if data, err = json.MarshalIndent(vasp, "", "  "); err != nil {
		return err
	}

	ctx := reviewRequestContext{
		VID:     vasp.Id,
		Request: string(data),
		Token:   vasp.AdminVerificationToken,
	}

	var text, html string
	if text, err = execTemplateString(reviewRequestPlainText, ctx); err != nil {
		return err
	}
	if html, err = execTemplateString(reviewRequestHTML, ctx); err != nil {
		return err
	}

	if err = s.sendEmail("TRISA Admins", s.conf.AdminEmail, reviewRequestSubject, text, html); err != nil {
		return err
	}

	return nil
}

// RejectRegistrationEmail sends a notification to all VASP contacts that their
// registration status is rejected without certificate issuance and explains why.
func (s *Server) RejectRegistrationEmail(vasp *pb.VASP, reason string) (err error) {
	ctx := &rejectRegistrationContext{
		VASP:   vasp.Id,
		Reason: reason,
	}

	var contacts = []*pb.Contact{
		vasp.Contacts.Technical, vasp.Contacts.Administrative,
		vasp.Contacts.Billing, vasp.Contacts.Legal,
	}

	for _, contact := range contacts {
		if contact != nil && contact.Verified {
			ctx.Name = contact.Name
			var text, html string
			if text, err = execTemplateString(rejectRegistrationPlainText, ctx); err != nil {
				return err
			}
			if html, err = execTemplateString(rejectRegistrationHTML, ctx); err != nil {
				return err
			}

			if err = s.sendEmail(contact.Name, contact.Email, rejectRegistrationSubject, text, html); err != nil {
				return err
			}
		}
	}

	return nil
}

// DeliverCertificatesEmail sends the PKCS12 encrypted certificate files to the VASP
// contacts as an attachment, completing the certificate issuance process.
func (s *Server) DeliverCertificatesEmail(vasp *pb.VASP, path string) (err error) {
	ctx := &deliverCertsContext{
		VASP:         vasp.Id,
		CommonName:   vasp.CommonName,
		SerialNumber: hex.EncodeToString(vasp.IdentityCertificate.SerialNumber),
		Endpoint:     vasp.TrisaEndpoint,
	}

	var contacts = []*pb.Contact{
		vasp.Contacts.Technical, vasp.Contacts.Administrative,
		vasp.Contacts.Billing, vasp.Contacts.Legal,
	}

	for _, contact := range contacts {
		if contact != nil && contact.Verified {
			ctx.Name = contact.Name
			var text, html string
			if text, err = execTemplateString(deliverCertsPlainText, ctx); err != nil {
				return err
			}
			if html, err = execTemplateString(deliverCertsHTML, ctx); err != nil {
				return err
			}

			if err = s.sendEmailAttachment(contact.Name, contact.Email, deliverCertsSubject, text, html, path); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Server) sendEmail(recipient, emailaddr, subject, text, html string) (err error) {
	message := mail.NewSingleEmail(
		mail.NewEmail("TRISA Directory Service", s.conf.ServiceEmail),
		subject,
		mail.NewEmail(recipient, emailaddr),
		text, html,
	)

	var rep *rest.Response
	if rep, err = s.email.Send(message); err != nil {
		return err
	}

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		return errors.New(rep.Body)
	}

	return nil
}

func (s *Server) sendEmailAttachment(recipient, emailaddr, subject, text, html, path string) (err error) {
	message := mail.NewSingleEmail(
		mail.NewEmail("TRISA Directory Service", s.conf.ServiceEmail),
		subject,
		mail.NewEmail(recipient, emailaddr),
		text, html,
	)

	// Read and encode the attachment data
	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		return err
	}
	encoded := base64.StdEncoding.EncodeToString(data)

	// Create the attachment
	// TODO: detect mimetype rather than assuming zip
	attach := mail.NewAttachment()
	attach.SetContent(encoded)
	attach.SetType("application/zip")
	attach.SetFilename(filepath.Base(path))
	attach.SetDisposition("attachment")
	message.AddAttachment(attach)

	var rep *rest.Response
	if rep, err = s.email.Send(message); err != nil {
		return err
	}

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		return errors.New(rep.Body)
	}

	return nil
}

func execTemplateString(t *template.Template, ctx interface{}) (_ string, err error) {
	buf := new(strings.Builder)
	if err = t.Execute(buf, ctx); err != nil {
		return "", err
	}
	return buf.String(), nil
}

type verifyContactContext struct {
	Name  string
	Token string
	VID   string
}

var verifyContactSubject = "Verify Email Address"

// VerifyContact Plain Text Content Template
var verifyContactPlainText = template.Must(template.New("verifyContactPlainText").Parse(`
Hello {{ .Name }},

Thank you for submitting a TRISA TestNet VASP registration request. To begin the
verification process, please submit the following email verification token using the
VerifyEmail RPC in the TRISA directory service protocol:

ID: {{ .VID }}
Token: {{ .Token }}

This can be done with the trisads CLI utility or using the protocol buffers library in
the programming language of your choice.

Note that we're working to create a URL endpoint for the vaspdirectory.net site to
simplify the verification process. We're sorry about the inconvenience of this method at
the early stage of the TRISA Test Net.

Best Regards,
The TRISA Directory Service`))

// VerifyContact HTML Content Template
var verifyContactHTML = template.Must(template.New("verifyContactHTML").Parse(`
<p>Hello {{ .Name }},</p>

<p>Thank you for submitting a TRISA TestNet VASP registration request. To begin the
verification process, please submit the following email verification token using the
VerifyEmail RPC in the TRISA directory service protocol:</p>

<ul>
	<li>ID: <strong>{{ .VID }}</strong></li>
	<li>Token: <strong>{{ .Token }}</strong></li>
</ul>

<p>This can be done with the trisads CLI utility or using the protocol buffers library in
the programming language of your choice.</p>

<p>Note that we're working to create a URL endpoint for the
<a href="https://vaspdirectory.net/">vaspdirectory.net</a> site to simplify the
verification process. We're sorry about the inconvenience of this method at the early
stage of the TRISA Test Net.</p>

<p>Best Regards,<br />
The TRISA Directory Service</p>`))

type reviewRequestContext struct {
	VID     string
	Token   string
	Request string
}

var reviewRequestSubject = "Please Review TRISA TestNet Directory Registration Request"

// VerifyContact Plain Text Content Template
var reviewRequestPlainText = template.Must(template.New("reviewRequestPlainText").Parse(`
Hello TRISA Admin,

We have received a new registration request from a VASP that needs to be reviewed. The
requestor has verified their email address and received a PKCS12 password to decrypt a
certificate that will be generated if you approve this request. The request JSON is:

{{ .Request }}

To verify or reject the registration request, use the following metadata with the
trisads verify command:

ID: {{ .VID }}
Token: {{ .Token }}

Note that we're working to create a URL endpoint for the vaspdirectory.net site to
simplify the verification process. We're sorry about the inconvenience of this method at
the early stage of the TRISA Test Net.

Best Regards,
The TRISA Directory Service`))

// VerifyContact HTML Content Template
var reviewRequestHTML = template.Must(template.New("reviewRequestHTML").Parse(`
<p>Hello TRISA Admin,</p>

<p>We have received a new registration request from a VASP that needs to be reviewed.
The requestor has verified their email address and received a PKCS12 password to decrypt
a certificate that will be generated if you approve this request. The request JSON is:</p>

<pre>{{ .Request }}</pre>

<p>To verify or reject the registration request, use the following metadata with the
<code>trisads verify</code> command:</p>

<ul>
	<li>ID: <strong>{{ .VID }}</strong></li>
	<li>Token: <strong>{{ .Token }}</strong></li>
</ul>

<p>Note that we're working to create a URL endpoint for the
<a href="https://vaspdirectory.net/">vaspdirectory.net</a> site to simplify the
verification process. We're sorry about the inconvenience of this method at the early
stage of the TRISA TestNet.</p>

<p>Best Regards,<br />
The TRISA Directory Service</p>`))

type rejectRegistrationContext struct {
	Name   string
	VASP   string
	Reason string
}

var rejectRegistrationSubject = "TRISA TestNet Directory Registration Status"

// VerifyContact Plain Text Content Template
var rejectRegistrationPlainText = template.Must(template.New("rejectRegistrationPlainText").Parse(`
Hello {{ .Name }},

Unfortunately we have had to reject your TRISA TestNet Directory Service registration
request. The reason the request was rejected is as follows:

ID: {{ .VASP }}

{{ .Reason }}

Please note that in order to fix the problems mentioned above, you'll have to submit a
registration appeal. Please send any questions to admin@trisa.io.

Best Regards,
The TRISA Directory Service`))

// VerifyContact HTML Content Template
var rejectRegistrationHTML = template.Must(template.New("rejectRegistrationHTML").Parse(`
<p>Hello {{ .Name }},</p>

<p>Unfortunately we have had to reject your TRISA TestNet Directory Service registration
request (<strong>{{ .VASP}}</strong>). The reason the request was rejected is as
follows:</p>

<p>{{ .Reason }}</p>

<p>Please note that in order to fix the problems mentioned above, you'll have to submit
a registration appeal. Please send any questions to
<a href="mailto:admin@trisa.io">admin@trisa.io</a>.</p>

<p>Best Regards,<br />
The TRISA Directory Service</p>`))

type deliverCertsContext struct {
	Name         string
	VASP         string
	CommonName   string
	SerialNumber string
	Endpoint     string
}

var deliverCertsSubject = "TRISA TestNet Directory Registration Status"

// VerifyContact Plain Text Content Template
var deliverCertsPlainText = template.Must(template.New("deliverCertsPlainText").Parse(`
Hello {{ .Name }},

Your TRISA TestNet registration has been approved! You've been validated as a member of
the TRISA Test Network as an integration partner. Attached to this email are your
PKCS12 encrypted certificates so that you can implement the TRISA P2P protocol using
mTLS with other network integration partners.

The primary details of your directory entry are as follows:

ID: {{ .VASP }}
Common Name: {{ .CommonName }}
Serial Number: {{ .SerialNumber }}
Endpoint: {{ .Endpoint }}

To decrypt your certificates, you will need the PKCS12 password that you received when
you verified your email. Note that the first contact to verify their email received the
password. To decrypt the certificates on the command line, you can use openssl as
follows:

openssl pkcs12 -in INFILE.p12 -out OUTFILE.crt -nodes

For more information on integrating with the TestNet, please see our documentation at
https://trisatest.net/. If you have any questions, you may contact us at admin@trisa.io.

Best Regards,
The TRISA Directory Service`))

// VerifyContact HTML Content Template
var deliverCertsHTML = template.Must(template.New("deliverCertsHTML").Parse(`
<p>Hello {{ .Name }},</p>

<p>Your TRISA TestNet registration has been approved! You've been validated as a member
of the TRISA Test Network as an <em>integration partner</em>. Attached to this email are
your PKCS12 encrypted certificates so that you can implement the TRISA P2P protocol
using mTLS with other network integration partners.</p>

<p>The primary details of your directory entry are as follows:</p>

<ul>
	<li><strong>ID:</strong> {{ .VASP }}</li>
	<li><strong>Common Name:</strong> {{ .CommonName }}</li>
	<li><strong>Serial Number:</strong> {{ .SerialNumber }}</li>
	<li><strong>Endpoint:</strong> {{ .Endpoint }}</li>
</ul>

To decrypt your certificates, you will need the PKCS12 password that you received when
you verified your email. Note that <strong>the first contact to verify their email
received the password</strong>. To decrypt the certificates on the command line, you can
use <code>openssl</code> as follows:

<pre>$ openssl pkcs12 -in INFILE.p12 -out OUTFILE.crt -nodes</pre>

<p>For more information on integrating with the TestNet, please see our documentation at
<a href="https://trisatest.net/">trisatest.net</a>. If you have any questions, you may
contact us at <a href="mailto:admin@trisa.io">admin@trisa.io</a>.</p>

<p>Best Regards,<br />
The TRISA Directory Service</p>`))
