package databases

import (
	"context"
	"database/sql"

	"github.com/go-redis/redis"
	"github.com/trixky/hypertube/api-auth/sqlc"
)

type SQLQueries interface {
	Create42ExternalUser(ctx context.Context, arg sqlc.Create42ExternalUserParams) (sqlc.User, error)
	GetUserBy42Id(ctx context.Context, id42 sql.NullInt32) (sqlc.User, error)
	CreateGoogleExternalUser(ctx context.Context, arg sqlc.CreateGoogleExternalUserParams) (sqlc.User, error)
	GetUserByGoogleId(ctx context.Context, idGoogle sql.NullString) (sqlc.User, error)
	GetInternalUserByCredentials(ctx context.Context, arg sqlc.GetInternalUserByCredentialsParams) (sqlc.User, error)
	UpdateUserPassword(ctx context.Context, arg sqlc.UpdateUserPasswordParams) error
	GetInternalUserByEmail(ctx context.Context, email string) (sqlc.User, error)
	CreateInternalUser(ctx context.Context, arg sqlc.CreateInternalUserParams) (sqlc.User, error)
}

type RedisTokenInfo struct {
	Id       int64
	External string
}

type REDISQueries interface {
	AddToken(user_id int64, token string, external string) error
	RetrieveToken(token string) (*RedisTokenInfo, error)
	AddPasswordToken(user_id int64, token string) error
	RetrievePasswordToken(token string, delete bool) (*RedisTokenInfo, error)
}

type Databases struct {
	// SQL
	SqlDatabase *sql.DB
	SqlcQueries SQLQueries
	// REDIS
	RedisDatabase *redis.Client
	RedisQueries  REDISQueries
}

const (
	EXTERNAL_none   = "none"
	EXTERNAL_42     = "42"
	EXTERNAL_GOOGLE = "google"
)

var DBs Databases
