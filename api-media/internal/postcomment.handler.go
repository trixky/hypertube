package internal

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/trixky/hypertube/api-media/databases"
	pb "github.com/trixky/hypertube/api-media/proto"
	"github.com/trixky/hypertube/api-media/sqlc"
	"github.com/trixky/hypertube/api-media/utils"
)

func (s *MediaServer) PostComment(ctx context.Context, in *pb.PostCommentRequest) (*pb.PostCommentResponse, error) {
	if err := utils.RequireLogin(ctx); err != nil {
		return nil, err
	}

	// Check if media exists
	media, err := databases.DBs.SqlcQueries.GetMediaByID(ctx, int64(in.MediaId))
	if err != nil {
		return nil, err
	}

	// Add the comment
	comment, err := databases.DBs.SqlcQueries.CreateComment(ctx, sqlc.CreateCommentParams{
		UserID:  1,
		MediaID: int32(media.ID),
		Content: in.Content,
	})
	if err != nil {
		return nil, err
	}

	created_at, _ := ptypes.TimestampProto(comment.CreatedAt)
	return &pb.PostCommentResponse{
		Id:   uint64(comment.ID),
		Date: created_at,
	}, nil
}
