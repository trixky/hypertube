package internal

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/trixky/hypertube/api-scrapper/databases"
	pb "github.com/trixky/hypertube/api-scrapper/proto"
	st "github.com/trixky/hypertube/api-scrapper/sites"
	"github.com/trixky/hypertube/api-scrapper/sqlc"
	"github.com/trixky/hypertube/api-scrapper/utils"
	ut "github.com/trixky/hypertube/api-scrapper/utils"
	grpcMetadata "google.golang.org/grpc/metadata"
)

// Only torrent.ID and torrent.FullURL is required
func refreshTorrent(ctx context.Context, torrent *sqlc.Torrent) (response *pb.RefreshResponse, err error) {
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

			return &pb.RefreshResponse{
				TorrentId: uint32(torrent.ID),
				Seed:      torrent_update.Seed,
				Leech:     torrent_update.Leech,
			}, nil
		}
	}

	log.Println("No scrapper matched for", torrent.FullUrl)

	return &pb.RefreshResponse{}, nil

}

func (s *ScrapperServer) Refresh(in *pb.RefreshRequest, out pb.ScrapperService_RefreshServer) error {
	ctx := out.Context()
	md, ok := grpcMetadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("missing args")
		return nil
	}

	refresh := md.Get("refresh")
	log.Println("refresh:", refresh)

	// Find the media
	media, err := databases.DBs.SqlcQueries.GetMediaByID(ctx, int64(in.MediaId))
	if err != nil {
		return err
	}

	// Get the torrents associated with the Media
	media_id := int32(media.ID)
	torrents, err := databases.DBs.SqlcQueries.GetMediaTorrents(ctx, ut.MakeNullInt32(&media_id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		} else {
			return err
		}
	}

	// Call ScrapperService.Refresh on each torrents and stream the response
	// Only update a torrent if it has been updated more than 30min ago
	// Always send the seed and leech count even if it wasn't updated
	scrapped_last := false
	for _, torrent := range torrents {
		if time.Now().Sub(torrent.LastUpdate.Time).Minutes() >= 30 {
			if scrapped_last {
				time.Sleep(time.Duration(1) * time.Second)
			}
			response, err := refreshTorrent(ctx, &sqlc.Torrent{
				ID:      torrent.ID.Int64,
				FullUrl: torrent.FullUrl.String,
			})
			if err != nil {
				return err
			}
			out.Send(response)
			scrapped_last = true
		} else {
			out.Send(&pb.RefreshResponse{
				TorrentId: uint32(torrent.ID.Int64),
				Seed:      torrent.Seed.Int32,
				Leech:     torrent.Leech.Int32,
			})
			scrapped_last = false
		}
	}

	return nil
}
