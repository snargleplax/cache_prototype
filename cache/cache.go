package cache

import (
	"container/heap"
	"sync"
	"time"
)

// A Cache is an associative in-memory store.
type Cache[T any] struct {
	m    sync.Map
	heap heap.Interface
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
func NewCache[T any](conf Config) *Cache[T] {
	var c Cache[T]
	c.heap = &expiryHeap[T]{keys: []string{}, c: &c}
	go c.checkExpiry(conf.ExpireCheck)
	return &c
}

// Put stores value v in Cache c under key k, replacing any existing value. The
// value is retained until the provided TTL expires.
func (c *Cache[T]) Put(k string, v T, ttl time.Duration) {
	c.m.Store(k, entry{v: v, expireAfter: time.Now().Add(ttl)})
	c.heap.Push(k)
}

// Get returns the value, if any, stored in Cache c under key k. If no value is
// stored, it returns nil.
func (c *Cache[T]) Get(k string) any {
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
func (c *Cache[T]) checkExpiry(cadence time.Duration) {
	for {
		for i := 0; i < c.heap.Len(); i++ {
			k := heap.Pop(c.heap).(string)
			v, _ := c.m.Load(k)
			if v == nil {
				continue
			}
			ent := v.(entry)
			if !ent.expired() {
				heap.Push(c.heap, k)
				break
			}
			c.m.Delete(k)
		}
		time.Sleep(cadence)
	}
}

// An expiryHeap is a "container/heap".Heap of Cache keys, sorted by expiry (oldest pops first).
type expiryHeap[T any] struct {
	keys []string
	c    *Cache[T]
}

func (h expiryHeap[T]) Len() int { return len(h.keys) }
func (h expiryHeap[T]) Less(i, j int) bool {
	lhs, _ := h.c.m.Load(h.keys[i])
	rhs, _ := h.c.m.Load(h.keys[j])
	return lhs.(entry).expireAfter.Before(rhs.(entry).expireAfter)
}
func (h expiryHeap[T]) Swap(i, j int) { h.keys[i], h.keys[j] = h.keys[j], h.keys[i] }
func (h *expiryHeap[T]) Push(x any)   { h.keys = append(h.keys, x.(string)) }
func (h *expiryHeap[T]) Pop() any {
	old := h.keys
	n := len(old)
	x := old[n-1]
	h.keys = old[0 : n-1]
	return x
}
