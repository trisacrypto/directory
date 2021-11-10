package utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/utils"
)

func TestRewire(t *testing.T) {
	m := &models.CertificateRequest{
		Id:           "foo",
		CommonName:   "foo.example.com",
		CreationDate: "2021-04-08T12:42:33Z",
	}

	out, err := utils.Rewire(m)
	require.NoError(t, err, "could not rewire protobuf message")
	require.Contains(t, out, "id")
	require.Equal(t, m.Id, out["id"])
	require.Contains(t, out, "common_name")
	require.Equal(t, m.CommonName, out["common_name"])
	require.Contains(t, out, "creation_date")
	require.Equal(t, m.CreationDate, out["creation_date"])
}
