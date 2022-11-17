package bff

import (
	"errors"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/idna"
)

// Allow most domain names, including IDNs and muliple subdomains.
var validDomain = regexp.MustCompile(`^((xn--)?[a-z0-9][a-z0-9-_]*[a-z0-9]{0,1}\.(xn--)?)+([a-z0-9\-]+|[a-z0-9-]+\.[a-z]{2,})$`)

// Normalize a domain name for matching purposes.
func NormalizeDomain(domain string) (string, error) {
	// Attempt to parse as a URL, if successful then extract the domain without the
	// scheme or port.
	if u, err := url.ParseRequestURI(domain); err == nil {
		domain = u.Hostname()
	}

	// Convert to lowercase and remove trailing characters.
	domain = strings.ToLower(strings.TrimSpace(domain))

	// Convert international domain names to ASCII
	var err error
	if domain, err = idna.ToASCII(domain); err != nil {
		return "", err
	}

	return domain, nil
}

// Validation of domain names which also accepts internationalized domain names. This
// assumes that the domain name has already been normalized and converted to an ASCII
// compatible format.
func ValidateDomain(domain string) error {
	if !validDomain.MatchString(domain) {
		return errors.New("invalid domain name")
	}
	return nil
}
