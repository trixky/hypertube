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

func (s *MediaServer) GetComments(ctx context.Context, in *pb.GetCommentsRequest) (*pb.GetCommentsResponse, error) {
	if _, err := utils.RequireLogin(ctx); err != nil {
		return nil, err
	}

	// Find the media
	media, err := queries.SqlcQueries.GetMediaByID(ctx, int64(in.MediaId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "no media with this id")
		} else {
			log.Println(err)
			return nil, err
		}
	}

	// Load comments
	comments, err := queries.SqlcQueries.GetMediaComments(ctx, int32(media.ID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &pb.GetCommentsResponse{}, nil
		} else {
			log.Println(err)
			return nil, err
		}
	}

	// Convert comments
	response := pb.GetCommentsResponse{
		TotalComments: uint32(len(comments)),
		Comments:      make([]*pb.Comment, 0),
	}
	for _, comment := range comments {
		created_at, _ := ptypes.TimestampProto(comment.CreatedAt.Time)
		response.Comments = append(response.Comments, &pb.Comment{
			Id: uint64(comment.ID.Int64),
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
