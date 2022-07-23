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

func TestInternalRecoverPassword(t *testing.T) {
	server := &AuthServer{}

	type Input struct {
		Usermame  string
		Firstname string
		Lastname  string
		Email     string
		Password  string
	}

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
				Email: databases.InitialsSqlMockUsers[0].Email,
			},
			error_expected: false,
		},
		{
			input: &proto.InternalRecoverPasswordRequest{
				Email: databases.InitialsSqlMockUsers[1].Email,
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
