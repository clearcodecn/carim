package tcp

import (
	"context"
	"github.com/clearcodecn/carim/pkg/atomic"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

type Server struct {
	ln net.Listener

	option Options

	logger logrus.FieldLogger

	isShutdown *atomic.Bool

	ctx context.Context
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


	}
}
