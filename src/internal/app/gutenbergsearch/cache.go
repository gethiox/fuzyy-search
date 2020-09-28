package gutenbergsearch

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
}

type memCache struct {
	cache *cache.Cache
}

func (m *memCache) Get(key string) (interface{}, bool) {
	return m.cache.Get(key)
}

func (m *memCache) Set(key string, value interface{}) {
	m.cache.Set(key, value, cache.DefaultExpiration)
}

type dummyCache struct{}

func (d dummyCache) Get(key string) (interface{}, bool) { return nil, false }
func (d dummyCache) Set(key string, value interface{})  {}

func NewCache(enabled bool, expiration, cleanupInterval time.Duration) Cache {
	if !enabled {
		return &dummyCache{}
	}

	return &memCache{cache.New(expiration, cleanupInterval)}
}
