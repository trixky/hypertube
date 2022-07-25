package internal

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	ut "github.com/trixky/hypertube/.shared/utils"
	pb "github.com/trixky/hypertube/api-scrapper/proto"
	"github.com/trixky/hypertube/api-scrapper/queries"
	"github.com/trixky/hypertube/api-scrapper/scrapper"
	st "github.com/trixky/hypertube/api-scrapper/sites"
	"github.com/trixky/hypertube/api-scrapper/sqlc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TorrentCreationResult struct {
	torrent sqlc.Torrent
	files   []sqlc.TorrentFile
}

var Categories []string = []string{
	"movies",
	// "shows",
}

var ErrorBackOff []int32 = []int32{1, 3, 5}

func TorrenToSQL(torrent *pb.UnprocessedTorrent) sqlc.CreateTorrentParams {
	return sqlc.CreateTorrentParams{
		Name:            torrent.Name,
		Type:            torrent.Type.String(),
		FullUrl:         torrent.FullUrl,
		DescriptionHtml: ut.MakeNullString(torrent.DescriptionHtml),
		Leech:           ut.MakeNullInt32(&torrent.Leech),
		Magnet:          ut.MakeNullString(torrent.Magnet),
		Seed:            ut.MakeNullInt32(&torrent.Seed),
		Size:            ut.MakeNullString(torrent.Size),
		TorrentUrl:      ut.MakeNullString(torrent.TorrentUrl),
	}
}

func TorrentToProto(creation_result *TorrentCreationResult) (converted_torrent pb.Torrent) {
	torrent := creation_result.torrent
	var torrent_type pb.MediaCategory
	if torrent.Type == Categories[0] {
		torrent_type = pb.MediaCategory_CATEGORY_MOVIE
	} else {
		torrent_type = pb.MediaCategory_CATEGORY_SERIE
	}
	var size *string
	if torrent.Size.Valid {
		size = &torrent.Size.String
	}
	var upload_time *timestamppb.Timestamp
	if torrent.UploadTime.Valid {
		upload_time, _ = ptypes.TimestampProto(torrent.UploadTime.Time)
	}
	var description *string
	if torrent.DescriptionHtml.Valid {
		description = &torrent.DescriptionHtml.String
	}
	var magnet *string
	if torrent.Magnet.Valid {
		magnet = &torrent.Magnet.String
	}

	// Convert the DB Torrent to a protobuf message
	converted_torrent = pb.Torrent{
		Id:              uint32(torrent.ID),
		Name:            torrent.Name,
		Type:            torrent_type,
		FullUrl:         torrent.FullUrl,
		TorrentUrl:      torrent.FullUrl,
		Seed:            torrent.Seed.Int32,
		Leech:           torrent.Leech.Int32,
		Size:            size,
		UploadTime:      upload_time,
		DescriptionHtml: description,
		Magnet:          magnet,
		Files:           make([]*pb.TorrentFile, 0, len(creation_result.files)),
	}

	// Add files
	for _, torrent_file := range creation_result.files {
		var path *string
		if torrent_file.Path.Valid {
			path = &torrent_file.Path.String
		}
		var size *string
		if torrent_file.Size.Valid {
			size = &torrent_file.Size.String
		}
		converted_torrent.Files = append(converted_torrent.Files, &pb.TorrentFile{
			Name: torrent_file.Name,
			Path: path,
			Size: size,
		})
	}

	return
}

func UpdateOrCreateTorrent(ctx context.Context, scrapper *scrapper.Scrapper, torrent *pb.UnprocessedTorrent) (created bool, creation_result TorrentCreationResult, err error) {
	db_torrent, err := queries.SqlcQueries.GetTorrentByURL(ctx, torrent.FullUrl)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return
		} else {
			err = nil
		}
	}
	var imdb_id string

	if db_torrent.ID == 0 {
		log.Println("Inserting new torrent for", torrent.FullUrl)
		scrapper.ScrapeSingle(torrent)
		if torrent.ImdbId != nil {
			imdb_id = *torrent.ImdbId
		}

		created_torrent, err_creation := queries.SqlcQueries.CreateTorrent(ctx, TorrenToSQL(torrent))
		if err_creation != nil {
			err = err_creation
			return
		}
		creation_result.torrent = created_torrent
		db_torrent = created_torrent

		time.Sleep(time.Second)
		created = true
	} else {
		log.Println("Updating torrent", db_torrent.ID)
		creation_result.torrent = db_torrent
		scrapped := false

		// Update some informations if the torrent was not fully fetched originally
		if !db_torrent.TorrentUrl.Valid && !db_torrent.Magnet.Valid {
			scrapper.ScrapeSingle(torrent)

			queries.SqlcQueries.SetTorrentInformations(ctx, sqlc.SetTorrentInformationsParams{
				ID:              db_torrent.ID,
				DescriptionHtml: ut.MakeNullString(torrent.DescriptionHtml),
				Magnet:          ut.MakeNullString(torrent.Magnet),
				Size:            ut.MakeNullString(torrent.Size),
				TorrentUrl:      ut.MakeNullString(torrent.TorrentUrl),
			})

			if torrent.ImdbId != nil && *torrent.ImdbId != "" {
				imdb_id = *torrent.ImdbId
			}
			scrapped = true
		}

		// Always update peers
		err = queries.SqlcQueries.SetTorrentPeers(ctx, sqlc.SetTorrentPeersParams{
			ID:    db_torrent.ID,
			Seed:  ut.MakeNullInt32(&torrent.Seed),
			Leech: ut.MakeNullInt32(&torrent.Leech),
		})
		if err != nil {
			return
		}

		if scrapped {
			time.Sleep(time.Second)
		}
		created = false
	}

	// Find IMDB informations or associated the media with an ID
	if imdb_id != "" && !db_torrent.MediaID.Valid {
		media_id, err_find := st.InsertOrGetMedia(imdb_id)
		if err_find != nil {
			err = err_find
			return
		}
		if media_id > 0 {
			err = queries.SqlcQueries.AddTorrentMediaId(ctx, sqlc.AddTorrentMediaIdParams{
				ID:      db_torrent.ID,
				MediaID: ut.MakeNullInt32(&media_id),
			})
			if err != nil {
				return
			}
			db_torrent.MediaID.Int32 = media_id
			db_torrent.MediaID.Valid = true
		}
	}

	if created {
		// Try to find an IMDB
		if imdb_id == "" && !db_torrent.MediaID.Valid {
			media_id, err_find := st.TryInsertOrGetMedia(db_torrent.Name)
			if err_find != nil {
				err = err_find
				return
			}
			if media_id > 0 {
				err = queries.SqlcQueries.AddTorrentMediaId(ctx, sqlc.AddTorrentMediaIdParams{
					ID:      db_torrent.ID,
					MediaID: ut.MakeNullInt32(&media_id),
				})
				if err != nil {
					return
				}
				db_torrent.MediaID.Int32 = media_id
				db_torrent.MediaID.Valid = true
			}
		}

		// Create associated files
		for _, file := range torrent.Files {
			created_file, err_file := queries.SqlcQueries.AddTorrentFile(ctx, sqlc.AddTorrentFileParams{
				TorrentID: int32(creation_result.torrent.ID),
				Name:      file.Name,
				Path:      ut.MakeNullString(file.Path),
				Size:      ut.MakeNullString(file.Size),
			})
			if err_file != nil {
				err = err_file
				return
			}
			creation_result.files = append(creation_result.files, created_file)
		}
	}

	return
}
