package tcp

import (
	"context"
	"github.com/clearcodecn/carim/pkg/atomic"
	"github.com/clearcodecn/carim/pkg/network"
	"github.com/clearcodecn/carim/proto/protocol"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
	"time"
)

type Server struct {
	ln net.Listener

	option Options

	logger logrus.FieldLogger

	isShutdown *atomic.Bool

	ctx   context.Context
	mutex sync.RWMutex

	connections map[string][]network.Connection
}

func NewServer(opts ...Option) *Server {
	var opt Options
	for _, o := range opts {
		o(&opt)
	}

	server := new(Server)

	server.option = opt

	return server
}

func (s *Server) Init() error {
	ln, err := net.Listen("tcp", s.option.address)
	if err != nil {
		return err
	}

	s.ln = ln

	return nil
}

func (s *Server) Serve(stop chan struct{}) error {
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				s.logger.WithError(err).Errorf("accept failed,retrying in %v", tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}

		go s.acceptConnection(conn)
	}
}

func (s *Server) acceptConnection(conn net.Conn) {
	var (
		readBuf  = make(chan *protocol.Message, s.option.readBufferSize)
		writeBuf = make(chan *protocol.Message, s.option.writeBufferSize)
	)
	connection, err := NewConn(s.ctx, conn, s, readBuf, writeBuf)
	if err != nil {
		s.logger.WithError(err).Errorf("create connection failed")
		return
	}

	defer func() {
		if s.OnConnectionClose != nil {
			s.OnConnectionClose(connection)
		}
		connection.Close()
	}()

	_, user, err := connection.Authentication()
	if err != nil {
		s.logger.WithError(err).Errorf("authenticate failed")
		return
	}

	s.mutex.Lock()
	if _, ok := s.connections[user.ID()]; !ok {
		s.connections[user.ID()] = make([]network.Connection, 0)
	}
	s.connections[user.ID()] = append(s.connections[user.ID()], connection)
	s.mutex.Unlock()

}
