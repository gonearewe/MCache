package tcp

import (
	"log"
	"net"

	"github.com/gonearewe/MCache/cache"
)

const ServerPort = ":8999"

type Server struct {
	cache.Cache
}

func NewServer(cache cache.Cache) *Server {
	return &Server{cache}
}

func (s *Server) Listen() {
	listener, err := net.Listen("tcp", ServerPort)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go s.process(conn)
	}
}
