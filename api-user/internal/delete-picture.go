package internal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/trixky/hypertube/.shared/databases"
	"github.com/trixky/hypertube/.shared/utils"
	pb "github.com/trixky/hypertube/api-user/proto"
	"github.com/trixky/hypertube/api-user/queries"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserServer) DeletePicture(ctx context.Context, in *pb.DeletePictureRequest) (*pb.DeletePictureResponse, error) {
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

	// -------------------- delete file
	err = os.Remove(StoragePath + fmt.Sprint(user.ID) + "." + user.Extension.String)
	if err != nil {
		// Only log the error, there is valid reasons that the delete could fail
		log.Println("Error delete user", user.ID, "picture", err)
	}

	// -------------------- delete the reference
	if err := queries.SqlcQueries.DeleteUserPicture(context.Background(), token_info.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete picture in the database")
	}

	return &pb.DeletePictureResponse{}, nil
}
