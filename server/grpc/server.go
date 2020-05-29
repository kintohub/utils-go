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

type RegisterServiceHandler func(s *grpc.Server)

func RunServer(grpcPort, grpcWebPort, corsAllowedHosts string, handlers ...RegisterServiceHandler) {
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

	for _, handler := range handlers {
		handler(server)
	}

	go runGrpcServer(server, grpcPort)

	// This must go last as its a blocking call
	runGrpcWebServer(server, grpcWebPort, corsAllowedHosts)
}

func runGrpcServer(server *grpc.Server, port string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":"+port))
	if err != nil {
		log.Panic().Err(err).Msgf("failed to listen to tcp port.")
	}

	log.Info().Msgf("Listening to :%s for grpc connection requests", port)

	_ = server.Serve(lis)
}

func runGrpcWebServer(server *grpc.Server, port, coresAllowedHosts string) {
	wrappedGrpc := grpcweb.WrapServer(server,
		grpcweb.WithAllowedRequestHeaders([]string{coresAllowedHosts}),
		grpcweb.WithOriginFunc(func(origin string) bool {
			return true
		}))

	log.Info().Msgf("Listening to :%s for grpc-web connection requests", port)

	panic(http.ListenAndServe(":"+port, http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if wrappedGrpc.IsGrpcWebRequest(req) {
			wrappedGrpc.ServeHTTP(resp, req)
		}

		resp.Header().Set("Access-Control-Allow-Origin", coresAllowedHosts)
		resp.Header().Set("Access-Control-Allow-Headers", coresAllowedHosts)

		resp.WriteHeader(http.StatusOK)
	})))
}
