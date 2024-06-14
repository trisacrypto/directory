package trtl

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	engine "github.com/rotationalio/honu/engines"
	"github.com/rs/zerolog/log"
	bff "github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/store/config"
	storeerrors "github.com/trisacrypto/directory/pkg/store/errors"
	"github.com/trisacrypto/directory/pkg/store/index"
	"github.com/trisacrypto/directory/pkg/store/iterator"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/trisacrypto/directory/pkg/utils"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Open a connection to the Trtl database.
func Open(conf config.StoreConfig) (store *Store, err error) {
	store = &Store{}
	if store.conn, err = Connect(conf); err != nil {
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
func (s *Store) ListVASPs(ctx context.Context) iterator.DirectoryIterator {
	return &vaspIterator{
		NewTrtlStreamingIterator(s.client, wire.NamespaceVASPs),
	}
}

// SearchVASPs is intended to specifically identify a VASP (rather than as a browsing
// functionality). As such it is primarily a filtering search rather than an inclusive
// search. The query can contain a one or more name or website terms. Names are prefixed
// matched to the index and websites are hostname matched. The query can contain one or
// more country and category filters as well, which reduce the number of search results.
func (s *Store) SearchVASPs(ctx context.Context, query map[string]interface{}) (vasps []*gds.VASP, err error) {
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
			if vasp, err = s.RetrieveVASP(ctx, id); err != nil {
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
func (s *Store) CreateVASP(ctx context.Context, v *gds.VASP) (id string, err error) {
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

	// Check the uniqueness constraints
	// NOTE: website removed as uniqueness constraint in SC-4483
	if _, ok := s.names.Find(v.CommonName); ok {
		return "", storeerrors.ErrDuplicateEntity
	}

	var data []byte
	if data, err = proto.Marshal(v); err != nil {
		return "", err
	}

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
func (s *Store) RetrieveVASP(ctx context.Context, id string) (v *gds.VASP, err error) {
	key := []byte(id)

	ctx, cancel := utils.WithDeadline(ctx)
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
func (s *Store) UpdateVASP(ctx context.Context, v *gds.VASP) (err error) {
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
	o, err := s.RetrieveVASP(ctx, v.Id)
	if err != nil {
		return err
	}

	// Check the uniqueness constraints
	// NOTE: website removed as uniqueness constraint in SC-4483
	if id, ok := s.names.Find(v.CommonName); ok && id != v.Id {
		return storeerrors.ErrDuplicateEntity
	}

	// Update the VASP record
	// This must be inside the lock so that there is no race condition between the index
	// and the stored index inside of the database.
	ctx, cancel := utils.WithDeadline(ctx)
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
func (s *Store) DeleteVASP(ctx context.Context, id string) error {
	key := []byte(id)

	// Critical section (optimizing for safety rather than speed)
	s.Lock()
	defer s.Unlock()

	// Lookup the record in order to remove in order to remove data from the indices.
	// This must be inside the lock to ensure indices are updated correctly with what
	// is on disk. However, this doesn't prevent a concurrency issue with another
	// GDS replica interacting with trtl.
	o, err := s.RetrieveVASP(ctx, id)
	if err != nil {
		if err == storeerrors.ErrEntityNotFound {
			return nil
		}
		return err
	}

	ctx, cancel := utils.WithDeadline(ctx)
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

// ListCerts returns all certificates that are currently in the store.
func (s *Store) ListCerts(ctx context.Context) iterator.CertificateIterator {
	return &certIterator{
		NewTrtlStreamingIterator(s.client, wire.NamespaceCerts),
	}
}

// CreateCert and assign a new ID and return the version.
func (s *Store) CreateCert(ctx context.Context, c *models.Certificate) (id string, err error) {
	if c.Id != "" {
		return "", storeerrors.ErrIDAlreadySet
	}

	// Create UUID for record
	// TODO: check uniqueness of the ID
	c.Id = uuid.New().String()
	key := []byte(c.Id)

	var data []byte
	if data, err = proto.Marshal(c); err != nil {
		return "", err
	}

	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()
	request := &pb.PutRequest{
		Key:       key,
		Value:     data,
		Namespace: wire.NamespaceCerts,
	}
	if reply, err := s.client.Put(ctx, request); err != nil || !reply.Success {
		if err == nil {
			err = storeerrors.ErrProtocol
		}
		return "", err
	}

	return c.Id, nil
}

// RetrieveCert returns a certificate by certificate ID.
func (s *Store) RetrieveCert(ctx context.Context, id string) (c *models.Certificate, err error) {
	if id == "" {
		return nil, storeerrors.ErrEntityNotFound
	}

	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()
	request := &pb.GetRequest{
		Key:       []byte(id),
		Namespace: wire.NamespaceCerts,
	}
	var reply *pb.GetReply
	if reply, err = s.client.Get(ctx, request); err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, storeerrors.ErrEntityNotFound
		}
		return nil, err
	}

	c = new(models.Certificate)
	if err = proto.Unmarshal(reply.Value, c); err != nil {
		return nil, err
	}

	return c, nil
}

// UpdateCert can create or update a certificate. The certificate should be as
// complete as possible, including an ID generated by the caller.
func (s *Store) UpdateCert(ctx context.Context, c *models.Certificate) (err error) {
	if c.Id == "" {
		return storeerrors.ErrIncompleteRecord
	}

	var data []byte
	key := []byte(c.Id)
	if data, err = proto.Marshal(c); err != nil {
		return err
	}

	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()
	request := &pb.PutRequest{
		Key:       key,
		Value:     data,
		Namespace: wire.NamespaceCerts,
	}
	if reply, err := s.client.Put(ctx, request); err != nil || !reply.Success {
		if err == nil {
			err = storeerrors.ErrProtocol
		}
		return err
	}
	return nil
}

// DeleteCert removes a certificate from the store.
func (s *Store) DeleteCert(ctx context.Context, id string) (err error) {
	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()
	request := &pb.DeleteRequest{
		Key:       []byte(id),
		Namespace: wire.NamespaceCerts,
	}
	if reply, err := s.client.Delete(ctx, request); err != nil || !reply.Success {
		if err == nil {
			err = storeerrors.ErrProtocol
		}
		return err
	}
	return nil
}

//===========================================================================
// CertificateRequestStore Implementation
//===========================================================================

// ListCertReqs returns all certificate requests that are currently in the store.
func (s *Store) ListCertReqs(ctx context.Context) iterator.CertificateRequestIterator {
	return &certReqIterator{
		NewTrtlStreamingIterator(s.client, wire.NamespaceCertReqs),
	}
}

// CreateCertReq and assign a new ID and return the version.
func (s *Store) CreateCertReq(ctx context.Context, r *models.CertificateRequest) (id string, err error) {
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

	ctx, cancel := utils.WithDeadline(ctx)
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
func (s *Store) RetrieveCertReq(ctx context.Context, id string) (r *models.CertificateRequest, err error) {
	if id == "" {
		return nil, storeerrors.ErrEntityNotFound
	}

	ctx, cancel := utils.WithDeadline(ctx)
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
func (s *Store) UpdateCertReq(ctx context.Context, r *models.CertificateRequest) (err error) {
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

	ctx, cancel := utils.WithDeadline(ctx)
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
func (s *Store) DeleteCertReq(ctx context.Context, id string) (err error) {
	ctx, cancel := utils.WithDeadline(ctx)
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

//===========================================================================
// AnnouncementStore Implementation
//===========================================================================

// RetrieveAnnouncementMonth returns the announcement month "crate" for the given month
// timestamp in the format YYYY-MM.
func (s *Store) RetrieveAnnouncementMonth(ctx context.Context, date string) (m *bff.AnnouncementMonth, err error) {
	if date == "" {
		return nil, storeerrors.ErrEntityNotFound
	}

	// Get the key by creating an intermediate announcement month to ensure that
	// validation and key creation always happens the same way.
	var key []byte
	m = &bff.AnnouncementMonth{Date: date}
	if key, err = m.Key(); err != nil {
		return nil, err
	}

	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()
	request := &pb.GetRequest{
		Key:       key,
		Namespace: wire.NamespaceAnnouncements,
	}
	var reply *pb.GetReply
	if reply, err = s.client.Get(ctx, request); err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, storeerrors.ErrEntityNotFound
		}
		return nil, err
	}

	if err = proto.Unmarshal(reply.Value, m); err != nil {
		return nil, err
	}
	return m, nil
}

// UpdateAnnouncementMonth creates a new announcement month "crate" if it doesn't
// already exist or replaces the existing record.
func (s *Store) UpdateAnnouncementMonth(ctx context.Context, m *bff.AnnouncementMonth) (err error) {
	if m.Date == "" {
		return storeerrors.ErrIncompleteRecord
	}

	// Get the key by creating an intermediate announcement month to ensure that
	// validation and key creation always happens the same way.
	var key []byte
	if key, err = m.Key(); err != nil {
		return err
	}

	// Update the modified timestamp
	m.Modified = time.Now().Format(time.RFC3339Nano)
	if m.Created == "" {
		m.Created = m.Modified
	}

	var data []byte
	if data, err = proto.Marshal(m); err != nil {
		return err
	}

	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()
	request := &pb.PutRequest{
		Key:       key,
		Value:     data,
		Namespace: wire.NamespaceAnnouncements,
	}
	if reply, err := s.client.Put(ctx, request); err != nil || !reply.Success {
		if err == nil {
			err = storeerrors.ErrProtocol
		}
		return err
	}
	return nil
}

// DeleteAnnouncementMonth removes an announcement month "crate" from the store.
func (s *Store) DeleteAnnouncementMonth(ctx context.Context, date string) (err error) {
	// Get the key by creating an intermediate announcement month to ensure that
	// validation and key creation always happens the same way.
	var key []byte
	m := &bff.AnnouncementMonth{Date: date}
	if key, err = m.Key(); err != nil {
		return err
	}

	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()
	request := &pb.DeleteRequest{
		Key:       key,
		Namespace: wire.NamespaceAnnouncements,
	}
	if reply, err := s.client.Delete(ctx, request); err != nil || !reply.Success {
		if err == nil {
			err = storeerrors.ErrProtocol
		}
		return err
	}
	return nil
}

//===========================================================================
// ActivityStore Implementation
//===========================================================================

// RetrieveActivityMonth returns the activity month record for the given date
// timestamp in the format YYYY-MM.
func (s *Store) RetrieveActivityMonth(ctx context.Context, date string) (m *bff.ActivityMonth, err error) {
	if date == "" {
		return nil, storeerrors.ErrEntityNotFound
	}

	// Get the key by creating an intermediate activity month to ensure that
	// validation and key creation always happens the same way.
	m = &bff.ActivityMonth{Date: date}
	key, err := m.Key()
	if err != nil {
		return nil, err
	}

	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()
	request := &pb.GetRequest{
		Key:       key,
		Namespace: wire.NamespaceActivities,
	}
	var reply *pb.GetReply
	if reply, err = s.client.Get(ctx, request); err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, storeerrors.ErrEntityNotFound
		}
		return nil, err
	}

	if err = proto.Unmarshal(reply.Value, m); err != nil {
		return nil, err
	}
	return m, nil
}

// UpdateActivityMonth creates a new activity month record if it doesn't already
// exist or replaces the existing record.
func (s *Store) UpdateActivityMonth(ctx context.Context, m *bff.ActivityMonth) (err error) {
	if m.Date == "" {
		return storeerrors.ErrIncompleteRecord
	}

	// Get the key by creating an intermediate activity month to ensure that
	// validation and key creation always happens the same way.
	var key []byte
	if key, err = m.Key(); err != nil {
		return err
	}

	// Update the modified timestamp
	m.Modified = time.Now().Format(time.RFC3339Nano)
	if m.Created == "" {
		m.Created = m.Modified
	}

	var data []byte
	if data, err = proto.Marshal(m); err != nil {
		return err
	}

	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()
	request := &pb.PutRequest{
		Key:       key,
		Value:     data,
		Namespace: wire.NamespaceActivities,
	}
	if reply, err := s.client.Put(ctx, request); err != nil || !reply.Success {
		if err == nil {
			err = storeerrors.ErrProtocol
		}
		return err
	}
	return nil
}

// DeleteActivityMonth removes an activity month record from the store.
func (s *Store) DeleteActivityMonth(ctx context.Context, date string) (err error) {
	// Get the key by creating an intermediate activity month to ensure that
	// validation and key creation always happens the same way.
	var key []byte
	m := &bff.ActivityMonth{Date: date}
	if key, err = m.Key(); err != nil {
		return err
	}

	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()
	request := &pb.DeleteRequest{
		Key:       key,
		Namespace: wire.NamespaceActivities,
	}
	if reply, err := s.client.Delete(ctx, request); err != nil || !reply.Success {
		if err == nil {
			err = storeerrors.ErrProtocol
		}
		return err
	}
	return nil
}

//===========================================================================
// OrganizationStore Implementation
//===========================================================================

// ListOrganizations returns an iterator to retrieve all organizations.
func (s *Store) ListOrganizations(ctx context.Context) iterator.OrganizationIterator {
	return &organizationIterator{
		NewTrtlStreamingIterator(s.client, wire.NamespaceOrganizations),
	}
}

// CreateOrganization creates a new organization record in the store, assigning a
// unique ID if it doesn't exist and setting the created and modified timestamps.
func (s *Store) CreateOrganization(ctx context.Context, o *bff.Organization) (id string, err error) {
	// Create the organization ID if not provided
	if o.Id == "" {
		o.Id = uuid.New().String()
	}

	// Set the created and modified timestamps
	ts := time.Now().Format(time.RFC3339Nano)
	o.Created = ts
	o.Modified = ts

	var data []byte
	if data, err = proto.Marshal(o); err != nil {
		return "", err
	}

	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()
	request := &pb.PutRequest{
		Key:       o.Key(),
		Value:     data,
		Namespace: wire.NamespaceOrganizations,
	}
	if reply, err := s.client.Put(ctx, request); err != nil || !reply.Success {
		if err == nil {
			err = storeerrors.ErrProtocol
		}
		return "", err
	}
	return o.Id, nil
}

// RetrieveOrganization retrieves an organization record from the store by UUID.
func (s *Store) RetrieveOrganization(ctx context.Context, id uuid.UUID) (o *bff.Organization, err error) {
	if id == uuid.Nil {
		return nil, storeerrors.ErrEntityNotFound
	}

	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()
	request := &pb.GetRequest{
		Key:       id[:],
		Namespace: wire.NamespaceOrganizations,
	}
	var reply *pb.GetReply
	if reply, err = s.client.Get(ctx, request); err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, storeerrors.ErrEntityNotFound
		}
		return nil, err
	}

	o = new(bff.Organization)
	if err = proto.Unmarshal(reply.Value, o); err != nil {
		return nil, err
	}
	return o, nil
}

// UpdateOrganization updates an organization record in the store by replacing the
// existing record.
func (s *Store) UpdateOrganization(ctx context.Context, o *bff.Organization) (err error) {
	if o.Id == "" {
		return storeerrors.ErrEntityNotFound
	}

	// Update the modified timestamp
	o.Modified = time.Now().Format(time.RFC3339Nano)
	if o.Created == "" {
		o.Created = o.Modified
	}

	var data []byte
	if data, err = proto.Marshal(o); err != nil {
		return err
	}

	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()
	request := &pb.PutRequest{
		Key:       o.Key(),
		Value:     data,
		Namespace: wire.NamespaceOrganizations,
	}
	if reply, err := s.client.Put(ctx, request); err != nil || !reply.Success {
		if err == nil {
			err = storeerrors.ErrProtocol
		}
		return err
	}
	return nil
}

// DeleteOrganization deletes an organization record from the store by UUID.
func (s *Store) DeleteOrganization(ctx context.Context, id uuid.UUID) (err error) {
	if id == uuid.Nil {
		return storeerrors.ErrEntityNotFound
	}

	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()
	request := &pb.DeleteRequest{
		Key:       id[:],
		Namespace: wire.NamespaceOrganizations,
	}
	if reply, err := s.client.Delete(ctx, request); err != nil || !reply.Success {
		if err == nil {
			err = storeerrors.ErrProtocol
		}
		return err
	}
	return nil
}

// TODO: Delete the trtl ContactStore implementation
//===========================================================================
// ContactStore Implementation
//===========================================================================

func (s *Store) ListContacts(ctx context.Context) []*models.Contact {
	return nil
}

// CreateContact creates a new Contact record in the store, using the contact's
// email as a unique ID.
func (s *Store) CreateContact(ctx context.Context, c *models.Contact) (_ string, err error) {
	if c == nil || c.Email == "" {
		return "", storeerrors.ErrIncompleteRecord
	}

	// Update management timestamps and record metadata
	c.Created = time.Now().Format(time.RFC3339)
	c.Modified = c.Created

	// Marshal the Contact
	var data []byte
	if data, err = proto.Marshal(c); err != nil {
		return "", err
	}

	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()

	// TODO: determine the best way to ensure uniqueness of the key
	// Create and store the PutRequest
	key := []byte(models.NormalizeEmail(c.Email))
	request := &pb.PutRequest{
		Key:       key,
		Value:     data,
		Namespace: wire.NamespaceContacts,
	}
	if reply, err := s.client.Put(ctx, request); err != nil || !reply.Success {
		if err == nil {
			err = storeerrors.ErrProtocol
		}
		return "", err
	}

	return c.Email, nil
}

// RetrieveContact returns a contact request by contact email.
func (s *Store) RetrieveContact(ctx context.Context, email string) (c *models.Contact, err error) {
	if email == "" {
		return nil, storeerrors.ErrEntityNotFound
	}

	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()

	key := []byte(models.NormalizeEmail(email))
	request := &pb.GetRequest{
		Key:       key,
		Namespace: wire.NamespaceContacts,
	}
	var reply *pb.GetReply
	if reply, err = s.client.Get(ctx, request); err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, storeerrors.ErrEntityNotFound
		}
		return nil, err
	}

	c = new(models.Contact)
	if err = proto.Unmarshal(reply.Value, c); err != nil {
		return nil, err
	}

	return c, nil
}

// UpdateContact can create or update a contact request. The request should be as
// complete as possible, including an email provided by the caller.
func (s *Store) UpdateContact(ctx context.Context, c *models.Contact) (err error) {
	if c == nil || c.Email == "" {
		return storeerrors.ErrIncompleteRecord
	}

	var data []byte
	c.Modified = time.Now().Format(time.RFC3339)
	if data, err = proto.Marshal(c); err != nil {
		return err
	}

	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()

	key := []byte(models.NormalizeEmail(c.Email))
	request := &pb.PutRequest{
		Key:       key,
		Value:     data,
		Namespace: wire.NamespaceContacts,
	}
	if reply, err := s.client.Put(ctx, request); err != nil || !reply.Success {
		if err == nil {
			err = storeerrors.ErrProtocol
		}
		return err
	}
	return nil
}

// DeleteContact deletes a contact record from the store by email.
func (s *Store) DeleteContact(ctx context.Context, email string) error {
	if email == "" {
		return storeerrors.ErrEntityNotFound
	}

	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()

	key := []byte(models.NormalizeEmail(email))
	request := &pb.DeleteRequest{
		Key:       key,
		Namespace: wire.NamespaceContacts,
	}
	if reply, err := s.client.Delete(ctx, request); err != nil || !reply.Success {
		if err == nil {
			err = storeerrors.ErrProtocol
		}
		return err
	}
	return nil
}

//===========================================================================
// EmailStore Implementation
//===========================================================================

// List all of the emails in the database.
func (s *Store) ListEmails(ctx context.Context) iterator.EmailIterator {
	return &emailIterator{
		NewTrtlStreamingIterator(s.client, wire.NamespaceEmails),
	}
}

// CreateEmail creates a new Email record in the store using the normalized unique email
// as the key. If the email already exists in the database, an error is returned.
func (s *Store) CreateEmail(ctx context.Context, c *models.Email) (_ string, err error) {
	if c == nil || c.Email == "" {
		return "", storeerrors.ErrIncompleteRecord
	}

	// Validate the email model
	if err = c.Validate(); err != nil {
		return "", err
	}

	// Update management timestamps and record metadata
	c.Created = time.Now().Format(time.RFC3339)
	c.Modified = c.Created

	// Create the Put request to save the email in trtl; create email requires that the
	// email doesn't exist, otherwise an AlreadyExists error is returned. This is
	// handled by the Honu transaction exists invariant.
	request := &pb.PutRequest{
		Key:       []byte(models.NormalizeEmail(c.Email)),
		Namespace: wire.NamespaceEmails,
		Options: &pb.Options{
			RequireNotExists: true,
		},
	}

	// Serialize the protocol buffer into the request
	if request.Value, err = proto.Marshal(c); err != nil {
		return "", err
	}

	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()

	// Execute the Put request to save the data in trtl
	if reply, err := s.client.Put(ctx, request); err != nil || !reply.Success {
		if err == nil {
			err = storeerrors.ErrProtocol
		}

		if serr, ok := status.FromError(err); ok {
			// Unfortunately there isn't a better way to check if the error is already
			// exists without string matching.
			if serr.Code() == codes.FailedPrecondition && serr.Message() == engine.ErrAlreadyExists.Error() {
				return "", storeerrors.ErrEmailExists
			}
		}
		return "", err
	}
	return c.Email, nil
}

// RetrieveContact returns a contact request by contact email.
func (s *Store) RetrieveEmail(ctx context.Context, email string) (c *models.Email, err error) {
	if strings.TrimSpace(email) == "" {
		return nil, storeerrors.ErrEntityNotFound
	}

	request := &pb.GetRequest{
		Key:       []byte(models.NormalizeEmail(email)),
		Namespace: wire.NamespaceEmails,
	}

	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()

	var reply *pb.GetReply
	if reply, err = s.client.Get(ctx, request); err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, storeerrors.ErrEntityNotFound
		}
		return nil, err
	}

	c = &models.Email{}
	if err = proto.Unmarshal(reply.Value, c); err != nil {
		return nil, err
	}
	return c, nil
}

