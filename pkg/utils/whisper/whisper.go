package whisper

import (
	"context"
	"errors"
	"time"

	whisper "github.com/rotationalio/whisper/pkg/api/v1"
)

// Takes in a PKCS12 secret and
func CreateWhisperLink(secret string, password string, accesses int, expirationDate time.Time) (link string, err error) {
	if secret == "" {
		return "", errors.New("a secret is required to generate a Whisper link")
	}

	if password == "" {
		return "", errors.New("a password is required to generate a Whisper link")
	}
	if accesses <= 0 {
		accesses = 3
	}

	var client whisper.Service
	endpoint := "https://api.whisper.rotational.dev"
	if client, err = whisper.New(endpoint); err != nil {
		return "", err
	}

	until := time.Until(expirationDate)
	lifetime := whisper.Duration(until)

	request := &whisper.CreateSecretRequest{
		Secret:   secret,
		Password: password,
		Accesses: accesses,
		Lifetime: lifetime,
		IsBase64: false,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var reply *whisper.CreateSecretReply
	if reply, err = client.CreateSecret(ctx, request); err != nil {
		return "", err
	}

	link = endpoint + "/secret/" + reply.Token
	return link, nil
}
