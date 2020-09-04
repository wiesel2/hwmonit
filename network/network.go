package network

import (
	context "context"
	"fmt"
	"hwmonit/common"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	grpc "google.golang.org/grpc"
)

type Server struct {
}

var logger = common.GetLogger()

func (s *Server) Serve(address string, rm HWMonitServer) error {
	lis, errL := net.Listen("tcp", address)
	if errL != nil {
		logger.Fatal(fmt.Sprintf("Start serivce with port:%s failed", address))
	}
	httpServer := provideHTTP(address, rm)
	logger.Infof("address: %s，", address)

	err := httpServer.Serve(lis)
	return err
}

func provideHTTP(endpoint string, srv HWMonitServer) *http.Server {
	grpcServer := grpc.NewServer()
	RegisterHWMonitServer(grpcServer, srv)
	ctx := context.Background()
	gwmux := runtime.NewServeMux()

	dopts := []grpc.DialOption{grpc.WithInsecure()}
	// register grpc-gateway pb
	if err := RegisterHWMonitHandlerFromEndpoint(ctx, gwmux, endpoint, dopts); err != nil {
		log.Printf("Failed to register gw server: %v\n", err)
	}

	// http服务
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	return &http.Server{
		Addr:    endpoint,
		Handler: grpcHandlerFunc(grpcServer, mux),
	}
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
