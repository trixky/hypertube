package server

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/trixky/hypertube/api-auth/databases"
	pb "github.com/trixky/hypertube/api-auth/proto"
	"github.com/trixky/hypertube/api-auth/sanitizer"
	"github.com/trixky/hypertube/api-auth/sqlc"
	"google.golang.org/grpc/codes"
	grpcMetadata "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func sanitizeInternalRegister(in *pb.InternalRegisterRequest) error {
	if err := sanitizer.SanitizeEmail(in.GetEmail()); err != nil { // email
		return err
	}
	if err := sanitizer.SanitizeName(in.GetUsername(), sanitizer.LABEL_username); err != nil { // username
		return err
	}
	if err := sanitizer.SanitizeName(in.GetFirstname(), sanitizer.LABEL_firstname); err != nil { // firstname
		return err
	}
	if err := sanitizer.SanitizeName(in.GetLastname(), sanitizer.LABEL_lastname); err != nil { // lastname
		return err
	}
	if err := sanitizer.SanitizeSHA256Password(in.GetPassword()); err != nil { // password
		return err
	}

	return nil
}

func (s *AuthServer) InternalRegister(ctx context.Context, in *pb.InternalRegisterRequest) (*pb.InternalLoginResponse, error) {
	// ----------------------- (temp) headers
	md, ok := grpcMetadata.FromIncomingContext(ctx)

	if !ok {
		log.Println("arg 111")
		return nil, nil
	}

	fmt.Println("method:", md.Get("method"))
	fmt.Println("pattern:", md.Get("pattern"))

	fmt.Println("email", in.GetEmail())
	fmt.Println("username", in.GetUsername())
	fmt.Println("firstname", in.GetFirstname())
	fmt.Println("lastname", in.GetLastname())
	fmt.Println("password", in.GetPassword())

	// -------------------- sanitize
	if err := sanitizeInternalRegister(in); err != nil {
		return nil, err
	}

	// -------------------- db
	new_user, err := databases.DBs.SqlcQueries.CreateUser(context.Background(), sqlc.CreateUserParams{
		Email:     in.GetEmail(),
		Username:  in.GetUsername(),
		Firstname: in.GetFirstname(),
		Lastname:  in.GetLastname(),
		Password:  in.GetPassword(),
	})

	if err != nil {
		if databases.ErrorIsDuplication(err) {
			return nil, status.Errorf(codes.AlreadyExists, "email is already in use")
		}
		return nil, status.Errorf(codes.Internal, "user creation failed")
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
