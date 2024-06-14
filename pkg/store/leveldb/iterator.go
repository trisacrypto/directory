package leveldb

import (
	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	bff "github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/proto"
)

type iterWrapper struct {
	iter iterator.Iterator
}

type vaspIterator struct {
	iterWrapper
}

type certIterator struct {
	iterWrapper
}

type certReqIterator struct {
	iterWrapper
}

type organizationIterator struct {
	iterWrapper
}

type emailIterator struct {
	iterWrapper
}

func (i *iterWrapper) Next() bool {
	return i.iter.Next()
}

func (i *iterWrapper) Prev() bool {
	return i.iter.Prev()
}

func (i *iterWrapper) Error() error {
	return i.iter.Error()
}

func (i *iterWrapper) Release() {
	i.iter.Release()
}

func (i *vaspIterator) VASP() (*pb.VASP, error) {
	vasp := new(pb.VASP)
	if err := proto.Unmarshal(i.iter.Value(), vasp); err != nil {
		sentry.Error(nil).Err(err).Str("type", wire.NamespaceVASPs).Str("key", string(i.iter.Key())).Msg("corrupted data encountered")
		return nil, err
	}
	return vasp, nil
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

func (i *vaspIterator) Id() string {
	// The VASP ID is prefix + uuid so strip off the prefix and return the string
	key := i.iter.Key()
	return string(key[len(preVASPs):])
}

func (i *vaspIterator) SeekId(vaspID string) bool {
	key := vaspKey(vaspID)
	return i.iter.Seek(key)
}

func (i *certIterator) Cert() (*models.Certificate, error) {
	r := new(models.Certificate)
	if err := proto.Unmarshal(i.iter.Value(), r); err != nil {
		sentry.Error(nil).Err(err).Str("type", wire.NamespaceCerts).Str("key", string(i.iter.Key())).Msg("corrupted data encountered")
		return nil, err
	}
	return r, nil
}

func (i *certIterator) All() (certs []*models.Certificate, err error) {
	certs = make([]*models.Certificate, 0)
	defer i.iter.Release()
	for i.iter.Next() {
		c := new(models.Certificate)
		if err = proto.Unmarshal(i.iter.Value(), c); err != nil {
			return nil, err
		}
		certs = append(certs, c)
	}

	if err = i.iter.Error(); err != nil {
		return nil, err
	}

	return certs, nil
}

func (i *certReqIterator) CertReq() (*models.CertificateRequest, error) {
	r := new(models.CertificateRequest)
	if err := proto.Unmarshal(i.iter.Value(), r); err != nil {
		sentry.Error(nil).Err(err).Str("type", wire.NamespaceCertReqs).Str("key", string(i.iter.Key())).Msg("corrupted data encountered")
		return nil, err
	}
	return r, nil
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

func (i *organizationIterator) ID() string {
	// The orgID is prefix + uuid so strip off the prefix and parse the UUID string
	key := i.iter.Key()
	orgID, err := uuid.FromBytes(key[len(preOrganizations):])
	if err != nil {
		panic(err)
	}
	return orgID.String()
}

func (i *organizationIterator) Organization() (o *bff.Organization, err error) {
	o = new(bff.Organization)
	if err = proto.Unmarshal(i.iter.Value(), o); err != nil {
		return nil, err
	}
	return o, nil
}

func (i *emailIterator) Email() (c *models.Email, err error) {
	c = &models.Email{}
	if err = proto.Unmarshal(i.iter.Value(), c); err != nil {
		return nil, err
	}
	return c, nil
}
