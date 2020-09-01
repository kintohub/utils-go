package grpc

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/kintohub/utils-go/klog"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

func ValidateGrpcRequest(v validation.Validatable) error {
	klog.Debugf("Request received - %v", v)

	if err := v.Validate(); err != nil {
		data, _ := json.Marshal(&err)
		klog.ErrorWithErr(err, "error during the validation of the request")
		return grpcStatus.Error(grpcCodes.InvalidArgument, string(data))
	}

	return nil
}
