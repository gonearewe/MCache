package http

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type cacheHandler struct {
	*Server
}

func (s *Server) cacheHandler() http.Handler {
	return &cacheHandler{s}
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
