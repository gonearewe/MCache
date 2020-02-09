package http

import (
	"net/http"

	"github.com/gonearewe/MCache/cache"
)

type Server struct {
	cache.Cache
}

func NewServer(cache cache.Cache) *Server {
	return &Server{cache}
}

func (s *Server) Listen() {
	http.Handle("/cache/", s.cacheHandler())
	http.Handle("/status", s.statusHandler())
	http.ListenAndServe(":8999", nil)
}
