package server

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/trixky/hypertube/api-auth/databases"
	pb "github.com/trixky/hypertube/api-auth/proto"
	"github.com/trixky/hypertube/api-auth/sanitizer"
	"github.com/trixky/hypertube/api-auth/sqlc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// const expiration = time.Hour * 24 * 1000

func sanitizeInternalLogin(in *pb.InternalLoginRequest) error {
	if err := sanitizer.SanitizeEmail(in.GetEmail()); err != nil { // email
		return err
	}
	if err := sanitizer.SanitizeSHA256Password(in.GetPassword()); err != nil { // password
		return err
	}

	return nil
}

func (s *AuthServer) InternalLogin(ctx context.Context, in *pb.InternalLoginRequest) (*pb.InternalLoginResponse, error) {
	// -------------------- sanitize
	if err := sanitizeInternalLogin(in); err != nil {
		return nil, err
	}

	// -------------------- db
	new_user, err := databases.DBs.SqlcQueries.GetUserByCredentials(context.Background(), sqlc.GetUserByCredentialsParams{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	})

	log.Println("apres le getuserbycredentials:")
	log.Println(new_user.ID)
	log.Println(err)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.PermissionDenied, "incorrect email and/or password")
		}
		return nil, status.Errorf(codes.Internal, "connection failed")
	}

	// -------------------- cache
	token := uuid.New().String() // token generation

	if err := databases.AddToken(new_user.ID, token); err != nil {
		return nil, status.Errorf(codes.Internal, "token generation failed")
	}

	return &pb.InternalLoginResponse{
		Token: token,
	}, nil
}
