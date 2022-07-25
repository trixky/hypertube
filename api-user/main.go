package main

import (
	"log"

	"github.com/trixky/hypertube/.shared/environment"
	"github.com/trixky/hypertube/api-user/databases"
	"github.com/trixky/hypertube/api-user/internal"
)

func init() {
	log.Println("------------------------- INIT api-user")

	// Set environment config
	environment_config := environment.Config{
		ENV_grpc_port:         "API_USER_GRPC_PORT",
		ENV_grpc_gateway_port: "API_USER_GRPC_GATEWAY_PORT",
	}

	environment.Postgres.GetAll()                // Get postgres environment
	environment.Redis.GetAll()                   // Get redis environment
	environment.Grpc.GetAll(&environment_config) // Get grpc environment

	databases.InitDBs()       // Init DBs
	internal.NewGrpcServers() // Init internal servers
}

func main() {
	log.Println("------------------------- START api-user")
	select {} // Keep alive
}
