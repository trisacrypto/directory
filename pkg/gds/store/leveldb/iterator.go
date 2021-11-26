package leveldb

import (
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/store/wire"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/proto"
)

type iterWrapper struct {
	iter iterator.Iterator
}

type vaspIterator struct {
	iterWrapper
}

type certReqIterator struct {
	iterWrapper
}

func (i *iterWrapper) Next() bool {
	return i.iter.Next()
}

func (i *iterWrapper) Error() error {
	return i.iter.Error()
}

func (i *iterWrapper) Release() {
	i.iter.Release()
}

func (i *vaspIterator) VASP() *pb.VASP {
	vasp := new(pb.VASP)
	if err := proto.Unmarshal(i.iter.Value(), vasp); err != nil {
		log.Error().Err(err).Str("type", wire.NamespaceVASPs).Str("key", string(i.iter.Key())).Msg("corrupted data encountered")
		return nil
	}
	return vasp
}

func (i *vaspIterator) All() (vasps []*pb.VASP, err error) {
	vasps = make([]*pb.VASP, 0)
	defer i.iter.Release()

	for i.iter.Next() {
		vasp := new(pb.VASP)
		if err = proto.Unmarshal(i.iter.Value(), vasp); err != nil {
			return nil, err
		}
		vasps = append(vasps, vasp)
	}

	if err = i.iter.Error(); err != nil {
		return nil, err
	}
	return vasps, nil
}

func (i *certReqIterator) CertReq() *models.CertificateRequest {
	r := new(models.CertificateRequest)
	if err := proto.Unmarshal(i.iter.Value(), r); err != nil {
		log.Error().Err(err).Str("type", wire.NamespaceCertReqs).Str("key", string(i.iter.Key())).Msg("corrupted data encountered")
		return nil
	}
	return r
}

func (i *certReqIterator) All() (reqs []*models.CertificateRequest, err error) {
	reqs = make([]*models.CertificateRequest, 0)
	defer i.iter.Release()
	for i.iter.Next() {
		r := new(models.CertificateRequest)
		if err = proto.Unmarshal(i.iter.Value(), r); err != nil {
			return nil, err
		}
		reqs = append(reqs, r)
	}

	if err = i.iter.Error(); err != nil {
		return nil, err
	}

	return reqs, nil
}
