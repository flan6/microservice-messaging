//go:generate protoc --go_out=../pb --go_opt=paths=source_relative --proto_path=../proto api.proto --go-grpc_out=../pb --go-grpc_opt=paths=source_relative
package server

import (
	"context"

	"github.com/flan6/microservice-messaging/internal/channel"
	pb "github.com/flan6/microservice-messaging/internal/generator/api/rpc/pb"
)

type Api struct {
	c channel.Channel

	pb.UnimplementedApiServer
}

func NewApi(c channel.Channel) Api {
	return Api{c: c}
}

func (a Api) Ping(_ context.Context, _ *pb.Empty) (*pb.Empty, error) {
	return nil, nil
}
