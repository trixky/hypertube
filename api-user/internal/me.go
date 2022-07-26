package internal

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/trixky/hypertube/.shared/databases"
	"github.com/trixky/hypertube/.shared/sanitizer"
	"github.com/trixky/hypertube/.shared/utils"
	pb "github.com/trixky/hypertube/api-user/proto"
	"github.com/trixky/hypertube/api-user/queries"
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
	fmt.Println("@@@@@@@@@@ 1")
	// -------------------- sanitize
	if err := sanitizeGetMe(in); err != nil {
		return nil, err
	}
	fmt.Println("@@@@@@@@@@ 2")

	// -------------------- cache
	token_info, err := databases.RetrieveToken(in.GetToken())

	fmt.Println("@@@@@@@@@@ 3")
	if err != nil {

		return nil, status.Errorf(codes.Unauthenticated, "token retrieving failed")
	}
	fmt.Println("@@@@@@@@@@ 4: ", token_info.Id)

	// -------------------- db
	user, err := queries.SqlcQueries.GetUserById(context.Background(), token_info.Id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "user infos retrieving failed")
	}

	fmt.Println("@@@@@@@@@@ 5")
	me, err := utils.HeaderCookieUserGeneration(utils.User{
		Id:        int(user.ID),
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		External:  token_info.External,
	}, false)
	fmt.Println("@@@@@@@@@@ 6")

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cookie generation failed")
	}
	fmt.Println("@@@@@@@@@@ 7")

	return &pb.UserInfoResponse{
		UserInfo: base64.StdEncoding.EncodeToString([]byte(me.Value)),
	}, nil
}
