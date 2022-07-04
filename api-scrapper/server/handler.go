package server

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/trixky/hypertube/api-scrapper/postgres"
	pb "github.com/trixky/hypertube/api-scrapper/proto"
	"github.com/trixky/hypertube/api-scrapper/scrapper"
	st "github.com/trixky/hypertube/api-scrapper/sites"
	"github.com/trixky/hypertube/api-scrapper/sqlc"
	ut "github.com/trixky/hypertube/api-scrapper/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TorrentCreationResult struct {
	torrent sqlc.Torrent
	files   []sqlc.TorrentFile
}

var categories []string = []string{
	"movies",
	// "shows",
}

func torrenToSQL(torrent *pb.UnprocessedTorrent) sqlc.CreateTorrentParams {
	return sqlc.CreateTorrentParams{
		Name:            torrent.Name,
		Type:            torrent.Type.String(),
		FullUrl:         torrent.FullUrl,
		DescriptionHtml: ut.MakeNullString(torrent.DescriptionHtml),
		Leech:           ut.MakeNullInt32(torrent.Leech),
		Magnet:          ut.MakeNullString(torrent.Magnet),
		Seed:            ut.MakeNullInt32(torrent.Seed),
		Size:            ut.MakeNullString(torrent.Size),
		TorrentUrl:      ut.MakeNullString(torrent.TorrentUrl),
	}
}

func torrentToProto(creation_result *TorrentCreationResult) (converted_torrent pb.Torrent) {
	torrent := creation_result.torrent
	var torrent_type pb.MediaCategory
	if torrent.Type == categories[0] {
		torrent_type = pb.MediaCategory_CATEGORY_MOVIE
	} else {
		torrent_type = pb.MediaCategory_CATEGORY_SERIE
	}
	var seed *int32
	if torrent.Seed.Valid {
		seed = &torrent.Seed.Int32
	}
	var leech *int32
	if torrent.Leech.Valid {
		leech = &torrent.Leech.Int32
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
		Seed:            seed,
		Leech:           leech,
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

func updateOrCreateTorrent(ctx context.Context, scrapper *scrapper.Scrapper, torrent *pb.UnprocessedTorrent) (created bool, creation_result TorrentCreationResult, err error) {
	db_torrent, err := postgres.DB.SqlcQueries.GetTorrentByURL(ctx, torrent.FullUrl)
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

		created_torrent, err_creation := postgres.DB.SqlcQueries.CreateTorrent(ctx, torrenToSQL(torrent))
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
		if !db_torrent.TorrentUrl.Valid {
			scrapper.ScrapeSingle(torrent)

			postgres.DB.SqlcQueries.SetTorrentInformations(ctx, sqlc.SetTorrentInformationsParams{
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
		err = postgres.DB.SqlcQueries.SetTorrentPeers(ctx, sqlc.SetTorrentPeersParams{
			ID:    db_torrent.ID,
			Seed:  ut.MakeNullInt32(torrent.Seed),
			Leech: ut.MakeNullInt32(torrent.Leech),
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
			err = postgres.DB.SqlcQueries.AddTorrentMediaId(ctx, sqlc.AddTorrentMediaIdParams{
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
				err = postgres.DB.SqlcQueries.AddTorrentMediaId(ctx, sqlc.AddTorrentMediaIdParams{
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
			created_file, err_file := postgres.DB.SqlcQueries.AddTorrentFile(ctx, sqlc.AddTorrentFileParams{
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

func (s *ScrapperServer) ScrapeAll(request *pb.ScrapeRequest, out pb.ScrapperService_ScrapeAllServer) error {
	ctx := context.Background()
	log.Println("Scrape All", request)

	for _, scrapper := range st.Scrappers {
		for _, category := range categories {
			var page uint32 = 1
			for {
				start := time.Now()
				page_result, err := scrapper.ScrapeList(category, page)
				if err != nil {
					return err
				}
				new_torrents := make([]*pb.Torrent, 0, 30)

				// Update each existing torrents or add them to the database
				for _, torrent := range page_result.Torrents {
					created, created_torrent, err := updateOrCreateTorrent(ctx, &scrapper, torrent)
					if err != nil {
						return err
					}
					if created {
						converted_torrent := torrentToProto(&created_torrent)
						new_torrents = append(new_torrents, &converted_torrent)
					}
				}

				// Send the new torrents on the stream
				if err := out.Send(&pb.ScrapeResponse{
					MsDuration: uint32(time.Since(start).Milliseconds()),
					Torrents:   new_torrents,
				}); err != nil {
					return err
				}

				// Update NextPage to loop or complete the job
				page = page_result.NextPage
				if page == 0 {
					break
				}
				time.Sleep(time.Second)
			}
		}
	}

	log.Println("Done ScrapeAll")

	return nil
}

func (s *ScrapperServer) ScrapeLatest(request *pb.ScrapeLatestRequest, out pb.ScrapperService_ScrapeLatestServer) error {
	ctx := context.Background()
	log.Println("Scrape Latest", request)

	for _, scrapper := range st.Scrappers {
		for _, category := range categories {
			var page uint32 = 1
			has_existing := false
			for {
				start := time.Now()
				page_result, err := scrapper.ScrapeList(category, page)
				if err != nil {
					return err
				}

				// Update each existing torrents or add them to the database
				new_torrents := make([]*pb.Torrent, 0, 30)
				for _, torrent := range page_result.Torrents {
					created, created_torrent, err := updateOrCreateTorrent(ctx, &scrapper, torrent)
					if err != nil {
						return err
					}
					if created {
						converted_torrent := torrentToProto(&created_torrent)
						new_torrents = append(new_torrents, &converted_torrent)
					} else {
						has_existing = true
					}
				}

				// Send the new torrents on the stream
				if err := out.Send(&pb.ScrapeResponse{
					MsDuration: uint32(time.Since(start).Milliseconds()),
					Torrents:   new_torrents,
				}); err != nil {
					return err
				}

				// Update NextPage to loop or complete the job
				page = page_result.NextPage
				if page == 0 || has_existing {
					break
				}
				time.Sleep(time.Second)
			}
		}
	}

	log.Println("Done ScrapeLatest")

	return nil
}
