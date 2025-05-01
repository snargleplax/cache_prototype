package cache

import (
	"sync"
	"time"
)

// A Cache is an associative in-memory store.
type Cache struct {
	m sync.Map
}

// An entry is an element stored in a Cache. It includes a record of its own
// expiry time.
type entry struct {
	v           any
	expireAfter time.Time
}

func (e entry) expired() bool {
	return time.Now().After(e.expireAfter)
}

// A Config expresses configuration options for a Cache.
type Config struct {
	// ExpireCheck is the desired approximate cadence of automatic expiry
	// checking. Setting this to a lower value than required may negatively
	// impact system performance.
	ExpireCheck time.Duration
}

// DefaultConfig is a sane default Config.
var DefaultConfig Config = Config{
	ExpireCheck: time.Minute,
}

// NewCache creates and returns a new, empty Cache.
func NewCache(conf Config) *Cache {
	var c Cache
	go c.checkExpiry(conf.ExpireCheck)
	return &c
}

// Put stores value v in Cache c under key k, replacing any existing value. The
// value is retained until the provided TTL expires.
func (c *Cache) Put(k string, v any, ttl time.Duration) {
	c.m.Store(k, entry{v: v, expireAfter: time.Now().Add(ttl)})
}

// Get returns the value, if any, stored in Cache c under key k. If no value is
// stored, it returns nil.
func (c *Cache) Get(k string) any {
	v, ok := c.m.Load(k)
	if !ok {
		return nil
	}
	ent := v.(entry)
	if ent.expired() {
		c.m.Delete(k)
		return nil
	}
	return ent.v
}

// checkExpiry is an endless loop that checks for and evicts expired cache
// entries. Checks are scheduled according to the specified cadence, to avoid a
// wasteful busy-loop.
func (c *Cache) checkExpiry(cadence time.Duration) {
	for {
		c.m.Range(func(k any, v any) bool {
			if v != nil && v.(entry).expired() {
				c.m.Delete(k)
			}
			return true
		})
		time.Sleep(cadence)
	}
}
