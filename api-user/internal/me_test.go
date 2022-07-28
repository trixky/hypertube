package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/trixky/hypertube/.shared/databases"
	"github.com/trixky/hypertube/.shared/environment"
	_test "github.com/trixky/hypertube/.shared/test"
	"github.com/trixky/hypertube/api-user/proto"
	"github.com/trixky/hypertube/api-user/queries"
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

	databases.InitPostgres() // Init postgres
	databases.InitRedis()    // Init redis
	queries.InitSqlc()       // Init sqlc queries
}

func TestGetMe(t *testing.T) {
	server := &UserServer{}

	tests := []struct {
		input           *proto.GetMeRequest
		user_id         string
		corrupted_token bool
		invalid_token   bool
		error_expected  bool
	}{
		// ------------------------- Failed expected
		{ // Token missing #1
			input:          &proto.GetMeRequest{},
			user_id:        "201",
			error_expected: true,
		},
		{ // Token invalid
			input: &proto.GetMeRequest{
				Token: "576100d0-0c2b-11ed-861d-0242ac120002",
			},
			user_id:        "201",
			invalid_token:  true,
			error_expected: true,
		},

		{ // Token corrupted
			input: &proto.GetMeRequest{
				Token: "528ff7b4-0c2b-11ed-861d-0242ac120002",
			},
			user_id:         "202",
			corrupted_token: true,
			error_expected:  true,
		},
		// ------------------------- Success expected
		{
			input: &proto.GetMeRequest{
				Token: "f944c98c-0c2a-11ed-861d-0242ac120002",
			},
			user_id:        "201",
			error_expected: false,
		},
		{
			input: &proto.GetMeRequest{
				Token: "fef1aab2-0c2a-11ed-861d-0242ac120002",
			},
			user_id:        "202",
			error_expected: false,
		},
		{
			input: &proto.GetMeRequest{
				Token: "03554a5a-0c2b-11ed-861d-0242ac120002",
			},
			user_id:        "203",
			error_expected: false,
		},
		{
			input: &proto.GetMeRequest{
				Token: "075cedd8-0c2b-11ed-861d-0242ac120002",
			},
			user_id:        "204",
			error_expected: false,
		},
	}

	for _, test := range tests {
		ctx := context.Background()

		if len(test.input.Token) > 0 {
			key := databases.REDIS_PATTERN_KEY_token + databases.REDIS_SEPARATOR + test.input.Token + databases.REDIS_SEPARATOR + test.user_id

			databases.Redis.Set(key, databases.REDIS_EXTERNAL_none, 0)

			if test.invalid_token {
				test.input.Token = _test.InvalidateToken(test.input.Token)
			} else if test.corrupted_token {
				test.input.Token = _test.CorrupteToken(test.input.Token)
			}
		}

		if _, err := server.GetMe(ctx, test.input); (err != nil) != test.error_expected {
			t.Fatal(err)
		}
	}
}
