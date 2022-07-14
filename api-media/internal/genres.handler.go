package internal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/trixky/hypertube/api-media/databases"
	pb "github.com/trixky/hypertube/api-media/proto"
	"github.com/trixky/hypertube/api-media/utils"
)

func (s *MediaServer) Genres(ctx context.Context, in *pb.GenresRequest) (*pb.GenresResponse, error) {
	if err := utils.RequireLogin(ctx); err != nil {
		return nil, err
	}

	fmt.Println("List Genres", in)

	genres, err := databases.DBs.SqlcQueries.GetGenres(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &pb.GenresResponse{}, err
		} else {
			err = nil
		}
	}
	response := pb.GenresResponse{
		Genres: make([]*pb.Genre, 0, len(genres)),
	}
	for _, genre := range genres {
		response.Genres = append(response.Genres, &pb.Genre{
			Id:   int32(genre.ID),
			Name: genre.Name,
		})
	}

	return &response, nil
}
