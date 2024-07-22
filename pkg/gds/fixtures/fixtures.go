package fixtures

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/store"
	"github.com/trisacrypto/directory/pkg/store/config"
	trtlstore "github.com/trisacrypto/directory/pkg/store/trtl"
	"github.com/trisacrypto/directory/pkg/trtl"
	trtlmock "github.com/trisacrypto/directory/pkg/trtl/mock"
	"github.com/trisacrypto/directory/pkg/utils"
	"github.com/trisacrypto/directory/pkg/utils/bufconn"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
)

type StoreType uint8

const (
	StoreUnknown StoreType = iota
	StoreLevelDB
	StoreTrtl
)

type FixtureType uint8

const (
	Unknown FixtureType = iota
	Empty
	Small
	Full
)

var (
	smallDBSubset = map[string]map[string]struct{}{
		wire.NamespaceVASPs: {
			"charliebank": {},
			"delta":       {},
			"hotel":       {},
		},
		wire.NamespaceContacts: {
			"adam@example.com":  {},
			"bruce@example.com": {},
		},
	}
)

func New(fixturesPath, dbPath string, stype StoreType) (lib *Library, err error) {
	lib = &Library{
		fixturesPath: fixturesPath,
		dbPath:       dbPath,
		stype:        stype,
	}

	if err = lib.LoadReferenceFixtures(); err != nil {
		return nil, err
	}

	if lib.stype == StoreTrtl {
		if err = lib.SetupTrtl(); err != nil {
			return nil, err
		}
	}

	return lib, nil
}

type Library struct {
	ftype        FixtureType
	fixtures     map[string]map[string]interface{}
	fixturesPath string
	dbPath       string
	stype        StoreType
	trtl         *trtl.Server
	trtlListener *bufconn.GRPCListener
}

func (lib *Library) Fixtures() map[string]map[string]interface{} {
	return lib.fixtures
}

func (lib *Library) GetContact(name string) (contact *models.Contact, err error) {
	var ok bool
	if contact, ok = lib.fixtures[wire.NamespaceContacts][name].(*models.Contact); !ok {
		return nil, fmt.Errorf("could not retrieve contact %s from fixtures", name)
	}

	return contact, nil
}

func (lib *Library) GetVASP(name string) (vasp *pb.VASP, err error) {
	var ok bool
	if vasp, ok = lib.fixtures[wire.NamespaceVASPs][name].(*pb.VASP); !ok {
		return nil, fmt.Errorf("could not retrieve VASP %s from fixtures", name)
	}

	return vasp, nil
}

func (lib *Library) GetCert(name string) (cert *models.Certificate, err error) {
	var ok bool
	if cert, ok = lib.fixtures[wire.NamespaceCerts][name].(*models.Certificate); !ok {
		return nil, fmt.Errorf("could not retrieve certificate %s from fixtures", name)
	}

	return cert, nil
}

func (lib *Library) GetCertReq(name string) (certReq *models.CertificateRequest, err error) {
	var ok bool
	if certReq, ok = lib.fixtures[wire.NamespaceCertReqs][name].(*models.CertificateRequest); !ok {
		return nil, fmt.Errorf("could not retrieve certificate request %s from fixtures", name)
	}

	return certReq, nil
}

func (lib *Library) FixtureType() FixtureType {
	return lib.ftype
}

func (lib *Library) StoreType() StoreType {
	return lib.stype
}

func (lib *Library) DBPath() string {
	return lib.dbPath
}

func (lib *Library) Load(t FixtureType) (err error) {
	// Check if the requested fixture type is already loaded
	if lib.ftype == t {
		return nil
	}

	// The trtl server must be shutdown before reloading the fixtures
	if err = lib.ShutdownTrtl(); err != nil {
		return err
	}

	// Recreate the database directory
	if err = lib.ResetDB(); err != nil {
		return err
	}

	// Create a new trtl server, this must be done before generating the database
	if lib.stype == StoreTrtl {
		if err = lib.SetupTrtl(); err != nil {
			return err
		}
	}

	// Generate the database with the new fixtures
	if err = lib.GenerateDB(t); err != nil {
		return err
	}

	lib.ftype = t
	return nil
}

