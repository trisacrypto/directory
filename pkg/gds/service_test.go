package gds_test

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/honu"
	"github.com/rotationalio/honu/options"
	"github.com/rs/zerolog/log"
	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/suite"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/trisacrypto/directory/pkg/gds"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/emails"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/store"
	trtlstore "github.com/trisacrypto/directory/pkg/gds/store/trtl"
	"github.com/trisacrypto/directory/pkg/trtl"
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
	fakesFixturePath = filepath.Join("testdata", "fakes.tgz")
	smallDBSubset    = map[string]map[string]struct{}{
		vasps: {
			"charliebank": {},
			"delta":       {},
		},
	}
	dbPath = filepath.Join("testdata", "db")
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
	ftype        fixtureType
	fixtures     map[string]map[string]interface{}
	svc          *gds.Service
	stype        storeType
	trtl         *trtl.Server
	trtlListener *bufconn.GRPCListener
	grpc         *bufconn.GRPCListener
	conf         *config.Config
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
	gin.SetMode(gin.TestMode)

	// Load the reference fixtures into memory for use in testing
	s.loadReferenceFixtures()

	// Start the Trtl server if required
	if s.stype == storeTrtl {
		s.SetupTrtl()
	}

	// Generate the fixture database
	s.generateDB(empty)

	// Start with an empty fixtures service
	s.LoadEmptyFixtures()
}

// SetupGDS starts the GDS server
// Run this inside the test methods after loading the appropriate fixtures
func (s *gdsTestSuite) SetupGDS() {

	// Using a bufconn listener allows us to avoid network requests
	s.grpc = bufconn.New(bufSize)
	go s.svc.GetGDS().Run(s.grpc.Listener)
}

// SetupMembers starts the Members server
// Run this inside the test methods after loading the appropriate fixtures
func (s *gdsTestSuite) SetupMembers() {

	// Using a bufconn listener allows us to avoid network requests
	s.grpc = bufconn.New(bufSize)
	go s.svc.GetMembers().Run(s.grpc.Listener)
}

// SetupTrtl starts the Trtl server, which must be done before calling generateDB.
func (s *gdsTestSuite) SetupTrtl() {
	var err error
	require := s.Require()

	conf := trtl.MockConfig()
	conf.Database.URL = "leveldb:///" + dbPath

	// Mark as processed since the config wasn't loaded from the envrionment
	conf, err = conf.Mark()
	require.NoError(err)

	// Start the Trtl server
	s.trtl, err = trtl.New(conf)
	require.NoError(err, "could not start Trtl server")

	// Using a bufconn listener allows us to avoid network requests
	s.trtlListener = bufconn.New(bufSize)
	go s.trtl.Run(s.trtlListener.Listener)

	// Connect to the running Trtl server
	require.NoError(s.trtlListener.Connect())
}

// Helper function to shutdown any previously running GDS or Members servers and release the gRPC connection
func (s *gdsTestSuite) shutdownServers() {
	// Shutdown old GDS and Members servers, if they exist
	if s.svc != nil {
		if err := s.svc.GetGDS().Shutdown(); err != nil {
			log.Warn().Err(err).Msg("could not shutdown GDS server to start new one")
		}
		if err := s.svc.GetMembers().Shutdown(); err != nil {
			log.Warn().Err(err).Msg("could not shutdown Members server to start new one")
		}
	}
	if s.grpc != nil {
		s.grpc.Release()
	}

	// Shutdown the Trtl server if it is running
	if s.trtl != nil {
		if err := s.trtl.Shutdown(); err != nil {
			log.Warn().Err(err).Msg("could not shutdown Trtl server to start new one")
		}
	}
	if s.trtlListener != nil {
		s.trtlListener.Release()
	}
}

func (s *gdsTestSuite) TearDownSuite() {
	if s.svc != nil && s.svc.GetStore() != nil {
		s.svc.GetStore().Close()
	}

	s.shutdownServers()
	os.RemoveAll(dbPath)
}

func TestGDSLevelDB(t *testing.T) {
	s := new(gdsTestSuite)
	s.stype = storeLevelDB
	suite.Run(t, s)
}

func TestGDSTrtl(t *testing.T) {
	s := new(gdsTestSuite)
	s.stype = storeTrtl
	suite.Run(t, s)
}

