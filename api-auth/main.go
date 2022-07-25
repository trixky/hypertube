package main

import (
	"log"

	"github.com/trixky/hypertube/.shared/environment"
	"github.com/trixky/hypertube/api-auth/databases"
	"github.com/trixky/hypertube/api-auth/external"
	"github.com/trixky/hypertube/api-auth/internal"
)

func init() {
	log.Println("------------------------- INIT api-auth")

	// Set environment config
	environment_config := environment.Config{
		ENV_grpc_port:         "API_AUTH_GRPC_PORT",
		ENV_grpc_gateway_port: "API_AUTH_GRPC_GATEWAY_PORT",
		ENV_http_port:         "API_AUTH_HTTP_PORT",
	}

	environment.Postgres.GetAll()                // Get postgres environment
	environment.Redis.GetAll()                   // Get redis environment
	environment.Grpc.GetAll(&environment_config) // Get grpc environment
	environment.Http.GetAll(&environment_config) // Get http environment
	environment.Api42.GetAll()                   // Get 42 api environment
	environment.ApiGoogle.GetAll()               // Get google api environment
	environment.Outlook.GetAll()                 // Get outlook environment

	databases.InitDBs()       // Init DBs
	internal.NewGrpcServers() // Init internal servers
	external.NewHttpServer()  // Init external servers
}

func main() {
	log.Println("------------------------- START api-auth")
	select {} // Keep alive
}
