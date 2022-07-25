package databases

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
	"github.com/trixky/hypertube/.shared/environment"
)

func ErrorIsDuplication(err error) bool {
	return strings.Contains(err.Error(), "duplicate")
}

const (
	postgres_driver = "postgres"
)

type PostgresConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

func (d *PostgresConfig) CompileConfig() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", d.Host, d.Port, d.User, d.Password, d.Dbname)
}

var SqlDatabase *sql.DB

func InitPostgres() error {
	config := PostgresConfig{
		Driver:   postgres_driver,
		Host:     environment.Postgres.Host,
		Port:     environment.Postgres.Port,
		User:     environment.Postgres.User,
		Password: environment.Postgres.Password,
		Dbname:   environment.Postgres.DBname,
	}

	log.Printf("start connection to postgres on %s:%d\n", environment.Postgres.Host, environment.Postgres.Port)

	sql_database, err := sql.Open(config.Driver, config.CompileConfig())

	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	if err = sql_database.Ping(); err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	SqlDatabase = sql_database

	return nil
}
