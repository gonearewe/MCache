package main

import (
	"github.com/gonearewe/MCache/cache"
	"github.com/gonearewe/MCache/tcp"
)

func main() {
	c := cache.NewCache()
	tcp.NewServer(c).Listen()
}
