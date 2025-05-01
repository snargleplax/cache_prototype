package cache

// A Cache is an associative in-memory store.
type Cache struct {
	m map[string]any
}

// NewCache creates and returns a new, empty Cache.
func NewCache() *Cache {
	return &Cache{
		m: make(map[string]any),
	}
}

// Put stores value v in Cache c under key k, replacing any existing value.
func (c *Cache) Put(k string, v any) {
	c.m[k] = v
}

// Get returns the value, if any, stored in Cache c under key k. If no value is stored, it returns nil.
func (c *Cache) Get(k string) any {
	return c.m[k]
}
