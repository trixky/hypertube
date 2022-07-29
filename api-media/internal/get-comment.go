package internal

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/trixky/hypertube/api-media/proto"
	"github.com/trixky/hypertube/api-media/queries"
	"github.com/trixky/hypertube/api-media/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *MediaServer) GetComment(ctx context.Context, in *pb.GetCommentRequest) (*pb.GlobalComment, error) {
	_, err := utils.RequireLogin(ctx)
	if err != nil {
		return nil, err
	}

	// Check if the comment exists
	comment, err := queries.SqlcQueries.GetCommentById(ctx, int64(in.CommentId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "no comment with this id")
		} else {
			log.Println(err)
			return nil, err
		}
	}

	created_at, _ := ptypes.TimestampProto(comment.CreatedAt.Time)
	return &pb.GlobalComment{
		Id:      uint64(comment.ID.Int64),
		MediaId: uint64(comment.MediaID.Int32),
		User: &pb.CommentUser{
			Id:   uint64(comment.UserID.Int32),
			Name: comment.Username,
		},
		Content: comment.Content.String,
		Date:    created_at,
	}, nil
}
