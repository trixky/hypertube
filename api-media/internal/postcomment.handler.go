package internal

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/trixky/hypertube/api-media/databases"
	pb "github.com/trixky/hypertube/api-media/proto"
	"github.com/trixky/hypertube/api-media/sqlc"
	"github.com/trixky/hypertube/api-media/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *MediaServer) PostComment(ctx context.Context, in *pb.PostCommentRequest) (*pb.PostCommentResponse, error) {
	user, err := utils.RequireLogin(ctx)
	if err != nil {
		return nil, err
	}

	// Check if media exists
	media, err := databases.DBs.SqlcQueries.GetMediaByID(ctx, int64(in.MediaId))
	if err != nil {
		return nil, err
	}

	// Sanitize comment, only check the length
	if len(in.Content) < 2 || len(in.Content) > 500 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid comment content length")
	}

	// Add the comment
	comment, err := databases.DBs.SqlcQueries.CreateComment(ctx, sqlc.CreateCommentParams{
		UserID:  int32(user.ID),
		MediaID: int32(media.ID),
		Content: in.Content,
	})
	if err != nil {
		return nil, err
	}

	created_at, _ := ptypes.TimestampProto(comment.CreatedAt)
	return &pb.PostCommentResponse{
		Id: uint64(comment.ID),
		User: &pb.CommentUser{
			Id:   uint64(user.ID),
			Name: user.Username,
		},
		Content: in.Content,
		Date:    created_at,
	}, nil
}