func (lib *Library) ShutdownTrtl() (err error) {
	if lib.trtl != nil {
		if err = lib.trtl.Shutdown(); err != nil {
			return err
		}
	}
	if lib.trtlListener != nil {
		lib.trtlListener.Release()
	}
	return nil
}

func (lib *Library) Close() (err error) {
	if err = lib.ShutdownTrtl(); err != nil {
		return err
	}

	if err = os.RemoveAll(lib.dbPath); err != nil {
		return err
	}
	return nil
}

func (lib *Library) Reset() {
	lib.ftype = Unknown
}

func (lib *Library) ResetDB() (err error) {
	if err = os.RemoveAll(lib.dbPath); err != nil {
		return err
	}

	var db *leveldb.DB
	if db, err = leveldb.OpenFile(lib.dbPath, nil); err != nil {
		return err
	}
	return db.Close()
}

// SetupTrtl starts the Trtl server, which must be done before calling GenerateDB.
func (lib *Library) SetupTrtl() (err error) {
	conf := trtlmock.Config()
	conf.Database.URL = "leveldb:///" + lib.dbPath

	// Mark as processed since the config wasn't loaded from the envrionment
	if conf, err = conf.Mark(); err != nil {
		return err
	}

	// Start the Trtl server
	if lib.trtl, err = trtl.New(conf); err != nil {
		return err
	}

	// Using a bufconn listener allows us to avoid network requests
	lib.trtlListener = bufconn.New("")
	go lib.trtl.Run(lib.trtlListener.Listener)

	// Connect to the running Trtl server
	if err = lib.trtlListener.Connect(context.Background()); err != nil {
		return err
	}
	return nil
}

func (lib *Library) ConnectTrtl(ctx context.Context) (conn *grpc.ClientConn, err error) {
	if lib.trtlListener == nil {
		return nil, nil
	}

	if err = lib.trtlListener.Connect(ctx); err != nil {
		return nil, err
	}
	return lib.trtlListener.Conn, nil
}

// Load reference fixtures into the fixtures map; the reference fixtures contain all
// possible fixtures that can be in any fixture database. The fixtures themselves are
// a mapping of namespace to a map of ID strings to the unmarshaled value. E.g. to get
// the expected value for VASP 1234 you would use s.fixtures[vasps]["1234"]
func (lib *Library) LoadReferenceFixtures() (err error) {
	// Extract the gzip archive
	var root string
	if root, err = utils.ExtractGzip(lib.fixturesPath, "testdata", true); err != nil {
		return err
	}
	defer os.RemoveAll(root)

	// Create the reference fixtures map
	lib.fixtures = make(map[string]map[string]interface{})
	for _, namespace := range []string{wire.NamespaceContacts, wire.NamespaceVASPs, wire.NamespaceCerts, wire.NamespaceCertReqs} {
		lib.fixtures[namespace] = make(map[string]interface{})
	}

	// Load the JSON fixtures into the fixtures map
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".json" {
			return nil
		}

		prefix := filepath.Base(filepath.Dir(path))     // prefix is the directory in the fixture
		key := strings.TrimSuffix(info.Name(), ".json") // key is the name of the file in the fixture

		// Unmarshal the JSON into the global fixtures map.
		var data []byte
		if data, err = os.ReadFile(path); err != nil {
			return err
		}

		switch prefix {
		case wire.NamespaceContacts:
			contact := &models.Contact{}
			if err = protojson.Unmarshal(data, contact); err != nil {
				return err
			}
			lib.fixtures[wire.NamespaceContacts][key] = contact
		case wire.NamespaceVASPs:
			vasp := &pb.VASP{}
			if err = protojson.Unmarshal(data, vasp); err != nil {
				return err
			}
			lib.fixtures[wire.NamespaceVASPs][key] = vasp
		case wire.NamespaceCerts:
			cert := &models.Certificate{}
			if err = protojson.Unmarshal(data, cert); err != nil {
				return err
			}
			lib.fixtures[wire.NamespaceCerts][key] = cert
		case wire.NamespaceCertReqs:
			cert := &models.CertificateRequest{}
			if err = protojson.Unmarshal(data, cert); err != nil {
				return err
			}
			lib.fixtures[wire.NamespaceCertReqs][key] = cert
		default:
			return fmt.Errorf("unrecognized prefix for file: %s", info.Name())
		}
		return nil
	})
}

