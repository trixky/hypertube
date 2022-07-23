package internal

import (
	"context"
	"log"
	"testing"

	"github.com/trixky/hypertube/api-auth/environment"
	initializer "github.com/trixky/hypertube/api-auth/initializer/mock"
	"github.com/trixky/hypertube/api-auth/proto"
)

func init() {
	log.Println("------------------------- INIT api-auth (TEST)")
	environment.E.GetAll() // Get environment variables
	initializer.InitDBs()  // Init DBs
}

func TestInternalApplyRecoverPassword(t *testing.T) {
	server := &AuthServer{}

	type Input struct {
		Usermame  string
		Firstname string
		Lastname  string
		Email     string
		Password  string
	}

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
				NewPassword:   "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
			},
			error_expected: false,
		},
		{
			input: &proto.InternalApplyRecoverPasswordRequest{
				PasswordToken: "eeee318f-1b68-49f8-9ad7-e8695ad114a9",
				NewPassword:   "cd180755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
			},
			error_expected: false,
		},
	}

	for _, test := range tests {
		if _, err := server.InternalApplyRecoverPassword(context.Background(), test.input); (err != nil) != test.error_expected {
			t.Fatal(err)
		}
	}
}
