package store

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

// OpenLevelDB directory Store at the specified path. This is the default storage provider.
func OpenLevelDB(uri string) (Store, error) {
	dsn, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("could not parse leveldb uri: %s", err)
	}

	db, err := leveldb.OpenFile(dsn.Path, nil)
	if err != nil {
		return nil, err
	}

	store := &ldbStore{db: db}
	if err = store.sync(); err != nil {
		return nil, err
	}

	// Perform a reindex if the local indices are null or empty. In the case where the
	// store has no data, this won't be harmful - but in the case where the stored index
	// has been corrupted, this should repair it.
	if len(store.names) == 0 || len(store.websites) == 0 || len(store.countries) == 0 || len(store.categories) == 0 {
		log.Info().Msg("reindexing to recover from empty indices")
		if err = store.Reindex(); err != nil {
			return nil, err
		}
	}

	return store, nil
}

// Errors that may occur during LevelDB operations
var (
	ErrCorruptedSequence = errors.New("primary key sequence is invalid")
	ErrCorruptedIndex    = errors.New("search indices are invalid")
	ErrIncompleteRecord  = errors.New("record is missing required fields")
	ErrEntityNotFound    = errors.New("entity not found")
	ErrDuplicateEntity   = errors.New("entity unique constraints violated")
)

// keys and prefixes for leveldb buckets and indices
var (
	keyAutoSequence  = []byte("sequence::pks")
	keyNameIndex     = []byte("index::names")
	keyWebsiteIndex  = []byte("index::websites")
	keyCountryIndex  = []byte("index::countries")
	keyCategoryIndex = []byte("index::categories")
	preVASPs         = []byte("vasps::")
	preCertReqs      = []byte("certreqs::")
)

// Implements Store for some basic LevelDB operations and simple protocol buffer storage.
type ldbStore struct {
	sync.RWMutex
	db         *leveldb.DB
	pkseq      sequence       // autoincrement sequence for ID values
	names      uniqueIndex    // case insensitive name index
	websites   uniqueIndex    // website/url index
	countries  containerIndex // lookup vasps in a specific country
	categories containerIndex // lookup vasps based on specified categories
}

//===========================================================================
// Store Implementation
//===========================================================================

// Close the database, allowing no further interactions. This method also synchronizes
// the indices to ensure that they are saved between sessions.
func (s *ldbStore) Close() error {
	defer s.db.Close()
	if err := s.sync(); err != nil {
		return err
	}
	return nil
}

//===========================================================================
// DirectoryStore Implementation
//===========================================================================

// Create a VASP into the directory. This method requires the VASP to have a unique
// name and ignores any ID fields that are set on the VASP, instead assigning new IDs.
func (s *ldbStore) Create(v *pb.VASP) (id string, err error) {
	// Create UUID for record
	// TODO: check uniqueness of the ID
	v.Id = uuid.New().String()

	// Ensure a common name exists for the uniqueness constraint
	// NOTE: other validation should have been performed in advance
	if name := normalize(v.CommonName); name == "" {
		return "", ErrIncompleteRecord
	}

	// Update management timestamps and record metadata
	v.LastUpdated = time.Now().Format(time.RFC3339)
	if v.FirstListed == "" {
		v.FirstListed = v.LastUpdated
	}
	if v.Version == 0 {
		v.Version = 1
	}

	// Critical section (optimizing for safety rather than speed)
	s.Lock()
	defer s.Unlock()

	// Check the uniqueness constraint
	if _, ok := s.names.find(v.CommonName, normalize); ok {
		return "", ErrDuplicateEntity
	}

	var data []byte
	key := s.vaspKey(v.Id)
	if data, err = proto.Marshal(v); err != nil {
		return "", err
	}

	if err = s.db.Put(key, data, nil); err != nil {
		return "", err
	}

	// Update indices after successful insert
	if err = s.insertIndices(v); err != nil {
		return "", err
	}
	return v.Id, nil
}

// Retrieve a VASP record by id; returns an error if the record does not exist.
func (s *ldbStore) Retrieve(id string) (v *pb.VASP, err error) {
	var val []byte
	key := s.vaspKey(id)
	if val, err = s.db.Get(key, nil); err != nil {
		if err == leveldb.ErrNotFound {
			return v, ErrEntityNotFound
		}
		return v, err
	}

	v = new(pb.VASP)
	if err = proto.Unmarshal(val, v); err != nil {
		return v, err
	}

	return v, nil
}

