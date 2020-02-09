package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	Cache
}

type cacheHandler struct {
	*Server
}

type statusHandler struct {
	*Server
}

func NewServer(cache Cache) *Server {
	return &Server{cache}
}

func (s *Server) Listen() {
	http.Handle("/cache/", s.cacheHandler())
	http.Handle("/status", s.statusHandler())
	http.ListenAndServe(":2000", nil)
}

func (s *Server) cacheHandler() http.Handler {
	return &cacheHandler{s}
}

func (s *Server) statusHandler() http.Handler {
	return &statusHandler{s}
}

func (h *cacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.EscapedPath(), "/") // /cache/key
	if len(s) < 3 || len(s[2]) == 0 {            // wrong URL or empty key
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var key = s[2]
	switch r.Method {
	case http.MethodPut:
		val, _ := ioutil.ReadAll(r.Body)
		if len(val) != 0 {
			if err := h.Set(key, val); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		return

	case http.MethodGet:
		val, err := h.Get(key)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(val) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if _, err := w.Write(val); err != nil {
			log.Println(err)
		}
		return

	case http.MethodDelete:
		err := h.Del(key)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (h *statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if data, err := json.Marshal(h.GetStatus()); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		_, _ = w.Write(data)
		return
	}
}
