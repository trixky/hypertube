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

func (s *MediaServer) GetAllComments(ctx context.Context, in *pb.GetAllCommentsRequest) (*pb.GetAllCommentsResponse, error) {
	_, err := utils.RequireLogin(ctx)
	if err != nil {
		return nil, err
	}

	// Get the last 15 comments
	comments, err := queries.SqlcQueries.GetLastComments(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &pb.GetAllCommentsResponse{}, nil
		} else {
			log.Println(err)
			return nil, status.Errorf(codes.Internal, "failed to load comments")
		}
	}

	// Convert comments
	response := pb.GetAllCommentsResponse{
		TotalComments: uint32(len(comments)),
		Comments:      make([]*pb.GlobalComment, 0),
	}
	for _, comment := range comments {
		created_at, _ := ptypes.TimestampProto(comment.CreatedAt.Time)
		response.Comments = append(response.Comments, &pb.GlobalComment{
			Id:      uint64(comment.ID.Int64),
			MediaId: uint64(comment.MediaID.Int32),
			User: &pb.CommentUser{
				Id:   uint64(comment.UserID.Int32),
				Name: comment.Username,
			},
			Content: comment.Content.String,
			Date:    created_at,
		})
	}

	return &response, nil
}
