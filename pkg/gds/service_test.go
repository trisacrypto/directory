package gds_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/suite"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/trisacrypto/directory/pkg/gds"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/emails"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/store"
	"github.com/trisacrypto/directory/pkg/utils"
	"github.com/trisacrypto/directory/pkg/utils/bufconn"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	vasps    = "vasps"
	certreqs = "certreqs"
	index    = "index"
	bufSize  = 1024 * 1024
)

var (
	update             = flag.Bool("update", false, "update the gzipped test databases")
	fakesFixturePath   = filepath.Join("testdata", "fakes.tgz")
	fullDBFixturePath  = filepath.Join("testdata", "db.tgz")
	smallDBFixturePath = filepath.Join("testdata", "smalldb.tgz")
	smallDBSubset      = map[string]map[string]struct{}{
		vasps: {
			"charliebank": {},
			"delta":       {},
		},
	}
)

// The GDS Test Suite provides mock functionality and database fixtures for the services
// defined in the GDS package. Most tests in this package should be methods of the Test
// Suite. On startup, the reference fixtures are loaded from the `fakes.tgz` dataset and
// a mock service is created with an empty database. If tests require fixtures to be
// loaded they should call the loadFixtures() or loadSmallFixtures() methods to point
// the mock service at a database that has those fixtures or to loadEmptyFixtures() if
// they require an empty database. If the test modifies the database they should defer
// a call to resetFixtures, resetSmallFixtures, or resetEmptyFixtures as necessary.
//
// Tests should use accessor methods such as s.svc.GetAdmin() or s.svc.GetStore() to
// access the internals of the service and services for testing purposes.
type gdsTestSuite struct {
	suite.Suite
	ftype    fixtureType
	fixtures map[string]map[string]interface{}
	svc      *gds.Service
	dbPaths  map[fixtureType]string
	grpc     *bufconn.GRPCListener
	conf     *config.Config
}

// SetConfig allows a custom config to be specified by the tests.
// Note that loadFixtures() needs to be called in order for the config to be used.
func (s *gdsTestSuite) SetConfig(conf config.Config) {
	s.conf = &conf
}

// ResetConfig back to the default.
func (s *gdsTestSuite) ResetConfig() {
	s.conf = nil
}

func (s *gdsTestSuite) SetupSuite() {
	var err error
	require := s.Require()
	gin.SetMode(gin.TestMode)

	// Create a bufconn listener to use for gRPC requests
	s.grpc = bufconn.New(bufSize)

	// Create database paths to unpack fixtures to
	s.dbPaths = make(map[fixtureType]string)
	for _, ftype := range []fixtureType{empty, small, full} {
		s.dbPaths[ftype], err = ioutil.TempDir("testdata", "db-*")
		require.NoError(err)
	}

	// Load the reference fixtures into memory for use in testing
	s.loadReferenceFixtures()

	// Generate the databases if required (depends on the loaded reference fixtures)
	if *update || !pathExists(fullDBFixturePath) || !pathExists(smallDBFixturePath) {
		s.generateDB()
	}

	// Start with an empty fixtures service
	s.LoadEmptyFixtures()
}

func (s *gdsTestSuite) TearDownSuite() {
	if s.svc != nil && s.svc.GetStore() != nil {
		s.svc.GetStore().Close()
	}

	for _, path := range s.dbPaths {
		os.RemoveAll(path)
		log.Info().Str("path", path).Msg("cleaned up database fixture")
	}
	s.grpc.Release()
}

func TestGDS(t *testing.T) {
	suite.Run(t, new(gdsTestSuite))
}

