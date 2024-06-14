package trtl_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	bff "github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/models/v1"
	store "github.com/trisacrypto/directory/pkg/store/trtl"
	"github.com/trisacrypto/directory/pkg/trtl"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/utils/bufconn"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	storeerrors "github.com/trisacrypto/directory/pkg/store/errors"
)

const (
	metaRegion = "tauceti"
	metaOwner  = "taurian"
	metaPID    = 8
	bufSize    = 1024 * 1024
)

type trtlStoreTestSuite struct {
	suite.Suite
	tmpdb string
	conf  *config.Config
	trtl  *trtl.Server
	grpc  *bufconn.GRPCListener
}

func (s *trtlStoreTestSuite) SetupSuite() {
	require := s.Require()

	// Discard logging from the application to focus on test logs
	// NOTE: ConsoleLog MUST be false otherwise this will be overwritten
	logger.Discard()

	var err error
	s.tmpdb = s.T().TempDir()

	conf := config.Config{
		Maintenance: false,
		BindAddr:    ":4436",
		LogLevel:    logger.LevelDecoder(zerolog.WarnLevel),
		ConsoleLog:  false,
		Database: config.DatabaseConfig{
			URL:           fmt.Sprintf("leveldb:///%s", s.tmpdb),
			ReindexOnBoot: false,
		},
		Replica: config.ReplicaConfig{
			Enabled: false, // Replica is tested in the replica package
			PID:     metaPID,
			Region:  metaRegion,
			Name:    metaOwner,
		},
		MTLS: config.MTLSConfig{
			Insecure: true,
		},
	}

	// Mark config as processed
	conf, err = conf.Mark()
	require.NoError(err)
	s.conf = &conf

	s.trtl, err = trtl.New(*s.conf)
	require.NoError(err)

	s.grpc = bufconn.New(bufSize, "")
	go s.trtl.Run(s.grpc.Listener)
}

func (s *trtlStoreTestSuite) TearDownSuite() {
	require := s.Require()

	// Shutdown the trtl server if it is running
	// This should shutdown all the running services and close the database
	// Note that Shutdown should be graceful and not shutdown anything not running.
	if s.trtl != nil {
		require.NoError(s.trtl.Shutdown())
	}

	// Shutdown the gRPC connection if it's running
	if s.grpc != nil {
		s.grpc.Release()
	}

	// Cleanup the tmpdb and delete any stray files
	if s.tmpdb != "" {
		os.RemoveAll(s.tmpdb)
	}

	// Reset all of the test suite variables
	s.tmpdb = ""
	s.grpc = nil
	s.trtl = nil

	logger.ResetLogger()
}

func TestTrtlStore(t *testing.T) {
	suite.Run(t, new(trtlStoreTestSuite))
}

// Tests all the directory store methods for interacting with VASPs on the Trtl DB.
func (s *trtlStoreTestSuite) TestDirectoryStore() {
	require := s.Require()

	// Load the VASP fixture
	data, err := os.ReadFile("../testdata/vasp.json")
	require.NoError(err)

	alice := &pb.VASP{}
	err = protojson.Unmarshal(data, alice)
	require.NoError(err)

	// Validate the VASP record loaded correctly and is partial
	require.NotEmpty(alice.CommonName)
	require.NotEmpty(alice.TrisaEndpoint)
	require.NoError(alice.Validate(true))
	require.Empty(alice.Id)

	// Inject bufconn connection into the store
	require.NoError(s.grpc.Connect(context.Background()))
	defer s.grpc.Close()

	db, err := store.NewMock(s.grpc.Conn)
	require.NoError(err)

	// Initially there should be no VASPs
	iter := db.ListVASPs(context.Background())
	require.False(iter.Next())
	iter.Release()

	// Should get a not found error trying to retrieve a VASP that doesn't exist
	_, err = db.RetrieveVASP(context.Background(), "12345")
	require.EqualError(err, storeerrors.ErrEntityNotFound.Error())

	// Attempt to Create the VASP
	id, err := db.CreateVASP(context.Background(), alice)
	require.NoError(err)
	require.NotEmpty(id)

	// Should not be able to create a duplicate VASP
	id2, err := db.CreateVASP(context.Background(), alice)
	require.EqualError(err, storeerrors.ErrDuplicateEntity.Error())
	require.Empty(id2)

	// Attempt to Retrieve the VASP
	alicer, err := db.RetrieveVASP(context.Background(), id)
	require.NoError(err)
	require.Equal(id, alicer.Id)
	require.Equal(alicer.FirstListed, alicer.LastUpdated)
	require.NotEmpty(alicer.LastUpdated)
	require.NotEmpty(alicer.Version)
	require.Equal(uint64(1), alicer.Version.Version)

	// Ensure the modification time rolls over to the next second for comparison
	time.Sleep(1 * time.Second)

	// Update the VASP
	alicer.Entity.Name.NameIdentifiers[0].LegalPersonName = "AliceLiteCoin, LLC"
	alicer.VerificationStatus = pb.VerificationState_VERIFIED
	alicer.VerifiedOn = "2021-06-30T10:40:40Z"
	err = db.UpdateVASP(context.Background(), alicer)
	require.NoError(err)

	alicer, err = db.RetrieveVASP(context.Background(), id)
	require.NoError(err)
	require.Equal(id, alicer.Id)
	require.NotEmpty(alicer.LastUpdated)
	require.NotEqual(alicer.FirstListed, alicer.LastUpdated)
	require.NotEmpty(alicer.Version)
	require.Equal(uint64(2), alicer.Version.Version)
	require.Equal(alicer.VerificationStatus, pb.VerificationState_VERIFIED)

	// Delete the VASP
	err = db.DeleteVASP(context.Background(), id)
	require.NoError(err)
	alicer, err = db.RetrieveVASP(context.Background(), id)
	require.ErrorIs(err, storeerrors.ErrEntityNotFound)
	require.Empty(alicer)

	// Add a few more VASPs
	err = createVASPs(db, 10, 1)
	require.NoError(err)

	// Test listing all of the VASPs
	reqs, err := db.ListVASPs(context.Background()).All()
	require.NoError(err)
	require.Len(reqs, 10)

	// Test seeking to a specific VASP
	key := reqs[5].Id
	iter = db.ListVASPs(context.Background())
	require.True(iter.SeekId(key))
	v, err := iter.VASP()
	require.NoError(err)
	require.NoError(iter.Error())
	require.Equal(key, v.Id)

	// Test Prev() and Next() interactions
	require.False(iter.Prev(), "should move behind the first VASP")
	require.True(iter.Next(), "should move to the first VASP")
	first, err := iter.VASP()
	require.NoError(err)
	require.NotNil(first)
	require.Equal(key, first.Id, "should be the first VASP")
	require.True(iter.Next(), "should move to the second VASP")
	second, err := iter.VASP()
	require.NoError(err)
	require.NotNil(second)
	require.NotEqual(key, second.Id, "should be the second VASP")

	// Consume the rest of the iterator
	for iter.Next() {
		v, err := iter.VASP()
		require.NoError(err)
		require.NotNil(v)
	}
	require.NoError(iter.Error())
	iter.Release()

	// Test iterating over all the VASPs
	var niters int
	iter = db.ListVASPs(context.Background())
	for iter.Next() {
		require.NotEmpty(iter.VASP())
		niters++
	}
	require.NoError(iter.Error())
	iter.Release()
	require.Equal(10, niters)

	// Create enough VASPs to exceed the page size
	err = createVASPs(db, 100, 11)
	require.NoError(err)

	// Test listing all of the VASPs
	reqs, err = db.ListVASPs(context.Background()).All()
	require.NoError(err)
	require.Len(reqs, 110)

	// Cleanup database
	require.NoError(deleteVASPs(db), "could not cleanup database")
}

