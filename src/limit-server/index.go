package limitserver

import (
	"context"
	"net"

	proto "github.com/rpsoftech/bullion-server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type LimitServerServer struct{}

func (LimitServerServer) PlaceLimit(context.Context, *proto.UplinkPlaceLimitRequest) (*proto.UplinkPlaceLimitResponse, error) {
	panic("unimplemented")
}
func (LimitServerServer) testEmbeddedByValue() {}
func (LimitServerServer) PlaceLimitStream(grpc.BidiStreamingServer[proto.UplinkPlaceLimitRequest, proto.UplinkPlaceLimitResponse]) error {
	panic("unimplemented")
}

func (LimitServerServer) mustEmbedUnimplementedLimitServerServer() {
	panic("unimplemented")
}

func Start() {
	lis, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}
	srv := grpc.NewServer(
	// grpc.UnaryInterceptor()
	)
	// proto.re
	proto.RegisterLimitServerServer(srv, &LimitServerServer{})
	reflection.Register(srv)

	if e := srv.Serve(lis); e != nil {
		panic(err)
	}
}

// func (s *LimitServer) mustEmbedUnimplementedLimitServerServer() {}
// func (s *LimitServer) PlaceLimitStream(a grpc.BidiStreamingServer[*proto.UplinkPlaceLimitRequest, *proto.UplinkPlaceLimitResponse]) error {
// 	return nil
// }

// // func (s *LimitServer) PlaceLimit(context.Context, *UplinkPlaceLimitRequest) (*UplinkPlaceLimitResponse, error) {
// // }

// func (s *LimitServer) PlaceLimit(_ context.Context, request *proto.UplinkPlaceLimitRequest) (*proto.UplinkPlaceLimitResponse, error) {
// 	bullionId, weight, price := request.GetBullionId(), request.GetWeight(), request.GetPrice()
// 	println(bullionId, weight, price)
// 	return &proto.UplinkPlaceLimitResponse{
// 		ReqId:   request.GetReqId(),
// 		Success: true,
// 		Message: "",
// 	}, nil
// }
