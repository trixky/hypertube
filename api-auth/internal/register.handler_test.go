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

func TestInternalRegister(t *testing.T) {
	server := &AuthServer{}

	type Input struct {
		Usermame  string
		Firstname string
		Lastname  string
		Email     string
		Password  string
	}

	tests := []struct {
		input          *proto.InternalRegisterRequest
		error_expected bool
	}{
		// ------------------------- Failed expected
		{ // Input missing #1
			input:          nil,
			error_expected: true,
		},
		{ // Input missing #2
			input:          &proto.InternalRegisterRequest{},
			error_expected: true,
		},
		{ // Username missing
			input: &proto.InternalRegisterRequest{
				Firstname: "input_firstname",
				Lastname:  "input_lastname",
				Email:     "a@b.c",
				Password:  "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
			},
			error_expected: true,
		},
		{ // Username corrupted
			input: &proto.InternalRegisterRequest{
				Username:  "x",
				Firstname: "input_firstname",
				Lastname:  "input_lastname",
				Email:     "a@b.c",
				Password:  "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
			},
			error_expected: true,
		},
		{ // Firstname missing
			input: &proto.InternalRegisterRequest{
				Username: "input_username",
				Lastname: "input_lastname",
				Email:    "a@b.c",
				Password: "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
			},
			error_expected: true,
		},
		{ // Firstname corrupted
			input: &proto.InternalRegisterRequest{
				Firstname: "x",
				Username:  "input_username",
				Lastname:  "input_lastname",
				Email:     "a@b.c",
				Password:  "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
			},
			error_expected: true,
		},
		{ // Lastname missing
			input: &proto.InternalRegisterRequest{
				Username:  "input_username",
				Firstname: "input_firstname",
				Email:     "a@b.c",
				Password:  "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
			},
			error_expected: true,
		},
		{ // Lastname corrupted
			input: &proto.InternalRegisterRequest{
				Username:  "input_username",
				Firstname: "input_firstname",
				Lastname:  "x",
				Email:     "a@b.c",
				Password:  "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
			},
			error_expected: true,
		},
		{ // Email missing
			input: &proto.InternalRegisterRequest{
				Username:  "input_username",
				Firstname: "input_firstname",
				Lastname:  "input_lastname",
				Password:  "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
			},
			error_expected: true,
		},
		{ // Email corrupted
			input: &proto.InternalRegisterRequest{
				Username:  "input_username",
				Firstname: "input_firstname",
				Lastname:  "input_lastname",
				Email:     "x",
				Password:  "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
			},
			error_expected: true,
		},
		{ // Password missing
			input: &proto.InternalRegisterRequest{
				Username:  "input_username",
				Firstname: "input_firstname",
				Lastname:  "input_lastname",
				Email:     "a@b.c",
			},
			error_expected: true,
		},
		{ // Password corrupted (one char missing)
			input: &proto.InternalRegisterRequest{
				Username:  "input_username",
				Firstname: "input_firstname",
				Lastname:  "input_lastname",
				Email:     "a@b.c",
				Password:  "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f36",
			},
			error_expected: true,
		},
		// ------------------------- Success expected
		{
			input: &proto.InternalRegisterRequest{
				Username:  "input_username_1",
				Firstname: "input_firstname_1",
				Lastname:  "input_lastname_1",
				Email:     "a1@b.c",
				Password:  "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
			},
			error_expected: false,
		},
		{
			input: &proto.InternalRegisterRequest{
				Username:  "input_username_2",
				Firstname: "input_firstname_2",
				Lastname:  "input_lastname_2",
				Email:     "a2@b.c",
				Password:  "cd180755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
			},
			error_expected: false,
		},
		{
			input: &proto.InternalRegisterRequest{
				Username:  "input_username_3",
				Firstname: "input_firstname_3",
				Lastname:  "input_lastname_3",
				Email:     "a3@b.c",
				Password:  "a4032963bb30edb4ea035afe180d2504a590700ba619f61fdae751dbbddf44f9",
			},
			error_expected: false,
		},
	}

	for _, test := range tests {
		if _, err := server.InternalRegister(context.Background(), test.input); (err != nil) != test.error_expected {
			t.Fatal(err)
		}
	}
}
