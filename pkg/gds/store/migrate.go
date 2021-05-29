package store

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	v1beta1gds "github.com/trisacrypto/directory/pkg/gds/models/v1"
	v1alpha1 "github.com/trisacrypto/directory/pkg/trisads/pb/models/v1alpha1"
	v1beta1 "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	vaspPrefix  = "vasps::"
	careqPrefix = "certreqs::"
	indexPrefix = "index::"
)

// Migrate is a temporary command to move v1alpha1 data to v1beta1 data, potentially
// through an intermediate JSON step so that the data can be edited while being migrated.
func Migrate(src, srcFormat, dst, dstFormat string) (err error) {

	var (
		loader migrateLoader
		dumper migrateDumper
	)

	switch srcFormat {
	case "leveldb":
		if loader, err = newldbMigrateLoader(src); err != nil {
			return err
		}
	case "json":
		loader = &jsonMigrator{path: src}
	default:
		return fmt.Errorf("unknown src format %q", srcFormat)
	}

	switch dstFormat {
	case "leveldb":
		if dumper, err = newldbMigrateDumper(dst); err != nil {
			return err
		}
	case "json":
		dumper = &jsonMigrator{path: dst}
	default:
		return fmt.Errorf("unknown dst format %q", dstFormat)
	}

	return loader.Load(dumper)
}

type migrateLoader interface {
	Load(dumper migrateDumper) error
}

type migrateDumper interface {
	// Dump takes a key and a protocol buffer message to save to dis
	Dump(string, proto.Message) error
}

// Loads v1alpha1 models from leveldb and converts them to v1beta1
type leveldbLoader struct {
	db *leveldb.DB
}

// Dumps v1beta1 models into leveldb
type leveldbDumper struct {
	db *leveldb.DB
}

// Loads and dumps v1beta1 models from and to protojson in a directory on disk
type jsonMigrator struct {
	path string
}

func newldbMigrateLoader(src string) (m *leveldbLoader, err error) {
	m = &leveldbLoader{}
	if m.db, err = leveldb.OpenFile(src, &opt.Options{ErrorIfMissing: true, ReadOnly: true}); err != nil {
		return nil, err
	}
	return m, nil
}

func newldbMigrateDumper(dst string) (m *leveldbDumper, err error) {
	m = &leveldbDumper{}
	if m.db, err = leveldb.OpenFile(dst, &opt.Options{ErrorIfExist: true}); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *leveldbLoader) Load(d migrateDumper) (err error) {
	// Close the database after we're done loading
	defer m.db.Close()

	// Iterate through all the key/value pairs in the database
	iter := m.db.NewIterator(nil, nil)
	defer iter.Release()

dbscan:
	for iter.Next() {
		// Get the Key/Value Pair
		key := string(iter.Key())
		data := iter.Value()

		// Load the v1alpha1 protocol buffer message based on type
		// Convert to a v1beta1 protocol buffer message also based on type
		var msg proto.Message
		if strings.HasPrefix(key, vaspPrefix) {
			vasp := &v1alpha1.VASP{}
			if err = proto.Unmarshal(data, vasp); err != nil {
				return fmt.Errorf("could not unmarshal %q into v1alpha1.VASP record", key)
			}

			if msg, err = convertVASP(vasp); err != nil {
				return err
			}
		} else if strings.HasPrefix(key, careqPrefix) {
			careq := &v1alpha1.CertificateRequest{}
			if err = proto.Unmarshal(data, careq); err != nil {
				return fmt.Errorf("could not unmarshal %q into v1alpha1.CertificateRequest record", key)
			}

			if msg, err = convertCertReq(careq); err != nil {
				return err
			}
		} else if strings.HasPrefix(key, indexPrefix) {
			// Skipping index message
			fmt.Printf("skipping key %q\n", key)
			continue dbscan
		} else {
			return fmt.Errorf("key with unknown prefix, could not parse: %q", key)
		}

		// Dump the protocol buffer message
		if err = d.Dump(key, msg); err != nil {
			return err
		}
	}

	return iter.Error()
}

func (m *jsonMigrator) Dump(key string, msg proto.Message) error {
	opts := protojson.MarshalOptions{
		Multiline:       true,
		Indent:          "  ",
		AllowPartial:    true,
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
	}

	data, err := opts.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not marshal %q: %s", key, err)
	}

	path := filepath.Join(m.path, key+".json")
	if err = ioutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("could not write to %s: %s", path, err)
	}

	fmt.Printf("wrote %q to %s\n", key, path)
	return nil
}

func (m *jsonMigrator) Load(d migrateDumper) error {
	paths, err := filepath.Glob(filepath.Join(m.path, "*.json"))
	if err != nil {
		return fmt.Errorf("could not find json files to load: %s", err)
	}

	pbjson := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: false,
	}

	for _, path := range paths {
		var msg proto.Message
		key := strings.TrimSuffix(filepath.Base(path), ".json")
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("could not read %s: %s", key, err)
		}

		if strings.HasPrefix(key, vaspPrefix) {
			vasp := &v1beta1.VASP{}
			if err = pbjson.Unmarshal(data, vasp); err != nil {
				return err
			}
			msg = vasp
		} else if strings.HasPrefix(key, careqPrefix) {
			careq := &v1beta1gds.CertificateRequest{}
			if err = pbjson.Unmarshal(data, careq); err != nil {
				return err
			}
			msg = careq
		} else {
			return fmt.Errorf("unknown key format %q", key)
		}

		if err = d.Dump(key, msg); err != nil {
			return err
		}
	}
	return nil
}

