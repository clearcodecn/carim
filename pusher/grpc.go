package pusher

import (
	"github.com/clearcodecn/car"
	"github.com/clearcodecn/carim/proto"
	"google.golang.org/grpc"
)

type ImServer struct {
	imService  *car.Server
	grpcServer *grpc.Server

	config Config
}


type Config struct {
	WebsocketIp   string
	WebsocketPort string

	GrpcIp   string
	GrpcPort string
}

func NewImServer(config Config) *ImServer {

	is := new(ImServer)
	grpcServer := grpc.NewServer()
	is.grpcServer = grpcServer

	cs := car.NewServer()
	is.imService = cs
	is.config = config

	proto.RegisterCarImPusherServer(is.grpcServer, is)

	return is
}