func (s *trtlStoreTestSuite) TestCertificateStore() {
	require := s.Require()

	// Load the VASP record from testdata
	data, err := os.ReadFile("../testdata/cert.json")
	s.NoError(err)

	cert := &models.Certificate{}
	err = protojson.Unmarshal(data, cert)
	s.NoError(err)

	// Verify the certificate is loaded correctly
	s.Empty(cert.Id)
	s.NotEmpty(cert.Request)
	s.NotEmpty(cert.Vasp)
	s.Equal(models.CertificateState_ISSUED, cert.Status)
	s.NotEmpty(cert.Details)
	s.NotEmpty(cert.Details.NotAfter)
	s.NotEmpty(cert.Details.NotBefore)

	// Inject bufconn connection into the store
	require.NoError(s.grpc.Connect(context.Background()))
	defer s.grpc.Close()

	db, err := store.NewMock(s.grpc.Conn)
	require.NoError(err)

	// Initially there should be no Certs
	iter := db.ListCerts(context.Background())
	require.False(iter.Next())
	iter.Release()

	// Should get a not found error trying to retrieve a Cert that doesn't exist
	_, err = db.RetrieveCert(context.Background(), "12345")
	require.EqualError(err, storeerrors.ErrEntityNotFound.Error())

	// Attempt to Create the Cert
	id, err := db.CreateCert(context.Background(), cert)
	s.NoError(err)

	// Attempt to Retrieve the Cert
	crr, err := db.RetrieveCert(context.Background(), id)
	s.NoError(err)
	s.NotNil(crr)
	s.Equal(id, crr.Id)
	s.Equal(cert.Request, crr.Request)
	s.Equal(cert.Vasp, crr.Vasp)
	s.Equal(cert.Status, crr.Status)
	s.True(proto.Equal(cert.Details, crr.Details))

	// Attempt to save a certificate with an ID on it
	icrr := &models.Certificate{
		Id:      uuid.New().String(),
		Request: cert.Request,
		Vasp:    crr.Vasp,
		Status:  models.CertificateState_ISSUED,
		Details: cert.Details,
	}
	_, err = db.CreateCert(context.Background(), icrr)
	s.ErrorIs(err, storeerrors.ErrIDAlreadySet)

	// Sleep for a second to roll over the clock for the modified time stamp
	time.Sleep(1 * time.Second)

	// Update the Cert
	crr.Status = models.CertificateState_REVOKED
	err = db.UpdateCert(context.Background(), crr)
	s.NoError(err)

	crr, err = db.RetrieveCert(context.Background(), id)
	s.NoError(err)
	s.Equal(id, crr.Id)
	s.Equal(models.CertificateState_REVOKED, crr.Status)

	// Attempt to update a certificate request with no Id on it
	cert.Id = ""
	s.ErrorIs(db.UpdateCert(context.Background(), cert), storeerrors.ErrIncompleteRecord)

	// Delete the Cert
	err = db.DeleteCert(context.Background(), id)
	s.NoError(err)
	crr, err = db.RetrieveCert(context.Background(), id)
	s.ErrorIs(err, storeerrors.ErrEntityNotFound)
	s.Empty(crr)

	// Add a few more certificate requests
	for i := 0; i < 10; i++ {
		crr := &models.Certificate{
			Request: uuid.New().String(),
			Vasp:    uuid.New().String(),
			Status:  models.CertificateState_ISSUED,
			Details: &pb.Certificate{
				SerialNumber: []byte(uuid.New().String()),
			},
		}
		_, err := db.CreateCert(context.Background(), crr)
		s.NoError(err)
	}

	// Test listing all of the certificates
	certs, err := db.ListCerts(context.Background()).All()
	s.NoError(err)
	s.Len(certs, 10)

	// Test Prev() and Next() interactions
	iter = db.ListCerts(context.Background())
	require.False(iter.Prev(), "should move behind the first Cert")
	require.True(iter.Next(), "should move to the first Cert")
	first, err := iter.Cert()
	require.NoError(err)
	require.NotNil(first)
	require.True(iter.Next(), "should move to the second Cert")
	second, err := iter.Cert()
	require.NoError(err)
	require.NotNil(second)
	require.NotEqual(first.Id, second.Id)

	// Consume the rest of the iterator
	for iter.Next() {
		cert, err := iter.Cert()
		require.NoError(err)
		require.NotNil(cert)
	}
	require.NoError(iter.Error())
	iter.Release()

	// Create enough Certs to exceed the page size
	for i := 0; i < 100; i++ {
		crr := &models.Certificate{
			Request: uuid.New().String(),
			Vasp:    uuid.New().String(),
			Status:  models.CertificateState_EXPIRED,
			Details: &pb.Certificate{
				SerialNumber: []byte(uuid.New().String()),
			},
		}
		_, err := db.CreateCert(context.Background(), crr)
		s.NoError(err)
	}

	// Test listing all of the Cert
	certs, err = db.ListCerts(context.Background()).All()
	require.NoError(err)
	require.Len(certs, 110)
}

