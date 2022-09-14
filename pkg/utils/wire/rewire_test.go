package wire_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/utils/wire"
)

func TestRewire(t *testing.T) {
	m := &models.CertificateRequest{
		Id:           "foo",
		CommonName:   "foo.example.com",
		CreationDate: "2021-04-08T12:42:33Z",
	}

	out, err := wire.Rewire(m)
	require.NoError(t, err, "could not rewire protobuf message")
	require.Contains(t, out, "id")
	require.Equal(t, m.Id, out["id"])
	require.Contains(t, out, "common_name")
	require.Equal(t, m.CommonName, out["common_name"])
	require.Contains(t, out, "creation_date")
	require.Equal(t, m.CreationDate, out["creation_date"])
}

func TestUnwire(t *testing.T) {
	m := map[string]interface{}{
		"id":            "foo",
		"common_name":   "foo.example.com",
		"creation_date": "2021-04-08T12:42:33Z",
	}

	cr := &models.CertificateRequest{}
	err := wire.Unwire(m, cr)
	require.NoError(t, err, "could not unwire json message")

	require.Equal(t, m["id"], cr.Id)
	require.Equal(t, m["common_name"], cr.CommonName)
	require.Equal(t, m["creation_date"], cr.CreationDate)
	require.Empty(t, cr.BatchName)
}
