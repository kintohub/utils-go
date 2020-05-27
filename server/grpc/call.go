package grpc

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

const (
	contextMethodNameKey = "method"
)

// Currently enriches grpc call with method name only atm
func unaryEnrichCallInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
	resp interface{}, err error) {

	return handler(context.WithValue(ctx, contextMethodNameKey, info.FullMethod), req)
}

// Currently enriches grpc call with method name only atm
func streamEnrichCallInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {

	return handler(srv, &grpc_middleware.WrappedServerStream{
		ServerStream:   ss,
		WrappedContext: context.WithValue(ss.Context(), contextMethodNameKey, info.FullMethod),
	})
}
