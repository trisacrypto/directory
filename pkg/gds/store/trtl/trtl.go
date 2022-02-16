package trtl

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/gds/client"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	storeerrors "github.com/trisacrypto/directory/pkg/gds/store/errors"
	"github.com/trisacrypto/directory/pkg/gds/store/index"
	"github.com/trisacrypto/directory/pkg/gds/store/iterator"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Open a connection to the Trtl database.
func Open(profile *client.TrtlProfile) (store *Store, err error) {
	store = &Store{}
	if store.conn, err = profile.Connect(); err != nil {
		return nil, err
	}
	store.client = pb.NewTrtlClient(store.conn)

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

	// Run background go routine to periodically checkpoint index to disk.
	// NOTE: the leveldb store does this in the backup go routine.
	// TODO: configure (enable/disable) this functionality and shutdown
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for {
			<-ticker.C
			if err := store.sync(); err != nil {
				log.Error().Err(err).Msg("could not synchronize indices")
			}
		}
	}()

	return store, nil
}

// Store implements the store.Store interface for the Trtl replicated database.
type Store struct {
	sync.RWMutex
	conn       *grpc.ClientConn
	client     pb.TrtlClient
	names      index.SingleIndex // case insensitive name index
	websites   index.SingleIndex // website/url index
	countries  index.MultiIndex  // lookup vasps in a specific country
	categories index.MultiIndex  // lookup vasps based on specified categories
}

func withContext(ctx context.Context) (context.Context, context.CancelFunc) {
	// TODO: Timeout should be configurable.
	return context.WithTimeout(ctx, time.Second*30)
}

//===========================================================================
// Store Implementation
//===========================================================================

// Close the connection to the database.
func (s *Store) Close() error {
	defer s.conn.Close()
	if err := s.sync(); err != nil {
		return err
	}
	return nil
}

//===========================================================================
// DirectoryStore Implementation
//===========================================================================

// ListVASPs returns an iterator over all VASPs in the database.
func (s *Store) ListVASPs() iterator.DirectoryIterator {
	return &vaspIterator{
		NewTrtlStreamingIterator(s.client, wire.NamespaceVASPs),
	}
}

func (s *Store) SearchVASPs(query map[string]interface{}) (vasps []*gds.VASP, err error) {
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
					// NOTE: safe to remove during map iteration
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
					// NOTE: safe to remove during map iteration
					delete(records, record)
				}
			}
		}
	}
	s.RUnlock()

	// Perform the lookup of records if there are any
	if len(records) > 0 {
		vasps = make([]*gds.VASP, 0, len(records))
		for id := range records {
			var vasp *gds.VASP
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

// CreateVASP into the directory. This method requires the VASP to have a unique
// name and ignores any ID fields that are set on the VASP, instead assigning new IDs.
func (s *Store) CreateVASP(v *gds.VASP) (id string, err error) {
	// Create UUID for record
	if v.Id == "" {
		v.Id = uuid.New().String()
	}
	key := []byte(v.Id)

	// Ensure a common name exists for the uniqueness constraint
	// NOTE: other validation should have been performed in advance, including a check
	// for common name presence, this is to ensure that the index is correctly updated.
	if cn := index.Normalize(v.CommonName); cn == "" {
		return "", storeerrors.ErrIncompleteRecord
	}

	// Update management timestamps and record metadata
	v.LastUpdated = time.Now().Format(time.RFC3339)
	if v.FirstListed == "" {
		v.FirstListed = v.LastUpdated
	}
	if v.Version == nil || v.Version.Version == 0 {
		v.Version = &gds.Version{Version: 1}
	}

	// Critical section (optimizing for safety rather than speed)
	// TODO: if trtl also managed the index, this section would be far faster
	s.Lock()
	defer s.Unlock()

	// Check the uniqueness constraint
	if _, ok := s.names.Find(v.CommonName); ok {
		fmt.Printf("%+v\n", s.names)
		return "", storeerrors.ErrDuplicateEntity
	}

	var data []byte
	if data, err = proto.Marshal(v); err != nil {
		return "", err
	}

	ctx, cancel := withContext(context.Background())
	defer cancel()
	request := &pb.PutRequest{
		Key:       key,
		Value:     data,
		Namespace: wire.NamespaceVASPs,
	}
	if reply, err := s.client.Put(ctx, request); err != nil || !reply.Success {
		if err == nil {
			err = storeerrors.ErrProtocol
		}
		return "", err
	}

	// Update indices after a successful insert
	if err = s.insertIndices(v); err != nil {
		return "", err
	}
	return v.Id, nil
}

// RetrieveVASP record by id. Returns ErrEntityNotFound if the record does not exist.
func (s *Store) RetrieveVASP(id string) (v *gds.VASP, err error) {
	key := []byte(id)

	ctx, cancel := withContext(context.Background())
	defer cancel()
	request := &pb.GetRequest{
		Key:       key,
		Namespace: wire.NamespaceVASPs,
	}
	var reply *pb.GetReply
	if reply, err = s.client.Get(ctx, request); err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, storeerrors.ErrEntityNotFound
		}
		return nil, err
	}

	v = new(gds.VASP)
	if err = proto.Unmarshal(reply.Value, v); err != nil {
		return nil, err
	}

	return v, nil
}

