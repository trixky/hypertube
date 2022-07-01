package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/trixky/hypertube/api-auth/postgres"
	pb "github.com/trixky/hypertube/api-auth/proto"
	"github.com/trixky/hypertube/api-auth/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	host            = "0.0.0.0"
	database_driver = "postgres"

	env_grpc_port         = "API_AUTH_GRPC_PORT"
	env_http_port         = "API_AUTH_HTTP_PORT"
	env_postgres_host     = "POSTGRES_HOST"
	env_postgres_port     = "POSTGRES_PORT"
	env_postgres_user     = "POSTGRES_USER"
	env_postgres_password = "POSTGRES_PASSWORD"
	env_postgres_db       = "POSTGRES_DB"
)

type Env struct {
	GrpcPort         int
	HttpPort         int
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
}

func (e *Env) GetAll() {
	const (
		port_max = 65535
		port_min = 1000
	)

	// --------- get GRPC Port
	grpc_port, err := strconv.Atoi(os.Getenv(env_grpc_port))

	if err != nil {
		log.Fatalf("%s environement variable is corrupted or missing", env_grpc_port)
	} else if grpc_port < port_min || grpc_port > port_max {
		log.Fatalf("port need to be included between %d and %d (%d)", port_min, port_max, grpc_port)
	}

	e.GrpcPort = grpc_port

	// --------- get HTTP Port
	http_port, err := strconv.Atoi(os.Getenv(env_http_port))

	if err != nil {
		log.Fatalf("%s environement variable is corrupted or missing", env_http_port)
	} else if http_port < port_min || http_port > port_max {
		log.Fatalf("port need to be included between %d and %d (%d)", port_min, port_max, http_port)
	}

	e.HttpPort = http_port

	// --------- get PostgresHost
	if e.PostgresHost = os.Getenv(env_postgres_host); len(e.PostgresHost) == 0 {
		log.Fatalf("%s environement variable missing", env_postgres_host)
	}

	// --------- get PostgresPort
	postgres_port, err := strconv.Atoi(os.Getenv(env_postgres_port))

	if err != nil {
		log.Fatalf("%s environement variable is corrupted or missing", env_postgres_port)
	} else if postgres_port < port_min || postgres_port > port_max {
		log.Fatalf("port need to be included between %d and %d (%d)", port_min, port_max, postgres_port)
	}

	e.PostgresPort = postgres_port

	// --------- get PostgresUser
	if e.PostgresUser = os.Getenv(env_postgres_user); len(e.PostgresUser) == 0 {
		log.Fatalf("%s environement variable missing", env_postgres_user)
	}

	// --------- get PostgresPassword
	if e.PostgresPassword = os.Getenv(env_postgres_password); len(e.PostgresPassword) == 0 {
		log.Fatalf("%s environement variable missing", env_postgres_password)
	}

	// --------- get PostgresDB
	if e.PostgresDB = os.Getenv(env_postgres_db); len(e.PostgresDB) == 0 {
		log.Fatalf("%s environement variable missing", env_postgres_db)
	}
}

func readEnv() (env Env) {
	env.GetAll()

	return
}

func main() {
	env := readEnv()

	grpc_addr := host + ":" + strconv.Itoa(env.GrpcPort)

	if err := postgres.Init(postgres.Config{
		Driver:   database_driver,
		Host:     env.PostgresHost,
		Port:     env.PostgresPort,
		User:     env.PostgresUser,
		Password: env.PostgresPassword,
		Dbname:   env.PostgresDB,
	}); err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	log.Println("connected to postgres")

	listen, err := net.Listen("tcp", grpc_addr)

	if err != nil {
		log.Fatalf("failed to listen on: %v\n", err)
	}

	log.Printf("listening on %s\n", grpc_addr)

	s := grpc.NewServer()

	pb.RegisterAuthServiceServer(s, &server.AuthServer{})

	go func() {
		log.Fatalf("failed to serve on: %v\n", s.Serve(listen))
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		grpc_addr,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	http_addr := ":" + strconv.Itoa(env.HttpPort)

	gwmux := runtime.NewServeMux()
	err = pb.RegisterAuthServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    http_addr,
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090 /", http_addr)
	log.Fatalln(gwServer.ListenAndServe())
}
