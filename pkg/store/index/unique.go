package index

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
)

// Unique indices map a key to a specified value, usually a search term or unique value
// to a record ID. This index prevents duplicates across records in the database and is
// also used for fast lookups matching values to documents.
type Unique map[string]string

// Add an entry to the unique index, normalizing if necessary. If the entry is already
// in the index, it will not be modified and false will be returned. See overwrite if
// the entry needs to be replaced.
func (c Unique) Add(key, value string, norm Normalizer) bool {
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
func (c Unique) Overwrite(key, value string, norm Normalizer) bool {
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

func (c Unique) Remove(key string, norm Normalizer) bool {
	if norm != nil {
		key = norm(key)
	}

	if _, ok := c[key]; !ok {
		return false
	}

	delete(c, key)
	return true
}

func (c Unique) Find(key string, norm Normalizer) (value string, ok bool) {
	if norm != nil {
		key = norm(key)
	}
	value, ok = c[key]
	return value, ok
}

func (c Unique) Reverse(value string, norm Normalizer) ([]string, bool) {
	results := make([]string, 0, 1)
	for key, val := range c {
		if value == val {
			results = append(results, key)
		}
	}

	return results, len(results) > 0
}

// Dump a unique index to a byte representation for storage on disk.
func (c Unique) Dump() (_ []byte, err error) {
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
func (c Unique) Load(data []byte) (err error) {
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

// Normalized unique maintains the normalizer with the unique index.
type normalizedUnique struct {
	index  Unique
	norm   Normalizer
	search Searcher
}

func (c *normalizedUnique) Add(key, value string) bool {
	return c.index.Add(key, value, c.norm)
}

func (c *normalizedUnique) Reverse(value string) ([]string, bool) {
	return c.index.Reverse(value, c.norm)
}

func (c *normalizedUnique) Find(key string) (string, bool) {
	return c.index.Find(key, c.norm)
}

func (c *normalizedUnique) Remove(key string) bool {
	return c.index.Remove(key, c.norm)
}

func (c *normalizedUnique) Overwrite(key, value string) bool {
	return c.index.Overwrite(key, value, c.norm)
}

func (c *normalizedUnique) Load(data []byte) error {
	return c.index.Load(data)
}

func (c *normalizedUnique) Dump() ([]byte, error) {
	return c.index.Dump()
}

func (c *normalizedUnique) Len() int {
	return len(c.index)
}

func (c *normalizedUnique) Empty() bool {
	return len(c.index) == 0
}

func (c *normalizedUnique) Search(query map[string]interface{}) []string {
	return c.search(query)
}
