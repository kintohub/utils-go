package grpc

import (
	"context"
	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func unaryLoggingInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
	resp interface{}, err error) {
	logger := createRequestLogger(info.FullMethod)

	logger.Debug().Msg("...starting to process new grpc request")

	return handler(logger.WithContext(ctx), req)
}

func streamLoggingInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	logger := createRequestLogger(info.FullMethod)

	logger.Debug().Msg("...starting to process new grpc stream")

	return handler(srv, &grpc_middleware.WrappedServerStream{
		ServerStream:   ss,
		WrappedContext: logger.WithContext(ss.Context()),
	})
}

func createRequestLogger(requestName string) zerolog.Logger {
	return log.With().
		Caller(). // For all calls that do not have errors - simple stack trace
		Str("requestId", uuid.New().String()).
		Str("requestName", requestName).
		Logger()
}
