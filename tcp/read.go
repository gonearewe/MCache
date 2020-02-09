package tcp

import (
	"bufio"
	"io"
)

// command = op key | key val
// op = 'S' | "G' | 'D" (set, get and delete)
// key = byte-array
// val = byte-array

// byte-array = length content
// length = 1*DIGIT (2 bytes)
// content = *OCTET

// response = byte-array | error
// error = 255 255 byte-array

// readKeyAndVal parses and returns key from given Reader.
func (s *Server) readKey(r *bufio.Reader) (string, error) {
	k, err := readByteArray(r)
	if err != nil {
		return "", err
	}

	return string(k), nil
}

// readKeyAndVal parses and returns key and val from given Reader.
func (s *Server) readKeyAndVal(r *bufio.Reader) (string, []byte, error) {
	k, err := readByteArray(r)
	if err != nil {
		return "", nil, err
	}

	v, err := readByteArray(r)
	if err != nil {
		return "", nil, err
	}

	return string(k), v, nil
}

func readByteArray(r *bufio.Reader) ([]byte, error) {
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

// errCode returns a slice(len: 2) as a response error code.
func errCode() []byte {
	return []byte{255, 255}
}

// byteArrayLength calculates the length of given byte array, and put length
// info into a slice(len: 2) according to the protocol.
func byteArrayLength(array []byte) []byte {
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
