package secrets_test

import (
	"math/rand"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/secrets"
)

func TestCreateToken(t *testing.T) {
	vals := make(map[string]struct{})
	for i := 0; i < 5000; i++ {
		val := secrets.CreateToken(8)
		require.Len(t, val, 8)
		vals[val] = struct{}{}
	}
	require.Len(t, vals, 5000, "there is a very small chance that a duplicate value is generated")
}

func TestAlpha(t *testing.T) {
	// Test creating different random strings at different lengths
	for i := 0; i < 10000; i++ {
		len := rand.Intn(512) + 1
		alpha := secrets.Alpha(len)
		require.Len(t, alpha, len)
		require.Regexp(t, regexp.MustCompile(`[a-zA-Z]+`), alpha)
	}

	vals := make(map[string]struct{})
	for i := 0; i < 10000; i++ {
		val := secrets.Alpha(16)
		vals[val] = struct{}{}
	}
	require.Len(t, vals, 10000, "there is a very low chance that a duplicate value was generated")
}

func TestAlphaNumeric(t *testing.T) {
	// Test creating different random strings at different lengths
	for i := 0; i < 10000; i++ {
		len := rand.Intn(512) + 1
		alpha := secrets.AlphaNumeric(len)
		require.Len(t, alpha, len)
		require.Regexp(t, regexp.MustCompile(`[a-zA-Z0-9]+`), alpha)
	}

	vals := make(map[string]struct{})
	for i := 0; i < 10000; i++ {
		val := secrets.AlphaNumeric(16)
		vals[val] = struct{}{}
	}
	require.Len(t, vals, 10000, "there is a very low chance that a duplicate value was generated")
}

func TestCryptoRandInt(t *testing.T) {
	nums := make(map[uint64]struct{})
	for i := 0; i < 10000; i++ {
		val := secrets.CryptoRandInt()
		nums[val] = struct{}{}
	}
	require.Len(t, nums, 10000, "there is a very low chance that a duplicate value was generated")
}
