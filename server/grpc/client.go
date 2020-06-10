package grpc

import (
	"crypto/x509"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func getDialOptionSecurity(isTLS bool) grpc.DialOption {
	dialOption := grpc.WithInsecure()

	if isTLS {
		// https://grpc.io/docs/guides/auth/#authenticate-with-google
		pool, _ := x509.SystemCertPool()
		creds := credentials.NewClientTLSFromCert(pool, "")
		dialOption = grpc.WithTransportCredentials(creds)

	}

	return dialOption
}

func CreateConnectionOrDie(host string, isTLS bool) *grpc.ClientConn {
	conn, err := grpc.Dial(host, getDialOptionSecurity(isTLS))

	if err != nil {
		log.Panic().Msgf("could not create grpc connection to %v - %v", host, err)
	}

	return conn
}

func CreateConnectionWithMaxMsgSizeOrDie(host string, isTLS bool, maxMsgSizeInBytes int) *grpc.ClientConn {
	conn, err := grpc.Dial(
		host,
		getDialOptionSecurity(isTLS),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSizeInBytes)),
	)

	if err != nil {
		log.Panic().Msgf("could not create grpc connection to %v - %v", host, err)
	}

	return conn
}
