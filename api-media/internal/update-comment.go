package internal

import (
	"context"
	"database/sql"
	"errors"
	"log"

	pb "github.com/trixky/hypertube/api-media/proto"
	"github.com/trixky/hypertube/api-media/queries"
	"github.com/trixky/hypertube/api-media/sqlc"
	"github.com/trixky/hypertube/api-media/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *MediaServer) UpdateComment(ctx context.Context, in *pb.UpdateCommentRequest) (*pb.UpdateCommentResponse, error) {
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
			return nil, status.Errorf(codes.Internal, "failed to get comment")
		}
	}

	// Check if the comment was posted by the user
	if comment.UserID.Int32 != int32(user.ID) {
		return nil, status.Errorf(codes.PermissionDenied, "you can only update your own comments")
	}

	// Sanitize comment, only check the length
	if len(in.Content) < 2 || len(in.Content) > 500 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid comment content length")
	}

	// Add the comment
	err = queries.SqlcQueries.UpdateComment(ctx, sqlc.UpdateCommentParams{
		ID:      comment.ID.Int64,
		Content: in.Content,
	})
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to update comment")
	}

	return &pb.UpdateCommentResponse{}, nil
}
