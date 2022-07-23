package internal

import (
	"context"
	"testing"

	initializer "github.com/trixky/hypertube/api-auth/databases"
	"github.com/trixky/hypertube/api-auth/environment"
	"github.com/trixky/hypertube/api-auth/proto"
)

func init() {
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
