package internal

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/trixky/hypertube/api-media/databases"
	"github.com/trixky/hypertube/api-media/finder"
	pb "github.com/trixky/hypertube/api-media/proto"
	"github.com/trixky/hypertube/api-media/utils"
)

func (s *MediaServer) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	if err := utils.RequireLogin(ctx); err != nil {
		return nil, err
	}

	user_locale := utils.GetLocale(ctx)

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
		column := strings.ToLower(*in.SortBy)
		if column == "id" || column == "name" || column == "year" || column == "duration" {
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
		thumbnail := media.Thumbnail.String
		pb_media := pb.Media{
			Id:          uint32(media.ID),
			Type:        pb.MediaCategory_CATEGORY_MOVIE,
			Description: media.Description.String,
			Year:        uint32(media.Year.Int32),
			Duration:    &media.Duration.Int32,
			Names:       make([]*pb.MediaName, 0),
			Genres:      make([]string, 0),
			Thumbnail:   &thumbnail,
			Rating:      &rating,
		}

		// Load relations
		names, err := databases.DBs.SqlcQueries.GetMediaNames(ctx, int32(media.ID))
		if err != nil {
			return nil, err
		}
		for _, name := range names {
			if utils.NameMatchLocale(&user_locale, name.Lang) {
				pb_media.Names = append(pb_media.Names, &pb.MediaName{
					Lang:  name.Lang,
					Title: name.Name,
				})
			}
		}
		genres, err := databases.DBs.SqlcQueries.GetMediaGenres(ctx, int32(media.ID))
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
