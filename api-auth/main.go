package main

import (
	"log"
	"strconv"

	"github.com/trixky/hypertube/api-auth/databases"
	"github.com/trixky/hypertube/api-auth/environment"
	"github.com/trixky/hypertube/api-auth/external"
	"github.com/trixky/hypertube/api-auth/internal"
)

const (
	host            = "0.0.0.0"
	postgres_driver = "postgres"
)

func main() {
	environment.E.GetAll()
	// ------------- postgres
	log.Printf("start connection to postgres on %s:%d\n", environment.E.PostgresHost, environment.E.PostgresPort)
	if err := databases.InitPosgres(databases.PostgresConfig{
		Driver:   postgres_driver,
		Host:     environment.E.PostgresHost,
		Port:     environment.E.PostgresPort,
		User:     environment.E.PostgresUser,
		Password: environment.E.PostgresPassword,
		Dbname:   environment.E.PostgresDB,
	}); err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	// ------------- redis
	log.Println("start connection to redis on default address")
	if err := databases.InitRedis(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	// ------------- grpc
	grpc_addr := host + ":" + strconv.Itoa(environment.E.GrpcPort)

	go func() {
		log.Fatalf("failed to serve grpc on: %v\n", internal.NewGrpcServer(grpc_addr))
	}()

	// ------------- grpc-gateway
	grpc_gateway_addr := ":" + strconv.Itoa(environment.E.GrpcGatewayPort)

	go func() {
		log.Fatalf("failed to serve grpc-gateway on: %v\n", internal.NewGrpcGatewayServer(grpc_gateway_addr, grpc_addr))
	}()

	// ------------- http
	http_addr := ":" + strconv.Itoa(environment.E.HttpPort)

	log.Printf("start to serve http services on %s\n", http_addr)
	log.Fatalf("failed to serve http services on %s: %v", http_addr, external.NewHttpServer(http_addr))
}
