package internal

import (
	"context"
	"testing"

	"github.com/trixky/hypertube/.shared/databases"
	"github.com/trixky/hypertube/.shared/environment"
	_test "github.com/trixky/hypertube/.shared/test"
	"github.com/trixky/hypertube/api-media/proto"
	"github.com/trixky/hypertube/api-media/queries"
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

func TestSearch(t *testing.T) {
	server := &MediaServer{}

	query := "42"
	var page uint32 = 1
	year := "year"
	var yearNumber uint32 = 199
	sortAsc := "asc"
	var rating float32 = 7.2

	tests := []struct {
		token           string
		input           *proto.SearchRequest
		corrupted_token bool
		invalid_token   bool
		error_expected  bool
	}{
		// ------------------------- Failed expected
		{ // Token missing #1
			input:          &proto.SearchRequest{},
			error_expected: true,
		},
		{ // Token invalid
			token:          "576100d0-0c2b-11ed-861d-0242ac120002",
			input:          &proto.SearchRequest{},
			invalid_token:  true,
			error_expected: true,
		},
		{ // Token corrupted
			token:           "528ff7b4-0c2b-11ed-861d-0242ac120002",
			input:           &proto.SearchRequest{},
			corrupted_token: true,
			error_expected:  true,
		},
		// ------------------------- Success expected
		{
			token:          "f944c98c-0c2a-11ed-861d-0242ac120002",
			input:          &proto.SearchRequest{},
			error_expected: false,
		},
		{
			token: "f944c98c-0c2a-11ed-861d-0242ac120002",
			input: &proto.SearchRequest{
				SortBy: &year,
			},
			error_expected: false,
		},
		{
			token: "f944c98c-0c2a-11ed-861d-0242ac120002",
			input: &proto.SearchRequest{
				SortOrder: &sortAsc,
			},
			error_expected: false,
		},
		{
			token: "f944c98c-0c2a-11ed-861d-0242ac120002",
			input: &proto.SearchRequest{
				Year: &yearNumber,
			},
			error_expected: false,
		},
		{
			token: "f944c98c-0c2a-11ed-861d-0242ac120002",
			input: &proto.SearchRequest{
				Rating: &rating,
			},
			error_expected: false,
		},
		{
			token: "f944c98c-0c2a-11ed-861d-0242ac120002",
			input: &proto.SearchRequest{
				Query: &query,
			},
			error_expected: false,
		},
		{
			token: "f944c98c-0c2a-11ed-861d-0242ac120002",
			input: &proto.SearchRequest{
				GenreIds: []int32{1, 2},
			},
			error_expected: false,
		},
		{
			token: "f944c98c-0c2a-11ed-861d-0242ac120002",
			input: &proto.SearchRequest{
				Query:     &query,
				Page:      &page,
				SortBy:    &year,
				SortOrder: &sortAsc,
				Year:      &yearNumber,
				Rating:    &rating,
				GenreIds:  []int32{1, 2},
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

		if _, err := server.Search(ctx, test.input); (err != nil) != test.error_expected {
			t.Fatal(err)
		}
	}
}
