package grpc

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rs/zerolog/log"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

func ValidateGrpcRequest(v validation.Validatable) error {
	if err := v.Validate(); err != nil {
		data, _ := json.Marshal(&err)
		log.Error().Err(err).Msgf("error during the validation of the request")
		return grpcStatus.Error(grpcCodes.InvalidArgument, string(data))
	}

	log.Debug().Msgf("Request received - %v", v)

	return nil
}
