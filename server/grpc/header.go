package grpc

import (
	"context"
	"github.com/kintohub/utils-go/server/utils"
	"google.golang.org/grpc/metadata"
	"strings"
)

func GetAuthBearerTokenFromHeader(ctx context.Context) (string, *utils.Error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return "", utils.NewError(
			utils.StatusCode_InternalServerError, "could not parse grpc metadata from grpc context?")
	}

	const grpcAuthorizationHeaderKey = "authorization"
	authorizationArray := md.Get(grpcAuthorizationHeaderKey)
	arrLen := len(authorizationArray)

	// default empty for public requests
	token := ""
	if arrLen > 1 {
		return "", utils.NewError(utils.StatusCode_BadRequest,
			"invalid authorization metadata - can only have one authorization header!")
	} else if arrLen == 1 {
		const bearerTokenPrefix = "Bearer "
		token = strings.TrimPrefix(authorizationArray[0], bearerTokenPrefix)
	}

	return token, nil
}
