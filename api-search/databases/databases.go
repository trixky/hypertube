package databases

import (
	"database/sql"

	"github.com/trixky/hypertube/api-search/sqlc"
)

type Databases struct {
	SqlDatabase *sql.DB
	SqlcQueries *sqlc.Queries
}

var DBs Databases
