package internal

import (
	"context"
	"database/sql"
	"encoding/base64"

	"github.com/google/uuid"
	"github.com/trixky/hypertube/api-auth/databases"
	pb "github.com/trixky/hypertube/api-auth/proto"
	"github.com/trixky/hypertube/api-auth/sanitizer"
	"github.com/trixky/hypertube/api-auth/sqlc"
	"github.com/trixky/hypertube/api-auth/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func sanitizeInternalRegister(in *pb.InternalRegisterRequest) error {
	if err := sanitizer.SanitizeEmail(in.GetEmail()); err != nil { // email
		return err
	}
	if err := sanitizer.SanitizeUsername(in.GetUsername()); err != nil { // username
		return err
	}
	if err := sanitizer.SanitizeFirstname(in.GetFirstname()); err != nil { // firstname
		return err
	}
	if err := sanitizer.SanitizeLastname(in.GetLastname()); err != nil { // lastname
		return err
	}
	if err := sanitizer.SanitizeSHA256Password(in.GetPassword()); err != nil { // password
		return err
	}

	return nil
}

func (s *AuthServer) InternalRegister(ctx context.Context, in *pb.InternalRegisterRequest) (*pb.GenericConnectionResponse, error) {
	// -------------------- sanitize
	if err := sanitizeInternalRegister(in); err != nil {
		return nil, err
	}

	// -------------------- db
	new_user, err := databases.DBs.SqlcQueries.CreateInternalUser(context.Background(), sqlc.CreateInternalUserParams{
		Email:     in.GetEmail(),
		Username:  in.GetUsername(),
		Firstname: in.GetFirstname(),
		Lastname:  in.GetLastname(),
		Password: sql.NullString{
			String: in.GetPassword(),
			Valid:  true,
		},
	})

	if err != nil {
		if databases.ErrorIsDuplication(err) {
			return nil, status.Errorf(codes.AlreadyExists, "email is already in use")
		}
		return nil, status.Errorf(codes.Internal, "user creation failed")
	}

	// -------------------- cache
	token := uuid.New().String() // token generation

	if err := databases.AddToken(new_user.ID, token, databases.EXTERNAL_none); err != nil {
		return nil, status.Errorf(codes.Internal, "token generation failed")
	}

	me, err := utils.HeaderCookieMeGeneration(utils.CookieMe{
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

	return &pb.GenericConnectionResponse{
		Token: token,
		Me:    base64.StdEncoding.EncodeToString([]byte(me.Value)),
	}, nil
}
