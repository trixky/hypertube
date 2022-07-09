package internal

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/trixky/hypertube/api-scrapper/databases"
	pb "github.com/trixky/hypertube/api-scrapper/proto"
	st "github.com/trixky/hypertube/api-scrapper/sites"
	"github.com/trixky/hypertube/api-scrapper/sqlc"
	"github.com/trixky/hypertube/api-scrapper/utils"
	grpcMetadata "google.golang.org/grpc/metadata"
)

func (s *ScrapperServer) Refresh(ctx context.Context, in *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	md, ok := grpcMetadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("missing args")
		return nil, nil
	}

	refresh := md.Get("refresh")
	log.Println("refresh:", refresh)

	// Check if the torrent exists and get it's URL
	torrent, err := databases.DBs.SqlcQueries.GetTorrentByID(ctx, int64(in.Id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &pb.RefreshResponse{}, nil
		} else {
			return &pb.RefreshResponse{}, err
		}
	}

	// Scrape with scrapeSingle on the corresponding scrapper
	for _, scrapper := range st.Scrappers {
		if scrapper.CanUpdate(torrent.FullUrl) {
			torrent_update := pb.UnprocessedTorrent{
				FullUrl: torrent.FullUrl,
			}
			scrapper.ScrapeSingle(&torrent_update)

			// Update the seeds and leeches
			err = databases.DBs.SqlcQueries.SetTorrentPeers(ctx, sqlc.SetTorrentPeersParams{
				ID:    torrent.ID,
				Seed:  utils.MakeNullInt32(&torrent_update.Seed),
				Leech: utils.MakeNullInt32(&torrent_update.Leech),
			})
			if err != nil {
				return nil, err
			}

			log.Println("Done Refresh")
			return &pb.RefreshResponse{
				Id:    uint32(torrent.ID),
				Seed:  torrent_update.Seed,
				Leech: torrent_update.Leech,
			}, nil
		}
	}

	log.Println("No scrapper matched for", torrent.FullUrl)

	return &pb.RefreshResponse{}, nil

}
