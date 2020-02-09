package main

import (
	"github.com/gonearewe/MCache/cache"
	"github.com/gonearewe/MCache/http"
)

func main() {
	c := cache.NewCache()
	http.NewServer(c).Listen()
}
