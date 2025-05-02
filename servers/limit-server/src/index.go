package limit_server

import (
	"context"
	"io"
	"net"

	limitserver "github.com/rpsoftech/bullion-server/protocode/limit-server"
	interceptor "github.com/rpsoftech/bullion-server/servers/limit-server/src/Interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type LimitServerService struct {
	*limitserver.LimitServerServer
}

func (LimitServerService) mustEmbedUnimplementedLimitServerServer() {
	panic("unimplemented")
}

func Start() {
	lis, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.OAuth2Interceptor),
	// grpc.MaxConcurrentStreams(3),
	// grpc.
	// grpc.UnaryInterceptor()
	)
	// re
	limitserver.RegisterLimitServerServer(srv, &LimitServerService{})
	reflection.Register(srv)

	if e := srv.Serve(lis); e != nil {
		panic(e)
	}
}

func (s *LimitServerService) PlaceLimitStream(stream grpc.BidiStreamingServer[limitserver.UplinkPlaceLimitRequest, limitserver.UplinkPlaceLimitResponse]) error {
	// panic("unimplemented")
	for {
		request, err := stream.Recv()
		bullionId, weight, price := request.GetBullionId(), request.GetWeight(), request.GetPrice()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		println(bullionId, weight, price)
		stream.Send(&limitserver.UplinkPlaceLimitResponse{
			ReqId:   request.GetReqId(),
			Success: true,
			Message: "",
		})
		// for _, note := range s.routeNotes[key] {
		// 	if err := stream.Send(note); err != nil {
		// 		return err
		// 	}
		// }
	}
}

func (s *LimitServerService) PlaceOrderHedgingServer(stream grpc.BidiStreamingServer[limitserver.PlaceOrderHedgingResponse, limitserver.PlaceOrderHedgingRequest]) error {
	panic("unimplemented")
}

func (s *LimitServerService) PlaceLimit(_ context.Context, request *limitserver.UplinkPlaceLimitRequest) (*limitserver.UplinkPlaceLimitResponse, error) {
	bullionId, weight, price := request.GetBullionId(), request.GetWeight(), request.GetPrice()
	println(bullionId, weight, price)
	return &limitserver.UplinkPlaceLimitResponse{
		ReqId:   request.GetReqId(),
		Success: true,
		Message: "",
	}, nil
}
