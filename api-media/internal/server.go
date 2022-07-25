package internal

import (
	"context"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/trixky/hypertube/.shared/environment"
	md "github.com/trixky/hypertube/.shared/middlewares"
	sutils "github.com/trixky/hypertube/.shared/utils"
	pb "github.com/trixky/hypertube/api-media/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MediaServer struct {
	pb.MediaServiceServer
}

// NewGrpcServer create a new GRPC server
func NewGrpcServer() (string, *grpc.Server) {
	grpc_port := ":" + strconv.Itoa(environment.Grpc.Port)
	grpc_addr := environment.DEFAULT_HOST + grpc_port

	log.Printf("start to serve grpc services on \t\t%s\n", grpc_addr)

	listen, err := net.Listen("tcp", grpc_addr)
	if err != nil {
		log.Fatalf("failed to serve grpc: %v\n", err)
	}

	s := grpc.NewServer()

	pb.RegisterMediaServiceServer(s, &MediaServer{})

	go func() {
		log.Fatalf("failed to serve grpc: %v\n", s.Serve(listen))
	}()

	return grpc_addr, s
}

// NewGrpcGatewayServer create a new GRPC gateway server
func NewGrpcGatewayServer(grpc_addr string) {
	grpc_gateway_addr := ":" + strconv.Itoa(environment.Grpc.GatewayPort)

	log.Printf("start to serve grpc gateway services on \t%s\n", grpc_gateway_addr)

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
	err = pb.RegisterMediaServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatal(err)
	}

	// Create the HTTP sever
	gwServer := &http.Server{
		Addr:    grpc_gateway_addr,
		Handler: sutils.AllowCORS(gwmux),
	}

	go func() {
		log.Fatalf("failed to serve grpc-gateway: %v\n", gwServer.ListenAndServe())
	}()
}

// NewGrpcServers create all GRPC servers
func NewGrpcServers() *grpc.Server {
	grpc_addr, grpc_server := NewGrpcServer() // GRPC
	NewGrpcGatewayServer(grpc_addr)           // GRPC GATEWAY

	return grpc_server
}
