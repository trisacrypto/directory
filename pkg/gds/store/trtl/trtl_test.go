package trtl_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	store "github.com/trisacrypto/directory/pkg/gds/store/trtl"
	"github.com/trisacrypto/directory/pkg/trtl"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/utils/bufconn"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/encoding/protojson"

	storeerrors "github.com/trisacrypto/directory/pkg/gds/store/errors"
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

	var err error
	s.tmpdb, err = ioutil.TempDir("../testdata", "db-*")
	require.NoError(err)

	conf := config.Config{
		Maintenance: false,
		BindAddr:    ":4436",
		LogLevel:    logger.LevelDecoder(zerolog.DebugLevel),
		ConsoleLog:  true,
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

	trtl, err := trtl.New(*s.conf)
	require.NoError(err)

	s.grpc = bufconn.New(bufSize)
	go trtl.Run(s.grpc.Listener)
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
	require.NoError(s.grpc.Connect())
	defer s.grpc.Close()

	db, err := store.NewMock(s.grpc.Conn)
	require.NoError(err)

	// Initially there should be no VASPs
	// TODO: Test snapshot isolation iteration
	iter := db.ListVASPs()
	require.False(iter.Next())

	// Should get a not found error trying to retrieve a VASP that doesn't exist
	_, err = db.RetrieveVASP("12345")
	require.EqualError(err, storeerrors.ErrEntityNotFound.Error())

	// Attempt to Create the VASP
	id, err := db.CreateVASP(alice)
	require.NoError(err)
	require.NotEmpty(id)

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
		_, err := db.CreateVASP(vasp)
		require.NoError(err)
	}

	// Test listing all of the VASPs
	reqs, err := db.ListVASPs().All()
	require.NoError(err)
	require.Len(reqs, 10)

	// Test seeking to a specific VASP
	key := reqs[5].Id
	iter = db.ListVASPs()
	require.True(iter.SeekId(key))
	require.True(iter.Next())
	v, err := iter.VASP()
	require.NoError(err)
	require.NoError(iter.Error())
	require.Equal(key, v.Id)

	// Test that Prev() and Next() work properly
	require.False(iter.Prev())
	require.True(iter.Next())
	next, err := iter.VASP()
	require.NoError(err)
	require.NotNil(next)
	require.NotEqual(key, next.Id)
	require.True(iter.Prev())
	prev, err := iter.VASP()
	require.NoError(err)
	require.NotNil(prev)
	require.Equal(key, prev.Id)
	require.True(iter.Next())
	next, err = iter.VASP()
	require.NoError(err)
	require.NotNil(next)
	require.NotEqual(key, next.Id)

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
	for i := 0; i < 100; i++ {
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
		_, err := db.CreateVASP(vasp)
		require.NoError(err)
	}

	// Test listing all of the VASPs
	reqs, err = db.ListVASPs().All()
	require.NoError(err)
	require.Len(reqs, 110)
}
