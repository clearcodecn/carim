package tcp

import (
	"bufio"
	"context"
	"errors"
	"github.com/clearcodecn/carim/pkg/network"
	"github.com/clearcodecn/carim/proto/protocol"
	"net"
	time "time"
)

type Conn struct {
	id     string
	conn   net.Conn
	ctx    context.Context
	server *Server

	readBuf  *bufio.Reader
	writeBuf *bufio.Writer
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

func (c *Conn) ID() string {
	return c.id
}

func (c *Conn) Authentication() (int, network.Identify, error) {
	if c.server.option.timeout != 0 {
		if err := c.SetTimeout(c.server.option.timeout); err != nil {
			return 0, nil, err
		}
	}
	var message = new(protocol.Message)
	n, err := c.ReadMessage(message)
	if err != nil {
		return 0, nil, err
	}
	if message.Operate != protocol.OpAuthenticate {
		return n, nil, errors.New("authenticate required")
	}
}

func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Conn) SetTimeout(d time.Duration) error {
	return c.conn.SetDeadline(time.Now().Add(d))
}

func (c *Conn) ReadMessage(message *protocol.Message) (n int, err error) {
	if message == nil {
		return 0, errors.New("message can't be nil")
	}
	return message.ReadFrom(c.readBuf)
}

func (c *Conn) WriteMessage(message *protocol.Message) (n int, err error) {
	return message.WriteTo(c.writeBuf)
}

func NewConn(ctx context.Context, conn net.Conn, s *Server) (*Conn, error) {
	cc := &Conn{
		ctx:    ctx,
		conn:   conn,
		server: s,
	}
	return cc, nil
}
