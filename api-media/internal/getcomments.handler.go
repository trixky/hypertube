package internal

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/golang/protobuf/ptypes"
	"github.com/trixky/hypertube/api-media/databases"
	pb "github.com/trixky/hypertube/api-media/proto"
	grpcMetadata "google.golang.org/grpc/metadata"
)

func (s *MediaServer) GetComments(ctx context.Context, in *pb.GetCommentsRequest) (*pb.GetCommentsResponse, error) {
	md, ok := grpcMetadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("missing args")
		return nil, nil
	}

	getComments := md.Get("getComments")
	log.Println("getComments:", getComments)

	// Find the media
	media, err := databases.DBs.SqlcQueries.GetMediaByID(ctx, int64(in.MediaId))
	if err != nil {
		return nil, err
	}

	// Load comments
	comments, err := databases.DBs.SqlcQueries.GetMediaComments(ctx, int32(media.ID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &pb.GetCommentsResponse{}, nil
		} else {
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
