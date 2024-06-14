package leveldb

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	bff "github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/models/v1"
	storeerrors "github.com/trisacrypto/directory/pkg/store/errors"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type leveldbTestSuite struct {
	suite.Suite
	path string
	db   *Store
}

func (s *leveldbTestSuite) SetupSuite() {
	// Discard logging from the application to focus on test logs
	// NOTE: ConsoleLog MUST be false otherwise this will be overridden
	logger.Discard()

	// Open the database in a temp directory
	var err error
	s.path = s.T().TempDir()
	s.db, err = Open(s.path)
	s.NoError(err)
}

func (s *leveldbTestSuite) TearDownSuite() {
	// Delete the temp directory when done
	err := os.RemoveAll(s.path)
	s.NoError(err)
	logger.ResetLogger()
}

func TestLevelDB(t *testing.T) {
	suite.Run(t, new(leveldbTestSuite))
}

func (s *leveldbTestSuite) TestDirectoryStore() {
	// Load the VASP record from testdata
	data, err := os.ReadFile("../testdata/vasp.json")
	s.NoError(err)

	alice := &pb.VASP{}
	err = protojson.Unmarshal(data, alice)
	s.NoError(err)

	// Validate the VASP record loaded correctly and is partial
	s.NotEmpty(alice.CommonName)
	s.NotEmpty(alice.TrisaEndpoint)
	s.NoError(alice.Validate(true))
	s.Empty(alice.Id)

	// Attempt to Create the VASP
	id, err := s.db.CreateVASP(context.Background(), alice)
	s.NoError(err)
	s.NotEmpty(id)

	// Attempt to Retrieve the VASP
	alicer, err := s.db.RetrieveVASP(context.Background(), id)
	s.NoError(err)
	s.Equal(id, alicer.Id)
	s.Equal(alicer.FirstListed, alicer.LastUpdated)
	s.NotEmpty(alicer.LastUpdated)
	s.NotEmpty(alicer.Version)
	s.Equal(uint64(1), alicer.Version.Version)

	// Ensure the modification time rolls over to the next second for comparison
	time.Sleep(1 * time.Second)

	// Update the VASP
	alicer.Entity.Name.NameIdentifiers[0].LegalPersonName = "AliceLiteCoin, LLC"
	alicer.VerificationStatus = pb.VerificationState_VERIFIED
	alicer.VerifiedOn = "2021-06-30T10:40:40Z"
	err = s.db.UpdateVASP(context.Background(), alicer)
	s.NoError(err)

	alicer, err = s.db.RetrieveVASP(context.Background(), id)
	s.NoError(err)
	s.Equal(id, alicer.Id)
	s.NotEmpty(alicer.LastUpdated)
	s.NotEqual(alicer.FirstListed, alicer.LastUpdated)
	s.NotEmpty(alicer.Version)
	s.Equal(uint64(2), alicer.Version.Version)
	s.Equal(alicer.VerificationStatus, pb.VerificationState_VERIFIED)

	// Delete the VASP
	err = s.db.DeleteVASP(context.Background(), id)
	s.NoError(err)
	alicer, err = s.db.RetrieveVASP(context.Background(), id)
	s.ErrorIs(err, storeerrors.ErrEntityNotFound)
	s.Empty(alicer)

	// Add a few more VASPs
	for i := 0; i < 10; i++ {
		vasp := &pb.VASP{
			Entity: &ivms101.LegalPerson{
				Name: &ivms101.LegalPersonName{
					NameIdentifiers: []*ivms101.LegalPersonNameId{
						{
							LegalPersonName:               fmt.Sprintf("Test %d", i+1),
							LegalPersonNameIdentifierType: ivms101.LegalPersonLegal,
						},
					},
				},
			},
			CommonName: fmt.Sprintf("trisa%d.test.net", i+1),
		}
		_, err := s.db.CreateVASP(context.Background(), vasp)
		s.NoError(err)
	}

	// Test listing all of the VASPs
	reqs, err := s.db.ListVASPs(context.Background()).All()
	s.NoError(err)
	s.Len(reqs, 10)

	// Test iterating over all the VASPs
	var niters int
	iter := s.db.ListVASPs(context.Background())
	for iter.Next() {
		s.NotEmpty(iter.VASP())
		niters++
	}
	s.NoError(iter.Error())
	iter.Release()
	s.Equal(10, niters)
}

