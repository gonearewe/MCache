package http

import (
	"net/http"

	"github.com/gonearewe/MCache/cache"
)

// Cache servers running on self-made protocol on TCP perform better,
// so HTTP side should work for administrator, providing managing measures,
// which I didn't implement here. This package is currently unused.

const ServerPort = ":8999"

type Server struct {
	cache.Cache
}

func NewServer(cache cache.Cache) *Server {
	return &Server{cache}
}

func (s *Server) Listen() {
	http.Handle("/cache/", s.cacheHandler())
	http.Handle("/status", s.statusHandler())
	http.ListenAndServe(ServerPort, nil)
}
