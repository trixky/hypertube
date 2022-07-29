package internal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/trixky/hypertube/.shared/databases"
	"github.com/trixky/hypertube/.shared/utils"
	pb "github.com/trixky/hypertube/api-user/proto"
	"github.com/trixky/hypertube/api-user/queries"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserServer) UploadPicture(ctx context.Context, in *pb.UploadPictureRequest) (*pb.UploadPictureResponse, error) {
	// -------------------- get token
	sanitized_token, err := utils.ExtractSanitizedTokenFromGrpcGatewayCookies("", ctx)

	if err != nil {
		return nil, err
	}

	// -------------------- cache
	token_info, err := databases.RetrieveToken(sanitized_token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token retrieving failed")
	}

	// -------------------- db
	user, err := queries.SqlcQueries.GetUserById(context.Background(), token_info.Id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "no user found with this id")
		}
		return nil, status.Errorf(codes.Internal, "user infos retrieving failed")
	}

	// -------------------- check uploaded file validity
	// -------------------- Save file in storage (shared volume)
	path := StoragePath + fmt.Sprint(user.ID) + "." + "png"
	err = os.WriteFile(path, in.GetPicture(), 0)

	return &pb.UploadPictureResponse{}, nil
}
