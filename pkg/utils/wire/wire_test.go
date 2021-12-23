package wire_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/trtl/peers/v1"
	. "github.com/trisacrypto/directory/pkg/utils/wire"
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
	in, err := ioutil.ReadFile("testdata/vasps/838b1f57-1646-488d-a231-d71d88681cfa.json")
	require.NoError(t, err)
	out, err := RemarshalJSON(NamespaceVASPs, in)
	require.NoError(t, err)
	msg, err := UnmarshalProto(NamespaceVASPs, out)
	require.NoError(t, err)
	vasp, ok := msg.(*pb.VASP)
	require.True(t, ok)

	// Make sure that the VASP is valid
	require.Equal(t, "838b1f57-1646-488d-a231-d71d88681cfa", vasp.Id)
	require.NoError(t, vasp.Validate(false))

	// Test unmarshal certificate requests
	in, err = ioutil.ReadFile("testdata/certreqs/87657b4d-e72c-4526-9332-c8fc56adb367.json")
	require.NoError(t, err)
	out, err = RemarshalJSON(NamespaceCertReqs, in)
	require.NoError(t, err)
	msg, err = UnmarshalProto(NamespaceCertReqs, out)
	require.NoError(t, err)
	certreq, ok := msg.(*models.CertificateRequest)
	require.True(t, ok)

	// Make sure CertificateRequest is valid
	require.Equal(t, "87657b4d-e72c-4526-9332-c8fc56adb367", certreq.Id)
	require.Equal(t, "838b1f57-1646-488d-a231-d71d88681cfa", certreq.Vasp)
	require.Equal(t, models.CertificateRequestState_COMPLETED, certreq.Status)

	// Test Unmarshal Peers
	in, err = ioutil.ReadFile("testdata/peers/8.json")
	require.NoError(t, err)
	out, err = RemarshalJSON(NamespaceReplicas, in)
	require.NoError(t, err)
	msg, err = UnmarshalProto(NamespaceReplicas, out)
	require.NoError(t, err)
	peer, ok := msg.(*peers.Peer)
	require.True(t, ok)

	// Make sure Peer is valid
	require.Equal(t, "0008", peer.Key())
	require.Equal(t, "localhost:4435", peer.Addr)

	// Test Unmarshal category index
	in, err = ioutil.ReadFile("testdata/index/categories.json")
	require.NoError(t, err)
	out, err = RemarshalJSON(NamespaceIndices, in)
	require.NoError(t, err)
	index, err := UnmarshalIndex(out)
	require.NoError(t, err)
	require.Equal(t, 6, len(index))

	// Test Unmarshal names index
	in, err = ioutil.ReadFile("testdata/index/names.json")
	require.NoError(t, err)
	out, err = RemarshalJSON(NamespaceIndices, in)
	require.NoError(t, err)
	index, err = UnmarshalIndex(out)
	require.NoError(t, err)
	require.Equal(t, 8, len(index))

	// Test Unmarshal Sequence
	out, err = RemarshalJSON(NamespaceSequence, []byte("42"))
	require.NoError(t, err)
	seq, err := UnmarshalSequence(out)
	require.NoError(t, err)
	require.Equal(t, uint64(42), seq)
}
