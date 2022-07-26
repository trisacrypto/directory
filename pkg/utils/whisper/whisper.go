package whisper

import (
	"context"
	"errors"
	"time"

	whisper "github.com/rotationalio/whisper/pkg/api/v1"
)

const accessDefault int = 3

// Takes in a PKCS12 secret and whisper password and returns the Whisper URL
// used to retrieve the secret.
func CreateWhisperLink(secret string, password string, accesses int, expirationDate time.Time) (link string, err error) {
	// Ensure the secret is not empty
	if secret == "" {
		return "", errors.New("a secret is required to generate a Whisper link")
	}

	// Ensure the password is not empty
	if password == "" {
		return "", errors.New("a password is required to generate a Whisper link")
	}

	// If accesses if not set, set it to the default
	if accesses <= 0 {
		accesses = accessDefault
	}

	// Ensure the expiration date is valid
	if expirationDate.Before(time.Now()) {
		return "", errors.New("the expiration date for the secret must be in the future")
	}

	// Create the whisper client
	var client whisper.Service
	endpoint := "https://api.whisper.rotational.dev"
	if client, err = whisper.New(endpoint); err != nil {
		return "", err
	}

	// Convert the expiration date into a whisper.Duration
	until := time.Until(expirationDate)
	lifetime := whisper.Duration(until)

	// Create the secret request
	request := &whisper.CreateSecretRequest{
		Secret:   secret,
		Password: password,
		Accesses: accesses,
		Lifetime: lifetime,
		IsBase64: false,
	}

	// Create a 30 second context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Call CreateSecretReply
	var reply *whisper.CreateSecretReply
	if reply, err = client.CreateSecret(ctx, request); err != nil {
		return "", err
	}

	// Create and return the Whisper link with the returned
	link = endpoint + "/secret/" + reply.Token
	return link, nil
}
