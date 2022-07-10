package internal

import (
	"log"
	"net"

	pb "github.com/trixky/hypertube/api-scrapper/proto"
	"google.golang.org/grpc"
)

type ScrapperServer struct {
	pb.ScrapperServiceServer
}

func NewGrpcServer(grpc_addr string) error {
	listen, err := net.Listen("tcp", grpc_addr)
	if err != nil {
		return err
	}

	s := grpc.NewServer()

	pb.RegisterScrapperServiceServer(s, &ScrapperServer{})

	log.Printf("start to serve grpc services on %s\n", grpc_addr)

	return s.Serve(listen)
}
