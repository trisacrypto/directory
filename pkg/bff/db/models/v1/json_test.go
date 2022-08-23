package models_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff/db/models/v1"
	"google.golang.org/protobuf/proto"
)

// Test that the registration form marshals and unmarshals correctly to and from JSON
func TestMarshalRegistrationForm(t *testing.T) {
	// Empty form should be marshaled correctly and contain default values
	form := &models.RegistrationForm{}
	data, err := json.Marshal(form)
	require.NoError(t, err, "error marshaling empty registration form to JSON")
	require.NotEqual(t, []byte("{}"), data, "missing fields should be populated in marshaled JSON")

	// Empty form should be unmarshaled correctly
	require.NoError(t, json.Unmarshal(data, form), "error unmarshaling empty registration form from JSON")
	require.True(t, proto.Equal(&models.RegistrationForm{}, form), "empty registration form should be unmarshaled correctly")

	// Form with data should be marshaled correctly
	form = &models.RegistrationForm{
		Website: "https://alice.example.com",
	}
	data, err = json.Marshal(form)
	require.NoError(t, err, "error marshaling registration form to JSON")

	// Form with data should be unmarshaled correctly
	require.NoError(t, json.Unmarshal(data, form), "error unmarshaling registration form from JSON")
}
