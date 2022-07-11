package environment

import (
	"fmt"
	"os"
	"strconv"
)

// https://github.com/grpc-ecosystem/grpc-gateway/blob/f046a4ebdc9be76e11c6239eaeba4e30e9e2444e/docs/docs/faq.md
// https://github.com/grpc-ecosystem/grpc-gateway/blob/master/examples/internal/gateway/main.go
// https://github.com/grpc-ecosystem/grpc-gateway/blob/0ebdfba80649c56b0da0777376c970d17c3c9540/examples/internal/gateway/handlers.go#L32

const (
	ENV_grpc_port         = "API_AUTH_GRPC_PORT"
	ENV_grpc_gateway_port = "API_AUTH_GRPC_GATEWAY_PORT"
	ENV_http_port         = "API_AUTH_HTTP_PORT"
	ENV_postgres_host     = "POSTGRES_HOST"
	ENV_postgres_port     = "POSTGRES_PORT"
	ENV_postgres_user     = "POSTGRES_USER"
	ENV_postgres_password = "POSTGRES_PASSWORD"
	ENV_postgres_db       = "POSTGRES_DB"
)

type Env struct {
	GrpcPort         int
	GrpcGatewayPort  int
	HttpPort         int
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	API42            Api42
	APIGoogle        ApiGoogle
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