// UpdateContact can create or update a contact request. The request should be as
// complete as possible, including an email provided by the caller.
func (s *Store) UpdateEmail(ctx context.Context, c *models.Email) (err error) {
	if c == nil || c.Email == "" {
		return storeerrors.ErrIncompleteRecord
	}

	// Validate the email model
	if err = c.Validate(); err != nil {
		return err
	}

	// Manage the updated and modified timestamps
	c.Modified = time.Now().Format(time.RFC3339)
	if c.Created == "" {
		c.Created = c.Modified
	}

	// RequireNotExists should be false for update semantics. We're also ensuring that
	// RequireExists is false so that this method can be used for create without an
	// existence check to match the leveldb store semantics.
	request := &pb.PutRequest{
		Key:       []byte(models.NormalizeEmail(c.Email)),
		Namespace: wire.NamespaceEmails,
		Options: &pb.Options{
			RequireExists:    false,
			RequireNotExists: false,
		},
	}

	// Serialize the protocol buffer into the request
	if request.Value, err = proto.Marshal(c); err != nil {
		return err
	}

	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()

	if reply, err := s.client.Put(ctx, request); err != nil || !reply.Success {
		if err == nil {
			return storeerrors.ErrProtocol
		}
		return err
	}

	return nil
}

// DeleteContact deletes an contact record from the store by email.
func (s *Store) DeleteEmail(ctx context.Context, email string) (err error) {
	if email == "" {
		return storeerrors.ErrEntityNotFound
	}

	request := &pb.DeleteRequest{
		Key:       []byte(models.NormalizeEmail(email)),
		Namespace: wire.NamespaceEmails,
	}

	ctx, cancel := utils.WithDeadline(ctx)
	defer cancel()

	var reply *pb.DeleteReply
	if reply, err = s.client.Delete(ctx, request); err != nil || !reply.Success {
		if err == nil {
			return storeerrors.ErrProtocol
		}
		return err
	}
	return nil
}

