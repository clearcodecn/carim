package server

import (
	"errors"
	"github.com/google/uuid"
	"net"
	"strings"
)

func (s *server) generateNewID() string {
	uid := uuid.New().String()
	ip := s.config.HostName
	if ip == "" {
		ip = s.config.Ip
	}
	return ip + ":" + s.config.Port + "|" + uid
}

func (s *server) extractID(id string) (host, port string, err error) {
	arr := strings.Split(id, "|")
	if len(arr) != 2 {
		return "", "", errors.New("invalid id")
	}

	return net.SplitHostPort(arr[0])
}
