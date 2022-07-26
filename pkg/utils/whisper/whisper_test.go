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
	accesses := 3
	oneWeek := time.Now().AddDate(0, 0, 7)
	link, err := whisper.CreateWhisperLink("", "not empty", accesses, oneWeek)
	require.Equal(t, link, "")
	require.EqualError(t, err, "a secret is required to generate a Whisper link")

	link, err = whisper.CreateWhisperLink("not empty", "", accesses, oneWeek)
	require.Equal(t, link, "")
	require.EqualError(t, err, "a password is required to generate a Whisper link")

	link, err = whisper.CreateWhisperLink("this is a secret", "password", accesses, oneWeek)
	slashIndex := strings.LastIndex(link, "/")
	url := link[:slashIndex+1]
	token := link[slashIndex+1:]
	require.Equal(t, url, "https://api.whisper.rotational.dev/secret/")
	require.NotNil(t, token)

	whisperClient, err := api.New("https://api.whisper.rotational.dev")
	require.NoError(t, err)

	for i := 1; i <= 4; i++ {
		fetchReply, err := whisperClient.FetchSecret(context.Background(), token, "password")
		require.NoError(t, err)
		require.Equal(t, fetchReply.Secret, "this is a secret")
		require.Equal(t, fetchReply.Accesses, i)
		if i == accesses {
			require.True(t, fetchReply.Destroyed)
		}
	}
}
