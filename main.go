package main

import (
	"context"
	"fmt"
	"hwmonit/common"
	"hwmonit/network"
	"hwmonit/resource"
	"hwmonit/resource/base"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

var logger = common.GetLogger()

func main() {
	// config := common.GetConfig()

	rm := resource.NewResourceManager()
	for _, name := range base.AllResource() {
		v1, ok := os.LookupEnv(name)
		if ok == false {
			continue
		}
		interval, err := strconv.ParseInt(v1, 10, 10)
		if err != nil {
			continue
		}

		if interval > 0 {
			rm.AddResource(name, int(interval))
		}

	}
	rm.Start()
	config := common.GetConfig()
	address := ":" + config["port"]
	lis, errL := net.Listen("tcp", address)
	if errL != nil {
		logger.Fatal(fmt.Sprintf("Start serivce with port:%s failed", config["port"]))
	}
	httpServer := provideHTTP(address, rm)
	logger.Infof("address: %s，", address)

	if err := httpServer.Serve(lis); err != nil {
		log.Fatal("ListenAndServe: ", err)
		rm.Stop()
	}
}

func provideHTTP(endpoint string, srv network.HWMonitServer) *http.Server {
	grpcServer := grpc.NewServer()
	network.RegisterHWMonitServer(grpcServer, srv)
	ctx := context.Background()
	gwmux := runtime.NewServeMux()

	dopts := []grpc.DialOption{grpc.WithInsecure()}
	// register grpc-gateway pb
	if err := network.RegisterHWMonitHandlerFromEndpoint(ctx, gwmux, endpoint, dopts); err != nil {
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
