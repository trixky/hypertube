package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/trixky/hypertube/api-scrapper/postgres"
	pb "github.com/trixky/hypertube/api-scrapper/proto"
	"github.com/trixky/hypertube/api-scrapper/scrapper"
	st "github.com/trixky/hypertube/api-scrapper/sites"
	"github.com/trixky/hypertube/api-scrapper/sqlc"
	grpcMetadata "google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TorrentCreationResult struct {
	torrent sqlc.Torrent
	files   []sqlc.TorrentFile
}

var categories []string = []string{"movies", "shows"}

func makeNullInt32(value *int32) (null_int32 sql.NullInt32) {
	if value == nil {
		return
	}
	null_int32.Int32 = *value
	null_int32.Valid = true
	return
}

func makeNullString(value *string) (null_string sql.NullString) {
	if value == nil {
		return
	}
	null_string.String = *value
	null_string.Valid = true
	return
}

func torrenToSQL(torrent *pb.UnprocessedTorrent) sqlc.CreateTorrentParams {
	return sqlc.CreateTorrentParams{
		Name:            torrent.Name,
		Type:            torrent.Type.String(),
		FullUrl:         torrent.FullUrl,
		DescriptionHtml: makeNullString(torrent.DescriptionHtml),
		ImdbID:          makeNullString(torrent.ImdbId),
		Leech:           makeNullInt32(torrent.Leech),
		Magnet:          makeNullString(torrent.Magnet),
		Seed:            makeNullInt32(torrent.Seed),
		Size:            makeNullString(torrent.Size),
		TorrentUrl:      makeNullString(torrent.TorrentUrl),
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
	var imdb_id *string
	if torrent.ImdbID.Valid {
		imdb_id = &torrent.ImdbID.String
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
		ImdbId:          imdb_id,
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
	existing_torrent, err := postgres.DB.SqlcQueries.GetTorrentByURL(ctx, torrent.FullUrl)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	} else {
		err = nil
	}
	if existing_torrent.ID == 0 {
		fmt.Println("Inserting new torrent for", torrent.FullUrl)
		scrapper.ScrapeSingle(torrent)
		created_torrent, err_creation := postgres.DB.SqlcQueries.CreateTorrent(ctx, torrenToSQL(torrent))
		if err_creation != nil {
			err = err_creation
			return
		}
		creation_result.torrent = created_torrent
		time.Sleep(time.Second)
		created = true
	} else {
		fmt.Println("Updating torrent", existing_torrent.ID)
		creation_result.torrent = existing_torrent
		scrapped := false
		if !existing_torrent.TorrentUrl.Valid {
			scrapper.ScrapeSingle(torrent)
			// Update some informations if the torrent was not fully fetched originally
			postgres.DB.SqlcQueries.SetTorrentInformations(ctx, sqlc.SetTorrentInformationsParams{
				ID:              existing_torrent.ID,
				DescriptionHtml: makeNullString(torrent.DescriptionHtml),
				ImdbID:          makeNullString(torrent.ImdbId),
				Magnet:          makeNullString(torrent.Magnet),
				Size:            makeNullString(torrent.Size),
				TorrentUrl:      makeNullString(torrent.TorrentUrl),
			})
			scrapped = true
		}
		err = postgres.DB.SqlcQueries.SetTorrentPeers(ctx, sqlc.SetTorrentPeersParams{
			ID:    existing_torrent.ID,
			Seed:  makeNullInt32(torrent.Seed),
			Leech: makeNullInt32(torrent.Leech),
		})
		if err != nil {
			return
		}
		if scrapped {
			time.Sleep(time.Second)
		}
		created = false
	}

	// Create associated files
	if created {
		for _, file := range torrent.Files {
			created_file, err_file := postgres.DB.SqlcQueries.AddTorrentFile(ctx, sqlc.AddTorrentFileParams{
				TorrentID: int32(creation_result.torrent.ID),
				Name:      file.Name,
				Path:      makeNullString(file.Path),
				Size:      makeNullString(file.Size),
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
	log.Printf("Scrape All %v\n", request)

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

	return nil
}

func (s *ScrapperServer) IdentifyAll(request *pb.IdentifyRequest, out pb.ScrapperService_IdentifyAllServer) error {
	log.Printf("Identify All %v\n", request)
	return nil
}

func (s *ScrapperServer) ScrapeLatest(request *pb.ScrapeLatestRequest, out pb.ScrapperService_ScrapeLatestServer) error {
	ctx := context.Background()
	log.Printf("Scrap Latest %v\n", request)

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

	return nil
}

func (s *ScrapperServer) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
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

func (s *ScrapperServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	md, ok := grpcMetadata.FromIncomingContext(ctx)

	if !ok {
		log.Println("missing args")
		return nil, nil
	}

	search := md.Get("get")
	fmt.Println("get:", search)

	return &pb.GetResponse{
		Id:          1,
		ImdbId:      "tt3456",
		Name:        "Movie",
		Description: "Movie ?",
		Year:        2000,
		TorrentPublicInformations: &pb.TorrentPublicInformations{
			Name:  "Movie [1080p]",
			Seed:  &[]int32{42}[0],
			Leech: &[]int32{42}[0],
			Size:  &[]string{"123456789"}[0],
		},
		Staffs: []*pb.Staff{
			&pb.Staff{
				Id:        1,
				ImdbId:    "tt1234",
				Name:      "Writer",
				Role:      "Writer",
				Thumbnail: "jpg",
				Url:       "http",
			},
		},
		Relations: []*pb.Relation{
			&pb.Relation{
				Id:        2,
				ImdbId:    "tt2345",
				Name:      "Movie 2",
				Thumbnail: "jpg",
			},
		},
		Duration:   &[]string{"1h42"}[0],
		Thumbnail:  &[]string{"jpg"}[0],
		Background: &[]string{"jpg"}[0],
		Genres:     &[]string{"movie,action"}[0],
		Rating:     &[]string{"80"}[0],
	}, nil
}
