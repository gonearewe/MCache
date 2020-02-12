package tcp

import (
	"bufio"
	"io"
	"log"
	"net"
)

type result struct {
	value []byte
	err   error
}

// process serves for conn as a routine, handling its requests.
func (s *Server) process(conn net.Conn) {
	var command = bufio.NewReader(conn)

	var resultsChan = make(chan chan *result, 5000)
	defer close(resultsChan) // when process quits naturally or not, reply routine receives a signal and quits
	go s.reply(conn, resultsChan)

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

		var resChan = make(chan *result, 1)
		resultsChan <- resChan
		switch op {
		case 'S':
			s.set(resChan, command)
		case 'G':
			s.get(resChan, command)
		case 'D':
			s.delete(resChan, command)
		default:
			log.Println("unknown op code: ", op)
			return
		}
	}
}

func (s *Server) set(res chan<- *result, r *bufio.Reader) {
	if k, v, err := s.readKeyAndVal(r); err != nil {
		res <- &result{value: nil, err: err}
		return
	} else {
		go func() {
			res <- &result{
				value: nil,
				err:   s.Set(k, v),
			}
		}()
	}
}

func (s *Server) get(res chan<- *result, r *bufio.Reader) {
	k, err := s.readKey(r)
	if err != nil {
		res <- &result{value: nil, err: err}
		return
	}

	go func() {
		v, e := s.Get(k)
		res <- &result{
			value: v,
			err:   e,
		}
	}()
}

func (s *Server) delete(res chan<- *result, r *bufio.Reader) {
	k, err := s.readKey(r)
	if err != nil {
		return
	}

	go func() {
		res <- &result{
			value: nil,
			err:   s.Del(k),
		}
	}()
}

func (s *Server) reply(conn net.Conn, resultsChan <-chan chan *result) {
	// the right to close connection is handed over to reply routine
	// since errors during response means connection failure,
	// then reply's father will fail when reading from connection and also quits
	defer conn.Close()
	for {
		resChan, open := <-resultsChan
		if !open {
			return
		}
		res := <-resChan

		if err := s.sendResponse(conn, res.err, res.value); err != nil {
			log.Println(err)
			return
		}
	}
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
