package index

import (
	"net/url"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/trisa/pkg/iso3166"
)

const searchPrefixMinLength = 3

// Normalize specifies how to modify a string value in preparation for search.
type Normalizer func(s string) string

// Normalize an index value in preparation for index storage and lookups.
func Normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// Normalize country attempts to return an ISO 3166-1 alpha-2 code for the specified
// country, then stores that as the index record. If the normalized record can't be
// found then standard normalization on the string is used.
func NormalizeCountry(s string) string {
	s = Normalize(s)
	code, err := iso3166.Find(s)
	if err != nil {
		return s
	}
	return code.Alpha2
}

// Normalize URL parses a URL and only returns the hostname, if an error occurs, then
// an empty string is returned so the URL doesn't get added to the index.
func NormalizeURL(s string) string {
	u, err := url.Parse(Normalize(s))
	if err != nil {
		return ""
	}
	return u.Hostname()
}

// Searcher is a function that takes a query, parses it and returns a list of the found
// values in the index. Searcher functions are produced by the ExactMatch, PrefixMatch,
// and ContainsRecord functions to be used with different Indices.
type Searcher func(query map[string]interface{}) []string

// ExactMatch returns a search function that exactly matches the normalized value in
// the normalized unique index. This is used primarily for the websites index.
func (c *normalizedUnique) ExactMatch(indexName string) Searcher {
	return func(query map[string]interface{}) (results []string) {
		terms, ok := ParseQuery(indexName, query, c.norm)
		if ok {
			log.Debug().Str("index", indexName).Strs("terms", terms).Msg("exact match search")
			results = make([]string, 0)
			for _, term := range terms {
				if result, ok := c.index[term]; ok {
					// Only collect unique results in sorted order
					results = insort(results, result)
				}
			}
		}
		return results
	}
}

// PrefixMatch returns a search function that will match a normalized value if the term
// is of a minimum prefix length and a normalized unique index value has the prefix.
func (c *normalizedUnique) PrefixMatch(indexName string, searchPrefixMinLength int) Searcher {
	return func(query map[string]interface{}) (results []string) {
		terms, ok := ParseQuery(indexName, query, c.norm)
		if ok {
			log.Debug().Str("index", indexName).Strs("terms", terms).Msg("prefix match search")
			results = make([]string, 0)
			for _, term := range terms {
				if result, ok := c.index[term]; ok {
					// exact match
					// Only collect unique results in sorted order
					results = insort(results, result)
				} else if len(term) >= searchPrefixMinLength {
					// prefix match
					for record, result := range c.index {
						if strings.HasPrefix(record, term) {
							// Only collect unique results in sorted order
							results = insort(results, result)
						}
					}
				}
			}
		}
		return results
	}
}

// UniqueRecords returns a search function that gathers all of the unique IDs for the
// terms specified in the search query.
func (c *normalizedContainer) ContainsRecord(indexName string) Searcher {
	return func(query map[string]interface{}) (results []string) {
		terms, ok := ParseQuery(indexName, query, c.norm)
		if ok {
			log.Debug().Str("index", indexName).Strs("terms", terms).Msg("contains record search")
			results = make([]string, 0)
			for _, term := range terms {
				if records, ok := c.index[term]; ok {
					// Add only unique records to the results in sorted order.
					for _, record := range records {
						results = insort(results, record)
					}
				}
			}
		}
		return results
	}
}

// Queries are maps that hold an index name (e.g. "name" or "country") and map it to an
// indexable key in the index. The query can be either a single string or a list of
// strings; the parse function extracts the appropriate type and returns a list of
// strings if the query for the particular index name exists.
//
// All values are normalized by the normalize function if provided.
func ParseQuery(index string, query map[string]interface{}, norm Normalizer) ([]string, bool) {
	// Check to ensure the index is in the query
	val, ok := query[index]
	if !ok {
		return nil, false
	}

	// Parse the query based on the query type
	switch vals := val.(type) {
	case []string:
		terms := make([]string, 0, len(vals))
		for _, term := range vals {
			if norm != nil {
				term = norm(term)
			}
			terms = insort(terms, term)
		}
		return terms, true

	case string:
		if norm != nil {
			vals = norm(vals)
		}
		return []string{vals}, true

	default:
		return nil, false
	}
}

func insort(arr []string, item string) []string {
	if len(arr) == 0 {
		arr = append(arr, item)
		return arr
	}

	i := sort.Search(len(arr), func(i int) bool { return arr[i] >= item })
	if i < len(arr) && arr[i] == item {
		// value is already in the array
		return arr
	}

	arr = append(arr, "")
	copy(arr[i+1:], arr[i:])
	arr[i] = item
	return arr
}
