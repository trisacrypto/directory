package mock

import (
	"errors"
	"fmt"
	"hash/fnv"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/sendgrid/rest"
	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/trisacrypto/directory/pkg/utils/emails"
)

var Emails [][]byte

func PurgeEmails() {
	Emails = nil
}

type SendGridClient struct {
	Storage string
}

func (c *SendGridClient) Send(msg *sgmail.SGMailV3) (rep *rest.Response, err error) {
	// Marshal the email struct into bytes
	data := sgmail.GetRequestBody(msg)
	if data == nil {
		return &rest.Response{
			StatusCode: http.StatusBadRequest,
			Body:       "invalid email data",
		}, errors.New("could not marshal email")
	}

	// Email needs to contain a From address
	if msg.From.Address == "" {
		return &rest.Response{
			StatusCode: http.StatusBadRequest,
			Body:       "no From address",
		}, errors.New("requires From address")
	}

	// Validate From address
	if _, err := sgmail.ParseEmail(msg.From.Address); err != nil {
		return &rest.Response{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("invalid From address: %s", msg.From.Address),
		}, err
	}

	// Email recipients are stored in Personalizations
	if len(msg.Personalizations) == 0 {
		return &rest.Response{
			StatusCode: http.StatusBadRequest,
			Body:       "no Personalization info",
		}, errors.New("requires Personalization info")
	}

	var toAddress string
	for _, p := range msg.Personalizations {
		// Email needs to contain at least one To address
		if len(p.To) == 0 {
			return &rest.Response{
				StatusCode: http.StatusBadRequest,
				Body:       "no To addresses",
			}, errors.New("requires To address")
		}

		for _, t := range p.To {
			// Validate To address
			if t.Address == "" {
				return &rest.Response{
					StatusCode: http.StatusBadRequest,
					Body:       "no To address",
				}, errors.New("empty To address")
			}

			var mail *sgmail.Email
			if mail, err = sgmail.ParseEmail(t.Address); err != nil {
				return &rest.Response{
					StatusCode: http.StatusBadRequest,
					Body:       fmt.Sprintf("invalid To address: %s", t.Address),
				}, err
			}
			toAddress = mail.Address
		}
	}

	// "Send" the email
	Emails = append(Emails, data)

	if c.Storage != "" {
		// Save the email to disk for manual inspection
		dir := filepath.Join(c.Storage, toAddress)
		if err = os.MkdirAll(dir, 0755); err != nil {
			return &rest.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       fmt.Sprintf("could not create archive directory at %s", dir),
			}, err
		}

		// Generate unique filename to avoid overwriting
		ts := time.Now().Format(time.RFC3339)
		h := fnv.New32()
		h.Write(data)
		path := filepath.Join(dir, fmt.Sprintf("%s-%d.mim", ts, h.Sum32()))
		if err = emails.WriteMIME(msg, path); err != nil {
			return &rest.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       fmt.Sprintf("could not archive email to %s", path),
			}, err
		}
	}

	return &rest.Response{StatusCode: http.StatusOK}, nil
}