func (s *gdsTestSuite) TestFixtures() {
	s.loadReferenceFixtures()

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

	for ftype := range expected {
		switch s.stype {
		case storeLevelDB:
			// Test the levelDB fixtures
			s.loadFixtures(ftype)

			// Close the database so we can open a leveldb connection to it
			s.svc.GetStore().Close()
			defer s.ResetFixtures()

			db, err := leveldb.OpenFile(dbPath, nil)
			require.NoError(err, "could not open levelDB at %s", dbPath)

			// Ensure database is closed if there is an error with the test
			defer db.Close()

			// Count the number of things in the database
			counts := s.countLevelDBFixtures(db)
			require.Equal(expected[ftype], counts)

			// Ensure the database is closed before next loop
			db.Close()

			s.ResetFixtures()
		case storeTrtl:
			// Test the Trtl fixtures
			s.loadFixtures(ftype)

			// Stop the Trtl server so we can open the database with Honu
			s.trtl.Shutdown()
			defer s.ResetFixtures()

			// Count the number of things in the database
			hdb, err := honu.Open("leveldb:///" + dbPath)
			require.NoError(err, "could not open Honu db at %s", dbPath)

			// Ensure database is closed if there is an error with the test
			defer hdb.Close()

			// Count the number of things in the database
			counts := s.countHonuFixtures(hdb)
			require.Equal(expected[ftype], counts)

			// Ensure the database is closed before next loop
			hdb.Close()
		default:
			require.Fail("unknown store type")
		}
	}
}

func (s *gdsTestSuite) countLevelDBFixtures(db *leveldb.DB) (counts map[string]int) {
	require := s.Require()
	counts = make(map[string]int)

	iter := db.NewIterator(nil, nil)
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

	require.NoError(iter.Error())
	iter.Release()
	return counts
}

