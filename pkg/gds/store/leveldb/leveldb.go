package leveldb

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	storeerrors "github.com/trisacrypto/directory/pkg/gds/store/errors"
	"github.com/trisacrypto/directory/pkg/gds/store/index"
	"github.com/trisacrypto/directory/pkg/gds/store/iterator"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/proto"
)

// Open LevelDB directory Store at the specified path.
func Open(path string) (*Store, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	store := &Store{db: db}
	if err = store.sync(); err != nil {
		return nil, err
	}

	// Perform a reindex if the local indices are null or empty. In the case where the
	// store has no data, this won't be harmful - but in the case where the stored index
	// has been corrupted, this should repair it.
	if store.names.Empty() || store.websites.Empty() || store.countries.Empty() || store.categories.Empty() {
		log.Info().Msg("reindexing to recover from empty indices")
		if err = store.Reindex(); err != nil {
			return nil, err
		}
	}

	return store, nil
}

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

// Store implements store.Store for some basic LevelDB operations and simple protocol
// buffer storage in a key/value database.
type Store struct {
	sync.RWMutex
	db         *leveldb.DB
	pkseq      index.Sequence    // autoincrement sequence for ID values
	names      index.SingleIndex // case insensitive name index
	websites   index.SingleIndex // website/url index
	countries  index.MultiIndex  // lookup vasps in a specific country
	categories index.MultiIndex  // lookup vasps based on specified categories
}

//===========================================================================
// Store Implementation
//===========================================================================

// Close the database, allowing no further interactions. This method also synchronizes
// the indices to ensure that they are saved between sessions.
func (s *Store) Close() error {
	defer s.db.Close()
	if err := s.sync(); err != nil {
		return err
	}
	return nil
}

//===========================================================================
// DirectoryStore Implementation
//===========================================================================

// CreateVASP into the directory. This method requires the VASP to have a unique
// name and ignores any ID fields that are set on the VASP, instead assigning new IDs.
func (s *Store) CreateVASP(v *pb.VASP) (id string, err error) {
	// Create UUID for record
	if v.Id == "" {
		v.Id = uuid.New().String()
	}
	key := vaspKey(v.Id)

	// Ensure a common name exists for the uniqueness constraint
	// NOTE: other validation should have been performed in advance
	if name := index.Normalize(v.CommonName); name == "" {
		return "", storeerrors.ErrIncompleteRecord
	}

	// Update management timestamps and record metadata
	v.LastUpdated = time.Now().Format(time.RFC3339)
	if v.FirstListed == "" {
		v.FirstListed = v.LastUpdated
	}
	if v.Version == nil || v.Version.Version == 0 {
		v.Version = &pb.Version{Version: 1}
	}

	// Critical section (optimizing for safety rather than speed)
	s.Lock()
	defer s.Unlock()

	// Check the uniqueness constraints
	// NOTE: website removed as uniqueness constraint in SC-4483
	if id, ok := s.names.Find(v.CommonName); ok && id != v.Id {
		return "", storeerrors.ErrDuplicateEntity
	}

	// It is not necessary for the marshal to be inside the lock, but we don't want to
	// do extra serialization work in memory if there is a duplicate entity, a check
	// which must be inside the lock.
	var data []byte
	if data, err = proto.Marshal(v); err != nil {
		return "", err
	}

	// This Put must be inside the lock to ensure the indices reflect what is in the db.
	if err = s.db.Put(key, data, nil); err != nil {
		return "", err
	}

	// Update indices after successful insert
	if err = s.insertIndices(v); err != nil {
		return "", err
	}
	return v.Id, nil
}

// RetrieveVASP record by id; returns an error if the record does not exist.
func (s *Store) RetrieveVASP(id string) (v *pb.VASP, err error) {
	var val []byte
	key := vaspKey(id)
	if val, err = s.db.Get(key, nil); err != nil {
		if err == leveldb.ErrNotFound {
			return nil, storeerrors.ErrEntityNotFound
		}
		return nil, err
	}

	v = new(pb.VASP)
	if err = proto.Unmarshal(val, v); err != nil {
		return nil, err
	}

	return v, nil
}

