package cache

type Cache interface {
	Set(key string, val []byte) error
	Get(key string) (val []byte, err error)
	Del(key string) error
	GetStatus() Status
}

// Status records some statistic about the cache.
type Status struct {
	Count   int64 // number of caches
	KeySize int64 // size of keys
	ValSize int64 // size of values
}

func NewCache(type_ string, ttl int) Cache {
	if type_ == "inmemory" {
		return newInMemoryCache(ttl)
	} else {
		return newRocksdbCache(ttl)
	}
}

func (s *Status) add(key string, val []byte) {
	s.Count++
	s.KeySize += int64(len(key))
	s.ValSize += int64(len(val))
}

func (s *Status) reduce(key string, val []byte) {
	s.Count--
	s.KeySize -= int64(len(key))
	s.ValSize -= int64(len(val))
}