// generateDB generates an updated database and compresses it to a gzip file.
// Note: This also generates a temporary directory which the suite teardown
// should clean up.
func (lib *Library) GenerateDB(ftype FixtureType) (err error) {
	// Data is required to generate the database
	if len(lib.fixtures) == 0 {
		return errors.New("there are no reference fixtures in the library")
	}

	if v, ok := lib.fixtures[wire.NamespaceVASPs]; !ok || len(v) == 0 {
		return errors.New("there are no VASPs in the fixtures library")
	}

	if c, ok := lib.fixtures[wire.NamespaceCerts]; !ok || len(c) == 0 {
		return errors.New("there are no certificates in the fixtures library")
	}

	if cr, ok := lib.fixtures[wire.NamespaceCertReqs]; !ok || len(cr) == 0 {
		return errors.New("there are no certificate requests in the fixtures library")
	}

	// No need to do anything with the empty database
	if ftype == Empty {
		return
	}

	// Create the Store object depending on the store type
	var db store.Store
	switch lib.stype {
	case StoreLevelDB:
		if db, err = store.Open(config.StoreConfig{
			URL:           "leveldb:///" + lib.dbPath,
			ReindexOnBoot: false,
		}); err != nil {
			return err
		}
		defer db.Close()
	case StoreTrtl:
		if err = lib.trtlListener.Connect(context.Background()); err != nil {
			return err
		}
		if db, err = trtlstore.NewMock(lib.trtlListener.Conn); err != nil {
			return err
		}
		defer db.Close()
	default:
		return errors.New("unknown store type")
	}

	// Loop through all the reference fixtures, adding them as necessary
	for namespace, items := range lib.fixtures {
	itemLoop:
		for key, item := range items {

			// If we're in small database mode, check if we should add the data
			if ftype == Small {
				if _, ok := smallDBSubset[namespace][key]; !ok {
					continue itemLoop
				}
			}

			// Add the fixture to the database, updating indices
			switch namespace {
			case wire.NamespaceContacts:
				if err = db.UpdateContact(context.Background(), item.(*models.Contact)); err != nil {
					return err
				}
			case wire.NamespaceVASPs:
				vasp := item.(*pb.VASP)
				var id string
				if id, err = db.CreateVASP(context.Background(), vasp); err != nil {
					return err
				}

				if vasp.Id != id {
					return fmt.Errorf("VASP ID mismatch after creation: %s != %s", vasp.Id, id)
				}
			case wire.NamespaceCerts:
				if err = db.UpdateCert(context.Background(), item.(*models.Certificate)); err != nil {
					return err
				}
			case wire.NamespaceCertReqs:
				if err = db.UpdateCertReq(context.Background(), item.(*models.CertificateRequest)); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unrecognized namespace: %s", namespace)
			}
		}
	}

	log.Info().Msg("successfully regenerated test database")
	return nil
}

