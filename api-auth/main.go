package main

import (
	"log"
	"net"

	pb "github.com/trixky/hypertube/api/proto"
	"github.com/trixky/hypertube/api/server"
	"google.golang.org/grpc"
)

func main() {
	addr := "0.0.0.0:5051"

	listen, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Failed to listen on: %v\n", err)
	}

	log.Printf("Listening on %s\n", addr)

	s := grpc.NewServer()

	pb.RegisterAuthServiceServer(s, &server.AuthServer{})

	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve on: %v\n", err)
	}
}
