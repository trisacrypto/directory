package whisper_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
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
	link, err := whisper.CreateSecretLink("", "not empty", accesses, oneWeek)
	require.Equal(t, link, "")
	require.EqualError(t, err, "a secret is required to generate a Whisper link")

	fixture := &api.CreateSecretReply{
		Token:   "abcdefghijklmnop",
		Expires: oneWeek,
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, r.Method, http.MethodPost)
		require.Equal(t, "/v1/secrets", r.URL.Path)

		in := new(api.CreateSecretRequest)
		err := json.NewDecoder(r.Body).Decode(in)
		require.NoError(t, err)

		require.Equal(t, in.Accesses, 3)
		require.NotNil(t, in.Lifetime)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer testServer.Close()

	err = whisper.ConnectClient(testServer.URL)
	defer whisper.ResetClient()
	require.NoError(t, err)

	link, err = whisper.CreateSecretLink("this is a secret", "password", accesses, oneWeek)
	require.NoError(t, err)
	require.Equal(t, link, "https://whisper.rotational.dev/secret/abcdefghijklmnop")
}

func TestCreateWhisperLinkLive(t *testing.T) {
	if os.Getenv("GDS_TEST_LIVE_WHISPER") != "1" {
		t.Skip()
	}
	// Pass in valid arguments and check that the returned URL is valid
	accesses := 3
	oneWeek := time.Now().AddDate(0, 0, 7)
	link, err := whisper.CreateSecretLink("this is a secret", "password", accesses, oneWeek)
	lastSlash := strings.LastIndex(link, "/")
	url := link[:lastSlash+1]
	token := link[lastSlash+1:]
	require.NoError(t, err)
	require.Equal(t, url, "https://whisper.rotational.dev/secret/")
	require.NotEmpty(t, token)

	// Create a Whisper client
	whisperClient, err := api.New("https://api.whisper.rotational.dev")
	defer whisper.ResetClient()
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
