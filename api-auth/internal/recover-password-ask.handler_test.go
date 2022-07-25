package internal

import (
	"context"
	"testing"

	"github.com/trixky/hypertube/.shared/environment"
	"github.com/trixky/hypertube/api-auth/databases"
	"github.com/trixky/hypertube/api-auth/proto"
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

func TestInternalRecoverPassword(t *testing.T) {
	server := &AuthServer{}

	tests := []struct {
		input          *proto.InternalRecoverPasswordRequest
		error_expected bool
	}{
		// ------------------------- Failed expected
		{ // Input missing #1
			input:          nil,
			error_expected: true,
		},
		{ // Input missing #2
			input:          &proto.InternalRecoverPasswordRequest{},
			error_expected: true,
		},
		{ // Input corrupted #1
			input: &proto.InternalRecoverPasswordRequest{
				Email: "a",
			},
			error_expected: true,
		},
		{ // Input corrupted #2
			input: &proto.InternalRecoverPasswordRequest{
				Email: "a@b",
			},
			error_expected: true,
		},
		{ // Input corrupted #3
			input: &proto.InternalRecoverPasswordRequest{
				Email: "a@b.",
			},
			error_expected: true,
		},
		{ // Input invalid #1
			input: &proto.InternalRecoverPasswordRequest{
				Email: "a@b.c",
			},
			error_expected: true,
		},
		{ // Input invalid #2
			input: &proto.InternalRecoverPasswordRequest{
				Email: "ab.cd@ef.gh",
			},
			error_expected: true,
		},
		// ------------------------- Success expected
		{
			input: &proto.InternalRecoverPasswordRequest{
				Email: "email.test_1@test.com",
			},
			error_expected: false,
		},
		{
			input: &proto.InternalRecoverPasswordRequest{
				Email: "email.test_2@test.com",
			},
			error_expected: false,
		},
		{
			input: &proto.InternalRecoverPasswordRequest{
				Email: "email.test_3@test.com",
			},
			error_expected: false,
		},
		{
			input: &proto.InternalRecoverPasswordRequest{
				Email: "email.test_4@test.com",
			},
			error_expected: false,
		},
	}

	for _, test := range tests {
		if _, err := server.InternalRecoverPassword(context.Background(), test.input); (err != nil) != test.error_expected {
			t.Fatal(err)
		}
	}
}
