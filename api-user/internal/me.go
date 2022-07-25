package internal

import (
	"context"
	"encoding/base64"

	"github.com/trixky/hypertube/.shared/sanitizer"
	"github.com/trixky/hypertube/.shared/utils"
	"github.com/trixky/hypertube/api-user/databases"
	pb "github.com/trixky/hypertube/api-user/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func sanitizeGetMe(in *pb.GetMeRequest) error {
	if err := sanitizer.SanitizeToken(in.GetToken()); err != nil { // email
		return err
	}

	return nil
}

func (s *UserServer) GetMe(ctx context.Context, in *pb.GetMeRequest) (*pb.UserInfoResponse, error) {
	// -------------------- sanitize
	if err := sanitizeGetMe(in); err != nil {
		return nil, err
	}

	// -------------------- cache
	token_info, err := databases.RetrieveToken(in.GetToken())

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token retrieving failed")
	}

	// -------------------- db
	user, err := databases.DBs.SqlcQueries.GetUserById(context.Background(), token_info.Id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "user infos retrieving failed")
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
