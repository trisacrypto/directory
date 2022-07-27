package whisper

import (
	"context"
	"errors"
	"sync"
	"time"

	whisper "github.com/rotationalio/whisper/pkg/api/v1"
)

const (
	accessDefault int    = 3
	endpoint      string = "https://api.whisper.rotational.dev"
)

var (
	client     whisper.Service
	initError  error
	initClient sync.Once
)

// Creates a Whisper secret and returns the link to the Whisper UI where a user can access the secret value. If the
// password is specified, then the user will need to enter the password to access the secret. The accesses and
// expiration date options limit how long the secret is available. By default these values are 3 accesses and 7 days
// from today. An error is returned if the secret cannot be created with Whisper.
func CreateSecretLink(secret string, password string, accesses int, expirationDate time.Time) (link string, err error) {
	// Ensure the secret is not empty
	if secret == "" {
		return "", errors.New("a secret is required to generate a Whisper link")
	}

	// If accesses if not set, set it to the default
	if accesses <= 0 {
		accesses = accessDefault
	}

	// Ensure the expiration date is valid
	if expirationDate.IsZero() {
		expirationDate = time.Now().AddDate(0, 0, 7)
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

	if err = ConnectClient(endpoint); err != nil {
		return "", err
	}

	// Call CreateSecretReply
	var reply *whisper.CreateSecretReply
	if reply, err = client.CreateSecret(ctx, request); err != nil {
		return "", err
	}

	// Create and return the Whisper link with the returned
	link = "https://whisper.rotational.dev/secret/" + reply.Token
	return link, nil
}

// Connects the client to the Whisper service, throws an
// error if the client is already connected.
func ConnectClient(endpoint string) error {
	// Create the whisper client, using sync.Once to ensure we only instantiate
	// a client once
	initClient.Do(func() {
		client, initError = whisper.New(endpoint)
	})
	if initError != nil {
		return initError
	}
	return nil
}

// Accessor for the client used in testing
func Client() whisper.Service {
	return client
}

// Resets the client connection
func ResetClient() {
	initClient = *new(sync.Once)
}
