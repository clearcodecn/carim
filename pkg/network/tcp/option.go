package tcp

import (
	"log"
	"net"
	"time"
)

type Option func(o *Options)

func WithTimeout(d time.Duration) Option {
	return func(options *Options) {
		options.timeout = d
	}
}

func WithAddress(addr string) Option {
	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		log.Fatal(err)
	}
	return func(o *Options) {
		o.address = addr
	}
}

type Options struct {
	timeout time.Duration
	address string

	writeBufferSize int
	readBufferSize  int
}

type ConnectionOption func(conn *Connection)

func ReadBufferSize(s int) ConnectionOption {
	return func(conn *Connection) {
		conn.readBufferSize = s
	}
}

func WriteBufferSize(s int) ConnectionOption {
	return func(conn *Connection) {
		conn.writeBufferSize = s
	}
}