func (s *gdsTestSuite) TestFixtures() {
	// Ensure all fixtures are loaded and extracted
	s.LoadFullFixtures()
	s.LoadSmallFixtures()
	s.LoadEmptyFixtures()

	// Close the empty fixtures so we can open a leveldb connection to it
	s.svc.GetStore().Close()
	defer s.ResetEmptyFixtures()

	require := s.Require()
	expected := map[fixtureType]map[string]int{
		empty: {},
		small: {vasps: 2},
		full:  {vasps: 14, certreqs: 10},
	}

	// Test the reference fixtures
	require.Len(s.fixtures, 2)
	require.Contains(s.fixtures, vasps)
	require.Contains(s.fixtures, certreqs)
	require.Len(s.fixtures[vasps], expected[full][vasps])
	require.Len(s.fixtures[certreqs], expected[full][certreqs])

	// Test the loaded database paths
	require.Len(s.dbPaths, 3)
	for ftype := range expected {
		require.Contains(s.dbPaths, ftype)
	}

	// Test the fixtures databases
	for ftype, dbPath := range s.dbPaths {
		db, err := leveldb.OpenFile(dbPath, nil)
		require.NoError(err, "could not open ftype %d db at %s", ftype, dbPath)

		// Ensure database is closed if there is an error with the test
		defer db.Close()

		// Count the number of things in the database
		counts := make(map[string]int)
		iter := db.NewIterator(nil, nil)
		defer iter.Release()

		for iter.Next() {
			// Fetch the key and split the namespace from the ID
			key := strings.Split(string(iter.Key()), "::")
			require.Len(key, 2, "key does not have a namespace prefix")

			// Ensure we can unmarshal the fixture
			var obj interface{}
			switch prefix := key[0]; prefix {
			case vasps:
				vasp := &pb.VASP{}
				require.NoError(proto.Unmarshal(iter.Value(), vasp))
				obj = vasp
			case certreqs:
				certreq := &models.CertificateRequest{}
				require.NoError(proto.Unmarshal(iter.Value(), certreq))
				obj = certreq
			case index:
				continue
			default:
				require.Fail("unrecognized object for namespace %q", prefix)
			}

			// Count occurrence of the key
			counts[key[0]]++

			// Test that the database fixture matches our reference
			s.CompareFixture(key[0], key[1], obj, false)
		}

		// Ensure we have the expected number of items in the database
		require.NoError(iter.Error())
		require.Equal(expected[ftype], counts)

		// Ensure the database is closed before next loop
		db.Close()
	}
}

//===========================================================================
// Custom Assertions
//===========================================================================

func (s *gdsTestSuite) CompareFixture(namespace, key string, obj interface{}, removeExtra bool) {
	var (
		ok bool
	)

	require := s.Require()

	_, ok = s.fixtures[namespace]
	require.True(ok, "unknown namespace %s", namespace)

	// Reset any time fields for the comparison and compare directly
	switch namespace {
	case vasps:
		var a *pb.VASP
		for _, f := range s.fixtures[namespace] {
			ref := f.(*pb.VASP)
			if ref.Id == key {
				a = ref
				break
			}
		}
		require.NotNil(a, "unknown VASP fixture %s", key)

		b, ok := obj.(*pb.VASP)
		require.True(ok, "obj is not a VASP object")

		// Remove time fields for comparison
		a.LastUpdated, b.LastUpdated = "", ""

		if removeExtra {
			a.Extra, b.Extra = nil, nil
			a.Contacts.Administrative.Extra, b.Contacts.Administrative.Extra = nil, nil
			a.Contacts.Technical.Extra, b.Contacts.Technical.Extra = nil, nil
			a.Contacts.Legal.Extra, b.Contacts.Legal.Extra = nil, nil
			a.Contacts.Billing.Extra, b.Contacts.Billing.Extra = nil, nil
		}

		require.True(proto.Equal(a, b), "vasps are not the same")

	case certreqs:
		var a *models.CertificateRequest
		for _, f := range s.fixtures[namespace] {
			ref := f.(*models.CertificateRequest)
			if ref.Id == key {
				a = ref
				break
			}
		}
		require.NotNil(a, "unknown CertificateRequest fixture %s", key)

		b, ok := obj.(*models.CertificateRequest)
		require.True(ok, "obj is not a CertificateRequest object")

		a.Modified, b.Modified = "", ""
		a.Created, b.Created = "", ""

		require.True(proto.Equal(a, b), "certreqs are not the same")

	default:
		require.Fail("unhandled namespace %s", namespace)
	}

}

type emailMeta struct {
	contact   *pb.Contact
	to        string
	from      string
	subject   string
	reason    string
	timestamp time.Time
}

// CheckEmails verifies that the provided email messages exist in both the email mock
// and the audit log on the contact, if the email was sent to a contact.
func (s *gdsTestSuite) CheckEmails(messages []*emailMeta) {
	require := s.Require()

	var sentEmails []*sgmail.SGMailV3

	// Check total number of emails sent
	require.Len(emails.MockEmails, len(messages))

	// Get emails from the mock
	for _, data := range emails.MockEmails {
		msg := &sgmail.SGMailV3{}
		require.NoError(json.Unmarshal(data, msg))
		sentEmails = append(sentEmails, msg)
	}

	for _, msg := range messages {
		// If the email was sent to a contact, check the audit log
		if msg.contact != nil {
			log, err := models.GetEmailLog(msg.contact)
			require.NoError(err)
			require.Len(log, 1, "contact %s has unexpected number of email logs", msg.contact.Email)
			require.Equal(msg.reason, log[0].Reason)
			ts, err := time.Parse(time.RFC3339, log[0].Timestamp)
			require.NoError(err)
			require.True(ts.Sub(msg.timestamp) < time.Minute, "timestamp in email log is too old")
		}

		expectedRecipient, err := mail.ParseAddress(msg.to)
		require.NoError(err)

		// Search for the sent email in the mock and check the metadata
		found := false
		for _, sent := range sentEmails {
			recipient, err := emails.GetRecipient(sent)
			require.NoError(err)
			if recipient == expectedRecipient.Address {
				found = true
				sender, err := mail.ParseAddress(msg.from)
				require.NoError(err)
				require.Equal(sender.Address, sent.From.Address)
				require.Equal(msg.subject, sent.Subject)
				break
			}
		}
		require.True(found, "email not sent for recipient %s", msg.to)
	}
}

