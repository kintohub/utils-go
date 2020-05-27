package grpc

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

func RunServer(grpcPort, grpcWebPort string, registerServicesFunc func(server *grpc.Server)) {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			unaryEnrichCallInterceptor,
			unaryLoggingInterceptor,
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(panicRecoveryHandler)),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			streamEnrichCallInterceptor,
			streamLoggingInterceptor,
			grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(panicRecoveryHandler)),
		)),
	)

	registerServicesFunc(server)

	go runGrpcServer(server, grpcPort)

	// This must go last as its a blocking call
	runGrpcWebServer(server, grpcWebPort)
}

func runGrpcServer(server *grpc.Server, port string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":"+port))
	if err != nil {
		log.Panic().Err(err).Msgf("failed to listen to tcp port.")
	}

	log.Info().Msgf("Listening to :%s for grpc connection requests", port)

	_ = server.Serve(lis)
}

func runGrpcWebServer(server *grpc.Server, port string) {
	wrappedGrpc := grpcweb.WrapServer(server,
		grpcweb.WithAllowedRequestHeaders([]string{"*"}),
		grpcweb.WithOriginFunc(func(origin string) bool {
			return true
		}))

	log.Info().Msgf("Listening to :%s for grpc-web connection requests", port)

	panic(http.ListenAndServe(":"+port, http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if wrappedGrpc.IsGrpcWebRequest(req) {
			wrappedGrpc.ServeHTTP(resp, req)
		}

		resp.Header().Set("Access-Control-Allow-Origin", "*")
		resp.Header().Set("Access-Control-Allow-Headers", "*")

		resp.WriteHeader(http.StatusOK)
	})))
}
