package internal

import (
	"context"
	"testing"

	initializer "github.com/trixky/hypertube/.shared/databases"
	"github.com/trixky/hypertube/.shared/environment"
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

	initializer.InitPostgres() // Init DBs
	initializer.InitRedis()
}

func TestInternalLogin(t *testing.T) {
	server := &AuthServer{}

	type Input struct {
		Usermame  string
		Firstname string
		Lastname  string
		Email     string
		Password  string
	}

	tests := []struct {
		input          *proto.InternalLoginRequest
		error_expected bool
	}{
		// ------------------------- Failed expected
		{ // Input missing #1
			input:          nil,
			error_expected: true,
		},
		{ // Input missing #2
			input:          &proto.InternalLoginRequest{},
			error_expected: true,
		},
		{ // Email missing
			input: &proto.InternalLoginRequest{
				Password: "1e8392fcefc860ef9714dcf4ad2249a995118c7a3bdbf4a96e8ffd7fe354c2e6",
			},
			error_expected: true,
		},
		{ // Email corrupted
			input: &proto.InternalLoginRequest{
				Email:    "x",
				Password: "1e8392fcefc860ef9714dcf4ad2249a995118c7a3bdbf4a96e8ffd7fe354c2e6",
			},
			error_expected: true,
		},
		{ // Email invalid
			input: &proto.InternalLoginRequest{
				Email:    "a@b.c",
				Password: "1e8392fcefc860ef9714dcf4ad2249a995118c7a3bdbf4a96e8ffd7fe354c2e6",
			},
			error_expected: true,
		},
		{ // Password missing
			input: &proto.InternalLoginRequest{
				Email: "a@b.c",
			},
			error_expected: true,
		},
		{ // Password corrupted
			input: &proto.InternalLoginRequest{
				Email:    "a@b.c",
				Password: "x",
			},
			error_expected: true,
		},
		{ // Password invalid
			input: &proto.InternalLoginRequest{
				Email:    "a@b.c",
				Password: "aa8392fcefc860ef9714dcf4ad2249a995118c7a3bdbf4a96e8ffd7fe354c2e6",
			},
			error_expected: true,
		},
		// ------------------------- Success expected
		{
			input: &proto.InternalLoginRequest{
				Email:    "email.test_1@test.com",
				Password: "1e8392fcefc860ef9714dcf4ad2249a995118c7a3bdbf4a96e8ffd7fe354c2e6",
			},
			error_expected: false,
		},
		{
			input: &proto.InternalLoginRequest{
				Email:    "email.test_2@test.com",
				Password: "1e8392fcefc860ef9714dcf4ad2249a995118c7a3bdbf4a96e8ffd7fe354c2e6",
			},
			error_expected: false,
		},
		{
			input: &proto.InternalLoginRequest{
				Email:    "email.test_3@test.com",
				Password: "1e8392fcefc860ef9714dcf4ad2249a995118c7a3bdbf4a96e8ffd7fe354c2e6",
			},
			error_expected: false,
		},
		{
			input: &proto.InternalLoginRequest{
				Email:    "email.test_4@test.com",
				Password: "1e8392fcefc860ef9714dcf4ad2249a995118c7a3bdbf4a96e8ffd7fe354c2e6",
			},
			error_expected: false,
		},
	}

	for _, test := range tests {
		if _, err := server.InternalLogin(context.Background(), test.input); (err != nil) != test.error_expected {
			t.Fatal(err)
		}
	}
}
