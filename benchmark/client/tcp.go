package client

import (
	"bufio"
	"errors"
	"net"

	"github.com/gonearewe/MCache/tcp"
)

type tcpClient struct {
	net.Conn
	reader *bufio.Reader // wrap a buf for reading from conn
}

func newTCPClient(server string) *tcpClient {
	conn, err := net.Dial("tcp", server+tcp.ServerPort)
	if err != nil {
		panic(err)
	}

	return &tcpClient{
		Conn:   conn,
		reader: bufio.NewReader(conn),
	}
}

func (c *tcpClient) Run(req *Request) {
	switch req.Type {
	case RequestGet:
		c.sendGet(req)
		req.Val, req.Error = c.recvResponse()
	case RequestSet:
		c.sendSet(req)
		_, req.Error = c.recvResponse()
	default:
		panic("unhandled request type")
	}
}

func (c *tcpClient) PipelineRun(reqs []*Request) {
	for _, req := range reqs { // send requests all at once
		switch req.Type {
		case RequestSet:
			c.sendSet(req)
		case RequestGet:
			c.sendGet(req)
		default:
			panic("unhandled request type")
		}
	}

	for _, req := range reqs { // receive responses in the same order
		switch req.Type {
		case RequestGet:
			req.Val, req.Error = c.recvResponse()
		case RequestSet:
			_, req.Error = c.recvResponse()
		default:
			panic("unhandled request type")
		}
	}
}

func (c *tcpClient) sendGet(req *Request) {
	key := []byte(req.Key)
	_, _ = c.Write([]byte{'G'})              // op
	_, _ = c.Write(tcp.ByteArrayLength(key)) // length
	_, _ = c.Write(key)                      // key
}

func (c *tcpClient) sendSet(req *Request) {
	key := []byte(req.Key)
	_, _ = c.Write([]byte{'S'})                  // op
	_, _ = c.Write(tcp.ByteArrayLength(key))     // key length
	_, _ = c.Write(key)                          // key
	_, _ = c.Write(tcp.ByteArrayLength(req.Val)) // value length
	_, _ = c.Write(req.Val)                      // value
}

func (c *tcpClient) recvResponse() ([]byte, error) {
	hasErr, err := hasErrCode(c.reader)
	if err != nil {
		return nil, err
	}

	if hasErr { // response reports error
		errInfo, err := tcp.ReadByteArray(c.reader)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(string(errInfo))
	}

	// response returns value we just query
	value, err := tcp.ReadByteArray(c.reader)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// hasErrCode reads one byte from given reader and tells if it's a error code.
func hasErrCode(r *bufio.Reader) (bool, error) {
	code, err := r.ReadByte()
	if err != nil {
		return false, err
	}

	if code == tcp.ResponseStatusCode(true) {
		return false, nil
	} else if code == tcp.ResponseStatusCode(false) {
		return true, nil
	} else {
		return false, errors.New("expecting response status code")
	}
}
