package limitserver

import (
	"context"
	"net"

	"github.com/gofiber/recipes/fiber-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type LimitServer struct{}

func Start() {
	lis, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}
	srv := grpc.NewServer()
	proto.RegisterAddServiceServer(srv, &LimitServer{})
	reflection.Register(srv)

	if e := srv.Serve(lis); e != nil {
		panic(err)
	}
}

func (s *LimitServer) Add(_ context.Context, request *proto.Request) (*proto.Response, error) {
	a, b := request.GetA(), request.GetB()

	result := a + b

	return &proto.Response{Result: result}, nil
}

func (s *LimitServer) Multiply(_ context.Context, request *proto.Request) (*proto.Response, error) {
	a, b := request.GetA(), request.GetB()

	result := a * b

	return &proto.Response{Result: result}, nil
}
