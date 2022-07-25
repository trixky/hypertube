package queries

import (
	"database/sql"

	"github.com/trixky/hypertube/api-media/sqlc"
)

var SqlcQueries *sqlc.Queries

func InitSqlc(db *sql.DB) {
	SqlcQueries = sqlc.New(db)
}