// Update the VASP entry by the VASP ID (required). This method simply overwrites the
// entire VASP record and does not update individual fields.
func (s *ldbStore) Update(v *pb.VASP) (err error) {
	if v.Id == "" {
		return ErrIncompleteRecord
	}

	// Ensure a common name exists for the uniqueness constraint
	// NOTE: other validation should have been performed in advance
	if name := normalize(v.CommonName); name == "" {
		return ErrIncompleteRecord
	}

	// Retrieve the original record to ensure that the indices are updated properly
	key := s.vaspKey(v.Id)
	o, err := s.Retrieve(v.Id)
	if err != nil {
		return err
	}

	// Update the record metadata
	v.Version++
	v.LastUpdated = time.Now().Format(time.RFC3339)
	if v.FirstListed == "" {
		v.FirstListed = v.LastUpdated
	}

	// Critical section (optimizing for safety rather than speed)
	s.Lock()
	defer s.Unlock()

	var val []byte
	if val, err = proto.Marshal(v); err != nil {
		return err
	}

	// Insert the new record
	if err = s.db.Put(key, val, nil); err != nil {
		return err
	}

	// Update indices to match new record, removing the old indices and updating the
	// new indices with any changes to ensure the index correctly reflects the state.
	// TODO: check errors from index removal in case the index must be reset.
	s.removeIndices(o)
	if err = s.insertIndices(v); err != nil {
		return err
	}
	return nil
}

// Destroy a record, removing it completely from the database and indices.
func (s *ldbStore) Destroy(id string) (err error) {
	key := s.vaspKey(id)

	// Lookup the record in order to remove data from indices
	record, err := s.Retrieve(id)
	if err != nil {
		if err == ErrEntityNotFound {
			return nil
		}
		return err
	}

	// LevelDB will not return an error if the entity does not exist
	if err = s.db.Delete(key, nil); err != nil {
		return err
	}

	// Critical section (optimizing for safety rather than speed)
	s.Lock()
	defer s.Unlock()

	// Remove the records from the indices
	if err = s.removeIndices(record); err != nil {
		return err
	}
	return nil
}

// Search uses the names and countries index to find VASPS that match the specified
// query. This is a very simple search and is not intended for robust usage. To find a
// VASP by name, a case insensitive search is performed if the query exists in
// any of the VASP entity names. Alternatively a list of names can be given or a country
// or list of countries for case-insensitive exact matches.
func (s *ldbStore) Search(query map[string]interface{}) (vasps []*pb.VASP, err error) {
	// A set of records that match the query and need to be fetched
	records := make(map[string]struct{})

	s.RLock()
	// Lookup by name
	names, ok := parseQuery("name", query, normalize)
	if ok {
		log.Debug().Strs("name", names).Msg("search name query")
		for _, name := range names {
			if id := s.names[name]; id != "" {
				records[id] = struct{}{}
			}
		}
	}

	// Lookup by website
	websites, ok := parseQuery("website", query, normalizeURL)
	if ok {
		log.Debug().Strs("website", websites).Msg("search website query")
		for _, website := range websites {
			if id := s.websites[website]; id != "" {
				records[id] = struct{}{}
			}
		}
	}

	// Filter by country
	// NOTE: if country is not in the index, no records will be returned
	countries, ok := parseQuery("country", query, normalizeCountry)
	if ok {
		for _, country := range countries {
			for record := range records {
				if !s.countries.contains(country, record, nil) {
					// Remove the found VASP since it is not in the country index
					// NOTE: safe to remove during map iteration: https://stackoverflow.com/questions/23229975/is-it-safe-to-remove-selected-keys-from-map-within-a-range-loop
					delete(records, record)
				}
			}
		}
	}

	// Filter by category
	// NOTE: if category is not in the index, no records will be returned
	categories, ok := parseQuery("category", query, normalize)
	if ok {
		for _, category := range categories {
			for record := range records {
				if !s.categories.contains(category, record, nil) {
					// Remove the found VASP since it is not in the country index
					// NOTE: safe to remove during map iteration: https://stackoverflow.com/questions/23229975/is-it-safe-to-remove-selected-keys-from-map-within-a-range-loop
					delete(records, record)
				}
			}
		}
	}

	s.RUnlock()

	// Perform the lookup of records if there are any
	if len(records) > 0 {
		vasps = make([]*pb.VASP, 0, len(records))
		for id := range records {
			var vasp *pb.VASP
			if vasp, err = s.Retrieve(id); err != nil {
				if err == ErrEntityNotFound {
					continue
				}
				return nil, err
			}
			vasps = append(vasps, vasp)
		}
	}

	return vasps, nil
}