func (s *trtlStoreTestSuite) TestCertificateRequestStore() {
	require := s.Require()

	// Load the VASP record from testdata
	data, err := os.ReadFile("../testdata/certreq.json")
	s.NoError(err)

	certreq := &models.CertificateRequest{}
	err = protojson.Unmarshal(data, certreq)
	s.NoError(err)

	// Verify the certificate request is loaded correctly
	s.Empty(certreq.Id)
	s.NotEmpty(certreq.Vasp)
	s.NotEmpty(certreq.CommonName)
	s.Equal(models.CertificateRequestState_INITIALIZED, certreq.Status)
	s.Empty(certreq.Created)
	s.Empty(certreq.Modified)

	// Inject bufconn connection into the store
	require.NoError(s.grpc.Connect(context.Background()))
	defer s.grpc.Close()

	db, err := store.NewMock(s.grpc.Conn)
	require.NoError(err)

	// Initially there should be no CertReqs
	iter := db.ListCertReqs(context.Background())
	require.False(iter.Next())
	iter.Release()

	// Should get a not found error trying to retrieve a CertReq that doesn't exist
	_, err = db.RetrieveCertReq(context.Background(), "12345")
	require.EqualError(err, storeerrors.ErrEntityNotFound.Error())

	// Attempt to Create the CertReq
	id, err := db.CreateCertReq(context.Background(), certreq)
	s.NoError(err)

	// Attempt to Retrieve the CertReq
	crr, err := db.RetrieveCertReq(context.Background(), id)
	s.NoError(err)
	s.Equal(id, crr.Id)
	s.NotEmpty(crr.Created)
	s.Equal(crr.Modified, crr.Created)
	s.Equal(certreq.Vasp, crr.Vasp)
	s.Equal(certreq.CommonName, crr.CommonName)

	// Attempt to save a certificate request with an ID on it
	icrr := &models.CertificateRequest{
		Id:         uuid.New().String(),
		Vasp:       crr.Vasp,
		CommonName: crr.CommonName,
		Status:     models.CertificateRequestState_INITIALIZED,
	}
	_, err = db.CreateCertReq(context.Background(), icrr)
	s.ErrorIs(err, storeerrors.ErrIDAlreadySet)

	// Sleep for a second to roll over the clock for the modified time stamp
	time.Sleep(1 * time.Second)

	// Update the CertReq
	crr.Status = models.CertificateRequestState_COMPLETED
	err = db.UpdateCertReq(context.Background(), crr)
	s.NoError(err)

	crr, err = db.RetrieveCertReq(context.Background(), id)
	s.NoError(err)
	s.Equal(id, crr.Id)
	s.Equal(models.CertificateRequestState_COMPLETED, crr.Status)
	s.NotEmpty(crr.Modified)
	s.NotEqual(crr.Modified, crr.Created)

	// Attempt to update a certificate request with no Id on it
	certreq.Id = ""
	s.ErrorIs(db.UpdateCertReq(context.Background(), certreq), storeerrors.ErrIncompleteRecord)

	// Delete the CertReq
	err = db.DeleteCertReq(context.Background(), id)
	s.NoError(err)
	crr, err = db.RetrieveCertReq(context.Background(), id)
	s.ErrorIs(err, storeerrors.ErrEntityNotFound)
	s.Empty(crr)

	// Add a few more certificate requests
	for i := 0; i < 10; i++ {
		crr := &models.CertificateRequest{
			Vasp:       uuid.New().String(),
			CommonName: fmt.Sprintf("trisa%d.example.com", i+1),
			Status:     models.CertificateRequestState_COMPLETED,
		}
		_, err := db.CreateCertReq(context.Background(), crr)
		s.NoError(err)
	}

	// Test listing all of the certificates
	reqs, err := db.ListCertReqs(context.Background()).All()
	s.NoError(err)
	s.Len(reqs, 10)

	// Test Prev() and Next() interactions
	iter = db.ListCertReqs(context.Background())
	require.False(iter.Prev(), "should move behind the first CertReq")
	require.True(iter.Next(), "should move to the first CertReq")
	first, err := iter.CertReq()
	require.NoError(err)
	require.NotNil(first)
	require.True(iter.Next(), "should move to the second CertReq")
	second, err := iter.CertReq()
	require.NoError(err)
	require.NotNil(second)
	require.NotEqual(first.Id, second.Id)

	// Consume the rest of the iterator
	for iter.Next() {
		certreq, err := iter.CertReq()
		require.NoError(err)
		require.NotNil(certreq)
	}
	require.NoError(iter.Error())
	iter.Release()

	// Create enough CertReqs to exceed the page size
	for i := 0; i < 100; i++ {
		crr := &models.CertificateRequest{
			Vasp:       uuid.New().String(),
			CommonName: fmt.Sprintf("trisa%d.example.com", i+1),
			Status:     models.CertificateRequestState_COMPLETED,
		}
		_, err := db.CreateCertReq(context.Background(), crr)
		s.NoError(err)
	}

	// Test listing all of the CertReqs
	reqs, err = db.ListCertReqs(context.Background()).All()
	require.NoError(err)
	require.Len(reqs, 110)
}

