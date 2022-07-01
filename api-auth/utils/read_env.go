package utils

import (
	"log"
	"os"
	"strconv"
)

// https://github.com/grpc-ecosystem/grpc-gateway/blob/f046a4ebdc9be76e11c6239eaeba4e30e9e2444e/docs/docs/faq.md
// https://github.com/grpc-ecosystem/grpc-gateway/blob/master/examples/internal/gateway/main.go
// https://github.com/grpc-ecosystem/grpc-gateway/blob/0ebdfba80649c56b0da0777376c970d17c3c9540/examples/internal/gateway/handlers.go#L32

const (
	ENV_grpc_port         = "API_AUTH_GRPC_PORT"
	ENV_http_port         = "API_AUTH_HTTP_PORT"
	ENV_postgres_host     = "POSTGRES_HOST"
	ENV_postgres_port     = "POSTGRES_PORT"
	ENV_postgres_user     = "POSTGRES_USER"
	ENV_postgres_password = "POSTGRES_PASSWORD"
	ENV_postgres_db       = "POSTGRES_DB"
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

func (e *Env) getAll() {
	const (
		port_max = 65535
		port_min = 1000
	)

	// --------- get GRPC Port
	grpc_port, err := strconv.Atoi(os.Getenv(ENV_grpc_port))

	if err != nil {
		log.Fatalf("%s environement variable is corrupted or missing", ENV_grpc_port)
	} else if grpc_port < port_min || grpc_port > port_max {
		log.Fatalf("port need to be included between %d and %d (%d)", port_min, port_max, grpc_port)
	}

	e.GrpcPort = grpc_port

	// --------- get HTTP Port
	http_port, err := strconv.Atoi(os.Getenv(ENV_http_port))

	if err != nil {
		log.Fatalf("%s environement variable is corrupted or missing", ENV_http_port)
	} else if http_port < port_min || http_port > port_max {
		log.Fatalf("port need to be included between %d and %d (%d)", port_min, port_max, http_port)
	}

	e.HttpPort = http_port

	// --------- get PostgresHost
	if e.PostgresHost = os.Getenv(ENV_postgres_host); len(e.PostgresHost) == 0 {
		log.Fatalf("%s environement variable missing", ENV_postgres_host)
	}

	// --------- get PostgresPort
	postgres_port, err := strconv.Atoi(os.Getenv(ENV_postgres_port))

	if err != nil {
		log.Fatalf("%s environement variable is corrupted or missing", ENV_postgres_port)
	} else if postgres_port < port_min || postgres_port > port_max {
		log.Fatalf("port need to be included between %d and %d (%d)", port_min, port_max, postgres_port)
	}

	e.PostgresPort = postgres_port

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

func ReadEnv() (env Env) {
	env.getAll()

	return
}