func (s *leveldbTestSuite) TestCertificateStore() {
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
	s.NotEmpty(cert.Details.NotBefore)
	s.NotEmpty(cert.Details.NotAfter)

	// Attempt to Create the Cert
	id, err := s.db.CreateCert(context.Background(), cert)
	s.NoError(err)

	// Attempt to Retrieve the Cert
	crr, err := s.db.RetrieveCert(context.Background(), id)
	s.NoError(err)
	s.Equal(id, crr.Id)
	s.Equal(cert.Request, crr.Request)
	s.Equal(cert.Vasp, crr.Vasp)
	s.Equal(cert.Status, crr.Status)
	s.True(proto.Equal(cert.Details, crr.Details))

	// Attempt to save a certificate with an ID on it
	icrr := &models.Certificate{
		Id:      uuid.New().String(),
		Request: crr.Request,
		Vasp:    crr.Vasp,
		Status:  models.CertificateState_ISSUED,
		Details: crr.Details,
	}
	_, err = s.db.CreateCert(context.Background(), icrr)
	s.ErrorIs(err, storeerrors.ErrIDAlreadySet)

	// Update the Cert
	crr.Status = models.CertificateState_REVOKED
	err = s.db.UpdateCert(context.Background(), crr)
	s.NoError(err)

	crr, err = s.db.RetrieveCert(context.Background(), id)
	s.NoError(err)
	s.Equal(id, crr.Id)
	s.Equal(models.CertificateState_REVOKED, crr.Status)

	// Attempt to update a certificate with no Id on it
	cert.Id = ""
	s.ErrorIs(s.db.UpdateCert(context.Background(), cert), storeerrors.ErrIncompleteRecord)

	// Delete the Cert
	err = s.db.DeleteCert(context.Background(), id)
	s.NoError(err)
	crr, err = s.db.RetrieveCert(context.Background(), id)
	s.ErrorIs(err, storeerrors.ErrEntityNotFound)
	s.Empty(crr)

	// Add a few more certificates
	for i := 0; i < 10; i++ {
		crr := &models.Certificate{
			Request: uuid.New().String(),
			Vasp:    uuid.New().String(),
			Status:  models.CertificateState_ISSUED,
			Details: &pb.Certificate{
				SerialNumber: []byte(uuid.New().String()),
			},
		}
		_, err := s.db.CreateCert(context.Background(), crr)
		s.NoError(err)
	}

	// Test listing all of the certificates
	certs, err := s.db.ListCerts(context.Background()).All()
	s.NoError(err)
	s.Len(certs, 10)

	// Test iterating over all the certificates
	var niters int
	iter := s.db.ListCerts(context.Background())
	for iter.Next() {
		s.NotEmpty(iter.Cert())
		niters++
	}
	s.NoError(iter.Error())
	iter.Release()
	s.Equal(10, niters)
}

func (s *leveldbTestSuite) TestCertificateRequestStore() {
	// Load the certreq record from testdata
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

	// Attempt to Create the CertReq
	id, err := s.db.CreateCertReq(context.Background(), certreq)
	s.NoError(err)

	// Attempt to Retrieve the CertReq
	crr, err := s.db.RetrieveCertReq(context.Background(), id)
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
	_, err = s.db.CreateCertReq(context.Background(), icrr)
	s.ErrorIs(err, storeerrors.ErrIDAlreadySet)

	// Sleep for a second to roll over the clock for the modified time stamp
	time.Sleep(1 * time.Second)

	// Update the CertReq
	crr.Status = models.CertificateRequestState_COMPLETED
	err = s.db.UpdateCertReq(context.Background(), crr)
	s.NoError(err)

	crr, err = s.db.RetrieveCertReq(context.Background(), id)
	s.NoError(err)
	s.Equal(id, crr.Id)
	s.Equal(models.CertificateRequestState_COMPLETED, crr.Status)
	s.NotEmpty(crr.Modified)
	s.NotEqual(crr.Modified, crr.Created)

	// Attempt to update a certificate request with no Id on it
	certreq.Id = ""
	s.ErrorIs(s.db.UpdateCertReq(context.Background(), certreq), storeerrors.ErrIncompleteRecord)

	// Delete the CertReq
	err = s.db.DeleteCertReq(context.Background(), id)
	s.NoError(err)
	crr, err = s.db.RetrieveCertReq(context.Background(), id)
	s.ErrorIs(err, storeerrors.ErrEntityNotFound)
	s.Empty(crr)

	// Add a few more certificate requests
	for i := 0; i < 10; i++ {
		crr := &models.CertificateRequest{
			Vasp:       uuid.New().String(),
			CommonName: fmt.Sprintf("trisa%d.example.com", i+1),
			Status:     models.CertificateRequestState_COMPLETED,
		}
		_, err := s.db.CreateCertReq(context.Background(), crr)
		s.NoError(err)
	}

	// Test listing all of the certificates
	reqs, err := s.db.ListCertReqs(context.Background()).All()
	s.NoError(err)
	s.Len(reqs, 10)

	// Test iterating over all the certificates
	var niters int
	iter := s.db.ListCertReqs(context.Background())
	for iter.Next() {
		s.NotEmpty(iter.CertReq())
		niters++
	}
	s.NoError(iter.Error())
	iter.Release()
	s.Equal(10, niters)
}

