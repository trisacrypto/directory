package wire_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/global/v1"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/peers/v1"
	. "github.com/trisacrypto/directory/pkg/gds/store/wire"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

func TestWire(t *testing.T) {
	// Without base64 encoding of keys
	key, err := DecodeKey("foo", false)
	require.NoError(t, err)
	require.Equal(t, "foo", EncodeKey(key, false))

	// With base64 encoding of keys
	key, err = DecodeKey("foo", true)
	require.NoError(t, err)
	require.Equal(t, "foo", EncodeKey(key, true))

	// The basic test pattern here is to load data from disk, then remarshal JSON
	// and testing unmarshal of the specific namespace.
	// Test unmarshal VASP records
	in, err := ioutil.ReadFile("testdata/vasps::838b1f57-1646-488d-a231-d71d88681cfa.json")
	require.NoError(t, err)
	out, err := RemarshalJSON(global.NamespaceVASPs, in)
	require.NoError(t, err)
	vasp, err := UnmarshalProto(global.NamespaceVASPs, out)
	require.NoError(t, err)
	_, ok := vasp.(*pb.VASP)
	require.True(t, ok)
	obj, err := UnmarshalObject(global.NamespaceVASPs, out, false)
	require.NoError(t, err)
	require.Equal(t, "vasps::838b1f57-1646-488d-a231-d71d88681cfa", obj.Key)
	require.Equal(t, uint64(8), obj.Version.Version)
	require.Empty(t, obj.Data)
	obj, err = UnmarshalObject(global.NamespaceVASPs, out, true)
	require.NoError(t, err)
	require.NotEmpty(t, obj.Data)

	// Test unmarshal certificate requests
	in, err = ioutil.ReadFile("testdata/certreqs::87657b4d-e72c-4526-9332-c8fc56adb367.json")
	require.NoError(t, err)
	out, err = RemarshalJSON(global.NamespaceCertReqs, in)
	require.NoError(t, err)
	certreq, err := UnmarshalProto(global.NamespaceCertReqs, out)
	require.NoError(t, err)
	_, ok = certreq.(*models.CertificateRequest)
	require.True(t, ok)
	obj, err = UnmarshalObject(global.NamespaceCertReqs, out, false)
	require.NoError(t, err)
	require.Equal(t, "certreqs::87657b4d-e72c-4526-9332-c8fc56adb367", obj.Key)
	require.Equal(t, uint64(3), obj.Version.Version)
	require.Empty(t, obj.Data)
	obj, err = UnmarshalObject(global.NamespaceCertReqs, out, true)
	require.NoError(t, err)
	require.NotEmpty(t, obj.Data)

	// Test Unmarshal Peers
	in, err = ioutil.ReadFile("testdata/peers::8.json")
	require.NoError(t, err)
	out, err = RemarshalJSON(global.NamespaceReplicas, in)
	require.NoError(t, err)
	peer, err := UnmarshalProto(global.NamespaceReplicas, out)
	require.NoError(t, err)
	_, ok = peer.(*peers.Peer)
	require.True(t, ok)
	obj, err = UnmarshalObject(global.NamespaceReplicas, out, false)
	require.NoError(t, err)
	require.Equal(t, "peers::8", obj.Key)
	require.Equal(t, uint64(1), obj.Version.Version)
	require.Empty(t, obj.Data)
	obj, err = UnmarshalObject(global.NamespaceReplicas, out, true)
	require.NoError(t, err)
	require.NotEmpty(t, obj.Data)

	// Test Unmarshal category index
	in, err = ioutil.ReadFile("testdata/index::categories.json")
	require.NoError(t, err)
	out, err = RemarshalJSON(global.NamespaceIndices, in)
	require.NoError(t, err)
	index, err := UnmarshalIndex(out)
	require.NoError(t, err)
	require.Equal(t, 6, len(index))

	// Test Unmarshal names index
	in, err = ioutil.ReadFile("testdata/index::names.json")
	require.NoError(t, err)
	out, err = RemarshalJSON(global.NamespaceIndices, in)
	require.NoError(t, err)
	index, err = UnmarshalIndex(out)
	require.NoError(t, err)
	require.Equal(t, 8, len(index))

	// Test Unmarshal Sequence
	out, err = RemarshalJSON(global.NamespaceSequence, []byte("42"))
	require.NoError(t, err)
	seq, err := UnmarshalSequence(out)
	require.NoError(t, err)
	require.Equal(t, uint64(42), seq)
}
