package internal

import (
	"context"
	"database/sql"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/trixky/hypertube/api-auth/databases"
	pb "github.com/trixky/hypertube/api-auth/proto"
	"github.com/trixky/hypertube/api-auth/sanitizer"
	"github.com/trixky/hypertube/api-auth/sqlc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func sanitizeInternalApplyRecoverPassword(in *pb.InternalApplyRecoverPasswordRequest) error {
	if err := sanitizer.SanitizeToken(in.GetPasswordToken()); err != nil { // email
		return err
	}
	if err := sanitizer.SanitizeSHA256Password(in.GetNewPassword()); err != nil { // password
		return err
	}

	return nil
}

func (s *AuthServer) InternalApplyRecoverPassword(ctx context.Context, in *pb.InternalApplyRecoverPasswordRequest) (*empty.Empty, error) {
	_empty := new(empty.Empty)
	// -------------------- sanitize
	if err := sanitizeInternalApplyRecoverPassword(in); err != nil {
		return _empty, err
	}

	// -------------------- cache
	password_token_info, err := databases.RetrievePasswordToken(in.GetPasswordToken(), true)
	if err != nil {
		return _empty, status.Errorf(codes.Internal, "password token generation failed")
	}

	// -------------------- db
	if err := databases.DBs.SqlcQueries.UpdateUserPassword(context.Background(), sqlc.UpdateUserPasswordParams{
		ID: password_token_info.Id,
		Password: sql.NullString{
			String: in.GetNewPassword(),
			Valid:  true,
		},
	}); err != nil {
		return _empty, status.Errorf(codes.Internal, "password update failed")
	}

	// -------------------- cache

	return _empty, nil
}
