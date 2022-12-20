package cache_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/utils/cache"
)

func TestDisabled(t *testing.T) {
	disabled := &cache.Disabled{}

	// Add should not panic
	disabled.Add("foo", "bar")

	// Get should return false and nil values
	val, ok := disabled.Get("foo")
	require.False(t, ok, "disabled cache should return false")
	require.Nil(t, val, "disabled cache should return nil values")

	// Remove should not panic
	disabled.Remove("foo")
}
