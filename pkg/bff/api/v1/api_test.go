package api_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
)

func TestFieldValidationError(t *testing.T) {
	verr := &models.ValidationError{
		Field: "name",
		Err:   "name is required",
		Index: 42,
	}

	ferr := api.NewFieldValidationError(verr)
	require.Equal(t, verr.Field, ferr.Field)
	require.Equal(t, verr.Err, ferr.Error)
	require.Equal(t, verr.Index, ferr.Index)

	ferrs := api.FromValidationErrors(verr)
	require.Len(t, ferrs, 1)
	require.Equal(t, ferr, ferrs[0])

	err := errors.New("something went wrong")
	ferr = api.NewFieldValidationError(err)
	require.Empty(t, ferr.Field)
	require.Equal(t, err.Error(), ferr.Error)
	require.Zero(t, ferr.Index)

	ferrs = api.FromValidationErrors(err)
	require.Len(t, ferrs, 1)
	require.Equal(t, ferr, ferrs[0])
}

func TestFromValidationErrors(t *testing.T) {
	verr := make(models.ValidationErrors, 0, 3)
	verr, _ = verr.Append(&models.ValidationError{Field: "name", Err: "name is required"})
	verr, _ = verr.Append(&models.ValidationError{Field: "colors", Err: "invalid hex color", Index: 3})
	verr, _ = verr.Append(&models.ValidationError{Field: "colors", Err: "invalid rgb color", Index: 5})

	ferrs := api.FromValidationErrors(verr)
	require.Len(t, ferrs, len(verr))
	for i, ferr := range ferrs {
		require.Equal(t, verr[i].Field, ferr.Field)
		require.Equal(t, verr[i].Err, ferr.Error)
		require.Equal(t, verr[i].Index, ferr.Index)
	}
}
