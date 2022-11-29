package server_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/sectigo/server"
)

func TestSerialNumber(t *testing.T) {
	sn := server.SerialNumber()
	sns := fmt.Sprintf("%X", sn)
	require.Len(t, sns, 32)

	sn2 := server.SerialNumber()
	sns2 := fmt.Sprintf("%X", sn2)
	require.NotEqual(t, sns, sns2)
}
