package leveldb

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/trisacrypto/trisa/pkg/iso3166"
)

//===========================================================================
// Index Types
//===========================================================================

// Index types provide in-memory map index and helper methods for lookups and
// constraints to improve the performance of the store without disk access. The indices
// are intended to be checkpointed and synchronized to disk regularly but can also be
// rebuilt by stores that are Indexers.
//
// The indices map normalized string values (lower case, extra whitespace removed) to
// base64 encoded []byte keys or to UUID strings depending on the context.
type uniqueIndex map[string]string
type containerIndex map[string][]string

// An auto-increment primary key sequence for generating monotonically increasing IDs
type sequence uint64

//===========================================================================
// Utility Functions
//===========================================================================

// Normalize specifies how to modify a string value in preparation for search.
type normalizer func(s string) string

// Normalize an index value in preparation for index storage and lookups.
func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// Normalize country attempts to return an ISO 3166-1 alpha-2 code for the specified
// country, then stores that as the index record. If the normalized record can't be
// found then standard normalization on the string is used.
func normalizeCountry(s string) string {
	code, err := iso3166.Find(s)
	if err != nil {
		return normalize(s)
	}
	return code.Alpha2
}

// Normalize URL parses a URL and only returns the hostname, if an error occurs, then
// an empty string is returned so the URL doesn't get added to the index.
func normalizeURL(s string) string {
	u, err := url.Parse(normalize(s))
	if err != nil {
		return ""
	}
	return u.Hostname()
}

// Queries are maps that hold an index name (e.g. "name" or "country") and map it to an
// indexable key in the index. The query can be either a single string or a list of
// strings; the parse function extracts the appropriate type and returns a list of
// strings if the query for the particular index name exists.
// All values are normalized by the normalize function if provided.
func parseQuery(index string, query map[string]interface{}, norm normalizer) ([]string, bool) {
	val, ok := query[index]
	if !ok {
		return nil, false
	}

	if vals, ok := val.([]string); ok {
		for i := range vals {
			if norm != nil {
				vals[i] = norm(vals[i])
			}
		}
		return vals, true
	}

	if vals, ok := val.(string); ok {
		if norm != nil {
			vals = norm(vals)
		}
		return []string{vals}, true
	}

	return nil, false
}

//===========================================================================
// Unique Index
//===========================================================================

// Add an entry to the unique index, normalizing if necessary. If the entry is already
// in the index, it will not be modified and false will be returned. See overwrite if
// the entry needs to be replaced.
func (c uniqueIndex) add(key, value string, norm normalizer) bool {
	if norm != nil {
		key = norm(key)
	}

	// Do not add empty strings to the index
	if key == "" {
		return false
	}

	if _, ok := c[key]; ok {
		return false
	}

	c[key] = value
	return true
}

// Add or replace an entry that exists in the index.
func (c uniqueIndex) overwrite(key, value string, norm normalizer) bool {
	if norm != nil {
		key = norm(key)
	}

	// Do not add empty strings to the index
	if key == "" {
		return false
	}

	c[key] = value
	return true
}

func (c uniqueIndex) rm(key string, norm normalizer) bool {
	if norm != nil {
		key = norm(key)
	}

	if _, ok := c[key]; !ok {
		return false
	}

	delete(c, key)
	return true
}

func (c uniqueIndex) find(key string, norm normalizer) (value string, ok bool) {
	if norm != nil {
		key = norm(key)
	}
	value, ok = c[key]
	return value, ok
}

func (c uniqueIndex) reverse(value string, norm normalizer) ([]string, bool) {
	if norm != nil {
		value = norm(value)
	}

	results := make([]string, 0, 1)
	for key, val := range c {
		if value == val {
			results = append(results, key)
		}
	}

	return results, len(results) > 0
}

// Dump a unique index to a byte representation for storage on disk.
func (c uniqueIndex) Dump() (_ []byte, err error) {
	// Create a compressed writer to encode JSON into
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	// Marshal the JSON representation of the index
	encoder := json.NewEncoder(gz)
	if err = encoder.Encode(c); err != nil {
		return nil, fmt.Errorf("could not encode index: %s", err)
	}

	if err = gz.Close(); err != nil {
		return nil, fmt.Errorf("could not compress index: %s", err)
	}

	return buf.Bytes(), nil
}