// UpdateVASP by the VASP ID (required). This method simply overwrites the
// entire VASP record and does not update individual fields.
func (s *Store) UpdateVASP(v *pb.VASP) (err error) {
	if v.Id == "" {
		return storeerrors.ErrIncompleteRecord
	}
	key := vaspKey(v.Id)

	// Ensure a common name exists for the uniqueness constraint
	// NOTE: other validation should have been performed in advance
	if name := index.Normalize(v.CommonName); name == "" {
		return storeerrors.ErrIncompleteRecord
	}

	// Update management timestamps and record metadata
	v.Version.Version++
	v.LastUpdated = time.Now().Format(time.RFC3339)
	if v.FirstListed == "" {
		v.FirstListed = v.LastUpdated
	}

	var val []byte
	if val, err = proto.Marshal(v); err != nil {
		return err
	}

	// Critical section (optimizing for safety rather than speed)
	s.Lock()
	defer s.Unlock()

	// Retrieve the original record to ensure that the indices are updated properly
	// This must be inside the lock so that the database indices are currently updated.
	o, err := s.RetrieveVASP(v.Id)
	if err != nil {
		return err
	}

	// Check the uniqueness constraints
	// NOTE: website removed as uniqueness constraint in SC-4483
	if id, ok := s.names.Find(v.CommonName); ok && id != v.Id {
		return storeerrors.ErrDuplicateEntity
	}

	// Insert the new record
	// This must be inside the lock so that the indices reflect what is currently in
	// the database and there is no race condition between the retrieve and put.
	if err = s.db.Put(key, val, nil); err != nil {
		return err
	}

	// Update indices to match new record, removing the old indices and updating the
	// new indices with any changes to ensure the index correctly reflects the state.
	if err = s.removeIndices(o); err != nil {
		// NOTE: if this error is triggered, admins may want to reindex the database
		log.Error().Err(err).Msg("could not remove previous indices on update: reindex required")
	}
	if err = s.insertIndices(v); err != nil {
		return err
	}
	return nil
}

// DeleteVASP record, removing it completely from the database and indices.
func (s *Store) DeleteVASP(id string) (err error) {
	key := vaspKey(id)

	// Critical section (optimizing for safety rather than speed)
	s.Lock()
	defer s.Unlock()

	// Lookup the record in order to remove data from indices, this must be inside the
	// lock to ensure the indices are correctly updated with what is on disk.
	record, err := s.RetrieveVASP(id)
	if err != nil {
		if err == storeerrors.ErrEntityNotFound {
			return nil
		}
		return err
	}

	// LevelDB will not return an error if the entity does not exist
	if err = s.db.Delete(key, nil); err != nil {
		return err
	}

	// Remove the records from the indices
	if err = s.removeIndices(record); err != nil {
		return err
	}
	return nil
}

// ListVASPs returns all of the VASPs in the database
func (s *Store) ListVASPs() iterator.DirectoryIterator {
	return &vaspIterator{
		iterWrapper{
			iter: s.db.NewIterator(util.BytesPrefix(preVASPs), nil),
		},
	}
}