func (s *leveldbTestSuite) TestAnnouncementStore() {
	// Load the announcement month record from testdata
	data, err := os.ReadFile("../testdata/announcements.json")
	s.NoError(err)

	month := &bff.AnnouncementMonth{}
	err = protojson.Unmarshal(data, month)
	s.NoError(err)

	// Verify the announcement month is loaded correctly
	s.NotEmpty(month.Date)
	s.NotEmpty(month.Announcements)
	s.Empty(month.Created)
	s.Empty(month.Modified)

	// Create the announcement month
	s.NoError(s.db.UpdateAnnouncementMonth(context.Background(), month))

	// Attempt to Retrieve the announcement month
	m, err := s.db.RetrieveAnnouncementMonth(context.Background(), month.Date)
	s.NoError(err)
	s.Equal(month.Date, m.Date)
	s.NotEmpty(m.Created)
	s.Equal(m.Modified, m.Created)
	s.Len(m.Announcements, len(month.Announcements))

	// Attempt to Retrieve a non-existent announcement month
	_, err = s.db.RetrieveAnnouncementMonth(context.Background(), "")
	s.ErrorIs(err, storeerrors.ErrEntityNotFound)
	_, err = s.db.RetrieveAnnouncementMonth(context.Background(), "2022-01-01")
	s.Error(err)
	_, err = s.db.RetrieveAnnouncementMonth(context.Background(), "2021-01")
	s.ErrorIs(err, storeerrors.ErrEntityNotFound)

	// Attempt to save an announcement month without a date on it
	month.Date = ""
	err = s.db.UpdateAnnouncementMonth(context.Background(), month)
	s.ErrorIs(err, storeerrors.ErrIncompleteRecord)

	// Sleep to advance the clock for the modified timestamp
	time.Sleep(1 * time.Millisecond)

	// Update the announcement month
	m.Announcements[0].Title = "Happy New Year!"
	err = s.db.UpdateAnnouncementMonth(context.Background(), m)
	s.NoError(err)

	m, err = s.db.RetrieveAnnouncementMonth(context.Background(), m.Date)
	s.NoError(err)
	s.Equal("Happy New Year!", m.Announcements[0].Title)
	s.NotEmpty(m.Modified)
	s.NotEqual(m.Modified, m.Created)

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
	s.NoError(s.db.UpdateAnnouncementMonth(context.Background(), month))

	// Test that we can still retrieve both months
	january, err := s.db.RetrieveAnnouncementMonth(context.Background(), "2022-01")
	s.NoError(err)
	s.Equal("Happy New Year!", january.Announcements[0].Title)

	february, err := s.db.RetrieveAnnouncementMonth(context.Background(), "2022-02")
	s.NoError(err)
	s.Equal("Happy Groundhog Day", february.Announcements[0].Title)

	// Delete an announcement month
	s.NoError(s.db.DeleteAnnouncementMonth(context.Background(), "2022-01"))

	// Should not be able to retrieve the deleted announcement month
	_, err = s.db.RetrieveAnnouncementMonth(context.Background(), "2022-01")
	s.ErrorIs(err, storeerrors.ErrEntityNotFound)
}

