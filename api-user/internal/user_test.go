package internal

import (
	"context"
	"strconv"
	"testing"

	"github.com/trixky/hypertube/.shared/databases"
	"github.com/trixky/hypertube/.shared/environment"
	_test "github.com/trixky/hypertube/.shared/test"
	"github.com/trixky/hypertube/api-user/proto"
	"github.com/trixky/hypertube/api-user/queries"
	"google.golang.org/grpc/metadata"
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

func TestGetUser(t *testing.T) {
	server := &UserServer{}

	tests := []struct {
		input           *proto.GetUserRequest
		token           string
		corrupted_token bool
		invalid_token   bool
		error_expected  bool
	}{
		// ------------------------- Failed expected
		{ // Token missing #1
			input:          nil, // nil == 0
			error_expected: true,
		},
		{ // Token missing #2
			input: &proto.GetUserRequest{
				Id: 0,
			},
			error_expected: true,
		},
		{ // Token missing #3
			input: &proto.GetUserRequest{
				Id: 201,
			},
			error_expected: true,
		},
		{ // Token invalid
			input: &proto.GetUserRequest{
				Id: 201,
			},
			token:          "66085e64-0c15-11ed-861d-0242ac120002",
			invalid_token:  true,
			error_expected: true,
		},
		{ // Token corrupted
			input: &proto.GetUserRequest{
				Id: 201,
			},
			token:           "a9cc2fc6-0c16-11ed-861d-0242ac120002",
			corrupted_token: true,
			error_expected:  true,
		},
		// ------------------------- Success expected
		{
			input: &proto.GetUserRequest{
				Id: 201,
			},
			token:          "33e529ae-0c14-11ed-861d-0242ac120002",
			error_expected: false,
		},
		{
			input: &proto.GetUserRequest{
				Id: 202,
			},
			token:          "4050a3e4-0c14-11ed-861d-0242ac120002",
			error_expected: false,
		},
		{
			input: &proto.GetUserRequest{
				Id: 203,
			},
			token:          "c57a31ea-0125-4a9e-b602-fb9faf7e3f45",
			error_expected: false,
		},
	}

	for _, test := range tests {
		ctx := context.Background()

		if len(test.token) > 0 {

			key := databases.REDIS_PATTERN_KEY_token + databases.REDIS_SEPARATOR + test.token + databases.REDIS_SEPARATOR + strconv.Itoa(int(test.input.GetId()))

			if test.invalid_token {
				test.token = _test.InvalidateToken(test.token)
			} else if test.corrupted_token {
				test.token = _test.CorrupteToken(test.token)
			}

			ctx = metadata.NewIncomingContext(ctx, metadata.New(map[string]string{
				"grpcgateway-cookie": "token=" + test.token + ";",
			}))

			databases.Redis.Set(key, databases.REDIS_EXTERNAL_none, 0)
		}

		if _, err := server.GetUser(ctx, test.input); (err != nil) != test.error_expected {
			t.Fatal(err)
		}
	}
}
