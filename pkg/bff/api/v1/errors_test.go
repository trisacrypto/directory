package api_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/directory/pkg/bff/api/v1"
)

type IntString int

func (t IntString) String() string {
	return fmt.Sprintf("%04x", int(t))
}

type JMap map[string]string

func (j JMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string(j))
}

func TestErrorResponse(t *testing.T) {
	testCases := []struct {
		err      interface{}
		expected string
	}{
		{fmt.Errorf("something %s happened", "bad"), "something bad happened"},
		{errors.New("godzilla is here!"), "godzilla is here!"},
		{"this is a simple string", "this is a simple string"},
		{IntString(42), "002a"},
		{JMap{"color": "red"}, "{\"color\":\"red\"}"},
		{42, "unhandled error response"},
	}

	for _, tc := range testCases {
		rep := ErrorResponse(tc.err)
		require.False(t, rep.Success, "expected error reply to be success false")
		require.Equal(t, tc.expected, rep.Error, "unexpected result")
	}
}

func TestNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	NotFound(ctx)

	result := r.Result()
	defer result.Body.Close()
	require.Equal(t, result.StatusCode, http.StatusNotFound)
	require.Equal(t, "application/json; charset=utf-8", result.Header.Get("Content-Type"))

	var data map[string]interface{}
	err := json.NewDecoder(result.Body).Decode(&data)
	require.NoError(t, err)

}

func TestNotAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	NotAllowed(ctx)

	result := r.Result()
	defer result.Body.Close()
	require.Equal(t, result.StatusCode, http.StatusMethodNotAllowed)
	require.Equal(t, "application/json; charset=utf-8", result.Header.Get("Content-Type"))

	var data map[string]interface{}
	err := json.NewDecoder(result.Body).Decode(&data)
	require.NoError(t, err)
}
