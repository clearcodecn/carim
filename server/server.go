package server

import (
	"net"
	"sync"
)

type Server interface {
	Init(...Option) error

	Options() Options

	Serve() error

	Shutdown() error

	String() string
}

type server struct {
	stopChan chan chan error

	wg *sync.WaitGroup

	options Options

	started bool
	sync.RWMutex
}

func (s *server) Init(options Options) error {
	var opts []TransportOptions
	opts = append(opts, TransportAddress(s.options.Address))

	if err := s.options.Transport.Init(opts...); err != nil {
		return err
	}
}

func (s *server) Serve() error {
	s.RLock()
	if s.started {
		s.RUnlock()
		return nil
	}
	s.RUnlock()

	go func() {
	loop:
		for {
			select {
			case ch := <-s.stopChan:
				ch <- ln.Close()
				break loop
			default:
			}

			conn, err := ln.Accept()
			if err != nil {
				return
			}
			s.wg.Add(1)
			go s.accept(conn)
		}
	}()
	return nil
}

func (s *server) Shutdown() error {
	s.RLock()
	if !s.started {
		s.RUnlock()
		return nil
	}
	s.RUnlock()

	ch := make(chan error)
	s.stopChan <- ch

	err := <-ch
	s.Lock()
	s.started = false
	s.Unlock()

	return err
}

func (s *server) accept(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

}
