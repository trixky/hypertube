package internal

import (
	"database/sql"
	"errors"
	"log"
	"time"

	sutils "github.com/trixky/hypertube/.shared/utils"
	"github.com/trixky/hypertube/api-media/external"
	pb "github.com/trixky/hypertube/api-media/proto"
	"github.com/trixky/hypertube/api-media/queries"
	"github.com/trixky/hypertube/api-media/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *MediaServer) Refresh(in *pb.RefreshMediaRequest, out pb.MediaService_RefreshServer) error {
	if _, err := utils.RequireLogin(out.Context()); err != nil {
		return err
	}

	ctx := out.Context()

	// Find the media
	media, err := queries.SqlcQueries.GetMediaByID(ctx, int64(in.MediaId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return status.Errorf(codes.NotFound, "no media with this id")
		} else {
			log.Println(err)
			return status.Errorf(codes.Internal, "failed to get media")
		}
	}

	// Get the torrents associated with the Media
	media_id := int32(media.ID)
	torrents, err := queries.SqlcQueries.GetMediaTorrents(ctx, sutils.MakeNullInt32(&media_id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		} else {
			log.Println(err)
			return status.Errorf(codes.Internal, "failed to get media torrents")
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
				log.Println(err)
				return status.Errorf(codes.Internal, "failed to scrape torrent")
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
