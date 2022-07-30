package internal

import (
	"context"
	"database/sql"
	"errors"
	"log"

	pb "github.com/trixky/hypertube/api-media/proto"
	"github.com/trixky/hypertube/api-media/queries"
	"github.com/trixky/hypertube/api-media/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *MediaServer) DeleteComment(ctx context.Context, in *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	user, err := utils.RequireLogin(ctx)
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

	// Check if the comment was posted by the user
	if comment.UserID.Int32 != int32(user.ID) {
		return nil, status.Errorf(codes.PermissionDenied, "you can only delete your own comments")
	}

	// Delete the comment
	if err := queries.SqlcQueries.DeleteComment(ctx, comment.ID.Int64); err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.DeleteCommentResponse{}, nil
}
