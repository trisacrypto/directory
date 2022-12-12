package cache

// Disabled implements the Cache interface but implements no-op methods that return
// false. This allows callers to interact with a cache without having to check if it's
// enabled.
type Disabled struct{}

// Disabled implements the Cache interface.
var _ Cache = &Disabled{}

func (d *Disabled) Get(key interface{}) (data interface{}, ok bool) {
	return nil, false
}

func (d *Disabled) Add(key interface{}, data interface{}) {}

func (d *Disabled) Remove(key interface{}) {}
