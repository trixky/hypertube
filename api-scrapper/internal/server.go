package internal

import (
	"log"
	"net"
	"strconv"

	"github.com/trixky/hypertube/.shared/environment"
	pb "github.com/trixky/hypertube/api-scrapper/proto"
	"google.golang.org/grpc"
)

type ScrapperServer struct {
	pb.ScrapperServiceServer
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

	pb.RegisterScrapperServiceServer(s, &ScrapperServer{})

	go func() {
		log.Fatalf("failed to serve grpc: %v\n", s.Serve(listen))
	}()

	return grpc_addr, s
}
