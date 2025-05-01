package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCacheGetNothing(t *testing.T) {
	// Given a key and a cache with no value for that key
	k := "test_key"
	c := NewCache(DefaultConfig)

	// When I retrieve from the cache using that key
	got := c.Get(k)

	// Then I get nothing back
	require.Nil(t, got)
}

func TestCachePutAndGet(t *testing.T) {
	// Given a key and a cache with no value for that key
	k := "test_key"
	c := NewCache(DefaultConfig)

	// When I put an integer value into the cache under that key
	v := 42
	c.Put(k, v, time.Hour)
	// And I retrieve from the cache using that key
	got := c.Get(k)

	// Then I get the same value I had put in
	require.Equal(t, v, got)
}

func TestCacheExpireTTLOnRead(t *testing.T) {
	// Given a value cached for some TTL
	k := "test_key"
	cadence := time.Hour
	c := NewCache(Config{ExpireCheck: cadence})
	v := 42
	ttl := 100 * time.Millisecond
	c.Put(k, v, ttl)

	// When I wait for a duration > the TTL (but not past the expire check cadence)
	time.Sleep(ttl + 1)
	// And I attempt to retrieve the cached value
	got := c.Get(k)

	// Then I get nothing back
	require.Nil(t, got)
}

func TestCacheExpireTTLInBackground(t *testing.T) {
	// Given a value cached for some TTL
	k := "test_key"
	cadence := 10 * time.Millisecond
	c := NewCache(Config{ExpireCheck: cadence})
	v := 42
	ttl := 100 * time.Millisecond
	c.Put(k, v, ttl)

	// When I wait for a duration > the TTL (plus the expire check cadence)
	time.Sleep(ttl + cadence)

	// Then the internal storage of the cache no longer contains the value.
	_, ok := c.m.Load(k)
	require.False(t, ok, "value is still stored in map")
}
