package main

type Cache interface {
	Set(key string, val []byte) error
	Get(key string) (val []byte, err error)
	Del(key string) error
	GetStatus() Status
}

type Status struct {
	Count   int64 // number of caches
	KeySize int64 // size of keys
	ValSize int64 // size of values
}

func NewCache() Cache {
	return newInMemoryCache()
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
