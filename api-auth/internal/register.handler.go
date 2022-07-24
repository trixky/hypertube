package internal

import (
	"context"
	"database/sql"
	"encoding/base64"
	"log"

	"github.com/google/uuid"
	"github.com/trixky/hypertube/.shared/sanitizer"
	"github.com/trixky/hypertube/.shared/utils"
	"github.com/trixky/hypertube/api-auth/databases"
	"github.com/trixky/hypertube/api-auth/email"
	pb "github.com/trixky/hypertube/api-auth/proto"
	"github.com/trixky/hypertube/api-auth/sqlc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// sanitizeInternalRegister sanitizes inputs for "InternalRegister" route
func sanitizeInternalRegister(in *pb.InternalRegisterRequest) error {
	if err := sanitizer.SanitizeEmail(in.GetEmail()); err != nil { // Email
		return err
	}
	if err := sanitizer.SanitizeUsername(in.GetUsername()); err != nil { // Username
		return err
	}
	if err := sanitizer.SanitizeFirstname(in.GetFirstname()); err != nil { // Firstname
		return err
	}
	if err := sanitizer.SanitizeLastname(in.GetLastname()); err != nil { // Lastname
		return err
	}
	if err := sanitizer.SanitizeSHA256Password(in.GetPassword()); err != nil { // Password
		return err
	}

	return nil
}

// InternalRegister Handles the "InternalRegister" route
func (s *AuthServer) InternalRegister(ctx context.Context, in *pb.InternalRegisterRequest) (*pb.GenericConnectionResponse, error) {
	// -------------------- Sanitize
	if err := sanitizeInternalRegister(in); err != nil {
		return nil, err
	}

	// -------------------- DB
	new_user, err := databases.DBs.SqlcQueries.CreateInternalUser(context.Background(), sqlc.CreateInternalUserParams{
		Email:     in.GetEmail(),
		Username:  in.GetUsername(),
		Firstname: in.GetFirstname(),
		Lastname:  in.GetLastname(),
		Password: sql.NullString{
			String: utils.EncryptPassword(in.GetPassword()),
			Valid:  true,
		},
	})

	if err != nil {
		if databases.ErrorIsDuplication(err) {
			return nil, status.Errorf(codes.AlreadyExists, "email is already in use")
		}
		return nil, status.Errorf(codes.Internal, "user creation failed")
	}

	// -------------------- Cache
	// Generate the token
	token := uuid.New().String()

	if err := databases.DBs.RedisQueries.AddToken(new_user.ID, token, databases.EXTERNAL_none); err != nil {
		return nil, status.Errorf(codes.Internal, "token generation failed")
	}

	me, err := utils.HeaderCookieUserGeneration(utils.User{
		Id:        int(new_user.ID),
		Username:  new_user.Username,
		Firstname: new_user.Firstname,
		Lastname:  new_user.Lastname,
		Email:     new_user.Email,
		External:  databases.EXTERNAL_none,
	}, false)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cookie generation failed")
	}

	// -------------------- Email
	go func() {
		if err := email.SendRegistrationConfirmation(in.GetEmail()); err != nil {
			log.Printf("sending registration email to [%s] failed: %s\n", in.GetEmail(), err.Error())
		}
	}()

	return &pb.GenericConnectionResponse{
		Token:    token,
		UserInfo: base64.StdEncoding.EncodeToString([]byte(me.Value)),
	}, nil
}
