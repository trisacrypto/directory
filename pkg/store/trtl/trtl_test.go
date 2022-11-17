package trtl_test

import (
	"context"
	"fmt"
	"io/ioutil"
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
	// NOTE: ConsoleLog MUST be false otherwise this will be overriden
	logger.Discard()

	var err error
	s.tmpdb, err = ioutil.TempDir("../testdata", "db-*")
	require.NoError(err)

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
	data, err := ioutil.ReadFile("../testdata/vasp.json")
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
	iter := db.ListVASPs()
	require.False(iter.Next())
	iter.Release()

	// Should get a not found error trying to retrieve a VASP that doesn't exist
	_, err = db.RetrieveVASP("12345")
	require.EqualError(err, storeerrors.ErrEntityNotFound.Error())

	// Attempt to Create the VASP
	id, err := db.CreateVASP(alice)
	require.NoError(err)
	require.NotEmpty(id)

	// Should not be able to create a duplicate VASP
	id2, err := db.CreateVASP(alice)
	require.EqualError(err, storeerrors.ErrDuplicateEntity.Error())
	require.Empty(id2)

	// Attempt to Retrieve the VASP
	alicer, err := db.RetrieveVASP(id)
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
	err = db.UpdateVASP(alicer)
	require.NoError(err)

	alicer, err = db.RetrieveVASP(id)
	require.NoError(err)
	require.Equal(id, alicer.Id)
	require.NotEmpty(alicer.LastUpdated)
	require.NotEqual(alicer.FirstListed, alicer.LastUpdated)
	require.NotEmpty(alicer.Version)
	require.Equal(uint64(2), alicer.Version.Version)
	require.Equal(alicer.VerificationStatus, pb.VerificationState_VERIFIED)

	// Delete the VASP
	err = db.DeleteVASP(id)
	require.NoError(err)
	alicer, err = db.RetrieveVASP(id)
	require.ErrorIs(err, storeerrors.ErrEntityNotFound)
	require.Empty(alicer)

	// Add a few more VASPs
	err = createVASPs(db, 10, 1)
	require.NoError(err)

	// Test listing all of the VASPs
	reqs, err := db.ListVASPs().All()
	require.NoError(err)
	require.Len(reqs, 10)

	// Test seeking to a specific VASP
	key := reqs[5].Id
	iter = db.ListVASPs()
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
	iter = db.ListVASPs()
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
	reqs, err = db.ListVASPs().All()
	require.NoError(err)
	require.Len(reqs, 110)

	// Cleanup database
	require.NoError(deleteVASPs(db), "could not cleanup database")
}

func (s *trtlStoreTestSuite) TestCertificateStore() {
	require := s.Require()

	// Load the VASP record from testdata
	data, err := ioutil.ReadFile("../testdata/cert.json")
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
	iter := db.ListCerts()
	require.False(iter.Next())
	iter.Release()

	// Should get a not found error trying to retrieve a Cert that doesn't exist
	_, err = db.RetrieveCert("12345")
	require.EqualError(err, storeerrors.ErrEntityNotFound.Error())

	// Attempt to Create the Cert
	id, err := db.CreateCert(cert)
	s.NoError(err)

	// Attempt to Retrieve the Cert
	crr, err := db.RetrieveCert(id)
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
	_, err = db.CreateCert(icrr)
	s.ErrorIs(err, storeerrors.ErrIDAlreadySet)

	// Sleep for a second to roll over the clock for the modified time stamp
	time.Sleep(1 * time.Second)

	// Update the Cert
	crr.Status = models.CertificateState_REVOKED
	err = db.UpdateCert(crr)
	s.NoError(err)

	crr, err = db.RetrieveCert(id)
	s.NoError(err)
	s.Equal(id, crr.Id)
	s.Equal(models.CertificateState_REVOKED, crr.Status)

	// Attempt to update a certificate request with no Id on it
	cert.Id = ""
	s.ErrorIs(db.UpdateCert(cert), storeerrors.ErrIncompleteRecord)

	// Delete the Cert
	err = db.DeleteCert(id)
	s.NoError(err)
	crr, err = db.RetrieveCert(id)
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
		_, err := db.CreateCert(crr)
		s.NoError(err)
	}

	// Test listing all of the certificates
	certs, err := db.ListCerts().All()
	s.NoError(err)
	s.Len(certs, 10)

	// Test Prev() and Next() interactions
	iter = db.ListCerts()
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
		_, err := db.CreateCert(crr)
		s.NoError(err)
	}

	// Test listing all of the Cert
	certs, err = db.ListCerts().All()
	require.NoError(err)
	require.Len(certs, 110)
}

