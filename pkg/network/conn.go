package network

import (
	"github.com/clearcodecn/carim/proto/protocol"
	"net"
	"time"
)

type Conn interface {
	Close() error

	ID() string

	Authentication() (Identify, error)

	LocalAddr() net.Addr

	RemoteAddr() net.Addr

	SetTimeout(d time.Duration) error

	ReadMessage(message *protocol.Message) (n int, err error)

	WriteMessage(message *protocol.Message) (n int, err error)
}