func (s *trtlStoreTestSuite) TestAnnouncementStore() {
	require := s.Require()

	// Load the announcement month record from testdata
	data, err := os.ReadFile("../testdata/announcements.json")
	require.NoError(err)

	month := &bff.AnnouncementMonth{}
	err = protojson.Unmarshal(data, month)
	require.NoError(err)

	// Verify the announcement month is loaded correctly
	require.NotEmpty(month.Date)
	require.NotEmpty(month.Announcements)
	require.Empty(month.Created)
	require.Empty(month.Modified)

	// Inject bufconn connection into the store
	require.NoError(s.grpc.Connect(context.Background()))
	defer s.grpc.Close()

	db, err := store.NewMock(s.grpc.Conn)
	require.NoError(err)

	// Create the announcement month
	require.NoError(db.UpdateAnnouncementMonth(context.Background(), month))

	// Attempt to Retrieve the announcement month
	m, err := db.RetrieveAnnouncementMonth(context.Background(), month.Date)
	require.NoError(err)
	require.Equal(month.Date, m.Date)
	require.NotEmpty(m.Created)
	require.Equal(m.Modified, m.Created)
	require.Len(m.Announcements, len(month.Announcements))

	// Attempt to Retrieve a non-existent announcement month
	_, err = db.RetrieveAnnouncementMonth(context.Background(), "")
	require.ErrorIs(err, storeerrors.ErrEntityNotFound)
	_, err = db.RetrieveAnnouncementMonth(context.Background(), "2022-01-01")
	require.Error(err)
	_, err = db.RetrieveAnnouncementMonth(context.Background(), "2021-01")
	require.ErrorIs(err, storeerrors.ErrEntityNotFound)

	// Attempt to save an announcement month without a date on it
	month.Date = ""
	err = db.UpdateAnnouncementMonth(context.Background(), month)
	require.ErrorIs(err, storeerrors.ErrIncompleteRecord)

	// Sleep to advance the clock for the modified timestamp
	time.Sleep(1 * time.Millisecond)

	// Update the announcement month
	m.Announcements[0].Title = "Happy New Year!"
	err = db.UpdateAnnouncementMonth(context.Background(), m)
	require.NoError(err)

	m, err = db.RetrieveAnnouncementMonth(context.Background(), m.Date)
	require.NoError(err)
	require.Equal("Happy New Year!", m.Announcements[0].Title)
	require.NotEmpty(m.Modified)
	require.NotEqual(m.Modified, m.Created)

	// Add another announcement month
	month = &bff.AnnouncementMonth{
		Date: "2022-02",
		Announcements: []*bff.Announcement{
			{
				Title:    "Happy Groundhog Day",
				Body:     "The groundhog saw his shadow, so we have six more weeks of winter.",
				PostDate: "2022-02-02",
				Author:   "phil@punxsutawney.com",
			},
		},
	}
	require.NoError(db.UpdateAnnouncementMonth(context.Background(), month))

	// Test that we can still retrieve both months
	january, err := db.RetrieveAnnouncementMonth(context.Background(), "2022-01")
	require.NoError(err)
	require.Equal("Happy New Year!", january.Announcements[0].Title)

	february, err := db.RetrieveAnnouncementMonth(context.Background(), "2022-02")
	require.NoError(err)
	require.Equal("Happy Groundhog Day", february.Announcements[0].Title)

	// Delete an announcement month
	require.NoError(db.DeleteAnnouncementMonth(context.Background(), "2022-01"))

	// Should not be able to retrieve the deleted announcement month
	_, err = db.RetrieveAnnouncementMonth(context.Background(), "2022-01")
	require.ErrorIs(err, storeerrors.ErrEntityNotFound)
}

