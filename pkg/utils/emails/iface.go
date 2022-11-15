package emails

import (
	"github.com/sendgrid/rest"
	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
)

// EmailClient is an interface that can be implemented by SendGrid email clients.
type EmailClient interface {
	Send(email *sgmail.SGMailV3) (*rest.Response, error)
}
