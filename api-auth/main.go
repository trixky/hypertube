package main

import (
	"log"
	"net"
	"os"
	"strconv"

	pb "github.com/trixky/hypertube/api/proto"
	"github.com/trixky/hypertube/api/server"
	"google.golang.org/grpc"
)

const (
	env_port = "API_AUTH_PORT"
)

type Env struct {
	Port int
}

func (e *Env) GetAll() {
	const (
		port_max = 65535
		port_min = 1000
	)

	port, err := strconv.Atoi(os.Getenv(env_port))

	if err != nil {
		log.Fatalf("%s environement variable is corrupted or missing", env_port)
	}

	if port < port_min || port > port_max {
		log.Fatalf("port need to be included between %d and %d (%d)", port_min, port_max, port)
	}

	e.Port = port
}

func readEnv() (env Env) {
	env.GetAll()

	return
}

func main() {
	env := readEnv()

	addr := "0.0.0.0:" + strconv.Itoa(env.Port)

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
