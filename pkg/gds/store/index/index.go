package index

// Index types provide in-memory map index and helper methods for lookups and
// constraints to improve the performance of the store without disk access. The indices
// are intended to be checkpointed and synchronized to disk regularly but can also be
// rebuilt by stores that are Indexers.
//
// The indices map normalized string values (lower case, extra whitespace removed) to
// base64 encoded []byte keys or to UUID strings depending on the context.
type Index interface {
	Len() int
	Empty() bool
	Search(query map[string]interface{}) []string
	Add(key, value string) bool
	Reverse(value string) ([]string, bool)
}

// SingleIndex is not the best term for this but we will refactor this during the Trtl
// indexing sprint. This index is used to denote a mapping of indexed value to record ID.
type SingleIndex interface {
	Index
	Serializer
	Find(key string) (string, bool)
	Remove(key string) bool
	Overwrite(key, value string) bool
}

// MultiIndex is not the best term for this but we will refactor this during the Trtl
// indexing sprint. This index is used to done a mapping of indexed value to multiple
// record IDs.
type MultiIndex interface {
	Index
	Serializer
	Find(key string) ([]string, bool)
	Remove(key, value string) bool
	Contains(key, value string) bool
}

// Serializer allows indices to be loaded and dumped to disk.
type Serializer interface {
	Load([]byte) error
	Dump() ([]byte, error)
}

func NewNamesIndex() SingleIndex {
	idx := &normalizedUnique{
		index: make(Unique),
		norm:  Normalize,
	}

	idx.search = idx.PrefixMatch("name", searchPrefixMinLength)
	return idx
}

func NewWebsiteIndex() SingleIndex {
	idx := &normalizedUnique{
		index: make(Unique),
		norm:  NormalizeURL,
	}

	idx.search = idx.ExactMatch("website")
	return idx
}

func NewCountryIndex() MultiIndex {
	idx := &normalizedContainer{
		index: make(Container),
		norm:  NormalizeCountry,
	}

	idx.search = idx.ContainsRecord("country")
	return idx
}

func NewCategoryIndex() MultiIndex {
	idx := &normalizedContainer{
		index: make(Container),
		norm:  Normalize,
	}

	idx.search = idx.ContainsRecord("category")
	return idx
}
