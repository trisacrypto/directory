/*
Package wire is a store helper package that handles "wire" representations of the
objects being managed in the store - e.g. marshalling and unmarshalling data from
protocol buffers or json data and performing common operations across multiple data
types that would otherwise require a switch statement and type checking.
*/
package wire

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/trtl/peers/v1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	ErrCannotReplicate = errors.New("object in namespace cannot be replicated")
	ErrObjectNotFound  = errors.New("object not found in local store")
)

// Namespace constants for all managed objects in GDS
const (
	NamespaceVASPs         = "vasps"
	NamespaceCerts         = "certs"
	NamespaceCertReqs      = "certreqs"
	NamespaceReplicas      = "peers"
	NamespaceIndices       = "index"
	NamespaceSequence      = "sequence"
	NamespaceAnnouncements = "announcements"
	NamespaceActivities    = "activities"
	NamespaceOrganizations = "organizations"
	NamespaceContacts      = "contacts"
	NamespaceEmails        = "emails"
)

// Namespaces defines all possible namespaces that GDS manages
// TODO: do we need to add the BFF namespaces here as well?
var Namespaces = [4]string{NamespaceVASPs, NamespaceCertReqs, NamespaceReplicas, NamespaceContacts}

// UnmarshalProto expects protocol buffer data and unmarshals it to the correct type
// based on the namespace. This is a utility function for dealing with the various
// namespaces and types that GDS manages and is not a substitute for direct unmarshaling.
func UnmarshalProto(namespace string, data []byte) (_ proto.Message, err error) {
	switch namespace {
	case NamespaceVASPs:
		vasp := &pb.VASP{}
		if err = proto.Unmarshal(data, vasp); err != nil {
			return nil, fmt.Errorf("could not unmarshal %s to %T: %s", namespace, vasp, err)
		}
		return vasp, nil
	case NamespaceCerts:
		cert := &models.Certificate{}
		if err = proto.Unmarshal(data, cert); err != nil {
			return nil, fmt.Errorf("could not unmarshal %s to %T: %s", namespace, cert, err)
		}
		return cert, nil
	case NamespaceCertReqs:
		certreq := &models.CertificateRequest{}
		if err = proto.Unmarshal(data, certreq); err != nil {
			return nil, fmt.Errorf("could not unmarshal %s to %T: %s", namespace, certreq, err)
		}
		return certreq, nil
	case NamespaceReplicas:
		peer := &peers.Peer{}
		if err = proto.Unmarshal(data, peer); err != nil {
			return nil, fmt.Errorf("could not unmarshal %s to %T: %s", namespace, peer, err)
		}
		return peer, nil
	case NamespaceContacts:
		contact := &models.Contact{}
		if err = proto.Unmarshal(data, contact); err != nil {
			return nil, fmt.Errorf("could not unmarshal %s to %T: %s", namespace, contact, err)
		}
		return contact, nil
	default:
		return nil, fmt.Errorf("unknown namespaces %q", namespace)
	}
}

// UnmarshalIndex extracts a map[string]interface{} from the gzip compressed index.
func UnmarshalIndex(data []byte) (index map[string]interface{}, err error) {
	buf := bytes.NewBuffer(data)
	index = make(map[string]interface{})

	var gz *gzip.Reader
	if gz, err = gzip.NewReader(buf); err != nil {
		return nil, fmt.Errorf("could not decompress data: %s", err)
	}

	decoder := json.NewDecoder(gz)
	if err = decoder.Decode(&index); err != nil {
		return nil, fmt.Errorf("could not decode json data: %s", err)
	}

	return index, nil
}

// UnmarshalSequence extracts a uint64 from the binary data
func UnmarshalSequence(data []byte) (seq uint64, err error) {
	var n int
	if seq, n = binary.Uvarint(data); n <= 0 {
		return 0, errors.New("could not parse sequence")
	}
	return seq, nil
}

// RemarshalJSON is an odd utility, it takes raw JSON data and converts it to the
// appropriate type for database storage, e.g. marshaled protocol buffers or compressed
// json for an index. This is primarily used to take JSON data from disk and put it into
// a form that UnmarshalProto or UnmarshalIndex can use.
func RemarshalJSON(namespace string, in []byte) (out []byte, err error) {
	jsonpb := &protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	switch namespace {
	case NamespaceVASPs:
		vasp := &pb.VASP{}
		if err = jsonpb.Unmarshal(in, vasp); err != nil {
			return nil, fmt.Errorf("could not unmarshal json %s into %T: %s", namespace, vasp, err)
		}
		return proto.Marshal(vasp)
	case NamespaceCerts:
		cert := &models.Certificate{}
		if err = jsonpb.Unmarshal(in, cert); err != nil {
			return nil, fmt.Errorf("could not unmarshal json %s into %T: %s", namespace, cert, err)
		}
		return proto.Marshal(cert)
	case NamespaceCertReqs:
		certreq := &models.CertificateRequest{}
		if err = jsonpb.Unmarshal(in, certreq); err != nil {
			return nil, fmt.Errorf("could not unmarshal json %s into %T: %s", namespace, certreq, err)
		}
		return proto.Marshal(certreq)
	case NamespaceReplicas:
		peer := &peers.Peer{}
		if err = jsonpb.Unmarshal(in, peer); err != nil {
			return nil, fmt.Errorf("could not unmarshal json %s into %T: %s", namespace, peer, err)
		}
		return proto.Marshal(peer)
	case NamespaceIndices:
		// For now, we're just compressing the JSON data, not checking if it is the correct type for the index
		// TODO: should we handle indices better?
		buf := &bytes.Buffer{}
		gz := gzip.NewWriter(buf)
		if _, err = gz.Write(in); err != nil {
			return nil, fmt.Errorf("could not compress index: %s", err)
		}
		gz.Close()
		return buf.Bytes(), nil
	case NamespaceSequence:
		var seq uint64
		if err = json.Unmarshal(in, &seq); err != nil {
			return nil, fmt.Errorf("could not unmarshal json %s into %T: %s", namespace, seq, err)
		}
		out = make([]byte, binary.MaxVarintLen64)
		binary.PutUvarint(out, seq)
		return out, nil
	default:
		return nil, fmt.Errorf("unknown namespace %q: cannot remarshal json", namespace)
	}
}

// DecodeKey returns the byte representation of the key, base64 decoding if necessary.
func DecodeKey(keys string, b64decode bool) (key []byte, err error) {
	if b64decode {
		if key, err = base64.RawStdEncoding.DecodeString(keys); err != nil {
			return nil, err
		}
		return key, nil
	}
	return []byte(keys), nil
}

// EncodeKey returns the string representation of the key, base64 encoding if necessary.
func EncodeKey(key []byte, b64encode bool) string {
	if b64encode {
		return base64.RawStdEncoding.EncodeToString(key)
	}
	return string(key)
}
