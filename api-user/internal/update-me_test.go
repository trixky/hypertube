package internal

import (
	"context"
	"testing"

	"github.com/trixky/hypertube/.shared/environment"
	_test "github.com/trixky/hypertube/.shared/test"
	"github.com/trixky/hypertube/api-user/databases"
	"github.com/trixky/hypertube/api-user/proto"
)

func init() {
	// Set environment config
	environment_config := environment.Config{
		ENV_grpc_port:         "API_AUTH_GRPC_PORT",
		ENV_grpc_gateway_port: "API_AUTH_GRPC_GATEWAY_PORT",
		ENV_http_port:         "API_AUTH_HTTP_PORT",
	}

	environment.Postgres.GetAll()                // Get postgres environment
	environment.Redis.GetAll()                   // Get redis environment
	environment.Grpc.GetAll(&environment_config) // Get grpc environment
	environment.Http.GetAll(&environment_config) // Get http environment
	environment.Api42.GetAll()                   // Get 42 api environment
	environment.ApiGoogle.GetAll()               // Get google api environment
	environment.Outlook.GetAll()                 // Get outlook environment

	databases.InitDBs() // Init DBs
}

func TestUpdateMe(t *testing.T) {
	server := &UserServer{}

	tests := []struct {
		input           *proto.UpdateMeRequest
		user_id         string
		corrupted_token bool
		invalid_token   bool
		error_expected  bool
	}{
		// ------------------------- Failed expected
		{ // Token missing
			input:          &proto.UpdateMeRequest{},
			user_id:        "202",
			error_expected: true,
		},
		{ // Token invalid
			input: &proto.UpdateMeRequest{
				Token: "b7f3844a-0c26-11ed-861d-0242ac120002",
			},
			user_id:        "202",
			invalid_token:  true,
			error_expected: true,
		},
		{ // Token corrupted
			input: &proto.UpdateMeRequest{
				Token: "aefa81d6-0c26-11ed-861d-0242ac120002",
			},
			user_id:         "202",
			corrupted_token: true,
			error_expected:  true,
		},
		{ // Username corrupted
			input: &proto.UpdateMeRequest{
				Token:           "aa67f090-0c26-11ed-861d-0242ac120002",
				Username:        "x",
				Firstname:       "firstname_change",
				Lastname:        "lastname_change",
				Email:           "email_change@test.com",
				CurrentPassword: "1e8392fcefc860ef9714dcf4ad2249a995118c7a3bdbf4a96e8ffd7fe354c2e6",
				NewPassword:     "0000062e6c992abe652e575881931b22a29676b328d6cc00d04f1659888d40a1",
			},
			user_id:        "202",
			error_expected: true,
		},
		{ // Firstname corrupted
			input: &proto.UpdateMeRequest{
				Token:           "a3a04dca-0c26-11ed-861d-0242ac120002",
				Username:        "username_change",
				Firstname:       "x",
				Lastname:        "lastname_change",
				Email:           "email_change@test.com",
				CurrentPassword: "1e8392fcefc860ef9714dcf4ad2249a995118c7a3bdbf4a96e8ffd7fe354c2e6",
				NewPassword:     "0000062e6c992abe652e575881931b22a29676b328d6cc00d04f1659888d40a1",
			},
			user_id:        "202",
			error_expected: true,
		},
		{ // Email corrupted
			input: &proto.UpdateMeRequest{
				Token:           "9c91a59c-0c26-11ed-861d-0242ac120002",
				Username:        "username_change",
				Firstname:       "firstname_change",
				Lastname:        "lastname_change",
				Email:           "x",
				CurrentPassword: "1e8392fcefc860ef9714dcf4ad2249a995118c7a3bdbf4a96e8ffd7fe354c2e6",
				NewPassword:     "0000062e6c992abe652e575881931b22a29676b328d6cc00d04f1659888d40a1",
			},
			user_id:        "202",
			error_expected: true,
		},
		{ // Current assword invalid
			input: &proto.UpdateMeRequest{
				Token:           "9875c330-0c26-11ed-861d-0242ac120002",
				Username:        "username_change",
				Firstname:       "firstname_change",
				Lastname:        "lastname_change",
				Email:           "email_change@test.com",
				CurrentPassword: "0e8392fcefc860ef9714dcf4ad2249a995118c7a3bdbf4a96e8ffd7fe354c2e6",
				NewPassword:     "1e8392fcefc860ef9714dcf4ad2249a995118c7a3bdbf4a96e8ffd7fe354c2e6",
			},
			user_id:        "202",
			error_expected: true,
		},
		{ // Current password corrupted
			input: &proto.UpdateMeRequest{
				Token:           "978aaae4-0c26-11ed-861d-0242ac120002",
				Username:        "username_change",
				Firstname:       "firstname_change",
				Lastname:        "lastname_change",
				Email:           "email_change@test.com",
				CurrentPassword: "@e8392fcefc860ef9714dcf4ad2249a995118c7a3bdbf4a96e8ffd7fe354c2e6",
				NewPassword:     "1e8392fcefc860ef9714dcf4ad2249a995118c7a3bdbf4a96e8ffd7fe354c2e6",
			},
			user_id:        "202",
			error_expected: true,
		},
		{ // New password corrupted
			input: &proto.UpdateMeRequest{
				Token:           "930f87c8-0c26-11ed-861d-0242ac120002",
				Username:        "username_change",
				Firstname:       "firstname_change",
				Lastname:        "lastname_change",
				Email:           "email_change@test.com",
				CurrentPassword: "1e8392fcefc860ef9714dcf4ad2249a995118c7a3bdbf4a96e8ffd7fe354c2e6",
				NewPassword:     "@1e8392fcefc860ef9714dcf4ad2249a995118c7a3bdbf4a96e8ffd7fe354c2e6",
			},
			user_id:        "202",
			error_expected: true,
		},
		// ------------------------- Success expected
		{
			input: &proto.UpdateMeRequest{
				Token:           "fdaf39e0-0c29-11ed-861d-0242ac120002",
				Username:        "username_change",
				Firstname:       "firstname_change",
				Lastname:        "lastname_change",
				Email:           "email_change@test.com",
				CurrentPassword: "1e8392fcefc860ef9714dcf4ad2249a995118c7a3bdbf4a96e8ffd7fe354c2e6",
				NewPassword:     "0000062e6c992abe652e575881931b22a29676b328d6cc00d04f1659888d40a1",
			},
			user_id:        "202",
			error_expected: false,
		},
	}

	for _, test := range tests {
		ctx := context.Background()

		if len(test.input.Token) > 0 {
			key := databases.REDIS_PATTERN_KEY_token + databases.REDIS_SEPARATOR + test.input.Token + databases.REDIS_SEPARATOR + test.user_id

			databases.DBs.Redis.Set(key, databases.EXTERNAL_none, 0)

			if test.invalid_token {
				test.input.Token = _test.InvalidateToken(test.input.Token)
			} else if test.corrupted_token {
				test.input.Token = _test.CorrupteToken(test.input.Token)
			}
		}

		if _, err := server.UpdateMe(ctx, test.input); (err != nil) != test.error_expected {
			t.Fatal(err)
		}
	}
}
