package databases

import (
	"database/sql"

	"github.com/go-redis/redis"
	"github.com/trixky/hypertube/api-user/sqlc"
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