func (s *trtlStoreTestSuite) TestCertificateRequestStore() {
	require := s.Require()

	// Load the VASP record from testdata
	data, err := ioutil.ReadFile("../testdata/certreq.json")
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
	iter := db.ListCertReqs()
	require.False(iter.Next())
	iter.Release()

	// Should get a not found error trying to retrieve a CertReq that doesn't exist
	_, err = db.RetrieveCertReq("12345")
	require.EqualError(err, storeerrors.ErrEntityNotFound.Error())

	// Attempt to Create the CertReq
	id, err := db.CreateCertReq(certreq)
	s.NoError(err)

	// Attempt to Retrieve the CertReq
	crr, err := db.RetrieveCertReq(id)
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
	_, err = db.CreateCertReq(icrr)
	s.ErrorIs(err, storeerrors.ErrIDAlreadySet)

	// Sleep for a second to roll over the clock for the modified time stamp
	time.Sleep(1 * time.Second)

	// Update the CertReq
	crr.Status = models.CertificateRequestState_COMPLETED
	err = db.UpdateCertReq(crr)
	s.NoError(err)

	crr, err = db.RetrieveCertReq(id)
	s.NoError(err)
	s.Equal(id, crr.Id)
	s.Equal(models.CertificateRequestState_COMPLETED, crr.Status)
	s.NotEmpty(crr.Modified)
	s.NotEqual(crr.Modified, crr.Created)

	// Attempt to update a certificate request with no Id on it
	certreq.Id = ""
	s.ErrorIs(db.UpdateCertReq(certreq), storeerrors.ErrIncompleteRecord)

	// Delete the CertReq
	err = db.DeleteCertReq(id)
	s.NoError(err)
	crr, err = db.RetrieveCertReq(id)
	s.ErrorIs(err, storeerrors.ErrEntityNotFound)
	s.Empty(crr)

	// Add a few more certificate requests
	for i := 0; i < 10; i++ {
		crr := &models.CertificateRequest{
			Vasp:       uuid.New().String(),
			CommonName: fmt.Sprintf("trisa%d.example.com", i+1),
			Status:     models.CertificateRequestState_COMPLETED,
		}
		_, err := db.CreateCertReq(crr)
		s.NoError(err)
	}

	// Test listing all of the certificates
	reqs, err := db.ListCertReqs().All()
	s.NoError(err)
	s.Len(reqs, 10)

	// Test Prev() and Next() interactions
	iter = db.ListCertReqs()
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
		_, err := db.CreateCertReq(crr)
		s.NoError(err)
	}

	// Test listing all of the CertReqs
	reqs, err = db.ListCertReqs().All()
	require.NoError(err)
	require.Len(reqs, 110)
}

