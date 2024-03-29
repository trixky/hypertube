package internal

import (
	"context"
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

	databases.InitPostgres() // Init postgres
	databases.InitRedis()    // Init redis
	queries.InitSqlc()       // Init sqlc queries
}

func TestGetUserMovies(t *testing.T) {
	server := &UserServer{}

	// var page uint32 = 1
	tests := []struct {
		token           string
		input           *proto.GetUserMoviesRequest
		corrupted_token bool
		invalid_token   bool
		error_expected  bool
	}{
		// ------------------------- Failed expected
		{ // Token missing #1
			input:          &proto.GetUserMoviesRequest{},
			error_expected: true,
		},
		{ // Token invalid
			token:          "576100d0-0c2b-11ed-861d-0242ac120002",
			input:          &proto.GetUserMoviesRequest{},
			invalid_token:  true,
			error_expected: true,
		},
		{ // Token corrupted
			token:           "528ff7b4-0c2b-11ed-861d-0242ac120002",
			input:           &proto.GetUserMoviesRequest{},
			corrupted_token: true,
			error_expected:  true,
		},
		{ // User does not exists
			token: "f944c98c-0c2a-11ed-861d-0242ac120002",
			input: &proto.GetUserMoviesRequest{
				UserId: 42424242,
			},
			error_expected: true,
		},
		// ------------------------- Success expected
		{
			token: "f944c98c-0c2a-11ed-861d-0242ac120002",
			input: &proto.GetUserMoviesRequest{
				UserId: 201,
			},
			error_expected: false,
		},
	}

	for _, test := range tests {
		ctx := context.Background()

		if len(test.token) > 0 {
			if test.invalid_token {
				test.token = _test.InvalidateToken(test.token)
			} else if test.corrupted_token {
				test.token = _test.CorrupteToken(test.token)
			}

			ctx = metadata.NewIncomingContext(ctx, metadata.New(map[string]string{
				"grpcgateway-cookie": "token=" + test.token + ";",
			}))
		}

		if _, err := server.GetUserMovies(ctx, test.input); (err != nil) != test.error_expected {
			t.Fatal(err)
		}
	}
}
