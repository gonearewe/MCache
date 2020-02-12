package main

import (
	"github.com/gonearewe/MCache/cache"
	"github.com/gonearewe/MCache/server/tcp"
)

func main() {
	c := cache.NewCache("inmemory")
	tcp.NewServer(c).Listen()
}
