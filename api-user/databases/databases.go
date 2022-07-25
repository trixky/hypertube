package databases

import (
	"database/sql"
	"log"

	"github.com/go-redis/redis"
	"github.com/trixky/hypertube/.shared/environment"
	"github.com/trixky/hypertube/api-user/sqlc"
)

const (
	postgres_driver = "postgres"
)

type Databases struct {
	SqlDatabase *sql.DB
	SqlcQueries *sqlc.Queries
	Redis       *redis.Client
}

const (
	EXTERNAL_none      = "none"
	EXTERNAL_42        = "42"
	EXTERNAL_undefined = "undefined"
)

var DBs Databases

func InitDBs() {
	// ------------- POSTGRES
	log.Printf("start connection to postgres on %s:%d\n", environment.Postgres.Host, environment.Postgres.Port)

	if err := InitPostgres(PostgresConfig{
		Driver:   postgres_driver,
		Host:     environment.Postgres.Host,
		Port:     environment.Postgres.Port,
		User:     environment.Postgres.User,
		Password: environment.Postgres.Password,
		Dbname:   environment.Postgres.DBname,
	}); err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	// ------------- REDIS
	log.Println("start connection to redis on default address")
	if err := InitRedis(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
}
