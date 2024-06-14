package leveldb

import (
	"context"
	"errors"

	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func (s *Store) countPrefix(prefix []byte) (count uint64, err error) {
	iter := s.db.NewIterator(util.BytesPrefix(prefix), &opt.ReadOptions{DontFillCache: true})
	defer iter.Release()

	for iter.Next() {
		count++
	}
	if err = iter.Error(); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Store) CountVASPs(ctx context.Context) (uint64, error) {
	return s.countPrefix(preVASPs)
}

func (s *Store) CountCertReqs(context.Context) (uint64, error) {
	return s.countPrefix(preCertReqs)
}

func (s *Store) CountCerts(context.Context) (uint64, error) {
	return s.countPrefix(preCerts)
}

func (s *Store) CountAnnouncementMonths(context.Context) (uint64, error) {
	return 0, errors.New("cannot count announcement month in leveldb")
}

func (s *Store) CountActivityMonth(context.Context) (uint64, error) {
	return 0, errors.New("cannot count activity month in leveldb")
}

func (s *Store) CountOrganizations(context.Context) (uint64, error) {
	return s.countPrefix(preOrganizations)
}

func (s *Store) CountContacts(context.Context) (uint64, error) {
	return s.countPrefix(preContacts)
}

func (s *Store) CountEmails(context.Context) (uint64, error) {
	return s.countPrefix(preEmails)
}