// SearchVASPs uses the names and countries index to find VASPs that match the specified
// query. This is a very simple search and is not intended for robust usage. To find a
// VASP by name, a case insensitive search is performed if the query exists in
// any of the VASP entity names. If there is not an exact match a prefix lookup is used
// so long as the prefix > 3 characters. The search also looks up website matches by
// parsing urls to match hostnames rather than scheme or path. Finally the query is
// filtered by country and category.
func (s *Store) SearchVASPs(query map[string]interface{}) (vasps []*pb.VASP, err error) {
	// A set of records that match the query and need to be fetched
	records := make(map[string]struct{})

	s.RLock()
	// Search the name index
	for _, result := range s.names.Search(query) {
		records[result] = struct{}{}
	}

	// Lookup by website
	for _, result := range s.websites.Search(query) {
		records[result] = struct{}{}
	}

	// Filter by country
	// NOTE: if country is not in the index, no records will be returned
	countries, ok := index.ParseQuery("country", query, index.NormalizeCountry)
	if ok {
		for _, country := range countries {
			for record := range records {
				if !s.countries.Contains(country, record) {
					// Remove the found VASP since it is not in the country index
					// NOTE: safe to remove during map iteration: https://stackoverflow.com/questions/23229975/is-it-safe-to-remove-selected-keys-from-map-within-a-range-loop
					delete(records, record)
				}
			}
		}
	}

	// Filter by category
	// NOTE: if category is not in the index, no records will be returned
	categories, ok := index.ParseQuery("category", query, index.Normalize)
	if ok {
		for _, category := range categories {
			for record := range records {
				if !s.categories.Contains(category, record) {
					// Remove the found VASP since it is not in the category index
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
			if vasp, err = s.RetrieveVASP(id); err != nil {
				if err == storeerrors.ErrEntityNotFound {
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

// ListCertReqs returns all certificate requests that are currently in the store.
func (s *Store) ListCertReqs() iterator.CertificateIterator {
	return &certReqIterator{
		iterWrapper{
			iter: s.db.NewIterator(util.BytesPrefix(preCertReqs), nil),
		},
	}
}

// CreateCertReq and assign a new ID and return the version.
func (s *Store) CreateCertReq(r *models.CertificateRequest) (id string, err error) {
	if r.Id != "" {
		return "", storeerrors.ErrIDAlreadySet
	}

	// Create UUID for record
	// TODO: check uniqueness of the ID
	r.Id = uuid.New().String()

	// Update management timestamps and record metadata
	r.Created = time.Now().Format(time.RFC3339)
	if r.Modified == "" {
		r.Modified = r.Created
	}

	var data []byte
	key := careqKey(r.Id)
	if data, err = proto.Marshal(r); err != nil {
		return "", err
	}

	if err = s.db.Put(key, data, nil); err != nil {
		return "", err
	}

	return r.Id, nil
}

// RetrieveCertReq returns a certificate request by certificate request ID.
func (s *Store) RetrieveCertReq(id string) (r *models.CertificateRequest, err error) {
	if id == "" {
		return nil, storeerrors.ErrEntityNotFound
	}

	var val []byte
	if val, err = s.db.Get(careqKey(id), nil); err != nil {
		if err == leveldb.ErrNotFound {
			return nil, storeerrors.ErrEntityNotFound
		}
		return nil, err
	}

	r = new(models.CertificateRequest)
	if err = proto.Unmarshal(val, r); err != nil {
		return nil, err
	}

	return r, nil
}

// UpdateCertReq can create or update a certificate request. The request should be as
// complete as possible, including an ID generated by the caller.
func (s *Store) UpdateCertReq(r *models.CertificateRequest) (err error) {
	if r.Id == "" {
		return storeerrors.ErrIncompleteRecord
	}

	// Update management timestamps and record metadata
	r.Modified = time.Now().Format(time.RFC3339)
	if r.Created == "" {
		r.Created = r.Modified
	}

	var data []byte
	key := careqKey(r.Id)
	if data, err = proto.Marshal(r); err != nil {
		return err
	}

	if err = s.db.Put(key, data, nil); err != nil {
		return err
	}

	return nil
}

// DeleteCertReq removes a certificate request from the store.
func (s *Store) DeleteCertReq(id string) (err error) {
	// LevelDB will not return an error if the entity does not exist
	key := careqKey(id)
	if err = s.db.Delete(key, nil); err != nil {
		return err
	}
	return nil
}

//===========================================================================
// Key Handlers
//===========================================================================

// creates a []byte key from the vasp id using a prefix to act as a leveldb bucket
func makeKey(prefix []byte, id string) (key []byte) {
	buf := []byte(id)
	key = make([]byte, 0, len(prefix)+len(buf))
	key = append(key, prefix...)
	key = append(key, buf...)
	return key
}

// creates a []byte key from the vasp id using a prefix to act as a leveldb bucket
func vaspKey(id string) (key []byte) {
	return makeKey(preVASPs, id)
}

// creates a []byte key from the cert request id using a prefix to act as a leveldb bucket
func careqKey(id string) (key []byte) {
	return makeKey(preCertReqs, id)
}

//===========================================================================
// Indexer
//===========================================================================

// Reindex rebuilds the name and country indices for the server and synchronizes them
// back to disk to ensure they're complete and accurate.
func (s *Store) Reindex() (err error) {
	names := index.NewNamesIndex()
	websites := index.NewWebsiteIndex()
	countries := index.NewCountryIndex()
	categories := index.NewCategoryIndex()

	iter := s.db.NewIterator(util.BytesPrefix(preVASPs), nil)
	defer iter.Release()

	for iter.Next() {
		vasp := new(pb.VASP)
		if err = proto.Unmarshal(iter.Value(), vasp); err != nil {
			return err
		}

		// Update name index
		names.Add(vasp.CommonName, vasp.Id)
		for _, name := range vasp.Entity.Names() {
			names.Add(name, vasp.Id)
		}

		// Update website index
		websites.Add(vasp.Website, vasp.Id)

		// Update country index
		countries.Add(vasp.Entity.CountryOfRegistration, vasp.Id)
		for _, addr := range vasp.Entity.GeographicAddresses {
			countries.Add(addr.Country, vasp.Id)
		}

		// Update category index
		categories.Add(vasp.BusinessCategory.String(), vasp.Id)
		for _, vaspCategory := range vasp.VaspCategories {
			categories.Add(vaspCategory, vasp.Id)
		}
	}

	if err = iter.Error(); err != nil {
		return err
	}

	s.Lock()
	if !names.Empty() {
		s.names = names
	}

	if !websites.Empty() {
		s.websites = websites
	}

	if !countries.Empty() {
		s.countries = countries
	}

	if !categories.Empty() {
		s.categories = categories
	}
	s.Unlock()

	if err = s.sync(); err != nil {
		return err
	}

	log.Debug().
		Int("names", s.names.Len()).
		Int("websites", s.websites.Len()).
		Int("countries", s.countries.Len()).
		Int("categories", s.categories.Len()).
		Msg("reindex complete")
	return nil
}

//===========================================================================
// Backup
//===========================================================================

// Backup copies the leveldb database to a new directory and archives it as gzip tar.
// See: https://github.com/wbolster/plyvel/issues/46
func (s *Store) Backup(path string) (err error) {
	// Before backup ensure the indices are sync'd to disk
	if err = s.sync(); err != nil {
		return fmt.Errorf("could not synchronize indices prior to backup: %s", err)
	}

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

	var narchived uint64
	if narchived, err = CopyDB(s.db, arcdb); err != nil {
		return fmt.Errorf("could not write all records to archive database, wrote %d records: %s", narchived, err)
	}
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

// CopyDB is a utility function to copy all the records from one leveldb database
// object to another.
func CopyDB(src *leveldb.DB, dst *leveldb.DB) (ncopied uint64, err error) {
	// Create a new batch write to the destination database, writing every 100 records
	// as we iterate over all of the data in the source database.
	var nrows uint64
	batch := new(leveldb.Batch)
	iter := src.NewIterator(nil, nil)
	for iter.Next() {
		nrows++
		batch.Put(iter.Key(), iter.Value())

		if nrows%100 == 0 {
			if err = dst.Write(batch, &opt.WriteOptions{Sync: true}); err != nil {
				return ncopied, fmt.Errorf("could not write next 100 rows after %d rows: %s", ncopied, err)
			}
			batch.Reset()
			ncopied += 100
		}
	}

	// Release the iterator and check for errors, just in case we didn't write anything
	iter.Release()
	if err = iter.Error(); err != nil {
		return ncopied, fmt.Errorf("could not iterate over GDS store: %s", err)
	}

	// Write final rows to the database
	if err = dst.Write(batch, &opt.WriteOptions{Sync: true}); err != nil {
		return ncopied, fmt.Errorf("could not write final %d rows after %d rows: %s", nrows-ncopied, ncopied, err)
	}
	batch.Reset()
	ncopied += (nrows - ncopied)
	return ncopied, nil
}

//===========================================================================
// Indices and Synchronization
//===========================================================================

func (s *Store) insertIndices(v *pb.VASP) (err error) {
	s.names.Add(v.CommonName, v.Id)
	for _, name := range v.Entity.Names() {
		s.names.Add(name, v.Id)
	}

	s.websites.Add(v.Website, v.Id)

	s.countries.Add(v.Entity.CountryOfRegistration, v.Id)
	for _, addr := range v.Entity.GeographicAddresses {
		s.countries.Add(addr.Country, v.Id)
	}

	s.categories.Add(v.BusinessCategory.String(), v.Id)
	for _, vaspCategory := range v.VaspCategories {
		s.categories.Add(vaspCategory, v.Id)
	}

	return nil
}

func (s *Store) removeIndices(v *pb.VASP) (err error) {
	s.names.Remove(v.CommonName)
	for _, name := range v.Entity.Names() {
		s.names.Remove(name)
	}

	s.websites.Remove(v.Website)

	s.countries.Remove(v.Entity.CountryOfRegistration, v.Id)
	for _, addr := range v.Entity.GeographicAddresses {
		s.countries.Remove(addr.Country, v.Id)
	}

	s.categories.Remove(v.BusinessCategory.String(), v.Id)
	for _, vaspCategory := range v.VaspCategories {
		s.categories.Remove(vaspCategory, v.Id)
	}
	return nil
}

// sync all indices with the underlying database
func (s *Store) sync() (err error) {
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
		Int("names", s.names.Len()).
		Int("websites", s.websites.Len()).
		Int("countries", s.countries.Len()).
		Int("categories", s.categories.Len()).
		Msg("indices synchronized")
	return nil
}

// sync the autoincrement sequence with the leveldb auto sequence key
func (s *Store) seqsync() (err error) {
	var pk index.Sequence
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
		return storeerrors.ErrCorruptedSequence
	}

	log.Debug().Uint64("sequence", uint64(s.pkseq)).Msg("cached primary key sequence to disk")
	return nil
}

// sync the names index with the leveldb names key
func (s *Store) syncnames() (err error) {
	var val []byte

	// Critical section (optimizing for safety rather than speed)
	s.Lock()
	defer s.Unlock()

	if s.names == nil {
		// Create the index to load it from disk
		s.names = index.NewNamesIndex()

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
			return storeerrors.ErrCorruptedIndex
		}
	}

	// Put the current names back to the database
	if !s.names.Empty() {
		if val, err = s.names.Dump(); err != nil {
			log.Error().Err(err).Msg("could not marshal names index")
			return storeerrors.ErrCorruptedIndex
		}

		if err = s.db.Put(keyNameIndex, val, nil); err != nil {
			log.Error().Err(err).Msg("could not put names index")
			return storeerrors.ErrCorruptedIndex
		}
	}

	log.Debug().Int("size", len(val)).Msg("names index checkpointed")
	return nil
}

// sync the websites index with the leveldb websites key
func (s *Store) syncwebsites() (err error) {
	var val []byte

	// Critical section (optimizing for safety rather than speed)
	s.Lock()
	defer s.Unlock()

	if s.websites == nil {
		// Create the index to load it from disk
		s.websites = index.NewWebsiteIndex()

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
			return storeerrors.ErrCorruptedIndex
		}
	}

	// Put the current websites back to the database
	if !s.websites.Empty() {
		if val, err = s.websites.Dump(); err != nil {
			log.Error().Err(err).Msg("could not marshal websites index")
			return storeerrors.ErrCorruptedIndex
		}

		if err = s.db.Put(keyWebsiteIndex, val, nil); err != nil {
			log.Error().Err(err).Msg("could not put websites index")
			return storeerrors.ErrCorruptedIndex
		}
	}

	log.Debug().Int("size", len(val)).Msg("websites index checkpointed")
	return nil
}

// sync the countries index with the leveldb countries key
func (s *Store) synccountries() (err error) {
	var val []byte

	// Critical section (optimizing for safety rather than speed)
	s.Lock()
	defer s.Unlock()

	if s.countries == nil {
		// Create the countries index an dload from the database
		s.countries = index.NewCountryIndex()

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
			return storeerrors.ErrCorruptedIndex
		}
	}

	if !s.countries.Empty() {
		// Put the current countries back to the database
		if val, err = s.countries.Dump(); err != nil {
			log.Error().Err(err).Msg("could not marshal country index")
			return storeerrors.ErrCorruptedIndex
		}

		if err = s.db.Put(keyCountryIndex, val, nil); err != nil {
			log.Error().Err(err).Msg("could not put country index")
			return storeerrors.ErrCorruptedIndex
		}
	}

	log.Debug().Int("size", len(val)).Msg("country index checkpointed")
	return nil
}

// sync the categories index with the leveldb categories key
func (s *Store) synccategories() (err error) {
	var val []byte

	// Critical section (optimizing for safety rather than speed)
	s.Lock()
	defer s.Unlock()

	if s.categories == nil {
		// Create the categories index and load from the database
		s.categories = index.NewCategoryIndex()

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
			return storeerrors.ErrCorruptedIndex
		}
	}

	if !s.categories.Empty() {
		// Put the current categories back to the database
		if val, err = s.categories.Dump(); err != nil {
			log.Error().Err(err).Msg("could not marshal categories index")
			return storeerrors.ErrCorruptedIndex
		}

		if err = s.db.Put(keyCategoryIndex, val, nil); err != nil {
			log.Error().Err(err).Msg("could not put categories index")
			return storeerrors.ErrCorruptedIndex
		}
	}

	log.Debug().Int("size", len(val)).Msg("categories index checkpointed")
	return nil
}