func (s *leveldbTestSuite) TestActivityStore() {
	// Load the activity month record from testdata
	data, err := os.ReadFile("../testdata/activity.json")
	s.NoError(err)

	month := &bff.ActivityMonth{}
	err = protojson.Unmarshal(data, month)
	s.NoError(err)

	// Verify the activity month is loaded correctly
	s.NotEmpty(month.Date)
	s.NotEmpty(month.Days)
	s.Empty(month.Created)
	s.Empty(month.Modified)

	// Create the activity month
	s.NoError(s.db.UpdateActivityMonth(context.Background(), month))

	// Attempt to Retrieve the activity month
	m, err := s.db.RetrieveActivityMonth(context.Background(), month.Date)
	s.NoError(err)
	s.Equal(month.Date, m.Date)
	s.NotEmpty(m.Created)
	s.Equal(m.Modified, m.Created)
	s.Len(m.Days, len(month.Days))

	// Attempt to Retrieve a non-existent activity month
	_, err = s.db.RetrieveActivityMonth(context.Background(), "")
	s.ErrorIs(err, storeerrors.ErrEntityNotFound)
	_, err = s.db.RetrieveActivityMonth(context.Background(), "2022-01-01")
	s.Error(err)
	_, err = s.db.RetrieveActivityMonth(context.Background(), "2021-01")
	s.ErrorIs(err, storeerrors.ErrEntityNotFound)

	// Attempt to save an activity month without a date on it
	month.Date = ""
	err = s.db.UpdateActivityMonth(context.Background(), month)
	s.ErrorIs(err, storeerrors.ErrIncompleteRecord)

	// Sleep to advance the clock for the modified timestamp
	time.Sleep(1 * time.Millisecond)

	// Update the activity month
	m.Days[0].Activity.Mainnet["lookup"] = 10
	err = s.db.UpdateActivityMonth(context.Background(), m)
	s.NoError(err)

	m, err = s.db.RetrieveActivityMonth(context.Background(), m.Date)
	s.NoError(err)
	s.Equal(uint64(10), m.Days[0].Activity.Mainnet["lookup"])
	s.NotEmpty(m.Modified)
	s.NotEqual(m.Modified, m.Created)

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
	s.NoError(s.db.UpdateActivityMonth(context.Background(), month))

	// Test that we can still retrieve both months
	january, err := s.db.RetrieveActivityMonth(context.Background(), "2023-01")
	s.NoError(err)
	s.Equal(uint64(1), january.Days[0].Activity.Testnet["lookup"])

	february, err := s.db.RetrieveActivityMonth(context.Background(), "2023-02")
	s.NoError(err)
	s.Equal(uint64(20), february.Days[0].Activity.Testnet["lookup"])

	// Delete an activity month
	s.NoError(s.db.DeleteActivityMonth(context.Background(), "2023-01"))

	// Should not be able to retrieve the deleted activity month
	_, err = s.db.RetrieveActivityMonth(context.Background(), "2023-01")
	s.ErrorIs(err, storeerrors.ErrEntityNotFound)
}

