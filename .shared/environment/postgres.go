package environment

import (
	"log"
	"os"
)

const (
	ENV_postgres_host     = "POSTGRES_HOST"
	ENV_postgres_port     = "POSTGRES_PORT"
	ENV_postgres_user     = "POSTGRES_USER"
	ENV_postgres_password = "POSTGRES_PASSWORD"
	ENV_postgres_db       = "POSTGRES_DB"
)

type env_postgres struct {
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
}

// GetAll read all needed enviornment variables
func (e *env_postgres) GetAll() {
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
}

var Postgres = env_postgres{}