// CompareFixture returns True if the given object matches the object in the fixtures
// library. This is used in tests to verify the correctness of the reference library,
// and to verify endpoints that return unmodified objects from the database.
func (lib *Library) CompareFixture(namespace, key string, obj interface{}, removeExtra, removeSerials bool) (matches bool, err error) {
	var (
		ok bool
	)

	if _, ok = lib.fixtures[namespace]; !ok {
		return false, fmt.Errorf("unknown namespace %s", namespace)
	}

	// Reset any time fields for the comparison and compare directly
	switch namespace {
	case wire.NamespaceContacts:
		var a *models.Contact
		for _, f := range lib.fixtures[namespace] {
			ref := f.(*models.Contact)
			if ref.Email == key {
				a = ref
				break
			}
		}
		if a == nil {
			return false, fmt.Errorf("unknown contact fixture %s", key)
		}

		var b *models.Contact
		if b, ok = obj.(*models.Contact); !ok {
			return false, errors.New("obj is not a Contact object")
		}

		// Remove time fields for comparison
		a.Created, b.Created = "", ""
		a.Modified, b.Modified = "", ""
		a.VerifiedOn, b.VerifiedOn = "", ""

		return proto.Equal(a, b), nil

	case wire.NamespaceVASPs:
		var a *pb.VASP
		for _, f := range lib.fixtures[namespace] {
			ref := f.(*pb.VASP)
			if ref.Id == key {
				a = ref
				break
			}
		}
		if a == nil {
			return false, fmt.Errorf("unknown VASP fixture %s", key)
		}

		var b *pb.VASP
		if b, ok = obj.(*pb.VASP); !ok {
			return false, errors.New("obj is not a VASP object")
		}

		// Remove time fields for comparison
		a.LastUpdated, b.LastUpdated = "", ""

		if removeExtra {
			a.Extra, b.Extra = nil, nil
			iter := models.NewContactIterator(a.Contacts)
			for iter.Next() {
				contact, _ := iter.Value()
				contact.Extra = nil
			}

			iter = models.NewContactIterator(b.Contacts)
			for iter.Next() {
				contact, _ := iter.Value()
				contact.Extra = nil
			}
		}

		if removeSerials {
			a.IdentityCertificate.SerialNumber, b.IdentityCertificate.SerialNumber = nil, nil

			for _, cert := range a.SigningCertificates {
				cert.SerialNumber = nil
			}

			for _, cert := range b.SigningCertificates {
				cert.SerialNumber = nil
			}
		}

		return proto.Equal(a, b), nil

	case wire.NamespaceCerts:
		var a *models.Certificate
		for _, f := range lib.fixtures[namespace] {
			ref := f.(*models.Certificate)
			if ref.Id == key {
				a = ref
				break
			}
		}
		if a == nil {
			return false, fmt.Errorf("unknown cert fixture %s", key)
		}

		var b *models.Certificate
		if b, ok = obj.(*models.Certificate); !ok {
			return false, errors.New("obj is not a Certificate object")
		}
		return proto.Equal(a, b), nil

	case wire.NamespaceCertReqs:
		var a *models.CertificateRequest
		for _, f := range lib.fixtures[namespace] {
			ref := f.(*models.CertificateRequest)
			if ref.Id == key {
				// Avoid modifying the object in the fixtures map
				a = ref
				break
			}
		}
		if a == nil {
			return false, fmt.Errorf("unknown certreq fixture %s", key)
		}

		var b *models.CertificateRequest
		if b, ok = obj.(*models.CertificateRequest); !ok {
			return false, errors.New("obj is not a CertificateRequest object")
		}

		a.Modified, b.Modified = "", ""
		a.Created, b.Created = "", ""

		return proto.Equal(a, b), nil

	default:
		return false, fmt.Errorf("unrecognized namespace: %s", namespace)
	}
}

func RemarshalProto(namespace string, obj map[string]interface{}) (_ protoreflect.ProtoMessage, err error) {
	var data []byte
	if data, err = json.Marshal(obj); err != nil {
		return nil, err
	}

	jsonpb := &protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	switch namespace {
	case wire.NamespaceContacts:
		contact := &models.Contact{}
		if err = jsonpb.Unmarshal(data, contact); err != nil {
			return nil, err
		}
		return contact, nil
	case wire.NamespaceVASPs:
		vasp := &pb.VASP{}
		if err = jsonpb.Unmarshal(data, vasp); err != nil {
			return nil, err
		}
		return vasp, nil
	case wire.NamespaceCerts:
		cert := &models.Certificate{}
		if err = jsonpb.Unmarshal(data, cert); err != nil {
			return nil, err
		}
		return cert, nil
	case wire.NamespaceCertReqs:
		certreq := &models.CertificateRequest{}
		if err = jsonpb.Unmarshal(data, certreq); err != nil {
			return nil, err
		}
		return certreq, nil
	default:
		return nil, fmt.Errorf("unknown namespace %q", namespace)
	}
}

// ClearContactEmailLogs clears the contact email logs on a VASP object. Tests which
// assert against state on the contact email logs should call this method to ensure
// that the logs are empty before reaching the test point.
func ClearContactEmailLogs(vasp *pb.VASP) (err error) {
	contacts := vasp.Contacts
	iter := models.NewContactIterator(contacts)
	for iter.Next() {
		contact, _ := iter.Value()
		extra := &models.GDSContactExtraData{}
		if contact.Extra != nil {
			if err = contact.Extra.UnmarshalTo(extra); err != nil {
				return err
			}
			extra.EmailLog = nil
			if contact.Extra, err = anypb.New(extra); err != nil {
				return err
			}
		}
	}
	return nil
}
