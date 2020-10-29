package protocol

import (
	"bufio"
	"encoding/binary"
	"errors"
	"github.com/gogo/protobuf/proto"
	"io"
)

const (
	OpKnown = iota
	OpMessage
	OpCommand
	OpHeartBeat
)

const (
	headerSize  = 2
	maxBodySize = 2 << 15
)

var (
	ErrBodyTooLong = errors.New("body too long")
)

// ReadFrom read from buf and marshal message to protocol
func (m *Message) ReadFrom(r *bufio.Reader) (int64, error) {
	var (
		header     = make([]byte, headerSize)
		bodyLength int
		err        error
		body       []byte
		n          int
	)
	_, err = io.ReadFull(r, header)
	bodyLength = int(binary.BigEndian.Uint16(header))
	body = make([]byte, bodyLength)
	n, err = io.ReadFull(r, body)
	if err != nil {
		return headerSize + int64(n), err
	}
	err = proto.Unmarshal(body, m)
	return int64(headerSize + bodyLength), err
}

func (m *Message) WriteTo(w *bufio.Writer) (int64, error) {
	msg, err := proto.Marshal(m)
	if err != nil {
		return 0, err
	}
	size := len(msg)

	if size > maxBodySize {
		return 0, ErrBodyTooLong
	}
	header := make([]byte, headerSize)
	binary.BigEndian.PutUint16(header, uint16(size))
	n, err := w.Write(header)
	if err != nil {
		return 0, err
	}
	nn, err := w.Write(msg)
	if err != nil {
		return int64(n), err
	}
	if err = w.Flush(); err != nil {
		return 0, err
	}

	return int64(nn + n), nil
}
