package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/trixky/hypertube/api-search/postgres"
	pb "github.com/trixky/hypertube/api-search/proto"
	grpcMetadata "google.golang.org/grpc/metadata"
)

func (s *SearchServer) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	md, ok := grpcMetadata.FromIncomingContext(ctx)

	if !ok {
		log.Println("missing args")
		return nil, nil
	}

	search := md.Get("search")
	fmt.Println("search:", search)

	return &pb.SearchResponse{
		Page:   1,
		Medias: []*pb.Media{},
	}, nil
}

func (s *SearchServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	md, ok := grpcMetadata.FromIncomingContext(ctx)

	if !ok {
		log.Println("missing args")
		return nil, nil
	}

	search := md.Get("get")
	fmt.Println("get:", search)

	var duration int32 = 0
	var thumbnail string = "https://www.themoviedb.org/t/p/w300_and_h450_bestv2/yFwFp5QVHazxTklKGiJ0G59pVab.jpg"
	var rating float32 = 7.2
	return &pb.GetResponse{
		Media: &pb.Media{
			Id:          1,
			Type:        pb.MediaCategory_CATEGORY_MOVIE,
			Description: "Movie",
			Year:        2000,
			Names: []*pb.MediaName{
				{
					Lang:  "FR",
					Title: "Movie",
				},
			},
			Genres:    []string{"Action"},
			Duration:  &duration,
			Thumbnail: &thumbnail,
			Rating:    &rating,
		},
	}, nil
}

func (s *SearchServer) Genres(ctx context.Context, in *pb.GenresRequest) (*pb.GenresResponse, error) {
	fmt.Println("List Genres", in)

	genres, err := postgres.DB.SqlcQueries.GetGenres(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
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
