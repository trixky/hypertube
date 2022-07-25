package internal

import (
	"context"
	"log"

	"github.com/trixky/hypertube/.shared/utils"
	pb "github.com/trixky/hypertube/api-scrapper/proto"
	"github.com/trixky/hypertube/api-scrapper/queries"
	st "github.com/trixky/hypertube/api-scrapper/sites"
	"github.com/trixky/hypertube/api-scrapper/sqlc"
)

func (s *ScrapperServer) RefreshTorrent(ctx context.Context, in *pb.RefreshTorrentRequest) (*pb.RefreshTorrentResponse, error) {
	// Check if the torrent exists
	torrent, err := queries.SqlcQueries.GetTorrentByID(ctx, int64(in.TorrentId))
	if err != nil {
		return nil, err
	}

	// Scrape with scrapeSingle on the corresponding scrapper
	for _, scrapper := range st.Scrappers {
		if scrapper.CanUpdate(torrent.FullUrl) {
			torrent_update := pb.UnprocessedTorrent{
				FullUrl: torrent.FullUrl,
			}
			scrapper.ScrapeSingle(&torrent_update)

			// Update the seeds and leeches
			err = queries.SqlcQueries.SetTorrentPeers(ctx, sqlc.SetTorrentPeersParams{
				ID:    torrent.ID,
				Seed:  utils.MakeNullInt32(&torrent_update.Seed),
				Leech: utils.MakeNullInt32(&torrent_update.Leech),
			})
			if err != nil {
				return nil, err
			}

			return &pb.RefreshTorrentResponse{
				TorrentId: uint32(torrent.ID),
				Seed:      torrent_update.Seed,
				Leech:     torrent_update.Leech,
			}, nil
		}
	}

	log.Println("No scrapper matched for", torrent.FullUrl)

	return &pb.RefreshTorrentResponse{}, nil
}
