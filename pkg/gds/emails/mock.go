package emails

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"time"

	"github.com/sendgrid/rest"
	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
)

var MockEmails [][]byte

func PurgeMockEmails() {
	MockEmails = nil
}

type mockSendGridClient struct {
	storage string
}

func WriteMIME(msg *sgmail.SGMailV3, path string) (err error) {
	type emailMetadata struct {
		From    string   `json:"from"`
		To      []string `json:"to"`
		Subject string   `json:"subject"`
	}

	// Create a buffer to store the MIME data
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Create the metadata header
	header := textproto.MIMEHeader{}
	header.Set("Content-Type", "application/json")
	part, err := writer.CreatePart(header)
	if err != nil {
		writer.Close()
		return err
	}

	// Construct the metadata header
	metadata := emailMetadata{
		From:    msg.From.Address,
		Subject: msg.Subject,
	}
	for _, p := range msg.Personalizations {
		for _, r := range p.To {
			metadata.To = append(metadata.To, r.Address)
		}
	}

	// Write the metadata header
	var b []byte
	if b, err = json.Marshal(metadata); err != nil {
		writer.Close()
		return err
	}
	if _, err = part.Write(b); err != nil {
		writer.Close()
		return err
	}

	// Write the email content sections
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

	// Write the attachment sections
	for _, a := range msg.Attachments {
		header := textproto.MIMEHeader{}
		header.Set("Content-Type", a.Type)
		header.Set("Content-Disposition", a.Disposition)
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

	// Save the file to disk
	writer.Close()
	if err = os.WriteFile(path, body.Bytes(), 0644); err != nil {
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
	MockEmails = append(MockEmails, data)

	if c.storage != "" {
		// Save the email to disk for manual inspection
		dir := filepath.Join(c.storage, toAddress)
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
		if err = WriteMIME(msg, path); err != nil {
			return &rest.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       fmt.Sprintf("could not archive email to %s", path),
			}, err
		}
	}

	return &rest.Response{StatusCode: http.StatusOK}, nil
}

func GetRecipient(msg *sgmail.SGMailV3) (recipient string, err error) {
	for _, p := range msg.Personalizations {
		for _, t := range p.To {
			recipient = t.Address
			return recipient, nil
		}
	}
	return "", errors.New("no recipient found for email")
}
