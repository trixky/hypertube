package initializer

import (
	"log"

	databases "github.com/trixky/hypertube/api-auth/databases/prod"
	"github.com/trixky/hypertube/api-auth/environment"
)

const (
	host            = "0.0.0.0"
	postgres_driver = "postgres"
)

func InitDBs() {
	// ------------- POSTGRES
	log.Printf("start connection to postgres on %s:%d\n", environment.E.PostgresHost, environment.E.PostgresPort)
	if err := databases.InitPostgres(databases.PostgresConfig{
		Driver:   postgres_driver,
		Host:     environment.E.PostgresHost,
		Port:     environment.E.PostgresPort,
		User:     environment.E.PostgresUser,
		Password: environment.E.PostgresPassword,
		Dbname:   environment.E.PostgresDB,
	}); err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	// ------------- REDIS
	log.Println("start connection to redis on default address")
	if err := databases.InitRedis(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
}
