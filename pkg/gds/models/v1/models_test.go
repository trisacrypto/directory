package models_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/global/v1"
	. "github.com/trisacrypto/directory/pkg/gds/models/v1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/proto"
)

func TestVASPExtra(t *testing.T) {
	// Ensures that the Get/Set methods on the VASP do not overwrite values other than
	// the values the method is intended to interact with.
	vasp := &pb.VASP{}

	// Attempt to get an admin verification token on a nil extra
	token, err := GetAdminVerificationToken(vasp)
	require.NoError(t, err)
	require.Equal(t, "", token)

	// Attempt to get metadata from a nil extra
	meta, deletedOn, err := GetMetadata(vasp)
	require.NoError(t, err)
	require.Empty(t, meta)
	require.True(t, deletedOn.IsZero())

	// Attempt to set admin verification token on a nil extra
	err = SetAdminVerificationToken(vasp, "pontoonboatz")
	require.NoError(t, err)

	// Should be able to fetch the token
	token, err = GetAdminVerificationToken(vasp)
	require.NoError(t, err)
	require.Equal(t, "pontoonboatz", token)

	// Attempt to set object metadata on non-nil extra
	meta = &global.Object{
		Key:       "foo",
		Namespace: "awesome",
		Version: &global.Version{
			Pid:     8,
			Version: 31,
			Region:  "Delaware",
		},
		Region: "Monte Cristo",
		Owner:  "The Count",
	}

	err = SetMetadata(vasp, meta, deletedOn)
	require.NoError(t, err)

	// Should be able to fetch metadata
	obj, ts, err := GetMetadata(vasp)
	require.NoError(t, err)
	require.True(t, proto.Equal(meta, obj))
	require.Equal(t, deletedOn, ts)

	// Check to ensure that admin verification token has not been overwritten
	token, err = GetAdminVerificationToken(vasp)
	require.NoError(t, err)
	require.Equal(t, "pontoonboatz", token)

	// Set a new admin verification token and check it doesn't override object meta
	err = SetAdminVerificationToken(vasp, "pontoonboat")
	require.NoError(t, err)

	// Should be able to fetch metadata
	obj, ts, err = GetMetadata(vasp)
	require.NoError(t, err)
	require.True(t, proto.Equal(meta, obj))
	require.Equal(t, deletedOn, ts)

	// Attempt to set metadata on a nil extra
	vasp.Extra = nil
	err = SetMetadata(vasp, meta, deletedOn)
	require.NoError(t, err)
}
