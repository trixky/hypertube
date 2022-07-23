package databases

import (
	"context"
	"database/sql"
	"log"

	"github.com/go-redis/redis"
	"github.com/trixky/hypertube/api-auth/environment"
	"github.com/trixky/hypertube/api-auth/sqlc"
)

const (
	postgres_driver = "postgres"
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

func InitDBs() {
	// ------------- POSTGRES
	log.Printf("start connection to postgres on %s:%d\n", environment.E.PostgresHost, environment.E.PostgresPort)

	if err := InitPostgres(PostgresConfig{
		Driver:   postgres_driver,
		Host:     environment.E.PostgresHost,
		Port:     environment.E.PostgresPort,
		User:     environment.E.PostgresUser,
		Password: environment.E.PostgresPassword,
		Dbname:   environment.E.PostgresDB,
	}); err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	// ------------- REDIS
	log.Println("start connection to redis on default address")
	if err := InitRedis(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
}
