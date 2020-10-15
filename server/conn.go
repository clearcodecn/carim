package server

import (
	"encoding/binary"
	"encoding/json"
	"github.com/clearcodecn/carim/proto"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"time"
)

type Channel struct {
	conn      *websocket.Conn
	id        string
	writeChan chan []byte
}

type ChannelOption struct {
	id           string
	WriteBufSize int
	Timeout      time.Duration
}

func newChannel(wsconn *websocket.Conn, opt ChannelOption) *Channel {
	conn := new(Channel)
	conn.conn = wsconn
	conn.id = opt.id
	conn.writeChan = make(chan []byte, opt.WriteBufSize)
	go func() {
		for {
			for b := range conn.writeChan {
				if opt.Timeout != 0 {
					conn.conn.SetWriteDeadline(time.Now().Add(opt.Timeout))
				}
				err := conn.conn.WriteMessage(websocket.BinaryMessage, b)
				if err != nil {
					break
				}
				conn.conn.SetWriteDeadline(time.Time{})
			}
		}
	}()

	return conn
}

func (c *Channel) Close() error {
	close(c.writeChan)
	return c.conn.Close()
}

type Command int

const (
	CommandClose Command = iota + 1
	CommandPushToID
)

func (s *server) handleWebsocket(ctx *gin.Context) {
	wsConn, err := s.upgrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithError(400, err)
		return
	}
	id := s.generateNewID()

	conn := newChannel(wsConn, ChannelOption{
		id:           id,
		WriteBufSize: 1024,
		Timeout:      30 * time.Second,
	})
	s.channelMutex.Lock()
	s.channel[id] = conn
	s.channelMutex.Unlock()

	defer func() {
		conn.Close()
		s.channelMutex.Lock()
		delete(s.channel, id)
		s.channelMutex.Unlock()
	}()

	for {
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			break
		}
		cmd := Command(binary.BigEndian.Uint16(message[:2]))
		switch cmd {
		case CommandClose:
			break
		case CommandSendToUser: // handlers.
			var msg = proto.PushToId{}
			if err := decodeMessage(message[2:], &msg); err != nil {
				return
			}
			go s.pushToID(conn, &msg)
		}
	}
}

func decodeMessage(data []byte, msg interface{}) error {
	return json.Unmarshal(data, msg)
}
