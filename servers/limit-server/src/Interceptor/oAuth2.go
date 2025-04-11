package interceptor

import (
	"context"
	"fmt"
	"strings"

	"github.com/rpsoftech/bullion-server/src/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

// func(ctx context.Context, req any, info *UnaryServerInfo, handler UnaryHandler) (resp any, err error)

func OAuth2Interceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

	// return nil, nil
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	if !valid(md[env.RequestTokenHeaderKey]) {
		return nil, errInvalidToken
	}
	m, err := handler(ctx, req)
	if err != nil {
		logger("RPC failed with error: %v", err)
	}
	return m, err
}
func logger(format string, a ...any) {
	fmt.Printf("LOG:\t"+format+"\n", a...)
}

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	println(authorization)
	// Perform the token validation here. For the sake of this example, the code
	// here forgoes any of the usual OAuth2 token validation and instead checks
	// for a token matching an arbitrary string.
	return token == "some-secret-token"
}
