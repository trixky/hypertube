package internal

import (
	"context"
	"encoding/base64"

	"github.com/trixky/hypertube/api-user/databases"
	pb "github.com/trixky/hypertube/api-user/proto"
	"github.com/trixky/hypertube/api-user/sanitizer"
	"github.com/trixky/hypertube/api-user/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func sanitizeMe(in *pb.GetMeRequest) error {
	if err := sanitizer.SanitizeToken(in.GetToken()); err != nil { // email
		return err
	}

	return nil
}

func (s *UserServer) GetMe(ctx context.Context, in *pb.GetMeRequest) (*pb.MeResponse, error) {
	// -------------------- sanitize
	if err := sanitizeMe(in); err != nil {
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

	me, err := utils.HeaderCookieMeGeneration(utils.CookieMe{
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

	return &pb.MeResponse{
		Me: base64.StdEncoding.EncodeToString([]byte(me.Value)),
	}, nil
}
