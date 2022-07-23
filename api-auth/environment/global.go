package environment

import (
	"log"
	"os"
)

// GetAll read all needed enviornment variables
func (e *Env) GetAll() {
	// --------- Get GrpcPort
	if grpc_port, err := read_port(ENV_grpc_port); err != nil {
		log.Fatal(err)
	} else {
		e.GrpcPort = grpc_port
	}

	// --------- Get GrpcGatewayPort
	if http_port, err := read_port(ENV_grpc_gateway_port); err != nil {
		log.Fatal(err)
	} else {
		e.GrpcGatewayPort = http_port
	}

	// --------- Get HttpPort
	if redirect_port, err := read_port(ENV_http_port); err != nil {
		log.Fatal(err)
	} else {
		e.HttpPort = redirect_port
	}

	// --------- Get PostgresPort
	if postgres_port, err := read_port(ENV_postgres_port); err != nil {
		log.Fatal(err)
	} else {
		e.PostgresPort = postgres_port
	}

	// --------- Get PostgresHost
	if e.PostgresHost = os.Getenv(ENV_postgres_host); len(e.PostgresHost) == 0 {
		log.Fatalf("%s %s", ENV_postgres_host, environement_variable_missing)
	}

	// --------- Get PostgresUser
	if e.PostgresUser = os.Getenv(ENV_postgres_user); len(e.PostgresUser) == 0 {
		log.Fatalf("%s %s", ENV_postgres_user, environement_variable_missing)
	}

	// --------- Get PostgresPassword
	if e.PostgresPassword = os.Getenv(ENV_postgres_password); len(e.PostgresPassword) == 0 {
		log.Fatalf("%s %s", ENV_postgres_password, environement_variable_missing)
	}

	// --------- Get PostgresDB
	if e.PostgresDB = os.Getenv(ENV_postgres_db); len(e.PostgresDB) == 0 {
		log.Fatalf("%s %s", ENV_postgres_db, environement_variable_missing)
	}

	e.API42.GetAll()         // Get 42 environment variables
	e.APIGoogle.GetAll()     // Get google environment variables
	e.OUTLOOKConfig.GetAll() // Get outlook environment variables
}
