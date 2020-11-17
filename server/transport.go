package server

import "github.com/clearcodecn/carim/proto/protocol"

type Socket interface {
	Recv(message *protocol.Message) error
	Send(message *protocol.Message) error
	Close() error
}

type Transport interface {
	Init(option ...TransportOptions) error
	Accept(func(socket Socket)) error
	Close() error
}

type TransportOptions func(option *TransportOption)

type TransportOption struct {
	Addr string
}

func TransportAddress(addr string) TransportOptions {
	return func(option *TransportOption) {
		option.Addr = addr
	}
}
