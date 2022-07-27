package internal

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/trixky/hypertube/.shared/databases"
	sutils "github.com/trixky/hypertube/.shared/utils"
	pb "github.com/trixky/hypertube/api-media/proto"
	"github.com/trixky/hypertube/api-media/queries"
	"github.com/trixky/hypertube/api-media/sqlc"
	"github.com/trixky/hypertube/api-media/utils"
	"google.golang.org/protobuf/encoding/protojson"
)

func (s *MediaServer) Search(ctx context.Context, in *pb.SearchRequest) (*pb.MediaList, error) {
	user, err := utils.RequireLogin(ctx)
	if err != nil {
		return nil, err
	}

	user_locale := sutils.GetLocale(ctx)

	// Check and set arguments for the query
	params := utils.FindMediasParams{
		SortColumn: "id",
		SortOrder:  "DESC",
	}
	page := int32(1)
	if in.Page != nil && *in.Page > 1 {
		page = int32(*in.Page)
		params.Offset = (page - 1) * utils.PerPage
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

	// Check cache
	path := params.ToString(user_locale.Lang, databases.REDIS_SEPARATOR)
	cache_results, err := queries.RetrieveSearch(&path)
	if err != nil {
		log.Println("error in redis cache", err)
	} else if cache_results != "" {
		response := pb.MediaList{}
		err = protojson.Unmarshal([]byte(cache_results), &response)
		if err != nil {
			log.Println("error in redis cache unmarshal", err)
		} else {
			return &response, nil
		}
	}

	// Count the Medias first
	medias_count, err := utils.CountMedias(ctx, params)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		return nil, err
	}
	if medias_count == 0 {
		response := pb.MediaList{}
		search := protojson.Format(response.ProtoReflect().Interface())
		err = queries.CacheSearch(&path, &search)
		if err != nil {
			log.Println("failed to save to redis cache", err)
		}
		return &response, nil
	}

	// Find all Medias
	medias, err := utils.FindMedias(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response := pb.MediaList{}
			search := protojson.Format(response.ProtoReflect().Interface())
			err = queries.CacheSearch(&path, &search)
			if err != nil {
				log.Println("failed to save to redis cache", err)
			}
			return &response, nil
		}
		log.Println(err)
		return nil, err
	}

	// Convert to proto
	pb_medias := make([]*pb.Media, 0)
	for _, media := range medias {
		media_id := int32(media.ID)
		rating := float32(media.Rating.Float64)
		thumbnail := media.Thumbnail.String
		pb_media := pb.Media{
			Id:        uint32(media.ID),
			Type:      pb.MediaCategory_CATEGORY_MOVIE,
			Year:      uint32(media.Year.Int32),
			Names:     make([]*pb.MediaName, 0),
			Thumbnail: &thumbnail,
			Rating:    &rating,
			Watched:   true,
		}

		// Check watched status
		// -- at least one torrent has a position
		_, err := queries.SqlcQueries.WatchedMedia(ctx, sqlc.WatchedMediaParams{
			UserID:  int32(user.ID),
			MediaID: sutils.MakeNullInt32(&media_id),
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				pb_media.Watched = false
			} else {
				log.Println(err)
				return nil, err
			}
		}

		// Load names
		names, err := queries.SqlcQueries.GetMediaNames(ctx, int32(media.ID))
		if err != nil {
			log.Println(err)
			return nil, err
		}
		for _, name := range names {
			if sutils.NameMatchLocale(&user_locale, name.Lang) {
				pb_media.Names = append(pb_media.Names, &pb.MediaName{
					Lang:  name.Lang,
					Title: name.Name,
				})
			}
		}

		// Add everything to the response
		pb_medias = append(pb_medias, &pb_media)
	}

	response := pb.MediaList{
		Page:         uint32(page),
		Results:      uint32(len(pb_medias)),
		TotalResults: uint32(medias_count),
		Medias:       pb_medias,
	}

	// Save in redis
	search := protojson.Format(response.ProtoReflect().Interface())
	err = queries.CacheSearch(&path, &search)
	if err != nil {
		log.Println("failed to save to redis cache", err)
	}

	return &response, nil
}