//===========================================================================
// CertificateStore Implementation
//===========================================================================

// ListCertRequests returns all certificate requests that are currently in the store.
func (s *ldbStore) ListCertRequests() (reqs []*pb.CertificateRequest, err error) {
	reqs = make([]*pb.CertificateRequest, 0)
	iter := s.db.NewIterator(util.BytesPrefix(preCertReqs), nil)
	defer iter.Release()
	for iter.Next() {
		r := new(pb.CertificateRequest)
		if err = proto.Unmarshal(iter.Value(), r); err != nil {
			return nil, err
		}
		reqs = append(reqs, r)
	}

	if err = iter.Error(); err != nil {
		return nil, err
	}

	return reqs, nil
}

// GetCertRequest returns a certificate request by certificate request ID.
func (s *ldbStore) GetCertRequest(id string) (r *pb.CertificateRequest, err error) {
	if id == "" {
		return nil, ErrEntityNotFound
	}

	var val []byte
	if val, err = s.db.Get(s.careqKey(id), nil); err != nil {
		if err == leveldb.ErrNotFound {
			return nil, ErrEntityNotFound
		}
		return nil, err
	}

	r = new(pb.CertificateRequest)
	if err = proto.Unmarshal(val, r); err != nil {
		return nil, err
	}

	return r, nil
}

// SaveCertRequest can create or update a certificate request. The request should be as
// complete as possible, including an ID generated by the caller.
func (s *ldbStore) SaveCertRequest(r *pb.CertificateRequest) (err error) {
	if r.Id == "" {
		return ErrIncompleteRecord
	}

	// Update management timestamps and record metadata
	r.Modified = time.Now().Format(time.RFC3339)
	if r.Created == "" {
		r.Created = r.Modified
	}

	var data []byte
	key := s.careqKey(r.Id)
	if data, err = proto.Marshal(r); err != nil {
		return err
	}

	if err = s.db.Put(key, data, nil); err != nil {
		return err
	}

	return nil
}

// DeleteCertRequest removes a certificate request from the store.
func (s *ldbStore) DeleteCertRequest(id string) (err error) {
	// LevelDB will not return an error if the entity does not exist
	key := s.careqKey(id)
	if err = s.db.Delete(key, nil); err != nil {
		return err
	}
	return nil
}

//===========================================================================
// Key Handlers
//===========================================================================

// creates a []byte key from the vasp id using a prefix to act as a leveldb bucket
func (s *ldbStore) makeKey(prefix []byte, id string) (key []byte) {
	buf := []byte(id)
	key = make([]byte, 0, len(prefix)+len(buf))
	key = append(key, prefix...)
	key = append(key, buf...)
	return key
}

// creates a []byte key from the vasp id using a prefix to act as a leveldb bucket
func (s *ldbStore) vaspKey(id string) (key []byte) {
	return s.makeKey(preVASPs, id)
}

// creates a []byte key from the cert request id using a prefix to act as a leveldb bucket
func (s *ldbStore) careqKey(id string) (key []byte) {
	return s.makeKey(preCertReqs, id)
}

//===========================================================================
// Indexer
//===========================================================================

