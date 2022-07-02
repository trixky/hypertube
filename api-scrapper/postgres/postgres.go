package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
	"github.com/trixky/hypertube/api-scrapper/sqlc"
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
	SqlcQueries *sqlc.Queries
}

var DB Database

func Init(config Config) error {
	ctx := context.Background()

	pool, err := pgx.Connect(ctx, config.Compile_config())
	if err != nil {
		return err
	}

	if err = pool.Ping(ctx); err != nil {
		return err
	}

	DB.SqlcQueries = sqlc.New(pool)

	return nil
}
