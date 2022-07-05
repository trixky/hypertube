package databases

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/trixky/hypertube/api-scrapper/sqlc"
)

func ErrorIsDuplication(err error) bool {
	return strings.Contains(err.Error(), "duplicate")
}

type PostgresConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

func (d *PostgresConfig) Compile_config() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", d.Host, d.Port, d.User, d.Password, d.Dbname)
}

func InitPosgres(config PostgresConfig) error {
	sql_database, err := sql.Open(config.Driver, config.Compile_config())

	if err != nil {
		return err
	}

	if err = sql_database.Ping(); err != nil {
		return err
	}

	DBs.SqlDatabase = sql_database
	DBs.SqlcQueries = sqlc.New(sql_database)

	return nil
}