//===========================================================================
// Test fixtures management
//===========================================================================

type fixtureType uint8

const (
	unknown fixtureType = iota
	empty
	small
	full
)

// Load reference fixtures into s.fixtures; the reference fixtures contain all possible
// fixtures that can be in any fixture database. The fixtures themselves are a mapping
// of namespace to a map of ID strings to the unmarshaled value. E.g. to get the
// expected value for VASP 1234 you would use s.fixtures[vasps]["1234"]
func (s *gdsTestSuite) loadReferenceFixtures() {
	require := s.Require()

	// Extract the gzip archive
	root, err := utils.ExtractGzip(fakesFixturePath, "testdata", true)
	require.NoError(err)
	defer os.RemoveAll(root)

	// Create the reference fixtures map
	s.fixtures = make(map[string]map[string]interface{})
	for _, namespace := range []string{vasps, certreqs} {
		s.fixtures[namespace] = make(map[string]interface{})
	}

	// Load the JSON fixtures into the fixtures map
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		require.NoError(err)
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".json" {
			return nil
		}

		// Unmarshal the JSON into the global fixtures map.
		data, err := os.ReadFile(path)
		require.NoError(err)
		parts := strings.Split(strings.TrimSuffix(info.Name(), ".json"), "::")
		require.Len(parts, 2)

		prefix := parts[0]
		key := parts[1]

		switch prefix {
		case vasps:
			vasp := &pb.VASP{}
			err = protojson.Unmarshal(data, vasp)
			require.NoError(err)
			s.fixtures[vasps][key] = vasp
		case certreqs:
			cert := &models.CertificateRequest{}
			err = protojson.Unmarshal(data, cert)
			require.NoError(err)
			s.fixtures[certreqs][key] = cert
		default:
			return fmt.Errorf("unrecognized prefix for file: %s", info.Name())
		}
		return nil
	})
	require.NoError(err)
}

func (s *gdsTestSuite) loadFixtures(ftype fixtureType, fpath string) {
	// If we're already at the specified fixture type and no custom config is provided,
	// do nothing
	if s.ftype == ftype && s.conf == nil {
		log.Info().Uint8("ftype", uint8(ftype)).Str("path", fpath).Msg("CACHED FIXTURE")
		return
	}

	// Close the current service
	var err error
	require := s.Require()
	if s.svc != nil && s.svc.GetStore() != nil {
		if err := s.svc.GetStore().Close(); err != nil {
			log.Warn().Err(err).Msg("could not close service store to load new fixtures")
		}
	}

	// If we're given a fixture path to extract, check if dir is empty, otherwise extract it
	if fpath != "" {
		if !pathExists(filepath.Join(s.dbPaths[ftype], "CURRENT")) {
			if _, err := utils.ExtractGzip(fpath, s.dbPaths[ftype], false); err != nil {
				log.Warn().Err(err).Str("db", fpath).Msg("unable to extract test fixtures")
			}
			log.Info().Uint8("ftype", uint8(ftype)).Str("path", fpath).Str("dbpath", s.dbPaths[ftype]).Msg("FIXTURE EXTRACTED")
		}
	}

	// Shutdown old GDS server
	if s.svc != nil {
		if err := s.svc.GetGDS().Shutdown(); err != nil {
			log.Warn().Err(err).Msg("could not shutdown GDS server to start new one")
		}
	}
	if s.grpc != nil {
		s.grpc.Release()
	}

	// Use the custom config if specified
	var conf config.Config
	if s.conf != nil {
		conf = *s.conf
	} else {
		conf = gds.MockConfig()
	}

	// Create the new service with a database to the specified path
	conf.Database.URL = "leveldb:///" + s.dbPaths[ftype]
	s.svc, err = gds.NewMock(conf)
	require.NoError(err, "could not create mock GDS service")

	// Start the new GDS server using a bufconn listener to avoid network requests
	s.grpc = bufconn.New(bufSize)
	go s.svc.GetGDS().Run(s.grpc.Listener)

	s.ftype = ftype
	log.Info().Uint8("ftype", uint8(ftype)).Str("path", fpath).Msg("FIXTURE LOADED")
}