// Load a unique index from a byte representation on disk.
func (c uniqueIndex) Load(data []byte) (err error) {
	// Create a compressed reader to decode the JSON from.
	buf := bytes.NewBuffer(data)

	var gz *gzip.Reader
	if gz, err = gzip.NewReader(buf); err != nil {
		return fmt.Errorf("could not decompress index: %s", err)
	}

	decoder := json.NewDecoder(gz)
	if err = decoder.Decode(&c); err != nil {
		return fmt.Errorf("could not decode index: %s", err)
	}
	return nil
}

//===========================================================================
// Country Container Index
//===========================================================================

// Add a record ID to the index for the specified country. This method ensures that the
// list of records for the country is maintained in sorted order with no duplicates.
func (c containerIndex) add(key, value string, norm normalizer) bool {
	// normalize the key to enhance search
	if norm != nil {
		key = norm(key)
	}

	// do not add empty values
	if key == "" {
		return false
	}

	arr, ok := c[key]
	if !ok {
		c[key] = append(make([]string, 0, 10), value)
		return true
	}

	i := sort.Search(len(arr), func(i int) bool { return arr[i] >= value })
	if i < len(arr) && arr[i] == value {
		// value is already in the array
		return false
	}

	arr = append(arr, "")
	copy(arr[i+1:], arr[i:])
	arr[i] = value
	c[key] = arr
	return true
}

// Remove a value from the container specified by the key, normalizing the key if
// necessary. This method ensures that the containers are maintained in sorted order.
func (c containerIndex) rm(key, value string, norm normalizer) bool {
	// normalize the key to enhance search
	if norm != nil {
		key = norm(key)
	}

	arr, ok := c[key]
	if !ok {
		return false
	}

	i := sort.Search(len(arr), func(i int) bool { return arr[i] >= value })
	if i < len(arr) && arr[i] == value {
		copy(arr[i:], arr[i+1:])
		arr[len(arr)-1] = ""
		arr = arr[:len(arr)-1]
		c[key] = arr
		return true
	}
	return false
}

// Find all values for the specified key in the index, returns nil if it doesn't exist.
func (c containerIndex) find(key string, norm normalizer) (value []string, ok bool) {
	if norm != nil {
		key = norm(key)
	}
	value, ok = c[key]
	return value, ok
}

// Reverse find - find all keys that index the specified value.
func (c containerIndex) reverse(value string, norm normalizer) ([]string, bool) {
	if norm != nil {
		value = norm(value)
	}

	results := make([]string, 0)
	for key, arr := range c {
		i := sort.Search(len(arr), func(i int) bool { return arr[i] >= value })
		if i < len(arr) && arr[i] == value {
			results = append(results, key)
		}
	}

	return results, len(results) > 0
}

// Contains determines if the value is contained by the key index
func (c containerIndex) contains(key string, value string, norm normalizer) bool {
	// normalize the key to enhance search
	if norm != nil {
		key = norm(key)
	}

	// find the container for the index key
	arr, ok := c[key]
	if !ok {
		return false
	}

	// perform a binary search to determine if the value is in the array
	i := sort.Search(len(arr), func(i int) bool { return arr[i] >= value })
	if i < len(arr) && arr[i] == value {
		// value is already in the array
		return true
	}

	return false
}

// Dump a container index to a byte representation for storage on disk.
func (c containerIndex) Dump() (data []byte, err error) {
	// Create a compressed writer to encode JSON into
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	// Marshal the JSON representation of the index
	encoder := json.NewEncoder(gz)
	if err = encoder.Encode(c); err != nil {
		return nil, fmt.Errorf("could not encode index: %s", err)
	}

	if err = gz.Close(); err != nil {
		return nil, fmt.Errorf("could not compress index: %s", err)
	}

	return buf.Bytes(), nil
}

// Load a container index from a byte representation on disk.
func (c containerIndex) Load(data []byte) (err error) {
	// Create a compressed reader to decode the JSON from.
	buf := bytes.NewBuffer(data)

	var gz *gzip.Reader
	if gz, err = gzip.NewReader(buf); err != nil {
		return fmt.Errorf("could not decompress index: %s", err)
	}

	decoder := json.NewDecoder(gz)
	if err = decoder.Decode(&c); err != nil {
		return fmt.Errorf("could not decode index: %s", err)
	}
	return nil
}

//===========================================================================
// Sequence
//===========================================================================

func (c sequence) Dump() (data []byte, err error) {
	data = make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(data, uint64(c))
	return data, nil
}

func (c sequence) Load(data []byte) (s sequence, err error) {
	var n int
	var i uint64
	if i, n = binary.Uvarint(data); n <= 0 {
		return s, ErrCorruptedSequence
	}
	return sequence(i), nil
}
