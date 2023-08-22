package emails

import (
	"fmt"
	"strings"
	"text/template"

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

// InviteUserData to complete user invite email templates.
type InviteUserData struct {
	UserName     string // The name of the user being invited
	UserEmail    string // The email address of the user being invited
	InviterName  string // The name of the user sending the invite
	InviterEmail string // The email address of the user sending the invite
	Organization string // The name of the relevant organization
	InviteURL    string // The invite URL that the recipient uses to accept the invite
}

func (d InviteUserData) Subject() string {
	if d.InviterName != "" {
		return fmt.Sprintf(UserInviteWithNameRE, d.InviterName, d.Organization)
	}
	return fmt.Sprintf(UserInviteRE, d.Organization)
}

//===========================================================================
// Email Builders
//===========================================================================

// InviteUserEmail creates a new user invite email, ready for sending by rendering the
// text and html templates with the supplied data then constructing a sendgrid email.
func InviteUserEmail(sender, senderEmail, recipient, recipientEmail string, data InviteUserData) (msg *mail.SGMailV3, err error) {
	var text, html string
	if text, html, err = Render("invite_user", data); err != nil {
		return nil, err
	}

	return mail.NewSingleEmail(
		mail.NewEmail(sender, senderEmail),
		data.Subject(),
		mail.NewEmail(recipient, recipientEmail),
		text,
		html,
	), nil
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
