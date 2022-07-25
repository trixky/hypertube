package queries

import (
	"github.com/trixky/hypertube/.shared/databases"
	"github.com/trixky/hypertube/api-media/sqlc"
)

var SqlcQueries *sqlc.Queries

func InitSqlc() {
	SqlcQueries = sqlc.New(databases.SqlDatabase)
}
