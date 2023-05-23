package models_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/directory/pkg/bff/models/v1"
)

func TestValidationErrors(t *testing.T) {
	err := &ValidationError{Field: "testing", Err: "something went wrong", Index: 42}
	require.EqualError(t, err, "invalid field testing: something went wrong")

	var ok bool
	verrs := make(ValidationErrors, 0)
	verrs, ok = verrs.Append(err)
	require.True(t, ok)
	require.EqualError(t, verrs, "1 validation errors occurred:\ninvalid field testing: something went wrong")

	// Should be able to append nil
	uerrs, ok := verrs.Append(nil)
	require.True(t, ok)
	require.Equal(t, verrs, uerrs)

	// Should be able to append multiple errors
	uerrs = ValidationErrors{
		{Field: "name", Err: "name is required"},
		{Field: "age", Err: "age cannot be greater than 150"},
		{Field: "colors", Err: "could not parse hex color", Index: 3},
	}

	verrs, ok = verrs.Append(uerrs)
	require.True(t, ok)
	require.Len(t, verrs, 4)

	uerrs, ok = verrs.Append(ErrCollaboratorExists)
	require.False(t, ok)
	require.Equal(t, verrs, uerrs)
}
