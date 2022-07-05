package internal

import (
	"context"
	"database/sql"
	"errors"

	"github.com/trixky/hypertube/api-auth/databases"
	pb "github.com/trixky/hypertube/api-auth/proto"
	"github.com/trixky/hypertube/api-auth/sanitizer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func sanitizeInternalMe(in *pb.InternalMeRequest) error {
	if err := sanitizer.SanitizeToken(in.GetToken()); err != nil { // email
		return err
	}

	return nil
}

func (s *AuthServer) InternalMe(ctx context.Context, in *pb.InternalMeRequest) (*pb.InternalMeResponse, error) {
	// -------------------- sanitize
	if err := sanitizeInternalMe(in); err != nil {
		return nil, err
	}

	// -------------------- cache
	token_info, err := databases.RetrieveToken(in.GetToken())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "connection failed")
	}

	// -------------------- db
	user, err := databases.DBs.SqlcQueries.GetUserById(context.Background(), int64(token_info.Id))

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, status.Errorf(codes.Internal, "token is not linked to any user")
	}

	return &pb.InternalMeResponse{
		Id:        token_info.Id,
		Email:     user.Email,
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		External:  token_info.External,
	}, nil
}
