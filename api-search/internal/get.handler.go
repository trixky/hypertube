package internal

import (
	"context"
	"log"

	"github.com/trixky/hypertube/api-search/databases"
	pb "github.com/trixky/hypertube/api-search/proto"
	ut "github.com/trixky/hypertube/api-search/utils"
	grpcMetadata "google.golang.org/grpc/metadata"
)

func (s *SearchServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	md, ok := grpcMetadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("missing args")
		return nil, nil
	}

	get := md.Get("get")
	log.Println("get:", get)

	// Find the media
	media, err := databases.DBs.SqlcQueries.GetMediaByID(ctx, int64(in.Id))
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
	names, err := databases.DBs.SqlcQueries.GetMediaNames(ctx, int32(media.ID))
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		response.Media.Names = append(response.Media.Names, &pb.MediaName{
			Lang:  name.Lang,
			Title: name.Name,
		})
	}

	torrents, err := databases.DBs.SqlcQueries.GetMediaTorrents(ctx, ut.MakeNullInt32(&media_id))
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

	genres, err := databases.DBs.SqlcQueries.GetMediaGenres(ctx, int32(media.ID))
	if err != nil {
		return nil, err
	}
	for _, genre := range genres {
		response.Media.Genres = append(response.Media.Genres, genre.Name)
	}

	actors, err := databases.DBs.SqlcQueries.GetMediaActors(ctx, int32(media.ID))
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

	staffs, err := databases.DBs.SqlcQueries.GetMediaStaffs(ctx, int32(media.ID))
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
