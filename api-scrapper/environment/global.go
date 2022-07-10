package environment

import (
	"log"
	"os"
)

func (e *Env) GetAll() {

	// --------- get GrpcPort
	if grpc_port, err := read_port(ENV_grpc_port); err != nil {
		log.Fatal(err)
	} else {
		e.GrpcPort = grpc_port
	}

	// --------- get GrpcGatewayPort
	if http_port, err := read_port(ENV_grpc_gateway_port); err != nil {
		log.Fatal(err)
	} else {
		e.GrpcGatewayPort = http_port
	}

	// --------- get TMDB api key
	if e.TmdbApiKey = os.Getenv(ENV_tmdb_api_key); len(e.TmdbApiKey) == 0 {
		log.Fatalf("%s environement variable missing", ENV_tmdb_api_key)
	}

	// --------- get PostgresPort
	if postgres_port, err := read_port(ENV_postgres_port); err != nil {
		log.Fatal(err)
	} else {
		e.PostgresPort = postgres_port
	}

	// --------- get PostgresHost
	if e.PostgresHost = os.Getenv(ENV_postgres_host); len(e.PostgresHost) == 0 {
		log.Fatalf("%s environement variable missing", ENV_postgres_host)
	}

	// --------- get PostgresUser
	if e.PostgresUser = os.Getenv(ENV_postgres_user); len(e.PostgresUser) == 0 {
		log.Fatalf("%s environement variable missing", ENV_postgres_user)
	}

	// --------- get PostgresPassword
	if e.PostgresPassword = os.Getenv(ENV_postgres_password); len(e.PostgresPassword) == 0 {
		log.Fatalf("%s environement variable missing", ENV_postgres_password)
	}

	// --------- get PostgresDB
	if e.PostgresDB = os.Getenv(ENV_postgres_db); len(e.PostgresDB) == 0 {
		log.Fatalf("%s environement variable missing", ENV_postgres_db)
	}
}
