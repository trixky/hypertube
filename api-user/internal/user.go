package internal

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"

	"github.com/trixky/hypertube/.shared/databases"
	"github.com/trixky/hypertube/.shared/utils"
	pb "github.com/trixky/hypertube/api-user/proto"
	"github.com/trixky/hypertube/api-user/queries"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserServer) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.UserInfoResponse, error) {
	// -------------------- get token
	sanitized_token, err := utils.ExtractSanitizedTokenFromGrpcGatewayCookies("", ctx)

	if err != nil {
		return nil, err
	}

	// -------------------- cache
	if _, err := databases.RetrieveToken(sanitized_token); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token retrieving failed")
	}

	// -------------------- db
	user, err := queries.SqlcQueries.GetUserById(context.Background(), in.GetId())

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "no user found with this id")
		}
		return nil, status.Errorf(codes.Internal, "user infos retrieving failed")
	}

	_user, err := utils.HeaderCookieUserGeneration(utils.User{
		Id:        int(user.ID),
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
	}, false)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cookie generation failed")
	}

	return &pb.UserInfoResponse{
		UserInfo: base64.StdEncoding.EncodeToString([]byte(_user.Value)),
	}, nil
}
