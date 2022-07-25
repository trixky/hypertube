package utils

import (
	"context"

	sutils "github.com/trixky/hypertube/.shared/utils"
	"github.com/trixky/hypertube/api-media/databases"
	"github.com/trixky/hypertube/api-media/sqlc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RequireLogin(ctx context.Context) (*sqlc.User, error) {
	// -------------------- get token
	sanitized_token, err := sutils.ExtractSanitizedTokenFromGrpcGatewayCookies("", ctx)
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

	return &user, nil
}
