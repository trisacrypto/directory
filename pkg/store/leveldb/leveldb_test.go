package leveldb

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
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
	// NOTE: ConsoleLog MUST be false otherwise this will be overriden
	logger.Discard()

	path, err := ioutil.TempDir("", "gdsldbstore-*")
	s.NoError(err)

	// Open the database in a temp directory
	s.path = path
	s.db, err = Open(path)
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
	data, err := ioutil.ReadFile("../testdata/vasp.json")
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
	id, err := s.db.CreateVASP(alice)
	s.NoError(err)
	s.NotEmpty(id)

	// Attempt to Retrieve the VASP
	alicer, err := s.db.RetrieveVASP(id)
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
	err = s.db.UpdateVASP(alicer)
	s.NoError(err)

	alicer, err = s.db.RetrieveVASP(id)
	s.NoError(err)
	s.Equal(id, alicer.Id)
	s.NotEmpty(alicer.LastUpdated)
	s.NotEqual(alicer.FirstListed, alicer.LastUpdated)
	s.NotEmpty(alicer.Version)
	s.Equal(uint64(2), alicer.Version.Version)
	s.Equal(alicer.VerificationStatus, pb.VerificationState_VERIFIED)

	// Delete the VASP
	err = s.db.DeleteVASP(id)
	s.NoError(err)
	alicer, err = s.db.RetrieveVASP(id)
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
		_, err := s.db.CreateVASP(vasp)
		s.NoError(err)
	}

	// Test listing all of the VASPs
	reqs, err := s.db.ListVASPs().All()
	s.NoError(err)
	s.Len(reqs, 10)

	// Test iterating over all the VASPs
	var niters int
	iter := s.db.ListVASPs()
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
	s.NotEmpty(cert.Details.NotBefore)
	s.NotEmpty(cert.Details.NotAfter)

	// Attempt to Create the Cert
	id, err := s.db.CreateCert(cert)
	s.NoError(err)

	// Attempt to Retrieve the Cert
	crr, err := s.db.RetrieveCert(id)
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
	_, err = s.db.CreateCert(icrr)
	s.ErrorIs(err, storeerrors.ErrIDAlreadySet)

	// Update the Cert
	crr.Status = models.CertificateState_REVOKED
	err = s.db.UpdateCert(crr)
	s.NoError(err)

	crr, err = s.db.RetrieveCert(id)
	s.NoError(err)
	s.Equal(id, crr.Id)
	s.Equal(models.CertificateState_REVOKED, crr.Status)

	// Attempt to update a certificate with no Id on it
	cert.Id = ""
	s.ErrorIs(s.db.UpdateCert(cert), storeerrors.ErrIncompleteRecord)

	// Delete the Cert
	err = s.db.DeleteCert(id)
	s.NoError(err)
	crr, err = s.db.RetrieveCert(id)
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
		_, err := s.db.CreateCert(crr)
		s.NoError(err)
	}

	// Test listing all of the certificates
	certs, err := s.db.ListCerts().All()
	s.NoError(err)
	s.Len(certs, 10)

	// Test iterating over all the certificates
	var niters int
	iter := s.db.ListCerts()
	for iter.Next() {
		s.NotEmpty(iter.Cert())
		niters++
	}
	s.NoError(iter.Error())
	iter.Release()
	s.Equal(10, niters)
}

func (s *leveldbTestSuite) TestCertificateRequestStore() {
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

	// Attempt to Create the CertReq
	id, err := s.db.CreateCertReq(certreq)
	s.NoError(err)

	// Attempt to Retrieve the CertReq
	crr, err := s.db.RetrieveCertReq(id)
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
	_, err = s.db.CreateCertReq(icrr)
	s.ErrorIs(err, storeerrors.ErrIDAlreadySet)

	// Sleep for a second to roll over the clock for the modified time stamp
	time.Sleep(1 * time.Second)

	// Update the CertReq
	crr.Status = models.CertificateRequestState_COMPLETED
	err = s.db.UpdateCertReq(crr)
	s.NoError(err)

	crr, err = s.db.RetrieveCertReq(id)
	s.NoError(err)
	s.Equal(id, crr.Id)
	s.Equal(models.CertificateRequestState_COMPLETED, crr.Status)
	s.NotEmpty(crr.Modified)
	s.NotEqual(crr.Modified, crr.Created)

	// Attempt to update a certificate request with no Id on it
	certreq.Id = ""
	s.ErrorIs(s.db.UpdateCertReq(certreq), storeerrors.ErrIncompleteRecord)

	// Delete the CertReq
	err = s.db.DeleteCertReq(id)
	s.NoError(err)
	crr, err = s.db.RetrieveCertReq(id)
	s.ErrorIs(err, storeerrors.ErrEntityNotFound)
	s.Empty(crr)

	// Add a few more certificate requests
	for i := 0; i < 10; i++ {
		crr := &models.CertificateRequest{
			Vasp:       uuid.New().String(),
			CommonName: fmt.Sprintf("trisa%d.example.com", i+1),
			Status:     models.CertificateRequestState_COMPLETED,
		}
		_, err := s.db.CreateCertReq(crr)
		s.NoError(err)
	}

	// Test listing all of the certificates
	reqs, err := s.db.ListCertReqs().All()
	s.NoError(err)
	s.Len(reqs, 10)

	// Test iterating over all the certificates
	var niters int
	iter := s.db.ListCertReqs()
	for iter.Next() {
		s.NotEmpty(iter.CertReq())
		niters++
	}
	s.NoError(iter.Error())
	iter.Release()
	s.Equal(10, niters)
}
