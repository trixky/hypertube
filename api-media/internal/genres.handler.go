package internal

import (
	"context"
	"database/sql"
	"errors"
	"log"

	pb "github.com/trixky/hypertube/api-media/proto"
	"github.com/trixky/hypertube/api-media/queries"
	"github.com/trixky/hypertube/api-media/utils"
)

func (s *MediaServer) Genres(ctx context.Context, in *pb.GenresRequest) (*pb.GenresResponse, error) {
	if _, err := utils.RequireLogin(ctx); err != nil {
		return nil, err
	}

	genres, err := queries.SqlcQueries.GetGenres(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		} else {
			log.Println(err)
			return nil, err
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