func (s *leveldbTestSuite) TestOrganizationStore() {
	// Create a new organization in the database
	org := &bff.Organization{}
	id, err := s.db.CreateOrganization(context.Background(), org)
	s.NoError(err)

	// Verify that the created record has an ID and timestamps
	s.NotEmpty(org.Id)
	s.Equal(org.Id, id)
	s.NotEmpty(org.Created)
	s.Equal(org.Modified, org.Created)

	// Retrieve the organization by UUID
	uu, err := bff.ParseOrgID(org.Id)
	s.NoError(err)
	o, err := s.db.RetrieveOrganization(context.Background(), uu)
	s.NoError(err)
	s.True(proto.Equal(org, o), "retrieved organization does not match created organization")

	// Attempt to retrieve a non-existent organization
	_, err = s.db.RetrieveOrganization(context.Background(), uuid.Nil)
	s.ErrorIs(err, storeerrors.ErrEntityNotFound)
	_, err = s.db.RetrieveOrganization(context.Background(), uuid.New())
	s.ErrorIs(err, storeerrors.ErrEntityNotFound)

	// Sleep to advance the clock for the modified timestamp
	time.Sleep(1 * time.Millisecond)

	// Update the organization
	org.Name = "Alice Corp"
	err = s.db.UpdateOrganization(context.Background(), org)
	s.NoError(err)

	o, err = s.db.RetrieveOrganization(context.Background(), uu)
	s.NoError(err)
	s.Equal("Alice Corp", o.Name)
	s.NotEmpty(o.Modified)
	s.NotEqual(o.Modified, o.Created)

	// Attempt to update an organization with no Id on it
	org.Id = ""
	s.ErrorIs(s.db.UpdateOrganization(context.Background(), org), storeerrors.ErrIncompleteRecord)

	// Delete the organization
	err = s.db.DeleteOrganization(context.Background(), uu)
	s.NoError(err)
	_, err = s.db.RetrieveOrganization(context.Background(), uu)
	s.ErrorIs(err, storeerrors.ErrEntityNotFound)

	// Create a few organizations
	for i := 0; i < 10; i++ {
		org := &bff.Organization{}
		id, err := s.db.CreateOrganization(context.Background(), org)
		s.NoError(err)
		s.NotEmpty(id)
	}

	// Test Prev() and Next() interactions
	actualOrgs := 0
	iter := s.db.ListOrganizations(context.Background())
	s.False(iter.Prev(), "should move behind the first Organization")
	s.True(iter.Next(), "should move to the first Organization")
	first, err := iter.Organization()
	s.NoError(err)
	s.NotNil(first)
	actualOrgs++
	s.True(iter.Next(), "should move to the second Organization")
	second, err := iter.Organization()
	s.NoError(err)
	s.NotNil(second)
	s.NotEqual(first.Id, second.Id, "organizations should have unique IDs")
	actualOrgs++

	// Consume the rest of the iterator
	for iter.Next() {
		org, err := iter.Organization()
		s.NoError(err)
		s.NotNil(org)
		actualOrgs++
	}
	s.NoError(iter.Error())
	iter.Release()
	s.Equal(actualOrgs, 10, "iterator returned the wrong number of organizations")

	// Create enough organizations to exceed the default page size
	for i := 0; i < 100; i++ {
		org := &bff.Organization{}
		_, err := s.db.CreateOrganization(context.Background(), org)
		s.NoError(err)
	}

	// Consume all of the organizations
	iter = s.db.ListOrganizations(context.Background())
	actualOrgs = 0
	for iter.Next() {
		org, err := iter.Organization()
		s.NoError(err)
		s.NotNil(org)
		actualOrgs++
	}
	s.NoError(iter.Error())
	iter.Release()

	s.Equal(actualOrgs, 110, "iterator returned the wrong number of organizations")
}

func (s *leveldbTestSuite) TestContactStore() {
	// Make sure create errors with a nil contact
	email, err := s.db.CreateContact(context.Background(), nil)
	s.Empty(email)
	s.Equal(err, storeerrors.ErrIncompleteRecord)

	// Make sure create errors with a contact with an empty email
	contact := &models.Contact{
		Email: "",
	}
	email, err = s.db.CreateContact(context.Background(), contact)
	s.Empty(email)
	s.Equal(err, storeerrors.ErrIncompleteRecord)

	// Create a valid contact
	contact = &models.Contact{
		Email:      "testemail",
		Vasps:      []string{"foo", "bar"},
		Verified:   false,
		Token:      "testtoken",
		VerifiedOn: "",
	}
	email, err = s.db.CreateContact(context.Background(), contact)
	s.Equal(email, "testemail")
	s.NoError(err)

	// Make sure retrieve errors with an empty email
	var con *models.Contact
	con, err = s.db.RetrieveContact(context.Background(), "")
	s.Nil(con)
	s.Equal(err, storeerrors.ErrIncompleteRecord)

	// Make sure retrieve throws the proper error when a contact is not found
	con, err = s.db.RetrieveContact(context.Background(), "wrongemail")
	s.Nil(con)
	s.Equal(err, storeerrors.ErrEntityNotFound)

	// Retrieve the created contact
	con, err = s.db.RetrieveContact(context.Background(), "testemail")
	s.Equal(con.Vasps, contact.Vasps)
	s.Equal(con.Verified, contact.Verified)
	s.Equal(con.Token, contact.Token)
	s.NotEmpty(con.Created)
	s.NotEmpty(con.Modified)
	s.NoError(err)

	// Make sure update errors with a nil contact
	err = s.db.UpdateContact(context.Background(), nil)
	s.Equal(err, storeerrors.ErrIncompleteRecord)

	// Make sure update errors with a contact with an empty email
	contact = &models.Contact{
		Email: "",
	}
	err = s.db.UpdateContact(context.Background(), contact)
	s.Equal(err, storeerrors.ErrIncompleteRecord)

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
	err = s.db.UpdateContact(context.Background(), contact)
	s.NoError(err)

	// Retrieve the updated contact
	var updatedCon *models.Contact
	updatedCon, err = s.db.RetrieveContact(context.Background(), "testemail")
	s.NoError(err)
	s.Equal(updatedCon.Vasps, contact.Vasps)
	s.Equal(updatedCon.Verified, contact.Verified)
	s.Equal(updatedCon.Token, contact.Token)
	s.Equal(con.Created, updatedCon.Created)
	s.NotEqual(con.Modified, updatedCon.Modified)

	// Make sure delete errors with an empty email
	err = s.db.DeleteContact(context.Background(), "")
	s.Equal(err, storeerrors.ErrEntityNotFound)

	// Delete the created contact
	err = s.db.DeleteContact(context.Background(), "testemail")
	s.NoError(err)

	// Make sure the contact was deleted
	con, err = s.db.RetrieveContact(context.Background(), "testemail")
	s.Nil(con)
	s.Equal(err, storeerrors.ErrEntityNotFound)
}