func (s *trtlStoreTestSuite) TestAnnouncementStore() {
	require := s.Require()

	// Load the announcement month record from testdata
	data, err := ioutil.ReadFile("../testdata/announcements.json")
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
	require.NoError(db.UpdateAnnouncementMonth(month))

	// Attempt to Retrieve the announcement month
	m, err := db.RetrieveAnnouncementMonth(month.Date)
	require.NoError(err)
	require.Equal(month.Date, m.Date)
	require.NotEmpty(m.Created)
	require.Equal(m.Modified, m.Created)
	require.Len(m.Announcements, len(month.Announcements))

	// Attempt to Retrieve a non-existent announcement month
	_, err = db.RetrieveAnnouncementMonth("")
	require.ErrorIs(err, storeerrors.ErrEntityNotFound)
	_, err = db.RetrieveAnnouncementMonth("2022-01-01")
	require.Error(err)
	_, err = db.RetrieveAnnouncementMonth("2021-01")
	require.ErrorIs(err, storeerrors.ErrEntityNotFound)

	// Attempt to save an announcement month without a date on it
	month.Date = ""
	err = db.UpdateAnnouncementMonth(month)
	require.ErrorIs(err, storeerrors.ErrIncompleteRecord)

	// Sleep to advance the clock for the modified timestamp
	time.Sleep(1 * time.Millisecond)

	// Update the announcement month
	m.Announcements[0].Title = "Happy New Year!"
	err = db.UpdateAnnouncementMonth(m)
	require.NoError(err)

	m, err = db.RetrieveAnnouncementMonth(m.Date)
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
	require.NoError(db.UpdateAnnouncementMonth(month))

	// Test that we can still retrieve both months
	january, err := db.RetrieveAnnouncementMonth("2022-01")
	require.NoError(err)
	require.Equal("Happy New Year!", january.Announcements[0].Title)

	february, err := db.RetrieveAnnouncementMonth("2022-02")
	require.NoError(err)
	require.Equal("Happy Groundhog Day", february.Announcements[0].Title)

	// Delete an announcement month
	require.NoError(db.DeleteAnnouncementMonth("2022-01"))

	// Should not be able to retrieve the deleted announcement month
	_, err = db.RetrieveAnnouncementMonth("2022-01")
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
	id, err := db.CreateOrganization(org)
	require.NoError(err)

	// Verify that the created record has an ID and timestamps
	require.NotEmpty(org.Id)
	require.Equal(org.Id, id)
	require.NotEmpty(org.Created)
	require.Equal(org.Modified, org.Created)

	// Retrieve the organization by UUID
	uu, err := bff.ParseOrgID(org.Id)
	require.NoError(err)
	o, err := db.RetrieveOrganization(uu)
	require.NoError(err)
	require.True(proto.Equal(org, o), "retrieved organization does not match created organization")

	// Attempt to retrieve a non-existent organization
	_, err = db.RetrieveOrganization(uuid.Nil)
	require.ErrorIs(err, storeerrors.ErrEntityNotFound)
	_, err = db.RetrieveOrganization(uuid.New())
	require.ErrorIs(err, storeerrors.ErrEntityNotFound)

	// Sleep to advance the clock for the modified timestamp
	time.Sleep(1 * time.Millisecond)

	// Update the organization
	org.Name = "Alice Corp"
	err = db.UpdateOrganization(org)
	require.NoError(err)

	o, err = db.RetrieveOrganization(uu)
	require.NoError(err)
	require.Equal("Alice Corp", o.Name)
	require.NotEmpty(o.Modified)
	require.NotEqual(o.Modified, o.Created)

	// Attempt to update an organization with no Id on it
	org.Id = ""
	require.ErrorIs(db.UpdateOrganization(org), storeerrors.ErrEntityNotFound)

	// Delete the organization
	err = db.DeleteOrganization(uu)
	require.NoError(err)
	_, err = db.RetrieveOrganization(uu)
	require.ErrorIs(err, storeerrors.ErrEntityNotFound)
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

		if _, err := db.CreateVASP(vasp); err != nil {
			return err
		}
	}
	return nil
}

// TODO: Add Announcements and Organization tests

func deleteVASPs(db *store.Store) error {
	n := 0
	iter := db.ListVASPs()
	for iter.Next() {
		vasp, err := iter.VASP()
		if err != nil {
			iter.Release()
			return err
		}

		if err := db.DeleteVASP(vasp.Id); err != nil {
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
