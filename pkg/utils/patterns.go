package utils

import (
	"errors"
	"regexp"
	"strings"
)

// From: https://stackoverflow.com/a/3824105/488917
var cnre = regexp.MustCompile(`^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])(\.([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9]))*$`)

// Validate a common name. The common name should not be empty, nor start with an "*"
// (e.g. a DNS wildcard). It should not start with a - and each label should be no more
// than 63 octets long. The common name should not have a scheme e.g. https:// prefix
// and it shouldn't have a port, e.g. example.com:443. Parsing is primarily based on
// a regular expression match from the cnre pattern.
func ValidateCommonName(name string) (err error) {
	if name == "" {
		return errors.New("common name should not be empty")
	}

	if strings.HasPrefix(name, "*") {
		return errors.New("wildcards are not allowed in TRISA common names")
	}

	if !cnre.MatchString(name) {
		return errors.New("common name does not match domain name regular expression")
	}
	return nil
}