func (s *gdsTestSuite) countHonuFixtures(db *honu.DB) (counts map[string]int) {
	require := s.Require()
	counts = make(map[string]int)

	iter, err := db.Iter(nil, options.WithNamespace(vasps))
	require.NoError(err, "could not create Honu vasp iterator")
	for iter.Next() {
		vasp := &pb.VASP{}
		require.NoError(proto.Unmarshal(iter.Value(), vasp))
		counts[vasps]++
		s.CompareFixture(vasps, string(iter.Key()), vasp, false)
	}
	require.NoError(iter.Error())
	iter.Release()

	iter, err = db.Iter(nil, options.WithNamespace(certreqs))
	require.NoError(err, "could not create Honu certreq iterator")
	for iter.Next() {
		certreq := &models.CertificateRequest{}
		require.NoError(proto.Unmarshal(iter.Value(), certreq))
		counts[certreqs]++
		s.CompareFixture(certreqs, string(iter.Key()), certreq, false)
	}
	require.NoError(iter.Error())
	iter.Release()

	return counts
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
		var data pb.VASP
		for _, f := range s.fixtures[namespace] {
			ref := f.(*pb.VASP)
			if ref.Id == key {
				// Avoid modifying the object in the fixtures map
				data = *ref
				a = &data
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
			iter := models.NewContactIterator(a.Contacts, false, false)
			for iter.Next() {
				contact, _ := iter.Value()
				contact.Extra = nil
			}

			iter = models.NewContactIterator(b.Contacts, false, false)
			for iter.Next() {
				contact, _ := iter.Value()
				contact.Extra = nil
			}
		}

		require.True(proto.Equal(a, b), "vasps are not the same")

	case certreqs:
		var a *models.CertificateRequest
		var data models.CertificateRequest
		for _, f := range s.fixtures[namespace] {
			ref := f.(*models.CertificateRequest)
			if ref.Id == key {
				// Avoid modifying the object in the fixtures map
				data = *ref
				a = &data
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
	require.Len(emails.MockEmails, len(messages), "incorrect number of emails sent")

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

type storeType uint8

const (
	storeUnknown storeType = iota
	storeLevelDB
	storeTrtl
)

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

		prefix := filepath.Base(filepath.Dir(path))     // prefix is the directory in the fixture
		key := strings.TrimSuffix(info.Name(), ".json") // key is the name of the file in the fixture

		// Unmarshal the JSON into the global fixtures map.
		data, err := os.ReadFile(path)
		require.NoError(err)

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

// loadFixtures loads a new set of fixtures into the database. This method must respect
// the ftype variable on the test suite, which indicates which fixtures are currently
// loaded. If the ftype is different than the indicated fixture type, then this causes
// the current database to be completely overwritten with the indicated fixtures. Tests
// that require database fixtures should call the appropriate load method, either
// LoadFullFixtures, LoadSmallFixtures, or LoadEmptyFixtures to ensure that the correct
// fixtures are present before test execution.
func (s *gdsTestSuite) loadFixtures(ftype fixtureType) {
	// If we're already at the specified fixture type and no custom config is provided,
	// do nothing
	if s.ftype == ftype && s.conf == nil {
		log.Info().Uint8("ftype", uint8(ftype)).Msg("CACHED FIXTURE")
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

	s.shutdownServers()

	// Recreate the database directory
	require.NoError(os.RemoveAll(dbPath), "could not remove fixture database directory")
	db, err := leveldb.OpenFile(dbPath, nil)
	require.NoError(err, "could not open new fixture database")
	db.Close()

	// Use the custom config if specified
	var conf config.Config
	if s.conf != nil {
		conf = *s.conf
	} else {
		conf = gds.MockConfig()
	}

	// Store-specific handling of the database
	switch s.stype {
	case storeLevelDB:
		conf.Database.URL = "leveldb:///" + dbPath
	case storeTrtl:
		s.SetupTrtl()
	default:
		require.Fail("unrecognized store type")
	}
	s.generateDB(ftype)

	// Create the new service
	if s.trtlListener != nil {
		require.NoError(s.trtlListener.Connect())
		s.svc, err = gds.NewMock(conf, s.trtlListener.Conn)
	} else {
		s.svc, err = gds.NewMock(conf, nil)
	}

	require.NoError(err, "could not create mock GDS service")

	s.ftype = ftype
	log.Info().Uint8("ftype", uint8(ftype)).Msg("FIXTURE LOADED")
}

func (s *gdsTestSuite) LoadEmptyFixtures() {
	s.loadFixtures(empty)
}

// LoadFullFixtures loads the JSON test fixtures from disk and stores them in the dbFixtures map.
func (s *gdsTestSuite) LoadFullFixtures() {
	s.loadFixtures(full)
}

func (s *gdsTestSuite) LoadSmallFixtures() {
	s.loadFixtures(small)
}

// ResetFixtures uncaches the current database which causes the next call to
// loadFixtures to generate a new database that overwrites the current one. Tests that
// modify the database should call ResetFixtures to ensure that the fixtures are reset
// for the next test.
func (s *gdsTestSuite) ResetFixtures() {
	// Set the ftype to unknown to ensure that loadFixtures loads the fixture.
	s.ftype = unknown
}

// generateDB generates an updated database and compresses it to a gzip file.
// Note: This also generates a temporary directory which the suite teardown
// should clean up.
func (s *gdsTestSuite) generateDB(ftype fixtureType) {
	// Data is required to generate the database
	require := s.Require()
	require.NotEmpty(s.fixtures, "there are no reference fixtures to generate the db with")
	require.NotEmpty(s.fixtures[vasps], "there are no reference vasps fixtures to generate the db with")
	require.NotEmpty(s.fixtures[certreqs], "there are no reference certreqs fixtures to generate the db with")

	// No need to do anything with the empty database
	if ftype == empty {
		return
	}

	// Create the Store object depending on the store type
	var db store.Store
	var err error
	switch s.stype {
	case storeLevelDB:
		db, err = store.Open(config.DatabaseConfig{
			URL:           "leveldb:///" + dbPath,
			ReindexOnBoot: false,
		})
		require.NoError(err, "could not open leveldb store")
		defer db.Close()
	case storeTrtl:
		require.NoError(s.trtlListener.Connect())
		db, err = trtlstore.NewMock(s.trtlListener.Conn)
		require.NoError(err, "could not open trtl store")
		defer db.Close()
	default:
		require.Fail("unrecognized store type")
	}

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
				id, err := db.CreateVASP(vasp)
				require.NoError(err, "could not insert VASP into store")
				require.Equal(vasp.Id, id)
			case certreqs:
				err = db.UpdateCertReq(item.(*models.CertificateRequest))
				require.NoError(err, "could not insert CertificateRequest into store")
			default:
				require.Fail("unrecognized object for namespace %q", namespace)
			}
		}
	}

	// Close the database to sync the indices
	require.NoError(db.Close(), "could not close store")

	log.Info().Msg("successfully regenerated test database")
}

// SetVerificationStatus sets the verification status of a VASP fixture on the
// database. This is useful for testing VerificationState checks without having to use
// multiple fixtures.
func (s *gdsTestSuite) SetVerificationStatus(id string, status pb.VerificationState) {
	require := s.Require()

	// Retrieve the VASP from the database
	vasp, err := s.svc.GetStore().RetrieveVASP(id)
	require.NoError(err, "VASP not found in database")

	// Set the verification status and write back to the database
	vasp.VerificationStatus = status
	require.NoError(s.svc.GetStore().UpdateVASP(vasp), "could not update VASP")
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
