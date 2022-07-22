package databases

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/trixky/hypertube/api-auth/databases"
	"github.com/trixky/hypertube/api-auth/sqlc"
)

type PostgresConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

// Compile_config compile all postgres config infos to one connectione string line
func (d *PostgresConfig) Compile_config() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", d.Host, d.Port, d.User, d.Password, d.Dbname)
}

// ---------------------------------- INIT

// InitPostgres intisalizes the postgres connection
func InitPostgres(config PostgresConfig) error {
	sql_database, err := sql.Open(config.Driver, config.Compile_config())

	if err != nil {
		return err
	}

	// Test the connection with a ping
	if err = sql_database.Ping(); err != nil {
		return err
	}

	databases.DBs.SqlDatabase = sql_database
	databases.DBs.SqlcQueries = sqlc.New(sql_database)

	return nil
}
