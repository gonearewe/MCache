package cache

// #include "c.h"
// #include <stdlib.h>
// #cgo CFLAGS: -I${SRCDIR}/../rocksdb/include
// #cgo LDFLAGS: -L${SRCDIR}/../rocksdb -lrocksdb -lz -lpthread -lsnappy -lm -lstdc++ -O3
import "C"
import (
	"errors"
	"regexp"
	"runtime"
	"strconv"
	"unsafe"
)

type rocksdbCache struct {
	db *C.rocksdb_t
	ro *C.rocksdb_readoptions_t
	wo *C.rocksdb_writeoptions_t
	e  *C.char
	ch chan *pair
}

type pair struct {
	k string
	v []byte
}

func newRocksdbCache(ttl int) *rocksdbCache {
	const path = "/home/heathcliff/Projects/database" // database location

	options := C.rocksdb_options_create()
	C.rocksdb_options_increase_parallelism(options, C.int(runtime.NumCPU()))
	C.rocksdb_options_set_create_if_missing(options, 1)

	var e *C.char
	db := C.rocksdb_open_with_ttl(options, C.CString(path), C.int(ttl), &e)
	if e != nil {
		panic(C.GoString(e))
	}
	C.rocksdb_options_destroy(options)
	c := make(chan *pair, 5000)
	wo := C.rocksdb_writeoptions_create()
	go write_func(db, c, wo)

	return &rocksdbCache{db, C.rocksdb_readoptions_create(), wo, e, c}
}

func (c *rocksdbCache) Get(key string) ([]byte, error) {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	var length C.size_t
	v := C.rocksdb_get(c.db, c.ro, k, C.size_t(len(key)), &length, &c.e)
	if c.e != nil {
		return nil, errors.New(C.GoString(c.e))
	}
	defer C.free(unsafe.Pointer(v))
	return C.GoBytes(unsafe.Pointer(v), C.int(length)), nil
}

func (c *rocksdbCache) Del(key string) error {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	C.rocksdb_delete(c.db, c.wo, k, C.size_t(len(key)), &c.e)
	if c.e != nil {
		return errors.New(C.GoString(c.e))
	}
	return nil
}

func (c *rocksdbCache) GetStatus() Status {
	k := C.CString("rocksdb.aggregated-table-properties")
	defer C.free(unsafe.Pointer(k))
	v := C.rocksdb_property_value(c.db, k)
	defer C.free(unsafe.Pointer(v))
	p := C.GoString(v)
	r := regexp.MustCompile(`([^;]+)=([^;]+);`)
	s := Status{}
	for _, submatches := range r.FindAllStringSubmatch(p, -1) {
		if submatches[1] == " # entries" {
			s.Count, _ = strconv.ParseInt(submatches[2], 10, 64)
		} else if submatches[1] == " raw key size" {
			s.KeySize, _ = strconv.ParseInt(submatches[2], 10, 64)
		} else if submatches[1] == " raw value size" {
			s.ValSize, _ = strconv.ParseInt(submatches[2], 10, 64)
		}
	}
	return s
}
