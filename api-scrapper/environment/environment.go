package environment

import (
	"fmt"
	"os"
	"strconv"
)

const (
	ENV_grpc_port         = "API_SCRAPPER_GRPC_PORT"
	ENV_grpc_gateway_port = "API_SCRAPPER_GRPC_GATEWAY_PORT"
	ENV_tmdb_api_key      = "TMDB_API_KEY"
	ENV_postgres_host     = "POSTGRES_HOST"
	ENV_postgres_port     = "POSTGRES_PORT"
	ENV_postgres_user     = "POSTGRES_USER"
	ENV_postgres_password = "POSTGRES_PASSWORD"
	ENV_postgres_db       = "POSTGRES_DB"
)

type Env struct {
	GrpcPort         int
	GrpcGatewayPort  int
	TmdbApiKey       string
	HttpPort         int
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
}

var E = Env{}

func read_port(name string) (int, error) {
	const (
		port_max = 65535
		port_min = 1000
	)

	port, err := strconv.Atoi(os.Getenv(name))

	if err != nil {
		return 0, fmt.Errorf("%s environement variable is corrupted or missing", name)
	} else if port < port_min || port > port_max {
		return 0, fmt.Errorf("port need to be included between %d and %d (%d)", port_min, port_max, port)
	}

	return port, nil
}
