package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/trixky/hypertube/api-search/finder"
	"github.com/trixky/hypertube/api-search/postgres"
	pb "github.com/trixky/hypertube/api-search/proto"
	ut "github.com/trixky/hypertube/api-search/utils"
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

	// Check and set arguments for the query
	params := finder.FindMediasParams{
		SortColumn: "id",
		SortOrder:  "ASC",
	}
	page := int32(1)
	if in.Page != nil && *in.Page > 1 {
		page = int32(*in.Page)
		params.Offset = (page - 1) * finder.PerPage
	}
	if in.Query != nil && *in.Query != "" {
		params.SearchQuery = true
		params.Query = *in.Query
	}
	if in.Year != nil && *in.Year > 0 {
		params.SearchYear = true
		params.Year = int32(*in.Year)
	}
	if in.Rating != nil && *in.Rating > 0 && *in.Rating <= 10 {
		params.SearchRating = true
		params.Rating = float64(*in.Rating)
	}
	if len(in.GenreIds) > 0 {
		params.SearchGenres = true
		params.Genres = in.GenreIds
	}
	if in.SortBy != nil {
		column := strings.ToLower(*in.SortOrder)
		if column == "id" || column == "year" || column == "duration" {
			params.SortColumn = column
		}
	}
	if in.SortOrder != nil {
		order := strings.ToUpper(*in.SortOrder)
		if order == "ASC" || order == "DESC" {
			params.SortOrder = order
		}
	}

	// Count the Medias first
	medias_count, err := finder.CountMedias(ctx, params)
	if err != nil {
		return nil, err
	}
	if medias_count == 0 {
		return &pb.SearchResponse{}, nil
	}

	// Find all Medias
	medias, err := finder.FindMedias(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &pb.SearchResponse{}, nil
		}
		return nil, err
	}

	// Convert to proto
	pb_medias := make([]*pb.Media, 0)
	for _, media := range medias {
		rating := float32(media.Rating.Float64)
		pb_media := pb.Media{
			Id:          uint32(media.ID),
			Type:        pb.MediaCategory_CATEGORY_MOVIE,
			Description: media.Description.String,
			Year:        uint32(media.Year.Int32),
			Duration:    &media.Duration.Int32,
			Names:       make([]*pb.MediaName, 0),
			Genres:      make([]string, 0),
			Thumbnail:   &media.Thumbnail.String,
			Rating:      &rating,
		}

		// Find some relations
		names, err := postgres.DB.SqlcQueries.GetMediaNames(ctx, int32(media.ID))
		if err != nil {
			return nil, err
		}
		for _, name := range names {
			pb_media.Names = append(pb_media.Names, &pb.MediaName{
				Lang:  name.Lang,
				Title: name.Name,
			})
		}
		genres, err := postgres.DB.SqlcQueries.GetMediaGenres(ctx, int32(media.ID))
		if err != nil {
			return nil, err
		}
		for _, genre := range genres {
			pb_media.Genres = append(pb_media.Genres, genre.Name)
		}

		// Add everything to the response
		pb_medias = append(pb_medias, &pb_media)
	}

	return &pb.SearchResponse{
		Page:         uint32(page),
		Results:      uint32(len(pb_medias)),
		TotalResults: uint32(medias_count),
		Medias:       pb_medias,
	}, nil
}

func (s *SearchServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	md, ok := grpcMetadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("missing args")
		return nil, nil
	}

	get := md.Get("get")
	fmt.Println("get:", get)

	// Find the media
	media, err := postgres.DB.SqlcQueries.GetMediaByID(ctx, int64(in.Id))
	if err != nil {
		return nil, err
	}

	// Construct the response
	media_id := int32(media.ID)
	rating := float32(media.Rating.Float64)
	response := pb.GetResponse{
		Media: &pb.Media{
			Id:          uint32(media.ID),
			Type:        pb.MediaCategory_CATEGORY_MOVIE,
			Description: media.Description.String,
			Year:        uint32(media.Year.Int32),
			Duration:    &media.Duration.Int32,
			Names:       make([]*pb.MediaName, 0),
			Genres:      make([]string, 0),
			Thumbnail:   &media.Thumbnail.String,
			Rating:      &rating,
		},
		Torrents: make([]*pb.TorrentPublicInformations, 0),
		Staffs:   make([]*pb.Staff, 0),
		Actors:   make([]*pb.Actor, 0),
	}

	// Find relations
	names, err := postgres.DB.SqlcQueries.GetMediaNames(ctx, int32(media.ID))
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		response.Media.Names = append(response.Media.Names, &pb.MediaName{
			Lang:  name.Lang,
			Title: name.Name,
		})
	}

	torrents, err := postgres.DB.SqlcQueries.GetMediaTorrents(ctx, ut.MakeNullInt32(&media_id))
	if err != nil {
		return nil, err
	}
	for _, torrent := range torrents {
		response.Torrents = append(response.Torrents, &pb.TorrentPublicInformations{
			Id:    int32(torrent.ID.Int64),
			Name:  torrent.Name.String,
			Seed:  &torrent.Seed.Int32,
			Leech: &torrent.Leech.Int32,
		})
	}

	genres, err := postgres.DB.SqlcQueries.GetMediaGenres(ctx, int32(media.ID))
	if err != nil {
		return nil, err
	}
	for _, genre := range genres {
		response.Media.Genres = append(response.Media.Genres, genre.Name)
	}

	actors, err := postgres.DB.SqlcQueries.GetMediaActors(ctx, int32(media.ID))
	if err != nil {
		return nil, err
	}
	for _, actor := range actors {
		response.Actors = append(response.Actors, &pb.Actor{
			Id:        int32(actor.ID),
			Name:      actor.Name,
			Thumbnail: actor.Thumbnail.String,
			Character: actor.Character.String,
		})
	}

	staffs, err := postgres.DB.SqlcQueries.GetMediaStaffs(ctx, int32(media.ID))
	if err != nil {
		return nil, err
	}
	for _, staff := range staffs {
		response.Staffs = append(response.Staffs, &pb.Staff{
			Id:        int32(staff.ID),
			Name:      staff.Name,
			Thumbnail: staff.Thumbnail.String,
			Role:      staff.Role.String,
		})
	}

	return &response, nil
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