func (s *trtlStoreTestSuite) TestActivityStore() {
	require := s.Require()

	// Load the activity month record from testdata
	data, err := os.ReadFile("../testdata/activity.json")
	require.NoError(err)

	month := &bff.ActivityMonth{}
	err = protojson.Unmarshal(data, month)
	require.NoError(err)

	// Inject bufconn connection into the store
	require.NoError(s.grpc.Connect(context.Background()))
	defer s.grpc.Close()

	// Verify the activity month is loaded correctly
	require.NotEmpty(month.Date)
	require.NotEmpty(month.Days)
	require.Empty(month.Created)
	require.Empty(month.Modified)

	db, err := store.NewMock(s.grpc.Conn)
	require.NoError(err)

	// Create the activity month
	require.NoError(db.UpdateActivityMonth(context.Background(), month))

	// Attempt to Retrieve the activity month
	m, err := db.RetrieveActivityMonth(context.Background(), month.Date)
	require.NoError(err)
	require.Equal(month.Date, m.Date)
	require.NotEmpty(m.Created)
	require.Equal(m.Modified, m.Created)
	require.Len(m.Days, len(month.Days))

	// Attempt to Retrieve a non-existent activity month
	_, err = db.RetrieveActivityMonth(context.Background(), "")
	require.ErrorIs(err, storeerrors.ErrEntityNotFound)
	_, err = db.RetrieveActivityMonth(context.Background(), "2022-01-01")
	require.Error(err)
	_, err = db.RetrieveActivityMonth(context.Background(), "2021-01")
	require.ErrorIs(err, storeerrors.ErrEntityNotFound)

	// Attempt to save an activity month without a date on it
	month.Date = ""
	err = db.UpdateActivityMonth(context.Background(), month)
	require.ErrorIs(err, storeerrors.ErrIncompleteRecord)

	// Sleep to advance the clock for the modified timestamp
	time.Sleep(1 * time.Millisecond)

	// Update the activity month
	m.Days[0].Activity.Mainnet["lookup"] = 10
	err = db.UpdateActivityMonth(context.Background(), m)
	require.NoError(err)

	m, err = db.RetrieveActivityMonth(context.Background(), m.Date)
	require.NoError(err)
	require.Equal(uint64(10), m.Days[0].Activity.Mainnet["lookup"])
	require.NotEmpty(m.Modified)
	require.NotEqual(m.Modified, m.Created)

	// Add another activity month
	month = &bff.ActivityMonth{
		Date: "2023-02",
		Days: []*bff.ActivityDay{
			{
				Date: "2023-02-01",
				Activity: &bff.ActivityCount{
					Testnet: map[string]uint64{
						"lookup": 20,
					},
				},
			},
		},
	}
	require.NoError(db.UpdateActivityMonth(context.Background(), month))

	// Test that we can still retrieve both months
	january, err := db.RetrieveActivityMonth(context.Background(), "2023-01")
	require.NoError(err)
	require.Equal(uint64(1), january.Days[0].Activity.Testnet["lookup"])

	february, err := db.RetrieveActivityMonth(context.Background(), "2023-02")
	require.NoError(err)
	require.Equal(uint64(20), february.Days[0].Activity.Testnet["lookup"])

	// Delete an activity month
	require.NoError(db.DeleteActivityMonth(context.Background(), "2023-01"))

	// Should not be able to retrieve the deleted activity month
	_, err = db.RetrieveActivityMonth(context.Background(), "2023-01")
	require.ErrorIs(err, storeerrors.ErrEntityNotFound)
}

func (s *trtlStoreTestSuite) TestOrganizationStore() {
	require := s.Require()

	// Inject bufconn connection into the store
	require.NoError(s.grpc.Connect(context.Background()))
	defer s.grpc.Close()

	db, err := store.NewMock(s.grpc.Conn)
	require.NoError(err)

	// Create a new organization in the database
	org := &bff.Organization{}
	id, err := db.CreateOrganization(context.Background(), org)
	require.NoError(err)

	// Verify that the created record has an ID and timestamps
	require.NotEmpty(org.Id)
	require.Equal(org.Id, id)
	require.NotEmpty(org.Created)
	require.Equal(org.Modified, org.Created)

	// Retrieve the organization by UUID
	uu, err := bff.ParseOrgID(org.Id)
	require.NoError(err)
	o, err := db.RetrieveOrganization(context.Background(), uu)
	require.NoError(err)
	require.True(proto.Equal(org, o), "retrieved organization does not match created organization")

	// Attempt to retrieve a non-existent organization
	_, err = db.RetrieveOrganization(context.Background(), uuid.Nil)
	require.ErrorIs(err, storeerrors.ErrEntityNotFound)
	_, err = db.RetrieveOrganization(context.Background(), uuid.New())
	require.ErrorIs(err, storeerrors.ErrEntityNotFound)

	// Sleep to advance the clock for the modified timestamp
	time.Sleep(1 * time.Millisecond)

	// Update the organization
	org.Name = "Alice Corp"
	err = db.UpdateOrganization(context.Background(), org)
	require.NoError(err)

	o, err = db.RetrieveOrganization(context.Background(), uu)
	require.NoError(err)
	require.Equal("Alice Corp", o.Name)
	require.NotEmpty(o.Modified)
	require.NotEqual(o.Modified, o.Created)

	// Attempt to update an organization with no Id on it
	org.Id = ""
	require.ErrorIs(db.UpdateOrganization(context.Background(), org), storeerrors.ErrEntityNotFound)

	// Delete the organization
	err = db.DeleteOrganization(context.Background(), uu)
	require.NoError(err)
	_, err = db.RetrieveOrganization(context.Background(), uu)
	require.ErrorIs(err, storeerrors.ErrEntityNotFound)

	// Create a few organizations
	for i := 0; i < 10; i++ {
		org := &bff.Organization{}
		id, err := db.CreateOrganization(context.Background(), org)
		require.NoError(err)
		require.NotEmpty(id)
	}

	// Test Prev() and Next() interactions
	actualOrgs := 0
	iter := db.ListOrganizations(context.Background())
	require.False(iter.Prev(), "should move behind the first Organization")
	require.True(iter.Next(), "should move to the first Organization")
	first, err := iter.Organization()
	require.NoError(err)
	require.NotNil(first)
	actualOrgs++
	require.True(iter.Next(), "should move to the second Organization")
	second, err := iter.Organization()
	require.NoError(err)
	require.NotNil(second)
	require.NotEqual(first.Id, second.Id, "organizations should have unique IDs")
	actualOrgs++

	// Consume the rest of the iterator
	for iter.Next() {
		org, err := iter.Organization()
		require.NoError(err)
		require.NotNil(org)
		actualOrgs++
	}
	require.NoError(iter.Error())
	iter.Release()
	require.Equal(actualOrgs, 10, "iterator returned the wrong number of organizations")

	// Create enough organizations to exceed the default page size
	for i := 0; i < 100; i++ {
		org := &bff.Organization{}
		_, err := db.CreateOrganization(context.Background(), org)
		require.NoError(err)
	}

	// Consume all of the organizations
	iter = db.ListOrganizations(context.Background())
	actualOrgs = 0
	for iter.Next() {
		org, err := iter.Organization()
		require.NoError(err)
		require.NotNil(org)
		actualOrgs++
	}
	require.NoError(iter.Error())
	iter.Release()

	require.Equal(actualOrgs, 110, "iterator returned the wrong number of organizations")
}

