package trtl

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/trisacrypto/directory/pkg/gds/client"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	storeerrors "github.com/trisacrypto/directory/pkg/gds/store/errors"
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
	return store, nil
}

// Store implements the store.Store interface for the Trtl replicated database.
type Store struct {
	conn   *grpc.ClientConn
	client pb.TrtlClient
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
	return s.conn.Close()
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

func (s *Store) SearchVASPs(query map[string]interface{}) ([]*gds.VASP, error) {
	// TODO: We need indexes for this.
	return nil, errors.New("not implemented")
}

// CreateVASP into the directory. This method requires the VASP to have a unique
// name and ignores any ID fields that are set on the VASP, instead assigning new IDs.
func (s *Store) CreateVASP(v *gds.VASP) (id string, err error) {
	// Create UUID for record
	if v.Id == "" {
		v.Id = uuid.New().String()
	}
	key := []byte(v.Id)

	// Update management timestamps and record metadata
	v.LastUpdated = time.Now().Format(time.RFC3339)
	if v.FirstListed == "" {
		v.FirstListed = v.LastUpdated
	}
	if v.Version == nil || v.Version.Version == 0 {
		v.Version = &gds.Version{Version: 1}
	}

	// TODO: Update the names index and enforce uniqueness constraint.

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

	// TODO: Check the uniqueness constraint.

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

	// TODO: Update the indices.

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
	return nil
}

// DeleteVASP record, removing it completely from the database and indices.
func (s *Store) DeleteVASP(id string) error {
	key := []byte(id)

	// TODO: Update the indices.

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
