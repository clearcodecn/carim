package network

import "net"

type Connection interface {
	net.Conn
}