func (s *trtlStoreTestSuite) TestContactStore() {
	require := s.Require()

	// Inject bufconn connection into the store
	require.NoError(s.grpc.Connect(context.Background()))
	defer s.grpc.Close()

	// Connect a mock store
	db, err := store.NewMock(s.grpc.Conn)
	require.NoError(err)

	// Make sure create errors with a nil contact
	email, err := db.CreateContact(context.Background(), nil)
	require.Empty(email)
	require.Equal(err, storeerrors.ErrIncompleteRecord)

	// Make sure create errors with a contact with an empty email
	contact := &models.Contact{
		Email: "",
	}
	email, err = db.CreateContact(context.Background(), contact)
	require.Empty(email)
	require.Equal(err, storeerrors.ErrIncompleteRecord)

	// Create a valid contact
	contact = &models.Contact{
		Email:      "testemail",
		Vasps:      []string{"foo", "bar"},
		Verified:   false,
		Token:      "testtoken",
		VerifiedOn: "",
	}
	email, err = db.CreateContact(context.Background(), contact)
	require.Equal(email, "testemail")
	require.NoError(err)

	// Make sure retrieve errors with an empty email
	var con *models.Contact
	con, err = db.RetrieveContact(context.Background(), "")
	require.Nil(con)
	require.Equal(err, storeerrors.ErrEntityNotFound)

	// Make sure retrieve throws the proper error when a contact is not found
	con, err = db.RetrieveContact(context.Background(), "wrongemail")
	require.Nil(con)
	require.Equal(err, storeerrors.ErrEntityNotFound)

	// Retrieve the created contact
	con, err = db.RetrieveContact(context.Background(), "testemail")
	require.Equal(con.Vasps, contact.Vasps)
	require.Equal(con.Verified, contact.Verified)
	require.Equal(con.Token, contact.Token)
	require.NotEmpty(con.Created)
	require.NotEmpty(con.Modified)
	require.NoError(err)

	// Make sure update errors with a nil contact
	err = db.UpdateContact(context.Background(), nil)
	require.Equal(err, storeerrors.ErrIncompleteRecord)

	// Make sure update errors with a contact with an empty email
	contact = &models.Contact{
		Email: "",
	}
	err = db.UpdateContact(context.Background(), contact)
	require.Equal(err, storeerrors.ErrIncompleteRecord)

	// Sleep to allow time for the modified field to change
	time.Sleep(time.Second)

	// Properly update the contact
	contact = &models.Contact{
		Email:      "testemail",
		Vasps:      []string{"bar", "foo"},
		Verified:   true,
		Token:      "newtoken",
		VerifiedOn: "",
		Created:    con.Created,
	}
	err = db.UpdateContact(context.Background(), contact)
	require.NoError(err)

	// Retrieve the updated contact
	var updatedCon *models.Contact
	updatedCon, err = db.RetrieveContact(context.Background(), "testemail")
	require.Equal(updatedCon.Vasps, contact.Vasps)
	require.Equal(updatedCon.Verified, contact.Verified)
	require.Equal(updatedCon.Token, contact.Token)
	require.NoError(err)
	require.Equal(con.Created, updatedCon.Created)
	require.NotEqual(con.Modified, updatedCon.Modified)

	// Make sure delete errors with an empty email
	err = db.DeleteContact(context.Background(), "")
	require.Equal(err, storeerrors.ErrEntityNotFound)

	// Make sure delete throws an error when the contact to delete isn't found
	err = db.DeleteContact(context.Background(), "wrongemail")
	require.EqualError(err, "rpc error: code = NotFound desc = not found")

	// Delete the created contact
	err = db.DeleteContact(context.Background(), "testemail")
	require.NoError(err)

	// Make sure the contact was deleted
	con, err = db.RetrieveContact(context.Background(), "testemail")
	require.Nil(con)
	require.Equal(err, storeerrors.ErrEntityNotFound)
}

