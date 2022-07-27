package internal

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/trixky/hypertube/.shared/databases"
	"github.com/trixky/hypertube/.shared/utils"
	pb "github.com/trixky/hypertube/api-user/proto"
	"github.com/trixky/hypertube/api-user/queries"
	"github.com/trixky/hypertube/api-user/sqlc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserServer) GetUserMovies(ctx context.Context, in *pb.GetUserMoviesRequest) (*pb.MediaList, error) {
	user_locale := utils.GetLocale(ctx)
	page := in.GetPage()

	// -------------------- get token
	sanitized_token, err := utils.ExtractSanitizedTokenFromGrpcGatewayCookies("", ctx)

	if err != nil {
		return nil, err
	}

	// -------------------- cache
	if _, err := databases.RetrieveToken(sanitized_token); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token retrieving failed")
	}

	// -------------------- get self user
	token_info, err := databases.RetrieveToken(sanitized_token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token retrieving failed")
	}
	self, err := queries.SqlcQueries.GetUserById(context.Background(), token_info.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "user infos retrieving failed")
	}

	// -------------------- get requested user
	user, err := queries.SqlcQueries.GetUserById(context.Background(), in.GetUserId())

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "no user found with this id")
		}
		return nil, status.Errorf(codes.Internal, "user infos retrieving failed")
	}

	// -------------------- count medias
	count, err := queries.SqlcQueries.CountUserMedias(ctx, int32(user.ID))
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return &pb.MediaList{}, nil
	}

	// -------------------- get the medias
	var offset uint32 = 0
	if page > 1 {
		offset = (page - 1) * 50
	}
	medias, err := queries.SqlcQueries.GetUserMedias(ctx, sqlc.GetUserMediasParams{
		UserID: int32(user.ID),
		Offset: int32(offset),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &pb.MediaList{}, nil
		} else {
			return nil, err
		}
	}

	// -------------------- get medias informations
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
			UserID:  int32(self.ID),
			MediaID: utils.MakeNullInt32(&media_id),
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
			if utils.NameMatchLocale(&user_locale, name.Lang) {
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
		TotalResults: uint32(count),
		Medias:       pb_medias,
	}

	return &response, nil
}
