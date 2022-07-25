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
	Host     string
	Port     int
	User     string
	Password string
	DBname   string
}

// GetAll read all needed enviornment variables
func (e *env_postgres) GetAll() {
	// --------- Get PostgresPort
	if postgres_port, err := ReadPort(ENV_postgres_port); err != nil {
		log.Fatal(err)
	} else {
		e.Port = postgres_port
	}

	// --------- Get Host
	if e.Host = os.Getenv(ENV_postgres_host); len(e.Host) == 0 {
		log.Fatalf("%s %s", ENV_postgres_host, environement_variable_missing)
	}

	// --------- Get User
	if e.User = os.Getenv(ENV_postgres_user); len(e.User) == 0 {
		log.Fatalf("%s %s", ENV_postgres_user, environement_variable_missing)
	}

	// --------- Get Password
	if e.Password = os.Getenv(ENV_postgres_password); len(e.Password) == 0 {
		log.Fatalf("%s %s", ENV_postgres_password, environement_variable_missing)
	}

	// --------- Get DBname
	if e.DBname = os.Getenv(ENV_postgres_db); len(e.DBname) == 0 {
		log.Fatalf("%s %s", ENV_postgres_db, environement_variable_missing)
	}
}

var Postgres = env_postgres{}
