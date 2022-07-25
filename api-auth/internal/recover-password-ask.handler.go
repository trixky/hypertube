package internal

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/trixky/hypertube/.shared/sanitizer"
	"github.com/trixky/hypertube/api-auth/email"
	pb "github.com/trixky/hypertube/api-auth/proto"
	"github.com/trixky/hypertube/api-auth/queries"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// sanitizeInternalRecoverPassword sanitizes inputs for "InternalRecoverPassword" route
func sanitizeInternalRecoverPassword(in *pb.InternalRecoverPasswordRequest) error {
	if err := sanitizer.SanitizeEmail(in.GetEmail()); err != nil { // email
		return err
	}

	return nil
}

// InternalRecoverPassword Handles the "InternalRecoverPassword" route
func (s *AuthServer) InternalRecoverPassword(ctx context.Context, in *pb.InternalRecoverPasswordRequest) (*empty.Empty, error) {

	_empty := new(empty.Empty)

	// -------------------- Sanitize
	if err := sanitizeInternalRecoverPassword(in); err != nil {
		return _empty, err
	}

	// -------------------- DB
	user, err := queries.SqlcQueries.GetInternalUserByEmail(context.Background(), in.GetEmail())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return _empty, status.Errorf(codes.NotFound, "no user finded with this email")
		}
		return _empty, status.Errorf(codes.Internal, "connection failed")
	}

	// -------------------- Cache
	password_token := uuid.New().String() // token generation

	if err := queries.AddPasswordToken(user.ID, password_token); err != nil {
		return _empty, status.Errorf(codes.Internal, "password token generation failed")
	}

	// -------------------- Email
	go func() {
		if err := email.SendPasswordToken(in.GetEmail(), password_token); err != nil {
			log.Printf("sending recover email to [%s] failed: %s\n", in.GetEmail(), err.Error())
		}
	}()

	return _empty, nil
}
