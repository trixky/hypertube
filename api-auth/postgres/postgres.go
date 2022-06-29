package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/trixky/hypertube/api-auth/sqlc"
)

type Config struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

func (d *Config) Compile_config() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", d.Host, d.Port, d.User, d.Password, d.Dbname)
}

type Database struct {
	SqlDatabase *sql.DB
	SqlcQueries *sqlc.Queries
}

var database Database

func Init(config Config) error {
	sql_database, err := sql.Open(config.Driver, config.Compile_config())

	if err != nil {
		return err
	}

	if err = sql_database.Ping(); err != nil {
		return err
	}

	database.SqlDatabase = sql_database
	database.SqlcQueries = sqlc.New(sql_database)

	return nil
}