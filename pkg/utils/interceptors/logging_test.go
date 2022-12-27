package interceptors_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/utils/interceptors"
)

func TestParseMethod(t *testing.T) {
	tests := []struct {
		FullMethod string
		service    string
		rpc        string
	}{
		{"package.v1.Service/Echo", "package.v1.Service", "Echo"},
		{"Service/Echo", "Service", "Echo"},
	}

	for _, tc := range tests {
		service, rpc := interceptors.ParseMethod(tc.FullMethod)
		require.Equal(t, tc.service, service, "unexpected service parsed from %q", tc.FullMethod)
		require.Equal(t, tc.rpc, rpc, "unexpected rpc parsed from %q", tc.FullMethod)
	}
}