// Reindex rebuilds the name and country indices for the server and synchronizes them
// back to disk to ensure they're complete and accurate.
func (s *ldbStore) Reindex() (err error) {
	names := make(uniqueIndex)
	websites := make(uniqueIndex)
	countries := make(containerIndex)
	categories := make(containerIndex)

	iter := s.db.NewIterator(util.BytesPrefix(preVASPs), nil)
	defer iter.Release()

	for iter.Next() {
		vasp := new(pb.VASP)
		if err = proto.Unmarshal(iter.Value(), vasp); err != nil {
			return err
		}

		names.add(vasp.CommonName, vasp.Id, normalize)
		for _, name := range vasp.Entity.Names() {
			names.add(name, vasp.Id, normalize)
		}

		websites.add(vasp.Website, vasp.Id, normalizeURL)
		countries.add(vasp.Entity.CountryOfRegistration, vasp.Id, normalizeCountry)
		categories.add(vasp.BusinessCategory.String(), vasp.Id, normalize)
		categories.add(vasp.VaspCategory.String(), vasp.Id, normalize)
	}

	if err = iter.Error(); err != nil {
		return err
	}

	s.Lock()
	if len(names) > 0 {
		s.names = names
	}

	if len(websites) > 0 {
		s.websites = websites
	}

	if len(countries) > 0 {
		s.countries = countries
	}

	if len(categories) > 0 {
		s.categories = categories
	}
	s.Unlock()

	if err = s.sync(); err != nil {
		return err
	}

	log.Debug().
		Int("names", len(s.names)).
		Int("websites", len(s.websites)).
		Int("countries", len(s.countries)).
		Int("categories", len(s.categories)).
		Msg("reindex complete")
	return nil
}

//===========================================================================
// Backup
//===========================================================================

// Backup copies the leveldb database to a new directory and archives it as gzip tar.
// See: https://github.com/wbolster/plyvel/issues/46
func (s *ldbStore) Backup(path string) (err error) {
	// Create the directory for the copied leveldb database
	archive := filepath.Join(path, time.Now().UTC().Format("gdsdb-200601021504"))
	if err = os.Mkdir(archive, 0744); err != nil {
		return fmt.Errorf("could not create archive directory: %s", err)
	}

	// Ensure the archive directory is cleaned up when the backup is complete
	defer func() {
		os.RemoveAll(archive)
	}()

	// Open a second leveldb database at the backup location
	arcdb, err := leveldb.OpenFile(archive, nil)
	if err != nil {
		return fmt.Errorf("could not open archive database: %s", err)
	}

	// Create a new batch write to the archive database, writing every 100 records as
	// we iterate over all of the data in the store database.
	var nrows, narchived uint64
	batch := new(leveldb.Batch)
	iter := s.db.NewIterator(nil, nil)
	for iter.Next() {
		nrows++
		batch.Put(iter.Key(), iter.Value())

		if nrows%100 == 0 {
			if err = arcdb.Write(batch, &opt.WriteOptions{Sync: true}); err != nil {
				return fmt.Errorf("could not write next 100 rows after %d rows: %s", narchived, err)
			}
			batch.Reset()
			narchived += 100
		}
	}

	// Release the iterator and check for errors, just in case we didn't write anything
	iter.Release()
	if err = iter.Error(); err != nil {
		return fmt.Errorf("could not iterate over gds store: %s", err)
	}

	// Write final rows to the database
	if err = arcdb.Write(batch, &opt.WriteOptions{Sync: true}); err != nil {
		return fmt.Errorf("could not write final %d rows after %d rows: %s", nrows-narchived, narchived, err)
	}
	batch.Reset()
	narchived += (nrows - narchived)
	log.Info().Uint64("records", narchived).Msg("leveldb archive complete")

	// Close the archive database
	if err = arcdb.Close(); err != nil {
		return fmt.Errorf("could not close archive database: %s", err)
	}

	// Create the compressed tar archive with tar and gzip
	out, err := os.OpenFile(archive+".tgz", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("could not create archive file: %s", err)
	}

	// Create write chains: gzip writes to disk, tar writes to gzip
	gz := gzip.NewWriter(out)
	tw := tar.NewWriter(gz)

	// Walk the archive directory, removing the prefix
	prefix := filepath.Dir(archive) + "/"
	err = filepath.Walk(archive, func(file string, fi os.FileInfo, _ error) error {
		// Generate tar header
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		// Provide real name without directory prefix (so that archive is self-contained)
		// See: https://golang.org/src/archive/tar/common.go?#L626
		header.Name = filepath.ToSlash(strings.TrimPrefix(file, prefix))

		// Write header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// If not a directory, write the file content
		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}

			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("could not walk and tar archive directory: %s", err)
	}

	// Produce the tar
	if err = tw.Close(); err != nil {
		return fmt.Errorf("could not create tar: %s", err)
	}

	// Produce the gzip
	if err = gz.Close(); err != nil {
		return fmt.Errorf("could not create gzip tar: %s", err)
	}

	return nil
}

