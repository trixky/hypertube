package internal

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"

	"github.com/google/uuid"
	"github.com/trixky/hypertube/api-auth/databases"
	pb "github.com/trixky/hypertube/api-auth/proto"
	"github.com/trixky/hypertube/api-auth/sanitizer"
	"github.com/trixky/hypertube/api-auth/sqlc"
	"github.com/trixky/hypertube/api-auth/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// sanitizeInternalLogin sanitizes inputs for "InternalLogin" route
func sanitizeInternalLogin(in *pb.InternalLoginRequest) error {
	if err := sanitizer.SanitizeEmail(in.GetEmail()); err != nil { // email
		return err
	}
	if err := sanitizer.SanitizeSHA256Password(in.GetPassword()); err != nil { // password
		return err
	}

	return nil
}

// InternalLogin Handles the "InternalLogin" route
func (s *AuthServer) InternalLogin(ctx context.Context, in *pb.InternalLoginRequest) (*pb.GenericConnectionResponse, error) {
	// -------------------- Sanitize
	if err := sanitizeInternalLogin(in); err != nil {
		return nil, err
	}

	// -------------------- DB
	user, err := databases.DBs.SqlcQueries.GetInternalUserByCredentials(context.Background(), sqlc.GetInternalUserByCredentialsParams{
		Email: in.GetEmail(),
		Password: sql.NullString{
			String: in.GetPassword(),
			Valid:  true,
		},
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.PermissionDenied, "incorrect email and/or password")
		}
		return nil, status.Errorf(codes.Internal, "connection failed")
	}

	// -------------------- Cache
	// Generates the token
	token := uuid.New().String()

	if err := databases.DBs.RedisQueries.AddToken(user.ID, token, databases.EXTERNAL_none); err != nil {
		return nil, status.Errorf(codes.Internal, "token generation failed")
	}

	me, err := utils.HeaderCookieMeGeneration(utils.User{
		Id:        int(user.ID),
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		External:  databases.EXTERNAL_none,
	}, false)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cookie generation failed")
	}

	return &pb.GenericConnectionResponse{
		Token:    token,
		UserInfo: base64.StdEncoding.EncodeToString([]byte(me.Value)),
	}, nil
}
