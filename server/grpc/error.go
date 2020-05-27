package grpc

import (
	"context"
	"github.com/kintohub/utils-go/server/utils"
	"github.com/rs/zerolog/log"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// Used here and inside logger.go
const seriousErrorMsg = "a serious error occurred. if the issue persists, please contact the site administrator."

func ConvertToGrpcError(ctx context.Context, error *utils.Error) error {
	if error.StatusCode >= utils.StatusCode_InternalServerError {
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
			Debug().
			Err(error.Error).
			Interface("statusCode", error.StatusCode).
			Msg(error.Message)
	}

	return grpcStatus.Error(convertStatusCodeToGrpcCode(error.StatusCode), error.Message)
}

func panicRecoveryHandler(v interface{}) error {
	log.Error().Msgf("[IMPORTANT] a uncaught panic occurred: %v", v)
	return grpcStatus.Error(grpcCodes.Internal, seriousErrorMsg)
}

// Referenced from https://github.com/grpc-ecosystem/grpc-gateway/blob/master/runtime/errors.go
func convertStatusCodeToGrpcCode(code utils.StatusCode) grpcCodes.Code {
	switch code {
	case utils.StatusCode_OK:
		return grpcCodes.OK
	case utils.StatusCode_RequestTimeout:
		return grpcCodes.Canceled
	case utils.StatusCode_InternalServerError:
		return grpcCodes.Internal
	case utils.StatusCode_BadRequest:
		return grpcCodes.InvalidArgument
	case utils.StatusCode_GatewayTimeout:
		return grpcCodes.DeadlineExceeded
	case utils.StatusCode_NotFound:
		return grpcCodes.NotFound
	case utils.StatusCode_Conflict:
		return grpcCodes.AlreadyExists
	case utils.StatusCode_Forbidden:
		return grpcCodes.PermissionDenied
	case utils.StatusCode_Unauthorized:
		return grpcCodes.Unauthenticated
	case utils.StatusCode_TooManyRequests:
		return grpcCodes.ResourceExhausted
	case utils.StatusCode_NotImplemented:
		return grpcCodes.Unimplemented
	case utils.StatusCode_ServiceUnavailable:
		return grpcCodes.Unavailable
	default:
		log.Warn().Msgf("unsupported kinto error code %v to grpc translation.", code)
		return grpcCodes.Unknown
	}
}
