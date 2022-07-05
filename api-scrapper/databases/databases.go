package databases

import (
	"database/sql"

	"github.com/trixky/hypertube/api-scrapper/sqlc"
)

type Databases struct {
	SqlDatabase *sql.DB
	SqlcQueries *sqlc.Queries
}

var DBs Databases