func (s *trtlStoreTestSuite) TestEmailStore() {
	ctx := context.Background()

	// Inject bufconn connection into the store
	s.NoError(s.grpc.Connect(context.Background()))
	defer s.grpc.Close()

	// Connect a mock store
	db, err := store.NewMock(s.grpc.Conn)
	s.NoError(err)

	s.Run("List", func() {
		cleanup := func() {
			iter := db.ListEmails(ctx)
			defer iter.Release()

			for iter.Next() {
				email, _ := iter.Email()
				db.DeleteEmail(ctx, email.Email)
			}
		}

		// Ensure the database is empty and emptied at the end of the test.
		cleanup()
		defer cleanup()

		// Test list empty database
		iter := db.ListEmails(ctx)
		s.False(iter.Next(), "expected next to be false in empty database")
		s.False(iter.Prev(), "expected prev to be false in empty database")
		iter.Release()
		s.NoError(iter.Error(), "expected no error during iteration")

		// Create some emails in the database
		for i := 1; i < 128; i++ {
			email := &models.Email{Email: fmt.Sprintf("person%03d@example.com", i), Token: "foo"}
			err := db.UpdateEmail(ctx, email)
			s.NoError(err, "could not insert email into database")
		}

		count := 0
		iter = db.ListEmails(ctx)
		for iter.Next() {
			_, err := iter.Email()
			s.NoError(err, "could not load email from iterator")
			count++
		}

		iter.Release()
		s.NoError(iter.Error(), "expected no error during iteration")
		s.Equal(127, count, "wrong iterator count returned")
	})

	s.Run("HappyCRUD", func() {
		email := &models.Email{
			Email:    "Gary Vespers <Gary.Vespers@example.com>",
			Vasps:    []string{"06c74ef0-0b5b-4df1-8fc5-53f6bb7044af"},
			Verified: false,
			Token:    "ZPr7O4YRes30ie5SZjjqGwkETM2qHD1nfUIbcIecXTg",
		}

		// Test create email record
		key, err := db.CreateEmail(ctx, email)
		s.NoError(err, "could not create valid email record")
		s.Equal("gary.vespers@example.com", key, "unexpected key returned")
		s.NotEmpty(email.Created, "expected the created timestamp to be added")
		s.NotEmpty(email.Modified, "expected the modified timestamp to be added")
		s.Equal("gary.vespers@example.com", email.Email, "expected the email to be normalized")
		s.Equal("Gary Vespers", email.Name, "expected the name to be parsed from the email")

		// Test retrieve email record
		record, err := db.RetrieveEmail(ctx, "Gary Vespers <Gary.Vespers@example.com>")
		s.NoError(err)
		s.True(proto.Equal(email, record))

		// Sleep to allow time for the modified field to change
		time.Sleep(time.Second)

		// Test updating the email record
		record.Verified = true
		record.VerifiedOn = time.Now().Format(time.RFC3339)
		record.Token = ""

		err = db.UpdateEmail(ctx, record)
		s.NoError(err, "could not update email record")
		s.NotEqual(record.Modified, email.Modified, "expected modified timestamp to be updated")

		updated, err := db.RetrieveEmail(ctx, "gary.vespers@example.com")
		s.NoError(err, "could not retrieve updated email")
		s.False(proto.Equal(email, updated))
		s.True(proto.Equal(record, updated))

		// Test deleting the email record
		err = db.DeleteEmail(ctx, "Gary.Vespers@example.com")
		s.NoError(err, "could not delete email record")

		// Ensure the email was deleted
		_, err = db.RetrieveEmail(ctx, "gary.vespers@example.com")
		s.ErrorIs(err, storeerrors.ErrEntityNotFound)
	})

	s.Run("CreateEmptyEmailError", func() {
		// Make sure create errors with a nil email
		email, err := db.CreateEmail(ctx, nil)
		s.Empty(email)
		s.ErrorIs(err, storeerrors.ErrIncompleteRecord)

		// Make sure create errors with an email empty string
		email, err = db.CreateEmail(ctx, &models.Email{Email: ""})
		s.Empty(email)
		s.ErrorIs(err, storeerrors.ErrIncompleteRecord)
	})

	s.Run("CreateDuplicate", func() {
		// Cannot create a duplicate record
		email := &models.Email{
			Name:  "Barb Fabrittle",
			Email: "barb@example.com",
			Token: "ZPr7O4YRes30ie5SZjjqGwkETM2qHD1nfUIbcIecXTg",
		}

		_, err := db.CreateEmail(ctx, email)
		s.NoError(err, "could not create first record")

		_, err = db.CreateEmail(ctx, email)
		s.ErrorIs(err, storeerrors.ErrEmailExists, "expected email exists error")
	})

	s.Run("CreateNameParsing", func() {
		// NOTE: test cases need unique emails since database isn't reset between tests
		testCases := []struct {
			email         *models.Email
			expectedName  string
			expectedEmail string
		}{
			{
				&models.Email{Name: "Frank Shadypants", Email: "Kelly Clarkberry <kelly.clarkberry@example.com>", Token: "foo"},
				"Frank Shadypants",
				"kelly.clarkberry@example.com",
			},
			{
				&models.Email{Name: "", Email: "Gillian.Redbottom@example.com", Token: "foo"},
				"",
				"gillian.redbottom@example.com",
			},
			{
				&models.Email{Name: "", Email: "Edward Boilermaker <EDWARD@BOILERMAKER.IO>", Token: "foo"},
				"Edward Boilermaker",
				"edward@boilermaker.io",
			},
		}

		for i, tc := range testCases {
			_, err := db.CreateEmail(ctx, tc.email)
			s.NoError(err, "test case %d failed", i)
			s.Equal(tc.expectedName, tc.email.Name, "test case %d failed", i)
			s.Equal(tc.expectedEmail, tc.email.Email, "test case %d failed", i)
		}
	})

	s.Run("RetrieveEmptyEmailError", func() {
		email, err := db.RetrieveEmail(ctx, "")
		s.Empty(email)
		s.ErrorIs(err, storeerrors.ErrEntityNotFound)
	})

	s.Run("RetrieveNotFound", func() {
		email, err := db.RetrieveEmail(ctx, "thisemaildoesnot@exist.com")
		s.Empty(email)
		s.ErrorIs(err, storeerrors.ErrEntityNotFound)
	})

	s.Run("UpdateEmptyEmailError", func() {
		// Make sure update errors with a nil email
		err := db.UpdateEmail(ctx, nil)
		s.ErrorIs(err, storeerrors.ErrIncompleteRecord)

		// Make sure update errors with an email empty string
		err = db.UpdateEmail(ctx, &models.Email{Email: ""})
		s.ErrorIs(err, storeerrors.ErrIncompleteRecord)
	})

	s.Run("DeleteEmptyEmailError", func() {
		// Make sure delete errors with an email empty string
		err := db.DeleteEmail(ctx, "")
		s.ErrorIs(err, storeerrors.ErrEntityNotFound)
	})
}

