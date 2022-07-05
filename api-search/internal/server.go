package internal

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/trixky/hypertube/api-search/proto"
	"github.com/trixky/hypertube/api-search/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SearchServer struct {
	pb.SearchServiceServer
}

func NewGrpcServer(grpc_addr string) error {
	listen, err := net.Listen("tcp", grpc_addr)
	if err != nil {
		return err
	}

	s := grpc.NewServer()

	pb.RegisterSearchServiceServer(s, &SearchServer{})

	log.Printf("start to serve grpc services on %s\n", grpc_addr)

	return s.Serve(listen)
}

func NewGrpcGatewayServer(grpc_gateway_addr string, grpc_addr string) error {
	conn, err := grpc.DialContext(
		context.Background(),
		grpc_addr,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	gwmux := runtime.NewServeMux(runtime.WithMetadata(
		basic_middleware,
	))

	err = pb.RegisterSearchServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		return err
	}

	gwServer := &http.Server{
		Addr:    grpc_gateway_addr,
		Handler: utils.AllowCORS(gwmux),
	}

	log.Printf("start to serve grpc-gateway services on %s\n", grpc_gateway_addr)

	return gwServer.ListenAndServe()
}
