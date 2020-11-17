package tcp

import (
	"bufio"
	"context"
	"errors"
	"github.com/clearcodecn/carim/pkg/atomic"
	"github.com/clearcodecn/carim/pkg/network"
	"github.com/clearcodecn/carim/proto/protocol"
	"github.com/google/uuid"
	"io"
	"net"
	time "time"
)

type ioBuffer struct {
	msg *protocol.Message
	n   int
	err error
}

type Connection struct {
	id     string
	conn   net.Conn
	ctx    context.Context
	server *Server

	readBuf  *bufio.Reader
	writeBuf *bufio.Writer

	readChan  chan *ioBuffer
	writeChan chan *ioBuffer

	isStop *atomic.Bool

	readBufSize  int
	writeBufSize int
}

func (c *Connection) Authentication() (int, network.Identify, error) {
	var message = new(protocol.Message)
	n, err := c.ReadMessage(message)
	if err != nil {
		return 0, nil, err
	}
	if message.Operate != protocol.OpAuthenticate {
		return n, nil, errors.New("authenticate required")
	}
	user, err := c.server.OnAuthenticate(message)
	return n, user, err
}

func (c *Connection) Device() network.Device {
	panic("implement me")
}

func (c *Connection) SetTimeout(d time.Duration) error {
	panic("implement me")
}

func (c *Connection) Close() error {
	return c.conn.Close()
}

func (c *Connection) ID() string {
	return c.id
}

func (c *Connection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Connection) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func (c *Connection) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *Connection) ReadMessage() (n int, err error) {
	message := new(protocol.Message)
	iobuf := &ioBuffer{
		msg: message,
	}

	return n, err
}

func (c *Connection) WriteMessage(message *protocol.Message) (n int, err error) {
	if c.server.option.timeout != 0 {
		if err := c.SetWriteDeadline(time.Now().Add(c.server.option.timeout)); err != nil {
			return 0, err
		}
	}
	return message.WriteTo(c.writeBuf)
}

type connOption struct {
	ctx    context.Context
	conn   net.Conn
	server *Server
}

func NewConn(opt connOption) (network.Connection, error) {
	id := uuid.New().String()
	cc := &Connection{
		id:        id,
		conn:      opt.conn,
		ctx:       opt.ctx,
		server:    opt.server,
		readBuf:   bufio.NewReader(opt.conn),
		writeBuf:  bufio.NewWriter(opt.conn),
		readChan:  make(chan *ioBuffer),
		writeChan: make(chan *ioBuffer),
		isStop:    atomic.New(),
	}

	go cc.writeLoop()
	go cc.readLoop()

	return cc, nil
}

func NewConnection(opts ...ConnectionOption) (network.Connection, error) {
	conn := new(Connection)
	for _, o := range opts {
		o(conn)
	}
}

func (c *Connection) writeLoop() {
	for {
		if c.isStop.IsSet() {
			return
		}
		select {
		case msg, ok := <-c.writeChan:
			if !ok {
				return
			}
			_, err := msg.WriteTo(c.writeBuf)
			if err != nil {
				if err != io.EOF {
					c.server.logger.WithError(err).Errorf("failed to write message")
				}
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Connection) readLoop() {
	for {
		if c.isStop.IsSet() {
			return
		}
		select {
		case msg, ok := <-c.readChan:
			if !ok {
				return
			}
			_, err := msg.ReadFrom(c.readBuf)
			if err != nil {
				if err != io.EOF {
					c.server.logger.WithError(err).Errorf("failed to read message")
				}
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}
