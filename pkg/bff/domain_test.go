package bff_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff"
)

func TestNormalizeDomain(t *testing.T) {
	cases := []struct {
		domain string
		result string
		name   string
	}{
		{"foo.com", "foo.com", "already valid domain"},
		{"https://foo.com", "foo.com", "parse domain from URL"},
		{"http://foo.bar.com", "foo.bar.com", "parse subdomains from URL"},
		{"http://foo.com:1234", "foo.com", "parse domain from URL with port"},
		{"FOO.COM", "foo.com", "convert to lowercase"},
		{" foo.com ", "foo.com", "remove leading and trailing spaces"},
		{"♡.com", "xn--c6h.com", "convert to ASCII"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result, err := bff.NormalizeDomain(c.domain)
			require.NoError(t, err, "error normalizing domain")
			require.Equal(t, c.result, result)
		})
	}
}

func TestValidateDomain(t *testing.T) {
	// TODO: This is a good candidate for fuzz testing
	cases := []struct {
		domain string
		valid  bool
		name   string
	}{
		// Valid domains
		{"foo.com", true, "valid domain"},
		{"foo.bar.com", true, "valid domain with subdomain"},
		{"foo-bar.com", true, "valid domain with hyphen"},
		{"foo.bar.bar2.com", true, "valid domain with two subdomains"},
		{"foo.bar.bar2.bar3.bar4.com", true, "valid domain with multiple subdomains"},
		{"google.com.au", true, "valid domain with two trailing characters"},
		{"xn--bcher-kva.ch", true, "valid domain with international characters"},
		{"foo.xn--com", true, "valid domain with international characters and subdomain"},
		{"foo.bar.xn--com", true, "valid domain with international characters and subdomain"},
		{"xn-fsqu00a.xn-0zwm56d", true, "valid domain with international characters and numbers"},
		{"trisa-a1234b1234c1234d1234e1234f123456.example.longdomains.com", true, "valid domain with long subdomain"},

		// Invalid domains
		{"foo", false, "invalid domain with only letters"},
		{"12345", false, "invalid domain with only numbers"},
		{"alice@example.com", false, "invalid domain with email address"},
		{".foo.com", false, "invalid domain with leading period"},
		{"foo.com.", false, "invalid domain with trailing period"},
		{"https://foo.com", false, "invalid domain with protocol"},
		{"foo bar.com", false, "invalid domain with space"},
		{"foo+bar?.com", false, "invalid domain with invalid characters"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := bff.ValidateDomain(c.domain)
			if c.valid {
				require.NoError(t, err, "expected %s to be a valid domain", c.domain)
			} else {
				require.Error(t, err, "expected %s to be an invalid domain", c.domain)
			}
		})
	}
}
