package fixtures_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/rotationalio/honu"
	"github.com/rotationalio/honu/options"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/trisacrypto/directory/pkg/gds/fixtures"
	"github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/proto"
)

func TestFixtures(t *testing.T) {
	// Discard logging from the application to focus on test logs
	// NOTE: ConsoleLog MUST be false otherwise this will be overriden
	logger.Discard()

	// Init the fixtures library
	fixturesPath := filepath.Join("..", "testdata", "fakes.tgz")
	dbPath := filepath.Join("testdata", "db")

	expected := map[fixtures.FixtureType]map[string]int{
		fixtures.Empty: {},
		fixtures.Small: {wire.NamespaceContacts: 2, wire.NamespaceVASPs: 3},
		fixtures.Full:  {wire.NamespaceContacts: 2, wire.NamespaceVASPs: 14, wire.NamespaceCerts: 3, wire.NamespaceCertReqs: 10},
	}

	// Load the leveldb fixtures and verify everything was loaded correctly
	lib, err := fixtures.New(fixturesPath, dbPath, fixtures.StoreLevelDB)
	require.NoError(t, err, "could not create fixtures library")
	defer lib.Close()
	verifyFixtures(t, lib, expected)

	// Close the library so we can load the Trtl fixtures
	require.NoError(t, lib.Close())

	// Load the Trtl fixtures and verify everything was loaded correctly
	lib, err = fixtures.New(fixturesPath, dbPath, fixtures.StoreTrtl)
	require.NoError(t, err, "could not create fixtures library")
	verifyFixtures(t, lib, expected)
}

func verifyFixtures(t *testing.T, lib *fixtures.Library, expected map[fixtures.FixtureType]map[string]int) {
	// Test the reference fixtures
	refs := lib.Fixtures()
	require.Len(t, refs, 4, "unexpected number of namespaces in fixtures")
	require.Contains(t, refs, wire.NamespaceContacts)
	require.Contains(t, refs, wire.NamespaceVASPs)
	require.Contains(t, refs, wire.NamespaceCerts)
	require.Contains(t, refs, wire.NamespaceCertReqs)
	require.Len(t, refs[wire.NamespaceContacts], expected[fixtures.Full][wire.NamespaceContacts])
	require.Len(t, refs[wire.NamespaceVASPs], expected[fixtures.Full][wire.NamespaceVASPs])
	require.Len(t, refs[wire.NamespaceCerts], expected[fixtures.Full][wire.NamespaceCerts])
	require.Len(t, refs[wire.NamespaceCertReqs], expected[fixtures.Full][wire.NamespaceCertReqs])

	// Validate VASP fixtures
	for name, obj := range refs[wire.NamespaceVASPs] {
		vasp, ok := obj.(*pb.VASP)
		require.True(t, ok, "could not marshal VASP record %s", name)
		require.NoError(t, vasp.Validate(true), "VASP record %s is invalid", name)
	}

	for ftype := range expected {
		switch lib.StoreType() {
		case fixtures.StoreLevelDB:
			// Test the levelDB fixtures
			require.NoError(t, lib.Load(ftype))
			defer lib.Reset()

			// Count the number of things in the database
			counts := countLevelDBFixtures(t, lib)
			require.Equal(t, expected[ftype], counts)
		case fixtures.StoreTrtl:
			// Test the Trtl fixtures
			require.NoError(t, lib.Load(ftype))
			defer lib.Reset()

			// Count the number of things in the database
			counts := countHonuFixtures(t, lib)
			require.Equal(t, expected[ftype], counts)
		default:
			require.Fail(t, "unknown store type")
		}
	}
}

