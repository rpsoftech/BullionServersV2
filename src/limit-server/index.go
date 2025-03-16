package limit_server

import (
	"context"
	"io"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type LimitServerService struct{}

func (LimitServerService) mustEmbedUnimplementedLimitServerServer() {
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
	// re
	RegisterLimitServerServer(srv, &LimitServerService{})
	reflection.Register(srv)

	if e := srv.Serve(lis); e != nil {
		panic(err)
	}
}

func (s *LimitServerService) PlaceLimitStream(stream grpc.BidiStreamingServer[UplinkPlaceLimitRequest, UplinkPlaceLimitResponse]) error {
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
		stream.Send(&UplinkPlaceLimitResponse{
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

func (s *LimitServerService) PlaceOrderHedgingServer(stream grpc.BidiStreamingServer[PlaceOrderHedgingResponse, PlaceOrderHedgingRequest]) error {
	panic("unimplemented")
}

func (s *LimitServerService) PlaceLimit(_ context.Context, request *UplinkPlaceLimitRequest) (*UplinkPlaceLimitResponse, error) {
	bullionId, weight, price := request.GetBullionId(), request.GetWeight(), request.GetPrice()
	println(bullionId, weight, price)
	return &UplinkPlaceLimitResponse{
		ReqId:   request.GetReqId(),
		Success: true,
		Message: "",
	}, nil
}