func (s *trtlStoreTestSuite) TestDirectoryContactStore() {
	ctx := context.Background()
	require := s.Require()

	// Inject bufconn connection into the store
	s.NoError(s.grpc.Connect(context.Background()))
	defer s.grpc.Close()

	// Connect a mock store
	db, err := store.NewMock(s.grpc.Conn)
	s.NoError(err)

	// Load fixtures into database for tests
	data, err := os.ReadFile("../testdata/altvasp.json")
	require.NoError(err)

	vasp := &pb.VASP{}
	err = protojson.Unmarshal(data, vasp)
	require.NoError(err)

	vaspID, err := db.CreateVASP(ctx, vasp)
	require.NoError(err, "could not create vasp fixture")

	admin := &models.Email{Name: vasp.Contacts.Administrative.Name, Email: vasp.Contacts.Administrative.Email, Verified: true, VerifiedOn: time.Now().Format(time.RFC3339)}
	tech := &models.Email{Name: vasp.Contacts.Technical.Name, Email: vasp.Contacts.Technical.Email, Verified: true, VerifiedOn: time.Now().Format(time.RFC3339)}
	legal := &models.Email{Name: vasp.Contacts.Legal.Name, Email: vasp.Contacts.Legal.Email, Verified: true, VerifiedOn: time.Now().Format(time.RFC3339)}
	billing := &models.Email{Name: vasp.Contacts.Billing.Name, Email: vasp.Contacts.Billing.Email, Verified: true, VerifiedOn: time.Now().Format(time.RFC3339)}

	for _, record := range []*models.Email{admin, tech, legal, billing} {
		err = db.UpdateEmail(ctx, record)
		require.NoError(err, "could not create associated contact email")
	}

	defer func() {
		require.NoError(db.DeleteVASP(ctx, vaspID), "could not delete darlene vasp")
		for _, record := range []*models.Email{admin, tech, legal, billing} {
			err = db.DeleteEmail(ctx, record.Email)
			require.NoError(err, "could not delete associated contact email")
		}
	}()

	s.Run("VASPContacts", func() {
		contacts, err := db.VASPContacts(ctx, vasp)
		require.NoError(err, "could not get vasp contacts for darlene")
		require.Equal(vasp.Id, contacts.VASP, "expected contacts VASP to match VASP")
		require.Equal(vasp.Contacts, contacts.Contacts, "expected the contacts to match the VASP")
		require.Len(contacts.Emails, 4, "expected two emails retrieved")
	})

	s.Run("RetrieveVASPContacts", func() {
		contacts, err := db.RetrieveVASPContacts(ctx, vaspID)
		require.NoError(err, "could not retrieve vasp contacts for darlene")
		require.Equal(vasp.Id, contacts.VASP, "expected contacts VASP to match VASP")
		require.True(proto.Equal(vasp.Contacts, contacts.Contacts), "expected the contacts to match the VASP")
		require.Len(contacts.Emails, 4, "expected two emails retrieved")
	})
}

func createVASPs(db *store.Store, num, startIndex int) error {
	countries := []string{"TV", "KY", "CC", "LT", "EH", "SC", "NU"}
	bcats := []pb.BusinessCategory{pb.BusinessCategoryBusiness, pb.BusinessCategoryNonCommercial, pb.BusinessCategoryPrivate}

	for i := 0; i < num; i++ {
		country := countries[i%len(countries)]
		vasp := &pb.VASP{
			Entity: &ivms101.LegalPerson{
				Name: &ivms101.LegalPersonName{
					NameIdentifiers: []*ivms101.LegalPersonNameId{
						{
							LegalPersonName:               fmt.Sprintf("Test VASP %04X", i+startIndex),
							LegalPersonNameIdentifierType: ivms101.LegalPersonLegal,
						},
					},
				},
				CountryOfRegistration: country,
			},
			Website:          fmt.Sprintf("https://test%04X.net/", i+startIndex),
			CommonName:       fmt.Sprintf("trisa%04d.test.net", i+startIndex),
			BusinessCategory: bcats[i%len(bcats)],
		}

		if _, err := db.CreateVASP(context.Background(), vasp); err != nil {
			return err
		}
	}
	return nil
}

// TODO: Add Announcements and Organization tests

func deleteVASPs(db *store.Store) error {
	n := 0
	iter := db.ListVASPs(context.Background())
	for iter.Next() {
		vasp, err := iter.VASP()
		if err != nil {
			iter.Release()
			return err
		}

		if err := db.DeleteVASP(context.Background(), vasp.Id); err != nil {
			iter.Release()
			return err
		}

		n++
	}

	iter.Release()
	if err := iter.Error(); err != nil {
		return err
	}

	// TODO: do better at managing empty indices
	db.DeleteIndices()
	return nil
}
