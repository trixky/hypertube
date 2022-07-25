package internal

import (
	"context"
	"database/sql"
	"encoding/base64"

	"github.com/trixky/hypertube/.shared/sanitizer"
	"github.com/trixky/hypertube/.shared/utils"
	"github.com/trixky/hypertube/api-user/databases"
	pb "github.com/trixky/hypertube/api-user/proto"
	"github.com/trixky/hypertube/api-user/sqlc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserServer) UpdateMe(ctx context.Context, in *pb.UpdateMeRequest) (*pb.UserInfoResponse, error) {
	// -------------------- get token
	sanitized_token, err := utils.ExtractSanitizedTokenFromGrpcGatewayCookies(in.GetToken(), ctx)

	if err != nil {
		return nil, err
	}

	// -------------------- cache
	token_info, err := databases.RetrieveToken(sanitized_token)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token retrieving failed")
	}

	// -------------------- db get
	user, err := databases.DBs.SqlcQueries.GetUserById(context.Background(), token_info.Id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "user infos retrieving failed")
	}

	// -------------------------------------------- update
	need_update := false
	// -------------------- update username
	username := in.GetUsername()

	if len(username) > 0 {
		if err := sanitizer.SanitizeUsername(username); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "username corrupted")
		}
		need_update = true

		user.Username = username
	}

	// -------------------- update firstname
	firstname := in.GetFirstname()

	if len(firstname) > 0 {
		if err := sanitizer.SanitizeFirstname(firstname); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "firstname corrupted")
		}
		need_update = true
		user.Firstname = firstname
	}

	// -------------------- update lastname
	lastname := in.GetLastname()

	if len(lastname) > 0 {
		if err := sanitizer.SanitizeLastname(lastname); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "lastname corrupted")
		}
		need_update = true
		user.Lastname = lastname
	}

	// -------------------- update credentials
	if token_info.External == databases.EXTERNAL_none {
		// -------------------- update email
		email := in.GetEmail()

		if len(email) > 0 {
			if err := sanitizer.SanitizeEmail(email); err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "email corrupted")
			}
			need_update = true
			user.Email = email
		}
		// -------------------- update password

		new_password := in.GetNewPassword()
		new_password_encrypted := utils.EncryptPassword(in.GetNewPassword())

		if len(new_password) > 0 {
			if err := sanitizer.SanitizeSHA256Password(new_password); err != nil { // password
				return nil, err
			}
			current_password := in.GetCurrentPassword()
			current_password_encrypted := utils.EncryptPassword(current_password)

			if current_password_encrypted != user.Password.String {
				return nil, status.Errorf(codes.PermissionDenied, "invalid current password")
			}

			if err := sanitizer.SanitizeSHA256Password(new_password); err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "new password corrupted")
			}

			need_update = true
			user.Password = sql.NullString{
				String: new_password_encrypted,
				Valid:  true,
			}
		}
	}

	// -------------------- db update
	if need_update {
		if err := databases.DBs.SqlcQueries.UpdateUser(context.Background(), sqlc.UpdateUserParams{
			ID:        user.ID,
			Username:  user.Username,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Email:     user.Email,
			Password:  user.Password,
		}); err != nil {
			if databases.ErrorIsDuplication(err) {
				return nil, status.Errorf(codes.AlreadyExists, "email is already in use")
			}
			return nil, status.Errorf(codes.Internal, "update failed")
		}
	}

	me, err := utils.HeaderCookieUserGeneration(utils.User{
		Id:        int(user.ID),
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		External:  token_info.External,
	}, false)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cookie generation failed")
	}

	return &pb.UserInfoResponse{
		UserInfo: base64.StdEncoding.EncodeToString([]byte(me.Value)),
	}, nil
}
