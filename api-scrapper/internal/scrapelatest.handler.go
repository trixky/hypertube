package internal

import (
	"context"
	"log"
	"time"

	pb "github.com/trixky/hypertube/api-scrapper/proto"
	st "github.com/trixky/hypertube/api-scrapper/sites"
)

func (s *ScrapperServer) ScrapeLatest(request *pb.ScrapeLatestRequest, out pb.ScrapperService_ScrapeLatestServer) error {
	ctx := context.Background()
	log.Println("Scrape Latest", request)

	for _, scrapper := range st.Scrappers {
		for _, category := range Categories {
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
					created, created_torrent, err := UpdateOrCreateTorrent(ctx, &scrapper, torrent)
					if err != nil {
						return err
					}
					if created {
						converted_torrent := TorrentToProto(&created_torrent)
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
