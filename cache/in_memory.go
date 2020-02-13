package cache

import (
	"sync"
	"time"
)

type inMemoryCache struct {
	cache  map[string]value
	lock   sync.RWMutex
	status Status
	ttl    time.Duration // time to live, used for FIFO GC
}

type value struct {
	val          []byte
	lastUsedTime time.Time // used for GC
}

// newInMemoryCache creates an in memory cache, ttl is the interval(measured by second)
// of GC; if ttl is zero, then no GC.
func newInMemoryCache(ttl int) *inMemoryCache {
	c := &inMemoryCache{
		cache:  make(map[string]value),
		lock:   sync.RWMutex{},
		status: Status{},
		ttl:    time.Duration(ttl) * time.Second,
	}

	if ttl > 0 { // run GC if ttl is set
		go c.gc()
	}

	return c
}

func (c *inMemoryCache) Set(key string, val []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if v, exist := c.cache[key]; exist {
		c.status.reduce(key, v.val)
	}

	c.cache[key] = value{val: val, lastUsedTime: time.Now()}
	c.status.add(key, val)

	return nil // always returns nil
}

func (c *inMemoryCache) Get(key string) (val []byte, err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.cache[key].val, nil
}

func (c *inMemoryCache) Del(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if v, exist := c.cache[key]; exist {
		delete(c.cache, key)
		c.status.reduce(key, v.val)
	}

	return nil
}

func (c *inMemoryCache) GetStatus() Status {
	return c.status
}

// gc works as a routine in a loop, in each loop, it sleeps for ttl and then traverses the whole
// cache, deleting expired ones.
func (c *inMemoryCache) gc() {
	for {
		time.Sleep(c.ttl) // GC interval is ttl

		c.lock.RLock()
		for k, v := range c.cache {
			c.lock.RUnlock()                                  // unlock RLock since Del will acquires a lock
			if v.lastUsedTime.Add(c.ttl).Before(time.Now()) { // expire
				_ = c.Del(k)
			}
			c.lock.RLock()
		}
		c.lock.RUnlock()
	}
}