func (s *leveldbTestSuite) TestEmailStore() {
	ctx := context.Background()

	s.Run("List", func() {
		cleanup := func() {
			iter := s.db.ListEmails(ctx)
			defer iter.Release()

			for iter.Next() {
				email, _ := iter.Email()
				s.db.DeleteEmail(ctx, email.Email)
			}
		}

		// Ensure the database is empty and emptied at the end of the test.
		cleanup()
		defer cleanup()

		// Test list empty database
		iter := s.db.ListEmails(ctx)
		s.False(iter.Next(), "expected next to be false in empty database")
		s.False(iter.Prev(), "expected prev to be false in empty database")
		iter.Release()
		s.NoError(iter.Error(), "expected no error during iteration")

		// Create some emails in the database
		for i := 1; i < 128; i++ {
			email := &models.Email{Email: fmt.Sprintf("person%03d@example.com", i), Token: "foo"}
			err := s.db.UpdateEmail(ctx, email)
			s.NoError(err, "could not insert email into database")
		}

		count := 0
		iter = s.db.ListEmails(ctx)
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
		key, err := s.db.CreateEmail(ctx, email)
		s.NoError(err, "could not create valid email record")
		s.Equal("gary.vespers@example.com", key, "unexpected key returned")
		s.NotEmpty(email.Created, "expected the created timestamp to be added")
		s.NotEmpty(email.Modified, "expected the modified timestamp to be added")
		s.Equal("gary.vespers@example.com", email.Email, "expected the email to be normalized")
		s.Equal("Gary Vespers", email.Name, "expected the name to be parsed from the email")

		// Test retrieve email record
		record, err := s.db.RetrieveEmail(ctx, "Gary Vespers <Gary.Vespers@example.com>")
		s.NoError(err)
		s.True(proto.Equal(email, record))

		// Sleep to allow time for the modified field to change
		time.Sleep(time.Second)

		// Test updating the email record
		record.Verified = true
		record.VerifiedOn = time.Now().Format(time.RFC3339)
		record.Token = ""

		err = s.db.UpdateEmail(ctx, record)
		s.NoError(err, "could not update email record")
		s.NotEqual(record.Modified, email.Modified, "expected modified timestamp to be updated")

		updated, err := s.db.RetrieveEmail(ctx, "gary.vespers@example.com")
		s.NoError(err, "could not retrieve updated email")
		s.False(proto.Equal(email, updated))
		s.True(proto.Equal(record, updated))

		// Test deleting the email record
		err = s.db.DeleteEmail(ctx, "Gary.Vespers@example.com")
		s.NoError(err, "could not delete email record")

		// Ensure the email was deleted
		_, err = s.db.RetrieveEmail(ctx, "gary.vespers@example.com")
		s.ErrorIs(err, storeerrors.ErrEntityNotFound)
	})

	s.Run("CreateEmptyEmailError", func() {
		// Make sure create errors with a nil email
		email, err := s.db.CreateEmail(ctx, nil)
		s.Empty(email)
		s.ErrorIs(err, storeerrors.ErrIncompleteRecord)

		// Make sure create errors with an email empty string
		email, err = s.db.CreateEmail(ctx, &models.Email{Email: ""})
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

		_, err := s.db.CreateEmail(ctx, email)
		s.NoError(err, "could not create first record")

		_, err = s.db.CreateEmail(ctx, email)
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
			_, err := s.db.CreateEmail(ctx, tc.email)
			s.NoError(err, "test case %d failed", i)
			s.Equal(tc.expectedName, tc.email.Name, "test case %d failed", i)
			s.Equal(tc.expectedEmail, tc.email.Email, "test case %d failed", i)
		}
	})

	s.Run("RetrieveEmptyEmailError", func() {
		email, err := s.db.RetrieveEmail(ctx, "")
		s.Empty(email)
		s.ErrorIs(err, storeerrors.ErrEntityNotFound)
	})

	s.Run("RetrieveNotFound", func() {
		email, err := s.db.RetrieveEmail(ctx, "thisemaildoesnot@exist.com")
		s.Empty(email)
		s.ErrorIs(err, storeerrors.ErrEntityNotFound)
	})

	s.Run("UpdateEmptyEmailError", func() {
		// Make sure update errors with a nil email
		err := s.db.UpdateEmail(ctx, nil)
		s.ErrorIs(err, storeerrors.ErrIncompleteRecord)

		// Make sure update errors with an email empty string
		err = s.db.UpdateEmail(ctx, &models.Email{Email: ""})
		s.ErrorIs(err, storeerrors.ErrIncompleteRecord)
	})

	s.Run("DeleteEmptyEmailError", func() {
		// Make sure delete errors with an email empty string
		err := s.db.DeleteEmail(ctx, "")
		s.ErrorIs(err, storeerrors.ErrEntityNotFound)
	})
}

