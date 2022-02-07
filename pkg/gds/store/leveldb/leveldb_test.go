package leveldb

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	storeerrors "github.com/trisacrypto/directory/pkg/gds/store/errors"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/encoding/protojson"
)

type leveldbTestSuite struct {
	suite.Suite
	path string
	db   *Store
}

func (s *leveldbTestSuite) SetupSuite() {
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
