package databases

import (
	"context"
	"database/sql"
	"errors"

	"github.com/trixky/hypertube/api-auth/databases"
	"github.com/trixky/hypertube/api-auth/sqlc"
	"github.com/trixky/hypertube/api-auth/utils"
)

type PostgresDatabaseMock struct {
	Users []sqlc.User
}

var InitialsSqlMockUsers = []sqlc.User{ // Don't touche me
	{ // -------------------------------- 0
		ID:        0,
		Username:  "ini_0_username",
		Firstname: "ini_0_firstname",
		Lastname:  "ini_0_lastname",
		Email:     "ini_0_email@ini.ini",
		ID42: sql.NullInt32{
			Valid: false,
		},
		IDGoogle: sql.NullString{
			Valid: false,
		},
		Password: sql.NullString{
			Valid:  true,
			String: utils.EncryptPassword("22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368"),
			// 9ba9862e6c992abe652e575881931b22a29676b328d6cc00d04f1659888d40a1
		},
	},
	{ // -------------------------------- 1
		ID:        0,
		Username:  "ini_1_username",
		Firstname: "ini_1_firstname",
		Lastname:  "ini_1_lastname",
		Email:     "ini_1_email@ini.ini",
		ID42: sql.NullInt32{
			Valid: false,
		},
		IDGoogle: sql.NullString{
			Valid: false,
		},
		Password: sql.NullString{
			Valid:  true,
			String: utils.EncryptPassword("cd180755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368"),
			// c47b37adeba3497e1de2005bf3d0f8c28806a813527ce9f3cb917ae896b41cd4
		},
	},
}

// ---------------------------------- INIT

// InitPostgresMock intisalizes the postgres (mock)
func InitPostgresMock() {
	databases.DBs.SqlDatabase = nil
	databases.DBs.SqlcQueries = PostgresDatabaseMock{
		Users: InitialsSqlMockUsers,
	}
}

// ---------------------------------- IMPLEMENTATION

// Create42ExternalUser (mock methode)
func (pm PostgresDatabaseMock) Create42ExternalUser(ctx context.Context, arg sqlc.Create42ExternalUserParams) (sqlc.User, error) {
	return sqlc.User{}, nil
}

// GetUserBy42Id (mock methode)
func (pm PostgresDatabaseMock) GetUserBy42Id(ctx context.Context, id42 sql.NullInt32) (sqlc.User, error) {
	return sqlc.User{}, nil
}

// CreateGoogleExternalUser (mock methode)
func (pm PostgresDatabaseMock) CreateGoogleExternalUser(ctx context.Context, arg sqlc.CreateGoogleExternalUserParams) (sqlc.User, error) {
	return sqlc.User{}, nil
}

// GetUserByGoogleId (mock methode)
func (pm PostgresDatabaseMock) GetUserByGoogleId(ctx context.Context, idGoogle sql.NullString) (sqlc.User, error) {
	return sqlc.User{}, nil
}

// GetInternalUserByCredentials (mock methode)
func (pm PostgresDatabaseMock) GetInternalUserByCredentials(ctx context.Context, arg sqlc.GetInternalUserByCredentialsParams) (sqlc.User, error) {
	for _, user := range pm.Users {
		if user.Email == arg.Email && user.Password == arg.Password {
			return user, nil
		}
	}

	return sqlc.User{}, errors.New("no user finded")
}

// UpdateUserPassword (mock methode)
func (pm PostgresDatabaseMock) UpdateUserPassword(ctx context.Context, arg sqlc.UpdateUserPasswordParams) error {
	return nil
}

// GetInternalUserByEmail (mock methode)
func (pm PostgresDatabaseMock) GetInternalUserByEmail(ctx context.Context, email string) (sqlc.User, error) {
	for _, user := range pm.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return sqlc.User{}, errors.New("no user finded")
}

// CreateInternalUser (mock methode)
func (pm PostgresDatabaseMock) CreateInternalUser(ctx context.Context, arg sqlc.CreateInternalUserParams) (sqlc.User, error) {
	new_id := int64(len(pm.Users))

	new_sqlc_user := sqlc.User{
		ID:        new_id,
		Username:  arg.Username,
		Firstname: arg.Firstname,
		Lastname:  arg.Lastname,
		Email:     arg.Email,
		ID42: sql.NullInt32{
			Valid: false,
		},
		IDGoogle: sql.NullString{
			Valid: false,
		},
		Password: arg.Password,
	}

	pm.Users = append(pm.Users, new_sqlc_user)

	return new_sqlc_user, nil
}
