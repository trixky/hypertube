package internal

import (
	"context"
	"testing"

	"github.com/trixky/hypertube/.shared/databases"
	"github.com/trixky/hypertube/.shared/environment"
	_test "github.com/trixky/hypertube/.shared/test"
	"github.com/trixky/hypertube/api-media/proto"
	"github.com/trixky/hypertube/api-media/queries"
	"github.com/trixky/hypertube/api-media/utils"
	"google.golang.org/grpc/metadata"
)

func init() {
	// Set environment config
	environment_config := environment.Config{
		ENV_grpc_port:         "API_MEDIA_GRPC_PORT",
		ENV_grpc_gateway_port: "API_MEDIA_GRPC_GATEWAY_PORT",
	}

	environment.Postgres.GetAll()                // Get postgres environment
	environment.Redis.GetAll()                   // Get redis environment
	environment.Grpc.GetAll(&environment_config) // Get grpc environment

	databases.InitPostgres() // Init postgres
	databases.InitRedis()    // Init redis
	queries.InitSqlc()       // Init sqlc queries
}

func TestUpdateComment(t *testing.T) {
	server := &MediaServer{}

	tests := []struct {
		token           string
		input           *proto.UpdateCommentRequest
		corrupted_token bool
		invalid_token   bool
		error_expected  bool
	}{
		// ------------------------- Failed expected
		{ // Token missing #1
			input:          &proto.UpdateCommentRequest{},
			error_expected: true,
		},
		{ // Token invalid
			token:          "576100d0-0c2b-11ed-861d-0242ac120002",
			input:          &proto.UpdateCommentRequest{},
			invalid_token:  true,
			error_expected: true,
		},
		{ // Token corrupted
			token:           "528ff7b4-0c2b-11ed-861d-0242ac120002",
			input:           &proto.UpdateCommentRequest{},
			corrupted_token: true,
			error_expected:  true,
		},
		{ // Comment doesn't exist
			token: "f944c98c-0c2a-11ed-861d-0242ac120002",
			input: &proto.UpdateCommentRequest{
				CommentId: 424242,
				Content:   "hello",
			},
			error_expected: true,
		},
		{ // Comment is too short
			token: "f944c98c-0c2a-11ed-861d-0242ac120002",
			input: &proto.UpdateCommentRequest{
				CommentId: 2,
				Content:   "h",
			},
			error_expected: true,
		},
		{ // Comment is too long
			token: "f944c98c-0c2a-11ed-861d-0242ac120002",
			input: &proto.UpdateCommentRequest{
				CommentId: 2,
				Content:   utils.RandomString(100_000),
			},
			error_expected: true,
		},
		// ------------------------- Success expected
		{
			token: "f944c98c-0c2a-11ed-861d-0242ac120002",
			input: &proto.UpdateCommentRequest{
				CommentId: 2,
				Content:   "hello",
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

		if _, err := server.UpdateComment(ctx, test.input); (err != nil) != test.error_expected {
			t.Fatal(err)
		}
	}
}
