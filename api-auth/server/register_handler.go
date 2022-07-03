package server

import (
	"context"
	"fmt"
	"log"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/trixky/hypertube/api-auth/databases"
	pb "github.com/trixky/hypertube/api-auth/proto"
	"github.com/trixky/hypertube/api-auth/sqlc"
	"google.golang.org/grpc/codes"
	grpcMetadata "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const expiration = time.Hour * 24 * 1000

func SanitizeInternalRegister(in *pb.InternalRegisterRequest) error {
	// ------------------------ email
	email := in.GetEmail()
	if len(email) == 0 {
		return status.Errorf(codes.InvalidArgument, "email missing")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return status.Errorf(codes.InvalidArgument, "email corrupted")
	}

	// ------------------------ username
	username := in.GetUsername()
	if len(username) == 0 {
		return status.Errorf(codes.InvalidArgument, "username missing")
	}
	if len(username) < 3 {
		return status.Errorf(codes.InvalidArgument, "username too short")
	}
	if len(username) > 30 {
		return status.Errorf(codes.InvalidArgument, "username too long")
	}

	// ------------------------ firstname
	firstname := in.GetFirstname()
	if len(firstname) == 0 {
		return status.Errorf(codes.InvalidArgument, "firstname missing")
	}
	if len(firstname) < 3 {
		return status.Errorf(codes.InvalidArgument, "firstname too short")
	}
	if len(firstname) > 30 {
		return status.Errorf(codes.InvalidArgument, "firstname too long")
	}

	// ------------------------ lastname
	lastname := in.GetLastname()
	if len(lastname) == 0 {
		return status.Errorf(codes.InvalidArgument, "lastname missing")
	}
	if len(lastname) < 3 {
		return status.Errorf(codes.InvalidArgument, "lastname too short")
	}
	if len(lastname) > 20 {
		return status.Errorf(codes.InvalidArgument, "lastname too long")
	}

	// ------------------------ password
	password := in.GetPassword()
	if len(password) == 0 {
		return status.Errorf(codes.InvalidArgument, "password missing")
	}
	if len(password) < 8 {
		return status.Errorf(codes.InvalidArgument, "password too short")
	}
	if len(password) > 30 {
		return status.Errorf(codes.InvalidArgument, "password too long")
	}

	const numeric = "0123456789"

	if !strings.ContainsAny(password, numeric) {
		return status.Errorf(codes.InvalidArgument, "password must contain at least one number")
	}

	const alpha = "qwertyuioplkjhgfdsazxcvbnm" // qwerty

	if !strings.ContainsAny(password, alpha) {
		return status.Errorf(codes.InvalidArgument, "password must contain at least one lowercase character")
	}
	if !strings.ContainsAny(password, strings.ToUpper(alpha)) {
		return status.Errorf(codes.InvalidArgument, "password must contain at least one uppercase character")
	}

	const special = " !@#$%^&*()-=_+[]{}\\|'\";:/?.>,<`~"

	if !strings.ContainsAny(password, special) {
		return status.Errorf(codes.InvalidArgument, "password must contain at least one special character")
	}

	return nil
}

func (s *AuthServer) InternalRegister(ctx context.Context, in *pb.InternalRegisterRequest) (*pb.InternalLoginResponse, error) {
	log.Printf("Greet function awas invoked with %v\n", in)

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
	log.Println("yolooooooooooo 0")
	if err := SanitizeInternalRegister(in); err != nil {
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
	// -------------------- logic
	token := uuid.New().String() // token generation

	// -------------------- cache
	if err := databases.AddToken(new_user.ID, token); err != nil {
		return nil, status.Errorf(codes.Internal, "token creation failed")
	}

	return &pb.InternalLoginResponse{
		Token: token,
	}, nil
}
