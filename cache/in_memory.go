package cache

import (
	"sync"
)

type inMemoryCache struct {
	cache  map[string][]byte
	lock   sync.RWMutex
	status Status
}

func newInMemoryCache() *inMemoryCache {
	return &inMemoryCache{
		cache:  make(map[string][]byte),
		lock:   sync.RWMutex{},
		status: Status{},
	}
}

func (c *inMemoryCache) Set(key string, val []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if v, exist := c.cache[key]; exist {
		c.status.reduce(key, v)
	}

	c.cache[key] = val
	c.status.add(key, val)

	return nil // always returns nil
}

func (c *inMemoryCache) Get(key string) (val []byte, err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.cache[key], nil
}

func (c *inMemoryCache) Del(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if val, exist := c.cache[key]; exist {
		delete(c.cache, key)
		c.status.reduce(key, val)
	}

	return nil
}

func (c *inMemoryCache) GetStatus() Status {
	return c.status
}
