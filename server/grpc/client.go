package grpc

import (
	"crypto/x509"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func CreateConnectionOrDie(host string, isTLS bool) *grpc.ClientConn {
	dialOption := grpc.WithInsecure()

	if isTLS {
		// https://grpc.io/docs/guides/auth/#authenticate-with-google
		pool, _ := x509.SystemCertPool()
		creds := credentials.NewClientTLSFromCert(pool, "")
		dialOption = grpc.WithTransportCredentials(creds)
	}

	conn, err := grpc.Dial(host, dialOption)

	if err != nil {
		log.Panic().Msgf("could not create grpc connection to %v - %v", host, err)
	}

	return conn
}