func (s *leveldbTestSuite) TestDirectoryContactStore() {
	ctx := context.Background()
	require := s.Require()

	// Load fixtures into database for tests
	data, err := os.ReadFile("../testdata/altvasp.json")
	require.NoError(err)

	vasp := &pb.VASP{}
	err = protojson.Unmarshal(data, vasp)
	require.NoError(err)

	vaspID, err := s.db.CreateVASP(ctx, vasp)
	require.NoError(err, "could not create vasp fixture")

	admin := &models.Email{Name: vasp.Contacts.Administrative.Name, Email: vasp.Contacts.Administrative.Email, Verified: true, VerifiedOn: time.Now().Format(time.RFC3339)}
	tech := &models.Email{Name: vasp.Contacts.Technical.Name, Email: vasp.Contacts.Technical.Email, Verified: true, VerifiedOn: time.Now().Format(time.RFC3339)}
	legal := &models.Email{Name: vasp.Contacts.Legal.Name, Email: vasp.Contacts.Legal.Email, Verified: true, VerifiedOn: time.Now().Format(time.RFC3339)}
	billing := &models.Email{Name: vasp.Contacts.Billing.Name, Email: vasp.Contacts.Billing.Email, Verified: true, VerifiedOn: time.Now().Format(time.RFC3339)}

	for _, record := range []*models.Email{admin, tech, legal, billing} {
		err = s.db.UpdateEmail(ctx, record)
		require.NoError(err, "could not create associated contact email")
	}

	defer func() {
		require.NoError(s.db.DeleteVASP(ctx, vaspID), "could not delete darlene vasp")
		for _, record := range []*models.Email{admin, tech, legal, billing} {
			err = s.db.DeleteEmail(ctx, record.Email)
			require.NoError(err, "could not delete associated contact email")
		}
	}()

	s.Run("VASPContacts", func() {
		contacts, err := s.db.VASPContacts(ctx, vasp)
		require.NoError(err, "could not get vasp contacts for darlene")
		require.Equal(vasp.Id, contacts.VASP, "expected the vasp ID to be on the contacts")
		require.Equal(vasp.Contacts, contacts.Contacts, "expected the contacts to match the VASP")
		require.Len(contacts.Emails, 4, "expected two emails retrieved")
	})

	s.Run("RetrieveVASPContacts", func() {
		contacts, err := s.db.RetrieveVASPContacts(ctx, vaspID)
		require.NoError(err, "could not retrieve vasp contacts for darlene")
		require.Equal(vasp.Id, contacts.VASP, "expected the vasp ID to be on the contacts")
		require.True(proto.Equal(vasp.Contacts, contacts.Contacts), "expected the contacts to match the VASP")
		require.Len(contacts.Emails, 4, "expected two emails retrieved")
	})
}
