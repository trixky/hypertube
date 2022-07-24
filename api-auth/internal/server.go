package internal

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	md "github.com/trixky/hypertube/.shared/middlewares"
	"github.com/trixky/hypertube/.shared/utils"
	pb "github.com/trixky/hypertube/api-auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthServer struct {
	pb.AuthServiceServer
}

// NewGrpcServer create a new GRPC server
func NewGrpcServer(grpc_addr string) *grpc.Server {
	listen, err := net.Listen("tcp", grpc_addr)
	if err != nil {
		log.Fatalf("failed to serve grpc: %v\n", err)
	}

	s := grpc.NewServer()

	// Register the authentification service
	pb.RegisterAuthServiceServer(s, &AuthServer{})

	go func() {
		log.Fatalf("failed to serve grpc: %v\n", s.Serve(listen))
	}()

	return s
}

// NewGrpcGatewayServer create a new GRPC gateway server
func NewGrpcGatewayServer(grpc_gateway_addr string, grpc_addr string) {
	conn, err := grpc.DialContext(
		context.Background(),
		grpc_addr,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create mux
	gwmux := runtime.NewServeMux(runtime.WithMetadata(
		md.GrpcMiddleware,
	))

	// Register the authentification service
	err = pb.RegisterAuthServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatal(err)
	}

	// Create the HTTP sever
	gwServer := &http.Server{
		Addr:    grpc_gateway_addr,
		Handler: utils.AllowCORS(gwmux),
	}

	go func() {
		log.Fatalf("failed to serve grpc-gateway: %v\n", gwServer.ListenAndServe())
	}()
}
