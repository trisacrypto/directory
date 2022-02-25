package index_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/store/index"
)

func TestSequence(t *testing.T) {
	var pk index.Sequence
	require.Equal(t, uint64(pk), uint64(0))

	// Perform some operations on the pk
	for i := 0; i < 3; i++ {
		pk++
	}

	require.Equal(t, uint64(pk), uint64(3))

	var id uint64
	if pk > 0 {
		id = uint64(pk)
	}

	for i := 0; i < 5; i++ {
		pk++
	}

	require.Equal(t, uint64(3), id)
	require.Equal(t, uint64(8), uint64(pk))
	require.Equal(t, index.Sequence(8), pk)

	data, err := pk.Dump()
	require.NoError(t, err)
	require.NotEmpty(t, data)

	for i := 0; i < 4; i++ {
		pk++
	}

	other, err := pk.Load(data)
	require.NoError(t, err)
	require.Equal(t, index.Sequence(12), pk)
	require.NotEqual(t, other, pk)
	require.Equal(t, index.Sequence(8), other)
}
