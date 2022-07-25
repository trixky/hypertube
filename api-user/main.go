package main

import (
	"log"
	"strconv"

	"github.com/trixky/hypertube/.shared/environment"
	"github.com/trixky/hypertube/api-user/databases"
	"github.com/trixky/hypertube/api-user/internal"
)

const (
	host            = "0.0.0.0"
	postgres_driver = "postgres"
)

func main() {
	environment.ReadAll()
	// ------------- postgres
	log.Printf("start connection to postgres on %s:%d\n", environment.Postgres.PostgresHost, environment.Postgres.PostgresPort)
	if err := databases.InitPosgres(databases.PostgresConfig{
		Driver:   postgres_driver,
		Host:     environment.Postgres.PostgresHost,
		Port:     environment.Postgres.PostgresPort,
		User:     environment.Postgres.PostgresUser,
		Password: environment.Postgres.PostgresPassword,
		Dbname:   environment.Postgres.PostgresDB,
	}); err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	// ------------- redis
	log.Println("start connection to redis on default address")
	if err := databases.InitRedis(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	// ------------- grpc
	grpc_addr := host + ":" + strconv.Itoa(environment.Grpc.GrpcPort)

	go func() {
		log.Fatalf("failed to serve grpc on: %v\n", internal.NewGrpcServer(grpc_addr))
	}()

	// ------------- grpc-gateway
	grpc_gateway_addr := ":" + strconv.Itoa(environment.Grpc.GrpcGatewayPort)

	log.Fatalf("failed to serve grpc-gateway on: %v\n", internal.NewGrpcGatewayServer(grpc_gateway_addr, grpc_addr))
}
