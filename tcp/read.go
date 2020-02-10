package tcp

import (
	"bufio"
	"io"
)

// command = op key | key val
// op = 'S' | "G' | 'D' (set, get and delete)
// key = byte-array
// val = byte-array

// byte-array = length content
// length = 1*DIGIT (2 bytes)
// content = *OCTET

// response = success | error
// success = 0 byte-array
// error = 255 byte-array

// readKeyAndVal parses and returns key from given Reader.
func (s *Server) readKey(r *bufio.Reader) (string, error) {
	k, err := ReadByteArray(r)
	if err != nil {
		return "", err
	}

	return string(k), nil
}

// readKeyAndVal parses and returns key and val from given Reader.
func (s *Server) readKeyAndVal(r *bufio.Reader) (string, []byte, error) {
	k, err := ReadByteArray(r)
	if err != nil {
		return "", nil, err
	}

	v, err := ReadByteArray(r)
	if err != nil {
		return "", nil, err
	}

	return string(k), v, nil
}

func ReadByteArray(r *bufio.Reader) ([]byte, error) {
	// we use two bytes to store length info(their sum)
	num1, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	num2, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	num := int(num1) + int(num2) // sum up

	var array = make([]byte, num)
	_, err = io.ReadFull(r, array)
	if err != nil {
		return nil, err
	}

	return array, nil
}

func ResponseStatusCode(ok bool) byte {
	if ok {
		return 0
	}
	return 255
}

// ByteArrayLength calculates the length of given byte array, and put length
// info into a slice(len: 2) according to the protocol.
func ByteArrayLength(array []byte) []byte {
	var res = make([]byte, 2)
	length := len(array)
	if length >= 255 {
		res[1] = 255
		res[0] = byte(length - 255)
		return res
	}

	res[1] = byte(length)
	return res
}