//===========================================================================
// Indices and Synchronization
//===========================================================================

func (s *ldbStore) insertIndices(v *pb.VASP) (err error) {
	s.names.add(v.CommonName, v.Id, normalize)
	for _, name := range v.Entity.Names() {
		s.names.add(name, v.Id, normalize)
	}

	s.websites.add(v.Website, v.Id, normalizeURL)

	s.countries.add(v.Entity.CountryOfRegistration, v.Id, normalizeCountry)

	s.categories.add(v.BusinessCategory.String(), v.Id, normalize)
	s.categories.add(v.VaspCategory.String(), v.Id, normalize)
	return nil
}

func (s *ldbStore) removeIndices(v *pb.VASP) (err error) {
	s.names.rm(v.CommonName, normalize)
	for _, name := range v.Entity.Names() {
		s.names.rm(name, normalize)
	}

	s.websites.rm(v.Website, normalizeURL)

	s.countries.rm(v.Entity.CountryOfRegistration, v.Id, normalizeCountry)

	s.categories.rm(v.BusinessCategory.String(), v.Id, normalize)
	s.categories.rm(v.VaspCategory.String(), v.Id, normalize)
	return nil
}

// sync all indices with the underlying database
func (s *ldbStore) sync() (err error) {
	if err = s.seqsync(); err != nil {
		return err
	}

	if err = s.syncnames(); err != nil {
		return err
	}

	if err = s.syncwebsites(); err != nil {
		return err
	}

	if err = s.synccountries(); err != nil {
		return err
	}

	if err = s.synccategories(); err != nil {
		return err
	}

	log.Debug().
		Int("names", len(s.names)).
		Int("websites", len(s.websites)).
		Int("countries", len(s.countries)).
		Int("categories", len(s.categories)).
		Msg("indices synchronized")
	return nil
}

// sync the autoincrement sequence with the leveldb auto sequence key
func (s *ldbStore) seqsync() (err error) {
	var pk sequence
	var data []byte
	if data, err = s.db.Get(keyAutoSequence, nil); err != nil {
		// If the auto sequence key is not found, simply leave pk to 0
		if err != leveldb.ErrNotFound {
			return err
		}
	} else {
		if pk, err = pk.Load(data); err != nil {
			return err
		}
	}

	// Critical section (optimizing for safety rather than speed)
	s.Lock()
	defer s.Unlock()

	// Local is behind database state, set and return
	if s.pkseq <= pk {
		s.pkseq = pk
		log.Debug().Uint64("sequence", uint64(s.pkseq)).Msg("updated primary key sequence from cache")
		return nil
	}

	//  Update the database with the local state
	if data, err = s.pkseq.Dump(); err != nil {
		log.Error().Err(err).Msg("could not put primary key sequence value")
		return err
	}
	if err = s.db.Put(keyAutoSequence, data, nil); err != nil {
		log.Error().Err(err).Msg("could not put primary key sequence value")
		return ErrCorruptedSequence
	}

	log.Debug().Uint64("sequence", uint64(s.pkseq)).Msg("cached primary key sequence to disk")
	return nil
}

// sync the names index with the leveldb names key
func (s *ldbStore) syncnames() (err error) {
	var val []byte

	// Critical section (optimizing for safety rather than speed)
	s.Lock()
	defer s.Unlock()

	if s.names == nil {
		// Create the index to load it from disk
		s.names = make(uniqueIndex)

		// fetch the names from the database
		if val, err = s.db.Get(keyNameIndex, nil); err != nil {
			if err == leveldb.ErrNotFound {
				return nil
			}
			log.Error().Err(err).Msg("could not fetch names index from database")
			return err
		}

		if err = s.names.Load(val); err != nil {
			log.Error().Err(err).Msg("could not unmarshal names index")
			return ErrCorruptedIndex
		}
	}

	// Put the current names back to the database
	if len(s.names) > 0 {
		if val, err = s.names.Dump(); err != nil {
			log.Error().Err(err).Msg("could not marshal names index")
			return ErrCorruptedIndex
		}

		if err = s.db.Put(keyNameIndex, val, nil); err != nil {
			log.Error().Err(err).Msg("could not put names index")
			return ErrCorruptedIndex
		}
	}

	log.Debug().Int("size", len(val)).Msg("names index checkpointed")
	return nil
}

