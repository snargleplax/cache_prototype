package cache

import "time"

// A Cache is an associative in-memory store.
type Cache struct {
	m map[string]entry
}

type entry struct {
	v           any
	expireAfter time.Time
}

// NewCache creates and returns a new, empty Cache.
func NewCache() *Cache {
	return &Cache{
		m: make(map[string]entry),
	}
}

// Put stores value v in Cache c under key k, replacing any existing value. The
// value is retained until the provided TTL expires.
func (c *Cache) Put(k string, v any, ttl time.Duration) {
	c.m[k] = entry{v: v, expireAfter: time.Now().Add(ttl)}
}

// Get returns the value, if any, stored in Cache c under key k. If no value is
// stored, it returns nil.
func (c *Cache) Get(k string) any {
	ent, ok := c.m[k]
	if !ok {
		return nil
	}
	if time.Now().After(ent.expireAfter) {
		delete(c.m, k)
		return nil
	}
	return ent.v
}
