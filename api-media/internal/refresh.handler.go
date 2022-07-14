package internal

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/trixky/hypertube/api-media/databases"
	"github.com/trixky/hypertube/api-media/external"
	pb "github.com/trixky/hypertube/api-media/proto"
	"github.com/trixky/hypertube/api-media/utils"
	ut "github.com/trixky/hypertube/api-media/utils"
	grpcMetadata "google.golang.org/grpc/metadata"
)

func (s *MediaServer) Refresh(in *pb.RefreshMediaRequest, out pb.MediaService_RefreshServer) error {
	if _, err := utils.RequireLogin(out.Context()); err != nil {
		return err
	}

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

			// Ask api-scrapper to do the actual scrapping
			response, err := external.ApiScrapper.RefreshTorrent(ctx, &pb.RefreshTorrentRequest{
				TorrentId: uint64(torrent.ID.Int64),
			})
			if err != nil {
				return err
			}

			out.Send(&pb.RefreshMediaResponse{
				TorrentId: response.TorrentId,
				Seed:      response.Seed,
				Leech:     response.Leech,
			})
			scrapped_last = true
		} else {
			out.Send(&pb.RefreshMediaResponse{
				TorrentId: uint32(torrent.ID.Int64),
				Seed:      torrent.Seed.Int32,
				Leech:     torrent.Leech.Int32,
			})
			scrapped_last = false
		}
	}

	return nil
}
