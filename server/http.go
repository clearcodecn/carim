package server

import (
	"github.com/clearcodecn/carim/proto"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type Config struct {
	Ip       string
	Port     string
	HostName string
}

type server struct {
	config Config

	engine *gin.Engine

	clients     map[string]proto.PusherClient
	clientMutex sync.Mutex

	checkOrigin func(r *http.Request) bool
	upgrade     *websocket.Upgrader

	channel      map[string]*Channel
	channelMutex sync.Mutex
}

func (s *server) StartHTTPServer() {
	s.engine.Run(":1111")
}

func (s *server) init() {
	g := gin.Default()
	s.engine = g

	s.engine.GET("/ws", s.handleWebsocket)
}
