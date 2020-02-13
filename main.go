package main

import (
	"flag"

	"github.com/gonearewe/MCache/cache"
	"github.com/gonearewe/MCache/server/tcp"
)

func main() {
	type_ := flag.String("t", "inmemory", "cache type")
	ttl := flag.Int("ttl", 100, "cache life time, measured by second") // default ttl is 100s
	c := cache.NewCache(*type_, *ttl)
	tcp.NewServer(c).Listen() // by default, it runs on TCP
}
