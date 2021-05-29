package gds

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/trisacrypto/directory/pkg/gds/config"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var chars = []rune("ABCDEFGHIJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz1234567890!#$%&'()*+,-./:;<=>?@[]^_`{|}~")

// CreateToken creates a variable length random token that can be used for passwords or API keys.
func CreateToken(length int) string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[random.Intn(len(chars))])
	}
	return b.String()
}

// SecretManager holds a client to the Google secret manager, and the path to the `parent` project for the secret manager.
type SecretManager struct {
	requestId string
	parent    string
	client    *secretmanager.Client
}

// NewSecretManager creates and returns a new secret manager client and an error if one occurs.
// Note that the `secretmanager` package leverages the GOOGLE_APPLICATION_CREDENTIALS
// environment variable which specifies the json path to the service account
// credentials, meaning that this function is a lightweight method for testing
// that the application can successfully connect to the secret manager API.
// However, this function does not validate the parent path.
func NewSecretManager(config config.Config, certRequest string) (sm *SecretManager, err error) {

	sm = &SecretManager{
		parent:    fmt.Sprintf("projects/%s", config.Secrets.Project),
		requestId: certRequest,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if sm.client, err = secretmanager.NewClient(ctx); err != nil {
		return nil, fmt.Errorf("could not connect to secret manager: %s", err)
	}

	return sm, nil
}

// CreateSecret creates a new secret in the Google Cloud Manager top-
// level directory using the `secret` name provided.
// This function returns an error if any occurs.
// Note: A secret is a logical wrapper around a collection of secret versions.
// To store a secret payload, you must first CreateSecret and then AddSecretVersion.
func (sm *SecretManager) CreateSecret(ctx context.Context, secret string) error {

	secretName := fmt.Sprintf("%s-%s", sm.requestId, secret)
	// Create an internal context, since a failed API call will result in infinite hang
	sctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Build the request.
	req := &secretmanagerpb.CreateSecretRequest{
		Parent:   sm.parent,
		SecretId: secretName,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}

	// Call the API. Note: We don't actually need the result that comes back from the API call
	// and not accessing it directly (e.g. logging plaintext, etc) provides added security
	_, err := sm.client.CreateSecret(sctx, req)
	if err != nil {
		// If the API call is malformed, it will hang until the internal context times out
		if errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		// If the secret already exists, that means the client already has a password set up
		// This is fine because we'll just create a new version with a new password for them
		// and CertMan will always look for the most recent secret version.
		serr, ok := status.FromError(err)
		if ok && serr.Code() == codes.AlreadyExists {
			return err
		}

		// If the error is something else, something went wrong.
		return err
	}
	return nil
}

// AddSecretVersion adds a new secret version to the given secret and the
// provided payload. Returns an error if one occurs.
// Note: to add a secret version, the secret must first be created using CreateSecret.
func (sm *SecretManager) AddSecretVersion(ctx context.Context, secret string, payload []byte) error {

	secretPath := fmt.Sprintf("%s/secrets/%s-%s", sm.parent, sm.requestId, secret)
	// Create an internal context, since a failed API call will result in infinite hang
	sctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Build the request.
	req := &secretmanagerpb.AddSecretVersionRequest{
		Parent: secretPath,
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	}

	// Call the API. Note: We don't actually need the result that comes back from the API call
	// and not accessing it directly (e.g. logging plaintext, etc) provides added security
	_, err := sm.client.AddSecretVersion(sctx, req)
	if err != nil {
		// If the API call is malformed, it will hang until the internal context times out
		if errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		// If the secret does not exist (e.g. has been deleted or hasn't been created yet)
		// we'll get a Not Found error
		serr, ok := status.FromError(err)
		if ok && serr.Code() == codes.NotFound {
			return err
		}

		// If the error is something else, something went wrong.
		return errors.New("unknown error: unable to create secret version")
	}

	return nil
}

// GetLatestVersion returns the payload for the latest version of the given secret,
// if one exists, else an error.
func (sm *SecretManager) GetLatestVersion(ctx context.Context, secret string) ([]byte, error) {

	versionPath := fmt.Sprintf("%s/secrets/%s-%s/versions/latest", sm.parent, sm.requestId, secret)

	// Create an internal context, since a failed API call will result in infinite hang
	sctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: versionPath,
	}

	// Call the API.
	result, err := sm.client.AccessSecretVersion(sctx, req)
	if err != nil {
		// If the API call is malformed, it will hang until the internal context times out
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}

		// If the error is something else, something went wrong.
		return nil, errors.New("unknown error: unable to access latest secret version")
	}

	return result.Payload.Data, nil
}

// DeleteSecret deletes the secret with the given the name, and all of its versions.
// Note: this is an irreversible operation. Any service or workload that attempts to
// access a deleted secret receives a Not Found error.
func (sm *SecretManager) DeleteSecret(ctx context.Context, secret string) error {

	secretPath := fmt.Sprintf("%s/secrets/%s-%s", sm.parent, sm.requestId, secret)

	// Create an internal context, since a failed API call will result in infinite hang
	sctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Build the request.
	req := &secretmanagerpb.DeleteSecretRequest{
		Name: secretPath,
	}

	// Call the API.
	err := sm.client.DeleteSecret(sctx, req)
	if err != nil {
		// If the API call is malformed, it will hang until the internal context times out
		if errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		// If the error is something else, something went wrong.
		return errors.New("unknown error: unable to delete secret")
	}
	return nil
}