//===========================================================================
// DirectoryStore Implementation
//===========================================================================

// VASPContacts implements a join mechanism to ensure that the contacts on the VASP
// (e.g. the technical, administrative, legal, and billing contacts) are connected with
// the email address records for those contacts.
func (s *Store) VASPContacts(ctx context.Context, vasp *gds.VASP) (_ *models.Contacts, err error) {
	// Identify all the normalized, unique emails that need to be retrieved.
	emails := make(map[string]struct{})
	vcards := vasp.Contacts

	if vcards.Administrative != nil && vcards.Administrative.Email != "" {
		emails[models.NormalizeEmail(vcards.Administrative.Email)] = struct{}{}
	}

	if vcards.Technical != nil && vcards.Technical.Email != "" {
		emails[models.NormalizeEmail(vcards.Technical.Email)] = struct{}{}
	}

	if vcards.Legal != nil && vcards.Legal.Email != "" {
		emails[models.NormalizeEmail(vcards.Legal.Email)] = struct{}{}
	}

	if vcards.Billing != nil && vcards.Billing.Email != "" {
		emails[models.NormalizeEmail(vcards.Billing.Email)] = struct{}{}
	}

	// Create the contacts record to return.
	contacts := &models.Contacts{
		VASP:     vasp.Id,
		Contacts: vcards,
		Emails:   make([]*models.Email, 0, len(emails)),
	}

	// Fetch the emails and add them to the contacts.
	// TODO: rather than doing a retrieve request for each email, batch the request.
	for email := range emails {
		var record *models.Email
		if record, err = s.RetrieveEmail(ctx, email); err != nil {
			return nil, err
		}
		contacts.Emails = append(contacts.Emails, record)
	}
	return contacts, nil
}

// This is a helper method to fetch the contacts for a VASP with only the vaspID.
func (s *Store) RetrieveVASPContacts(ctx context.Context, vaspID string) (_ *models.Contacts, err error) {
	// TODO: batch all requests rather than sending indvidual requests.
	var vasp *gds.VASP
	if vasp, err = s.RetrieveVASP(ctx, vaspID); err != nil {
		return nil, err
	}
	return s.VASPContacts(ctx, vasp)
}

func (s *Store) UpdateVASPContacts(ctx context.Context, vaspID string, contacts *models.Contacts) error {
	return errors.New("not implemented yet")
}
