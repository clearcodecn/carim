package server

import (
	"errors"
	"fmt"
	"github.com/clearcodecn/carim/proto"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type Config struct {
	Driver string
	DSN    string

	Key string

	Host string
	Port string

	PushClientAddr string
}

type HttpServer struct {
	engine *gin.Engine
	db     *xorm.Engine
	config Config

	pushClient proto.CarImPusherClient
}

func (h *HttpServer) Start() {
	hostPort := net.JoinHostPort(h.config.Host, h.config.Port)
	h.engine.Run(hostPort)
}

func NewHttpServer(config Config) *HttpServer {
	hs := new(HttpServer)
	hs.config = config

	client, err := grpc.Dial(config.PushClientAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	hs.pushClient = proto.NewCarImPusherClient(client)
	newGin(hs)
	newModel(hs, config)
	return hs
}

func newGin(server *HttpServer) {
	e := gin.Default()
	server.engine = e
	server.route()
}

func newModel(server *HttpServer, config Config) {
	e, err := xorm.NewEngine(config.Driver, config.DSN)
	if err != nil {
		log.Fatal(err)
	}
	if err := e.Ping(); err != nil {
		log.Fatal(err)
	}
	server.db = e

	server.syncTable()
}

func (h *HttpServer) route() {
	h.engine.GET("/", h.index)
	h.engine.POST("/login", h.login)
	h.engine.POST("/register", h.register)

	authGroup := h.engine.Group("/", h.authMiddleware)
	{
		authGroup.GET("/search-car-no", h.searchCarNo)
	}
}

// TODO:: create more useful unique car no.
func (h *HttpServer) uniqueCarNo() string {
	now := time.Now()
	sec := now.Unix()
	usec := now.UnixNano() % 0x100000
	return fmt.Sprintf("%s%08x%05x", "car-", sec, usec)
}

var (
	userCtxKey = "user-info"
)

func (h *HttpServer) authMiddleware(ctx *gin.Context) {
	//token := ctx.Request.Header.Get("Authorization")
	token, err := request.ParseFromRequest(ctx.Request, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.config.Key), nil
	})
	if err != nil {
		ctx.AbortWithError(401, errors.New("UnAuthorized"))
		return
	}
	id := token.Claims.(jwt.MapClaims)["id"]
	sess := h.db.NewSession()
	defer sess.Close()

	var user = new(User)
	ok, err := sess.Where("id = ?", id).Get(user)
	if err != nil {
		sql, args := sess.LastSQL()
		logrus.WithError(err).Warnf("query user failed: %s,%v", sql, args)
		ctx.AbortWithError(500, errors.New("server error"))
		return
	}
	if !ok {
		ctx.AbortWithError(401, errors.New("UnAuthorized"))
		return
	}
	ctx.Set(userCtxKey, user)
	ctx.Next()
}

func userFromContext(ctx *gin.Context) *User {
	user, ok := ctx.Get(userCtxKey)
	if ok {
		return user.(*User)
	}
	return nil
}
