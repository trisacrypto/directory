package whisper_test

import (
	"context"
	"strings"
	"testing"
	"time"

	api "github.com/rotationalio/whisper/pkg/api/v1"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/utils/whisper"
)

func TestCreateWhisperLink(t *testing.T) {
	// CreateWhisperLink should throw an error when passed an empty secret
	accesses := 3
	oneWeek := time.Now().AddDate(0, 0, 7)
	link, err := whisper.CreateWhisperLink("", "not empty", accesses, oneWeek)
	require.Equal(t, link, "")
	require.EqualError(t, err, "a secret is required to generate a Whisper link")

	// CreateWhisperLink should throw an error when passed an empty password
	link, err = whisper.CreateWhisperLink("not empty", "", accesses, oneWeek)
	require.Equal(t, link, "")
	require.EqualError(t, err, "a password is required to generate a Whisper link")

	// CreateWhisperLink should throw an error when passed an expiration date in the past
	yesterday := time.Now().Add(-time.Hour * 24)
	link, err = whisper.CreateWhisperLink("not empty", "not empty", accesses, yesterday)
	require.Equal(t, link, "")
	require.EqualError(t, err, "the expiration date for the secret must be in the future")

	// Pass in valid arguments and check that the returned URL is valid
	link, err = whisper.CreateWhisperLink("this is a secret", "password", accesses, oneWeek)
	lastSlash := strings.LastIndex(link, "/")
	url := link[:lastSlash+1]
	token := link[lastSlash+1:]
	require.Equal(t, url, "https://api.whisper.rotational.dev/secret/")
	require.NotNil(t, token)

	// Create a Whisper client
	whisperClient, err := api.New("https://api.whisper.rotational.dev")
	require.NoError(t, err)

	// Make sure that the returned token can be used to fetch the secret for the number of
	// times that the 'accesses' argument was set
	for i := 1; i <= accesses; i++ {
		fetchReply, err := whisperClient.FetchSecret(context.Background(), token, "password")
		require.NoError(t, err, token)
		require.Equal(t, fetchReply.Secret, "this is a secret")
		if i == accesses {
			require.True(t, fetchReply.Destroyed)
		}
	}
}