func (s *gdsTestSuite) LoadEmptyFixtures() {
	s.loadFixtures(empty, "")
}

// LoadFullFixtures loads the JSON test fixtures from disk and stores them in the dbFixtures map.
func (s *gdsTestSuite) LoadFullFixtures() {
	s.loadFixtures(full, fullDBFixturePath)
}

func (s *gdsTestSuite) LoadSmallFixtures() {
	s.loadFixtures(small, smallDBFixturePath)
}

func (s *gdsTestSuite) resetFixtures(ftype fixtureType) {
	var err error
	require := s.Require()

	// Delete the current fixture database
	os.RemoveAll(s.dbPaths[ftype])
	s.dbPaths[ftype], err = ioutil.TempDir("testdata", "db-*")
	require.NoError(err)

	// Set the ftype to unknown to ensure the loader loads the fixture
	s.ftype = unknown
}

func (s *gdsTestSuite) ResetEmptyFixtures() {
	s.resetFixtures(empty)
}

func (s *gdsTestSuite) ResetFullFixtures() {
	s.resetFixtures(full)
}

func (s *gdsTestSuite) ResetSmallFixtures() {
	s.resetFixtures(small)
}

// generateDB generates an updated database and compresses it to a gzip file.
// Note: This also generates a temporary directory which the suite teardown
// should clean up.
func (s *gdsTestSuite) generateDB() {
	// Data is required to generate the database
	require := s.Require()
	require.NotEmpty(s.fixtures, "there are no reference fixtures to generate the db with")
	require.NotEmpty(s.fixtures[vasps], "there are no reference vasps fixtures to generate the db with")
	require.NotEmpty(s.fixtures[certreqs], "there are no reference certreqs fixtures to generate the db with")

	// Loop through each database type and create the database
	for ftype, path := range s.dbPaths {
		// No need to do anything with the empty database
		if ftype == empty {
			continue
		}

		// Open a Store object to write fixtures to the database
		store, err := store.Open(config.DatabaseConfig{
			URL:           "leveldb:///" + path,
			ReindexOnBoot: false,
		})
		require.NoError(err)
		defer store.Close()

		// Loop through all the reference fixtures, adding them as necessary
		for namespace, items := range s.fixtures {
		itemLoop:
			for key, item := range items {

				// If we're in small database mode, check if we should add the data
				if ftype == small {
					if _, ok := smallDBSubset[namespace][key]; !ok {
						continue itemLoop
					}
				}

				// Add the fixture to the database, updating indices
				switch namespace {
				case vasps:
					vasp := item.(*pb.VASP)
					id, err := store.CreateVASP(vasp)
					require.NoError(err, "could not insert VASP into store")
					require.Equal(vasp.Id, id)
				case certreqs:
					err = store.UpdateCertReq(item.(*models.CertificateRequest))
					require.NoError(err, "could not insert CertificateRequest into store")
				default:
					require.Fail("unrecognized object for namespace %q", namespace)
				}
			}
		}

		// Close the database to sync the indices
		require.NoError(store.Close(), "could not close store")

		// Write the database to the fixture path
		var fixturePath string
		switch ftype {
		case small:
			fixturePath = smallDBFixturePath
		case full:
			fixturePath = fullDBFixturePath
		default:
			require.Fail("unrecognized database path for ftype %d", ftype)
		}

		require.NoError(utils.WriteGzip(path, fixturePath))
		log.Info().Str("db", fixturePath).Msg("successfully regenerated test database")
	}
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	// Note this will return true if there is another error like a permissions error.
	// Those errors will be caught when the file is unzipped.
	return true
}

func remarshalProto(namespace string, obj map[string]interface{}) (_ protoreflect.ProtoMessage, err error) {
	var data []byte
	if data, err = json.Marshal(obj); err != nil {
		return nil, err
	}

	jsonpb := &protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	switch namespace {
	case vasps:
		vasp := &pb.VASP{}
		if err = jsonpb.Unmarshal(data, vasp); err != nil {
			return nil, err
		}
		return vasp, nil
	case certreqs:
		certreq := &models.CertificateRequest{}
		if err = jsonpb.Unmarshal(data, certreq); err != nil {
			return nil, err
		}
		return certreq, nil
	default:
		return nil, fmt.Errorf("unknown namespace %q", namespace)
	}
}
