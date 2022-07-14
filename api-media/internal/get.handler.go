package internal

import (
	"context"
	"sort"

	"github.com/golang/protobuf/ptypes"
	"github.com/trixky/hypertube/api-media/databases"
	pb "github.com/trixky/hypertube/api-media/proto"
	"github.com/trixky/hypertube/api-media/sqlc"
	"github.com/trixky/hypertube/api-media/utils"
	ut "github.com/trixky/hypertube/api-media/utils"
)

var StaffOrder []string = []string{
	"Director",
	"Writer",
	"Comic Book",
	"Story",
	"Original Story",
	"Author",
	"Animation",
	"Original Music Composer",
	"In Memory Of",
	"Producer",
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func indexOf(slice []string, value string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return 999
}

func (s *MediaServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	if err := utils.RequireLogin(ctx); err != nil {
		return nil, err
	}

	user_locale := ut.GetLocale(ctx)

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
			Background:  &media.Background.String,
			Rating:      &rating,
		},
		Torrents: make([]*pb.TorrentPublicInformations, 0),
		Staffs:   make([]*pb.Staff, 0),
		Actors:   make([]*pb.Actor, 0),
		Comments: make([]*pb.Comment, 0),
	}

	// Find relations
	names, err := databases.DBs.SqlcQueries.GetMediaNames(ctx, int32(media.ID))
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		if ut.NameMatchLocale(&user_locale, name.Lang) {
			response.Media.Names = append(response.Media.Names, &pb.MediaName{
				Lang:  name.Lang,
				Title: name.Name,
			})
		}
	}

	torrents, err := databases.DBs.SqlcQueries.GetMediaTorrents(ctx, ut.MakeNullInt32(&media_id))
	if err != nil {
		return nil, err
	}
	for _, torrent := range torrents {
		size := torrent.Size.String
		seed := torrent.Seed.Int32
		leech := torrent.Leech.Int32
		response.Torrents = append(response.Torrents, &pb.TorrentPublicInformations{
			Id:    int32(torrent.ID.Int64),
			Name:  torrent.Name.String,
			Size:  &size,
			Seed:  &seed,
			Leech: &leech,
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
	// Actors are limited to 15 and are already ordered by the `cast_order` column
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
	// Sort by "importance" before de-duplicating
	sort.SliceStable(staffs, func(i, j int) bool {
		return indexOf(StaffOrder, staffs[i].Role.String) < indexOf(StaffOrder, staffs[j].Role.String)
	})
	// Merge duplicate staffs with multiple roles
	merged_staffs := make([]sqlc.GetMediaStaffsRow, 0)
	for _, staff := range staffs {
		added_role := false
		for index, existing_staff := range merged_staffs {
			if existing_staff.ID == staff.ID {
				if existing_staff.Role.String == "" {
					merged_staffs[index].Role.String = staff.Role.String
				} else {
					merged_staffs[index].Role.String = existing_staff.Role.String + ", " + staff.Role.String
				}
				added_role = true
				break
			}
		}
		if !added_role {
			merged_staffs = append(merged_staffs, staff)
		}
	}
	// If there is more than 15 staffs, trim them to the most "important" only
	for i := 0; i < min(len(merged_staffs), 15); i++ {
		staff := merged_staffs[i]
		response.Staffs = append(response.Staffs, &pb.Staff{
			Id:        int32(staff.ID),
			Name:      staff.Name,
			Thumbnail: staff.Thumbnail.String,
			Role:      staff.Role.String,
		})
	}

	// Load comments
	comments, err := databases.DBs.SqlcQueries.GetMediaComments(ctx, int32(media.ID))
	if err != nil {
		return nil, err
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
