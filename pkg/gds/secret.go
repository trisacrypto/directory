package gds

import (
	"context"
	"math/rand"
	"strings"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"google.golang.org/api/iterator"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"google.golang.org/protobuf/types/known/durationpb"
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

// PingManager checks to make sure we can create a new client.
// This validates IAM permissions to some extent.
func PingManager(parent string) error {

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return err
	}
	client.Close()
	return nil
}

// CreateSecret creates a new secret in the Google Cloud Manager top-
// level directory, specified as `parent`, using the `secretID` provided
// as the name, to expire after `expiration` seconds.
// The parent should be a path, e.g.
//     "projects/project-name"
// This function returns a string representation of the path where the
// new secret is stored, e.g.
//     "projects/projectID/secrets/secretID"
// and an error if any occurs.
// Note: A secret is a logical wrapper around a collection of secret versions.
// Secret versions hold the actual secret material.
func CreateSecret(parent string, secretID string, expiration int64) (string, error) {

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.CreateSecretRequest{
		Parent:   parent,
		SecretId: secretID,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
			Expiration: &secretmanagerpb.Secret_Ttl{
				Ttl: &durationpb.Duration{
					Seconds: expiration,
				},
			},
		},
	}

	// Call the API.
	result, err := client.CreateSecret(ctx, req)
	if err != nil {
		return "", err
	}
	return result.Name, nil
}

// AddSecretVersion adds a new secret version to the given secret path with the
// provided payload. The path should be the full path to the secret, e.g.
//     "projects/projectID/secrets/secretID"
// Returns an error if one occurs.
func AddSecretVersion(path string, payload []byte) error {

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AddSecretVersionRequest{
		Parent: path,
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	}

	// Call the API.
	_, err = client.AddSecretVersion(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

// AccessSecretVersion returns the payload for the given secret version if one
// exists. The `version` is the full path to the secret version, and can be a
// version number as a string (e.g. "5") or an alias (e.g. "latest"), i.e.
//     "projects/projectID/secrets/secretID/versions/latest"
//     "projects/projectID/secrets/secretID/versions/5"
func AccessSecretVersion(version string) ([]byte, error) {

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: version,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, err
	}

	return result.Payload.Data, nil
}

// DeleteSecret deletes the secret with the given `name`, and all of its versions.
// `name` should be the root path to the secret, e.g.:
//     "projects/projectID/secrets/secretID"
// This is an irreversible operation. Any service or workload that attempts to
// access a deleted secret receives a Not Found error.
func DeleteSecret(name string) error {

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.DeleteSecretRequest{
		Name: name,
	}

	// Call the API.
	if err := client.DeleteSecret(ctx, req); err != nil {
		return err
	}
	return nil
}

// ListSecrets retrieves the names of all secrets in the project,
// given the `parent`, e.g.:
//     "projects/my-project"
// It returns a slice of strings representing the paths to the retrieved secrets,
// and a matching slice of errors for each failed retrieval.
func ListSecrets(parent string) (secrets []string, errors []error) {

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return secrets, append(errors, err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.ListSecretsRequest{
		Parent: parent,
	}

	// Call the API.
	it := client.ListSecrets(ctx, req)

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			errors = append(errors, err)
			secrets = append(secrets, "")
			continue
		}
		secrets = append(secrets, resp.Name)
		errors = append(errors, nil)
	}
	return secrets, errors
}
