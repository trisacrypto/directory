package index

import (
	"encoding/binary"

	storeerrors "github.com/trisacrypto/directory/pkg/gds/store/errors"
)

// An auto-increment primary key Sequence for generating monotonically increasing IDs
type Sequence uint64

//===========================================================================
// Sequence
//===========================================================================

func (c Sequence) Dump() (data []byte, err error) {
	data = make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(data, uint64(c))
	return data, nil
}

func (c Sequence) Load(data []byte) (s Sequence, err error) {
	var n int
	var i uint64
	if i, n = binary.Uvarint(data); n <= 0 {
		return s, storeerrors.ErrCorruptedSequence
	}
	return Sequence(i), nil
}
