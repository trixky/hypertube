package main

import (
	"log"

	"github.com/trixky/hypertube/.shared/databases"
	"github.com/trixky/hypertube/.shared/environment"
	"github.com/trixky/hypertube/api-media/external"
	"github.com/trixky/hypertube/api-media/internal"
	"github.com/trixky/hypertube/api-media/queries"
)

func init() {
	log.Println("------------------------- INIT api-media")

	// Set environment config
	environment_config := environment.Config{
		ENV_grpc_port:         "API_MEDIA_GRPC_PORT",
		ENV_grpc_gateway_port: "API_MEDIA_GRPC_GATEWAY_PORT",
	}

	environment.Postgres.GetAll()                // Get postgres environment
	environment.Redis.GetAll()                   // Get redis environment
	environment.Grpc.GetAll(&environment_config) // Get grpc environment
	external.Scrapper.GetAll()                   // Get scrapper environment

	external.NewApiScrapperClient() // Init Scrapper client
	databases.InitPostgres()        // Init DBs
	databases.InitRedis()
	queries.InitSqlc()        // Init Sqlc queries
	internal.NewGrpcServers() // Init internal servers
}

func main() {
	log.Println("------------------------- START api-media")
	select {} // Keep alive
}
