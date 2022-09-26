package utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/utils"
)

func TestLatest(t *testing.T) {
	alpha := "2022-04-07T20:04:21.000Z"
	bravo := "2022-04-07T08:36:07.000Z"

	require.Equal(t, alpha, utils.Latest(alpha, ""))
	require.Equal(t, bravo, utils.Latest("", bravo))
	require.Empty(t, utils.Latest("", ""))
	require.Equal(t, alpha, utils.Latest(alpha, bravo))
}
