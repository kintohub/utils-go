package grpc

import (
	"context"
	"github.com/kintohub/utils-go/server"
	"github.com/rs/zerolog/log"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// Used here and inside logger.go
const seriousErrorMsg = "a serious error occurred. if the issue persists, please contact the site administrator."

func ConvertToGrpcError(ctx context.Context, error *server.Error) error {
	if error.StatusCode >= server.StatusCode_InternalServerError {
		log.Ctx(ctx).
			Error().
			Stack().
			Err(error.Error).
			Interface("statusCode", error.StatusCode).
			Msg(error.Message)

		// modify the true error message for the client to just be a serious error vs system error
		error.Message = seriousErrorMsg
	} else {
		log.Ctx(ctx).
			Error().
			Err(error.Error).
			Interface("statusCode", error.StatusCode).
			Msg(error.Message)
	}

	return grpcStatus.Error(convertStatusCodeToGrpcCode(error.StatusCode), error.Message)
}

func panicRecoveryHandler(v interface{}) error {
	log.Error().Stack().Msgf("[IMPORTANT] a uncaught panic occurred: %v", v)
	return grpcStatus.Error(grpcCodes.Internal, seriousErrorMsg)
}

// Referenced from https://github.com/grpc-ecosystem/grpc-gateway/blob/master/runtime/errors.go
func convertStatusCodeToGrpcCode(code server.StatusCode) grpcCodes.Code {
	switch code {
	case server.StatusCode_OK:
		return grpcCodes.OK
	case server.StatusCode_RequestTimeout:
		return grpcCodes.Canceled
	case server.StatusCode_InternalServerError:
		return grpcCodes.Internal
	case server.StatusCode_BadRequest:
		return grpcCodes.InvalidArgument
	case server.StatusCode_GatewayTimeout:
		return grpcCodes.DeadlineExceeded
	case server.StatusCode_NotFound:
		return grpcCodes.NotFound
	case server.StatusCode_Conflict:
		return grpcCodes.AlreadyExists
	case server.StatusCode_Forbidden:
		return grpcCodes.PermissionDenied
	case server.StatusCode_Unauthorized:
		return grpcCodes.Unauthenticated
	case server.StatusCode_TooManyRequests:
		return grpcCodes.ResourceExhausted
	case server.StatusCode_NotImplemented:
		return grpcCodes.Unimplemented
	case server.StatusCode_ServiceUnavailable:
		return grpcCodes.Unavailable
	default:
		log.Warn().Msgf("unsupported kinto error code %v to grpc translation.", code)
		return grpcCodes.Unknown
	}
}
