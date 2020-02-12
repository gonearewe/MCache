package tcp

import (
	"bufio"
	"io"
	"log"
	"net"
)

// process serves for conn as a routine, handling its requests.
func (s *Server) process(conn net.Conn) {
	defer conn.Close()
	var command = bufio.NewReader(conn)
	// handle requests in a loop, thus multiple requests in one connection is allowed
	for {
		op, err := command.ReadByte()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Println("request command's format error")
			return
		}
		// err == nil here

		switch op {
		case 'S':
			err = s.set(conn, command)
		case 'G':
			err = s.get(conn, command)
		case 'D':
			err = s.delete(conn, command)
		default:
			log.Println("unknown op code: ", op)
			return
		}

		if err != nil {
			log.Panicln(err)
			return
		}
	}
}

func (s *Server) set(conn net.Conn, r *bufio.Reader) error {
	if k, v, err := s.readKeyAndVal(r); err != nil {
		return err
	} else {
		return s.sendResponse(conn, s.Set(k, v), nil)
	}
}

func (s *Server) get(conn net.Conn, r *bufio.Reader) error {
	k, err := s.readKey(r)
	if err != nil {
		return err
	}

	v, err := s.Get(k)
	return s.sendResponse(conn, err, v)
}

func (s *Server) delete(conn net.Conn, r *bufio.Reader) error {
	k, err := s.readKey(r)
	if err != nil {
		return err
	}

	return s.sendResponse(conn, s.Del(k), nil)
}

// sendResponse writes response to conn, reporting error when err!=nil or reporting value.
func (s *Server) sendResponse(conn net.Conn, err error, value []byte) error {
	if err != nil {
		_, e := conn.Write(append([]byte{ResponseStatusCode(false)}, []byte(err.Error())...))
		return e
	}

	array := append(ByteArrayLength(value), value...)
	_, e := conn.Write(append([]byte{ResponseStatusCode(true)}, array...))
	return e
}
