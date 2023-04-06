package trtl

import (
	"context"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
	storeerrors "github.com/trisacrypto/directory/pkg/store/errors"
	"github.com/trisacrypto/directory/pkg/store/index"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/trisacrypto/directory/pkg/utils"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

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

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	cursor, err := s.client.Cursor(ctx, &pb.CursorRequest{Namespace: wire.NamespaceVASPs})
	if err != nil {
		return err
	}

	for {
		var pair *pb.KVPair
		if pair, err = cursor.Recv(); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		vasp := &gds.VASP{}
		if err = proto.Unmarshal(pair.Value, vasp); err != nil {
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

	if err = cursor.CloseSend(); err != nil {
		log.Error().Err(err).Msg("could not close the trtl cursor in reindex")
	}

	// Critical section
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
// Indices and Synchronization
//===========================================================================

func (s *Store) insertIndices(v *gds.VASP) error {
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

func (s *Store) removeIndices(v *gds.VASP) error {
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

// keys and prefixes for leveldb buckets and indices
var (
	keyNameIndex     = []byte("names")
	keyWebsiteIndex  = []byte("websites")
	keyCountryIndex  = []byte("countries")
	keyCategoryIndex = []byte("categories")
)

// Sync exposes the index synchronization functionality to tests, allowing them to sync
// a single index or all indices all at once.
func (s *Store) Sync(index string) error {
	switch index {
	case "name", "names":
		return s.syncnames()
	case "website", "websites":
		return s.syncwebsites()
	case "country", "countries":
		return s.synccountries()
	case "category", "categories":
		return s.synccategories()
	case "", "all":
		return s.sync()
	default:
		return fmt.Errorf(`unknown index %q, use empty string or "all" to sync all indices`, index)
	}
}

func (s *Store) sync() (err error) {
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

func (s *Store) syncnames() (err error) {
	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Critical section (optimizing for safety rather than speed)
	s.Lock()
	defer s.Unlock()

	if s.names == nil {
		// Create the index to load it from disk
		s.names = index.NewNamesIndex()

		// Fetch the data from the database
		var rep *pb.GetReply
		if rep, err = s.client.Get(ctx, &pb.GetRequest{Key: keyNameIndex, Namespace: wire.NamespaceIndices}); err != nil {
			if status.Code(err) == codes.NotFound {
				return nil
			}
			log.Error().Err(err).Msg("could not fetch names index from database")
			return err
		}

		if err = s.names.Load(rep.Value); err != nil {
			log.Error().Err(err).Msg("could not unmarshal names index")
			return storeerrors.ErrCorruptedIndex
		}
	}

	// Put the current names back to the database
	if !s.names.Empty() {
		var value []byte
		if value, err = s.names.Dump(); err != nil {
			log.Error().Err(err).Msg("could not marshal names index")
			return storeerrors.ErrCorruptedIndex
		}

		if rep, err := s.client.Put(ctx, &pb.PutRequest{Key: keyNameIndex, Value: value, Namespace: wire.NamespaceIndices}); err != nil || !rep.Success {
			if err == nil {
				err = storeerrors.ErrProtocol
			}
			log.Error().Err(err).Msg("could not put names index")
			return storeerrors.ErrCorruptedIndex
		}

		log.Debug().Int("size", len(value)).Msg("names index checkpointed")
	}
	return nil
}

func (s *Store) syncwebsites() (err error) {
	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Critical section (optimizing for safety rather than speed)
	s.Lock()
	defer s.Unlock()

	if s.websites == nil {
		// Create the index to load it from disk
		s.websites = index.NewWebsiteIndex()

		// Fetch the data from the database
		var rep *pb.GetReply
		if rep, err = s.client.Get(ctx, &pb.GetRequest{Key: keyWebsiteIndex, Namespace: wire.NamespaceIndices}); err != nil {
			if status.Code(err) == codes.NotFound {
				return nil
			}
			log.Error().Err(err).Msg("could not fetch websites index from database")
			return err
		}

		if err = s.websites.Load(rep.Value); err != nil {
			log.Error().Err(err).Msg("could not unmarshal websites index")
			return storeerrors.ErrCorruptedIndex
		}
	}

	// Put the current websites back to the database
	if !s.websites.Empty() {
		var value []byte
		if value, err = s.websites.Dump(); err != nil {
			log.Error().Err(err).Msg("could not marshal websites index")
			return storeerrors.ErrCorruptedIndex
		}

		if rep, err := s.client.Put(ctx, &pb.PutRequest{Key: keyWebsiteIndex, Value: value, Namespace: wire.NamespaceIndices}); err != nil || !rep.Success {
			if err == nil {
				err = storeerrors.ErrProtocol
			}
			log.Error().Err(err).Msg("could not put websites index")
			return storeerrors.ErrCorruptedIndex
		}

		log.Debug().Int("size", len(value)).Msg("websites index checkpointed")
	}
	return nil
}

func (s *Store) synccountries() (err error) {
	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Critical section (optimizing for safety rather than speed)
	s.Lock()
	defer s.Unlock()

	if s.countries == nil {
		// Create the index to load it from disk
		s.countries = index.NewCountryIndex()

		// Fetch the data from the database
		var rep *pb.GetReply
		if rep, err = s.client.Get(ctx, &pb.GetRequest{Key: keyCountryIndex, Namespace: wire.NamespaceIndices}); err != nil {
			if status.Code(err) == codes.NotFound {
				return nil
			}
			log.Error().Err(err).Msg("could not fetch countries index from database")
			return err
		}

		if err = s.countries.Load(rep.Value); err != nil {
			log.Error().Err(err).Msg("could not unmarshal countries index")
			return storeerrors.ErrCorruptedIndex
		}
	}

	// Put the current countries back to the database
	if !s.countries.Empty() {
		var value []byte
		if value, err = s.countries.Dump(); err != nil {
			log.Error().Err(err).Msg("could not marshal countries index")
			return storeerrors.ErrCorruptedIndex
		}

		if rep, err := s.client.Put(ctx, &pb.PutRequest{Key: keyCountryIndex, Value: value, Namespace: wire.NamespaceIndices}); err != nil || !rep.Success {
			if err == nil {
				err = storeerrors.ErrProtocol
			}
			log.Error().Err(err).Msg("could not put countries index")
			return storeerrors.ErrCorruptedIndex
		}

		log.Debug().Int("size", len(value)).Msg("countries index checkpointed")
	}
	return nil
}

func (s *Store) synccategories() (err error) {
	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Critical section (optimizing for safety rather than speed)
	s.Lock()
	defer s.Unlock()

	if s.categories == nil {
		// Create the index to load it from disk
		s.categories = index.NewCategoryIndex()

		// Fetch the data from the database
		var rep *pb.GetReply
		if rep, err = s.client.Get(ctx, &pb.GetRequest{Key: keyCategoryIndex, Namespace: wire.NamespaceIndices}); err != nil {
			if status.Code(err) == codes.NotFound {
				return nil
			}
			log.Error().Err(err).Msg("could not fetch categories index from database")
			return err
		}

		if err = s.categories.Load(rep.Value); err != nil {
			log.Error().Err(err).Msg("could not unmarshal categories index")
			return storeerrors.ErrCorruptedIndex
		}
	}

	// Put the current categories back to the database
	if !s.categories.Empty() {
		var value []byte
		if value, err = s.categories.Dump(); err != nil {
			log.Error().Err(err).Msg("could not marshal categories index")
			return storeerrors.ErrCorruptedIndex
		}

		if rep, err := s.client.Put(ctx, &pb.PutRequest{Key: keyCategoryIndex, Value: value, Namespace: wire.NamespaceIndices}); err != nil || !rep.Success {
			if err == nil {
				err = storeerrors.ErrProtocol
			}
			log.Error().Err(err).Msg("could not put categories index")
			return storeerrors.ErrCorruptedIndex
		}

		log.Debug().Int("size", len(value)).Msg("categories index checkpointed")
	}
	return nil
}

// GetNamesIndex for testing
func (s *Store) GetNamesIndex() index.SingleIndex {
	return s.names
}

// GetWebsitesIndex for testing
func (s *Store) GetWebsitesIndex() index.SingleIndex {
	return s.websites
}

// GetCountriesIndex for testing
func (s *Store) GetCountriesIndex() index.MultiIndex {
	return s.countries
}

// GetCategoriesIndex for testing
func (s *Store) GetCategoriesIndex() index.MultiIndex {
	return s.categories
}

// DeleteIndices for testing
// TODO: remove this function in favor of SC-3653
func (s *Store) DeleteIndices() (err error) {
	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	keys := [][]byte{keyNameIndex, keyWebsiteIndex, keyCategoryIndex, keyCountryIndex}
	for _, key := range keys {
		if _, err := s.client.Delete(ctx, &pb.DeleteRequest{Key: key, Namespace: wire.NamespaceIndices}); err != nil {
			log.Debug().Err(err).Msg("could not delete index")
		}
	}
	return nil
}
