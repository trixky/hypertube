package internal

import (
	"context"
	"log"
	"testing"

	databases "github.com/trixky/hypertube/api-auth/databases/mock"
	"github.com/trixky/hypertube/api-auth/environment"
	initializer "github.com/trixky/hypertube/api-auth/initializer/mock"
	"github.com/trixky/hypertube/api-auth/proto"
)

func init() {
	log.Println("------------------------- INIT api-auth (TEST)")
	environment.E.GetAll() // Get environment variables
	initializer.InitDBs()  // Init DBs
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
				Password: databases.InitialsSqlMockUsers[0].Password.String,
			},
			error_expected: true,
		},
		{ // Email corrupted
			input: &proto.InternalLoginRequest{
				Email:    "x",
				Password: databases.InitialsSqlMockUsers[0].Password.String,
			},
			error_expected: true,
		},
		{ // Email invalid
			input: &proto.InternalLoginRequest{
				Email:    databases.InitialsSqlMockUsers[0].Email + "x",
				Password: databases.InitialsSqlMockUsers[0].Password.String,
			},
			error_expected: true,
		},
		{ // Password missing
			input: &proto.InternalLoginRequest{
				Email: databases.InitialsSqlMockUsers[0].Email,
			},
			error_expected: true,
		},
		{ // Password corrupted
			input: &proto.InternalLoginRequest{
				Email:    databases.InitialsSqlMockUsers[0].Email,
				Password: "x",
			},
			error_expected: true,
		},
		{ // Password invalid
			input: &proto.InternalLoginRequest{
				Email:    databases.InitialsSqlMockUsers[0].Email,
				Password: databases.InitialsSqlMockUsers[0].Password.String[1:] + "0",
			},
			error_expected: true,
		},
		// ------------------------- Success expected
		{
			input: &proto.InternalLoginRequest{
				Email:    databases.InitialsSqlMockUsers[0].Email,
				Password: databases.InitialsSqlMockUsers[0].Password.String,
			},
			error_expected: false,
		},
		{
			input: &proto.InternalLoginRequest{
				Email:    databases.InitialsSqlMockUsers[1].Email,
				Password: databases.InitialsSqlMockUsers[1].Password.String,
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
