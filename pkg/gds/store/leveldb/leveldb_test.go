package leveldb

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/gds/global/v1"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
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
	// Use global versioner for these tests
	vm := &global.VersionManager{PID: 8, Owner: "8:mitchell", Region: "us-east-2c"}
	err := s.db.WithVersionManager(vm)
	s.NoError(err)

	// Load the VASP record from testdata
	data, err := ioutil.ReadFile("testdata/alice.json")
	s.NoError(err)

	alice := &pb.VASP{}
	err = protojson.Unmarshal(data, alice)
	s.NoError(err)

	// Attempt to Create the VASP
	id, err := s.db.CreateVASP(alice)
	s.NoError(err)

	// Attempt to Retrieve the VASP
	alicer, err := s.db.RetrieveVASP(id)
	s.NoError(err)
	s.Equal(id, alicer.Id)

	// Test the version of the retrieved object
	meta, deletedOn, err := models.GetMetadata(alicer)
	s.NoError(err)
	s.True(deletedOn.IsZero())
	s.True(proto.Equal(meta.Version, &global.Version{Pid: 8, Version: 1, Region: "us-east-2c"}))
	s.Empty(meta.Version.Parent)
	s.Equal(vm.Owner, meta.Owner)
	s.Equal(vm.Region, meta.Region)

	// Update the VASP
	alicer.VerificationStatus = pb.VerificationState_VERIFIED
	alicer.VerifiedOn = "2021-06-30T10:40:40Z"
	err = s.db.UpdateVASP(alicer)
	s.NoError(err)

	alicer, err = s.db.RetrieveVASP(id)
	s.NoError(err)
	s.Equal(id, alicer.Id)

	// Test the version of the retrieved object
	meta, deletedOn, err = models.GetMetadata(alicer)
	s.NoError(err)
	s.True(deletedOn.IsZero())
	s.True(proto.Equal(meta.Version, &global.Version{Pid: 8, Version: 2, Region: "us-east-2c", Parent: &global.Version{Pid: 8, Version: 1, Region: "us-east-2c"}}))

	// Delete the VASP
	err = s.db.DeleteVASP(id)
	s.NoError(err)
	alicer, err = s.db.RetrieveVASP(id)
	s.ErrorIs(err, ErrEntityNotFound)
	s.Empty(alicer)
}

func (s *leveldbTestSuite) TestCertificateStore() {
	// Use global versioner for these tests
	vm := &global.VersionManager{PID: 8, Owner: "8:mitchell", Region: "us-east-2c"}
	err := s.db.WithVersionManager(vm)
	s.NoError(err)

	// Load the VASP record from testdata
	data, err := ioutil.ReadFile("testdata/certreq.json")
	s.NoError(err)

	certreq := &models.CertificateRequest{}
	err = protojson.Unmarshal(data, certreq)
	s.NoError(err)

	// Attempt to Create the CertReq
	id, err := s.db.CreateCertReq(certreq)
	s.NoError(err)

	// Attempt to Retrieve the CertReq
	crr, err := s.db.RetrieveCertReq(id)
	s.NoError(err)
	s.Equal(id, crr.Id)

	// Test the version of the retrieved object
	s.NoError(err)
	s.Empty(crr.Deleted)
	s.True(proto.Equal(crr.Metadata.Version, &global.Version{Pid: 8, Version: 1, Region: "us-east-2c"}))
	s.Empty(crr.Metadata.Version.Parent)
	s.Equal(vm.Owner, crr.Metadata.Owner)
	s.Equal(vm.Region, crr.Metadata.Region)

	// Update the CertReq
	crr.Status = models.CertificateRequestState_COMPLETED
	err = s.db.UpdateCertReq(crr)
	s.NoError(err)

	crr, err = s.db.RetrieveCertReq(id)
	s.NoError(err)
	s.Equal(id, crr.Id)

	// Test the version of the retrieved object
	s.NoError(err)
	s.Empty(crr.Deleted)
	s.True(proto.Equal(crr.Metadata.Version, &global.Version{Pid: 8, Version: 2, Region: "us-east-2c", Parent: &global.Version{Pid: 8, Version: 1, Region: "us-east-2c"}}))

	// Delete the CertReq
	err = s.db.DeleteCertReq(id)
	s.NoError(err)
	crr, err = s.db.RetrieveCertReq(id)
	s.ErrorIs(err, ErrEntityNotFound)
	s.Empty(crr)
}
