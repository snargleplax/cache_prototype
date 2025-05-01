package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheGetNothing(t *testing.T) {
	// Given a key and a cache with no value for that key
	k := "test_key"
	c := NewCache()

	// When I retrieve from the cache using that key
	got := c.Get(k)

	// Then I get nothing back
	assert.Nil(t, got)
}

func TestCachePutAndGet(t *testing.T) {
	// Given a key and a cache with no value for that key
	k := "test_key"
	c := NewCache()

	// When I put an integer value into the cache under that key
	v := 42
	c.Put(k, v, time.Hour)
	// And I retrieve from the cache using that key
	got := c.Get(k)

	// Then I get the same value I had put in
	assert.Equal(t, v, got)
}

func TestCacheExpireTTL(t *testing.T) {
	// Given a value cached for some TTL
	k := "test_key"
	c := NewCache()
	v := 42
	ttl := 100 * time.Millisecond
	c.Put(k, v, ttl)

	// When I wait for a duration > the TTL
	time.Sleep(ttl + 1)
	// And I attempt to retrieve the cached value
	got := c.Get(k)

	// Then I get nothing back
	assert.Nil(t, got)
}
