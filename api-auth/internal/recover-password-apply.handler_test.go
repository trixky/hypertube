package internal

import (
	"context"
	"testing"

	"github.com/trixky/hypertube/.shared/databases"
	"github.com/trixky/hypertube/.shared/environment"
	"github.com/trixky/hypertube/api-auth/proto"
	"github.com/trixky/hypertube/api-auth/queries"
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

	databases.InitPostgres() // Init DBs
	databases.InitRedis()
}
func TestInternalApplyRecoverPassword(t *testing.T) {
	server := &AuthServer{}

	tests := []struct {
		input          *proto.InternalApplyRecoverPasswordRequest
		error_expected bool
	}{
		// ------------------------- Failed expected
		{ // Input missing #1
			input:          nil,
			error_expected: true,
		},
		{ // Input missing #2
			input:          &proto.InternalApplyRecoverPasswordRequest{},
			error_expected: true,
		},
		{ // Password token missing
			input: &proto.InternalApplyRecoverPasswordRequest{
				NewPassword: "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
			},
			error_expected: true,
		},
		{ // Password token corrupted #1
			input: &proto.InternalApplyRecoverPasswordRequest{
				PasswordToken: "x",
				NewPassword:   "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
			},
			error_expected: true,
		},
		{ // Password token corrupted #2
			input: &proto.InternalApplyRecoverPasswordRequest{
				PasswordToken: "57a31ea-0125-4a9e-b602-fb9faf7e3f45",
				NewPassword:   "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
			},
			error_expected: true,
		},
		{ // Password token corrupted #3
			input: &proto.InternalApplyRecoverPasswordRequest{
				PasswordToken: "c57a31ea.0125-4a9e-b602-fb9faf7e3f45",
				NewPassword:   "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
			},
			error_expected: true,
		},
		{ // New password token missing
			input: &proto.InternalApplyRecoverPasswordRequest{
				PasswordToken: "c57a31ea-0125-4a9e-b602-fb9faf7e3f45",
			},
			error_expected: true,
		},
		{ // New password token corrupted #1
			input: &proto.InternalApplyRecoverPasswordRequest{
				PasswordToken: "c57a31ea-0125-4a9e-b602-fb9faf7e3f45",
				NewPassword:   "x",
			},
			error_expected: true,
		},
		{ // New password token corrupted #2
			input: &proto.InternalApplyRecoverPasswordRequest{
				PasswordToken: "c57a31ea-0125-4a9e-b602-fb9faf7e3f45",
				NewPassword:   "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368a",
			},
			error_expected: true,
		},
		{ // New password token corrupted #3
			input: &proto.InternalApplyRecoverPasswordRequest{
				PasswordToken: "c57a31ea-0125-4a9e-b602-fb9faf7e3f45",
				NewPassword:   "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f36",
			},
			error_expected: true,
		},
		{ // Input invalid #1
			input: &proto.InternalApplyRecoverPasswordRequest{
				PasswordToken: "a57a31ea-0125-4a9e-b602-fb9faf7e3f45",
				NewPassword:   "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
			},
			error_expected: true,
		},
		{ // Input invalid #2
			input: &proto.InternalApplyRecoverPasswordRequest{
				PasswordToken: "c57a31ea-0125-4a9e-b602-fb9faf7e3f4b",
				NewPassword:   "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
			},
			error_expected: true,
		},
		{ // Input invalid #3
			input: &proto.InternalApplyRecoverPasswordRequest{
				PasswordToken: "a57a31ea-0125-4a9e-b602-fc9faf7e3f45",
				NewPassword:   "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
			},
			error_expected: true,
		},
		// ------------------------- Success expected
		{
			input: &proto.InternalApplyRecoverPasswordRequest{
				PasswordToken: "c57a31ea-0125-4a9e-b602-fb9faf7e3f45",
				NewPassword:   "eee392fcefc860ef9714dcf4ad2249a995118c7a3bdbf4a96e8ffd7fe354ceee",
			},
			error_expected: false,
		},
	}

	queries.AddPasswordToken(105, "c57a31ea-0125-4a9e-b602-fb9faf7e3f45")

	for _, test := range tests {
		if _, err := server.InternalApplyRecoverPassword(context.Background(), test.input); (err != nil) != test.error_expected {
			t.Fatal(err)
		}
	}
}
