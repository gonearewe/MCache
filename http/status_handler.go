package http

import (
	"encoding/json"
	"log"
	"net/http"
)

type statusHandler struct {
	*Server
}

func (s *Server) statusHandler() http.Handler {
	return &statusHandler{s}
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
