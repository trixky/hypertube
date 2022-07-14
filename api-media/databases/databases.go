package databases

import (
	"database/sql"

	"github.com/go-redis/redis"
	"github.com/trixky/hypertube/api-media/sqlc"
)

type Databases struct {
	SqlDatabase *sql.DB
	SqlcQueries *sqlc.Queries
	Redis       *redis.Client
}

var DBs Databases
