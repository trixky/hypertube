package main

import (
	"log"
	"os"
	"strconv"

	importer "github.com/trixky/hypertube/scripts/import-medias/importer"
	"github.com/trixky/hypertube/scripts/import-medias/postgres"
)

const (
	host            = "0.0.0.0"
	database_driver = "pgx"

	env_postgres_host     = "POSTGRES_HOST"
	env_postgres_port     = "POSTGRES_PORT"
	env_postgres_user     = "POSTGRES_USER"
	env_postgres_password = "POSTGRES_PASSWORD"
	env_postgres_db       = "POSTGRES_DB"
)

type Env struct {
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
}

func (e *Env) GetAll() {
	const (
		port_max = 65535
		port_min = 1000
	)

	// --------- get PostgresHost
	if e.PostgresHost = os.Getenv(env_postgres_host); len(e.PostgresHost) == 0 {
		log.Fatalf("%s environement variable missing", env_postgres_host)
	}

	// --------- get PostgresPort
	postgres_port, err := strconv.Atoi(os.Getenv(env_postgres_port))

	if err != nil {
		log.Fatalf("%s environement variable is corrupted or missing", env_postgres_port)
	} else if postgres_port < port_min || postgres_port > port_max {
		log.Fatalf("port need to be included between %d and %d (%d)", port_min, port_max, postgres_port)
	}

	e.PostgresPort = postgres_port

	// --------- get PostgresUser
	if e.PostgresUser = os.Getenv(env_postgres_user); len(e.PostgresUser) == 0 {
		log.Fatalf("%s environement variable missing", env_postgres_user)
	}

	// --------- get PostgresPassword
	if e.PostgresPassword = os.Getenv(env_postgres_password); len(e.PostgresPassword) == 0 {
		log.Fatalf("%s environement variable missing", env_postgres_password)
	}

	// --------- get PostgresDB
	if e.PostgresDB = os.Getenv(env_postgres_db); len(e.PostgresDB) == 0 {
		log.Fatalf("%s environement variable missing", env_postgres_db)
	}
}

func readEnv() (env Env) {
	env.GetAll()

	return
}

func main() {
	env := readEnv()

	if err := postgres.Init(postgres.Config{
		Driver:   database_driver,
		Host:     env.PostgresHost,
		Port:     env.PostgresPort,
		User:     env.PostgresUser,
		Password: env.PostgresPassword,
		Dbname:   env.PostgresDB,
	}); err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	log.Println("connected to postgres")

	importer.Import()
}