// UpdateVASP by the VASP ID (required). This method simply overwrites the
// entire VASP record and does not update individual fields.
func (s *Store) UpdateVASP(v *gds.VASP) (err error) {
	if v.Id == "" {
		return storeerrors.ErrIncompleteRecord
	}
	key := []byte(v.Id)

	// Ensure a common name exists for the uniqueness constraint
	if cn := index.Normalize(v.CommonName); cn == "" {
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
	// This must be inside the lock so that the database indices are consistent.
	// NOTE: the lock doesn't prevent concurrent writes from multiple GDS instances,
	// just from this instance ... we need trtl transactions to guarantee this (or
	// indices managed by trtl itself, since this is the only reason we're doing a Get).
	o, err := s.RetrieveVASP(v.Id)
	if err != nil {
		return err
	}

	// Update the VASP record
	// This must be inside the lock so that there is no race condition between the index
	// and the stored index inside of the database.
	ctx, cancel := withContext(context.Background())
	defer cancel()
	request := &pb.PutRequest{
		Key:       key,
		Value:     val,
		Namespace: wire.NamespaceVASPs,
	}
	if reply, err := s.client.Put(ctx, request); err != nil || !reply.Success {
		if err == nil {
			err = storeerrors.ErrProtocol
		}
		return err
	}

	// Update indices to match the new record, editing via remove old records and insert
	// new records to ensure that the index correctly reflects the state of the update
	// without having to determine exactly what changed.
	if err = s.removeIndices(o); err != nil {
		// NOTE: if this error is triggered, admins should reindex the database.
		log.Error().Err(err).Msg("could not remove previous indices on update: reindex required")
	}
	if err = s.insertIndices(v); err != nil {
		return err
	}
	return nil
}

// DeleteVASP record, removing it completely from the database and indices.
func (s *Store) DeleteVASP(id string) error {
	key := []byte(id)

	// Critical section (optimizing for safety rather than speed)
	s.Lock()
	defer s.Unlock()

	// Lookup the record in order to remove in order to remove data from the indices.
	// This must be inside the lock to ensure indices are updated correctly with what
	// is on disk. However, this doesn't prevent a concurrency issue with another
	// GDS replica interacting with trtl.
	o, err := s.RetrieveVASP(id)
	if err != nil {
		if err == storeerrors.ErrEntityNotFound {
			return nil
		}
		return err
	}

	ctx, cancel := withContext(context.Background())
	defer cancel()
	request := &pb.DeleteRequest{
		Key:       key,
		Namespace: wire.NamespaceVASPs,
	}
	if reply, err := s.client.Delete(ctx, request); err != nil || !reply.Success {
		if err == nil {
			err = storeerrors.ErrProtocol
		}
		return err
	}

	// Remove the deleted entity records from the indices
	if err = s.removeIndices(o); err != nil {
		return err
	}
	return nil
}

//===========================================================================
// CertificateStore Implementation
//===========================================================================

// ListCertReqs returns all certificate requests that are currently in the store.
func (s *Store) ListCertReqs() iterator.CertificateIterator {
	return &certReqIterator{
		NewTrtlStreamingIterator(s.client, wire.NamespaceCertReqs),
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

	key := []byte(r.Id)

	var data []byte
	if data, err = proto.Marshal(r); err != nil {
		return "", err
	}

	ctx, cancel := withContext(context.Background())
	defer cancel()
	request := &pb.PutRequest{
		Key:       key,
		Value:     data,
		Namespace: wire.NamespaceCertReqs,
	}
	if reply, err := s.client.Put(ctx, request); err != nil || !reply.Success {
		if err == nil {
			err = storeerrors.ErrProtocol
		}
		return "", err
	}

	return r.Id, nil
}

// RetrieveCertReq returns a certificate request by certificate request ID.
func (s *Store) RetrieveCertReq(id string) (r *models.CertificateRequest, err error) {
	if id == "" {
		return nil, storeerrors.ErrEntityNotFound
	}

	ctx, cancel := withContext(context.Background())
	defer cancel()
	request := &pb.GetRequest{
		Key:       []byte(id),
		Namespace: wire.NamespaceCertReqs,
	}
	var reply *pb.GetReply
	if reply, err = s.client.Get(ctx, request); err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, storeerrors.ErrEntityNotFound
		}
		return nil, err
	}

	r = new(models.CertificateRequest)
	if err = proto.Unmarshal(reply.Value, r); err != nil {
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
	key := []byte(r.Id)
	if data, err = proto.Marshal(r); err != nil {
		return err
	}

	ctx, cancel := withContext(context.Background())
	defer cancel()
	request := &pb.PutRequest{
		Key:       key,
		Value:     data,
		Namespace: wire.NamespaceCertReqs,
	}
	if reply, err := s.client.Put(ctx, request); err != nil || !reply.Success {
		if err == nil {
			err = storeerrors.ErrProtocol
		}
		return err
	}
	return nil
}

// DeleteCertReq removes a certificate request from the store.
func (s *Store) DeleteCertReq(id string) (err error) {
	ctx, cancel := withContext(context.Background())
	defer cancel()
	request := &pb.DeleteRequest{
		Key:       []byte(id),
		Namespace: wire.NamespaceCertReqs,
	}
	if reply, err := s.client.Delete(ctx, request); err != nil || !reply.Success {
		if err == nil {
			err = storeerrors.ErrProtocol
		}
		return err
	}
	return nil
}
