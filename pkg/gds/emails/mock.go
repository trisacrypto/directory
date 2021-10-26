package emails

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"

	"github.com/sendgrid/rest"
	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
)

var MockEmails [][]byte

func PurgeMockEmails() {
	MockEmails = nil
}

type mockSendGridClient struct{}

type emailMetadata struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
}

func write(msg *sgmail.SGMailV3) (err error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	header := textproto.MIMEHeader{}
	header.Set("Content-Type", "application/json")
	part, err := writer.CreatePart(header)
	if err != nil {
		writer.Close()
		return err
	}

	metadata := emailMetadata{
		From:    msg.From.Address,
		Subject: msg.Subject,
	}
	for _, p := range msg.Personalizations {
		for _, r := range p.To {
			metadata.To = append(metadata.To, r.Address)
		}
	}
	var b []byte
	if b, err = json.Marshal(metadata); err != nil {
		writer.Close()
		return err
	}
	if _, err = part.Write(b); err != nil {
		writer.Close()
		return err
	}

	for _, c := range msg.Content {
		header := textproto.MIMEHeader{}
		header.Set("Content-Type", c.Type)
		part, err := writer.CreatePart(header)
		if err != nil {
			writer.Close()
			return err
		}
		if _, err = part.Write([]byte(c.Value)); err != nil {
			writer.Close()
			return err
		}
	}

	for _, a := range msg.Attachments {
		header := textproto.MIMEHeader{}
		header.Set("Content-Type", a.Type)
		header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", a.Filename))
		part, err := writer.CreatePart(header)
		if err != nil {
			writer.Close()
			return err
		}
		if _, err = part.Write([]byte(a.Content)); err != nil {
			writer.Close()
			return err
		}
	}

	writer.Close()
	if err = os.WriteFile("test.txt", body.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

func (c *mockSendGridClient) Send(msg *sgmail.SGMailV3) (rep *rest.Response, err error) {
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

	// Email recepients are stored in Personalizations
	if len(msg.Personalizations) == 0 {
		return &rest.Response{
			StatusCode: http.StatusBadRequest,
			Body:       "no Personalization info",
		}, errors.New("requires Personalization info")
	}

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

			if _, err := sgmail.ParseEmail(t.Address); err != nil {
				return &rest.Response{
					StatusCode: http.StatusBadRequest,
					Body:       fmt.Sprintf("invalid To address: %s", t.Address),
				}, err
			}
		}
	}

	// "Send" the email
	MockEmails = append(MockEmails, data)
	if err = write(msg); err != nil {
		return &rest.Response{
			StatusCode: http.StatusBadRequest,
			Body:       "could not write email",
		}, err
	}

	return &rest.Response{StatusCode: http.StatusOK}, nil
}