func countLevelDBFixtures(t *testing.T, lib *fixtures.Library) (counts map[string]int) {
	db, err := leveldb.OpenFile(lib.DBPath(), nil)
	require.NoError(t, err, "could not open levelDB at %s", lib.DBPath())
	defer db.Close()

	counts = make(map[string]int)

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		// Fetch the key and split the namespace from the ID
		key := strings.Split(string(iter.Key()), "::")
		require.Len(t, key, 2, "key does not have a namespace prefix")

		// Ensure we can unmarshal the fixture
		var obj interface{}
		switch prefix := key[0]; prefix {
		case wire.NamespaceContacts:
			contact := &models.Contact{}
			require.NoError(t, proto.Unmarshal(iter.Value(), contact))
			obj = contact
		case wire.NamespaceVASPs:
			vasp := &pb.VASP{}
			require.NoError(t, proto.Unmarshal(iter.Value(), vasp))
			obj = vasp
		case wire.NamespaceCerts:
			cert := &models.Certificate{}
			require.NoError(t, proto.Unmarshal(iter.Value(), cert))
			obj = cert
		case wire.NamespaceCertReqs:
			certreq := &models.CertificateRequest{}
			require.NoError(t, proto.Unmarshal(iter.Value(), certreq))
			obj = certreq
		case wire.NamespaceIndices:
			continue
		default:
			require.Fail(t, "unrecognized object for namespace %q", prefix)
		}

		// Count occurrence of the key
		counts[key[0]]++

		// Test that the database fixture matches our reference
		match, err := lib.CompareFixture(key[0], key[1], obj, false)
		require.NoError(t, err, "could not compare leveldb fixture %s::%s to reference", key[0], key[1])
		require.True(t, match, "leveldb fixture %s::%s does not match reference", key[0], key[1])
	}

	require.NoError(t, iter.Error())
	iter.Release()
	return counts
}

func countHonuFixtures(t *testing.T, lib *fixtures.Library) (counts map[string]int) {
	// Stop the Trtl server so we can open the database with Honu
	require.NoError(t, lib.ShutdownTrtl())
	defer lib.SetupTrtl()

	db, err := honu.Open("leveldb:///" + lib.DBPath())
	require.NoError(t, err, "could not open Honu db at %s", lib.DBPath())
	defer db.Close()

	counts = make(map[string]int)

	iter, err := db.Iter(nil, options.WithNamespace(wire.NamespaceContacts))
	require.NoError(t, err, "could not create Honu contacts iterator")
	for iter.Next() {
		contact := &models.Contact{}
		require.NoError(t, proto.Unmarshal(iter.Value(), contact))
		counts[wire.NamespaceContacts]++
		lib.CompareFixture(wire.NamespaceContacts, string(iter.Key()), contact, false)
	}

	iter, err = db.Iter(nil, options.WithNamespace(wire.NamespaceVASPs))
	require.NoError(t, err, "could not create Honu vasp iterator")
	for iter.Next() {
		vasp := &pb.VASP{}
		require.NoError(t, proto.Unmarshal(iter.Value(), vasp))
		counts[wire.NamespaceVASPs]++
		lib.CompareFixture(wire.NamespaceVASPs, string(iter.Key()), vasp, false)
	}
	require.NoError(t, iter.Error())
	iter.Release()

	iter, err = db.Iter(nil, options.WithNamespace(wire.NamespaceCerts))
	require.NoError(t, err, "could not create Honu certs iterator")
	for iter.Next() {
		cert := &models.Certificate{}
		require.NoError(t, proto.Unmarshal(iter.Value(), cert))
		counts[wire.NamespaceCerts]++
		lib.CompareFixture(wire.NamespaceCerts, string(iter.Key()), cert, false)
	}
	require.NoError(t, iter.Error())
	iter.Release()

	iter, err = db.Iter(nil, options.WithNamespace(wire.NamespaceCertReqs))
	require.NoError(t, err, "could not create Honu certreq iterator")
	for iter.Next() {
		certreq := &models.CertificateRequest{}
		require.NoError(t, proto.Unmarshal(iter.Value(), certreq))
		counts[wire.NamespaceCertReqs]++
		lib.CompareFixture(wire.NamespaceCertReqs, string(iter.Key()), certreq, false)
	}
	require.NoError(t, iter.Error())
	iter.Release()

	return counts
}