func (m *leveldbDumper) Dump(key string, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not marshal protobuf: %s", err)
	}

	if err = m.db.Put([]byte(key), data, nil); err != nil {
		return fmt.Errorf("could not put %q: %s", key, err)
	}
	return nil
}

func convertVASP(vasp *v1alpha1.VASP) (*v1beta1.VASP, error) {
	// These values have changed from v1alpha1 to v1beta1 so we have to store them and unset from original
	trixo_safeguards_pii := vasp.Trixo.SafeguardsPii
	vasp.Trixo.SafeguardsPii = ""

	version := vasp.Version
	vasp.Version = 0

	vaspCategory := vasp.VaspCategory
	vasp.VaspCategory = 0

	// Deep clone using JSON intermediate -- very bad idea!
	// HACK: this is gross, don't repeat, for sanitation purposes only
	// Only works if names haven't changed - will work with field number changes ... hopefully
	// Will ignore extra fields ... hopefully
	// Will respect ENUM value strings ... hopefully
	mopts := protojson.MarshalOptions{
		Multiline:       false,
		AllowPartial:    false,
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: false,
	}
	data, err := mopts.Marshal(vasp)
	if err != nil {
		return nil, fmt.Errorf("could not serialize v1alpha1 message: %s", err)
	}

	uopts := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}
	converted := &v1beta1.VASP{}
	if err = uopts.Unmarshal(data, converted); err != nil {
		return nil, fmt.Errorf("could not deserialize v1beta1 message: %s", err)
	}

	// Set converted values onto the VASP record
	converted.Trixo.SafeguardsPii = parseBool(trixo_safeguards_pii)
	converted.Trixo.ComplianceThresholdCurrency = "USD"
	converted.Trixo.KycThresholdCurrency = "USD"
	converted.Version = &v1beta1.Version{Version: version}
	converted.VaspCategories = convertVaspCategory(vaspCategory)
	v1beta1gds.SetAdminVerificationToken(converted, vasp.AdminVerificationToken)

	if vasp.Contacts.Technical != nil && vasp.Contacts.Technical.Email != "" {
		v1beta1gds.SetContactVerification(converted.Contacts.Technical, vasp.Contacts.Technical.Token, vasp.Contacts.Technical.Verified)
	}

	if vasp.Contacts.Administrative != nil && vasp.Contacts.Administrative.Email != "" {
		v1beta1gds.SetContactVerification(converted.Contacts.Administrative, vasp.Contacts.Administrative.Token, vasp.Contacts.Administrative.Verified)
	}

	if vasp.Contacts.Legal != nil && vasp.Contacts.Legal.Email != "" {
		v1beta1gds.SetContactVerification(converted.Contacts.Legal, vasp.Contacts.Legal.Token, vasp.Contacts.Legal.Verified)
	}

	if vasp.Contacts.Billing != nil && vasp.Contacts.Billing.Email != "" {
		v1beta1gds.SetContactVerification(converted.Contacts.Billing, vasp.Contacts.Billing.Token, vasp.Contacts.Billing.Verified)
	}

	// Special ENUM case - due to zero valued enums
	if vasp.Entity != nil {
		for idx, nameID := range vasp.Entity.Name.NameIdentifiers {
			if nameID.LegalPersonNameIdentifierType == 0 {
				converted.Entity.Name.NameIdentifiers[idx].LegalPersonNameIdentifierType = 1
			}
		}
		for idx, nameID := range vasp.Entity.Name.LocalNameIdentifiers {
			if nameID.LegalPersonNameIdentifierType == 0 {
				converted.Entity.Name.NameIdentifiers[idx].LegalPersonNameIdentifierType = 1
			}
		}
		for idx, nameID := range vasp.Entity.Name.PhoneticNameIdentifiers {
			if nameID.LegalPersonNameIdentifierType == 0 {
				converted.Entity.Name.NameIdentifiers[idx].LegalPersonNameIdentifierType = 1
			}
		}
	}

	return converted, nil
}

func convertCertReq(careq *v1alpha1.CertificateRequest) (*v1beta1gds.CertificateRequest, error) {
	// Deep clone using JSON intermediate -- very bad idea!
	// HACK: this is gross, don't repeat, for sanitation purposes only
	// Only works if names haven't changed - will work with field number changes ... hopefully
	// Will ignore extra fields ... hopefully
	// Will respect ENUM value strings ... hopefully
	mopts := protojson.MarshalOptions{
		Multiline:       false,
		AllowPartial:    false,
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: false,
	}
	data, err := mopts.Marshal(careq)
	if err != nil {
		return nil, fmt.Errorf("could not serialize v1alpha1 message: %s", err)
	}

	uopts := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}
	converted := &v1beta1gds.CertificateRequest{}
	if err = uopts.Unmarshal(data, converted); err != nil {
		return nil, fmt.Errorf("could not deserialize v1beta1 message: %s", err)
	}
	return converted, nil
}

func parseBool(v string) bool {
	switch strings.ToLower(v) {
	case "y", "ye", "yes", "on", "1", "true":
		return true
	case "", "n", "no", "0", "off", "false":
		return false
	default:
		panic(fmt.Errorf("unknown bool value %q", v))
	}
}

func convertVaspCategory(v v1alpha1.VASPCategory) []string {
	switch v {
	case v1alpha1.VASPCategory_ATM:
		return []string{v1beta1.VASPCategoryKiosk}
	case v1alpha1.VASPCategory_EXCHANGE:
		return []string{v1beta1.VASPCategoryExchange}
	case v1alpha1.VASPCategory_HIGH_RISK_EXCHANGE:
		return []string{v1beta1.VASPCategoryExchange}
	default:
		return []string{}
	}
}
