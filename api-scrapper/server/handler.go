package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/trixky/hypertube/api-scrapper/postgres"
	pb "github.com/trixky/hypertube/api-scrapper/proto"
	"github.com/trixky/hypertube/api-scrapper/scrapper"
	st "github.com/trixky/hypertube/api-scrapper/sites"
	"github.com/trixky/hypertube/api-scrapper/sqlc"
	grpcMetadata "google.golang.org/grpc/metadata"
)

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

func updateOrCreateTorrent(ctx context.Context, scrapper *scrapper.Scrapper, torrent *pb.UnprocessedTorrent) (created bool, created_torrent sqlc.Torrent, err error) {
	existing_torrent, err := postgres.DB.SqlcQueries.GetTorrentByURL(ctx, torrent.FullUrl)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	}
	if existing_torrent.ID == 0 {
		fmt.Println("Inserting new torrent for", torrent.FullUrl)
		scrapper.ScrapeSingle(torrent)
		created_torrent, err = postgres.DB.SqlcQueries.CreateTorrent(ctx, torrenToSQL(torrent))
		if err != nil {
			return
		}
		time.Sleep(time.Second)
		created = true
	} else {
		fmt.Println("Updating torrent", existing_torrent.ID)
		created_torrent = existing_torrent
		scrapped := false
		if !existing_torrent.TorrentUrl.Valid {
			scrapper.ScrapeSingle(torrent)
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
	if created && len(torrent.Files) > 0 {
		for _, file := range torrent.Files {
			err = postgres.DB.SqlcQueries.AddTorrentFile(ctx, sqlc.AddTorrentFileParams{
				TorrentID: int32(created_torrent.ID),
				Name:      file.Name,
				Path:      makeNullString(file.Path),
				Size:      makeNullString(file.Size),
			})
			if err != nil {
				return
			}
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
				page_result, err := scrapper.ScrapeList(category, page)
				new_torrents := make([]*sqlc.Torrent, 0, 30)

				// Update each existing torrents or add them to the database
				for _, torrent := range page_result.Torrents {
					created, created_torrent, err := updateOrCreateTorrent(ctx, &scrapper, torrent)
					if err != nil {
						return err
					}
					if created {
						new_torrents = append(new_torrents, &created_torrent)
					}
				}

				// Send the new torrents on the stream
				if err == nil {
					if err := out.Send(&pb.ScrapeResponse{
						MsDuration: 0,
						// TODO Convert to pb.Torrent
						Torrents: []*pb.Torrent{}, /* new_torrents */
					}); err != nil {
						return err
					}
				} else {
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
				page_result, err := scrapper.ScrapeList(category, page)
				new_torrents := make([]*sqlc.Torrent, 0, 30)

				// Update each existing torrents or add them to the database
				for _, torrent := range page_result.Torrents {
					created, created_torrent, err := updateOrCreateTorrent(ctx, &scrapper, torrent)
					if err != nil {
						return err
					}
					if created {
						new_torrents = append(new_torrents, &created_torrent)
					} else {
						has_existing = true
					}
				}

				// Send the new torrents on the stream
				if err == nil {
					if err := out.Send(&pb.ScrapeResponse{
						MsDuration: 0,
						// TODO Convert to pb.Torrent
						Torrents: []*pb.Torrent{}, /* new_torrents */
					}); err != nil {
						return err
					}
				} else {
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