// sync the websites index with the leveldb websites key
func (s *ldbStore) syncwebsites() (err error) {
	var val []byte

	// Critical section (optimizing for safety rather than speed)
	s.Lock()
	defer s.Unlock()

	if s.websites == nil {
		// Create the index to load it from disk
		s.websites = make(uniqueIndex)

		// fetch the websites from the database
		if val, err = s.db.Get(keyWebsiteIndex, nil); err != nil {
			if err == leveldb.ErrNotFound {
				return nil
			}
			log.Error().Err(err).Msg("could not fetch websites index from database")
			return err
		}

		if err = s.websites.Load(val); err != nil {
			log.Error().Err(err).Msg("could not unmarshal websites index")
			return ErrCorruptedIndex
		}
	}

	// Put the current websites back to the database
	if len(s.websites) > 0 {
		if val, err = s.websites.Dump(); err != nil {
			log.Error().Err(err).Msg("could not marshal websites index")
			return ErrCorruptedIndex
		}

		if err = s.db.Put(keyWebsiteIndex, val, nil); err != nil {
			log.Error().Err(err).Msg("could not put websites index")
			return ErrCorruptedIndex
		}
	}

	log.Debug().Int("size", len(val)).Msg("websites index checkpointed")
	return nil
}

// sync the countries index with the leveldb countries key
func (s *ldbStore) synccountries() (err error) {
	var val []byte

	// Critical section (optimizing for safety rather than speed)
	s.Lock()
	defer s.Unlock()

	if s.countries == nil {
		// Create the countries index an dload from the database
		s.countries = make(containerIndex)

		// fetch the countries from the database
		if val, err = s.db.Get(keyCountryIndex, nil); err != nil {
			if err == leveldb.ErrNotFound {
				return nil
			}
			log.Error().Err(err).Msg("could fetch country index from database")
			return err
		}

		if err = s.countries.Load(val); err != nil {
			log.Error().Err(err).Msg("could not unmarshall country index")
			return ErrCorruptedIndex
		}
	}

	if len(s.countries) > 0 {
		// Put the current countries back to the database
		if val, err = s.countries.Dump(); err != nil {
			log.Error().Err(err).Msg("could not marshal country index")
			return ErrCorruptedIndex
		}

		if err = s.db.Put(keyCountryIndex, val, nil); err != nil {
			log.Error().Err(err).Msg("could not put country index")
			return ErrCorruptedIndex
		}
	}

	log.Debug().Int("size", len(val)).Msg("country index checkpointed")
	return nil
}

// sync the categories index with the leveldb categories key
func (s *ldbStore) synccategories() (err error) {
	var val []byte

	// Critical section (optimizing for safety rather than speed)
	s.Lock()
	defer s.Unlock()

	if s.categories == nil {
		// Create the categories index an dload from the database
		s.categories = make(containerIndex)

		// fetch the categories from the database
		if val, err = s.db.Get(keyCategoryIndex, nil); err != nil {
			if err == leveldb.ErrNotFound {
				return nil
			}
			log.Error().Err(err).Msg("could fetch categories index from database")
			return err
		}

		if err = s.categories.Load(val); err != nil {
			log.Error().Err(err).Msg("could not unmarshall categories index")
			return ErrCorruptedIndex
		}
	}

	if len(s.categories) > 0 {
		// Put the current categories back to the database
		if val, err = s.categories.Dump(); err != nil {
			log.Error().Err(err).Msg("could not marshal categories index")
			return ErrCorruptedIndex
		}

		if err = s.db.Put(keyCategoryIndex, val, nil); err != nil {
			log.Error().Err(err).Msg("could not put categories index")
			return ErrCorruptedIndex
		}
	}

	log.Debug().Int("size", len(val)).Msg("categories index checkpointed")
	return nil
}
