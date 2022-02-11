package index

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"sort"
)

// Container indices map a value to a set of values, most commonly an index to a set of
// record IDs that are related to that value. The container index is maintained in
// sorted order with no duplicates to save space and enhance ordered search.
type Container map[string][]string

// Add a record ID to the index for the specified value. This method ensures that the
// list of records is maintained in sorted order with no duplicates.
func (c Container) Add(key, value string, norm Normalizer) bool {
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
func (c Container) Remove(key, value string, norm Normalizer) bool {
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
func (c Container) Find(key string, norm Normalizer) (value []string, ok bool) {
	if norm != nil {
		key = norm(key)
	}
	value, ok = c[key]
	return value, ok
}

// Reverse find - find all keys that index the specified value.
func (c Container) Reverse(value string, norm Normalizer) ([]string, bool) {
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
func (c Container) Contains(key string, value string, norm Normalizer) bool {
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
func (c Container) Dump() (data []byte, err error) {
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
func (c Container) Load(data []byte) (err error) {
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

// Normalized container maintains the normalizer with the container index.
type normalizedContainer struct {
	index  Container
	norm   Normalizer
	search Searcher
}

func (c *normalizedContainer) Add(key, value string) bool {
	return c.index.Add(key, value, c.norm)
}

func (c *normalizedContainer) Reverse(value string) ([]string, bool) {
	return c.index.Reverse(value, c.norm)
}

func (c *normalizedContainer) Find(key string) ([]string, bool) {
	return c.index.Find(key, c.norm)
}

func (c *normalizedContainer) Remove(key, value string) bool {
	return c.index.Remove(key, value, c.norm)
}

func (c *normalizedContainer) Contains(key, value string) bool {
	return c.index.Contains(key, value, c.norm)
}

func (c *normalizedContainer) Load(data []byte) error {
	return c.index.Load(data)
}

func (c *normalizedContainer) Dump() ([]byte, error) {
	return c.index.Dump()
}

func (c *normalizedContainer) Len() int {
	return len(c.index)
}

func (c *normalizedContainer) Empty() bool {
	return len(c.index) == 0
}

func (c *normalizedContainer) Search(query map[string]interface{}) []string {
	return c.search(query)
}
